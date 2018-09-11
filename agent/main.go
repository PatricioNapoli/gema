package main

import (
	"gema/agent/dialer"
	"gema/agent/health"
	"log"
)

func main() {
	log.Printf("Starting GEMA agent.")

	go health.New(":80", "/health")

	d := dialer.New("/var/run/docker.sock")

	d.MonitorEvents(dialer.Filters{
		Label: []string{"gema.service"},
	})
}
