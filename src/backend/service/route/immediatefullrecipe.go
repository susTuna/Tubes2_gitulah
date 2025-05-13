package route

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"fmt"

	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func ImmediateFullRecipe(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Post("http://localhost:5761/fullrecipe/", "application/json", r.Body)

	if err != nil {
		w.WriteHeader(resp.StatusCode)
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

	resp, err = http.Get("http://localhost:5761/fullrecipe/" + fmt.Sprint(searchResponse.SearchID))

	if err != nil {
		w.WriteHeader(resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(body)
}
