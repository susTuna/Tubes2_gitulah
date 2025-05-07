package main

import (
	"fmt"
	"net/http"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Hello World!")

	database.Initialize()
	if !database.IsDefined() {
		database.Define()
		database.Seed()
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.ListenAndServe(":5761", r)
}
