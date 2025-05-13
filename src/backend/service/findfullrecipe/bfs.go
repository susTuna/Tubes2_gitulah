package findfullrecipe

import (
	"sync"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func WithSinglethreadedBFS(element schema.Element, count int, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		singlethreadedBFS(search, count, delay)
		cleanupInvalidCombinations(search.Root)

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

		var wg sync.WaitGroup
		setup := make(chan bool)

		go multithreadedBFS(search, search.Root, count, count, delay, &wg, setup)
		<-setup
		wg.Wait()
		cleanupInvalidCombinations(search.Root)

		search.Lock()
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Finished = true
		search.Unlock()
	}()

	return searchID
}

func multithreadedBFS(result *schema.SearchResult, node *schema.SearchNode, count int, topCount int, delay int, wg *sync.WaitGroup, setup chan bool) {
	wg.Add(1)
	defer wg.Done()
	setup <- true

	time.Sleep(time.Duration(delay) * time.Millisecond)

	node.RLock()
	recipes, _ := database.FindRecipeFor(int(node.Element.ID))
	node.RUnlock()

	if len(recipes) == 0 {
		node.RecipesFound = 1

		result.Lock()
		result.NodesSearched++
		result.Unlock()

		updateRecipeCounts(node)
		return
	}

	result.Root.RLock()
	rootContinue := result.Root.RecipesFound < topCount
	result.Root.RUnlock()

	node.RLock()
	for i := 0; i < len(recipes) && i < count && node.RecipesFound < count && rootContinue; i++ {
		node.RUnlock()

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

		setup1 := make(chan bool)
		setup2 := make(chan bool)

		go multithreadedBFS(result, combination.Ingredient1, max(count/len(recipes)/2, 1), topCount, delay, wg, setup1)
		go multithreadedBFS(result, combination.Ingredient2, max(count/len(recipes)/2, 1), topCount, delay, wg, setup2)

		<-setup1
		<-setup2

		result.Root.RLock()
		rootContinue = result.Root.RecipesFound < topCount
		result.Root.RUnlock()

		node.RLock()
	}
	node.RUnlock()
}

func singlethreadedBFS(result *schema.SearchResult, count int, delay int) {
	result.RLock()
	nodes := []*schema.SearchNode{result.Root}
	result.RUnlock()

	iterations := 0

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

				result.Lock()
				result.NodesSearched++
				result.Unlock()

				updateRecipeCounts(nodes[i])
				result.Root.RLock()
				continue
			}

			result.Root.RLock()
			rootContinue := result.Root.RecipesFound < count
			result.Root.RUnlock()

			nodes[i].RLock()
			nodeContinue := nodes[i].RecipesFound < max((count>>iterations)/len(recipes), 1)
			nodes[i].RUnlock()

			for j := 0; j < len(recipes) && j < max((count>>iterations)/len(recipes), 1) && rootContinue && nodeContinue; j++ {
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

				result.Root.RLock()
				rootContinue = result.Root.RecipesFound < count
				result.Root.RUnlock()

				nodes[i].RLock()
				nodeContinue = nodes[i].RecipesFound < max((count>>iterations)/len(recipes), 1)
				nodes[i].RUnlock()
			}

			result.Root.RLock()
		}
		result.Root.RUnlock()

		nodes = nextNodes
		iterations++

		result.Root.RLock()
	}
	result.Root.RUnlock()
}

func cleanupInvalidCombinations(node *schema.SearchNode) {
	node.RLock()
	if node.RecipesFound == 0 {
		if node.Parent != nil {
			result := node.Parent.Result
			combination := node.Parent

			result.Lock()
			newDeps := []*schema.Combination{}
			for _, dep := range result.Dependencies {
				if dep != combination {
					newDeps = append(newDeps, dep)
				}
			}
			result.Dependencies = newDeps
			result.Unlock()
		}
	} else {
		for _, combination := range node.Dependencies {
			cleanupInvalidCombinations(combination.Ingredient1)
			cleanupInvalidCombinations(combination.Ingredient2)
		}
	}
	node.RUnlock()
}
