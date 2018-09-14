package health

import (
	"fmt"
	"log"
	"net/http"
)

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK\n")
}

func New(addr string, path string) {
	http.HandleFunc(path, health)

	log.Printf("Starting health responder at %s", addr)
	http.ListenAndServe(addr, nil)
}
