package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"math/big"
	"sync"
	"time"
)

// AuthService handles user authentication and session management
type AuthService struct {
	mu       sync.RWMutex
	sessions map[string]*UserSession // token -> session
	users    map[string]*AuthUser     // phone -> user
	jwtKey   []byte
}

// AuthUser represents a user for authentication
type AuthUser struct {
	UserID      string
	Phone       string
	PIN         string
	Name        string
	Email       string
	Address     string
	DateOfBirth string
	SSN         string
	Accounts    []string
	CreatedAt   int64
}

// NewAuthService creates a new authentication service
func NewAuthService() *AuthService {
	return &AuthService{
		sessions: make(map[string]*UserSession),
		users:    make(map[string]*AuthUser),
		jwtKey:   []byte("your-secret-key-change-in-production"),
	}
}

// ============================================================================
// Login / Authentication
// ============================================================================

// Login authenticates a user and generates a token
func (as *AuthService) Login(phone, pin string) (*AuthResponse, error) {
	as.mu.RLock()
	user, exists := as.users[phone]
	as.mu.RUnlock()

	if !exists {
		return &AuthResponse{
			Success: false,
			Message: "User not found",
		}, fmt.Errorf("user not found: %s", phone)
	}

	// Verify PIN
	if user.PIN != pin {
		return &AuthResponse{
			Success: false,
			Message: "Invalid PIN",
		}, fmt.Errorf("invalid PIN for user: %s", phone)
	}

	// Generate token
	token := generateToken()

	// Create session
	now := time.Now().Unix()
	expiresAt := now + (24 * 60 * 60) // 24 hours

	session := &UserSession{
		Token:     token,
		UserID:    user.UserID,
		Phone:     user.Phone,
		Name:      user.Name,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}

	as.mu.Lock()
	as.sessions[token] = session
	as.mu.Unlock()

	log.Printf("User logged in: %s (%s)", user.Name, phone)

	return &AuthResponse{
		Success: true,
		Token:   token,
		UserID:  user.UserID,
		Name:    user.Name,
		Message: "Login successful",
	}, nil
}

// ValidateToken verifies a token and returns the session
func (as *AuthService) ValidateToken(token string) (*UserSession, error) {
	as.mu.RLock()
	session, exists := as.sessions[token]
	as.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("token not found")
	}

	// Check expiration
	if time.Now().Unix() > session.ExpiresAt {
		as.mu.Lock()
		delete(as.sessions, token)
		as.mu.Unlock()
		return nil, fmt.Errorf("token expired")
	}

	return session, nil
}

// Logout revokes a token
func (as *AuthService) Logout(token string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if session, exists := as.sessions[token]; exists {
		delete(as.sessions, token)
		log.Printf("User logged out: %s", session.Phone)
		return nil
	}

	return fmt.Errorf("token not found")
}

// ============================================================================
// User Registration & Sample Data
// ============================================================================

// RegisterUser creates a new user with initial sample data
func (as *AuthService) RegisterUser(phone, pin, name, email string) (*UserSeedResponse, error) {
	as.mu.Lock()
	defer as.mu.Unlock()

	// Check if user already exists
	if _, exists := as.users[phone]; exists {
		return &UserSeedResponse{
			Success: false,
			Message: "User already exists",
		}, fmt.Errorf("user already exists: %s", phone)
	}

	// Generate user ID
	userID := fmt.Sprintf("user_%s", randomString(8))

	// Create user
	user := &AuthUser{
		UserID:      userID,
		Phone:       phone,
		PIN:         pin,
		Name:        name,
		Email:       email,
		Address:     generateAddress(),
		DateOfBirth: generateDateOfBirth(),
		SSN:         generateSSN(),
		Accounts:    []string{},
		CreatedAt:   time.Now().Unix(),
	}

	// Create sample accounts with random data
	accounts := []SampleAccountData{}

	// Checking account
	checkingID := fmt.Sprintf("checking_%s", randomString(6))
	user.Accounts = append(user.Accounts, checkingID)
	accounts = append(accounts, SampleAccountData{
		AccountID:    checkingID,
		Type:         "Checking",
		Balance:      generateBalance(1000, 10000),
		Transactions: randomInt(5, 20),
	})

	// Savings account
	savingsID := fmt.Sprintf("savings_%s", randomString(6))
	user.Accounts = append(user.Accounts, savingsID)
	accounts = append(accounts, SampleAccountData{
		AccountID:    savingsID,
		Type:         "Savings",
		Balance:      generateBalance(5000, 50000),
		Transactions: randomInt(2, 10),
	})

	// Money Market account (optional)
	if randomInt(0, 1) == 1 {
		marketID := fmt.Sprintf("market_%s", randomString(6))
		user.Accounts = append(user.Accounts, marketID)
		accounts = append(accounts, SampleAccountData{
			AccountID:    marketID,
			Type:         "Money Market",
			Balance:      generateBalance(10000, 100000),
			Transactions: randomInt(1, 5),
		})
	}

	// Store user
	as.users[phone] = user

	log.Printf("New user registered: %s (%s) with %d accounts", name, phone, len(user.Accounts))

	return &UserSeedResponse{
		Success:   true,
		UserID:    userID,
		Phone:     phone,
		Name:      name,
		Accounts:  accounts,
		Message:   fmt.Sprintf("User %s registered successfully with %d accounts", name, len(user.Accounts)),
	}, nil
}

