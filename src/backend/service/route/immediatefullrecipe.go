package route

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func ImmediateFullRecipe(w http.ResponseWriter, r *http.Request) {
	var request schema.SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	if !request.Valid() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	element, err := database.FindElementById(request.Element)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var response schema.SearchResponse

	if request.Method == "dfs" && request.Threading == "single" {
		response.SearchID = findfullrecipe.WithSinglethreadedDFS(element, request.Count, request.Delay)
	} else if request.Method == "dfs" && request.Threading == "multi" {
		response.SearchID = -1
	} else if request.Method == "bfs" && request.Threading == "single" {
		response.SearchID = findfullrecipe.WithSinglethreadedBFS(element, request.Count, request.Delay)
	} else if request.Method == "bfs" && request.Threading == "multi" {
		response.SearchID = -1
	} else {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	searchResult := findfullrecipe.FindSearch(int(response.SearchID))

	if searchResult == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for !searchResult.Finished {
		time.Sleep(time.Millisecond * 1)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(searchResult.Serialize()))
}
