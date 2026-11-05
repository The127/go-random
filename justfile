# List the recipes.
default:
    @just --list

# Run the tests with the race detector.
test:
    go test -race ./...

# Lint the code.
lint:
    golangci-lint run ./...

# Format the code in place.
fmt:
    golangci-lint fmt ./...

# Check the dependencies for known vulnerabilities.
vuln:
    govulncheck ./...

# Check that every file has license and copyright information.
reuse:
    reuse lint

# Run every check that CI runs.
ci: lint test vuln reuse
