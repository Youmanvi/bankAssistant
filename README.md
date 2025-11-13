# AI Banking Assistant

A voice-driven, multi-agent AI banking assistant that handles customer interactions through natural conversation. Powered by FastAPI, OpenAI Swarm, Retell AI, and built for eventual Go orchestration.

## Quick Overview

```
📞 User input → Retell AI (voice provider)
                ↓
            Python FastAPI Backend (handles the intelligence)
                ├─ REST APIs (for Go orchestrator)
                ├─ WebSocket (Retell AI voice stream)
                ├─ LLM interaction (OpenAI)
                └─ Multi-agent system (Swarm)
                ↓
            Go Orchestrator (to be built - manages workflow)
                ↓
            Next.js Admin Dashboard (monitor calls)
```

## Getting Started

### Prerequisites
- **Python 3.10+** (backend)
- **Node.js 18+** (frontend)
- **API Keys**: OpenAI, Retell AI, Pinata (optional)

### Quick Start (10 minutes)

**Terminal 1 - Python Backend:**
```bash
cd backend
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt
cp .env.example .env
# Edit .env with your API keys
python main.py
# Backend runs on http://localhost:8000
```

**Terminal 2 - Go Orchestrator:**
```bash
cd orchestrator
cp .env.example .env
# Edit .env with your API keys
go run .
# Orchestrator runs on http://localhost:8001
```

**Terminal 3 - Next.js Frontend:**
```bash
cd frontend
npm install
npm run dev
# Frontend runs on http://localhost:3000
```

**Verify Installation:**
```bash
curl http://localhost:8000/health    # Backend
curl http://localhost:8001/health    # Orchestrator
```

**All components running!** You can now test the full system.

## How It Works

### Call Lifecycle

1. **User Calls**: User dials the banking assistant number
2. **Retell AI Answers**: Retell AI voice provider receives call
3. **Load Context**: Go orchestrator loads user information (name, accounts)
4. **User Speaks**: User says what they want ("Check my balance")
5. **Transcription**: Retell AI converts speech to text
6. **LLM Processing**: Python backend processes intent using OpenAI
7. **Agent Execution**: Specialized agent (Accounts, Payments, etc.) executes operation
8. **Database Query**: Fetch account balance or execute transfer
9. **Response Generation**: LLM formats natural language response
10. **Text-to-Speech**: Retell AI converts response to voice
11. **User Hears**: Call user hears response
12. **Call Ends**: Call terminates, logged for audit trail

### Agent System

The multi-agent system uses OpenAI Swarm framework:

**Triage Agent** (Entry Point)
- Analyzes user intent
- Routes to appropriate specialist agent

**Accounts Agent**
- Check balance: `check_balance(account_id)`
- View statements: `get_statements(account_id, month)`
- List accounts: `list_accounts(user_id)`

**Payments Agent**
- Transfer funds: `transfer_funds(from_account, to_account, amount)`
- Schedule payment: `schedule_payment(account, amount, date)`
- Cancel payment: `cancel_payment(payment_id)`

**Applications Agent**
- Apply for loan: `apply_for_loan(user_id, amount, term_months)`
- Apply for credit card: `apply_for_credit_card(user_id, limit)`

### Key Components Explained

#### Python FastAPI Backend
Handles all the intelligence:
- **REST APIs**: Stateless endpoints for Go orchestrator
- **WebSocket**: Real-time connection with Retell AI for voice
- **LLM Service**: Queries OpenAI GPT-4 for understanding and generation
- **Agent System**: Multi-agent orchestration using Swarm
- **Database**: Mock in-memory database (easily replaceable)

#### Go Orchestrator (To Be Built)
Manages the overall workflow:
- **Call Routing**: Routes Retell webhooks appropriately
- **State Management**: Tracks call through its lifecycle
- **User Context**: Loads and maintains user information during call
- **API Orchestration**: Calls Python APIs in correct sequence
- **Error Handling**: Retries failed operations, handles edge cases

