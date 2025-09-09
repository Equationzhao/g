# g - Enhanced ls Alternative

**Always reference these instructions first and fallback to search or bash commands only when you encounter unexpected information that does not match the info here.**

g is a feature-rich, customizable, and cross-platform `ls` alternative written in Go. It provides enhanced file listing with icons, Git integration, multiple layout options, and advanced sorting capabilities.

## Working Effectively

### Prerequisites and Setup
- Requires Go version >= 1.24.0 (project uses Go 1.24.0)
- Works on Linux, Windows, and macOS
- Repository uses `go mod` for dependency management

### Bootstrap and Build Commands
**ALWAYS run these commands in sequence for a fresh setup:**

```bash
# 1. Verify Go version (must be >= 1.24)
go version

# 2. Download dependencies and build (NEVER CANCEL: first build takes ~20 seconds)
time go build -v .
```

**Build timing expectations:**
- **NEVER CANCEL**: Initial build from scratch: ~12 seconds. Set timeout to 60+ seconds.
- **NEVER CANCEL**: Subsequent builds: ~1.5 seconds. Set timeout to 30+ seconds.
- **NEVER CANCEL**: Tests: ~11 seconds. Set timeout to 30+ seconds.
- **NEVER CANCEL**: Linting: ~60 seconds. Set timeout to 180+ seconds.

### Build Variants
The project supports multiple build configurations using Go build tags:

```bash
# Lite build (minimal dependencies, 7.4MB binary)
CGO_ENABLED=0 go build -ldflags="-s -w" -o g-lite .

# Full build (all features, 8.1MB binary)  
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy mounts" -o g-full .

# Fuzzy search only
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy" -o g-fuzzy .

# Mounts support only
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="mounts" -o g-mounts .
```

**Build tags:**
- `fuzzy`: Enables fuzzy search and path indexing (~500KB size impact)
- `mounts`: Enables mount point detection (~200KB size impact)

### Testing
```bash
# Run all tests (NEVER CANCEL: takes ~11 seconds)
time go test -v ./...
```

### Code Quality and CI Validation
**Always run these before committing changes:**

```bash
# Install formatting tool
go install mvdan.cc/gofumpt@latest

# Install linter (NEVER CANCEL: takes ~30 seconds to install)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Format code (required by CI)
export PATH=$PATH:~/go/bin
gofumpt -l -extra .  # Check formatting
gofumpt -w -extra .  # Fix formatting

# Lint code (NEVER CANCEL: takes ~60 seconds)
golangci-lint run ./... --timeout=3m
```

**CI Pipeline:** The project has three GitHub workflows that must pass:
- `go.yml`: Multi-platform builds and tests
- `gofumpt.yml`: Code formatting verification
- `lint.yml`: Static code analysis

## Application Usage and Validation

### Basic Functionality Testing
**Always test these scenarios after making changes:**

```bash
# 1. Test basic listing
./g .

# 2. Test with icons and formatting
./g --icon --long .

# 3. Test tree view
./g --tree --icon .

# 4. Test Git integration (if in git repo)
./g --git --icon .

# 5. Test table format
./g --table --size --time .

# 6. Test JSON output
./g --json . | head -10
```

### Comprehensive Validation Scenario
**Always run this complete end-to-end test after making significant changes:**

```bash
# 1. Create test directory and files
mkdir -p /tmp/g-validation-test && cd /tmp/g-validation-test
mkdir -p subdir
echo "test content" > file1.txt  
echo "another test" > subdir/file2.txt

# 2. Test core functionality with actual files
/path/to/g --tree --icon --size .  # Should show tree structure with sizes
/path/to/g --table --time .         # Should show tabular output with timestamps
/path/to/g --json . | jq '.'        # Should produce valid JSON
/path/to/g --recurse .              # Should list all files recursively

# 3. Verify output contains expected elements
# - Icons should be displayed if terminal supports them
# - File sizes should be human-readable (B, KiB, etc.)
# - Tree structure should use box-drawing characters
# - JSON should be valid and parseable

# 4. Clean up
cd /tmp && rm -rf g-validation-test
```

