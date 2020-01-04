package client

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func GetAgentMetrics(deviceConn net.Conn) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		"/metrics/agent",
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetDeviceMetrics(deviceConn net.Conn) (*http.Response, error) {
	req, _ := http.NewRequest(
		"GET",
		"/metrics/host",
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetServiceMetrics(deviceConn net.Conn, applicationID, service string, metricPath string, metricPort uint) (*http.Response, error) {
	serviceURL := url.URL{
		Path: fmt.Sprintf(
			"/applications/%s/services/%s/metrics",
			applicationID, service,
		),
	}

	query := serviceURL.Query()
	query.Set("path", base64.RawURLEncoding.EncodeToString([]byte(metricPath)))
	query.Set("port", strconv.Itoa(int(metricPort)))
	serviceURL.RawQuery = query.Encode()

	req, _ := http.NewRequest(
		"GET",
		serviceURL.RequestURI(),
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

func InitiateReboot(deviceConn net.Conn) (*http.Response, error) {
	req, _ := http.NewRequest(
		"POST",
		"/reboot",
		nil,
	)

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}
