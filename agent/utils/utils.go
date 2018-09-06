package utils

import "os"

func Handle(err error) {
	if err != nil {
		panic(err)
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