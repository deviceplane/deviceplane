package datadog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PostMetricsRequest struct {
	Series Series `json:"series"`
}

type Series []Metric

type Metric struct {
	Metric   string       `json:"metric"`
	Points   [][2]float32 `json:"points"`
	Type     string       `json:"type"`
	Interval *int64       `json:"interval,omitempty"`
	Host     string       `json:"host,omitempty"`
	Tags     []string     `json:"tags"`
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) PostMetrics(ctx context.Context, req PostMetricsRequest) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(
		fmt.Sprintf("https://api.datadoghq.com/api/v1/series?api_key=%s", c.apiKey),
		"application/json", bytes.NewBuffer(reqBytes),
	)
	if err != nil {
		return err
	}

	return resp.Body.Close()
}
