package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Orchestrator manages the overall orchestration service
type Orchestrator struct {
	config           *Config
	stateMachine     *CallStateMachine
	backendClient    *PythonBackendClient
	retellHandler    *RetellHandler
	httpServer       *http.Server
}

// NewOrchestrator creates a new orchestrator instance
func NewOrchestrator(config *Config) *Orchestrator {
	stateMachine := NewCallStateMachine()
	backendClient := NewPythonBackendClient(config.PythonBackendURL)
	retellHandler := NewRetellHandler(stateMachine, backendClient, config.RetellAPIKey)

	return &Orchestrator{
		config:        config,
		stateMachine: stateMachine,
		backendClient: backendClient,
		retellHandler: retellHandler,
	}
}

// Start starts the orchestrator server
func (o *Orchestrator) Start() error {
	// Verify Python backend is healthy
	healthy, err := o.backendClient.HealthCheck()
	if err != nil {
		log.Printf("Warning: Python backend health check failed: %v", err)
	} else if healthy {
		log.Printf("Python backend is healthy")
	} else {
		log.Printf("Warning: Python backend reported unhealthy status")
	}

	// Set up HTTP routes
	o.setupRoutes()

	// Create HTTP server
	o.httpServer = &http.Server{
		Addr:         o.config.Host + ":" + o.config.Port,
		Handler:      http.DefaultServeMux,
		ReadTimeout:  15,
		WriteTimeout: 15,
		IdleTimeout:  60,
	}

	// Start server
	log.Printf("Orchestrator starting on %s:%s", o.config.Host, o.config.Port)
	if err := o.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Stop gracefully stops the orchestrator server
func (o *Orchestrator) Stop() error {
	if o.httpServer != nil {
		log.Printf("Orchestrator shutting down...")
		return o.httpServer.Close()
	}
	return nil
}

// ============================================================================
// Routes Setup
// ============================================================================

// setupRoutes configures all HTTP routes
func (o *Orchestrator) setupRoutes() {
	// Health check
	http.HandleFunc("/health", o.handleHealth)

	// Retell AI webhook
	http.HandleFunc("/webhook", o.retellHandler.HandleWebhook)

	// Admin status endpoints
	http.HandleFunc("/admin/calls", o.handleAdminCalls)
	http.HandleFunc("/admin/call", o.retellHandler.GetCallStatus)

	// Orchestration endpoints - User & Context
	http.HandleFunc("/orchestrate/load-context", o.handleLoadContext)
	http.HandleFunc("/orchestrate/get-user", o.handleGetUser)
	http.HandleFunc("/orchestrate/get-user-profile", o.handleGetUserProfile)
	http.HandleFunc("/orchestrate/get-accounts", o.handleGetAccounts)

	// Orchestration endpoints - Accounts
	http.HandleFunc("/orchestrate/get-balance", o.handleGetBalance)
	http.HandleFunc("/orchestrate/get-statements", o.handleGetStatements)

	// Orchestration endpoints - Payments
	http.HandleFunc("/orchestrate/transfer", o.handleTransfer)

	// Orchestration endpoints - Applications
	http.HandleFunc("/orchestrate/apply-loan", o.handleApplyLoan)
	http.HandleFunc("/orchestrate/apply-credit-card", o.handleApplyCreditCard)
	http.HandleFunc("/orchestrate/get-application-status", o.handleGetApplicationStatus)

	log.Printf("Routes configured")
}

// ============================================================================
// Health Check Handler
// ============================================================================

// handleHealth responds to health check requests
func (o *Orchestrator) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "version": "0.1.0"}`))
}

// ============================================================================
// Admin Handlers
// ============================================================================

// handleAdminCalls returns all active calls
func (o *Orchestrator) handleAdminCalls(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	o.retellHandler.GetAllCalls(w, r)
}

// ============================================================================
// User & Context Orchestration Handlers
// ============================================================================

// handleGetUser retrieves user information
func (o *Orchestrator) handleGetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID string `json:"call_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		writeErrorJSON(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// Get user from Python backend
	response, err := o.backendClient.GetUser(req.UserID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting user: %v", err)
		return
	}

	// Update metadata if call_id provided
	if req.CallID != "" {
		o.stateMachine.UpdateMetadata(req.CallID, "user", response)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleGetUserProfile retrieves user profile
func (o *Orchestrator) handleGetUserProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		writeErrorJSON(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// Get user profile from Python backend
	response, err := o.backendClient.GetUserProfile(req.UserID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get user profile: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting user profile: %v", err)
		return
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleGetAccounts retrieves user accounts
func (o *Orchestrator) handleGetAccounts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID string `json:"call_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		writeErrorJSON(w, "Missing user_id", http.StatusBadRequest)
		return
	}

	// Get accounts from Python backend
	response, err := o.backendClient.GetUserAccounts(req.UserID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get user accounts: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting user accounts: %v", err)
		return
	}

	// Update metadata if call_id provided
	if req.CallID != "" {
		o.stateMachine.UpdateMetadata(req.CallID, "accounts", response.Accounts)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// ============================================================================
// Orchestration Handlers
// ============================================================================

// handleLoadContext loads user context for a call
func (o *Orchestrator) handleLoadContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID string `json:"call_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.UserID == "" {
		writeErrorJSON(w, "Missing call_id or user_id", http.StatusBadRequest)
		return
	}

	// Load user context
	if err := o.retellHandler.LoadUserContext(req.CallID, req.UserID); err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to load user context: %v", err), http.StatusInternalServerError)
		log.Printf("Error loading user context: %v", err)
		return
	}

	// Update call state to AWAITING_INTENT
	if err := o.stateMachine.UpdateState(req.CallID, AwaitingIntent); err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to update call state: %v", err), http.StatusInternalServerError)
		return
	}

	writeSuccessJSON(w, http.StatusOK, map[string]string{"status": "context loaded"})
}

