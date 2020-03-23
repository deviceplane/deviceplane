package service

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/deviceplane/deviceplane/pkg/agent/server/conncontext"
	"github.com/deviceplane/deviceplane/pkg/codes"
	"github.com/deviceplane/deviceplane/pkg/utils"
)

func (s *Service) connectTCP(w http.ResponseWriter, r *http.Request) {
	withPort(w, r, func(port int) {
		conn := conncontext.GetConn(r)

		localConn, err := net.Dial("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
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
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		url, err := url.Parse(fmt.Sprintf("http://localhost:%d", port))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req.RequestURI = ""
		req.URL = url

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), codes.StatusDeviceConnectionFailure)
			return
		}

		utils.ProxyResponse(w, resp)
	})
}
