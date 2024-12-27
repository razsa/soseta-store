package collection

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// CreateCollections creates the necessary collections if they don't already exist.
func CreateCollections(app core.App) error {
	// Create categories collection
	if err := createCategoriesCollection(app); err != nil {
		return err
	}

	// Create products collection
	if err := createProductsCollection(app); err != nil {
		return err
	}

	// Create orders collection
	if err := createOrdersCollection(app); err != nil {
		return err
	}

	// Create order_items collection
	if err := createOrderItemsCollection(app); err != nil {
		return err
	}

	slog.Info("All collections created successfully!")
	return nil
}

// createCategoriesCollection creates the "categories" collection.
func createCategoriesCollection(app core.App) error {
	// Check if the collection already exists
	existingCollection, err := app.FindCollectionByNameOrId("categories")
	if err == nil && existingCollection != nil {
		slog.Info("Collection 'categories' already exists")
		return nil
	}

	collection := core.NewBaseCollection("categories")

	// Add fields
	collection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	collection.Fields.Add(&core.TextField{
		Name: "description",
	})
	collection.Fields.Add(&core.FileField{
		Name:      "image",
		MaxSelect: 1,
		MimeTypes: []string{"image/jpg", "image/jpeg", "image/png", "image/gif"},
		MaxSize:   5242880,
	})

	// Save the collection
	if err := app.Save(collection); err != nil {
		slog.Error("Failed to create categories collection", "error", err)
		return err
	}

	slog.Info("Collection 'categories' created successfully")
	return nil
}

// createProductsCollection creates the "products" collection.
func createProductsCollection(app core.App) error {
	existingCollection, err := app.FindCollectionByNameOrId("products")
	if err == nil && existingCollection != nil {
		slog.Info("Collection 'products' already exists")
		return nil
	}

	collection := core.NewBaseCollection("products")

	// Add fields
	collection.Fields.Add(&core.TextField{
		Name:     "name",
		Required: true,
	})
	collection.Fields.Add(&core.TextField{
		Name: "description",
	})
	collection.Fields.Add(&core.NumberField{
		Name:     "price",
		Required: true,
	})
	collection.Fields.Add(&core.NumberField{
		Name:     "stock",
		Required: true,
	})

	// Add relation field to categories
	categoriesCollection, err := app.FindCollectionByNameOrId("categories")
	if err != nil {
		slog.Error("Failed to find categories collection", "error", err)
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "category",
		Required:     true,
		CollectionId: categoriesCollection.Id,
	})

	// Add file field for images
	collection.Fields.Add(&core.FileField{
		Name:      "images",
		MaxSelect: 5,
		MimeTypes: []string{"image/jpg", "image/jpeg", "image/png", "image/gif"},
		MaxSize:   5242880,
	})

	// Save the collection
	if err := app.Save(collection); err != nil {
		slog.Error("Failed to create products collection", "error", err)
		return err
	}

	slog.Info("Collection 'products' created successfully")
	return nil
}

// createOrdersCollection creates the "orders" collection.
func createOrdersCollection(app core.App) error {
	existingCollection, err := app.FindCollectionByNameOrId("orders")
	if err == nil && existingCollection != nil {
		slog.Info("Collection 'orders' already exists")
		return nil
	}

	collection := core.NewBaseCollection("orders")

	// Add relation field to users
	collection.Fields.Add(&core.RelationField{
		Name:         "user",
		Required:     true,
		CollectionId: "_pb_users_auth_",
	})

	// Add select field for status
	collection.Fields.Add(&core.SelectField{
		Name:     "status",
		Required: true,
		Values:   []string{"pending", "processing", "shipped", "delivered", "cancelled"},
	})

	// Add number field for total amount
	collection.Fields.Add(&core.NumberField{
		Name:     "total_amount",
		Required: true,
	})

	// Add text field for shipping address
	collection.Fields.Add(&core.TextField{
		Name:     "shipping_address",
		Required: true,
	})

	// Add select field for payment status
	collection.Fields.Add(&core.SelectField{
		Name:     "payment_status",
		Required: true,
		Values:   []string{"pending", "paid", "failed", "refunded"},
	})

	// Save the collection
	if err := app.Save(collection); err != nil {
		slog.Error("Failed to create orders collection", "error", err)
		return err
	}

	slog.Info("Collection 'orders' created successfully")
	return nil
}

