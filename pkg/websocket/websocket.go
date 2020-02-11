package websocket

import (
	"net/http"

	dpcontext "github.com/deviceplane/deviceplane/pkg/context"
	dphttp "github.com/deviceplane/deviceplane/pkg/http"
	"github.com/gorilla/websocket"
)

var (
	DefaultDialer = &Dialer{
		Dialer: websocket.DefaultDialer,
	}
)

type Conn struct {
	*websocket.Conn
}

type Dialer struct {
	*websocket.Dialer
}

func (d *Dialer) Dial(ctx *dpcontext.Context, urlStr string, requestHeader http.Header) (*Conn, *dphttp.Response, error) {
	conn, resp, err := d.Dialer.DialContext(ctx, urlStr, requestHeader)
	if err != nil {
		return nil, nil, err
	}
	return &Conn{
			Conn: conn,
		}, &dphttp.Response{
			Response: resp,
		}, nil
}
