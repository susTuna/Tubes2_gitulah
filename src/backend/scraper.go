package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gocolly/colly"
)

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

type RecipeType struct {
	Element     string
	ImgUrl1     string
	ImgUrl2     string
	Ingredient1 string
	Ingredient2 string
	Type        ElementType
}

func getElementType(index int) ElementType {
	switch index {
	case 1:
		return Starting
	case 2:
		return "" // Table[2] is special, we skip it (Ruins/Archeologist)
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

func main() {
	url := "https://little-alchemy.fandom.com/wiki/Elements_(Little_Alchemy_2)"
	var recipes []RecipeType

	c := colly.NewCollector(colly.AllowedDomains("little-alchemy.fandom.com"))
	tableIndex := 0
	elementCounter := 0
	recipeCounter := 0

	// each table (starting and tiers)
	c.OnHTML("table.list-table", func(table *colly.HTMLElement) {
		tableIndex++
		elementType := getElementType(tableIndex)
		if elementType == "" {
			return
		}

		// each element generated
		table.ForEach("tbody tr", func(_ int, h *colly.HTMLElement) {
			element := strings.TrimSpace(h.ChildText("td:first-of-type a"))
			if element == "" || element == "Time" || element == "Ruins" || element == "Archeologist" {
				return
			}

			elementCounter++
			fmt.Printf("\nElement[%v]: %-10s | %s\n", elementCounter, element, elementType)

			// each recipe to the element generated
			h.ForEach("td:nth-of-type(2) li", func(_ int, li *colly.HTMLElement) {
				recipeCounter++
				aTags := li.DOM.Find("a")

				if aTags.Length() < 2 {
					return
				}

				imgUrl1, _ := aTags.Eq(0).Find("img").Attr("data-src")
				imgUrl2, _ := aTags.Eq(2).Find("img").Attr("data-src")
				ingredient1 := strings.TrimSpace(aTags.Eq(1).Text())
				ingredient2 := strings.TrimSpace(aTags.Eq(3).Text())

				if ingredient1 == "Time" || ingredient2 == "Time" || ingredient1 == "Ruins" || ingredient2 == "Ruins" || ingredient1 == "Archeologist" || ingredient2 == "Archeologist" {
					return
				}

				r := RecipeType{
					Element:     element,
					ImgUrl1:     imgUrl1,
					ImgUrl2:     imgUrl2,
					Ingredient1: ingredient1,
					Ingredient2: ingredient2,
					Type:        elementType,
				}
				recipes = append(recipes, r)
				fmt.Printf("Recipe[%v]: %s + %s\n", recipeCounter, r.Ingredient1, r.Ingredient2)
				fmt.Printf("ImgUrl1: %s\n", r.ImgUrl1)
				fmt.Printf("ImgUrl2: %s\n", r.ImgUrl2)

			})
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Print("Visiting ", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Print(e.Error())
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatal(err)
	}
}