// createOrderItemsCollection creates the "order_items" collection.
func createOrderItemsCollection(app core.App) error {
	existingCollection, err := app.FindCollectionByNameOrId("order_items")
	if err == nil && existingCollection != nil {
		slog.Info("Collection 'order_items' already exists")
		return nil
	}

	collection := core.NewBaseCollection("order_items")

	// Add relation field to orders
	ordersCollection, err := app.FindCollectionByNameOrId("orders")
	if err != nil {
		slog.Error("Failed to find orders collection", "error", err)
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:          "order",
		Required:      true,
		CollectionId:  ordersCollection.Id,
		CascadeDelete: true,
	})

	// Add relation field to products
	productsCollection, err := app.FindCollectionByNameOrId("products")
	if err != nil {
		slog.Error("Failed to find products collection", "error", err)
		return err
	}
	collection.Fields.Add(&core.RelationField{
		Name:         "product",
		Required:     true,
		CollectionId: productsCollection.Id,
	})

	// Add number field for quantity
	collection.Fields.Add(&core.NumberField{
		Name:     "quantity",
		Required: true,
	})

	// Add number field for price at time of order
	collection.Fields.Add(&core.NumberField{
		Name:     "price_at_time",
		Required: true,
	})

	// Save the collection
	if err := app.Save(collection); err != nil {
		slog.Error("Failed to create order_items collection", "error", err)
		return err
	}

	slog.Info("Collection 'order_items' created successfully")
	return nil
}

// SeedExampleData seeds the database with example data.
func SeedExampleData(app *pocketbase.PocketBase) error {
	// Example categories
	categories := []map[string]interface{}{
		{"name": "Electronics", "description": "Electronic devices and accessories"},
		{"name": "Clothing", "description": "Fashion and accessories"},
		{"name": "Home & Garden", "description": "Home decor and gardening items"},
	}

	categoryIds := make(map[string]string)
	for _, cat := range categories {
		collection, err := app.FindCollectionByNameOrId("categories")
		if err != nil {
			return fmt.Errorf("failed to find categories collection: %v", err)
		}

		// Check if the category already exists
		filter := fmt.Sprintf("name = '%s'", cat["name"])
		existingRecords, err := app.FindRecordsByFilter(collection.Id, filter, "", 1, 0)
		if err != nil {
			return fmt.Errorf("failed to check for existing category: %v", err)
		}

		var record *core.Record
		if len(existingRecords) > 0 {
			// Category already exists, use the existing record
			record = existingRecords[0]
		} else {

			// Category doesn't exist, create a new one
			record = core.NewRecord(collection)
			record.Set("name", cat["name"])
			record.Set("description", cat["description"])

			if err := app.Save(record); err != nil {
				return fmt.Errorf("failed to create category: %v", err)
			}

		}
		categoryIds[cat["name"].(string)] = record.Id
	}

	// Seed Products
	products := []map[string]interface{}{
		{
			"name":        "Pepper Spray Keychain",
			"description": "Black pepper spray keychain",
			"price":       9.99,
			"stock":       50,
			"category":    categoryIds["Electronics"],
		},
		{
			"name":        "Bottle Opener Keychain",
			"description": "Blue bottle opener keychain",
			"price":       6.99,
			"stock":       30,
			"category":    categoryIds["Electronics"],
		},
		{
			"name":        "Cat Punch Keychain",
			"description": "Pink Cat keychain",
			"price":       5.99,
			"stock":       41,
			"category":    categoryIds["Electronics"],
		},
		{
			"name":        "Window Breaker Keychain",
			"description": "Red window breaker keychain",
			"price":       7.99,
			"stock":       75,
			"category":    categoryIds["Electronics"],
		},
	}

	// Create products
	for _, prod := range products {
		collection, err := app.FindCollectionByNameOrId("products")
		if err != nil {
			return fmt.Errorf("failed to find products collection: %v", err)
		}

		// Check if the product already exists
		filter := fmt.Sprintf("name = '%s'", prod["name"])
		existingRecords, err := app.FindRecordsByFilter(collection.Id, filter, "", 1, 0)
		if err != nil {
			return fmt.Errorf("failed to check for existing product: %v", err)
		}

		if len(existingRecords) > 0 {
			// Product already exists, skip creation
			continue
		}

		record := core.NewRecord(collection)
		record.Set("name", prod["name"])
		record.Set("description", prod["description"])
		record.Set("price", prod["price"])
		record.Set("stock", prod["stock"])
		record.Set("category", prod["category"])

		if err := app.Save(record); err != nil {
			return fmt.Errorf("failed to create product: %v", err)
		}
	}

	log.Println("Seeding completed successfully!")
	return nil
}