### Advanced Feature Validation
```bash
# Test fuzzy search (requires fuzzy build tag)
./g-full --fuzzy pattern

# Test mount information (requires mounts build tag)  
./g-full --mounts .

# Test recursive listing
./g --recurse directory/

# Test shell integration generation
./g --init bash
./g --init zsh
./g --init fish
```

### Performance Expectations
- Basic listing: ~0.015 seconds
- Tree view: < 1 second for typical directories
- Recursive operations: varies by directory size
- Application startup: Nearly instant

## Key Project Structure

### Repository Layout
```
/home/runner/work/g/g/
├── main.go              # Main entry point
├── go.mod               # Go module definition
├── justfile             # Build automation (just command)
├── internal/            # Internal Go packages
│   ├── cli/            # Command line interface
│   ├── display/        # Output formatting
│   ├── git/            # Git integration
│   ├── theme/          # Color themes
│   └── ...
├── .github/workflows/  # CI/CD pipelines
├── completions/        # Shell completions (bash, zsh, fish)
├── docs/               # Documentation
├── script/             # Development scripts
└── man/                # Manual pages
```

### Important Files
- `internal/cli/g.go`: Main CLI command definitions and logic
- `internal/display/`: All output format implementations
- `docs/BuildOption.md`: Detailed build configuration options
- `CONTRIBUTING.md`: Development workflow and commit standards

## Common Tasks

### Development Workflow
```bash
# 1. Make changes to source files
# 2. Format code (required by CI)
export PATH=$PATH:~/go/bin
gofumpt -w -extra .

# 3. Build and test (NEVER CANCEL: full workflow takes ~90 seconds)
go build .
go test -v ./...

# 4. Validate with lint (NEVER CANCEL: takes ~60 seconds)
golangci-lint run ./... --timeout=3m

# 5. Test functionality with real scenario
mkdir -p /tmp/validation && cd /tmp/validation
echo "test" > sample.txt
/path/to/g --tree --icon --size .
cd - && rm -rf /tmp/validation

# Complete workflow timing: ~90 seconds total
```

### Build Tools Reference
The project uses `justfile` for build automation, but core Go commands work directly:

```bash
# Instead of: just build-full
CGO_ENABLED=0 go build -ldflags="-s -w" -tags="fuzzy mounts" -o g-full .

# Instead of: just test  
go test -v ./...

# Instead of: just precheck
gofumpt -w -extra . && golangci-lint run ./...
```

### Shell Integration Setup
```bash
# Generate shell aliases
./g --init bash    # For bash
./g --init zsh     # For zsh  
./g --init fish    # For fish
./g --init powershell  # For PowerShell
```

## Troubleshooting

### Common Issues
- **"Go version too low"**: Ensure Go >= 1.24.0 is installed
- **Build fails**: Run `go mod tidy` to sync dependencies
- **Tests fail**: Ensure working directory is repository root
- **Lint fails**: Install golangci-lint compatible with Go 1.24+

### Platform-Specific Notes
- **Linux**: Full functionality available
- **macOS**: All features work, CGO enabled for Darwin builds in CI
- **Windows**: Core functionality works, some file attribute features limited

### Binary Size Expectations
- Default build: ~10.7MB
- Lite build (`-ldflags="-s -w"`): 7.4MB  
- Full build with all tags: 8.1MB

## Critical Reminders
- **NEVER CANCEL** any build, test, or lint command - full workflow takes up to 90 seconds
- Always set timeouts of 60+ seconds for builds, 30+ seconds for tests, 180+ seconds for linting
- Code formatting with gofumpt is mandatory - CI will fail without proper formatting
- Git status integration only works when run from within a Git repository
- Always test both lite and full build variants when making significant changes
- Use `/path/to/g` in validation scripts to reference your built binary location
- JSON output validation: always pipe through `jq '.'` to verify structure
- Test with actual files in `/tmp` directories for realistic validation scenarios

## Version Information
- Current version: v0.31.0
- Go version requirement: >= 1.24.0
- License: MIT License
- Repository: https://github.com/Equationzhao/g