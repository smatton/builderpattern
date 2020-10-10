package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

// type project struct {
// 	path string
// }

// type Project interface {
// 	Put()
// 	GetXml()
// 	GetJson()
// 	Delete()
// 	FromXmlFile()
// 	FromBytes()
// }

type requestBuilder struct {
	context context.Context
	path    string
	options []string
	method  string
}

type request struct {
	req     *http.Request
	path    string
	options string
	method  string
}

func (rb *requestBuilder) Context(context context.Context) RequestBuilder {
	rb.context = context
	return rb
}

func (rb *requestBuilder) Path(path string) RequestBuilder {
	rb.path = path
	return rb
}

func (rb *requestBuilder) Option(option string) RequestBuilder {
	rb.options = append(rb.options, option)
	return rb
}

func (rb *requestBuilder) Method(method string) RequestBuilder {
	rb.method = method
	return rb
}

func (rb *requestBuilder) Build() Request {

	opts := strings.Join(rb.options, "&")

	paths := strings.Builder{}
	paths.WriteString("?")
	paths.WriteString(opts)
	return &request{
		path:    paths.String(),
		options: opts,
		method:  rb.method,
	}
}

type Request interface {
	New(string) *http.Request
	Print()
}

func (r *request) New(url string) *http.Request {
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(url)
	urlBuilder.WriteString(r.path)

	req, _ := http.NewRequest(r.method, urlBuilder.String(), nil)
	return req
}

func (r *request) Print() {
	fmt.Println(r.path)
}

type RequestBuilder interface {
	Context(context.Context) RequestBuilder
	Path(string) RequestBuilder
	Option(string) RequestBuilder
	Build() Request
}

func New() ClientBuilder {
	return &clientBuilder{}
}

func NewRequest() RequestBuilder {
	return &requestBuilder{}
}

func main() {

	builder := New()

	client := builder.Host("localhost").
		Port(9910).
		Build()

	client.Print()

	request := NewRequest()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	newrequest := request.Context(ctx).Path("/projects/ppsd").
		Option("start=true").
		Option("format=csv").
		Build()

	newrequest.Print()
}
