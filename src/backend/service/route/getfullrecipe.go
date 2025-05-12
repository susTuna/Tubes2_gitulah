package route

import (
	"net/http"
	"strconv"

	"github.com/filbertengyo/Tubes2_gitulah/service/findfullrecipe"
)

func GetFullRecipe(w http.ResponseWriter, r *http.Request) {
	identifier, err := strconv.ParseInt(r.PathValue("identifier"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	searchResult := findfullrecipe.FindSearch(int(identifier))

	if searchResult == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(searchResult.Serialize()))
}
