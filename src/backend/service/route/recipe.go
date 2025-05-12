package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/filbertengyo/Tubes2_gitulah/database"
)

func Recipe(w http.ResponseWriter, r *http.Request) {
	identifier, err := strconv.ParseInt(r.PathValue("identifier"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := database.FindRecipeFor(int(identifier))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var body bytes.Buffer
	json.NewEncoder(&body).Encode(val)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(body.Bytes())
}
