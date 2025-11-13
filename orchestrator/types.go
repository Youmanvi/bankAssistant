package main

// ============================================================================
// Call State Management
// ============================================================================

// CallState represents the state of a call in the state machine
type CallState string

const (
	AwaitingCall        CallState = "AWAITING_CALL"
	CallStarted         CallState = "CALL_STARTED"
	AwaitingIntent      CallState = "AWAITING_INTENT"
	ProcessingRequest   CallState = "PROCESSING_REQUEST"
	GeneratingResponse  CallState = "GENERATING_RESPONSE"
	SpeakingResponse    CallState = "SPEAKING_RESPONSE"
	CallEnded           CallState = "CALL_ENDED"
)

// CallContext holds the state and metadata for a single call
type CallContext struct {
	CallID      string                 `json:"call_id"`
	UserID      string                 `json:"user_id"`
	PhoneNumber string                 `json:"phone_number"`
	State       CallState              `json:"state"`
	Metadata    map[string]interface{} `json:"metadata"`
	CreatedAt   int64                  `json:"created_at"`
	UpdatedAt   int64                  `json:"updated_at"`
}

// ============================================================================
// Retell AI Webhook Events
// ============================================================================

// RetellEventType represents types of Retell events
type RetellEventType string

const (
	EventCallStarted  RetellEventType = "call_started"
	EventCallEnded    RetellEventType = "call_ended"
	EventCallAnalyzed RetellEventType = "call_analyzed"
)

// RetellWebhookPayload represents the webhook payload from Retell AI
type RetellWebhookPayload struct {
	Event string                 `json:"event"`
	Data  map[string]interface{} `json:"data"`
}

// RetellCallData contains call information from Retell
type RetellCallData struct {
	CallID      string `json:"call_id"`
	PhoneNumber string `json:"phone_number"`
	RemotePhone string `json:"remote_phone"`
	Timestamp   int64  `json:"timestamp"`
}

// ============================================================================
// Python Backend API Models
// ============================================================================

// UserDetails represents user information from Python backend
type UserDetails struct {
	Name         string   `json:"name"`
	Accounts     []string `json:"accounts"`
	SSN          string   `json:"ssn"`
	Email        string   `json:"email"`
	Phone        string   `json:"phone"`
	Address      string   `json:"address"`
	DateOfBirth  string   `json:"date_of_birth"`
}

// UserResponse represents the response from GET /api/v1/users/{user_id}
type UserResponse struct {
	UserID  string      `json:"user_id"`
	Details UserDetails `json:"details"`
}

