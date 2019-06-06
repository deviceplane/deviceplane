package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/deviceplane/deviceplane/pkg/models"
)

const (
	projectsURL     = "projects"
	applicationsURL = "applications"
	releasesURL     = "releases"
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

func (c *Client) CreateRelease(ctx context.Context, project, application, config string) (*models.Release, error) {
	var release models.Release
	if err := c.post(ctx, models.CreateRelease{
		Config: config,
	}, &release, projectsURL, project, applicationsURL, application, releasesURL); err != nil {
		return nil, err
	}
	return &release, nil
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

	return json.NewDecoder(resp.Body).Decode(&out)
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

	return json.NewDecoder(resp.Body).Decode(&out)
}

func (c *Client) getURL(s ...string) string {
	return strings.Join(append([]string{c.url.String()}, s...), "/")
}
