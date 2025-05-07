package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Hello World!")

	if database.Initialize() != nil {
		os.Exit(1)
	}
	defer database.Close()

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
