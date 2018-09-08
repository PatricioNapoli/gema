package health

import (
	"fmt"
	"gema/backup/utils"
	"log"
	"net/http"
)

type Health struct {
	Status string
}

func health(w http.ResponseWriter, r *http.Request) {
	status := &Health{
		Status: "OK",
	}

	statusResponse := string(utils.ToJSON(status))
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, statusResponse)
}

func New(addr string, path string) {
	http.HandleFunc(path, health)

	log.Printf("Starting health responder at %s", addr)
	http.ListenAndServe(addr, nil)
}
