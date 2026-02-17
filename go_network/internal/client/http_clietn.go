package client

import (
    "fmt"
    "io"
    "net/http"
    "time"
)

type HTTPClient struct {
    client *http.Client
    baseURL string
}

func NewHTTPClient(baseURL string) *HTTPClient {
    return &HTTPClient{
        client: &http.Client{
            Timeout: 5 * time.Second,
        },
        baseURL: baseURL,
    }
}

func (c *HTTPClient) Get(path string) (string, error) {
    url := fmt.Sprintf("%s%s", c.baseURL, path)
    
    resp, err := c.client.Get(url)
    if err != nil {
        return "", fmt.Errorf("request failed: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("unexpected status: %s", resp.Status)
    }

    return string(body), nil
}

func (c *HTTPClient) Ping() error {
    _, err := c.Get("/ping")
    return err
}