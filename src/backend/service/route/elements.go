package route

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/filbertengyo/Tubes2_gitulah/database"
	"github.com/filbertengyo/Tubes2_gitulah/schema"
)

func Elements(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()

	var err error
	var start uint64 = 0
	var end uint64 = 20
	tiers := []int{}

	if str, ok := queries["start"]; ok {
		if len(str) > 0 {
			start, err = strconv.ParseUint(str[0], 10, 32)
		}
	}

	if str, ok := queries["end"]; ok {
		if len(str) > 0 {
			end, err = strconv.ParseUint(str[0], 10, 32)
		}
	}

	if strs, ok := queries["tiers"]; ok {
		for i := range strs {
			if tier, _err := strconv.ParseUint(strs[i], 10, 32); err != nil {
				err = _err
				break
			} else {
				tiers = append(tiers, int(tier))
			}
		}
	}

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if end < start {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	query := []schema.Element{}

	if len(tiers) > 0 {
		for i := range tiers {
			q, _err := database.FindElementInTier(int(start), int(end), tiers[i])

			if _err != nil {
				err = _err
				break
			}

			query = append(query, q...)
		}
	} else {
		query, err = database.Elements(int(start), int(end))
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	var body bytes.Buffer
	json.NewEncoder(&body).Encode(query)

	w.WriteHeader(http.StatusOK)
	w.Header().Add("Content-Type", "application/json")
	w.Write(body.Bytes())
}
