package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://api.gofirefly.io"
	defaultUserAgent = "terraform-provider-firefly"
)

// Config holds the configuration needed to interact with the Firefly API
type Config struct {
	AccessKey  string
	SecretKey  string
	APIURL     string
	HTTPClient *http.Client
}

// Client is a client for interacting with the Firefly API
type Client struct {
	baseURL    *url.URL
	userAgent  string
	httpClient *http.Client
	
	// Authentication
	accessKey  string
	secretKey  string
	authToken  string
	expiresAt  time.Time
	
	// Services
	Workspaces        *WorkspaceService
	Guardrails        *GuardrailService
	Projects          *ProjectService
	RunnersWorkspaces *RunnersWorkspaceService
	VariableSets      *VariableSetService
}

// AuthResponse represents the response from the login endpoint
type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
	TokenType   string `json:"tokenType"`
}

// NewClient creates a new client for interacting with the Firefly API
func NewClient(config Config) (*Client, error) {
	if config.AccessKey == "" {
		return nil, fmt.Errorf("access key is required")
	}
	
	if config.SecretKey == "" {
		return nil, fmt.Errorf("secret key is required")
	}
	
	baseURL := defaultBaseURL
	if config.APIURL != "" {
		baseURL = config.APIURL
	}
	
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid API URL: %s", err)
	}
	
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: time.Second * 30,
		}
	}
	
	c := &Client{
		baseURL:    parsedBaseURL,
		userAgent:  defaultUserAgent,
		httpClient: httpClient,
		accessKey:  config.AccessKey,
		secretKey:  config.SecretKey,
	}
	
	// Create service endpoints
	c.Workspaces = &WorkspaceService{client: c}
	c.Guardrails = &GuardrailService{client: c}
	c.Projects = &ProjectService{client: c}
	c.RunnersWorkspaces = &RunnersWorkspaceService{client: c}
	c.VariableSets = &VariableSetService{client: c}
	
	return c, nil
}

// ensureAuthenticated ensures the client has a valid authentication token
func (c *Client) ensureAuthenticated() error {
	// If we have a token and it's not expired, we're good
	if c.authToken != "" && time.Now().Before(c.expiresAt) {
		return nil
	}
	
	// Otherwise, we need to authenticate
	reqBody, err := json.Marshal(map[string]string{
		"accessKey": c.accessKey,
		"secretKey": c.secretKey,
	})
	if err != nil {
		return fmt.Errorf("error encoding login request: %s", err)
	}
	
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v2/login", c.baseURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("error creating login request: %s", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error executing login request: %s", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("error decoding login response: %s", err)
	}
	
	c.authToken = authResp.AccessToken
	c.expiresAt = time.Unix(authResp.ExpiresAt, 0)
	
	return nil
}

// doRequest sends an HTTP request and returns an HTTP response
func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	// Ensure we're authenticated
	if err := c.ensureAuthenticated(); err != nil {
		return nil, err
	}
	
	// Set common headers
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.authToken))
	req.Header.Set("User-Agent", c.userAgent)
	
	// Execute the request
	return c.httpClient.Do(req)
}

// newRequest creates a new HTTP request
func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	u, err := c.baseURL.Parse(path)
	if err != nil {
		return nil, err
	}
	
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err = json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}
	
	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}
	
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req, nil
}
