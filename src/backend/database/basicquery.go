package database

import (
	"strings"

	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/jackc/pgx/v5"
)

func FindElementById(id int) (schema.Element, error) {
	queryResult := QueryRow(`SELECT * FROM Elements WHERE id=$1`, id)

	var element schema.Element

	err := queryResult.Scan(&element.ID, &element.Name, &element.Tier, &element.ImageUrl)

	return element, err
}

func FindElementByName(name string) ([]schema.Element, error) {
	elements := []schema.Element{}

	queryResult, err := Query(`SELECT * FROM Elements WHERE LOWER(name) LIKE $1`, "%"+strings.ToLower(name)+"%")
	if err != nil {
		return elements, err
	}

	var i int32
	var n string
	var t int32
	var u string
	_, err = pgx.ForEachRow(queryResult, []any{&i, &n, &t, &u}, func() error {
		elements = append(elements, schema.Element{
			ID:       i,
			Name:     n,
			Tier:     t,
			ImageUrl: u,
		})
		return nil
	})

	return elements, err
}

func Elements(start int32, end int32) ([]schema.Element, error) {
	elements := []schema.Element{}

	queryResult, err := Query(`SELECT * FROM Elements ORDER BY name LIMIT $1 OFFSET $2`, end-start, start)
	if err != nil {
		return elements, err
	}

	var i int32
	var n string
	var t int32
	var u string
	_, err = pgx.ForEachRow(queryResult, []any{&i, &n, &t, &u}, func() error {
		elements = append(elements, schema.Element{
			ID:       i,
			Name:     n,
			Tier:     t,
			ImageUrl: u,
		})
		return nil
	})

	return elements, err
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
