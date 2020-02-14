package http

import (
	"io"
	"io/ioutil"
	"net/http"

	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	"github.com/pkg/errors"
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

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &Response{
			Response: resp,
		}, nil
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = errors.New(string(body))
	}

	return nil, errors.Wrapf(err, "status code %d", resp.StatusCode)
}

func (c *Client) Get(ctx *dpcontext.Context, url string) (*Response, error) {
	req, err := NewRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func Get(ctx *dpcontext.Context, url string) (*Response, error) {
	return DefaultClient.Get(ctx, url)
}
