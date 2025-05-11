package findfullrecipe

import (
	"sync"

	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

var searchesMutex sync.RWMutex
var searches map[int]*schema.SearchResult = make(map[int]*schema.SearchResult)
var searchID int = 0

func prepareSearch(element schema.Element) (int, *schema.SearchResult) {
	searchResult := &schema.SearchResult{
		Root: &schema.SearchNode{Element: element},
	}

	searchesMutex.Lock()
	defer searchesMutex.Unlock()

	searches[searchID] = searchResult
	searchID += 1

	return searchID - 1, searchResult
}

func FindSearch(searchID int) *schema.SearchResult {
	searchesMutex.RLock()
	defer searchesMutex.RUnlock()

	return searches[searchID]
}
