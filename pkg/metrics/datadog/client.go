package datadog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deviceplane/deviceplane/pkg/models"
)

func NewPoint(value float32) [2]interface{} {
	return [2]interface{}{
		time.Now().Unix(),
		value,
	}
}

type Client struct {
	apiKey string
}

func NewClient(apiKey string) *Client {
	return &Client{
		apiKey: apiKey,
	}
}

func (c *Client) PostMetrics(ctx context.Context, req models.DatadogPostMetricsRequest) error {
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
