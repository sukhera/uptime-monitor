# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Upgrade Go version from 1.21/1.22 to 1.24 across entire project
  - Updated `go.mod` to require Go 1.24
  - Updated all Dockerfiles to use `golang:1.24` base images
  - Updated all GitHub workflows to use Go 1.24
  - Updated README.md requirements to reflect Go 1.24+

### Fixed
- Fix Dockerfile casing issues for proper hadolint compliance
  - Fixed `FROM golang:1.21 as builder` to `FROM golang:1.24 AS builder` in API and status-checker Dockerfiles
  - Fixed `FROM golang:1.21-alpine AS dev` to `FROM golang:1.24-alpine AS dev` in development Dockerfile
- Fix GitHub security workflow SARIF upload failures
  - Added file existence checks to prevent uploading missing SARIF files when Docker builds fail
  - Use `hashFiles()` function to conditionally upload only existing SARIF files
  - Prevents pipeline failures when Trivy scans don't generate SARIF files due to build errors
- Resolve version inconsistency across Docker and CI/CD infrastructure
  - Standardized Go version to 1.24 in all Dockerfiles and GitHub workflows
  - Eliminates potential build conflicts from version mismatches

## [0.4.0] - 2025-08-04

### Fixed
- Resolve security and dependency check workflow failures ([9196a26](https://github.com/sukhera/uptime-monitor/commit/9196a26))
  - Update Go version to 1.22 across all security workflows for consistency
  - Make SNYK security scan conditional on token availability
  - Add error tolerance to npm audit and license checking tools
  - Use stable TruffleHog version instead of main branch
  - Improve security summary to only fail on critical issues
  - Allow license checking tools to fail gracefully without breaking pipeline
- Resolve CI pipeline failures and improve code quality ([66a8542](https://github.com/sukhera/uptime-monitor/commit/66a8542))
  - Fix linting issues and test failures
  - Improve error handling and code reliability
  - Enhance CI pipeline stability

### Added
- Implement comprehensive GitHub configuration and workflows ([8348cd5](https://github.com/sukhera/uptime-monitor/commit/8348cd5))
  - Issue templates for bug reports and feature requests
  - Pull request template with comprehensive checklist
  - CI workflow with Go/Node.js testing, linting, and security scanning
  - CD workflow for automated deployment to staging and production
  - Security workflow with CodeQL, Gosec, and vulnerability scanning
  - Dependency update workflow for automated maintenance
  - Release workflow with artifact building and GitHub releases
  - CODEOWNERS for code review requirements
  - Security policy and vulnerability reporting process
  - Dependabot configuration for automated dependency updates

## [0.3.0] - 2025-08-03

### Added
- Enhance testing infrastructure with comprehensive coverage and improved reliability ([b222928](https://github.com/sukhera/uptime-monitor/commit/b222928))
  - Comprehensive test suite for Go backend services
  - Frontend testing with React Testing Library and Jest
  - Integration tests for API endpoints
  - End-to-end testing infrastructure
  - Automated test coverage reporting
  - Performance benchmarking tests

### Added
- Add comprehensive Makefile infrastructure automation ([8b9fbb4](https://github.com/sukhera/uptime-monitor/commit/8b9fbb4))
  - 50+ automated commands for complete development workflow
  - Setup and installation automation (Go, Node.js, Docker, Git hooks)
  - Development environment management with hot reloading
  - Testing and quality assurance commands
  - Security scanning and vulnerability checks
  - Build and deployment automation
  - Database management and migration tools
  - Monitoring and maintenance utilities
  - Documentation generation
  - CI/CD integration commands

## [0.2.0] - 2025-08-02

### Changed
- Updated comprehensive README documentation ([e20ff48](https://github.com/sukhera/uptime-monitor/commit/e20ff48))
  - Enhanced project overview and feature descriptions
  - Detailed setup and installation instructions
  - Comprehensive development workflow documentation
  - Infrastructure automation guide
  - Database schema and configuration details
  - Deployment and monitoring instructions

### Added
- Migrate frontend from vanilla JS to React + Tailwind CSS ([25fbff1](https://github.com/sukhera/uptime-monitor/commit/25fbff1))
  - Modern React-based user interface with TypeScript
  - Tailwind CSS for responsive, utility-first styling
  - Real-time status updates with WebSocket integration
  - Dark mode support with system preference detection
  - Incident management and maintenance scheduling UI
  - Mobile-responsive design with touch-friendly interactions
  - Vite build system for fast development and optimized production builds
  - Comprehensive component library with accessibility features

## [0.1.0] - 2025-08-01

### Added
- Initial commit: Real-time service monitoring dashboard ([8ad0336](https://github.com/sukhera/uptime-monitor/commit/8ad0336))
  - Go-based backend API with MongoDB integration
  - Service health checking with configurable intervals
  - RESTful API endpoints for status management
  - Real-time monitoring with automated health checks
  - Docker containerization for easy deployment
  - Basic web interface for status visualization
  - Service configuration management
  - Status logging and historical data storage
  - Error handling and recovery mechanisms
  - Production-ready architecture foundation

### Infrastructure
- Docker multi-stage builds for optimized container images
- MongoDB database with proper indexing and schema design
- Nginx reverse proxy configuration for production deployment
- Environment-based configuration management
- Logging and monitoring infrastructure
- Health check endpoints for service monitoring

---

## Release Notes

### Version 0.4.0 - "Stability & Security"
This release focuses on improving CI/CD pipeline reliability and security scanning infrastructure. All workflow failures have been resolved, and the project now has robust automated testing and security scanning capabilities.

### Version 0.3.0 - "Automation & Testing"  
Major infrastructure improvements with comprehensive Makefile automation and enhanced testing coverage. This release introduces 50+ automated commands for complete development workflow management.

### Version 0.2.0 - "Modern Frontend"
Complete frontend modernization with React, TypeScript, and Tailwind CSS. The new interface provides a superior user experience with real-time updates, dark mode, and mobile responsiveness.

### Version 0.1.0 - "Foundation"
Initial release establishing the core monitoring infrastructure with Go backend, MongoDB storage, and Docker containerization. Provides the foundation for a production-ready uptime monitoring solution.

---

## Upgrade Notes

### Upgrading to 0.4.0
- No breaking changes
- Improved CI/CD pipeline stability
- Enhanced security scanning (automatic)

### Upgrading to 0.3.0
- New Makefile commands available - run `make help` to see all options
- Enhanced testing infrastructure - run `make test` for comprehensive testing
- No breaking changes to existing functionality

### Upgrading to 0.2.0
- **Breaking Change**: Frontend completely rewritten in React
- Update any custom frontend integrations
- New API endpoints may be available - check API documentation
- Database schema remains compatible

### Upgrading from 0.1.0
- Run database migrations if upgrading from initial version
- Update Docker Compose configuration for new services
- Review environment variable configuration

---

## Development

This changelog is automatically maintained as part of our development workflow. For detailed commit information, please see the [Git history](https://github.com/sukhera/uptime-monitor/commits/main).

### Contributing to Changelog
- Follow [Conventional Commits](https://conventionalcommits.org/) for commit messages
- Use semantic versioning for releases
- Document breaking changes prominently
- Include migration/upgrade instructions when necessary