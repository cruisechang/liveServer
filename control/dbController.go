package control

import (
	"io"
	"net/http"
	"time"
)

//BasicProcessor is parent struct for process.
type DBController struct {
	dbAPIHost string
}

func NewDBController(dbAPIHost string) *DBController {

	return &DBController{
		dbAPIHost: dbAPIHost,
	}

}

//db get/patch hall, get/patch room  , get/patch user
func (c *DBController) SetDBAPIHost(host string) {
	c.dbAPIHost = host
}

func (c *DBController) Do(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, c.dbAPIHost+path, body)
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	req.Header.Set("API-Key", "qwerASDFzxcv!@#$")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout:timeout,
	}
	return  client.Do(req)
}

func (c *DBController) DoWithFullPath(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	timeout := time.Duration(5 * time.Second)
	client := &http.Client{
		Timeout:timeout,
	}
	return client.Do(req)
}
