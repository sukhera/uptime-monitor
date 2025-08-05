# Mockery & GolangCI-Lint Implementation Guide

This document explains how **Mockery** (for Go mocks) and **GolangCI-Lint** (for linting) are implemented in this project. It covers configuration, usage, best practices, and troubleshooting for both tools.

---

## 1. Mockery Implementation

### What is Mockery?
[Mockery](https://github.com/vektra/mockery) is a Go code generator for creating mocks of interfaces, used for unit testing with libraries like `stretchr/testify`.

### How Mockery is Set Up

#### a. Configuration: `.mockery.yaml`
- All mock generation is controlled by `.mockery.yaml` at the project root.
- Example config:

```yaml
with-expecter: true
packages:
  github.com/sukhera/uptime-monitor/internal/infrastructure/database:
    config:
      inpackage: true
      mockname: "MockDatabaseInterface"
      filename: "MockDatabaseInterface.go"
    interfaces:
      Interface:
  # ... more packages/interfaces ...
```
- **Key options:**
  - `inpackage: true` — Mocks are generated in the same package as the interface.
  - `with-expecter: true` — Enables the modern expecter pattern for more readable test expectations.

#### b. Makefile Integration
- The Makefile provides easy commands:

```makefile
# Generate mocks
make generate-mocks

# Run tests with mocks (regenerates mocks first)
make test-with-mocks
```
- Example Makefile target:
```makefile
generate-mocks:
	@echo "$(BLUE)Generating mocks...$(NC)"
	@go install github.com/vektra/mockery/v2@latest
	@mockery
	@echo "$(GREEN)✓ Mocks generated!$(NC)"
```

#### c. Usage in Tests
- Import the generated mocks from the same package as the interface.
- Use the expecter pattern for clear, table-driven tests.
- Example:

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/sukhera/uptime-monitor/internal/infrastructure/database"
)

func TestSomething(t *testing.T) {
    mockDB := &database.MockDatabaseInterface{}
    mockDB.EXPECT().SomeMethod().Return(nil)
    // ...
}
```

#### d. Best Practices
- **Table-driven tests**: Use slices of test cases for coverage.
- **Expecter pattern**: Use `.EXPECT()` for clear, readable expectations.
- **Mocks in-package**: Keeps test dependencies simple and discoverable.
- **Regenerate mocks on interface change**: Always run `make generate-mocks` after changing interfaces.

#### e. Troubleshooting
- **Mocks not found**: Run `make generate-mocks`.
- **Mocks out of date**: Regenerate after interface changes.
- **Import errors**: Ensure you import the mock from the correct package (same as the interface).
- **Mockery config errors**: Check `.mockery.yaml` for typos or missing options.

---

## 2. GolangCI-Lint Implementation

### What is GolangCI-Lint?
[GolangCI-Lint](https://golangci-lint.run/) is a fast Go linters aggregator, running multiple linters in parallel.

### How GolangCI-Lint is Set Up

#### a. Configuration: `.golangci.yml`
- Located at the project root.
- Example config (v2):

```yaml
version: "2"
run:
  timeout: 5m
  modules-download-mode: readonly
linters:
  enable:
    - errcheck
    - govet
    - ineffassign
    - staticcheck
    - unused
    - gosec
    - bodyclose
    - copyloopvar
    - errname
    - makezero
    - misspell
    - rowserrcheck
    - sqlclosecheck
    - goconst
    - dupl
```
- **Note:** Only include linters supported by your installed version (see `golangci-lint help linters`).

#### b. Makefile Integration
- The Makefile provides a lint target:

```makefile
make lint
```
- Example target:
```makefile
lint:
	@echo "$(YELLOW)Running linters...$(NC)"
	@golangci-lint run --config .golangci.yml ./...
	@echo "$(GREEN)✓ Linting completed!$(NC)"
```

#### c. Usage
- Run `make lint` to check the whole project.
- Run `golangci-lint run --config .golangci.yml ./path/to/dir` to check a specific directory.

#### d. Best Practices
- **Run linting before every commit/PR**.
- **Fix all errors and warnings** before merging.
- **Keep `.golangci.yml` up to date** with only the linters you use.
- **Upgrade golangci-lint regularly** to get new linters and bug fixes.

#### e. Troubleshooting
- **Unknown linter error**: Remove or replace the linter in `.golangci.yml`.
- **Config version error**: Set `version: "2"` for golangci-lint v2+.
- **Formatter errors**: Remove `gofmt` and `goimports` from the `enable` list (they are not linters in v2).
- **IDE linter errors**: Some IDEs may expect a `version` field; always match your config to your installed golangci-lint version.

---

## 3. Practical Tips for New Contributors
- **Always regenerate mocks after changing interfaces.**
- **Run `make lint` and `make test-with-mocks` before pushing.**
- **If you see a linter or mock error, check the config and Makefile first.**
- **Keep your tools up to date:**
  - `go install github.com/vektra/mockery/v2@latest`
  - `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
- **Read the docs:**
  - [Mockery](https://github.com/vektra/mockery)
  - [GolangCI-Lint](https://golangci-lint.run/)

---

**This guide should help you quickly get up to speed with testing and linting in this project!**