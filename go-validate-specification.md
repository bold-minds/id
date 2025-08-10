# Review: Universal Go Project Analysis

## 🎯 Project Overview

**Review** is a comprehensive project analysis tool for Go projects that provides both local CLI usage and GitHub Action integration. It standardizes Go project review across repositories with consistent linting, testing, coverage analysis, and badge generation.

### Key Problems Solved
- Inconsistent validation standards across Go projects
- Complex setup for comprehensive Go project analysis
- Manual maintenance of project health badges
- Time-consuming local review setup
- Fragmented tooling across different repositories

### Solution Approach
- **Dual Distribution**: Local CLI tool + GitHub Action
- **Single Source of Truth**: One repository, multiple usage patterns
- **Comprehensive Analysis**: Linting, testing, coverage, documentation, security
- **Beautiful Output**: Colored terminal output with emojis and progress tracking
- **Badge Generation**: Automatic health badges for README files
- **Configurable Profiles**: Minimal, standard, and strict review levels

---

## 🚀 Repository Structure

```
go-validate/
├── action.yml                      # GitHub Action definition
├── Dockerfile                      # Container for GitHub Action
├── entrypoint.sh                   # GitHub Action entry point
├── review.sh                       # Core review script (bash)
├── install.sh                      # Local installation script
├── profiles/
│   ├── minimal.yml                 # Basic validation profile
│   ├── standard.yml                # Standard profile (recommended)
│   └── strict.yml                  # High-standard profile
├── templates/
│   ├── .golangci.yml               # Standard golangci-lint config
│   ├── .review.yml                 # Configuration template
│   └── github-workflow.yml         # Example GitHub workflow
├── examples/
│   ├── local-usage/
│   │   ├── basic/
│   │   ├── custom-config/
│   │   └── advanced/
│   ├── github-action/
│   │   ├── simple/
│   │   ├── matrix-builds/
│   │   └── badge-publishing/
│   └── integration/
│       ├── makefile-integration/
│       └── pre-commit-hooks/
├── docs/
│   ├── configuration.md
│   ├── profiles.md
│   ├── badge-system.md
│   ├── troubleshooting.md
│   └── contributing.md
├── tests/
│   ├── test-projects/              # Sample Go projects for testing
│   └── integration-tests.sh        # Test the validator itself
├── README.md
├── CHANGELOG.md
└── LICENSE
```

---

## 🎛️ Usage Patterns

### 1. Local CLI Installation & Usage

#### Installation Options:
```bash
# Option 1: Quick install script
curl -sSL https://raw.githubusercontent.com/yourusername/review/main/install.sh | bash

# Option 2: Manual download
wget https://raw.githubusercontent.com/yourusername/review/main/review.sh
chmod +x review.sh
sudo mv review.sh /usr/local/bin/review

# Option 3: Go install (if we add a Go wrapper)
go install github.com/yourusername/review/cmd/review@latest
```

#### Local Usage:
```bash
# Basic review
review

# Different modes
review local                        # Local development mode (default)
review ci                           # CI mode (stricter, no interactive prompts)

# Configuration options
COVERAGE_THRESHOLD=90 review        # Custom coverage threshold
SKIP_INTEGRATION=true review        # Skip integration tests
TEST_TIMEOUT=15m review             # Custom test timeout

# Review profiles
review --profile minimal            # Basic checks only
review --profile standard           # Recommended settings (default)
review --profile strict             # High standards (95% coverage, all linters)

# Custom configuration
review --config .my-config.yml      # Use custom config file

# Initialize in new project
review init                         # Create .review.yml and .golangci.yml
review init --with-github           # Also create GitHub Actions workflow
```

### 2. GitHub Action Integration

#### Basic Usage:
```yaml
# .github/workflows/validate.yml
name: Go Validation

on: [push, pull_request]

jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Validate Go Project
        uses: yourusername/go-validate-action@v1
        with:
          coverage-threshold: 85
          go-version: '1.24'
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Advanced Usage with Matrix:
```yaml
name: Go Review Matrix

on: [push, pull_request]

