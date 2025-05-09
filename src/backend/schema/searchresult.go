package schema

import (
	"encoding/json"
	"strings"
)

type Combination struct {
	Ingredient1 *SearchResult
	Ingredient2 *SearchResult
}

type SearchResult struct {
	Element       Element
	Dependencies  []Combination
	TimeTaken     int
	NodesSearched int
	RecipesFound  int
}

func (fr SearchResult) Serialize() string {
	var w *strings.Builder
	json.NewEncoder(w).Encode(fr.toSerializedIntermediate())
	return w.String()
}

type SerializedCombination struct {
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

func (fr SearchResult) toSerializedIntermediate() SerializedSearchResult {
	intermediate := SerializedSearchResult{
		Nodes:         []int{int(fr.Element.ID)},
		Dependencies:  []SerializedCombination{},
		TimeTaken:     fr.TimeTaken,
		NodesSearched: fr.NodesSearched,
		RecipesFound:  fr.RecipesFound,
	}

	for i := 0; i < len(fr.Dependencies); i++ {
		intermediate.Dependencies = append(intermediate.Dependencies, SerializedCombination{
			Dependency1: int(fr.Dependencies[i].Ingredient1.Element.ID),
			Dependency2: int(fr.Dependencies[i].Ingredient2.Element.ID),
		})

		intermediate1 := fr.Dependencies[i].Ingredient1.toSerializedIntermediate()
		intermediate2 := fr.Dependencies[i].Ingredient2.toSerializedIntermediate()

		intermediate.Nodes = append(intermediate.Nodes, intermediate1.Nodes...)
		intermediate.Nodes = append(intermediate.Nodes, intermediate2.Nodes...)

		intermediate.Dependencies = append(intermediate.Dependencies, intermediate1.Dependencies...)
		intermediate.Dependencies = append(intermediate.Dependencies, intermediate2.Dependencies...)
	}

	return intermediate
}
