# AI Banking Assistant

> ### Multi-Agent Conversational AI Banking System with Voice Integration

<p align="center">
  Frontend built with: <br>
  <img src=https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB alt="React">
  <img src=https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white alt="Typescript">
  <img src=https://img.shields.io/badge/next%20js-000000?style=for-the-badge&logo=nextdotjs&logoColor=white alt="Next">
  <img src=https://img.shields.io/badge/Tailwind_CSS-38B2AC?style=for-the-badge&logo=tailwind-css&logoColor=white alt="Tailwind">
  <br><br>
  Backend built with: <br>
  <img src=https://img.shields.io/badge/Python-3776AB?style=for-the-badge&logo=python&logoColor=white alt="Python">
  <img src=https://img.shields.io/badge/FastAPI-009688?style=for-the-badge&logo=fastapi&logoColor=white alt="FastAPI">
  <img src="https://img.shields.io/badge/OpenAI_Swarm-34a853?style=for-the-badge&logo=openai&logoColor=white" alt="OpenAI Swarm">
  <img src="https://img.shields.io/badge/Pinata_Cloud-F2E3A1?style=for-the-badge&logo=pinata&logoColor=black" alt="Pinata Cloud">
  <img src="https://img.shields.io/badge/Retell_AI-FF7F50?style=for-the-badge&logo=retell&logoColor=white" alt="Retell AI">
  <br>
</p>

## Overview

An innovative voice-driven banking assistant that demonstrates how conversational AI and multi-agent systems can democratize financial services. This proof-of-concept (POC) enables users to perform banking tasks through simple phone calls using natural language.

## What It Does

**AI Banking Assistant** provides a fully conversational banking experience accessible through voice. Users can manage their finances by speaking naturally to perform tasks such as:
- Checking account balances
- Transferring funds between accounts
- Scheduling payments
- Viewing transaction history
- Applying for loans or credit cards

Key features include:
- **Voice-Activated Services**: Conduct banking transactions hands-free via phone
- **Natural Language Processing**: Understand complex user intents and requests
- **Multi-Agent Orchestration**: Specialized agents handle different banking domains
- **Secure Operations**: Decentralized data storage with Pinata/IPFS
- **Admin Dashboard**: Comprehensive interface for monitoring transactions and managing data

## Architecture

### Technology Stack

**Frontend:**
- Next.js with TypeScript
- Tailwind CSS for styling
- Real-time WebSocket communication

**Backend:**
- FastAPI (Python) for main orchestration
- OpenAI Swarm for multi-agent orchestration
- Retell AI for voice processing
- Pinata Cloud for IPFS storage

**Planned:**
- Go orchestrator for replacing Python services

### Project Structure

```
bankAssistant/
├── client/                 # Next.js frontend
│   └── src/app/           # React components & pages
├── server/                 # FastAPI backend
│   ├── core/              # Type definitions & constants
│   ├── database/          # Data layer
│   ├── services/          # Business logic (LLM, Agents, Pinata)
│   ├── routers/           # API route handlers
│   ├── websocket/         # WebSocket connection management
│   ├── agents/            # AI agent implementations
│   └── utils/             # Helper functions
└── experiments/           # Development & testing
```

## How It Works

1. **Voice Input**: User calls the system and speaks naturally
2. **Speech-to-Text**: Retell AI converts voice to text
3. **Intent Understanding**: LLM identifies user intent
4. **Agent Routing**: Triage agent routes to appropriate specialized agent (Accounts, Payments, Applications)
5. **Execution**: Agent performs banking operation on mock database
6. **Response Generation**: LLM generates natural language response
7. **Text-to-Speech**: Retell AI converts response back to voice
8. **Audit Trail**: Receipts stored on IPFS via Pinata

## API Endpoints

### Health & Webhooks
```
GET  /health                           # Health check
GET  /status                           # Application status
POST /webhook                          # Retell AI webhooks
```

### WebSocket
```
WS   /llm-websocket/{call_id}         # Voice interaction with Retell
WS   /admin/ws?client_id=...          # Admin dashboard
```

### Admin Dashboard
```
GET  /admin/database                   # Get full database state
GET  /admin/calls                      # Get all call records
GET  /admin/status                     # Get admin service status
```

### API v1 (REST)
```
GET  /api/v1/users                     # List all users
GET  /api/v1/users/{user_id}           # Get user details
GET  /api/v1/users/{user_id}/accounts  # Get user's accounts

GET  /api/v1/accounts                  # List all accounts
GET  /api/v1/accounts/{id}             # Get account details
GET  /api/v1/accounts/{id}/balance     # Get account balance

GET  /api/v1/payments                  # List all payments
POST /api/v1/payments/transfer         # Transfer funds

POST /api/v1/applications/loan         # Apply for loan
POST /api/v1/applications/credit-card  # Apply for credit card
```

## Setup & Installation

### Prerequisites
- Python 3.9+
- Node.js 18+
- Retell AI API key
- OpenAI API key
- Pinata API key & secret

### Backend Setup

```bash
cd server
python -m venv venv
source venv/bin/activate  # Windows: venv\Scripts\activate
pip install -r requirements.txt

# Create .env file
echo "RETELL_API_KEY=your_key_here" >> .env
echo "OPENAI_API_KEY=your_key_here" >> .env
echo "PINATA_API_KEY=your_key_here" >> .env
echo "PINATA_API_SECRET=your_secret_here" >> .env

# Run server
python -m uvicorn main:app --reload
```

### Frontend Setup

```bash
cd client
npm install
npm run dev
```

Visit `http://localhost:3000` for the admin dashboard.

## Design & Code Organization

The codebase follows a modular, scalable architecture:

- **core/**: Type definitions and application constants
- **database/**: Data layer with mock data and CRUD operations
- **services/**: Business logic (LLM interaction, agent orchestration, IPFS)
- **routers/**: FastAPI route handlers organized by domain
- **agents/**: AI agent implementations using OpenAI Swarm
- **websocket/**: WebSocket connection management

This structure makes it easy to:
- Add new features without cluttering existing code
- Test components independently
- Prepare for Go orchestrator migration
- Scale to production systems

## Future Roadmap

### Phase 1: Foundation (Complete)
- Core multi-agent system
- Mock banking database
- FastAPI refactoring for modularity

### Phase 2: Enhancement
- Database persistence with SQLAlchemy
- Authentication & authorization
- Comprehensive logging
- Unit & integration tests

### Phase 3: Scalability
- Go orchestrator implementation
- Load balancing setup
- Caching layer (Redis)
- API rate limiting

### Phase 4: Production
- Real banking API integration
- Enhanced security measures
- Compliance (PCI-DSS, GDPR)
- Multi-language support

## Key Learnings

- **Voice Technology**: Complex challenges in creating natural, intuitive voice interactions
- **Multi-Agent Systems**: Effective agent design requires clear responsibilities and communication patterns
- **Security in Finance**: Importance of multi-layered security approaches
- **API Design**: Clean, versioned APIs prepare systems for future architectural changes

## Contributing

This is a proof-of-concept project. Contributions are welcome for:
- Enhanced agent prompts and behaviors
- Additional banking use cases
- Security improvements
- Frontend enhancements

## License

This project is provided as-is for educational and development purposes.

## Contact & Support

For questions or issues, please check the documentation or reach out through the project repository.