jobs:
  validate:
    strategy:
      matrix:
        go-version: ['1.22', '1.23', '1.24']
        os: [ubuntu-latest, macos-latest, windows-latest]
        profile: [standard, strict]
    
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/checkout@v4
      
      - name: Review Go Project
        uses: yourusername/review-action@v1
        with:
          go-version: ${{ matrix.go-version }}
          review-profile: ${{ matrix.profile }}
          badge-generation: true
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

#### Badge Publishing Workflow:
```yaml
- name: Review and Generate Badges
  uses: yourusername/review-action@v1
  with:
    badge-generation: true
    
- name: Publish Badges
  if: github.ref == 'refs/heads/main'
  run: |
    git config --local user.email "action@github.com"
    git config --local user.name "GitHub Action"
    git add .github/badges/
    git commit -m "Update project badges [skip ci]" || exit 0
    git push
```

---

## ⚙️ Configuration System

### Global Configuration: `~/.review.yml`
```yaml
# Global defaults for all projects
default_profile: "standard"
date_format: "2006-01-02"
badge_output_dir: ".github/badges"
github_token_env: "GITHUB_TOKEN"

# Default review settings
review:
  coverage_threshold: 80
  test_timeout: "10m"
  integration_tag: "integration"
  skip_integration: false

# Badge configuration
badges:
  enabled: true
  formats: ["shields.io"]
  colors:
    success: "brightgreen"
    warning: "yellow"
    error: "red"
    info: "blue"
```

### Per-Project Configuration: `.review.yml`
```yaml
# Override global settings for this project
profile: "strict"

validation:
  coverage_threshold: 95
  test_timeout: "20m"
  skip_integration: false
  custom_linters: ["gosec", "ineffassign", "misspell"]

# Project-specific badge configuration
badges:
  enabled: true
  output_dir: ".github/badges"
  include:
    - "golangci-lint"
    - "coverage"
    - "go-version"
    - "last-updated"
    - "security"

# Custom review steps
custom_steps:
  - name: "API Documentation"
    command: "swag init"
    required: false
  - name: "Database Migrations"
    command: "migrate -path ./migrations -database $DATABASE_URL up"
    required: false
    skip_in_ci: true

# Manual content for generated documentation
manual_sections:
  optimization_notes: |
    Recent performance improvements include...
  future_improvements:
    - "Implement connection pooling"
    - "Add Redis caching layer"
```

### Validation Profiles

#### Minimal Profile (`profiles/minimal.yml`):
```yaml
name: "minimal"
description: "Basic review for rapid development"

review:
  coverage_threshold: 50
  test_timeout: "5m"
  skip_integration: true
  
linting:
  enabled: true
  config: "minimal"
  
badges:
  enabled: false
```

#### Standard Profile (`profiles/standard.yml`):
```yaml
name: "standard"
description: "Recommended review for most projects"

review:
  coverage_threshold: 80
  test_timeout: "10m"
  skip_integration: false
  
linting:
  enabled: true
  config: "standard"
  security_checks: true
  
badges:
  enabled: true
  include: ["lint", "coverage", "go-version", "security"]
```

#### Strict Profile (`profiles/strict.yml`):
```yaml
name: "strict"
description: "High standards for production-ready projects"

review:
  coverage_threshold: 95
  test_timeout: "20m"
  skip_integration: false
  require_integration: true
  
linting:
  enabled: true
  config: "strict"
  security_checks: true
  performance_checks: true
  
documentation:
  require_package_docs: true
  require_function_docs: true
  
badges:
  enabled: true
  include: ["lint", "coverage", "go-version", "security", "docs"]
```

---

## 🔧 Core Review Pipeline

### Review Steps (Based on Current Script):
1. **🔍 Environment Check**
   - Go version validation
   - Git repository validation
   - Required tools availability

2. **🔍 Comprehensive Linting**
   - golangci-lint v2 with security checks
   - Code formatting validation
   - Import organization
   - Security vulnerability scanning

3. **🏗️ Build Validation**
   - Cross-platform build testing
   - Dependency resolution
   - Module validation

4. **🧪 Unit Tests**
   - Test execution with timeout
   - Race condition detection
   - Parallel test execution

5. **🔗 Integration Tests**
   - Tagged integration test execution
   - External dependency testing
   - End-to-end scenarios

