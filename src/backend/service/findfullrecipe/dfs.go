package findfullrecipe

import (
	"math"
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

		channels := multiDFSChannels{
			ready:        make(chan bool),
			finish:       make(chan int),
			redistribute: make(chan int),
			close:        make(chan bool),
		}

		multithreadedDFS(search, search.Root, channels, delay)

		<-channels.ready
		channels.redistribute <- count
		<-channels.finish
		channels.close <- true

		search.Lock()
		search.Finished = true
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Unlock()
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

	node.RLock()
	for i := 0; i < len(recipes) && count > 0; i++ {
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

		singlethreadedDFS(result, combination.Ingredient1, count, delay)

		combination.Ingredient1.RLock()
		adjacentCount := max(count/combination.Ingredient1.RecipesFound, 1)
		combination.Ingredient1.RUnlock()

		singlethreadedDFS(result, combination.Ingredient2, adjacentCount, delay)

		count -= combination.Ingredient1.RecipesFound * combination.Ingredient2.RecipesFound

		node.RLock()
	}
	node.RUnlock()
}

type multiDFSChannels struct {
	ready        chan bool
	finish       chan int
	redistribute chan int
	close        chan bool
}

func multithreadedDFS(result *schema.SearchResult, node *schema.SearchNode, channels multiDFSChannels, delay int) {
	time.Sleep(time.Duration(delay) * time.Millisecond)

	node.RLock()
	recipes, _ := database.FindRecipeFor(int(node.Element.ID))
	node.RUnlock()

	quota := 1
	i := 0

	channels.ready <- true

	if len(recipes) == 0 {
		node.RecipesFound = 1
		<-channels.redistribute
		channels.finish <- node.RecipesFound
		select {
		case <-channels.redistribute:
			channels.finish <- node.RecipesFound
			<-channels.close
		case <-channels.close:
		}
		return
	}

	for node.RecipesFound < quota {
		select {
		case quota = <-channels.redistribute:
		case <-channels.close:
			return
		}

		node.RLock()
		for ; i < len(recipes) && node.RecipesFound < quota; i++ {
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

			leftChannels := multiDFSChannels{
				ready:        make(chan bool),
				finish:       make(chan int),
				redistribute: make(chan int),
				close:        make(chan bool),
			}

			rightChannels := multiDFSChannels{
				ready:        make(chan bool),
				finish:       make(chan int),
				redistribute: make(chan int),
				close:        make(chan bool),
			}

			go multithreadedDFS(result, combination.Ingredient1, leftChannels, delay)
			go multithreadedDFS(result, combination.Ingredient2, rightChannels, delay)

			leftDistribution := math.Sqrt(float64(quota))
			rightDistribution := float64(quota) / leftDistribution

			ldr := int(math.Ceil(leftDistribution))
			rdr := int(math.Ceil(rightDistribution))

			<-leftChannels.ready
			leftChannels.redistribute <- ldr

			<-rightChannels.ready
			rightChannels.redistribute <- rdr

			leftFound := <-leftChannels.finish
			rightFound := <-rightChannels.finish

			if leftFound < ldr && rightFound >= rdr {
				rdr = quota / leftFound
				rightChannels.redistribute <- rdr
				rightFound = <-rightChannels.finish
			} else if rightFound < rdr && leftFound >= ldr {
				ldr = quota / rightFound
				leftChannels.redistribute <- ldr
				leftFound = <-leftChannels.finish
			}

			leftChannels.close <- true
			rightChannels.close <- true

			node.Lock()
			node.RecipesFound += leftFound * rightFound
			node.Unlock()

			node.RLock()
		}
		node.RUnlock()

		channels.finish <- node.RecipesFound
		if len(recipes) == len(node.Dependencies) {
			select {
			case <-channels.redistribute:
				channels.finish <- node.RecipesFound
				<-channels.close
			case <-channels.close:
			}
			return
		}
	}
}
