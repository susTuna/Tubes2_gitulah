package schema

import (
	"encoding/json"
	"strings"
)

type SearchResponse struct {
	SearchID int `json:"search_id"`
}

func (sr SearchResponse) Serialize() string {
	var w strings.Builder
	json.NewEncoder(&w).Encode(sr)
	return w.String()
}
