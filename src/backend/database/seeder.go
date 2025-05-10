package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

// Add ElementType definition
type ElementType string

const (
	Starting ElementType = "Starting"
	Tier1    ElementType = "Tier1"
	Tier2    ElementType = "Tier2"
	Tier3    ElementType = "Tier3"
	Tier4    ElementType = "Tier4"
	Tier5    ElementType = "Tier5"
	Tier6    ElementType = "Tier6"
	Tier7    ElementType = "Tier7"
	Tier8    ElementType = "Tier8"
	Tier9    ElementType = "Tier9"
	Tier10   ElementType = "Tier10"
	Tier11   ElementType = "Tier11"
	Tier12   ElementType = "Tier12"
	Tier13   ElementType = "Tier13"
	Tier14   ElementType = "Tier14"
	Tier15   ElementType = "Tier15"
)

func Seed() error {
	// Create a new collector
	c := colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"))

	// Map to store element names and their IDs
	elementMap := make(map[string]int32)
	var currentID int32 = 1
	tableIndex := 0

	// First pass: Scrape elements
	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++
		elementType := getElementType(tableIndex)
		if elementType == "" {
			return
		}

		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			name := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if name == "" || name == "Time" || name == "Ruins" || name == "Archeologist" {
				return
			}

			// Get the image URL from the first image in the row
			var imageURL string
			h.ForEach("td:first-of-type a img", func(_ int, img *colly.HTMLElement) {
				imageURL = img.Attr("data-src")
			})

			// Get tier number based on elementType
			tier := getTierNumber(elementType)

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
	})

	// First visit to scrape elements
	err := c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		return fmt.Errorf("failed to scrape elements: %v", err)
	}

	// Reset collector for second pass
	c = colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"))
	tableIndex = 0

	// Second pass: Scrape recipes
	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++
		elementType := getElementType(tableIndex)
		if elementType == "" {
			return
		}

		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			resultName := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if resultName == "" || resultName == "Time" || resultName == "Ruins" || resultName == "Archeologist" {
				return
			}

			// Process each recipe for the current element
			h.ForEach("td:nth-of-type(2) li", func(_ int, li *colly.HTMLElement) {
				aTags := li.DOM.Find("a")
				if aTags.Length() < 2 {
					return
				}

				// Get ingredient names
				ingredient1 := strings.TrimSpace(aTags.Eq(1).Text())
				ingredient2 := strings.TrimSpace(aTags.Eq(3).Text())

				if ingredient1 == "Time" || ingredient2 == "Time" ||
					ingredient1 == "Ruins" || ingredient2 == "Ruins" ||
					ingredient1 == "Archeologist" || ingredient2 == "Archeologist" {
					return
				}

				// Insert recipe if all elements exist in the map
				if resultID, ok := elementMap[resultName]; ok {
					if dep1ID, ok := elementMap[ingredient1]; ok {
						if dep2ID, ok := elementMap[ingredient2]; ok {
							_, err := Exec(`
                                INSERT INTO recipes (result_id, dependency1_id, dependency2_id)
                                VALUES ($1, $2, $3)
                            `, resultID, dep1ID, dep2ID)

							if err != nil {
								log.Printf("Error inserting recipe %s = %s + %s: %v",
									resultName, ingredient1, ingredient2, err)
							}
						}
					}
				}
			})
		})
	})

	// Second visit to scrape recipes
	err = c.Visit("https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)")
	if err != nil {
		return fmt.Errorf("failed to scrape recipes: %v", err)
	}

	return nil
}

func getElementType(index int) ElementType {
	switch index {
	case 1:
		return Starting
	case 2:
		// Special case, we skip it (Ruins/Archeologist)
		return ""
	case 3:
		return Tier1
	case 4:
		return Tier2
	case 5:
		return Tier3
	case 6:
		return Tier4
	case 7:
		return Tier5
	case 8:
		return Tier6
	case 9:
		return Tier7
	case 10:
		return Tier8
	case 11:
		return Tier9
	case 12:
		return Tier10
	case 13:
		return Tier11
	case 14:
		return Tier12
	case 15:
		return Tier13
	case 16:
		return Tier14
	case 17:
		return Tier15
	default:
		return ""
	}
}

func getTierNumber(elementType ElementType) int32 {
	switch elementType {
	case Starting:
		return 0
	case Tier1:
		return 1
	case Tier2:
		return 2
	case Tier3:
		return 3
	case Tier4:
		return 4
	case Tier5:
		return 5
	case Tier6:
		return 6
	case Tier7:
		return 7
	case Tier8:
		return 8
	case Tier9:
		return 9
	case Tier10:
		return 10
	case Tier11:
		return 11
	case Tier12:
		return 12
	case Tier13:
		return 13
	case Tier14:
		return 14
	case Tier15:
		return 15
	default:
		return -1
	}
}
