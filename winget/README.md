# WinGet Manifest Generation

This directory contains scripts to generate WinGet package manifests for the G project.

## Usage

To generate WinGet manifests:

```bash
./winget.sh
```

This script will:
1. Detect the latest git tag version
2. Generate checksums for Windows binaries (amd64, 386, arm64)
3. Create WinGet manifest files in the correct directory structure

## Generated Files

The script creates three manifest files:

- `Equationzhao.G.yaml` - Version manifest
- `Equationzhao.G.installer.yaml` - Installer manifest with download URLs and checksums
- `Equationzhao.G.locale.en-US.yaml` - Package metadata and description

## Submitting to WinGet

To submit the package to the official Microsoft winget-pkgs repository:

1. Fork https://github.com/microsoft/winget-pkgs
2. Copy the generated `manifests/` directory to your fork
3. Create a pull request to the main repository

## Local Testing

To test the manifest locally (requires Windows with WinGet):

```powershell
winget install --manifest manifests/e/Equationzhao/G/<version>
```

## Requirements

- Git repository with tags
- Windows binaries built in `../build/` directory:
  - `g-windows-amd64.exe`
  - `g-windows-386.exe` 
  - `g-windows-arm64.exe`
- `shasum` command for generating checksums