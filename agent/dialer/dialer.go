package dialer

import (
	"net"
	"gema/agent/utils"
	"net/url"
	"net/http"
	"fmt"
	"encoding/json"
	"bytes"
	"regexp"
	"io"
	"gema/agent/event"
)

var jsonRegex = regexp.MustCompile(`{.+}`)

type Filters struct {
	Label []string `json:"label"`
}

type Dialer struct {
	Addr string
	Handler *event.Handler
}

func New(addr string) *Dialer {
	return &Dialer{
		Addr: addr,
		Handler: event.New(),
	}
}

func (d *Dialer) MonitorEvents(filters Filters) {
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

		ev := &event.Event{}
		err = json.Unmarshal([]byte(j), ev)
		utils.Handle(err)

		go d.Handler.HandleEvent(ev)
	}
}

func (d *Dialer) dial() net.Conn {
	c, err := net.Dial("unix", d.Addr)
	utils.Handle(err)
	return c
}

func makeRequest(method string, Url string, filters Filters) http.Request {
	j, err := json.Marshal(filters)
	utils.Handle(err)

	js := string(j)

	u, _ := url.Parse(fmt.Sprintf("%s?filters=%s", Url, js))

	r := http.Request{
		Method: method,
		URL: u,

	}
	return r
}

func writeToConn(r http.Request, c net.Conn) {
	var buf bytes.Buffer
	r.Write(&buf)

	_, err := c.Write(buf.Bytes())
	utils.Handle(err)
}