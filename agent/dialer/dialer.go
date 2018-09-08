package dialer

import (
	"bytes"
	"fmt"
	"gema/agent/event"
	"gema/agent/utils"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
)

var jsonRegex = regexp.MustCompile(`{.+}`)

type Filters struct {
	Label []string `json:"label"`
}

type Dialer struct {
	Addr    string
	Handler *event.Handler
}

func New(addr string) *Dialer {
	return &Dialer{
		Addr:    addr,
		Handler: event.New(),
	}
}

func (d *Dialer) MonitorEvents(filters Filters) {
	log.Printf("Starting event monitoring of swarm.")

	c := d.dial()
	defer c.Close()

	r := makeRequest("GET", "http:/events", filters)

	writeToConn(r, c)
	d.readFromConn(c)
}

func (d *Dialer) readFromConn(r io.Reader) {
	buf := make([]byte, 2048)

	for {
		n, err := r.Read(buf[:])
		utils.Handle(err)

		match := jsonRegex.FindStringSubmatch(string(buf[0:n]))

		if len(match) == 0 {
			continue
		}

		j := match[0]

		ev := event.DefaultEvent()
		utils.FromJSON([]byte(j), ev)

		go d.Handler.HandleEvent(ev)
	}
}

func (d *Dialer) dial() net.Conn {
	log.Printf("Dialing %s", d.Addr)

	c, err := net.Dial("unix", d.Addr)
	utils.Handle(err)
	return c
}

func makeRequest(method string, Url string, filters Filters) http.Request {
	j := string(utils.ToJSON(filters))

	u, _ := url.Parse(fmt.Sprintf("%s?filters=%s", Url, j))

	r := http.Request{
		Method: method,
		URL:    u,
	}
	return r
}

func writeToConn(r http.Request, c net.Conn) {
	log.Printf("Writing to API: %s", r.URL.String())

	var buf bytes.Buffer
	r.Write(&buf)

	_, err := c.Write(buf.Bytes())
	utils.Handle(err)
}
