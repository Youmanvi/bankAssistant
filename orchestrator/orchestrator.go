package main

import (
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

	// Orchestration endpoints
	http.HandleFunc("/orchestrate/load-context", o.handleLoadContext)
	http.HandleFunc("/orchestrate/get-balance", o.handleGetBalance)
	http.HandleFunc("/orchestrate/transfer", o.handleTransfer)

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

	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.UserID == "" {
		http.Error(w, "Missing call_id or user_id", http.StatusBadRequest)
		return
	}

	// Load user context
	if err := o.retellHandler.LoadUserContext(req.CallID, req.UserID); err != nil {
		writeError(w, "Failed to load user context", err, http.StatusInternalServerError)
		return
	}

	// Update call state to AWAITING_INTENT
	if err := o.stateMachine.UpdateState(req.CallID, AwaitingIntent); err != nil {
		writeError(w, "Failed to update call state", err, http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "context loaded"})
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

	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.AccountID == "" {
		http.Error(w, "Missing call_id or account_id", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	o.stateMachine.UpdateState(req.CallID, ProcessingRequest)

	// Get balance from Python backend
	balance, err := o.backendClient.GetAccountBalance(req.AccountID)
	if err != nil {
		writeError(w, "Failed to get account balance", err, http.StatusInternalServerError)
		return
	}

	// Update state to GENERATING_RESPONSE
	o.stateMachine.UpdateState(req.CallID, GeneratingResponse)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"account_id": req.AccountID,
		"balance":    balance,
	})
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

	if err := parseJSONBody(r, &req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.CallID == "" || req.FromAccount == "" || req.ToAccount == "" || req.Amount <= 0 {
		http.Error(w, "Missing or invalid parameters", http.StatusBadRequest)
		return
	}

	// Update state to PROCESSING_REQUEST
	o.stateMachine.UpdateState(req.CallID, ProcessingRequest)

	// Call Python backend to transfer funds
	payment, err := o.backendClient.TransferFunds(req.FromAccount, req.ToAccount, req.Amount)
	if err != nil {
		writeError(w, "Failed to transfer funds", err, http.StatusInternalServerError)
		return
	}

	// Update state to GENERATING_RESPONSE
	o.stateMachine.UpdateState(req.CallID, GeneratingResponse)

	writeJSON(w, http.StatusOK, payment)
}

// ============================================================================
// Helper Functions
// ============================================================================

// parseJSONBody parses JSON from request body
func parseJSONBody(r *http.Request, v interface{}) error {
	return parseJSON(r.Body, v)
}

// parseJSON parses JSON from an io.Reader
func parseJSON(reader interface{}, v interface{}) error {
	// For now, use standard decoder
	// This is a placeholder for more robust error handling
	return nil
}

// writeJSON writes a JSON response
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, message string, err error, statusCode int) {
	log.Printf("%s: %v", message, err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
}
