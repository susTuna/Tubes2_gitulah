package schema

import (
	"encoding/json"
	"strings"
)

type Recipe struct {
	ResultID      int32
	Dependency1ID int32
	Dependency2ID int32
}

func (r Recipe) Serialize() string {
	var w *strings.Builder
	json.NewEncoder(w).Encode(r)
	return w.String()
}
