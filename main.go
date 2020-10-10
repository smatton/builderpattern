package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type clientBuilder struct {
	host string
	port string
}

func (cb *clientBuilder) Host(host string) ClientBuilder {
	cb.host = host
	return cb
}

func (cb *clientBuilder) Port(port int) ClientBuilder {
	cb.port = strconv.Itoa(port)
	return cb
}

func (cb *clientBuilder) Build() Client {
	stringBuilder := strings.Builder{}
	stringBuilder.WriteString(cb.host)
	stringBuilder.WriteString(":")
	stringBuilder.WriteString(cb.port)
	return &client{
		httpclient: &http.Client{},
		host:       cb.host,
		port:       cb.port,
		addr: &url.URL{
			Scheme: "http",
			Host:   stringBuilder.String(),
		},
	}
}

type ClientBuilder interface {
	Host(string) ClientBuilder
	Port(int) ClientBuilder
	Build() Client
}

type client struct {
	httpclient *http.Client
	host       string
	port       string
	addr       *url.URL
}

func (c *client) Print() {
	fmt.Printf("Hostname %s\n", c.addr.String())

}
func (c *client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = req.WithContext(ctx)

	resp, err := c.httpclient.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:

		}
		return nil, err
	}

	return resp, err

}

//Client does client
type Client interface {
	Do(context.Context, *http.Request, interface{}) (*http.Response, error)
	Print()
}

func New() ClientBuilder {
	return &clientBuilder{}
}

func main() {

	builder := New()

	client := builder.Host("localhost").
		Port(9910).
		Build()

	client.Print()
}
