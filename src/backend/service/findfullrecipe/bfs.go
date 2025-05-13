package findfullrecipe

import (
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

type searchJob struct {
	node  *schema.SearchNode
	count int
	delay int
}

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

		numWorkers := 4
		jobs := make(chan searchJob, 100)
		results := make(chan bool, 100)
		done := make(chan bool)

		for i := 0; i < numWorkers; i++ {
			go worker(search, jobs, results)
		}

		go func() {
			jobs <- searchJob{
				node:  search.Root,
				count: count,
				delay: delay,
			}

			for {
				select {
				case <-results:
					if search.Root.RecipesFound >= count {
						done <- true
						return
					}
				case <-time.After(2 * time.Second):
					done <- true
					return
				}
			}
		}()

		<-done
		close(jobs)
		close(results)

		cleanupInvalidCombinations(search.Root)

		search.Lock()
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Finished = true
		search.Unlock()
	}()

	return searchID
}

func worker(result *schema.SearchResult, jobs chan searchJob, results chan bool) {
	for job := range jobs {
		processNode(result, job.node, job.count, job.delay, jobs, results)
	}
}

func processNode(result *schema.SearchResult, node *schema.SearchNode, count int, delay int, jobs chan searchJob, results chan bool) {
	recipes, err := database.FindRecipeFor(int(node.Element.ID))
	if err != nil {
		results <- false
		return
	}

	if len(recipes) == 0 {
		node.Lock()
		node.RecipesFound = 1
		node.Unlock()
		updateRecipeCounts(node)
		results <- true
		return
	}

	for _, recipe := range recipes {
		result.Root.RLock()
		if result.Root.RecipesFound >= count {
			result.Root.RUnlock()
			results <- true
			return
		}
		result.Root.RUnlock()

		time.Sleep(time.Duration(delay) * time.Millisecond)

		result.Lock()
		result.NodesSearched++
		result.Unlock()

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

		node.Lock()
		node.Dependencies = append(node.Dependencies, &combination)
		node.Unlock()

		jobs <- searchJob{
			node:  combination.Ingredient1,
			count: count,
			delay: delay,
		}
		jobs <- searchJob{
			node:  combination.Ingredient2,
			count: count,
			delay: delay,
		}
	}

	results <- true
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
