package ranger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Client represents the Ranger client
type Client struct {
	BaseURL  string
	Username string
	Password string
}

// NewClient creates a new Ranger client
func NewClient(baseURL, username string, password string) *Client {
	return &Client{
		BaseURL:  baseURL,
		Username: username,
		Password: password,
	}
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.Username, c.Password)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	if err := handleNon2xxResponse(resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func handleNon2xxResponse(resp *http.Response) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed, status code: %d, response: %s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *Client) GetPolicies(serviceName ...string) ([]Policy, error) {
	if len(serviceName) > 1 {
		return nil, fmt.Errorf("only one service name can be provided, got %d", len(serviceName))
	}

	uri := fmt.Sprintf("%s/service/public/v2/api/policy", c.BaseURL)
	if len(serviceName) > 0 && serviceName[0] != "" {
		encodedServiceName := url.QueryEscape(serviceName[0])
		uri += fmt.Sprintf("?serviceName=%s", encodedServiceName)
	}

	req, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, fmt.Errorf("error creating get policies request to %s: %w", uri, err)
	}

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making get policies request to %s: %w", uri, err)
	}

	defer resp.Body.Close()

	var policies []Policy

	if err := json.NewDecoder(resp.Body).Decode(&policies); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return policies, nil
}

func (c *Client) CreatePolicy(p *Policy) (*Policy, error) {
	policyJSON, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error marshalling policy: %w", err)
	}

	uri := fmt.Sprintf("%s/service/public/v2/api/policy", c.BaseURL)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(policyJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating create policy request to %s: %w", uri, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making create policy request to %s: %w", uri, err)
	}
	defer resp.Body.Close()

	var createdPolicy Policy
	if err := json.NewDecoder(resp.Body).Decode(&createdPolicy); err != nil {
		return nil, fmt.Errorf("error decoding create policy response: %w", err)
	}

	return &createdPolicy, nil
}

func (c *Client) DeletePolicy(serviceName string, policyName string) error {
	encodedPolicyName := url.QueryEscape(policyName)
	encodedServiceName := url.QueryEscape(serviceName)
	uri := fmt.Sprintf("%s/service/public/v2/api/policy?servicename=%s&policyname=%s", c.BaseURL, encodedServiceName, encodedPolicyName)

	req, err := http.NewRequest("DELETE", uri, nil)
	if err != nil {
		return fmt.Errorf("error creating delete request to uri %s: %w", uri, err)
	}

	_, err = c.doRequest(req)
	if err != nil {
		return fmt.Errorf("error making delete policy request to uri %s: %w", uri, err)
	}

	return nil
}

func (c *Client) UpdatePolicy(policy *Policy) (*Policy, error) {
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return nil, fmt.Errorf("error marshalling policy: %w", err)
	}

	uri := fmt.Sprintf("%s/service/public/v2/api/policy/%d", c.BaseURL, policy.ID)

	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(policyJSON))
	if err != nil {
		return nil, fmt.Errorf("error creating update policy request to %s: %w", uri, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.doRequest(req)
	if err != nil {
		return nil, fmt.Errorf("error making update policy request to %s: %w", uri, err)
	}
	defer resp.Body.Close()

	var updatedPolicy Policy
	if err := json.NewDecoder(resp.Body).Decode(&updatedPolicy); err != nil {
		return nil, fmt.Errorf("error decoding update policy response: %w", err)
	}

	return &updatedPolicy, nil
}