6. **📊 Coverage Analysis**
   - Configurable coverage thresholds
   - Package-level coverage reporting
   - Coverage badge generation

7. **📚 Documentation Validation**
   - Go doc validation
   - README.md checks
   - API documentation verification

8. **🧹 Final Validation**
   - Git status checks
   - Uncommitted changes detection
   - Final cleanup

9. **🏷️ Badge Generation**
   - Dynamic badge creation
   - GitHub API integration
   - Security status monitoring

---

## 🐳 GitHub Action Implementation

### action.yml:
```yaml
name: 'Go Project Reviewer'
description: 'Comprehensive analysis for Go projects including linting, testing, coverage, and badge generation'
author: 'yourusername'

branding:
  icon: 'check-circle'
  color: 'green'

inputs:
  coverage-threshold:
    description: 'Minimum coverage percentage required'
    required: false
    default: '80'
  
  test-timeout:
    description: 'Test execution timeout'
    required: false
    default: '10m'
  
  skip-integration:
    description: 'Skip integration tests'
    required: false
    default: 'false'
  
  go-version:
    description: 'Go version to use'
    required: false
    default: '1.24'
  
  review-profile:
    description: 'Review profile (minimal, standard, strict)'
    required: false
    default: 'standard'
  
  badge-generation:
    description: 'Generate badge JSON files'
    required: false
    default: 'true'
  
  config-file:
    description: 'Path to custom config file'
    required: false
    default: '.go-validate.yml'

outputs:
  review-result:
    description: 'Overall review result (passed/failed)'
  coverage-percentage:
    description: 'Code coverage percentage achieved'
  lint-issues:
    description: 'Number of linting issues found'
  test-results:
    description: 'Test execution summary'

runs:
  using: 'docker'
  image: 'Dockerfile'
```

### Dockerfile:
```dockerfile
FROM golang:1.24-alpine

# Install system dependencies
RUN apk add --no-cache \
    bash \
    git \
    curl \
    jq \
    github-cli \
    make

# Install golangci-lint v2
RUN go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

# Copy validation scripts
COPY review.sh /usr/local/bin/review
COPY entrypoint.sh /entrypoint.sh
COPY profiles/ /usr/local/share/review/profiles/
COPY templates/ /usr/local/share/review/templates/

# Set permissions
RUN chmod +x /usr/local/bin/review /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
```

### entrypoint.sh:
```bash
#!/bin/bash
set -e

# Map GitHub Action inputs to environment variables
export COVERAGE_THRESHOLD="${INPUT_COVERAGE_THRESHOLD:-80}"
export TEST_TIMEOUT="${INPUT_TEST_TIMEOUT:-10m}"
export SKIP_INTEGRATION="${INPUT_SKIP_INTEGRATION:-false}"
export REVIEW_PROFILE="${INPUT_REVIEW_PROFILE:-standard}"
export BADGE_GENERATION="${INPUT_BADGE_GENERATION:-true}"
export CONFIG_FILE="${INPUT_CONFIG_FILE:-.go-validate.yml}"

# Set up Go version if specified
if [[ -n "${INPUT_GO_VERSION}" && "${INPUT_GO_VERSION}" != "1.24" ]]; then
    echo "🔧 Setting up Go ${INPUT_GO_VERSION}..."
    # Install specific Go version if needed
fi

# Load profile configuration
if [[ -f "/usr/local/share/review/profiles/${REVIEW_PROFILE}.yml" ]]; then
    echo "📋 Loading review profile: ${REVIEW_PROFILE}"
    # Load profile settings
fi

# Set GitHub Action outputs
set_output() {
    echo "$1=$2" >> $GITHUB_OUTPUT
}

# Run validation with CI mode
echo "🚀 Starting Go project review (GitHub Action mode)..."
if review ci; then
    set_output "review-result" "passed"
    echo "✅ Review completed successfully"
else
    set_output "review-result" "failed"
    echo "❌ Review failed"
    exit 1
fi

# Extract and set additional outputs
if [[ -f ".github/badges/coverage.json" ]]; then
    coverage=$(jq -r '.message' .github/badges/coverage.json | sed 's/%//')
    set_output "coverage-percentage" "$coverage"
fi

if [[ -f ".github/badges/golangci-lint.json" ]]; then
    lint_issues=$(jq -r '.message' .github/badges/golangci-lint.json | grep -oE '[0-9]+' || echo "0")
    set_output "lint-issues" "$lint_issues"
fi
```

