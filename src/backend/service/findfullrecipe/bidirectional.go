package findfullrecipe

import (
	"sync"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

type bidirectionalChannels struct {
	ready  chan bool
	finish chan bool
	found  chan bool
	close  chan bool
}

func WithSinglethreadedBidirectional(element schema.Element, target schema.Element, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		startNode := search.Root
		targetNode := &schema.SearchNode{Element: target}

		singlethreadedBidirectional(search, startNode, targetNode, delay)

		search.Lock()
		search.Finished = true
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Unlock()
	}()

	return searchID
}

func WithMultithreadedBidirectional(element schema.Element, target schema.Element, delay int) int {
	searchID, search := prepareSearch(element)

	go func() {
		start := time.Now()

		channels := bidirectionalChannels{
			ready:  make(chan bool),
			finish: make(chan bool),
			found:  make(chan bool),
			close:  make(chan bool),
		}

		startNode := search.Root
		targetNode := &schema.SearchNode{Element: target}

		go multithreadedBidirectional(search, startNode, targetNode, channels, delay)

		<-channels.ready
		<-channels.finish
		channels.close <- true

		search.Lock()
		search.Finished = true
		search.TimeTaken = int(time.Since(start).Milliseconds())
		search.Unlock()
	}()

	return searchID
}

func singlethreadedBidirectional(result *schema.SearchResult, startNode *schema.SearchNode, targetNode *schema.SearchNode, delay int) {
	forwardVisited := make(map[int]bool)
	backwardVisited := make(map[int]bool)

	var forwardQueue []*schema.SearchNode
	var backwardQueue []*schema.SearchNode

	forwardQueue = append(forwardQueue, startNode)
	backwardQueue = append(backwardQueue, targetNode)

	forwardVisited[int(startNode.Element.ID)] = true
	backwardVisited[int(targetNode.Element.ID)] = true

	for len(forwardQueue) > 0 && len(backwardQueue) > 0 {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		currentForward := forwardQueue[0]
		forwardQueue = forwardQueue[1:]

		recipes, _ := database.FindRecipeFor(int(currentForward.Element.ID))

		for _, recipe := range recipes {
			result.Lock()
			result.NodesSearched++
			result.Unlock()

			ingredient1, _ := database.FindElementById(int(recipe.Dependency1ID))
			ingredient2, _ := database.FindElementById(int(recipe.Dependency2ID))

			if backwardVisited[int(ingredient1.ID)] || backwardVisited[int(ingredient2.ID)] {
				return
			}

			if !forwardVisited[int(ingredient1.ID)] {
				newNode := &schema.SearchNode{Element: ingredient1}
				forwardQueue = append(forwardQueue, newNode)
				forwardVisited[int(ingredient1.ID)] = true
			}

			if !forwardVisited[int(ingredient2.ID)] {
				newNode := &schema.SearchNode{Element: ingredient2}
				forwardQueue = append(forwardQueue, newNode)
				forwardVisited[int(ingredient2.ID)] = true
			}
		}

		currentBackward := backwardQueue[0]
		backwardQueue = backwardQueue[1:]

		dependentRecipes, _ := database.FindRecipesUsingElement(int(currentBackward.Element.ID))

		for _, recipe := range dependentRecipes {
			result.Lock()
			result.NodesSearched++
			result.Unlock()

			resultElement, _ := database.FindElementById(int(recipe.ResultID))

			if forwardVisited[int(resultElement.ID)] {
				return
			}

			if !backwardVisited[int(resultElement.ID)] {
				newNode := &schema.SearchNode{Element: resultElement}
				backwardQueue = append(backwardQueue, newNode)
				backwardVisited[int(resultElement.ID)] = true
			}
		}
	}
}

func multithreadedBidirectional(result *schema.SearchResult, startNode *schema.SearchNode, targetNode *schema.SearchNode, channels bidirectionalChannels, delay int) {
	forwardVisited := sync.Map{}
	backwardVisited := sync.Map{}

	var forwardQueue []*schema.SearchNode
	var backwardQueue []*schema.SearchNode

	forwardQueue = append(forwardQueue, startNode)
	backwardQueue = append(backwardQueue, targetNode)

	forwardVisited.Store(startNode.Element.ID, true)
	backwardVisited.Store(targetNode.Element.ID, true)

	channels.ready <- true

	var wg sync.WaitGroup
	foundPath := false

	for len(forwardQueue) > 0 && len(backwardQueue) > 0 && !foundPath {
		time.Sleep(time.Duration(delay) * time.Millisecond)

		wg.Add(2)

		go func() {
			defer wg.Done()
			if len(forwardQueue) > 0 {
				currentForward := forwardQueue[0]
				forwardQueue = forwardQueue[1:]

				recipes, _ := database.FindRecipeFor(int(currentForward.Element.ID))

				for _, recipe := range recipes {
					result.Lock()
					result.NodesSearched++
					result.Unlock()

					ingredient1, _ := database.FindElementById(int(recipe.Dependency1ID))
					ingredient2, _ := database.FindElementById(int(recipe.Dependency2ID))

					if _, exists := backwardVisited.Load(ingredient1.ID); exists {
						foundPath = true
						return
					}
					if _, exists := backwardVisited.Load(ingredient2.ID); exists {
						foundPath = true
						return
					}

					if _, visited := forwardVisited.LoadOrStore(ingredient1.ID, true); !visited {
						newNode := &schema.SearchNode{Element: ingredient1}
						forwardQueue = append(forwardQueue, newNode)
					}

					if _, visited := forwardVisited.LoadOrStore(ingredient2.ID, true); !visited {
						newNode := &schema.SearchNode{Element: ingredient2}
						forwardQueue = append(forwardQueue, newNode)
					}
				}
			}
		}()

		go func() {
			defer wg.Done()
			if len(backwardQueue) > 0 {
				currentBackward := backwardQueue[0]
				backwardQueue = backwardQueue[1:]

				dependentRecipes, _ := database.FindRecipesUsingElement(int(currentBackward.Element.ID))

				for _, recipe := range dependentRecipes {
					result.Lock()
					result.NodesSearched++
					result.Unlock()

					resultElement, _ := database.FindElementById(int(recipe.ResultID))

					if _, exists := forwardVisited.Load(resultElement.ID); exists {
						foundPath = true
						return
					}

					if _, visited := backwardVisited.LoadOrStore(resultElement.ID, true); !visited {
						newNode := &schema.SearchNode{Element: resultElement}
						backwardQueue = append(backwardQueue, newNode)
					}
				}
			}
		}()

		wg.Wait()

		select {
		case <-channels.close:
			return
		default:
			if foundPath {
				channels.finish <- true
				return
			}
		}
	}

	channels.finish <- true
}
