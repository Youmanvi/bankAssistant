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

// User represents a user from the Python backend
type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Accounts []string `json:"accounts"`
	SSN     string   `json:"ssn"`
	Email   string   `json:"email"`
	Phone   string   `json:"phone"`
}

// Account represents a bank account from the Python backend
type Account struct {
	ID       string             `json:"id"`
	UserID   string             `json:"user_id"`
	Balance  float64            `json:"balance"`
	Type     string             `json:"type"`
	Statements map[string]string `json:"statements"`
}

// Payment represents a payment transaction
type Payment struct {
	ID          string  `json:"id"`
	FromAccount string  `json:"from_account"`
	ToAccount   string  `json:"to_account"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
	Status      string  `json:"status"`
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
