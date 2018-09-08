package utils

import (
	"encoding/json"
	"log"
	"os"
)

func Handle(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}

func ReadFile(path string) ([]byte, int) {
	file, err := os.Open(path)
	Handle(err)

	data := make([]byte, 2048)
	count, err := file.Read(data)
	Handle(err)

	return data, count
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
