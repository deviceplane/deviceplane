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

	ErrNonSuccessResponse = errors.New("non-2xx status code")
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, errors.WithMessagef(ErrNonSuccessResponse, "code: %d, body: %s", resp.StatusCode, string(body))
	}

	return &Response{
		Response: resp,
	}, nil
}
