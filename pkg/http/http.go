package http

import (
	"io"
	"net/http"

	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
)

var (
	DefaultClient = &Client{
		Client: http.DefaultClient,
	}
)

type Request struct {
	*http.Request
}

func NewRequest(ctx *dpcontext.Context, method, url string, body io.Reader) (*Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	return &Request{
		Request: req,
	}, nil
}

type Response struct {
	*http.Response
}

type Client struct {
	*http.Client
}

func (c *Client) Do(req *Request) (*Response, error) {
	resp, err := c.Client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	return &Response{
		Response: resp,
	}, nil
}
