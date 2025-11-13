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
func (c *PythonBackendClient) GetUser(userID string) (*UserResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s", c.baseURL, userID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	var response UserResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling user response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal user response: %w", err)
	}

	return &response, nil
}

// GetUserProfile retrieves user profile from the Python backend
func (c *PythonBackendClient) GetUserProfile(userID string) (*UserProfileResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/profile", c.baseURL, userID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user profile: %w", err)
	}

	var response UserProfileResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling user profile response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal user profile response: %w", err)
	}

	return &response, nil
}

// GetUserAccounts retrieves user accounts from the Python backend
func (c *PythonBackendClient) GetUserAccounts(userID string) (*UserAccountsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/users/%s/accounts", c.baseURL, userID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get user accounts: %w", err)
	}

	var response UserAccountsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling accounts response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal accounts response: %w", err)
	}

	return &response, nil
}

// ============================================================================
// Account API Calls
// ============================================================================

// GetAccount retrieves full account information from the Python backend
func (c *PythonBackendClient) GetAccount(accountID string) (*AccountResponse, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s", c.baseURL, accountID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	var response AccountResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling account response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal account response: %w", err)
	}

	return &response, nil
}

// GetAccountBalance retrieves account balance from the Python backend
func (c *PythonBackendClient) GetAccountBalance(accountID string) (*BalanceResponse, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s/balance", c.baseURL, accountID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balance: %w", err)
	}

	var response BalanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling balance response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal balance response: %w", err)
	}

	return &response, nil
}

// GetAccountStatements retrieves account statements from the Python backend
func (c *PythonBackendClient) GetAccountStatements(accountID string) (*StatementsResponse, error) {
	url := fmt.Sprintf("%s/api/v1/accounts/%s/statements", c.baseURL, accountID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get account statements: %w", err)
	}

	var response StatementsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling statements response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal statements response: %w", err)
	}

	return &response, nil
}

// ============================================================================
// Payment API Calls
// ============================================================================

// TransferFunds initiates a payment transfer via the Python backend
func (c *PythonBackendClient) TransferFunds(fromAccount, toAccount string, amount float64) (*PaymentResponse, error) {
	url := fmt.Sprintf("%s/api/v1/payments/transfer", c.baseURL)

	payload := map[string]interface{}{
		"from_account": fromAccount,
		"to_account":   toAccount,
		"amount":       amount,
	}

	body, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to transfer funds: %w", err)
	}

	var response PaymentResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling payment response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal payment response: %w", err)
	}

	return &response, nil
}

// ============================================================================
// Application API Calls
// ============================================================================

// ApplyForLoan initiates a loan application via the Python backend
func (c *PythonBackendClient) ApplyForLoan(userID string, loanAmount float64, loanPurpose string, termYears int) (*LoanApplicationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/applications/loan", c.baseURL)

	payload := map[string]interface{}{
		"user_id":      userID,
		"loan_amount":  loanAmount,
		"loan_purpose": loanPurpose,
		"term_years":   termYears,
	}

	body, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to apply for loan: %w", err)
	}

	var response LoanApplicationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling loan application response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal loan application response: %w", err)
	}

	return &response, nil
}

// ApplyForCreditCard initiates a credit card application via the Python backend
func (c *PythonBackendClient) ApplyForCreditCard(userID, cardType string, creditLimit float64) (*CreditCardApplicationResponse, error) {
	url := fmt.Sprintf("%s/api/v1/applications/credit-card", c.baseURL)

	payload := map[string]interface{}{
		"user_id":      userID,
		"card_type":    cardType,
		"credit_limit": creditLimit,
	}

	body, err := c.makeRequest("POST", url, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to apply for credit card: %w", err)
	}

	var response CreditCardApplicationResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling credit card application response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal credit card application response: %w", err)
	}

	return &response, nil
}

// GetApplicationStatus retrieves the status of an application
func (c *PythonBackendClient) GetApplicationStatus(applicationID string) (*ApplicationStatusResponse, error) {
	url := fmt.Sprintf("%s/api/v1/applications/%s", c.baseURL, applicationID)

	body, err := c.makeRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get application status: %w", err)
	}

	var response ApplicationStatusResponse
	if err := json.Unmarshal(body, &response); err != nil {
		log.Printf("Error unmarshaling application status response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal application status response: %w", err)
	}

	return &response, nil
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
