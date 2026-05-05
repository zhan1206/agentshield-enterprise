# Contributing to AgentShield Enterprise

Thank you for your interest in contributing! 🛡️

## Getting Started

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/amazing-feature`
3. Commit your changes: `git commit -m 'Add amazing feature'`
4. Push to the branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

## Development Setup

```bash
# Rust sandbox engine
cd sandbox && cargo build

# Go backend
cd backend && go mod tidy && go run ./cmd/server

# Python AI engine
cd ai && pip install -r requirements.txt

# Frontend
cd frontend && npm install && npm run dev
```

## Code Standards

- **Rust**: Follow `rustfmt` + `clippy`
- **Go**: Follow `gofmt` + `golint`
- **Python**: Follow PEP 8, use type hints
- **TypeScript**: Strict mode, ESLint + Prettier

## Security Contributions

We especially welcome:
- New threat detection rules and patterns
- Additional Agent framework adapters
- Security audit improvements
- Compliance rule additions

## Reporting Security Vulnerabilities

Please do NOT open public issues for security vulnerabilities. Email: security@agentshield.dev

## License

By contributing, you agree your code will be licensed under Apache 2.0.
