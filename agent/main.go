package main

import (
	"gema/agent/health"
	"gema/agent/dialer"
)


func main() {
	go health.New()

	d := dialer.New("/var/run/docker.sock")

	d.MonitorEvents(dialer.Filters{
		Label: []string{"gema.expose"},
	})
}