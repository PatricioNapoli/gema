package health

import (
	"net/http"
	"fmt"
	"encoding/json"
)

type Health struct {
	Status string
}

func health(w http.ResponseWriter, r *http.Request) {
	status := &Health{
		Status: "OK",
	}
	b, err := json.Marshal(status)
	if err != nil {
		panic(err)
	}

	statusResponse := string(b)
	print(statusResponse)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, statusResponse)
}

func New() {
	http.HandleFunc("/gema/health", health)
	http.ListenAndServe(":8081", nil)
}