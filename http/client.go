package http

import (
	"io"
	"net/http"
	"time"
)

type Client struct {
	address    string
	port       string
	targetAddr string
	client     *http.Client
	postURI    string
	postQuery  string
}

func NewClient(address, port string, connectTimeout, handshakeTimeout, requestTimeout int) (*Client, error) {
	c := &Client{}
	c.SetTargetAddress(address, port)

	var netTransport = &http.Transport{
		//Dial: (&net.Dialer{
		//	Timeout: connectTimeout * time.Second,
		//}).Dial,
		TLSHandshakeTimeout: time.Duration(handshakeTimeout) * time.Second,
	}

	c.client = &http.Client{
		Timeout:   time.Duration(requestTimeout) * time.Second,
		Transport: netTransport,
	}

	return c, nil
}

func (c *Client) SetTargetAddress(addr, port string) {
	c.address = addr
	c.port = port
	c.targetAddr = "http://" + addr + ":" + port
}

func (c *Client) Do(method, path string, body io.Reader) (*http.Response, error) {

	req, err := http.NewRequest(method, c.targetAddr+path, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	return c.client.Do(req)
}

func (c *Client) DoWithFullPath(method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	return c.client.Do(req)
}

func DoWithFullPath(method, path string, body io.Reader, timeoutSecond int) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {

	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")

	timeout := time.Second * time.Duration(timeoutSecond)
	client := &http.Client{
		Timeout: timeout,
	}
	return client.Do(req)
}
