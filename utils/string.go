package utils

import (
	"encoding/json"
	"strings"
)

type StringArray string

func (r StringArray) MarshalJSON() ([]byte, error) {
	items := []string{}
	if string(r) != "" {
		items = strings.Split(string(r), ",")
	}
	for _, item := range items {
		item = strings.TrimSpace(item)
	}
	return json.Marshal(items)
}
