package route

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"fmt"

	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func ImmediateFullRecipe(w http.ResponseWriter, r *http.Request) {
	port := os.Getenv("BACKEND_HOST_PORT")

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/fullrecipe/", port), "application/json", r.Body)

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

	resp, err = http.Get(fmt.Sprintf("http://localhost:%s/fullrecipe/", port) + fmt.Sprint(searchResponse.SearchID))

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
