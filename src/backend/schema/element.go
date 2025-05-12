package schema

import (
	"encoding/json"
	"strings"
)

type Element struct {
	ID       int32  `json:"id"`
	Name     string `json:"name"`
	Tier     int32  `json:"tier"`
	ImageUrl string `json:"image_url"`
}

func (e Element) Serialize() string {
	var w strings.Builder
	json.NewEncoder(&w).Encode(e)
	return w.String()
}
