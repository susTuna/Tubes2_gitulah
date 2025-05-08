package database

import (
	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/jackc/pgx/v5"
)

func FindElement(id int) (schema.Element, error) {
	queryResult := QueryRow(`SELECT * FROM Elements WHERE id=$1`, id)

	var element schema.Element

	err := queryResult.Scan(&element.ID, &element.Name, &element.ImageUrl)

	return element, err
}

func FindRecipeFor(elementID int) ([]schema.Recipe, error) {
	recipes := []schema.Recipe{}

	queryResult, err := Query(`SELECT * FROM Recipes WHERE result_id=$1`, elementID)
	if err != nil {
		return recipes, err
	}

	var i, j, k int32
	_, err = pgx.ForEachRow(queryResult, []any{&i, &j, &k}, func() error {
		recipes = append(recipes, schema.Recipe{
			ResultID:      i,
			Dependency1ID: j,
			Dependency2ID: k,
		})
		return nil
	})

	return recipes, err
}
