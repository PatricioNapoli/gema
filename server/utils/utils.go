package utils

import (
	"encoding/json"
)

func Handle(err error) {
	if err != nil {
		panic(err)
	}
}

func ToJSON(v interface{}) []byte {
	j, err := json.Marshal(v)
	Handle(err)
	return j
}

func FromJSON(d []byte, v interface{}) {
	err := json.Unmarshal(d, v)
	Handle(err)
}
