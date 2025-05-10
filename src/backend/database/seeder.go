package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

func Seed() error {
	// Create a new collector
	c := colly.NewCollector()

	// Map to store element names and their IDs
	elementMap := make(map[string]int32)
	var currentID int32 = 1

	// First pass: Scrape elements
	c.OnHTML("table.wikitable tbody tr", func(e *colly.HTMLElement) {
		// Skip header row
		if e.ChildText("th") != "" {
			return
		}

		name := strings.TrimSpace(e.ChildText("td:nth-child(1)"))
		if name == "" {
			return
		}

		// Get the image URL
		imageURL := ""
		e.ForEach("td:nth-child(1) img", func(_ int, img *colly.HTMLElement) {
			imageURL = img.Attr("src")
		})

		// Determine tier (starting elements are tier 1, others are tier > 1)
		tier := 2
		if e.ChildText("td:nth-child(3)") == "Starting Element" {
			tier = 1
		}

		// Insert element into database
		_, err := Exec(`
            INSERT INTO elements (id, name, tier, image_url)
            VALUES ($1, $2, $3, $4)
        `, currentID, name, tier, imageURL)

		if err != nil {
			log.Printf("Error inserting element %s: %v", name, err)
			return
		}

		elementMap[name] = currentID
		currentID++
	})

	// Visit the page to scrape elements
	err := c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		return fmt.Errorf("failed to scrape elements: %v", err)
	}

	// Second pass: Scrape recipes
	c.OnHTML("table.wikitable tbody tr", func(e *colly.HTMLElement) {
		if e.ChildText("th") != "" {
			return
		}

		resultName := strings.TrimSpace(e.ChildText("td:nth-child(1)"))
		if resultName == "" {
			return
		}

		// Get combinations
		combinations := strings.Split(e.ChildText("td:nth-child(3)"), "\n")
		for _, combo := range combinations {
			if combo == "Starting Element" {
				continue
			}

			parts := strings.Split(combo, "+")
			if len(parts) != 2 {
				continue
			}

			dep1 := strings.TrimSpace(parts[0])
			dep2 := strings.TrimSpace(parts[1])

			// Insert recipe if all elements exist
			if resultID, ok := elementMap[resultName]; ok {
				if dep1ID, ok := elementMap[dep1]; ok {
					if dep2ID, ok := elementMap[dep2]; ok {
						_, err := Exec(`
                            INSERT INTO recipes (result_id, dependency1_id, dependency2_id)
                            VALUES ($1, $2, $3)
                        `, resultID, dep1ID, dep2ID)

						if err != nil {
							log.Printf("Error inserting recipe %s = %s + %s: %v", resultName, dep1, dep2, err)
						}
					}
				}
			}
		}
	})

	// Visit the page again to scrape recipes
	err = c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		return fmt.Errorf("failed to scrape recipes: %v", err)
	}

	// Print database contents for verification
	PrintDatabaseContents()

	return nil
}

func PrintDatabaseContents() {
	fmt.Println("\n=== Database Contents ===")

	// Print elements
	fmt.Println("\nElements:")
	elements, err := Elements(0, 1000) // Get all elements
	if err != nil {
		log.Printf("Error fetching elements: %v", err)
		return
	}
	for _, elem := range elements {
		fmt.Printf("ID: %d, Name: %s, Tier: %d, ImageURL: %s\n",
			elem.ID, elem.Name, elem.Tier, elem.ImageUrl)
	}

	// Print recipes
	fmt.Println("\nRecipes:")
	rows, err := Query("SELECT r.result_id, e1.name, e2.name, e3.name FROM recipes r " +
		"JOIN elements e1 ON r.result_id = e1.id " +
		"JOIN elements e2 ON r.dependency1_id = e2.id " +
		"JOIN elements e3 ON r.dependency2_id = e3.id")
	if err != nil {
		log.Printf("Error fetching recipes: %v", err)
		return
	}

	var resultID int32
	var resultName, dep1Name, dep2Name string
	for rows.Next() {
		err := rows.Scan(&resultID, &resultName, &dep1Name, &dep2Name)
		if err != nil {
			log.Printf("Error scanning recipe row: %v", err)
			continue
		}
		fmt.Printf("Result: %s (%d) = %s + %s\n", resultName, resultID, dep1Name, dep2Name)
	}
}
