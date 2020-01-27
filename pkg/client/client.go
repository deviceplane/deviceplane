package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
	"github.com/function61/holepunch-server/pkg/wsconnadapter"
	"github.com/gorilla/websocket"
)

const (
	projectsURL     = "projects"
	applicationsURL = "applications"
	releasesURL     = "releases"
	devicesURL      = "devices"
	sshURL          = "ssh"
	executeURL      = "execute"
	rebootURL       = "reboot"
	bundleURL       = "bundle"
	metricsURL      = "metrics"
	servicesURL     = "services"
	membershipsURL  = "memberships"
)

type Client struct {
	url        *url.URL
	accessKey  string
	httpClient *http.Client
}

func NewClient(url *url.URL, accessKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		url:        url,
		accessKey:  accessKey,
		httpClient: httpClient,
	}
}

func (c *Client) CreateProject(ctx context.Context, name string) (*models.Project, error) {
	var project models.Project
	if err := c.post(ctx, models.Project{Name: name}, &project, projectsURL); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) CreateApplication(ctx context.Context, project string, name string) (*models.Application, error) {
	var application models.Application
	if err := c.post(ctx, models.Application{Name: name}, &application, projectsURL, project, applicationsURL); err != nil {
		return nil, err
	}
	return &application, nil
}

func (c *Client) ListProjects(ctx context.Context, project string) ([]models.ProjectFull, error) {
	var memberships []models.MembershipFull1
	if err := c.get(ctx, &memberships, membershipsURL+"?full"); err != nil {
		return nil, err
	}

	var projects []models.ProjectFull
	for _, m := range memberships {
		projects = append(projects, m.Project)
	}
	return projects, nil
}

func (c *Client) ListApplications(ctx context.Context, project string) ([]models.Application, error) {
	var applications []models.Application
	if err := c.get(ctx, &applications, projectsURL, project, applicationsURL); err != nil {
		return nil, err
	}
	return applications, nil
}

func (c *Client) ListDevices(ctx context.Context, filters []models.Filter, project string) ([]models.Device, error) {
	var devices []models.Device

	urlValues := url.Values{}
	for _, filter := range filters {
		bytes, err := json.Marshal(filter)
		if err != nil {
			return nil, err
		}

		b64Filter := base64.StdEncoding.EncodeToString(bytes)
		urlValues.Add("filter", b64Filter)
	}

	var queryString string
	if encoded := urlValues.Encode(); encoded != "" {
		queryString = "?" + encoded
	}

	if err := c.get(ctx, &devices, projectsURL, project, devicesURL+queryString); err != nil {
		return nil, err
	}
	return devices, nil
}

func (c *Client) GetApplication(ctx context.Context, project, application string) (*models.Application, error) {
	var app models.Application
	if err := c.get(ctx, &app, projectsURL, project, applicationsURL, application+"?full"); err != nil {
		return nil, err
	}
	return &app, nil
}

func (c *Client) GetDevice(ctx context.Context, project, device string) (*models.Device, error) {
	var d models.Device
	if err := c.get(ctx, &d, projectsURL, project, devicesURL, device+"?full"); err != nil {
		return nil, err
	}
	return &d, nil
}

func (c *Client) GetDeviceMetrics(ctx context.Context, project, device string) (*string, error) {
	var rawOpenMetrics string
	if err := c.get(ctx, &rawOpenMetrics, projectsURL, project, devicesURL, device, metricsURL, "host"); err != nil {
		return nil, err
	}
	return &rawOpenMetrics, nil
}

func (c *Client) GetServiceMetrics(ctx context.Context, project, device, application, service string) (*string, error) {
	var rawOpenMetrics string
	if err := c.get(ctx, &rawOpenMetrics, projectsURL, project, devicesURL, device, applicationsURL, application, servicesURL, service, metricsURL); err != nil {
		return nil, err
	}
	return &rawOpenMetrics, nil
}

func (c *Client) GetLatestRelease(ctx context.Context, project, application string) (*models.Release, error) {
	var release models.Release
	if err := c.get(ctx, &release, projectsURL, project, applicationsURL, application, releasesURL, "latest"); err != nil {
		return nil, err
	}
	return &release, nil
}

func (c *Client) CreateRelease(ctx context.Context, project, application, yamlConfig string) (*models.Release, error) {
	var release models.Release
	if err := c.post(ctx, models.CreateReleaseRequest{
		RawConfig: yamlConfig,
	}, &release, projectsURL, project, applicationsURL, application, releasesURL); err != nil {
		return nil, err
	}
	return &release, nil
}

func (c *Client) RebootDevice(ctx context.Context, project, device string) error {
	if err := c.post(ctx, []byte{}, nil, projectsURL, project, devicesURL, device, rebootURL); err != nil {
		return err
	}
	return nil
}

func (c *Client) Execute(ctx context.Context, project, deviceID, command string) (*models.ExecuteResponse, error) {
	var executeResponse models.ExecuteResponse
	if err := c.post(ctx, command, &executeResponse, projectsURL, project, devicesURL, deviceID, executeURL); err != nil {
		return nil, err
	}
	return &executeResponse, nil
}

func (c *Client) InitiateSSH(ctx context.Context, project, deviceID string) (net.Conn, error) {
	req, err := http.NewRequestWithContext(ctx, "", "", nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.accessKey, "")

	wsConn, _, err := websocket.DefaultDialer.Dial(getWebsocketURL(c.url, projectsURL, project, devicesURL, deviceID, sshURL), req.Header)
	if err != nil {
		return nil, err
	}

	return wsconnadapter.New(wsConn), nil
}

func (c *Client) get(ctx context.Context, out interface{}, s ...string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", getURL(c.url, s...), nil)
	if err != nil {
		return err
	}

	return c.performRequest(req, out)
}

func (c *Client) post(ctx context.Context, in, out interface{}, s ...string) error {
	var reqBytes []byte

	switch v := in.(type) {
	case string:
		reqBytes = []byte(v)
	default:
		var err error
		reqBytes, err = json.Marshal(in)
		if err != nil {
			return err
		}
	}

	reader := bytes.NewReader(reqBytes)

	req, err := http.NewRequestWithContext(ctx, "POST", getURL(c.url, s...), reader)
	if err != nil {
		return err
	}

	return c.performRequest(req, out)
}

func (c *Client) performRequest(req *http.Request, out interface{}) error {
	req.SetBasicAuth(c.accessKey, "")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return c.handleResponse(resp, out)
}

func (c *Client) handleResponse(resp *http.Response, out interface{}) error {
	switch resp.StatusCode {
	case http.StatusOK:
		switch o := out.(type) {
		case *string:
			bytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			*o = string(bytes)
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(&out)
	case http.StatusBadRequest, http.StatusNotFound:
		bytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bytes))
	default:
		return fmt.Errorf("%d %s", resp.StatusCode, resp.Status)
	}
}

func getURL(u *url.URL, s ...string) string {
	return strings.Join(append([]string{u.String()}, s...), "/")
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
