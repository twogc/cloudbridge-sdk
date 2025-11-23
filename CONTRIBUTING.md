# Contributing to CloudBridge SDK

Thank you for your interest in contributing to CloudBridge SDK. This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and professional environment for all contributors.

## Getting Started

### Prerequisites

- Go 1.21 or higher (for Go SDK)
- Python 3.9 or higher (for Python SDK)
- Node.js 18 or higher (for JavaScript SDK)
- Git
- GitHub account

### Development Setup

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone git@github.com:YOUR_USERNAME/cloudbridge-sdk.git
   cd cloudbridge-sdk
   ```
3. Add upstream remote:
   ```bash
   git remote add upstream git@github.com:twogc/cloudbridge-sdk.git
   ```
4. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### Making Changes

1. Keep changes focused and atomic
2. Write clear commit messages
3. Add tests for new functionality
4. Update documentation as needed
5. Ensure all tests pass before submitting

### Commit Message Format

```
type(scope): brief description

Detailed explanation of changes (if needed)

Fixes #123
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Test additions or changes
- `refactor`: Code refactoring
- `chore`: Maintenance tasks

Example:
```
feat(go): add support for mesh networking

Implements mesh network join/leave operations with automatic
peer discovery and connection management.

Fixes #45
```

### Testing

#### Go SDK
```bash
cd go
go test ./...
go test -race ./...
go test -cover ./...
```

#### Python SDK
```bash
cd python
python -m pytest tests/
python -m pytest --cov=cloudbridge tests/
```

#### JavaScript SDK
```bash
cd javascript
npm test
npm run test:coverage
```

### Code Style

#### Go
- Follow standard Go formatting (`gofmt`, `goimports`)
- Use `golangci-lint` for linting
- Write idiomatic Go code
- Add comments for exported functions

#### Python
- Follow PEP 8 style guide
- Use type hints
- Format with `black`
- Lint with `pylint` and `mypy`

#### JavaScript
- Use ESLint configuration provided
- Format with Prettier
- Use TypeScript for type safety
- Follow Airbnb style guide

### Documentation

- Update README.md if adding user-facing features
- Add inline code documentation
- Update API reference documentation
- Include examples for new features

## Submitting Changes

### Pull Request Process

1. Update your branch with latest upstream:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

3. Create pull request on GitHub

4. Fill out the PR template completely

5. Wait for review and address feedback

### Pull Request Guidelines

- Keep PRs focused on a single feature or fix
- Include tests for new functionality
- Update documentation as needed
- Ensure CI passes
- Request review from maintainers

### PR Title Format

```
[SDK] type: brief description
```

Examples:
- `[Go SDK] feat: add mesh networking support`
- `[Python SDK] fix: handle connection timeout correctly`
- `[Docs] docs: update authentication guide`

## Reporting Issues

### Bug Reports

Include:
- SDK version
- Operating system and version
- Go/Python/Node.js version
- Steps to reproduce
- Expected behavior
- Actual behavior
- Error messages or logs

### Feature Requests

Include:
- Clear description of the feature
- Use cases and benefits
- Proposed API or interface
- Any implementation considerations

## Development Guidelines

### API Design Principles

1. **Simplicity**: Make common tasks simple
2. **Consistency**: Maintain consistent patterns across SDKs
3. **Type Safety**: Use strong typing where possible
4. **Error Handling**: Provide clear, actionable error messages
5. **Documentation**: Document all public APIs

### Security

- Never commit credentials or tokens
- Use secure defaults
- Validate all inputs
- Follow security best practices
- Report security issues privately to security@2gc.ru

### Performance

- Profile code for performance bottlenecks
- Avoid unnecessary allocations
- Use connection pooling
- Implement efficient retry logic
- Add benchmarks for critical paths

## Release Process

Maintainers will handle releases. The process includes:

1. Version bump in all SDKs
2. Update CHANGELOG.md
3. Create GitHub release
4. Publish packages to registries
5. Update documentation

## Questions and Support

- GitHub Discussions: For questions and discussions
- GitHub Issues: For bug reports and feature requests
- Email: dev@2gc.ru for private inquiries

## License

By contributing to CloudBridge SDK, you agree that your contributions will be licensed under the MIT License.

## Recognition

Contributors will be acknowledged in:
- CONTRIBUTORS.md file
- Release notes
- Project documentation

Thank you for contributing to CloudBridge SDK!
