package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/filesystem"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
}

func importProducts(app *pocketbase.PocketBase) error {
	// Read the JSON file
	jsonData, err := os.ReadFile("products.json")
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}

	// Parse JSON into products slice
	var products []Product
	if err := json.Unmarshal(jsonData, &products); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	// Get the products collection
	collection, err := app.FindCollectionByNameOrId("products")
	if err != nil {
		return fmt.Errorf("error finding collection: %v", err)
	}

	// Create a default category if it doesn't exist
	categoriesCollection, err := app.FindCollectionByNameOrId("categories")
	if err != nil {
		return fmt.Errorf("error finding categories collection: %v", err)
	}

	// First create the Electronics category if it doesn't exist
	categoryRecord := core.NewRecord(categoriesCollection)
	categoryRecord.Set("name", "Electronics")
	categoryRecord.Set("description", "Electronic devices and accessories")

	if err := app.Save(categoryRecord); err != nil {
		log.Printf("Note: Category might already exist or there was an error: %v", err)
		// Try to find existing category
		categoryRecords, err := app.FindFirstRecordByData(categoriesCollection.Id, "name", "Electronics")
		if err != nil {
			return fmt.Errorf("error finding/creating category: %v", err)
		}
		categoryRecord = categoryRecords
	}

	// Import each product
	for _, product := range products {
		record := core.NewRecord(collection)

		// Set the basic fields
		record.Set("name", product.Name)
		record.Set("description", product.Description)
		record.Set("price", product.Price)
		record.Set("stock", 100) // Default stock value
		record.Set("category", categoryRecord.Id)

		// Handle the image - download from URL
		imageFile, err := filesystem.NewFileFromURL(context.Background(), product.Image)
		if err != nil {
			log.Printf("Error downloading image for product %d: %v", product.ID, err)
			continue
		}
		record.Set("images", []*filesystem.File{imageFile})

		// Save the record
		if err := app.Save(record); err != nil {
			log.Printf("Error saving product %d: %v", product.ID, err)
			continue
		}

		fmt.Printf("Successfully imported product: %s\n", product.Name)
	}

	fmt.Println("Import completed!")
	return nil
}
