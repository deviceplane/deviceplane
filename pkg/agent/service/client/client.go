package client

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

func GetDeviceMetrics(deviceConn net.Conn) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		"/metrics",
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetServiceMetrics(deviceConn net.Conn, applicationID, service string) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"/applications/%s/services/%s/metrics",
			applicationID, service,
		),
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetImagePullProgress(deviceConn net.Conn, applicationID, service string) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"/applications/%s/services/%s/imagepullprogress",
			applicationID, service,
		),
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func InitiateSSH(deviceConn net.Conn) error {
	req, _ := http.NewRequest("POST", "/ssh", nil)
	return req.Write(deviceConn)
}
