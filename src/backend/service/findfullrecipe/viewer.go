package findfullrecipe

import "github.com/filbertengyo/Tubes2_gitulah/schema"

var searches map[int]*schema.SearchResult = make(map[int]*schema.SearchResult)
var searchID int = 0

func prepareSearch(element schema.Element) (int, *schema.SearchResult) {
	searchResult := &schema.SearchResult{
		Element: element,
	}

	searches[searchID] = searchResult
	searchID += 1

	return searchID - 1, searchResult
}

func FindSearch(searchID int) *schema.SearchResult {
	return searches[searchID]
}
