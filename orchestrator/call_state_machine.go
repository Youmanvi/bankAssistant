package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// CallStateMachine manages the lifecycle of a call
type CallStateMachine struct {
	mu    sync.RWMutex
	calls map[string]*CallContext
}

// NewCallStateMachine creates a new call state machine
func NewCallStateMachine() *CallStateMachine {
	return &CallStateMachine{
		calls: make(map[string]*CallContext),
	}
}

// CreateCall creates a new call context
func (csm *CallStateMachine) CreateCall(callID, userID, phoneNumber string) (*CallContext, error) {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	if _, exists := csm.calls[callID]; exists {
		return nil, fmt.Errorf("call already exists: %s", callID)
	}

	now := time.Now().Unix()
	call := &CallContext{
		CallID:      callID,
		UserID:      userID,
		PhoneNumber: phoneNumber,
		State:       AwaitingCall,
		Metadata:    make(map[string]interface{}),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	csm.calls[callID] = call
	log.Printf("Call created: %s (User: %s, Phone: %s)", callID, userID, phoneNumber)

	return call, nil
}

// GetCall retrieves a call context by ID
func (csm *CallStateMachine) GetCall(callID string) (*CallContext, error) {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	call, exists := csm.calls[callID]
	if !exists {
		return nil, fmt.Errorf("call not found: %s", callID)
	}

	return call, nil
}

// UpdateState transitions a call to a new state
func (csm *CallStateMachine) UpdateState(callID string, newState CallState) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	call, exists := csm.calls[callID]
	if !exists {
		return fmt.Errorf("call not found: %s", callID)
	}

	// Validate state transition
	if !isValidTransition(call.State, newState) {
		return fmt.Errorf("invalid state transition: %s -> %s", call.State, newState)
	}

	oldState := call.State
	call.State = newState
	call.UpdatedAt = time.Now().Unix()

	log.Printf("Call state updated: %s [%s -> %s]", callID, oldState, newState)

	return nil
}

// UpdateMetadata updates call metadata
func (csm *CallStateMachine) UpdateMetadata(callID string, key string, value interface{}) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	call, exists := csm.calls[callID]
	if !exists {
		return fmt.Errorf("call not found: %s", callID)
	}

	call.Metadata[key] = value
	call.UpdatedAt = time.Now().Unix()

	return nil
}

// DeleteCall removes a call from the state machine
func (csm *CallStateMachine) DeleteCall(callID string) error {
	csm.mu.Lock()
	defer csm.mu.Unlock()

	if _, exists := csm.calls[callID]; !exists {
		return fmt.Errorf("call not found: %s", callID)
	}

	delete(csm.calls, callID)
	log.Printf("Call deleted: %s", callID)

	return nil
}

// GetAllCalls returns all active calls
func (csm *CallStateMachine) GetAllCalls() []*CallContext {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	calls := make([]*CallContext, 0, len(csm.calls))
	for _, call := range csm.calls {
		calls = append(calls, call)
	}

	return calls
}

// CountCalls returns the number of active calls
func (csm *CallStateMachine) CountCalls() int {
	csm.mu.RLock()
	defer csm.mu.RUnlock()

	return len(csm.calls)
}

// ============================================================================
// State Transition Validation
// ============================================================================

// isValidTransition checks if a state transition is valid
func isValidTransition(from, to CallState) bool {
	validTransitions := map[CallState][]CallState{
		AwaitingCall:       {CallStarted},
		CallStarted:        {AwaitingIntent},
		AwaitingIntent:     {ProcessingRequest},
		ProcessingRequest:  {GeneratingResponse},
		GeneratingResponse: {SpeakingResponse},
		SpeakingResponse:   {AwaitingIntent, CallEnded},
		CallEnded:          {AwaitingCall},
	}

	validStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, state := range validStates {
		if state == to {
			return true
		}
	}

	return false
}
