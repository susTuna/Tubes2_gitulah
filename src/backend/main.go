package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/go-chi/chi/v5"
)

func main() {
	time.Sleep(5 * time.Second)
	fmt.Println("Hello World!")

	if err := database.Initialize(); err != nil {
		fmt.Println("An error occured!")
		fmt.Println(err.Error())
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

	if err := http.ListenAndServe(":5761", r); err != nil {
		fmt.Println("An error occured!")
		fmt.Println(err.Error())
	}
	fmt.Println("Goodbye World!")
}
