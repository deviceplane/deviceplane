package conncontext

import (
	"context"
	"net"
	"net/http"
)

type key string

var ConnKey = key("http-conn")

func SaveConn(ctx context.Context, c net.Conn) context.Context {
	return context.WithValue(ctx, ConnKey, c)
}

func GetConn(r *http.Request) net.Conn {
	return r.Context().Value(ConnKey).(net.Conn)
}
