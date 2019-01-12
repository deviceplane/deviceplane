package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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
	url        string
	httpClient *http.Client
}

func NewClient(url string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	return &Client{
		url:        url,
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

func (c *Client) CreateApplication(ctx context.Context, projectID string) (*models.Application, error) {
	var application models.Application
	if err := c.post(ctx, struct{}{}, &application, projectID, applicationsURL); err != nil {
		return nil, err
	}
	return &application, nil
}

func (c *Client) CreateRelease(ctx context.Context, projectID, applicationID, config string) (*models.Release, error) {
	var release models.Release
	if err := c.post(ctx, models.CreateRelease{
		Config: config,
	}, &release, projectID, applicationsURL, applicationID, releasesURL); err != nil {
		return nil, err
	}
	return &release, nil
}

func (c *Client) GetBundle(ctx context.Context, projectID string) (*models.Bundle, error) {
	var bundle models.Bundle
	if err := c.get(ctx, &bundle, projectID, bundleURL); err != nil {
		return nil, err
	}
	return &bundle, nil
}

func (c *Client) get(ctx context.Context, out interface{}, s ...string) error {
	resp, err := c.httpClient.Get(c.getURL(s...))
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

	resp, err := c.httpClient.Post(c.getURL(s...), "application/json", reader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(&out)
}

func (c *Client) getURL(s ...string) string {
	return strings.Join(append([]string{c.url}, s...), "/")
}
