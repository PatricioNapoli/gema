package main

import (
	"gema/agent/health"
	"net"
	"io"
	"net/http"
	"net/url"
	"bytes"
)

func reader(r io.Reader) {
	buf := make([]byte, 2048)
	for {
		n, err := r.Read(buf[:])
		if err != nil {
			panic(err)
		}
		println(string(buf[0:n]))
	}
}

func main() {
	go health.New()

	c, err := net.Dial("unix", "/var/run/docker.sock")
	if err != nil {
		panic(err)
	}
	defer c.Close()

	u, _:= url.Parse("http:/events")
	r := http.Request{
		Method: "GET",
		Host: "v1.24",
		URL: u,
	}

	var buf bytes.Buffer
	r.Write(&buf)

	_, err = c.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}

	reader(c)
}