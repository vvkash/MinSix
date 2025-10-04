# Contributing to Minsix

Thank you for your interest in contributing to Minsix! This document provides guidelines and instructions for contributing.

## Development Setup

1. **Fork and clone the repository**
   ```bash
   git clone https://github.com/yourusername/minsix.git
   cd minsix
   ```

2. **Run the setup script**
   ```bash
   make setup
   ```

3. **Get an Alchemy API key**
   - Sign up at [Alchemy](https://www.alchemy.com/)
   - Create a new app for Ethereum Mainnet
   - Copy your API key to `backend/.env`

4. **Start development environment**
   ```bash
   make dev
   ```

## Project Structure

```
minsix/
├── backend/              # Go backend service
│   ├── cmd/             # Application entrypoints
│   ├── internal/        # Internal packages
│   │   ├── database/    # Database layer
│   │   ├── detector/    # Fraud detection engine
│   │   ├── ethereum/    # Ethereum client
│   │   ├── handlers/    # HTTP handlers
│   │   ├── models/      # Data models
│   │   └── websocket/   # WebSocket server
│   └── migrations/      # Database migrations
├── frontend/            # Next.js frontend
│   └── src/
│       ├── app/         # Next.js app directory
│       ├── components/  # React components
│       └── types/       # TypeScript types
└── scripts/             # Utility scripts
```

## Code Style

### Go (Backend)
- Follow standard Go conventions
- Run `go fmt` before committing
- Use meaningful variable names
- Add comments for exported functions

### TypeScript (Frontend)
- Use TypeScript strict mode
- Follow React best practices
- Use functional components with hooks
- Run `npm run lint` before committing

## Making Changes

### 1. Create a Branch
```bash
git checkout -b feature/your-feature-name
```

### 2. Make Your Changes
- Write clear, concise commit messages
- Keep commits atomic and focused
- Add tests for new features
- Update documentation as needed

### 3. Test Your Changes

**Backend:**
```bash
cd backend
go test ./...
```

**Frontend:**
```bash
cd frontend
npm run build
```

### 4. Submit a Pull Request
- Provide a clear description of changes
- Reference any related issues
- Ensure CI passes
- Request review from maintainers

## Adding New Fraud Detection Heuristics

To add a new detection heuristic:

1. **Add detection logic in `backend/internal/detector/detector.go`:**
   ```go
   func (fd *FraudDetector) checkYourHeuristic(tx *models.Transaction) bool {
       // Your detection logic
       return false
   }
   ```

2. **Call it in `AnalyzeTransaction`:**
   ```go
   if suspicious := fd.checkYourHeuristic(tx); suspicious {
       reasons = append(reasons, "Your reason here")
       riskScore += 15
   }
   ```

3. **Add tests:**
   ```go
   func TestCheckYourHeuristic(t *testing.T) {
       // Test cases
   }
   ```

4. **Update documentation** in README.md

## Database Migrations

To add a new migration:

1. Create a new migration file:
   ```bash
   touch backend/migrations/002_your_migration.sql
   ```

2. Add SQL statements:
   ```sql
   -- Add your migration SQL
   ALTER TABLE transactions ADD COLUMN new_field TEXT;
   ```

3. Update the migration runner if needed

## Frontend Components

When adding new components:

1. Create component file in `frontend/src/components/`
2. Use TypeScript for type safety
3. Follow existing component patterns
4. Make it responsive (mobile-first)
5. Use Tailwind CSS for styling

## Performance Considerations

- **Database queries**: Use indexes, avoid N+1 queries
- **WebSocket**: Handle reconnection gracefully
- **Frontend**: Use React.memo for expensive renders
- **API**: Implement pagination for large datasets

## Security

- Never commit API keys or secrets
- Sanitize user inputs
- Use parameterized queries
- Keep dependencies updated
- Report security issues privately

## Testing

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Aim for >70% code coverage

### Integration Tests
- Test end-to-end workflows
- Use test database
- Clean up test data

### Manual Testing
- Test on different networks (mainnet, testnet)
- Verify WebSocket reconnection
- Check mobile responsiveness

## Documentation

- Update README.md for user-facing changes
- Add inline comments for complex logic
- Update API documentation
- Include examples where helpful

## Getting Help

- **Questions**: Open a GitHub Discussion
- **Bugs**: Open a GitHub Issue
- **Security**: Email security@minsix.io
- **Chat**: Join our Discord (coming soon)

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Help others learn and grow
- Focus on what's best for the project

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes for significant contributions
- Annual contributor highlights

Thank you for contributing to Minsix!
