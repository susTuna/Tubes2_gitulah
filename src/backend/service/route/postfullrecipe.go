package route

import (
	"encoding/json"
	"net/http"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func PostFullRecipe(w http.ResponseWriter, r *http.Request) {
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
		response.SearchID = findfullrecipe.WithMultithreadedDFS(element, request.Count, request.Delay)
	} else if request.Method == "bfs" && request.Threading == "single" {
		response.SearchID = findfullrecipe.WithSinglethreadedBFS(element, request.Count, request.Delay)
	} else if request.Method == "bfs" && request.Threading == "multi" {
		response.SearchID = findfullrecipe.WithMultithreadedBFS(element, request.Count, request.Delay)
		//} else if request.Method == "bidirectional" && request.Threading == "single" {
		//	response.SearchID = findfullrecipe.WithSinglethreadedBidirectional(element, request.Count, request.Delay)
		//} else if request.Method == "bidirectional" && request.Threading == "multi" {

	} else {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(response.Serialize()))
}
