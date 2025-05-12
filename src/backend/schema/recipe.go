package schema

import (
	"encoding/json"
	"strings"
)

type Recipe struct {
	ResultID      int32 `json:"result_id"`
	Dependency1ID int32 `json:"dependency1_id"`
	Dependency2ID int32 `json:"dependency2_id"`
}

func (r Recipe) Serialize() string {
	var w strings.Builder
	json.NewEncoder(&w).Encode(r)
	return w.String()
}
