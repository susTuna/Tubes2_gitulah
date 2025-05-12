package schema

type SearchRequest struct {
	Element   int    `json:"element"`
	Method    string `json:"method"`
	Count     int    `json:"count"`
	Delay     int    `json:"delay"`
	Threading string `json:"threading"`
}

func (sr *SearchRequest) Valid() bool {
	return ((sr.Method == "dfs" || sr.Method == "bfs") &&
		(sr.Count > 0) &&
		(sr.Delay >= 0) &&
		(sr.Threading == "single" || sr.Threading == "multi"))
}
