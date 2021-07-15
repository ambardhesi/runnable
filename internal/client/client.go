package client

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/ambardhesi/runnable/internal/server"
	"github.com/go-resty/resty/v2"
)

type Config struct {
	ServerAddress  string
	CertFilePath   string
	KeyFilePath    string
	CaCertFilePath string
}

type Client struct {
	Config     Config
	HttpClient *resty.Client
}

func NewClient(config Config) (*Client, error) {
	rClient := resty.New()
	tlsConfig, err := GetTLSConfig(config.CertFilePath,
		config.KeyFilePath, config.CaCertFilePath)

	if err != nil {
		fmt.Printf("Failed to create http client %v\n", err)
		return nil, err
	}
	rClient.SetTLSClientConfig(tlsConfig)

	return &Client{
		HttpClient: rClient,
		Config:     config,
	}, nil
}

func (c *Client) StartJob(cmd string) (*string, error) {
	request := server.StartJobRequest{
		Command: cmd,
	}

	resp, err := c.HttpClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(request).
		Post(c.Config.ServerAddress + "/job")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil

}

func (c *Client) StopJob(jobID string) error {
	resp, err := c.HttpClient.R().
		Post(c.Config.ServerAddress + "/job/" + jobID + "/stop")

	if err != nil {
		return err
	}

	if resp.StatusCode() != http.StatusOK {
		return errors.New(string(resp.Body()))
	}

	return nil
}

func (c *Client) Get(jobID string) (*string, error) {
	resp, err := c.HttpClient.R().
		Get(c.Config.ServerAddress + "/job/" + jobID)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}

func (c *Client) GetLogs(jobID string) (*string, error) {
	resp, err := c.HttpClient.R().
		Get(c.Config.ServerAddress + "/job/" + jobID + "/logs")

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(resp.Body()))
	}

	body := string(resp.Body())
	return &body, nil
}
