package search

import "encoding/json"

func decodeCategory(s string) []string {
	out := []string{}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}
