package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	bundleURL       = "bundle"
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

func (c *Client) CreateProject(ctx context.Context) (*models.Project, error) {
	var project models.Project
	if err := c.post(ctx, struct{}{}, &project, projectsURL); err != nil {
		return nil, err
	}
	return &project, nil
}

func (c *Client) CreateApplication(ctx context.Context, project string) (*models.Application, error) {
	var application models.Application
	if err := c.post(ctx, struct{}{}, &application, projectsURL, project, applicationsURL); err != nil {
		return nil, err
	}
	return &application, nil
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

func (c *Client) Execute(ctx context.Context, project, deviceID, command string) (*models.ExecuteResponse, error) {
	var executeResponse models.ExecuteResponse
	if err := c.post(ctx, command, &executeResponse, projectsURL, project, devicesURL, deviceID, executeURL); err != nil {
		return nil, err
	}
	return &executeResponse, nil
}

func (c *Client) InitiateSSH(ctx context.Context, project, deviceID string) (net.Conn, error) {
	req, err := http.NewRequest("", "", nil)
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
	req, err := http.NewRequest("GET", getURL(c.url, s...), nil)
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

	req, err := http.NewRequest("POST", getURL(c.url, s...), reader)
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
		return json.NewDecoder(resp.Body).Decode(&out)
	case http.StatusBadRequest, http.StatusNotFound:
		bytes, _ := ioutil.ReadAll(resp.Body)
		return errors.New(string(bytes))
	default:
		return errors.New(resp.Status)
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