---

## 🏷️ Badge System

### Generated Badges:
1. **golangci-lint.json** - Linting issues count
2. **coverage.json** - Code coverage percentage
3. **go-version.json** - Go version used
4. **last-updated.json** - Last commit date
5. **security.json** - Combined Dependabot + Code Scanning alerts
6. **review.json** - Overall review status

### Badge URLs for README:
```markdown
![Review](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/yourusername/yourrepo/main/.github/badges/review.json)
![Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/yourusername/yourrepo/main/.github/badges/coverage.json)
![Lint](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/yourusername/yourrepo/main/.github/badges/golangci-lint.json)
![Security](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/yourusername/yourrepo/main/.github/badges/security.json)
![Last Updated](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/yourusername/yourrepo/main/.github/badges/last-updated.json)
```

---

## 🚀 Development Roadmap

### Phase 1: Core Implementation
- [ ] Extract and generalize current review script
- [ ] Create repository structure
- [ ] Implement configuration system
- [ ] Add review profiles (minimal, standard, strict)
- [ ] Create installation script

### Phase 2: GitHub Action
- [ ] Create Dockerfile and action.yml
- [ ] Implement entrypoint script
- [ ] Add input/output handling
- [ ] Test with sample repositories
- [ ] Publish to GitHub Actions Marketplace

### Phase 3: Enhanced Features
- [ ] Custom review step support
- [ ] Plugin architecture
- [ ] Multiple output formats (JSON, XML, SARIF)
- [ ] Integration with popular CI/CD systems
- [ ] Performance benchmarking and optimization

### Phase 4: Community & Documentation
- [ ] Comprehensive documentation
- [ ] Usage examples and tutorials
- [ ] Community contribution guidelines
- [ ] Integration guides for popular frameworks
- [ ] Video tutorials and demos

### Phase 5: Advanced Features
- [ ] Web dashboard for validation results
- [ ] Slack/Discord notifications
- [ ] Historical trend analysis
- [ ] Team/organization-wide policies
- [ ] Custom badge templates

---

## 🎯 Key Benefits

### For Individual Developers:
- ✅ **Consistent Standards**: Same validation across all projects
- ✅ **Time Saving**: No setup required, just run one command
- ✅ **Beautiful Output**: Clear, colorful feedback with progress tracking
- ✅ **Flexible Configuration**: Adapt to project needs with profiles
- ✅ **Badge Generation**: Professional README badges automatically

### For Teams:
- ✅ **Standardization**: Enforce consistent code quality across team
- ✅ **CI/CD Integration**: Drop-in GitHub Action for all repositories
- ✅ **Visibility**: Real-time project health badges
- ✅ **Onboarding**: New team members get consistent tooling
- ✅ **Maintenance**: Single source of truth for validation logic

### For Open Source:
- ✅ **Community Standards**: Encourage best practices
- ✅ **Contributor Experience**: Clear review feedback
- ✅ **Project Health**: Visible quality metrics
- ✅ **Easy Adoption**: Simple integration for any Go project
- ✅ **Marketplace Presence**: Discoverable in GitHub Actions

---

## 📋 Success Metrics

### Adoption Metrics:
- GitHub Action usage (installations, runs)
- Local CLI downloads and usage
- Community contributions (issues, PRs, discussions)
- Documentation views and engagement

### Quality Metrics:
- Review accuracy and reliability
- Performance (execution time, resource usage)
- User satisfaction (feedback, ratings)
- Bug reports and resolution time

### Impact Metrics:
- Projects using review showing improved code quality
- Reduction in common Go project issues
- Community adoption of standardized practices
- Integration with other Go ecosystem tools

---

This specification provides a complete blueprint for creating a universal Go project analysis tool that serves both local development and CI/CD needs, making Go project review consistent, comprehensive, and accessible to the entire community.
