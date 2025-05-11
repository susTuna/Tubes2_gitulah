package schema

import (
	"encoding/json"
	"strings"
	"sync"
)

type Combination struct {
	Result      *SearchNode
	Ingredient1 *SearchNode
	Ingredient2 *SearchNode
}

type SearchNode struct {
	Element      Element
	Parent       *Combination
	Dependencies []*Combination
	RecipesFound int
	lock         sync.RWMutex
}

type SearchResult struct {
	Root          *SearchNode
	TimeTaken     int
	NodesSearched int
	lock          sync.RWMutex
}

type SerializedCombination struct {
	Result      int
	Dependency1 int
	Dependency2 int
}

type SerializedSearchResult struct {
	Nodes         []int
	Dependencies  []SerializedCombination
	TimeTaken     int
	NodesSearched int
	RecipesFound  int
}

func (fr *SearchNode) Lock() {
	fr.lock.Lock()
}

func (fr *SearchNode) Unlock() {
	fr.lock.Unlock()
}

func (fr *SearchNode) RLock() {
	fr.lock.RLock()
}

func (fr *SearchNode) RUnlock() {
	fr.lock.RUnlock()
}

func (fr *SearchResult) Serialize() string {
	var w *strings.Builder
	json.NewEncoder(w).Encode(fr.toSerializedIntermediate())
	return w.String()
}

func (fr *SearchResult) Lock() {
	fr.lock.Lock()
}

func (fr *SearchResult) Unlock() {
	fr.lock.Unlock()
}

func (fr *SearchResult) RLock() {
	fr.lock.RLock()
}

func (fr *SearchResult) RUnlock() {
	fr.lock.RUnlock()
}

func (fr *SearchResult) toSerializedIntermediate() SerializedSearchResult {
	intermediate := SerializedSearchResult{
		Nodes:         []int{int(fr.Root.Element.ID)},
		Dependencies:  []SerializedCombination{},
		TimeTaken:     fr.TimeTaken,
		NodesSearched: fr.NodesSearched,
		RecipesFound:  fr.Root.RecipesFound,
	}

	for i := 0; i < len(fr.Root.Dependencies); i++ {
		intermediate.Nodes = append(intermediate.Nodes,
			int(fr.Root.Dependencies[i].Ingredient1.Element.ID),
			int(fr.Root.Dependencies[i].Ingredient2.Element.ID),
		)

		intermediate.Dependencies = append(intermediate.Dependencies, SerializedCombination{
			Result:      0,
			Dependency1: len(intermediate.Nodes) - 2,
			Dependency2: len(intermediate.Nodes) - 1,
		})

		search1 := SearchResult{Root: fr.Root.Dependencies[i].Ingredient1}
		search2 := SearchResult{Root: fr.Root.Dependencies[i].Ingredient2}

		intermediate1 := search1.toSerializedIntermediate()
		intermediate2 := search2.toSerializedIntermediate()

		for j := 0; j < len(intermediate1.Dependencies); j++ {
			intermediate1.Dependencies[j].Result += len(intermediate.Nodes)
			intermediate1.Dependencies[j].Dependency1 += len(intermediate.Nodes)
			intermediate1.Dependencies[j].Dependency2 += len(intermediate.Nodes)
		}

		intermediate.Dependencies = append(intermediate.Dependencies, intermediate1.Dependencies...)

		for j := 0; j < len(intermediate2.Dependencies); j++ {
			intermediate2.Dependencies[j].Result += len(intermediate.Nodes)
			intermediate2.Dependencies[j].Dependency1 += len(intermediate.Nodes)
			intermediate2.Dependencies[j].Dependency2 += len(intermediate.Nodes)
		}

		intermediate.Dependencies = append(intermediate.Dependencies, intermediate2.Dependencies...)
	}

	return intermediate
}