// UserProfileResponse represents the response from GET /api/v1/users/{user_id}/profile
type UserProfileResponse struct {
	UserID  string `json:"user_id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// AccountDetails represents account information
type AccountDetails struct {
	Balance    float64            `json:"balance"`
	Type       string             `json:"type"`
	Statements map[string]string  `json:"statements"`
}

// AccountResponse represents the response from GET /api/v1/accounts/{account_id}
type AccountResponse struct {
	AccountID string         `json:"account_id"`
	Details   AccountDetails `json:"details"`
}

// BalanceResponse represents the response from GET /api/v1/accounts/{account_id}/balance
type BalanceResponse struct {
	AccountID string  `json:"account_id"`
	Balance   float64 `json:"balance"`
}

// StatementsResponse represents the response from GET /api/v1/accounts/{account_id}/statements
type StatementsResponse struct {
	AccountID  string            `json:"account_id"`
	Statements map[string]string `json:"statements"`
}

// UserAccountsResponse represents the response from GET /api/v1/users/{user_id}/accounts
type UserAccountsResponse struct {
	UserID       string                    `json:"user_id"`
	AccountCount int                       `json:"account_count"`
	Accounts     map[string]AccountDetails `json:"accounts"`
}

// PaymentResponse represents the response from POST /api/v1/payments/transfer
type PaymentResponse struct {
	Status      string  `json:"status"`
	PaymentID   string  `json:"payment_id"`
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
}

// LoanApplicationRequest represents a loan application
type LoanApplicationRequest struct {
	UserID      string  `json:"user_id"`
	LoanAmount  float64 `json:"loan_amount"`
	LoanPurpose string  `json:"loan_purpose"`
	TermYears   int     `json:"term_years"`
}

// LoanApplicationResponse represents the response from POST /api/v1/applications/loan
type LoanApplicationResponse struct {
	Status          string  `json:"status"`
	ApplicationID   string  `json:"application_id"`
	UserID          string  `json:"user_id"`
	LoanAmount      float64 `json:"loan_amount"`
	LoanPurpose     string  `json:"loan_purpose"`
	TermYears       int     `json:"term_years"`
	Message         string  `json:"message"`
}

// CreditCardApplicationRequest represents a credit card application
type CreditCardApplicationRequest struct {
	UserID      string  `json:"user_id"`
	CardType    string  `json:"card_type"`
	CreditLimit float64 `json:"credit_limit"`
}

// CreditCardApplicationResponse represents the response from POST /api/v1/applications/credit-card
type CreditCardApplicationResponse struct {
	Status        string  `json:"status"`
	ApplicationID string  `json:"application_id"`
	UserID        string  `json:"user_id"`
	CardType      string  `json:"card_type"`
	CreditLimit   float64 `json:"credit_limit"`
	Message       string  `json:"message"`
}

// ApplicationStatusResponse represents the response from GET /api/v1/applications/{application_id}
type ApplicationStatusResponse struct {
	ApplicationID string `json:"application_id"`
	Status        string `json:"status"`
	Message       string `json:"message"`
}

// ============================================================================
// Authentication Models
// ============================================================================

// AuthRequest represents a login request
type AuthRequest struct {
	Phone    string `json:"phone"`
	PIN      string `json:"pin"`
}

// AuthResponse represents a login response
type AuthResponse struct {
	Success   bool   `json:"success"`
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Message   string `json:"message"`
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID string `json:"user_id"`
	Phone  string `json:"phone"`
	Name   string `json:"name"`
	Exp    int64  `json:"exp"`
}

// UserSession represents an authenticated user session
type UserSession struct {
	Token     string
	UserID    string
	Phone     string
	Name      string
	CreatedAt int64
	ExpiresAt int64
}

// ============================================================================
// Sample Data Generation Models
// ============================================================================

// UserSeedRequest represents a request to generate a test user with sample data
type UserSeedRequest struct {
	Phone  string `json:"phone"`
	PIN    string `json:"pin"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

// UserSeedResponse represents the response from user seed creation
type UserSeedResponse struct {
	Success   bool                 `json:"success"`
	UserID    string               `json:"user_id"`
	Phone     string               `json:"phone"`
	Name      string               `json:"name"`
	Accounts  []SampleAccountData  `json:"accounts"`
	Message   string               `json:"message"`
}

// SampleAccountData represents account data for a test user
type SampleAccountData struct {
	AccountID  string  `json:"account_id"`
	Type       string  `json:"type"`
	Balance    float64 `json:"balance"`
	Transactions int    `json:"transactions"`
}

// ============================================================================
// HTTP Response Models
// ============================================================================

// HealthResponse represents a health check response
type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error      string `json:"error"`
	Detail     string `json:"detail,omitempty"`
	StatusCode int    `json:"status_code"`
}

// ============================================================================
// Internal Service Models
// ============================================================================

// OrchestrationRequest represents a request to orchestrate operations
type OrchestrationRequest struct {
	CallID   string                 `json:"call_id"`
	UserID   string                 `json:"user_id"`
	Intent   string                 `json:"intent"`
	Metadata map[string]interface{} `json:"metadata"`
}

// OrchestrationResponse represents a response from orchestration
type OrchestrationResponse struct {
	Success bool                   `json:"success"`
	Data    map[string]interface{} `json:"data"`
	Error   string                 `json:"error,omitempty"`
}
