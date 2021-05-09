package service

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"

	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (s *Service) connectTCP(w http.ResponseWriter, r *http.Request) {
	withPort(w, r, func(port int) {
		conn := conncontext.GetConn(r)

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

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
			println(err.Error())
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		req.RequestURI = ""
		req.URL.Scheme = "http"
		req.URL.Host = fmt.Sprintf("localhost:%d", port)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			println(err.Error())
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponse(w, resp)
	})
}
