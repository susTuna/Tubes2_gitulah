package findfullrecipe

import (
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func WithSinglethreadedBFS(element schema.Element, count int, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		singlethreadedBFS(search, count, delay)

		search.Lock()
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Finished = true
		search.Unlock()
	}()

	return searchID
}

func WithMultithreadedBFS(element schema.Element, count int, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		singlethreadedBFS(search, count, delay)

		search.Lock()
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Finished = true
		search.Unlock()
	}()

	return searchID
}

func singlethreadedBFS(result *schema.SearchResult, count int, delay int) {
	result.RLock()
	nodes := []*schema.SearchNode{result.Root}
	result.RUnlock()

	result.Root.RLock()
	for len(nodes) > 0 && result.Root.RecipesFound < count {
		result.Root.RUnlock()

		nextNodes := []*schema.SearchNode{}

		result.Root.RLock()
		for i := 0; i < len(nodes) && result.Root.RecipesFound < count; i++ {
			result.Root.RUnlock()

			nodes[i].RLock()
			recipes, _ := database.FindRecipeFor(int(nodes[i].Element.ID))
			nodes[i].RUnlock()

			if len(recipes) == 0 {
				nodes[i].Lock()
				nodes[i].RecipesFound = 1
				nodes[i].Unlock()

				updateRecipeCounts(nodes[i])
			}

			for j := 0; j < len(recipes); j++ {
				time.Sleep(time.Duration(delay) * time.Millisecond)

				result.Lock()
				result.NodesSearched++
				result.Unlock()

				ingredient1, _ := database.FindElementById(int(recipes[j].Dependency1ID))
				ingredient2, _ := database.FindElementById(int(recipes[j].Dependency2ID))
				combination := schema.Combination{
					Result:      nodes[i],
					Ingredient1: &schema.SearchNode{Element: ingredient1},
					Ingredient2: &schema.SearchNode{Element: ingredient2},
				}
				combination.Ingredient1.Parent = &combination
				combination.Ingredient2.Parent = &combination

				nodes[i].Lock()
				nodes[i].Dependencies = append(nodes[i].Dependencies, &combination)
				nodes[i].Unlock()

				nextNodes = append(nextNodes, combination.Ingredient1, combination.Ingredient2)
			}

			result.Root.RLock()
		}
		result.Root.RUnlock()

		nodes = nextNodes

		result.Root.RLock()
	}
	result.Root.RUnlock()

	result.RLock()
	nodes = []*schema.SearchNode{result.Root}
	result.RUnlock()

	for len(nodes) > 0 {
		nextNodes := []*schema.SearchNode{}

		for i := 0; i < len(nodes); i++ {
			nodes[i].RLock()
			if nodes[i].RecipesFound > 0 {
				for j := 0; j < len(nodes[i].Dependencies); j++ {
					nextNodes = append(nextNodes, nodes[i].Dependencies[j].Ingredient1, nodes[i].Dependencies[j].Ingredient2)
				}
			} else {
				combination := nodes[i].Parent
				result := combination.Result
				newDependency := []*schema.Combination{}

				result.RLock()
				for j := 0; j < len(result.Dependencies); j++ {
					if result.Dependencies[j] != combination {
						newDependency = append(newDependency, result.Dependencies[j])
					}
				}
				result.RUnlock()

				result.Lock()
				result.Dependencies = newDependency
				result.Unlock()
			}
			nodes[i].RUnlock()
		}

		nodes = nextNodes
	}
}