#### Retell AI (Third-Party)
Provides voice capabilities:
- **Phone Integration**: Receives incoming calls
- **Speech-to-Text**: Transcribes user speech to text
- **Text-to-Speech**: Converts responses back to voice
- **WebSocket Stream**: Bidirectional communication with Python

#### Next.js Frontend
Admin dashboard for monitoring:
- **Call Logs**: View all calls and interactions
- **User Management**: See user accounts and balances
- **Real-time Updates**: WebSocket connection for live updates
- **Analytics**: Call metrics and system health

## Project Structure

```
bankAssistant/
├── backend/                    # FastAPI Python service
│   ├── api/                    # REST endpoints (Go calls these)
│   ├── ws/                     # WebSocket handlers (voice)
│   ├── services/               # Business logic (LLM, agents, storage)
│   ├── agents/                 # AI agents (Swarm framework)
│   ├── db/                     # Database operations
│   ├── utils/                  # Helper functions
│   ├── models.py               # Pydantic data models
│   ├── constants.py            # Configuration & enums
│   ├── config.py               # Environment config
│   ├── main.py                 # FastAPI app entry point
│   └── requirements.txt        # Python dependencies
│
├── orchestrator/               # Go orchestrator service
│   ├── main.go                 # Entry point
│   ├── config.go               # Configuration loading
│   ├── types.go                # Type definitions
│   ├── call_state_machine.go   # Call lifecycle management
│   ├── python_client.go        # Python backend HTTP client
│   ├── retell_handler.go       # Retell webhook handling
│   ├── orchestrator.go         # Main orchestration service
│   ├── go.mod                  # Go module definition
│   ├── .env.example            # Environment template
│   ├── Dockerfile              # Container image
│   └── README.md               # Orchestrator guide
│
├── frontend/                   # Next.js React dashboard
│   ├── src/app/                # Next.js app
│   ├── src/components/         # React components
│   └── package.json
│
├── docs/                       # Technical documentation
│   ├── SYSTEM.md               # Complete system architecture
│   ├── STRUCTURE.md            # Detailed project organization
│   ├── API.md                  # API endpoint reference
│   ├── SETUP.md                # Local development setup
│   ├── CONVENTIONS.md          # Code standards
│   ├── GO_ORCHESTRATOR.md      # Go orchestrator guide
│   └── GO_EXAMPLES.md          # Go code patterns
│
├── scripts/                    # Utility scripts & examples
├── tests/                      # Test suite (unit & integration)
├── .env.example                # Environment template
├── .gitignore                  # Git ignore rules
├── docker-compose.yml          # Multi-service orchestration
└── README.md                   # This file (comprehensive guide)
```

### Directory Details

