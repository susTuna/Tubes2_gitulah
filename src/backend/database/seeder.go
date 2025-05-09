package database

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Element struct {
	ID   int
	Name string
	URL  string
}

func Seed() error {
	elements := make(map[string]Element)

	// Fetch main elements page
	res, err := http.Get("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		return fmt.Errorf("failed to fetch wiki page: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocument(res.Request.URL.String())
	if err != nil {
		return fmt.Errorf("failed to parse HTML: %v", err)
	}

	// First pass: Insert elements and store their IDs
	doc.Find("table.article-table tr").Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			return
		}

		name := strings.TrimSpace(s.Find("td:first-child").Text())
		imageURL := ""
		elementURL := ""

		if img := s.Find("td:first-child img"); img.Length() > 0 {
			imageURL, _ = img.Attr("src")
		}

		if link := s.Find("td:first-child a"); link.Length() > 0 {
			href, exists := link.Attr("href")
			if exists {
				elementURL = "https://little-alchemy.fandom.com" + href
			}
		}

		if name != "" {
			var elementID int
			err := QueryRow(`
                INSERT INTO Elements (name, image_url)
                VALUES ($1, $2)
                RETURNING id
            `, name, imageURL).Scan(&elementID)

			if err != nil {
				log.Printf("Failed to insert element %s: %v", name, err)
				return
			}

			elements[name] = Element{ID: elementID, Name: name, URL: elementURL}
			log.Printf("Inserted element: %s (ID: %d)", name, elementID)
		}
	})

	// Second pass: Scrape recipes for each element
	for _, element := range elements {
		if element.URL == "" {
			continue
		}

		// Add delay to avoid overwhelming the server
		time.Sleep(500 * time.Millisecond)

		res, err := http.Get(element.URL)
		if err != nil {
			log.Printf("Failed to fetch element page for %s: %v", element.Name, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Printf("Failed to parse element page for %s: %v", element.Name, err)
			continue
		}

		// Find the combinations section
		doc.Find(`div.combination-table table tr`).Each(func(i int, s *goquery.Selection) {
			if i == 0 {
				return
			}

			// Extract the two elements that make up the combination
			deps := s.Find("td a").Map(func(i int, s *goquery.Selection) string {
				return strings.TrimSpace(s.Text())
			})

			if len(deps) >= 2 {
				dep1, exists1 := elements[deps[0]]
				dep2, exists2 := elements[deps[1]]

				if exists1 && exists2 {
					// Insert recipe
					_, err := Exec(`
                        INSERT INTO Recipes (result_id, dependency1_id, dependency2_id)
                        VALUES ($1, $2, $3)
                        ON CONFLICT DO NOTHING
                    `, element.ID, dep1.ID, dep2.ID)

					if err != nil {
						log.Printf("Failed to insert recipe for %s (%s + %s): %v",
							element.Name, deps[0], deps[1], err)
					} else {
						log.Printf("Inserted recipe: %s = %s + %s", element.Name, deps[0], deps[1])
					}
				}
			}
		})

		res.Body.Close()
	}
	displayScrapingResults()
	return nil
}

// Add this function after existing Seed() function
func displayScrapingResults() {
	fmt.Println("\n=== Scraping Results Summary ===")

	var elementCount int
	var recipeCount int

	err := QueryRow("SELECT COUNT(*) FROM Elements").Scan(&elementCount)
	if err != nil {
		log.Printf("Error counting elements: %v", err)
		return
	}

	err = QueryRow("SELECT COUNT(*) FROM Recipes").Scan(&recipeCount)
	if err != nil {
		log.Printf("Error counting recipes: %v", err)
		return
	}

	fmt.Printf("\nTotal Elements: %d\n", elementCount)
	fmt.Printf("Total Recipes: %d\n", recipeCount)

	fmt.Println("\nSample Elements (First 5):")
	rows, err := Query(`
        SELECT id, name, image_url 
        FROM Elements 
        ORDER BY id 
        LIMIT 5
    `)
	if err != nil {
		log.Printf("Error querying elements: %v", err)
		return
	}
	defer rows.Close()

	fmt.Printf("\n%-5s | %-30s | %s\n", "ID", "Name", "Image URL")
	fmt.Println(strings.Repeat("-", 80))
	for rows.Next() {
		var id int
		var name, imageUrl string
		rows.Scan(&id, &name, &imageUrl)
		fmt.Printf("%-5d | %-30s | %s\n", id, name, imageUrl)
	}

	fmt.Println("\nSample Recipes (First 5):")
	recipeRows, err := Query(`
        SELECT 
            e1.name as result,
            e2.name as ingredient1,
            e3.name as ingredient2
        FROM Recipes r
        JOIN Elements e1 ON r.result_id = e1.id
        JOIN Elements e2 ON r.dependency1_id = e2.id
        JOIN Elements e3 ON r.dependency2_id = e3.id
        LIMIT 5
    `)
	if err != nil {
		log.Printf("Error querying recipes: %v", err)
		return
	}
	defer recipeRows.Close()

	fmt.Printf("\n%-30s = %-30s + %s\n", "Result", "Ingredient 1", "Ingredient 2")
	fmt.Println(strings.Repeat("-", 80))
	for recipeRows.Next() {
		var result, ing1, ing2 string
		recipeRows.Scan(&result, &ing1, &ing2)
		fmt.Printf("%-30s = %-30s + %s\n", result, ing1, ing2)
	}

	fmt.Println("\n=== End of Results ===")
}
