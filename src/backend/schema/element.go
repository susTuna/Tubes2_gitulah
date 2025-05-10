package schema

import (
	"encoding/json"
	"strings"
)

type Element struct {
	ID       int32
	Name     string
	Tier     int32
	ImageUrl string
}

func (e Element) Serialize() string {
	var w *strings.Builder
	json.NewEncoder(w).Encode(e)
	return w.String()
}
