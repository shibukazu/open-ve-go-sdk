package client

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

type CheckRequest struct {
	Validations []CheckRequestValidation `json:"validations"`
}

type CheckRequestValidation struct {
	ID        string                 `json:"id"`
	Variables map[string]interface{} `json:"variables"`
}

type CheckResponse struct {
	Results []CheckResponseResult `json:"results"`
}

type CheckResponseResult struct {
	ID      string `json:"id"`
	IsValid bool   `json:"isValid"`
	Message string `json:"message"`
}

type Config struct {
	Url     string
	Version string
}

func VersionOption(version string) func(*Config) {
	return func(c *Config) {
		c.Version = version
	}
}

func NewConfig(url string, options ...func(*Config)) *Config {
	config := &Config{
		Url:     url,
		Version: "v1",
	}

	for _, option := range options {
		option(config)
	}

	return config
}

type Client struct {
	config *Config
}

func NewClient(config *Config) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) Check(checkRequest *CheckRequest) (*CheckResponse, error) {
	requestBody, err := json.Marshal(checkRequest)
	if err != nil {
		return nil, err
	}

	url, err := url.JoinPath(c.config.Url, c.config.Version, "check")
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var checkResponse CheckResponse
	err = json.Unmarshal(body, &checkResponse)
	if err != nil {
		return nil, err
	}

	return &checkResponse, nil
}
