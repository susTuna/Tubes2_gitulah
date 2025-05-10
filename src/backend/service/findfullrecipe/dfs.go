package findfullrecipe

import (
	"sync"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func WithDFS(element schema.Element, count int, delay int, multithreaded bool) int {
	searchID, search := prepareSearch(element)
	if multithreaded {
		go multithreadedDFS(search, count, delay, new(sync.WaitGroup))
	} else {
		go singlethreadedDFS(search, count, delay)
	}
	return searchID
}

func singlethreadedDFS(search *schema.SearchResult, count int, delay int) {
	recipes, _ := database.FindRecipeFor(int(search.Element.ID))

	if len(recipes) == 0 {
		search.RecipesFound = 1
		return
	}

	for i := 0; i < len(recipes) && search.RecipesFound < count; i++ {
		search.NodesSearched++

		ingredient1, _ := database.FindElementById(int(recipes[i].Dependency1ID))
		ingredient2, _ := database.FindElementById(int(recipes[i].Dependency2ID))
		combination := schema.Combination{
			Ingredient1: &schema.SearchResult{Element: ingredient1},
			Ingredient2: &schema.SearchResult{Element: ingredient2},
		}

		singlethreadedDFS(combination.Ingredient1, count, delay)
		singlethreadedDFS(combination.Ingredient2, count, delay)

		search.RecipesFound += combination.Ingredient1.RecipesFound * combination.Ingredient2.RecipesFound
		search.Dependencies = append(search.Dependencies, combination)
	}
}

func multithreadedDFS(search *schema.SearchResult, count int, delay int, wg *sync.WaitGroup) {
	defer wg.Done()

	recipes, _ := database.FindRecipeFor(int(search.Element.ID))

	if len(recipes) == 0 {
		search.RecipesFound = 1
		return
	}

	for i := 0; i < len(recipes) && search.RecipesFound < count; i++ {
		search.NodesSearched++

		ingredient1, _ := database.FindElementById(int(recipes[i].Dependency1ID))
		ingredient2, _ := database.FindElementById(int(recipes[i].Dependency2ID))
		combination := schema.Combination{
			Ingredient1: &schema.SearchResult{Element: ingredient1},
			Ingredient2: &schema.SearchResult{Element: ingredient2},
		}

		wg := new(sync.WaitGroup)
		wg.Add(2)

		go multithreadedDFS(combination.Ingredient2, count, delay, wg)
		multithreadedDFS(combination.Ingredient1, count, delay, wg)

		wg.Wait()

		search.RecipesFound += combination.Ingredient1.RecipesFound * combination.Ingredient2.RecipesFound
		search.Dependencies = append(search.Dependencies, combination)
	}
}
