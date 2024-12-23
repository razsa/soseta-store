package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

// API response structures
type SearchProductsResponse struct {
	Items       []*core.Record `json:"items"`
	TotalItems  int            `json:"totalItems"`
	TotalPages  int            `json:"totalPages"`
	CurrentPage int            `json:"currentPage"`
}

type CartItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type CreateOrderRequest struct {
	CartItems       []CartItem `json:"cartItems"`
	ShippingAddress string     `json:"shippingAddress"`
}

func main() {
	app := pocketbase.New()

	// Add custom routes
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		// serves static files from the provided public dir (if exists)
		se.Router.GET("/{path...}", apis.Static(os.DirFS("./pb_public"), false))

		// Product search endpoint with filtering and pagination
		se.Router.GET("/api/products/search", func(e *core.RequestEvent) error {
			collection, err := app.FindCollectionByNameOrId("products")
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Collection not found"})
			}

			query := e.Request.URL.Query().Get("query")
			category := e.Request.URL.Query().Get("category")
			page := e.Request.URL.Query().Get("page")
			perPage := e.Request.URL.Query().Get("perPage")

			// Default pagination values
			pageNum := 1
			perPageNum := 20
			if page != "" {
				pageNum, _ = strconv.Atoi(page)
			}
			if perPage != "" {
				perPageNum, _ = strconv.Atoi(perPage)
			}
			offset := (pageNum - 1) * perPageNum

			filter := ""
			if query != "" {
				filter = "name ~ '" + query + "' || description ~ '" + query + "'"
			}
			if category != "" {
				if filter != "" {
					filter += " && "
				}
				filter += "category = '" + category + "'"
			}

			// Get total count first
			totalCount, err := app.FindRecordsByFilter(collection.Id, filter, "", 0, 0)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			// Then get paginated results with expanded fields
			records, err := app.FindRecordsByFilter(collection.Id, filter, "id", perPageNum, offset)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			return e.JSON(http.StatusOK, SearchProductsResponse{
				Items:       records,
				TotalItems:  len(totalCount),
				TotalPages:  (len(totalCount) + perPageNum - 1) / perPageNum,
				CurrentPage: pageNum,
			})
		})

		// Create order endpoint
		se.Router.POST("/api/orders", func(e *core.RequestEvent) error {
			// Require authentication
			authRecord, _ := e.Get("authRecord").(*core.Record)
			if authRecord == nil {
				return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			var req CreateOrderRequest
			if err := json.NewDecoder(e.Request.Body).Decode(&req); err != nil {
				return e.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
			}

			// Validate cart items
			if len(req.CartItems) == 0 {
				return e.JSON(http.StatusBadRequest, map[string]string{"error": "Cart is empty"})
			}

			// Start transaction
			err := app.RunInTransaction(func(txApp core.App) error {
				// Create order record
				collection, err := app.FindCollectionByNameOrId("orders")
				if err != nil {
					return err
				}
				order := core.NewRecord(collection)
				order.Set("user", authRecord.Id)
				order.Set("status", "pending")
				order.Set("shipping_address", req.ShippingAddress)
				order.Set("payment_status", "pending")

				totalAmount := 0.0

				// Process each cart item
				for _, item := range req.CartItems {
					product, err := app.FindRecordById("products", item.ProductId)
					if err != nil {
						return err
					}

					// Check stock
					currentStock := product.GetInt("stock")
					if currentStock < item.Quantity {
						return echo.NewHTTPError(http.StatusBadRequest, "Insufficient stock for product: "+product.GetString("name"))
					}

					// Update stock
					product.Set("stock", currentStock-item.Quantity)
					if err := txApp.Save(product); err != nil {
						return err
					}

					// Calculate item total
					price := product.GetFloat("price")
					totalAmount += price * float64(item.Quantity)

					// Create order item
					orderItemsCollection, err := app.FindCollectionByNameOrId("order_items")
					if err != nil {
						return err
					}
					orderItem := core.NewRecord(orderItemsCollection)
					orderItem.Set("order", order.Id)
					orderItem.Set("product", item.ProductId)
					orderItem.Set("quantity", item.Quantity)
					orderItem.Set("price_at_time", price)

					if err := txApp.Save(orderItem); err != nil {
						return err
					}
				}

				order.Set("total_amount", totalAmount)
				if err := txApp.Save(order); err != nil {
					return err
				}

				return nil
			})

			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			return e.JSON(http.StatusOK, map[string]string{"message": "Order created successfully"})
		})

		// Get user orders endpoint
		se.Router.GET("/api/user/orders", func(e *core.RequestEvent) error {
			authRecord, _ := e.Get("authRecord").(*core.Record)
			if authRecord == nil {
				return e.JSON(http.StatusUnauthorized, map[string]string{"error": "Authentication required"})
			}

			collection, err := app.FindCollectionByNameOrId("orders")
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": "Collection not found"})
			}

			filter := "user = '" + authRecord.Id + "'"
			records, err := app.FindRecordsByFilter(collection.Id, filter, "", 20, 0)
			if err != nil {
				return e.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}

			return e.JSON(http.StatusOK, records)
		})

		return se.Next()
	})

	// Create collections if they don't exist
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		collections := []string{
			// Categories Collection
			`{
				"name": "categories",
				"type": "base",
				"schema": [
					{
						"name": "name",
						"type": "text",
						"required": true,
						"unique": true
					},
					{
						"name": "description",
						"type": "text"
					},
					{
						"name": "image",
						"type": "file",
						"options": {
							"maxSelect": 1,
							"mimeTypes": ["image/jpg", "image/jpeg", "image/png", "image/gif"],
							"maxSize": 5242880
						}
					}
				]
			}`,
			// Products Collection
			`{
				"name": "products",
				"type": "base",
				"schema": [
					{
						"name": "name",
						"type": "text",
						"required": true
					},
					{
						"name": "description",
						"type": "text"
					},
					{
						"name": "price",
						"type": "number",
						"required": true,
						"min": 0
					},
					{
						"name": "stock",
						"type": "number",
						"required": true,
						"min": 0
					},
					{
						"name": "category",
						"type": "relation",
						"required": true,
						"options": {
							"collectionId": "categories",
							"cascadeDelete": false
						}
					},
					{
						"name": "images",
						"type": "file",
						"options": {
							"maxSelect": 5,
							"mimeTypes": ["image/jpg", "image/jpeg", "image/png", "image/gif"],
							"maxSize": 5242880
						}
					}
				]
			}`,
			// Orders Collection
			`{
				"name": "orders",
				"type": "base",
				"schema": [
					{
						"name": "user",
						"type": "relation",
						"required": true,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false
						}
					},
					{
						"name": "status",
						"type": "select",
						"required": true,
						"options": {
							"values": ["pending", "processing", "shipped", "delivered", "cancelled"]
						}
					},
					{
						"name": "total_amount",
						"type": "number",
						"required": true,
						"min": 0
					},
					{
						"name": "shipping_address",
						"type": "text",
						"required": true
					},
					{
						"name": "payment_status",
						"type": "select",
						"required": true,
						"options": {
							"values": ["pending", "paid", "failed", "refunded"]
						}
					}
				]
			}`,
			// Order Items Collection
			`{
				"name": "order_items",
				"type": "base",
				"schema": [
					{
						"name": "order",
						"type": "relation",
						"required": true,
						"options": {
							"collectionId": "orders",
							"cascadeDelete": true
						}
					},
					{
						"name": "product",
						"type": "relation",
						"required": true,
						"options": {
							"collectionId": "products",
							"cascadeDelete": false
						}
					},
					{
						"name": "quantity",
						"type": "number",
						"required": true,
						"min": 1
					},
					{
						"name": "price_at_time",
						"type": "number",
						"required": true,
						"min": 0
					}
				]
			}`,
		}

		for _, collectionConfig := range collections {
			collection := &core.Collection{}
			if err := collection.UnmarshalJSON([]byte(collectionConfig)); err != nil {
				return err
			}

			// Skip if collection already exists
			if _, err := app.FindCollectionByNameOrId(collection.Name); err == nil {
				continue
			}

			if err := app.Save(collection); err != nil {
				return err
			}
		}

		err := importProducts(app)
		if err != nil {
			log.Printf("Error importing products: %v", err)
		}

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

