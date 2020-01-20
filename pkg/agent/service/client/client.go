package client

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strconv"
)

func GetAgentMetrics(ctx context.Context, deviceConn net.Conn) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"/metrics/agent",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetDeviceMetrics(ctx context.Context, deviceConn net.Conn) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		"/metrics/host",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetServiceMetrics(ctx context.Context, deviceConn net.Conn, applicationID, service string, metricPath string, metricPort uint) (*http.Response, error) {
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

	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		serviceURL.RequestURI(),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func GetImagePullProgress(ctx context.Context, deviceConn net.Conn, applicationID, service string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf(
			"/applications/%s/services/%s/imagepullprogress",
			applicationID, service,
		),
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}

func InitiateSSH(ctx context.Context, deviceConn net.Conn) error {
	req, err := http.NewRequestWithContext(ctx, "POST", "/ssh", nil)
	if err != nil {
		return err
	}
	return req.Write(deviceConn)
}

func InitiateReboot(ctx context.Context, deviceConn net.Conn) (*http.Response, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		"/reboot",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if err := req.Write(deviceConn); err != nil {
		return nil, err
	}

	return http.ReadResponse(bufio.NewReader(deviceConn), req)
}
