package utils

import (
	"encoding/json"
	"time"
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

func DoEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}