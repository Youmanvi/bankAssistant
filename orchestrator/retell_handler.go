package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// RetellHandler handles Retell AI webhook events
type RetellHandler struct {
	stateMachine *CallStateMachine
	backendClient *PythonBackendClient
	apiKey       string
}

// NewRetellHandler creates a new Retell webhook handler
func NewRetellHandler(stateMachine *CallStateMachine, backendClient *PythonBackendClient, apiKey string) *RetellHandler {
	return &RetellHandler{
		stateMachine: stateMachine,
		backendClient: backendClient,
		apiKey: apiKey,
	}
}

// HandleWebhook processes incoming Retell AI webhooks
func (rh *RetellHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		log.Printf("Error reading webhook body: %v", err)
		return
	}
	defer r.Body.Close()

	// Verify the webhook signature
	if !rh.verifySignature(r, body) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Printf("Webhook signature verification failed")
		return
	}

	// Parse the webhook payload
	var payload RetellWebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		log.Printf("Error parsing webhook payload: %v", err)
		return
	}

	// Handle the webhook event
	switch payload.Event {
	case "call_started":
		rh.handleCallStarted(payload)
	case "call_ended":
		rh.handleCallEnded(payload)
	case "call_analyzed":
		rh.handleCallAnalyzed(payload)
	default:
		log.Printf("Unknown event type: %s", payload.Event)
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"received": true})
}

// ============================================================================
// Event Handlers
// ============================================================================

// handleCallStarted processes the call_started event
func (rh *RetellHandler) handleCallStarted(payload RetellWebhookPayload) {
	callID, ok := payload.Data["call_id"].(string)
	if !ok {
		log.Printf("Missing or invalid call_id in call_started event")
		return
	}

	phoneNumber, ok := payload.Data["phone_number"].(string)
	if !ok {
		phoneNumber = "unknown"
	}

	log.Printf("Call started: %s (Phone: %s)", callID, phoneNumber)

	// For now, we don't know the user ID until they authenticate
	// This will be updated once we have user context
	_, err := rh.stateMachine.CreateCall(callID, "unknown", phoneNumber)
	if err != nil {
		log.Printf("Error creating call: %v", err)
		return
	}

	// Update state to CALL_STARTED
	if err := rh.stateMachine.UpdateState(callID, CallStarted); err != nil {
		log.Printf("Error updating call state: %v", err)
		return
	}

	// Store raw event data in metadata
	rh.stateMachine.UpdateMetadata(callID, "retell_data", payload.Data)
}

// handleCallEnded processes the call_ended event
func (rh *RetellHandler) handleCallEnded(payload RetellWebhookPayload) {
	callID, ok := payload.Data["call_id"].(string)
	if !ok {
		log.Printf("Missing or invalid call_id in call_ended event")
		return
	}

	log.Printf("Call ended: %s", callID)

	// Update state to CALL_ENDED
	if err := rh.stateMachine.UpdateState(callID, CallEnded); err != nil {
		log.Printf("Error updating call state: %v", err)
		return
	}

	// Store end event data
	rh.stateMachine.UpdateMetadata(callID, "end_data", payload.Data)

	// Clean up the call after a delay (for audit trail)
	// In production, you might want to keep this for logging/analytics
	// For now, we'll keep it for analysis
}

// handleCallAnalyzed processes the call_analyzed event
func (rh *RetellHandler) handleCallAnalyzed(payload RetellWebhookPayload) {
	callID, ok := payload.Data["call_id"].(string)
	if !ok {
		log.Printf("Missing or invalid call_id in call_analyzed event")
		return
	}

	log.Printf("Call analyzed: %s", callID)

	// Store analyzed data in metadata
	rh.stateMachine.UpdateMetadata(callID, "analyzed_data", payload.Data)
}

// ============================================================================
// Signature Verification
// ============================================================================

// verifySignature verifies the Retell API signature
func (rh *RetellHandler) verifySignature(r *http.Request, body []byte) bool {
	signature := r.Header.Get("X-Retell-Signature")
	if signature == "" {
		return false
	}

	// Compute HMAC-SHA256
	hash := hmac.New(sha256.New, []byte(rh.apiKey))
	hash.Write(body)
	expectedSignature := hex.EncodeToString(hash.Sum(nil))

	// Compare signatures
	return hmac.Equal([]byte(signature), []byte(expectedSignature))
}

// ============================================================================
// Admin Status Endpoints
// ============================================================================

// GetCallStatus returns the current status of a call
func (rh *RetellHandler) GetCallStatus(w http.ResponseWriter, r *http.Request) {
	callID := r.URL.Query().Get("call_id")
	if callID == "" {
		http.Error(w, "Missing call_id parameter", http.StatusBadRequest)
		return
	}

	call, err := rh.stateMachine.GetCall(callID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Call not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(call)
}

// GetAllCalls returns information about all active calls
func (rh *RetellHandler) GetAllCalls(w http.ResponseWriter, r *http.Request) {
	calls := rh.stateMachine.GetAllCalls()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(calls),
		"calls": calls,
	})
}

// ============================================================================
// User Context Loading (Called during orchestration)
// ============================================================================

// LoadUserContext loads user context from the Python backend
func (rh *RetellHandler) LoadUserContext(callID, userID string) error {
	// Fetch user from Python backend
	user, err := rh.backendClient.GetUser(userID)
	if err != nil {
		log.Printf("Error loading user context: %v", err)
		return err
	}

	// Fetch user accounts
	accounts, err := rh.backendClient.GetUserAccounts(userID)
	if err != nil {
		log.Printf("Error loading user accounts: %v", err)
		return err
	}

	// Update call context with user info
	if err := rh.stateMachine.UpdateMetadata(callID, "user", user); err != nil {
		return err
	}

	if err := rh.stateMachine.UpdateMetadata(callID, "accounts", accounts); err != nil {
		return err
	}

	log.Printf("User context loaded for call %s (User: %s, Accounts: %d)", callID, userID, len(accounts))

	return nil
}
