package service

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (s *Service) connectTCP(w http.ResponseWriter, r *http.Request) {
	withPort(w, r, func(port int) {
		conn := conncontext.GetConn(r)

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			panic(err)
		}

		fmt.Println(port)

		go io.Copy(localConn, conn)
		io.Copy(conn, localConn)
	})
}

func (s *Service) connectHTTP(w http.ResponseWriter, r *http.Request) {
	withPort(w, r, func(port int) {
		conn := conncontext.GetConn(r)

		serverConn := httputil.NewServerConn(conn, nil)

		req, err := serverConn.Read()
		if err != nil {
			panic(err)
		}

		url, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))
		if err != nil {
			panic(err)
		}

		req.RequestURI = ""
		req.URL = url

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		utils.ProxyResponse(w, resp)
	})
}
