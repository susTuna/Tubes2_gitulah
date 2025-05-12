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

		done := make(chan bool)

		go bfssearchhandle(search, search.Root, count, delay, done)
		<-done
		cleanupInvalidCombinations(search.Root)

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

func bfssearchhandle(result *schema.SearchResult, node *schema.SearchNode, count int, delay int, done chan bool) {
	defer func() {
		done <- true
	}()
	recipes, err := database.FindRecipeFor(int(node.Element.ID))
	if err != nil {
		return
	}

	if len(recipes) == 0 {
		node.Lock()
		node.RecipesFound = 1
		node.Unlock()
		updateRecipeCounts(node)
		return
	}

	for _, recipe := range recipes {
		checkDone := make(chan bool, 1)
		go func() {
			result.Root.RLock()
			if result.Root.RecipesFound >= count {
				checkDone <- true
			} else {
				checkDone <- false
			}
			result.Root.RUnlock()
		}()

		select {
		case isDone := <-checkDone:
			if isDone {
				return
			}
		case <-time.After(100 * time.Millisecond):
		}

		time.Sleep(time.Duration(delay) * time.Millisecond)

		updateDone := make(chan bool, 1)
		go func() {
			result.Lock()
			result.NodesSearched++
			result.Unlock()
			updateDone <- true
		}()

		select {
		case <-updateDone:
		case <-time.After(100 * time.Millisecond):
		}

		// Get ingredients
		ingredient1, err1 := database.FindElementById(int(recipe.Dependency1ID))
		ingredient2, err2 := database.FindElementById(int(recipe.Dependency2ID))
		if err1 != nil || err2 != nil {
			continue
		}

		combination := schema.Combination{
			Result:      node,
			Ingredient1: &schema.SearchNode{Element: ingredient1},
			Ingredient2: &schema.SearchNode{Element: ingredient2},
		}
		combination.Ingredient1.Parent = &combination
		combination.Ingredient2.Parent = &combination

		nodeLockDone := make(chan bool, 1)
		go func() {
			node.Lock()
			node.Dependencies = append(node.Dependencies, &combination)
			node.Unlock()
			nodeLockDone <- true
		}()

		select {
		case <-nodeLockDone:
		case <-time.After(100 * time.Millisecond):
		}

		done1 := make(chan bool, 1)
		done2 := make(chan bool, 1)

		go func() {
			bfssearchhandle(result, combination.Ingredient1, count, delay, done1)
		}()
		go func() {
			bfssearchhandle(result, combination.Ingredient2, count, delay, done2)
		}()

		select {
		case <-done1:
		case <-time.After(5 * time.Second):
		}

		select {
		case <-done2:
		case <-time.After(5 * time.Second):
		}
	}
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
