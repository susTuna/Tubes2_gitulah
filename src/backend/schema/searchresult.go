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
	Finished      bool
	lock          sync.RWMutex
}

type SerializedCombination struct {
	Result      int `json:"result"`
	Dependency1 int `json:"dependency1"`
	Dependency2 int `json:"dependency2"`
}

type SerializedSearchResult struct {
	Nodes         []int                   `json:"nodes"`
	Dependencies  []SerializedCombination `json:"dependencies"`
	TimeTaken     int                     `json:"time_taken"`
	NodesSearched int                     `json:"nodes_searched"`
	RecipesFound  int                     `json:"recipes_found"`
	Finished      bool                    `json:"finsihed"`
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
	var w strings.Builder
	json.NewEncoder(&w).Encode(fr.toSerializedIntermediate())
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
	fr.RLock()
	intermediate := SerializedSearchResult{
		Nodes:         []int{int(fr.Root.Element.ID)},
		Dependencies:  []SerializedCombination{},
		TimeTaken:     fr.TimeTaken,
		NodesSearched: fr.NodesSearched,
		RecipesFound:  fr.Root.RecipesFound,
		Finished:      fr.Finished,
	}
	fr.RUnlock()

	fr.Root.RLock()
	dependencies := fr.Root.Dependencies
	fr.Root.RUnlock()

	for i := 0; i < len(dependencies); i++ {
		search1 := SearchResult{Root: dependencies[i].Ingredient1}
		search2 := SearchResult{Root: dependencies[i].Ingredient2}

		intermediate1 := search1.toSerializedIntermediate()
		intermediate2 := search2.toSerializedIntermediate()

		depIndex1 := len(intermediate.Nodes)
		intermediate.Nodes = append(intermediate.Nodes, intermediate1.Nodes...)

		depIndex2 := len(intermediate.Nodes)
		intermediate.Nodes = append(intermediate.Nodes, intermediate2.Nodes...)

		intermediate.Dependencies = append(intermediate.Dependencies, SerializedCombination{
			Result:      0,
			Dependency1: depIndex1,
			Dependency2: depIndex2,
		})

		for j := range intermediate1.Dependencies {
			intermediate1.Dependencies[j].Result += depIndex1
			intermediate1.Dependencies[j].Dependency1 += depIndex1
			intermediate1.Dependencies[j].Dependency2 += depIndex1
		}

		for j := range intermediate2.Dependencies {
			intermediate2.Dependencies[j].Result += depIndex2
			intermediate2.Dependencies[j].Dependency1 += depIndex2
			intermediate2.Dependencies[j].Dependency2 += depIndex2
		}

		intermediate.Dependencies = append(intermediate.Dependencies, intermediate1.Dependencies...)
		intermediate.Dependencies = append(intermediate.Dependencies, intermediate2.Dependencies...)
	}

	return intermediate
}
