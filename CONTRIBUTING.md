# Contributing to Matt

Thank you for your interest in contributing to Matt! We welcome contributions of all kinds, including bug fixes, new features, documentation improvements, and feedback.

## AI Contributions Welcome

We welcome contributions generated or assisted by AI tools. However, to ensure high quality and make evaluation easy for maintainers, please adhere to the following rules:

1. **Use Up-to-Date Code**: Ensure that your AI assistant is working with the latest code from the `main` branch. Avoid submitting legacy formats or outdated designs.
2. **Test Before Submitting**: Check that your changes compile and run correctly on your local machine. Never submit code blindly without testing it yourself.
3. **Run Verification Checks**: Make sure to run all local checks before opening a pull request:
   ```bash
   # Formats and cleans Go imports
   go fmt ./...

   # Analyzes code for potential issues
   go vet ./...

   # Runs the test suite
   go test -v ./...

   # Ensures the project builds successfully
   go build ./...
   ```
4. **Document the Changes**: Clearly describe what the AI-generated code changes, why it does it, and confirm that all checks passed successfully.

By verifying your code first, maintainers can easily evaluate and merge your pull request!

## How to Contribute

1. Fork the repository and create your branch from `main`.
2. Make your changes and write tests if applicable.
3. Verify that all linting, vetting, testing, and building checks pass cleanly.
4. Push your branch and submit a Pull Request.

## Code of Conduct

Please be respectful and constructive in all communication channels.