// func seedExampleData(app *pocketbase.PocketBase) error {
// 	// Example categories
// 	categories := []map[string]interface{}{
// 		{"name": "Electronics", "description": "Electronic devices and accessories"},
// 		{"name": "Clothing", "description": "Fashion and apparel"},
// 		{"name": "Books", "description": "Books and literature"},
// 		{"name": "Home & Garden", "description": "Home decor and gardening items"},
// 	}

// 	categoryIds := make(map[string]string)
// 	for _, cat := range categories {
// 		collection, err := app.FindCollectionByNameOrId("categories")
// 		if err != nil {
// 			return fmt.Errorf("failed to find categories collection: %v", err)
// 		}

// 		record := core.NewRecord(collection)
// 		record.Set("name", cat["name"])
// 		record.Set("description", cat["description"])

// 		if err := app.Save(record); err != nil {
// 			return fmt.Errorf("failed to create category: %v", err)
// 		}

// 		categoryIds[cat["name"].(string)] = record.Id
// 	}

// 	// Seed Products
// 	products := []map[string]interface{}{
// 		{
// 			"name":        "Smartphone X",
// 			"description": "Latest smartphone with advanced features",
// 			"price":       699.99,
// 			"stock":       50,
// 			"category":    categoryIds["Electronics"],
// 		},
// 		{
// 			"name":        "Laptop Pro",
// 			"description": "High-performance laptop for professionals",
// 			"price":       1299.99,
// 			"stock":       30,
// 			"category":    categoryIds["Electronics"],
// 		},
// 		{
// 			"name":        "Classic T-Shirt",
// 			"description": "Comfortable cotton t-shirt",
// 			"price":       24.99,
// 			"stock":       100,
// 			"category":    categoryIds["Clothing"],
// 		},
// 		{
// 			"name":        "Designer Jeans",
// 			"description": "Premium quality denim jeans",
// 			"price":       89.99,
// 			"stock":       75,
// 			"category":    categoryIds["Clothing"],
// 		},
// 		{
// 			"name":        "Programming Guide",
// 			"description": "Comprehensive programming book",
// 			"price":       49.99,
// 			"stock":       45,
// 			"category":    categoryIds["Books"],
// 		},
// 		{
// 			"name":        "Garden Tools Set",
// 			"description": "Complete set of essential garden tools",
// 			"price":       129.99,
// 			"stock":       25,
// 			"category":    categoryIds["Home & Garden"],
// 		},
// 	}

// 	// Create products
// 	for _, prod := range products {
// 		collection, err := app.FindCollectionByNameOrId("products")
// 		if err != nil {
// 			return fmt.Errorf("failed to find products collection: %v", err)
// 		}

// 		record := core.NewRecord(collection)
// 		record.Set("name", prod["name"])
// 		record.Set("description", prod["description"])
// 		record.Set("price", prod["price"])
// 		record.Set("stock", prod["stock"])
// 		record.Set("category", prod["category"])

// 		if err := app.Save(record); err != nil {
// 			return fmt.Errorf("failed to create product: %v", err)
// 		}
// 	}

// 	log.Println("Seeding completed successfully!")
// 	return nil
// }
