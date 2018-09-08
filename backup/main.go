package main

import (
	"gema/backup/health"
	"log"
)

func main() {
	log.Printf("Starting GEMA backup.")

	health.New(":80", "/")
}
