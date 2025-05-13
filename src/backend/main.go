package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
	"github.com/filbertengyo/Tubes2_gitulah/service/middleware"
	"github.com/filbertengyo/Tubes2_gitulah/service/route"
)

func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("Hello World!")

	if err := database.Initialize(); err != nil {
		fmt.Println("An error occured during initialization!")
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer database.Close()

	if !database.IsDefined() {
		database.Define()

		if err := database.Seed(); err != nil {
			fmt.Println("An error occured during seeding!")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	findfullrecipe.InitializeSearchCleaner()
	defer findfullrecipe.DeinitializeSearchCleaner()

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if val, ok := q["msg"]; ok {
			fmt.Fprintf(w, "Echo: %v", val[0])
		} else {
			fmt.Fprint(w, "Hello World!")
		}
	})

	http.HandleFunc("GET /elements", route.Elements)

	http.HandleFunc("GET /elements/{identifier}", route.Element)

	http.HandleFunc("GET /elements/{identifier}/recipe", route.Recipe)

	http.HandleFunc("POST /fullrecipe/", route.PostFullRecipe)

	http.HandleFunc("GET /fullrecipe/immediate", route.ImmediateFullRecipe)

	http.HandleFunc("GET /fullrecipe/{identifier}", route.GetFullRecipe)

	handler := middleware.CORSMiddleware(http.DefaultServeMux)

	port := os.Getenv("BACKEND_HOST_PORT")

	if err := http.ListenAndServe(":"+port, handler); err != nil {
		fmt.Println("An error occured!")
		fmt.Println(err.Error())
	}

	fmt.Println("Goodbye World!")
}
