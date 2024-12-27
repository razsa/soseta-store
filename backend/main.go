package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"

	"e-com/collection" // Replace with your actual project name
)

// API response structures
type SearchProductsResponse struct {
	Items       []*core.Record `json:"items"`
	TotalItems  int            `json:"totalItems"`
	TotalPages  int            `json:"totalPages"`
	CurrentPage int            `json:"currentPage"`
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

		return se.Next()
	})

	// Create collections and seed example data
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		if err := collection.CreateCollections(app); err != nil {
			slog.Error("Failed to create collections", "error", err)
		}

		if err := collection.SeedExampleData(app); err != nil {
			slog.Error("Error seeding example data", "error", err)
		}

		return se.Next()
	})

	if err := app.Start(); err != nil {
		slog.Error("Failed to start the application", "error", err)
		os.Exit(1)
	}
}
