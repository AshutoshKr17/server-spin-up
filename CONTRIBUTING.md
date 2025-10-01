# Contributing to server-spin-up

Thank you for your interest in contributing to server-spin-up! We welcome contributions from everyone. This document will guide you through the contribution process.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before contributing.

## How Can I Contribute?

### Reporting Bugs

If you find a bug, please create an issue using the Bug Report template. Include:
- A clear and descriptive title
- Steps to reproduce the issue
- Expected behavior
- Actual behavior
- System information (OS, version, etc.)
- Any relevant logs or screenshots

### Suggesting Enhancements

We love new ideas! Create an issue using the Feature Request template and include:
- A clear and descriptive title
- Detailed description of the proposed feature
- Why this enhancement would be useful
- Any examples or mockups if applicable

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:
- `good first issue` - Simple issues perfect for newcomers
- `help wanted` - Issues where we need community help
- `beginner-friendly` - Great for first-time contributors

### Pull Request Process

1. **Fork the repository** and create your branch from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```
   or
   ```bash
   git checkout -b fix/your-bug-fix
   ```

2. **Make your changes**:
   - Write clear, commented code
   - Follow the existing code style
   - Update documentation if needed

3. **Test your changes**:
   - Ensure all existing tests pass
   - Add new tests if applicable
   - Test manually in different scenarios

4. **Commit your changes**:
   - Use clear and meaningful commit messages
   - Follow conventional commit format when possible:
     - `feat:` for new features
     - `fix:` for bug fixes
     - `docs:` for documentation changes
     - `style:` for formatting changes
     - `refactor:` for code refactoring
     - `test:` for adding tests
     - `chore:` for maintenance tasks

5. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```

6. **Open a Pull Request**:
   - Use the PR template
   - Link related issues
   - Provide a clear description of changes
   - Add screenshots/demos if applicable

## Development Setup

```bash
# Clone your fork
git clone https://github.com/your-username/server-spin-up.git
cd server-spin-up

# Add upstream remote
git remote add upstream https://github.com/original-owner/server-spin-up.git

# Install dependencies (if applicable)
# Add specific setup instructions here
```

## Style Guidelines

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature")
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests after the first line

### Code Style

- Write clean, readable, and maintainable code
- Comment your code when necessary
- Follow existing patterns in the codebase
- Keep functions small and focused

### Documentation

- Update README.md if you change functionality
- Comment complex logic
- Add JSDoc/docstrings to functions
- Update relevant documentation files

## Review Process

- Maintainers will review your PR
- Address any requested changes
- Once approved, a maintainer will merge your PR
- Your contribution will be credited in the project

## Getting Help

- Check existing issues and documentation
- Ask questions in issue comments
- Join our community discussions (if available)
- Tag maintainers if you need help

## Recognition

Contributors will be:
- Listed in our README
- Credited in release notes
- Part of our growing community!

## Hacktoberfest

This project participates in Hacktoberfest! To have your PR count:
- Make sure your PR is meaningful and not spammy
- Follow our contribution guidelines
- PR must be merged or labeled `hacktoberfest-accepted`

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to server-spin-up! 🎉

