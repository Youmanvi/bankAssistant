# Go Orchestrator

The Go orchestrator manages the call lifecycle and orchestrates communication between Retell AI (voice provider) and the Python FastAPI backend (business logic).

## Overview

```
Retell AI Webhook → Go Orchestrator → Python FastAPI Backend
    (call_started)   (state machine)    (REST API calls)
    (call_ended)    (load context)     (user info, accounts)
    (call_analyzed) (route requests)   (execute operations)
```

## Architecture

### Call State Machine

The orchestrator manages calls through 8 states:

```
AWAITING_CALL
    ↓
CALL_STARTED (Retell webhook fires)
    ↓
AWAITING_INTENT (Load user context, ready for input)
    ↓
PROCESSING_REQUEST (Routing user intent to Python backend)
    ↓
GENERATING_RESPONSE (Python LLM generates response)
    ↓
SPEAKING_RESPONSE (Retell AI speaks response to user)
    ↓
AWAITING_INTENT (Loop for follow-up) or CALL_ENDED
    ↓
AWAITING_CALL (Clean up and reset)
```

### Components

1. **Configuration** (`config.go`)
   - Load settings from `.env` file
   - Server host/port
   - Retell API key
   - Python backend URL

2. **Types** (`types.go`)
   - Call states and state machine types
   - Retell webhook payloads
   - Python backend models
   - HTTP response models

3. **Call State Machine** (`call_state_machine.go`)
   - Create, update, and track call states
   - Validate state transitions
   - Store call metadata
   - Concurrent-safe operations (mutex-protected)

4. **Python Backend Client** (`python_client.go`)
   - HTTP client for REST API calls to Python backend
   - Retry logic for failed requests
   - Load user context
   - Get account balance
   - Execute fund transfers

5. **Retell Handler** (`retell_handler.go`)
   - Receive and validate Retell webhooks
   - Verify webhook signatures (HMAC-SHA256)
   - Process call lifecycle events
   - Load user context from Python backend
   - Admin endpoints for monitoring calls

6. **Orchestrator** (`orchestrator.go`)
   - Main orchestration service
   - HTTP route handlers
   - Health checks
   - Admin endpoints
   - Orchestration endpoints for business logic

## Setup

### Prerequisites

