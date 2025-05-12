package findfullrecipe

import (
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func WithSinglethreadedDFS(element schema.Element, count int, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		singlethreadedDFS(search, search.Root, count, delay)

		search.Lock()
		search.Finished = true
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Unlock()
	}()

	return searchID
}

func WithMultithreadedDFS(element schema.Element, count int, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		singlethreadedDFS(search, search.Root, count, delay)

		search.TimeTaken = int(time.Since(start).Milliseconds())
	}()

	return searchID
}

func singlethreadedDFS(result *schema.SearchResult, node *schema.SearchNode, count int, delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	node.RLock()
	recipes, _ := database.FindRecipeFor(int(node.Element.ID))
	node.RUnlock()

	if len(recipes) == 0 {
		node.RecipesFound = 1
		updateRecipeCounts(node)
		return
	}

	result.RLock()
	for i := 0; i < len(recipes) && node.RecipesFound < count; i++ {
		result.RUnlock()

		result.Lock()
		result.NodesSearched++
		result.Unlock()

		ingredient1, _ := database.FindElementById(int(recipes[i].Dependency1ID))
		ingredient2, _ := database.FindElementById(int(recipes[i].Dependency2ID))
		combination := schema.Combination{
			Result:      node,
			Ingredient1: &schema.SearchNode{Element: ingredient1},
			Ingredient2: &schema.SearchNode{Element: ingredient2},
		}
		combination.Ingredient1.Parent = &combination
		combination.Ingredient2.Parent = &combination

		node.Lock()
		node.Dependencies = append(node.Dependencies, &combination)
		node.Unlock()

		singlethreadedDFS(result, combination.Ingredient1, count, delay)

		combination.Ingredient1.RLock()
		adjacentCount := max(count/combination.Ingredient1.RecipesFound, 1)
		combination.Ingredient2.RUnlock()

		singlethreadedDFS(result, combination.Ingredient2, adjacentCount, delay)

		result.RLock()
	}
	result.RUnlock()
}

func multithreadedDFS(result *schema.SearchResult, node *schema.SearchNode, count int, delay int) {

}
