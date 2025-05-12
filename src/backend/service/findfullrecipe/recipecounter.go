package findfullrecipe

import "github.com/filbertengyo/Tubes2_gitulah/schema"

func updateRecipeCounts(node *schema.SearchNode) {
	node.RLock()
	combination := node.Parent

	if combination == nil {
		node.RUnlock()
		return
	}

	result := combination.Result
	node.RUnlock()

	result.Lock()
	recipeCount := 0

	for i := 0; i < len(result.Dependencies); i++ {
		result.Dependencies[i].Ingredient1.RLock()
		result.Dependencies[i].Ingredient2.RLock()
		recipeCount += result.Dependencies[i].Ingredient1.RecipesFound * result.Dependencies[i].Ingredient2.RecipesFound
		result.Dependencies[i].Ingredient2.RUnlock()
		result.Dependencies[i].Ingredient1.RUnlock()
	}

	result.RecipesFound = recipeCount
	result.Unlock()

	updateRecipeCounts(result)
}