- **Go 1.21+** - [Download](https://golang.org/dl/)
- **PostgreSQL 12+** - [Download](https://www.postgresql.org/download/)
- **Python Backend** - Running on `http://localhost:8000`
- **Retell AI Account** - For API key
- **Environment Configuration** - `.env` file with API keys

### Database Setup (PostgreSQL)

The orchestrator uses PostgreSQL to store users, sessions, and accounts.

**On macOS:**
```bash
# Install PostgreSQL
brew install postgresql@15

# Start PostgreSQL service
brew services start postgresql@15

# Create database and user
psql -U postgres -f orchestrator/scripts_setup_db.sql
```

**On Ubuntu/Debian:**
```bash
# Install PostgreSQL
sudo apt-get install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql

# Create database and user (as postgres user)
sudo -u postgres psql -f orchestrator/scripts_setup_db.sql
```

**On Windows:**
```bash
# Download and install PostgreSQL from https://www.postgresql.org/download/windows/

# Run setup script
"C:\Program Files\PostgreSQL\15\bin\psql.exe" -U postgres -f orchestrator\scripts_setup_db.sql
```

**Note:** The application will automatically create all necessary tables on first run.

### Quick Start

```bash
# 1. Navigate to orchestrator directory
cd orchestrator

# 2. Create environment file
cp .env.example .env

# 3. Update .env with your configuration
# - Set Retell API key
# - Set PostgreSQL password (should match the one in setup script)
# - Configure other settings as needed

# 4. Download dependencies
go mod download

# 5. Run the orchestrator
go run .

# Expected output:
# Connected to PostgreSQL database: bankassistant
# Database schema initialized successfully
# Created 5 sample users
# Orchestrator starting on 0.0.0.0:8001
# Python backend is healthy
```

### Verify Installation

```bash
# In another terminal, check health
curl http://localhost:8001/health

# Expected response
{"status": "healthy", "version": "0.1.0"}
```

## Environment Variables

```env
# Server Configuration
ORCHESTRATOR_HOST=0.0.0.0          # Listen on all interfaces
ORCHESTRATOR_PORT=8001              # Server port

# Retell AI
RETELL_API_KEY=your_api_key        # From https://retellai.com/dashboard

# Python Backend
PYTHON_BACKEND_URL=http://localhost:8000

# Logging
LOG_LEVEL=INFO                      # DEBUG, INFO, WARN, ERROR

# State Storage
CALL_STATE_DB=memory                # "memory" (dev) or "redis" (prod)
```

## API Endpoints

### Retell Webhook

**Endpoint:** `POST /webhook`

Receives Retell AI webhook events:
- `call_started` - User initiates a call
- `call_ended` - Call terminates
- `call_analyzed` - Retell AI provides analysis

**Example:**
```bash
# Retell AI sends this when a call starts
POST http://localhost:8001/webhook

{
  "event": "call_started",
  "data": {
    "call_id": "call_123456",
    "phone_number": "+1234567890"
  }
}
```

### Load User Context

**Endpoint:** `POST /orchestrate/load-context`

Loads user information for a call.

**Request:**
```bash
curl -X POST http://localhost:8001/orchestrate/load-context \
  -H "Content-Type: application/json" \
  -d '{
    "call_id": "call_123456",
    "user_id": "user_1"
  }'
```

**Response:**
```json
{"status": "context loaded"}
```

### Get Account Balance

**Endpoint:** `POST /orchestrate/get-balance`

Retrieves account balance from Python backend.

**Request:**
```bash
curl -X POST http://localhost:8001/orchestrate/get-balance \
  -H "Content-Type: application/json" \
  -d '{
    "call_id": "call_123456",
    "account_id": "checking_1"
  }'
```

**Response:**
```json
{
  "account_id": "checking_1",
  "balance": 5000.50
}
```

### Transfer Funds

**Endpoint:** `POST /orchestrate/transfer`

Initiates a fund transfer.

**Request:**
```bash
curl -X POST http://localhost:8001/orchestrate/transfer \
  -H "Content-Type: application/json" \
  -d '{
    "call_id": "call_123456",
    "from_account": "checking_1",
    "to_account": "savings_1",
    "amount": 500.00
  }'
```

**Response:**
```json
{
  "id": "payment_789",
  "from_account": "checking_1",
  "to_account": "savings_1",
  "amount": 500.00,
  "status": "Completed"
}
```

### Admin Endpoints

**Get All Calls**
```bash
curl http://localhost:8001/admin/calls
```

**Get Call Status**
```bash
curl "http://localhost:8001/admin/call?call_id=call_123456"
```

## Development

### Code Structure

```
orchestrator/
├── main.go                 # Entry point
├── config.go              # Configuration loading
├── types.go               # Type definitions
├── call_state_machine.go  # Call lifecycle management
├── python_client.go       # Python backend HTTP client
├── retell_handler.go      # Retell webhook handling
├── orchestrator.go        # Main orchestration service
├── go.mod                 # Go module definition
├── go.sum                 # Dependency checksums
├── .env.example           # Environment template
├── Dockerfile             # Container image
└── README.md              # This file
```

### Building Locally

```bash
# Build binary
go build -o orchestrator

# Run binary
./orchestrator

# Or run directly
go run .
```

### Testing

```bash
# Run tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestCallStateMachine ./...
```

### Docker

**Build image:**
```bash
docker build -t bankassistant-orchestrator .
```

**Run container:**
```bash
docker run -p 8001:8001 \
  -e RETELL_API_KEY=your_key \
  -e PYTHON_BACKEND_URL=http://host.docker.internal:8000 \
  bankassistant-orchestrator
```

**With docker-compose:**
```bash
# In root bankAssistant directory
docker-compose up orchestrator
```

## Integration with Python Backend

The orchestrator calls the Python FastAPI backend REST APIs:

```
GET  /api/v1/users/{user_id}
GET  /api/v1/users/{user_id}/accounts
GET  /api/v1/accounts/{account_id}
GET  /api/v1/accounts/{account_id}/statements
POST /api/v1/payments/transfer
```

### Example Call Flow

```
1. User calls banking assistant → Retell AI
2. Retell sends call_started webhook → Go Orchestrator
3. Orchestrator creates call context with state AWAITING_CALL → CALL_STARTED
4. User authenticates (phone-based)
5. Orchestrator loads user context from Python backend
6. Orchestrator transitions state to AWAITING_INTENT
7. User says "Check my balance"
8. Orchestrator routes to Python backend
9. Python backend executes agents and returns response
10. Orchestrator transitions state to SPEAKING_RESPONSE
11. Retell AI speaks response to user
12. User hears response
13. Call ends → Retell sends call_ended webhook
14. Orchestrator transitions state to CALL_ENDED
```

## Common Development Tasks

### Add New HTTP Endpoint

1. Add handler function in `orchestrator.go`
2. Register route in `setupRoutes()`
3. Test endpoint with curl

Example:
```go
// Add handler
func (o *Orchestrator) handleNewFeature(w http.ResponseWriter, r *http.Request) {
    // Implementation
}

// Register in setupRoutes()
http.HandleFunc("/orchestrate/new-feature", o.handleNewFeature)
```

### Add Retry Logic

The `python_client.go` already includes retry logic. To adjust:

```go
client := NewPythonBackendClient(baseURL)
client.retries = 5  // Increase from 3 to 5
```

### Debug State Transitions

Enable debug logging:
```bash
LOG_LEVEL=DEBUG go run .
```

## Production Deployment

### Environment Variables

```env
# Use actual hostnames
ORCHESTRATOR_HOST=your.orchestrator.host
ORCHESTRATOR_PORT=8001

# Use production Retell key
RETELL_API_KEY=prod_key

# Use production Python backend
PYTHON_BACKEND_URL=https://api.bankassistant.com

# Use Redis for state storage
CALL_STATE_DB=redis
```

### Docker Compose

```yaml
version: '3'
services:
  orchestrator:
    build: ./orchestrator
    ports:
      - "8001:8001"
    env_file: .env
    depends_on:
      - backend
    networks:
      - bankassistant

  backend:
    build: ./backend
    ports:
      - "8000:8000"
    env_file: .env
    networks:
      - bankassistant

networks:
  bankassistant:
    driver: bridge
```

## Troubleshooting

### "Connection refused" error

**Problem:** Cannot connect to Python backend
```
Error: request failed: connection refused
```

**Solution:**
```bash
# Verify Python backend is running
curl http://localhost:8000/health

# If not running, start it
cd backend
python main.py
```

### "RETELL_API_KEY not found" error

**Solution:**
```bash
# Create .env file
cp .env.example .env

# Edit with your API key
nano .env
# or
notepad .env  # Windows
```

### "Port 8001 already in use"

**Solution:**
```bash
# Change port
export ORCHESTRATOR_PORT=8002  # macOS/Linux
set ORCHESTRATOR_PORT=8002     # Windows

go run .
```

## Performance Tips

1. **Concurrent Calls**
   - State machine is mutex-protected for safe concurrent access
   - Go's goroutines handle multiple calls efficiently

2. **Retry Logic**
   - Python client has exponential backoff retry logic
   - Adjust `retries` value based on backend stability

3. **Monitoring**
   - Use `/admin/calls` endpoint to monitor active calls
   - Log all state transitions for debugging

## Security

1. **Webhook Signature Verification**
   - All Retell webhooks are HMAC-SHA256 verified
   - Signature in `X-Retell-Signature` header

2. **API Key Management**
   - Store `RETELL_API_KEY` in `.env` (not in code)
   - Never commit `.env` to git
   - Use different keys for development/production

3. **CORS**
   - Configure CORS in Python backend, not orchestrator
   - Orchestrator acts as internal service only

## Next Steps

1. **[Set up Python backend](../backend/README.md)** if not already done
2. **Get Retell API key** from https://retellai.com/dashboard
3. **Configure `.env`** with your API keys
4. **Run orchestrator** locally for development
5. **Test webhook flow** using Postman or curl
6. **Monitor calls** using admin endpoints
7. **Deploy** using Docker Compose

## Support

- Check logs: `LOG_LEVEL=DEBUG go run .`
- Read Python backend docs: `../backend/README.md`
- Review Retell integration: `https://docs.retellai.com`
