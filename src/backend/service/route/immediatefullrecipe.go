package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func ImmediateFullRecipe(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post("localhost:5761/fullrecipe", "application/json", r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var searchResponse schema.SearchResponse
	err = json.NewDecoder(resp.Body).Decode(&searchResponse)
	resp.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for !findfullrecipe.FindSearch(searchResponse.SearchID).Finished {
		time.Sleep(1 * time.Millisecond)
	}

	GetFullRecipe(w, r)
}
