package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/deviceplane/deviceplane/pkg/models"
)

const (
	bundleURL = "bundle"
)

type Client struct {
	url        string
	projectID  string
	httpClient *http.Client

	deviceID  string
	accessKey string
}

func NewClient(url, projectID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
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

	req, err := http.NewRequest("POST", c.getURL("projects", c.projectID, "devices", "register"), reader)
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

func (c *Client) SetDeviceInfo(ctx context.Context, info models.SetDeviceInfoRequest) error {
	return c.post(ctx, info, nil, "projects", c.projectID, "devices", c.deviceID, "info")
}

func (c *Client) get(ctx context.Context, out interface{}, s ...string) error {
	req, err := http.NewRequest("GET", c.getURL(s...), nil)
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

	req, err := http.NewRequest("POST", c.getURL(s...), reader)
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

func (c *Client) getURL(s ...string) string {
	return strings.Join(append([]string{c.url}, s...), "/")
}
