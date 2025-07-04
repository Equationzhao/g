# Build Configuration

This document describes the optional features that can be enabled or disabled during compilation to control the binary size.

## Build Tags

The `g` CLI tool supports conditional compilation using Go build tags to include or exclude optional features:

### `fuzzy` tag
- **Purpose**: Enables fuzzy search and path indexing functionality
- **Dependencies**: `github.com/syndtr/goleveldb`, `github.com/sahilm/fuzzy`
- **Size impact**: ~500KB
- **Features affected**: 
  - `--fuzzy` flag for fuzzy path searching
  - Path indexing and index management commands
- **Usage**: `go build -tags="fuzzy" .`

### `mounts` tag  
- **Purpose**: Enables mount point detection and display
- **Dependencies**: `github.com/shirou/gopsutil/v3`
- **Size impact**: ~200KB  
- **Features affected**:
  - `--mounts` flag to show mount details for files
- **Usage**: `go build -tags="mounts" .`

## Build Examples

### Lite build (minimal size)
```bash
go build -ldflags="-s -w" -o g-lite .
```
- Size: ~7.4MB/7.0MiB for macOS
- Features: Core functionality only (no fuzzy search, no mount info)

### Full build (all features)
```bash  
go build -ldflags="-s -w" -tags="fuzzy mounts" -o g-full .
```
- Size: ~8.1MB/7.7MiB for macOS
- Features: All optional features enabled

### Custom builds
```bash
# Only fuzzy search
go build -ldflags="-s -w" -tags="fuzzy" -o g-fuzzy .

# Only mounts  
go build -ldflags="-s -w" -tags="mounts" -o g-mounts .
```

## Behavior without optional features

### Without `fuzzy` tag:
- `--fuzzy` flag is silently ignored (no error)
- No path indexing occurs
- Fuzzy path matching falls back to exact path matching

### Without `mounts` tag:
- `--mounts` flag is silently ignored (no mount info displayed)
- No system partition scanning occurs

This approach allows users to choose between a smaller binary size and optional functionality based on their needs.