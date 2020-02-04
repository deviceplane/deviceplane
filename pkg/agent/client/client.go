package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/websocket"
)

const (
	bundleURL = "bundle"
)

type Client struct {
	url        *url.URL
	projectID  string
	httpClient *http.Client

	deviceID  string
	accessKey string
}

func NewClient(url *url.URL, projectID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		url:        url,
		projectID:  projectID,
		httpClient: httpClient,
	}
}

func (c *Client) SetDeviceID(deviceID string) {
	c.deviceID = deviceID
}

func (c *Client) SetAccessKey(accessKey string) {
	c.accessKey = accessKey
}

func (c *Client) RegisterDevice(ctx context.Context, registrationToken string) (*models.RegisterDeviceResponse, error) {
	reqBytes, err := json.Marshal(models.RegisterDeviceRequest{
		DeviceRegistrationTokenID: registrationToken,
	})
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(reqBytes)

	req, err := http.NewRequest("POST", getURL(c.url, "projects", c.projectID, "devices", "register"), reader)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"status": resp.Status,
		"code":   resp.StatusCode,
		"body":   string(bytes),
	}).Debug("POST response")

	var registerDeviceResponse models.RegisterDeviceResponse
	if err := json.Unmarshal(bytes, &registerDeviceResponse); err != nil {
		return nil, err
	}

	return &registerDeviceResponse, nil
}

func (c *Client) GetBundle(ctx context.Context) (*models.Bundle, error) {
	var bundle models.Bundle
	if err := c.get(ctx, &bundle, "projects", c.projectID, "devices", c.deviceID, "bundle"); err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (c *Client) SetDeviceInfo(ctx context.Context, req models.SetDeviceInfoRequest) error {
	return c.post(ctx, req, nil, "projects", c.projectID, "devices", c.deviceID, "info")
}

func (c *Client) SendDeviceMetrics(ctx context.Context, req models.DatadogPostMetricsRequest) error {
	return c.post(ctx, req, nil, "projects", c.projectID, "devices", c.deviceID, "forwardmetrics", "device")
}

func (c *Client) SendServiceMetrics(ctx context.Context, req models.IntermediateServiceMetricsRequest) error {
	return c.post(ctx, req, nil, "projects", c.projectID, "devices", c.deviceID, "forwardmetrics", "service")
}

func (c *Client) SetDeviceApplicationStatus(ctx context.Context, applicationID string, req models.SetDeviceApplicationStatusRequest) error {
	return c.post(ctx, req, nil, "projects", c.projectID, "devices", c.deviceID, "applications", applicationID, "deviceapplicationstatuses")
}

func (c *Client) DeleteDeviceApplicationStatus(ctx context.Context, applicationID string) error {
	return c.delete(ctx, nil, "projects", c.projectID, "devices", c.deviceID, "applications", applicationID, "deviceapplicationstatuses")
}

func (c *Client) SetDeviceServiceStatus(ctx context.Context, applicationID, service string, req models.SetDeviceServiceStatusRequest) error {
	return c.post(ctx, req, nil, "projects", c.projectID, "devices", c.deviceID, "applications", applicationID, "services", service, "deviceservicestatuses")
}

func (c *Client) DeleteDeviceServiceStatus(ctx context.Context, applicationID, service string) error {
	return c.delete(ctx, nil, "projects", c.projectID, "devices", c.deviceID, "applications", applicationID, "services", service, "deviceservicestatuses")
}

func (c *Client) InitiateDeviceConnection(ctx context.Context) (net.Conn, error) {
	req, err := http.NewRequest("", "", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.accessKey, "")

	wsConn, _, err := websocket.DefaultDialer.Dial(getWebsocketURL(c.url, "projects", c.projectID, "devices", c.deviceID, "connection"), req.Header)
	if err != nil {
		return nil, err
	}

	return wsconnadapter.New(wsConn), nil
}

func (c *Client) Revdial(ctx context.Context, path string) (*websocket.Conn, *http.Response, error) {
	return websocket.DefaultDialer.Dial(getWebsocketURL(c.url, strings.TrimPrefix(path, "/")), nil)
}

func (c *Client) get(ctx context.Context, out interface{}, s ...string) error {
	req, err := http.NewRequest("GET", getURL(c.url, s...), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.accessKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"status": resp.Status,
		"code":   resp.StatusCode,
		"body":   string(bytes),
	}).Debug("GET response")

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &out)
}

func (c *Client) post(ctx context.Context, in, out interface{}, s ...string) error {
	reqBytes, err := json.Marshal(in)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(reqBytes)

	req, err := http.NewRequest("POST", getURL(c.url, s...), reader)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.accessKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"status": resp.Status,
		"code":   resp.StatusCode,
		"body":   string(bytes),
	}).Debug("POST response")

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &out)
}

func (c *Client) delete(ctx context.Context, out interface{}, s ...string) error {
	req, err := http.NewRequest("DELETE", getURL(c.url, s...), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(c.accessKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"status": resp.Status,
		"code":   resp.StatusCode,
		"body":   string(bytes),
	}).Debug("DELETE response")

	if len(bytes) == 0 {
		return nil
	}

	return json.Unmarshal(bytes, &out)
}

func getURL(url *url.URL, s ...string) string {
	return strings.Join(append([]string{url.String()}, s...), "/")
}

func getWebsocketURL(u *url.URL, s ...string) string {
	uCopy, _ := url.Parse(u.String())
	switch uCopy.Scheme {
	case "http":
		uCopy.Scheme = "ws"
	default:
		uCopy.Scheme = "wss"
	}
	return strings.Join(append([]string{uCopy.String()}, s...), "/")
}
