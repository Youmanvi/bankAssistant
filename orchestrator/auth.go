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
	mu   sync.RWMutex
	db   *Database
	jwtKey []byte
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
func NewAuthService(db *Database) *AuthService {
	return &AuthService{
		db:     db,
		jwtKey: []byte("your-secret-key-change-in-production"),
	}
}

// ============================================================================
// Login / Authentication
// ============================================================================

// Login authenticates a user and generates a token
func (as *AuthService) Login(phone, pin string) (*AuthResponse, error) {
	// Get user from database
	user, err := as.db.GetUserByPhone(phone)
	if err != nil {
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

	// Store session in database
	if err := as.db.CreateSession(session); err != nil {
		return &AuthResponse{
			Success: false,
			Message: "Failed to create session",
		}, fmt.Errorf("failed to create session: %w", err)
	}

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
	// Get session from database
	session, err := as.db.GetSessionByToken(token)
	if err != nil {
		return nil, fmt.Errorf("token not found")
	}

	// Check expiration
	if time.Now().Unix() > session.ExpiresAt {
		as.db.DeleteSession(token)
		return nil, fmt.Errorf("token expired")
	}

	// Fetch user info to populate session fields
	user, err := as.db.GetUserByID(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	session.Phone = user.Phone
	session.Name = user.Name

	return session, nil
}

// Logout revokes a token
func (as *AuthService) Logout(token string) error {
	// Get session to log user info
	session, err := as.db.GetSessionByToken(token)
	if err != nil {
		return fmt.Errorf("token not found")
	}

	// Delete session from database
	if err := as.db.DeleteSession(token); err != nil {
		return err
	}

	log.Printf("User logged out: %s", session.UserID)
	return nil
}

// ============================================================================
// User Registration & Sample Data
// ============================================================================

// RegisterUser creates a new user with initial sample data
func (as *AuthService) RegisterUser(phone, pin, name, email string) (*UserSeedResponse, error) {
	// Check if user already exists
	_, err := as.db.GetUserByPhone(phone)
	if err == nil {
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

	// Store user in database
	if err := as.db.CreateUser(user); err != nil {
		return &UserSeedResponse{
			Success: false,
			Message: "Failed to create user",
		}, fmt.Errorf("failed to create user: %w", err)
	}

	// Create sample accounts with random data
	accounts := []SampleAccountData{}

	// Checking account
	checkingID := fmt.Sprintf("checking_%s", randomString(6))
	checkingBalance := generateBalance(1000, 10000)
	user.Accounts = append(user.Accounts, checkingID)

	if err := as.db.CreateAccount(checkingID, userID, "Checking", checkingBalance); err != nil {
		log.Printf("Error creating checking account: %v", err)
	}

	accounts = append(accounts, SampleAccountData{
		AccountID:    checkingID,
		Type:         "Checking",
		Balance:      checkingBalance,
		Transactions: randomInt(5, 20),
	})

	// Savings account
	savingsID := fmt.Sprintf("savings_%s", randomString(6))
	savingsBalance := generateBalance(5000, 50000)
	user.Accounts = append(user.Accounts, savingsID)

	if err := as.db.CreateAccount(savingsID, userID, "Savings", savingsBalance); err != nil {
		log.Printf("Error creating savings account: %v", err)
	}

	accounts = append(accounts, SampleAccountData{
		AccountID:    savingsID,
		Type:         "Savings",
		Balance:      savingsBalance,
		Transactions: randomInt(2, 10),
	})

	// Money Market account (optional)
	if randomInt(0, 1) == 1 {
		marketID := fmt.Sprintf("market_%s", randomString(6))
		marketBalance := generateBalance(10000, 100000)
		user.Accounts = append(user.Accounts, marketID)

		if err := as.db.CreateAccount(marketID, userID, "Money Market", marketBalance); err != nil {
			log.Printf("Error creating money market account: %v", err)
		}

		accounts = append(accounts, SampleAccountData{
			AccountID:    marketID,
			Type:         "Money Market",
			Balance:      marketBalance,
			Transactions: randomInt(1, 5),
		})
	}

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
	return as.db.GetUserByPhone(phone)
}

// ListUsers returns all registered users
func (as *AuthService) ListUsers() []*AuthUser {
	users, err := as.db.ListAllUsers()
	if err != nil {
		log.Printf("Error listing users: %v", err)
		return []*AuthUser{}
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
