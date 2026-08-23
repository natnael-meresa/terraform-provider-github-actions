package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	baseURL    = "https://api.github.com"
	apiVersion = "2026-03-10"
)

type Client struct {
	token      string
	httpClient *http.Client
	baseURL    string
}

func NewClient(token string) *Client {
	return &Client{
		token:      token,
		httpClient: &http.Client{},
		baseURL:    baseURL,
	}
}

type Variable struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type APIError struct {
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return fmt.Sprintf(
		"Github API returned HTTP %d: %s",
		e.StatusCode,
		e.Message,
	)
}

func (c *Client) newRequest(
	ctx context.Context,
	method string,
	path string,
	body any,
) (*http.Request, error) {
	var bodyReader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}

		bodyReader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		method,
		c.baseURL+path,
		bodyReader,
	)

	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("X-Github-API-Version", apiVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func handleResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	defer resp.Body.Close()

	var githubError struct {
		Message string `json:"message"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubError); err != nil {
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
	}

	return &APIError{
		StatusCode: resp.StatusCode,
		Message:    githubError.Message,
	}
}

func (c *Client) CreateVariable(
	ctx context.Context,
	owner string,
	repo string,
	variable Variable,
) error {
	path := fmt.Sprintf(
		"/repos/%s/%s/actions/variables",
		owner,
		repo,
	)

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		path,
		variable,
	)

	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("create variable request failed: %w", err)
	}

	defer resp.Body.Close()

	return handleResponse(resp)
}

func (c *Client) GetVariable(
	ctx context.Context,
	owner string,
	repo string,
	name string,
) (*Variable, error) {
	path := fmt.Sprintf(
		"/repos/%s/%s/actions/variables/%s",
		owner,
		repo,
		name,
	)

	req, err := c.newRequest(
		ctx,
		http.MethodGet,
		path,
		nil,
	)

	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get variable request failed: %w", err)
	}

	defer resp.Body.Close()
	if err := handleResponse(resp); err != nil {
		return nil, err
	}

	var variable Variable

	if err := json.NewDecoder(resp.Body).Decode(&variable); err != nil {
		return nil, fmt.Errorf("decode variable response: %w", err)
	}

	return &variable, nil
}

func (c *Client) UpdateVariable(
	ctx context.Context,
	owner string,
	repo string,
	name string,
	value string,
) error {

	path := fmt.Sprintf(
		"/repos/%s/%s/actions/variables/%s",
		owner,
		repo,
		name,
	)

	body := Variable{
		Name:  name,
		Value: value,
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPatch,
		path,
		body,
	)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("update variable request failed: %w", err)
	}

	defer resp.Body.Close()

	return handleResponse(resp)
}

func (c *Client) DeleteVariable(
	ctx context.Context,
	owner string,
	repo string,
	name string,
) error {
	path := fmt.Sprintf(
		"/repos/%s/%s/actions/variables/%s",
		owner,
		repo,
		name,
	)

	req, err := c.newRequest(
		ctx,
		http.MethodDelete,
		path,
		nil,
	)

	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("delete variable request failed: %w", err)
	}

	defer resp.Body.Close()

	return handleResponse(resp)
}