**backend/** - FastAPI Python service
- `api/` - REST endpoints (GET/POST operations)
- `ws/` - WebSocket handlers (voice stream)
- `services/` - Business logic (LLM, agents, storage)
- `agents/` - AI agents using OpenAI Swarm framework
- `db/` - Database operations and CRUD methods
- `utils/` - Helper functions
- `models.py` - Pydantic data models for validation
- `constants.py` - Configuration constants and enums
- `config.py` - Environment variable management
- `main.py` - FastAPI application entry point

**orchestrator/** - Go orchestrator service
- `main.go` - Entry point and signal handling
- `config.go` - Configuration loading from environment
- `types.go` - Type definitions (CallState, RetellPayload, etc.)
- `call_state_machine.go` - Call lifecycle state management (8 states)
- `python_client.go` - HTTP client for Python backend REST APIs
- `retell_handler.go` - Retell webhook receiver and processor
- `orchestrator.go` - Main orchestration service and HTTP handlers
- Complete Dockerfile for containerization
- See `orchestrator/README.md` for full documentation

**frontend/** - Next.js React dashboard
- Interactive admin interface for monitoring calls
- Real-time updates via WebSocket
- User and account management views
- Call logs and analytics

**docs/** - Comprehensive technical documentation
- Architecture and system design
- Project organization and code structure
- API reference with examples
- Development setup instructions
- Code conventions and standards
- Go orchestrator building guide
- Go code examples and patterns

## Technology Stack

### Backend
- **FastAPI** - Modern async Python web framework
- **OpenAI Swarm** - Multi-agent orchestration
- **Pydantic** - Data validation
- **Python 3.10+** - Runtime

### Orchestrator
- **Go 1.21+** - High-performance statically-typed language
- **Standard Library** - net/http for HTTP server and client
- **godotenv** - Environment variable management
- Concurrent-safe state machine with mutexes

### Frontend
- **Next.js** - React framework
- **TypeScript** - Type safety
- **Tailwind CSS** - Styling

### External Services
- **Retell AI** - Voice processing & phone integration
- **OpenAI** - LLM (GPT-4)
- **Pinata** - IPFS/decentralized storage

## Architecture at a Glance

### Components

| Component | Role | Technology |
|-----------|------|-----------|
| **Retell AI** | Voice provider | Third-party service |
| **Go Orchestrator** | Webhook handler, call state machine, API orchestration | Go 1.21+ |
| **Python Backend** | REST API + LLM + agents | FastAPI + OpenAI Swarm |
| **Next.js Frontend** | Admin dashboard | React + TypeScript |

## REST API Endpoints

All API endpoints follow the `/api/v1/` prefix and return JSON responses.

### Health & Status
```bash
GET /health
# Response: {"status": "healthy", "version": "1.0.0"}

GET /status
# Response: {"status": "ok", "uptime": "2h 30m"}
```

### User Endpoints
```bash
GET /api/v1/users/{user_id}
# Get user information and linked accounts
# Response: {
#   "user_id": "USR001",
#   "name": "John Doe",
#   "phone": "+14155552671",
#   "email": "john@example.com",
#   "accounts": ["ACC001", "ACC002"]
# }
```

### Account Endpoints
```bash
GET /api/v1/accounts/{account_id}
# Get account details and balance
# Response: {
#   "account_id": "ACC001",
#   "type": "CHECKING",
#   "balance": 2500.00,
#   "statements": { "2024-11": [...] }
# }
```

### Payment Endpoints
```bash
POST /api/v1/payments/transfer
# Execute a fund transfer
# Request: {
#   "from_account": "ACC001",
#   "to_account": "ACC002",
#   "amount": 50.00
# }
# Response: {
#   "payment_id": "PAY001",
#   "status": "COMPLETED",
#   "amount": 50.00
# }
```

### Application Endpoints
```bash
POST /api/v1/applications/loan
# Submit a loan application
# Request: {
#   "user_id": "USR001",
#   "amount": 10000.00,
#   "term_months": 12
# }
# Response: {
#   "application_id": "APP001",
#   "status": "PENDING"
# }

POST /api/v1/applications/credit-card
# Submit a credit card application
# Request: {
#   "user_id": "USR001",
#   "credit_limit": 5000.00
# }
# Response: {
#   "application_id": "APP002",
#   "status": "PENDING"
# }
```

### WebSocket Endpoints
```
WS /llm-websocket/{call_id}
# Voice interaction with Retell AI (Python manages internally)

WS /admin/ws?client_id={id}
# Admin dashboard real-time updates
```

For detailed request/response formats and examples, see `docs/API.md`.

## Development Workflow

### Running Tests
```bash
cd backend
pytest tests/
```

### Running Backend in Debug Mode
```bash
cd backend
export DEBUG=true
python main.py
```

### Frontend Development
```bash
cd frontend
npm run dev
# Hot reload enabled
```

### Making Code Changes

1. **Naming & Style**:
   - Functions: snake_case (e.g., `process_transfer()`)
   - Classes: PascalCase (e.g., `TransferRequest`)
   - Constants: UPPER_SNAKE_CASE (e.g., `MAX_RETRIES`)
   - Max line length: 100 characters
   - Use type hints on all functions

2. **Code Organization**:
   - REST endpoints go in `backend/api/`
   - Business logic goes in `backend/services/`
   - AI agents go in `backend/agents/`
   - Database operations go in `backend/db/`
   - Data models go in `backend/models.py`

3. **Testing & Verification**:
   ```bash
   # Run tests
   cd backend
   pytest tests/

   # Check code style
   python -m py_compile backend/*.py

   # Type checking (optional)
   mypy backend/
   ```

4. **Commit Messages**:
   - Use clear, descriptive messages
   - Format: `<type>(<scope>): <subject>`
   - Types: feat, fix, docs, refactor, test

## Common Development Tasks

### Adding a New REST Endpoint

1. Create route in `backend/api/new_endpoint.py`:
   ```python
   from fastapi import APIRouter
   from models import Request, Response

   router = APIRouter(prefix="/new", tags=["new"])

   @router.post("/operation", response_model=Response)
   async def new_operation(request: Request) -> Response:
       """Operation description."""
       # Implementation
       pass
   ```

2. Define models in `backend/models.py`:
   ```python
   class Request(BaseModel):
       field: str
   ```

3. Register in `backend/main.py`:
   ```python
   from api import new_endpoint
   app.include_router(new_endpoint.router, prefix="/api/v1")
   ```

4. Test with:
   ```bash
   curl -X POST http://localhost:8000/api/v1/new/operation
   ```

### Adding a New Agent

1. Create `backend/agents/new_agent.py`:
   ```python
   from openai import Swarm

   class NewAgent:
       def __init__(self):
           self.name = "NewAgent"
           self.functions = [self.do_something]

       async def do_something(self, param: str) -> str:
           """Agent function."""
           return result
   ```

2. Register in `backend/services/agents.py`:
   ```python
   from agents.new_agent import NewAgent
   new_agent = NewAgent()
   ```

### Adding to the Database

1. Update model in `backend/models.py`
2. Add CRUD in `backend/db/service.py`
3. Add seed data in `backend/db/seed.py`

## Configuration

### Environment Variables

Copy `.env.example` to `.env` and fill in your API keys:

```bash
# OpenAI (Required)
OPENAI_API_KEY=sk_test_your_key_here
OPENAI_MODEL=gpt-4

# Retell AI (Required)
RETELL_API_KEY=your_retell_key_here
RETELL_API_BASE=https://api.retellai.com

# Pinata/IPFS (Optional)
PINATA_API_KEY=your_key_here
PINATA_SECRET_KEY=your_secret_here

# Server Configuration
APP_HOST=0.0.0.0
APP_PORT=8000
DEBUG=true
LOG_LEVEL=INFO

# Frontend
FRONTEND_URL=http://localhost:3000

# Database (Future)
# DATABASE_URL=postgresql://user:password@localhost/bankassistant

# Redis (Optional)
# REDIS_URL=redis://localhost:6379/0
```

### Getting API Keys

**OpenAI API Key:**
1. Go to https://platform.openai.com/api-keys
2. Create new secret key
3. Copy to `OPENAI_API_KEY`

**Retell AI API Key:**
1. Go to https://retellai.com/dashboard
2. Create new API key
3. Copy to `RETELL_API_KEY`

**Pinata API Key (Optional):**
1. Go to https://pinata.cloud
2. Create new API key
3. Copy to `PINATA_API_KEY` and `PINATA_SECRET_KEY`

## Supported Endpoints

All endpoints are fully functional and stateless REST APIs:

```
✓ GET  /health                        Health check
✓ GET  /api/v1/users/{user_id}       Get user data
✓ GET  /api/v1/accounts/{account_id} Get account info
✓ POST /api/v1/payments/transfer     Execute transfer
✓ POST /api/v1/applications/loan     Submit loan application
✓ POST /api/v1/applications/credit-card  Submit credit card app
✓ WS   /llm-websocket/{call_id}      Voice stream (WebSocket)
✓ WS   /admin/ws                     Admin dashboard updates
```

### Error Responses

When an error occurs, the API returns:
```json
{
  "status": "error",
  "error": "Human-readable error message",
  "details": {
    "field": "error details"
  }
}
```

Common errors:
- **400 Bad Request** - Invalid data in request
- **404 Not Found** - Resource doesn't exist
- **422 Validation Error** - Data validation failed
- **500 Server Error** - Unexpected server issue

For detailed endpoint documentation with examples, see `docs/API.md`.

## Database

### Current: In-Memory Mock Database

The development database is a simple Python dictionary stored in memory:

```python
{
  "users": {
    "+14155552671": {
      "user_id": "USR001",
      "name": "John Doe",
      "accounts": ["ACC001", "ACC002"]
    }
  },
  "accounts": {
    "ACC001": {
      "account_id": "ACC001",
      "type": "CHECKING",
      "balance": 2500.00,
      "statements": { "2024-11": [...] }
    }
  },
  "payments": {
    "PAY001": {
      "payment_id": "PAY001",
      "from_account": "ACC001",
      "to_account": "ACC002",
      "amount": 50.00,
      "status": "COMPLETED"
    }
  }
}
```

### Data Structure

**Users**
- `user_id`: Unique identifier
- `name`: Full name
- `phone`: Phone number (index)
- `email`: Email address
- `accounts`: List of linked account IDs

**Accounts**
- `account_id`: Unique identifier
- `type`: CHECKING, SAVINGS
- `balance`: Current balance
- `statements`: Monthly transaction history

**Payments**
- `payment_id`: Unique identifier
- `from_account`: Source account
- `to_account`: Destination account
- `amount`: Transfer amount
- `status`: PENDING, COMPLETED, FAILED
- `date`: Transaction date

### Future: Production Database

For production deployment, replace the mock database with PostgreSQL:

```python
# backend/db/service.py
import asyncpg

class DatabaseService:
    async def connect(self):
        self.pool = await asyncpg.create_pool(
            os.getenv("DATABASE_URL")
        )

    async def get_user(self, user_id: str):
        async with self.pool.acquire() as conn:
            return await conn.fetchrow(
                "SELECT * FROM users WHERE user_id = $1",
                user_id
            )
```

## Security

### Current (Development/POC)
- ✓ CORS enabled for development (localhost:3000, localhost:8000)
- ✓ Retell webhook signature verification
- ✓ Mock data only (no real accounts)
- ✓ No authentication required (safe for development)

### Production Requirements

**Authentication**
```python
# Add API key authentication for Go orchestrator
from fastapi.security import APIKeyHeader

security = APIKeyHeader(name="X-API-Key")

@app.get("/api/v1/users/{user_id}")
async def get_user(user_id: str, api_key: str = Depends(security)):
    if api_key != os.getenv("API_KEY"):
        raise HTTPException(status_code=403, detail="Invalid API key")
    return await db.get_user(user_id)
```

**Recommended Security Measures**
- [ ] **API Key Authentication**: Go ↔ Python
- [ ] **JWT Tokens**: Frontend session management
- [ ] **HTTPS/TLS**: All communication encrypted
- [ ] **Rate Limiting**: 100 requests per minute per IP
- [ ] **Input Validation**: Pydantic models validate all inputs
- [ ] **Audit Logging**: Log all transactions and changes
- [ ] **PCI-DSS Compliance**: For real banking operations
- [ ] **Secrets Management**: Use AWS Secrets Manager or similar

### Securing API Keys
```bash
# DO NOT commit .env file
echo ".env" >> .gitignore

# Use environment variables in production
export OPENAI_API_KEY="sk_prod_..."
export RETELL_API_KEY="prod_..."

# Rotate keys regularly
# Never log sensitive data
```

## Performance

### Typical Latency Per Component
- **Speech to Text**: ~500ms-2s (Retell AI)
- **LLM Processing**: ~1-3s (OpenAI)
- **Python API Call**: ~50-200ms
- **Text to Speech**: ~0.5-2s (Retell AI)
- **Total End-to-End**: ~3-8 seconds (acceptable for voice)

### Scaling Considerations

**Vertical Scaling (Single Instance)**
- Single Python instance handles ~100 concurrent calls
- Memory per call: ~50MB
- Total baseline: ~100MB + 5MB per agent

**Horizontal Scaling (Multiple Instances)**
```
Load Balancer (Nginx)
├─ Python Service 1 (port 8001)
├─ Python Service 2 (port 8002)
└─ Python Service 3 (port 8003)
```

**Caching Strategy**
```python
# Cache user data for 5 minutes
@cache(ttl=300)
async def get_user(user_id: str):
    return await db.get_user(user_id)
```

**Database Optimization**
- Index on user phone number
- Index on account IDs
- Batch read queries
- Connection pooling (asyncpg)

## Deployment

### Development Environment
```bash
# Local development
Python backend:  http://localhost:8000
React frontend:  http://localhost:3000
Go orchestrator: (not yet running)

# Access logs
tail -f backend/app.log
```

### Staging Environment
```bash
# Docker Compose for multi-service deployment
docker-compose -f docker-compose.staging.yml up

# Services behind Nginx reverse proxy
nginx → Python (8000) + Go (8001) + Frontend (3000)
```

### Production Environment
```
Internet
   ↓
CloudFlare (DNS + DDoS Protection)
   ↓
Application Load Balancer (AWS ALB)
   ├─ Go Orchestrator (x3) - Port 8001
   ├─ Python Backend (x3) - Port 8000
   └─ Next.js Frontend (x2) - Port 3000
   ↓
External Services
├─ RDS PostgreSQL (database)
├─ ElastiCache Redis (caching)
├─ S3 (file storage)
├─ CloudWatch (logging)
└─ Route 53 (DNS)
```

**Deployment Steps**
1. Build Docker images for each service
2. Push to Docker registry (ECR)
3. Deploy to Kubernetes or ECS
4. Configure load balancer
5. Setup auto-scaling policies
6. Configure monitoring and alerts

## Contributing

### Workflow
1. Create feature branch: `git checkout -b feature/your-feature`
2. Make changes following project structure
3. Write tests for new functionality
4. Run tests: `pytest tests/`
5. Commit with clear message: `git commit -m "feat(scope): description"`
6. Push and create pull request

### Code Standards

**Python Naming**
```python
# Functions: snake_case
def process_transfer(): pass

# Classes: PascalCase
class TransferRequest: pass

# Constants: UPPER_SNAKE_CASE
MAX_RETRIES = 3

# Private functions: leading underscore
def _internal_helper(): pass
```

**Code Style**
- Max 100 characters per line
- Use type hints: `def get_user(user_id: str) -> Optional[User]:`
- Use docstrings: `"""Description of what this does."""`
- Follow PEP 8 conventions

**Commit Message Format**
```
feat(accounts): add balance check endpoint

Implement GET /api/v1/accounts/{id}/balance to fetch
account balance information.

Closes #123
```

Types: `feat`, `fix`, `docs`, `refactor`, `test`, `chore`

### Testing

**Run all tests**
```bash
cd backend
pytest tests/
```

**Run specific test**
```bash
pytest tests/unit/test_api.py::TestTransfer::test_success
```

**With coverage**
```bash
pytest --cov=backend tests/
```

## Troubleshooting

### Backend Issues

**"Python not found" error**
```bash
# Verify Python 3.10+
python --version

# On Windows, try:
py --version

# On macOS/Linux with multiple versions:
python3 --version
```

**"ModuleNotFoundError" errors**
```bash
# Ensure virtual environment is activated
source venv/bin/activate  # macOS/Linux
venv\Scripts\activate     # Windows

# Reinstall dependencies
pip install -r backend/requirements.txt
```

**"Port 8000 already in use"**
```bash
# macOS/Linux: Kill process on port 8000
lsof -ti:8000 | xargs kill -9

# Windows: Find and kill process
netstat -ano | findstr :8000
taskkill /PID <PID> /F

# Or use different port
APP_PORT=8001 python main.py
```

**Backend starts but API returns errors**
```bash
# Check environment variables are set
cat .env | grep OPENAI
cat .env | grep RETELL

# Verify API keys are correct
# Test OpenAI connection:
curl -H "Authorization: Bearer $OPENAI_API_KEY" \
  https://api.openai.com/v1/models
```

### Frontend Issues

**"npm: command not found"**
```bash
# Reinstall Node.js
# Go to: https://nodejs.org

# Verify installation
npm --version
node --version
```

**"Port 3000 already in use"**
```bash
# macOS/Linux:
lsof -ti:3000 | xargs kill -9

# Windows:
netstat -ano | findstr :3000
taskkill /PID <PID> /F

# Or use different port:
PORT=3001 npm run dev
```

**"Module not found" error**
```bash
# Delete and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear Next.js cache
rm -rf .next
npm run dev
```

### Connection Issues

**"Cannot connect to backend"**
1. Verify backend is running: `curl http://localhost:8000/health`
2. Check backend port is 8000
3. Check firewall isn't blocking port 8000
4. Windows: Check Windows Defender Firewall settings

**"WebSocket connection failed"**
1. Verify Retell API key in `.env`
2. Check Retell account status at https://retellai.com/dashboard
3. Verify backend is running on correct port
4. Check browser console for detailed error message

**"CORS error from frontend"**
1. This shouldn't happen (CORS configured in backend)
2. Check backend is running on localhost:8000
3. Check frontend is on localhost:3000
4. Restart both services

## Project Status

### ✅ Completed Features
- ✓ FastAPI backend with multi-agent system
- ✓ REST API endpoints for banking operations
- ✓ WebSocket integration with Retell AI
- ✓ Next.js React admin dashboard
- ✓ OpenAI Swarm agent orchestration
- ✓ Pinata/IPFS storage integration
- ✓ Comprehensive documentation (7 docs files)
- ✓ Docker support
- ✓ Development environment setup

### 🚧 In Progress
- 🔄 Go orchestrator service (design complete)
- 🔄 Enhanced logging and monitoring
- 🔄 Production deployment configuration
- 🔄 Kubernetes manifests

### 📋 Planned for Future
- PostgreSQL database integration
- Redis caching layer
- Advanced agent routing (transfer between agents mid-conversation)
- Multi-language support
- Enhanced security (JWT, encryption, rate limiting)
- Comprehensive test suite (40+ tests)
- Load testing and benchmarks
- API versioning strategy
- GraphQL API alternative
- Mobile app frontend

## Resources & Documentation

### Quick Links
- **Main documentation**: See `/docs/` folder
- **System architecture**: `docs/SYSTEM.md`
- **API reference**: `docs/API.md`
- **Code structure**: `docs/STRUCTURE.md`
- **Development setup**: `docs/SETUP.md`
- **Code standards**: `docs/CONVENTIONS.md`
- **Go orchestrator**: `docs/GO_ORCHESTRATOR.md`
- **Go examples**: `docs/GO_EXAMPLES.md`

### External Resources
- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [OpenAI Swarm GitHub](https://github.com/openai/swarm)
- [Retell AI Documentation](https://docs.retellai.com)
- [Next.js Documentation](https://nextjs.org/docs)
- [Python Best Practices](https://pep8.org/)

### Getting Help
1. Check the documentation in `/docs/` folder
2. Review error messages carefully
3. Check troubleshooting section above
4. Search existing issues or discussions
5. Create a new issue with detailed description

## License

MIT License - See LICENSE file for details

---

**Last Updated**: November 2024
**Version**: 1.0.0
**Status**: Development

For additional technical details, see the documentation files in `/docs/` folder. For issues or questions, check the Troubleshooting section above or create an issue on the repository.