// GetUser retrieves user information by phone
func (as *AuthService) GetUser(phone string) (*AuthUser, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	user, exists := as.users[phone]
	if !exists {
		return nil, fmt.Errorf("user not found: %s", phone)
	}

	return user, nil
}

// ListUsers returns all registered users
func (as *AuthService) ListUsers() []*AuthUser {
	as.mu.RLock()
	defer as.mu.RUnlock()

	users := make([]*AuthUser, 0, len(as.users))
	for _, user := range as.users {
		users = append(users, user)
	}

	return users
}

// ============================================================================
// Helper Functions
// ============================================================================

// generateToken creates a random token
func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// randomString generates a random string of n characters
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		b[i] = letters[num.Int64()]
	}
	return string(b)
}

// randomInt generates a random integer between min and max
func randomInt(min, max int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return int(num.Int64()) + min
}

// generateBalance generates a random balance between min and max with cents
func generateBalance(min, max int) float64 {
	cents := randomInt(0, 99)
	whole := randomInt(min, max)
	return float64(whole) + float64(cents)/100
}

// generateAddress generates a random US address
func generateAddress() string {
	streets := []string{"Main St", "Oak Ave", "Maple Drive", "Elm Street", "Pine Road"}
	cities := []string{"Springfield", "Riverside", "Brookside", "Willowville", "Sunnydale"}
	states := []string{"CA", "NY", "TX", "FL", "IL"}

	street := randomInt(100, 9999)
	streetName := streets[randomInt(0, len(streets)-1)]
	city := cities[randomInt(0, len(cities)-1)]
	state := states[randomInt(0, len(states)-1)]
	zip := randomInt(10000, 99999)

	return fmt.Sprintf("%d %s, %s, %s %d", street, streetName, city, state, zip)
}

// generateDateOfBirth generates a random date of birth (18-75 years old)
func generateDateOfBirth() string {
	now := time.Now()
	age := randomInt(18, 75)
	year := now.Year() - age
	month := randomInt(1, 12)
	day := randomInt(1, 28)

	return fmt.Sprintf("%04d-%02d-%02d", year, month, day)
}

// generateSSN generates a random SSN (XXX-XX-XXXX format)
func generateSSN() string {
	return fmt.Sprintf("%03d-%02d-%04d",
		randomInt(0, 999),
		randomInt(0, 99),
		randomInt(0, 9999),
	)
}

// ============================================================================
// Sample Users
// ============================================================================

// CreateSampleUsers creates a set of demo users for testing
func (as *AuthService) CreateSampleUsers() error {
	sampleUsers := []struct {
		phone string
		pin   string
		name  string
		email string
	}{
		{"+14155552671", "1234", "Alice Johnson", "alice@example.com"},
		{"+14155552672", "5678", "Bob Smith", "bob@example.com"},
		{"+14155552673", "9012", "Carol Williams", "carol@example.com"},
		{"+14155552674", "3456", "David Brown", "david@example.com"},
		{"+14155552675", "7890", "Emma Davis", "emma@example.com"},
	}

	for _, su := range sampleUsers {
		_, err := as.RegisterUser(su.phone, su.pin, su.name, su.email)
		if err != nil {
			log.Printf("Warning: Could not create sample user %s: %v", su.name, err)
		}
	}

	log.Printf("Created %d sample users", len(sampleUsers))
	return nil
}
