package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// PythonBackendClient handles communication with the Python FastAPI backend
type PythonBackendClient struct {
	baseURL    string
	httpClient *http.Client
	retries    int
	timeout    time.Duration
}

// NewPythonBackendClient creates a new Python backend client
func NewPythonBackendClient(baseURL string) *PythonBackendClient {
	return &PythonBackendClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		retries: 3,
		timeout: 30 * time.Second,
	}
}

// ============================================================================
// User API Calls
// ============================================================================

// GetUser retrieves user information from the Python backend
func (c *PythonBackendClient) GetUser(userID string) (*User, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", c.baseURL, userID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var user User
	if err := json.Unmarshal(body, &user); err != nil {
		log.Printf("Error unmarshaling user response: %v", err)
		return nil, err
	}

	return &user, nil
}

// GetUserAccounts retrieves user accounts from the Python backend
func (c *PythonBackendClient) GetUserAccounts(userID string) ([]*Account, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/accounts", c.baseURL, userID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var accounts []*Account
	if err := json.Unmarshal(body, &accounts); err != nil {
		log.Printf("Error unmarshaling accounts response: %v", err)
		return nil, err
	}

	return accounts, nil
}

// ============================================================================
// Account API Calls
// ============================================================================

// GetAccountBalance retrieves account balance from the Python backend
func (c *PythonBackendClient) GetAccountBalance(accountID string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s", c.baseURL, accountID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	var account Account
	if err := json.Unmarshal(body, &account); err != nil {
		log.Printf("Error unmarshaling account response: %v", err)
		return 0, err
	}

	return account.Balance, nil
}

// GetAccountStatements retrieves account statements from the Python backend
func (c *PythonBackendClient) GetAccountStatements(accountID, month string) (map[string]string, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s/statements", c.baseURL, accountID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		log.Printf("Error unmarshaling statements response: %v", err)
		return nil, err
	}

	statements := make(map[string]string)
	if stmts, ok := data["statements"].(map[string]interface{}); ok {
		for k, v := range stmts {
			statements[k] = fmt.Sprintf("%v", v)
		}
	}

	return statements, nil
}

// ============================================================================
// Payment API Calls
// ============================================================================

// TransferFunds initiates a payment transfer via the Python backend
func (c *PythonBackendClient) TransferFunds(fromAccount, toAccount string, amount float64) (*Payment, error) {
	url := fmt.Sprintf("%s/api/v1/payments/transfer", c.baseURL)

	payload := map[string]interface{}{
		"from_account": fromAccount,
		"to_account":   toAccount,
		"amount":       amount,
	}

	body, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}

	var payment Payment
	if err := json.Unmarshal(body, &payment); err != nil {
		log.Printf("Error unmarshaling payment response: %v", err)
		return nil, err
	}

	return &payment, nil
}

// ============================================================================
// Health Check
// ============================================================================

// HealthCheck verifies the Python backend is running
func (c *PythonBackendClient) HealthCheck() (bool, error) {
	url := fmt.Sprintf("%s/health", c.baseURL)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return false, err
	}

	var response HealthResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling health response: %v", err)
		return false, err
	}

	return response.Status == "healthy", nil
}

// ============================================================================
// Internal HTTP Request Handler
// ============================================================================

// makeRequest makes an HTTP request with retry logic
func (c *PythonBackendClient) makeRequest(method, url string, payload interface{}) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.retries; attempt++ {
		if attempt > 0 {
			log.Printf("Retry attempt %d for %s %s", attempt, method, url)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		var req *http.Request
		var err error

		if payload != nil {
			jsonData, err := json.Marshal(payload)
			if err != nil {
				return nil, fmt.Errorf("error marshaling payload: %w", err)
			}

			req, err = http.NewRequest(method, url, bytes.NewBuffer(jsonData))
			if err != nil {
				return nil, fmt.Errorf("error creating request: %w", err)
			}

			req.Header.Set("Content-Type", "application/json")
		} else {
			req, err = http.NewRequest(method, url, nil)
			if err != nil {
				return nil, fmt.Errorf("error creating request: %w", err)
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		// Check for non-2xx status codes
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			lastErr = fmt.Errorf("backend returned status %d: %s", resp.StatusCode, string(body))
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.retries, lastErr)
}
