package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/filbertengyo/Tubes2_gitulah/database"
)

func Element(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	if val, ok := queries["type"]; ok {
		if len(val) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		} else if val[0] == "name" {
			elementByName(w, r)
		} else if val[0] == "id" {
			elementById(w, r)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		elementById(w, r)
	}
}

func elementById(w http.ResponseWriter, r *http.Request) {
	identifier, err := strconv.ParseInt(r.PathValue("identifier"), 10, 32)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	val, err := database.FindElementById(int(identifier))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(val.Serialize()))
}

func elementByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("identifier")

	val, err := database.FindElementByName(name)

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
