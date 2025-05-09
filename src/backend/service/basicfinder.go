package service

import (
	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func FindElementById(id int) (schema.Element, error) {
	return database.FindElementById(id)
}

func FindElementByName(name string) ([]schema.Element, error) {
	return database.FindElementByName(name)
}

func Elements(start int, end int) ([]schema.Element, error) {
	return database.Elements(int32(start), int32(end))
}

func FindRecipeFor(id int) ([]schema.Recipe, error) {
	return database.FindRecipeFor(id)
}