// handleGetBalance gets account balance
func (o *Orchestrator) handleGetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID    string `json:"call_id"`
		AccountID string `json:"account_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.AccountID == "" {
		writeErrorJSON(w, "Missing call_id or account_id", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	o.stateMachine.UpdateState(req.CallID, ProcessingRequest)

	// Get balance from Python backend
	response, err := o.backendClient.GetAccountBalance(req.AccountID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get account balance: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting account balance: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	o.stateMachine.UpdateState(req.CallID, GeneratingResponse)

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleGetStatements retrieves account statements
func (o *Orchestrator) handleGetStatements(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID    string `json:"call_id"`
		AccountID string `json:"account_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.AccountID == "" {
		writeErrorJSON(w, "Missing account_id", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST if call_id provided
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, ProcessingRequest)
	}

	// Get statements from Python backend
	response, err := o.backendClient.GetAccountStatements(req.AccountID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get account statements: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting account statements: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, GeneratingResponse)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleTransfer initiates a fund transfer
func (o *Orchestrator) handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID      string  `json:"call_id"`
		FromAccount string  `json:"from_account"`
		ToAccount   string  `json:"to_account"`
		Amount      float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.FromAccount == "" || req.ToAccount == "" || req.Amount <= 0 {
		writeErrorJSON(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	o.stateMachine.UpdateState(req.CallID, ProcessingRequest)

	// Call Python backend to transfer funds
	payment, err := o.backendClient.TransferFunds(req.FromAccount, req.ToAccount, req.Amount)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to transfer funds: %v", err), http.StatusInternalServerError)
		log.Printf("Error transferring funds: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	o.stateMachine.UpdateState(req.CallID, GeneratingResponse)

	writeSuccessJSON(w, http.StatusOK, payment)
}

// handleApplyLoan initiates a loan application
func (o *Orchestrator) handleApplyLoan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID      string  `json:"call_id"`
		UserID      string  `json:"user_id"`
		LoanAmount  float64 `json:"loan_amount"`
		LoanPurpose string  `json:"loan_purpose"`
		TermYears   int     `json:"term_years"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.LoanAmount <= 0 || req.TermYears <= 0 {
		writeErrorJSON(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, ProcessingRequest)
	}

	// Call Python backend to apply for loan
	response, err := o.backendClient.ApplyForLoan(req.UserID, req.LoanAmount, req.LoanPurpose, req.TermYears)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to apply for loan: %v", err), http.StatusInternalServerError)
		log.Printf("Error applying for loan: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, GeneratingResponse)
		o.stateMachine.UpdateMetadata(req.CallID, "application_id", response.ApplicationID)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleApplyCreditCard initiates a credit card application
func (o *Orchestrator) handleApplyCreditCard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID      string  `json:"call_id"`
		UserID      string  `json:"user_id"`
		CardType    string  `json:"card_type"`
		CreditLimit float64 `json:"credit_limit"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" || req.CreditLimit <= 0 {
		writeErrorJSON(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, ProcessingRequest)
	}

	// Call Python backend to apply for credit card
	response, err := o.backendClient.ApplyForCreditCard(req.UserID, req.CardType, req.CreditLimit)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to apply for credit card: %v", err), http.StatusInternalServerError)
		log.Printf("Error applying for credit card: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, GeneratingResponse)
		o.stateMachine.UpdateMetadata(req.CallID, "application_id", response.ApplicationID)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// handleGetApplicationStatus retrieves application status
func (o *Orchestrator) handleGetApplicationStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		CallID        string `json:"call_id"`
		ApplicationID string `json:"application_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErrorJSON(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.ApplicationID == "" {
		writeErrorJSON(w, "Missing application_id", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, ProcessingRequest)
	}

	// Get application status from Python backend
	response, err := o.backendClient.GetApplicationStatus(req.ApplicationID)
	if err != nil {
		writeErrorJSON(w, fmt.Sprintf("Failed to get application status: %v", err), http.StatusInternalServerError)
		log.Printf("Error getting application status: %v", err)
		return
	}

	// Update state to GENERATING_RESPONSE
	if req.CallID != "" {
		o.stateMachine.UpdateState(req.CallID, GeneratingResponse)
	}

	writeSuccessJSON(w, http.StatusOK, response)
}

// ============================================================================
// Helper Functions
// ============================================================================

// writeSuccessJSON writes a successful JSON response
func writeSuccessJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// writeErrorJSON writes an error JSON response
func writeErrorJSON(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":  message,
		"status": statusCode,
	})
}
