#!/bin/bash

# Generate WinGet manifest files for the g package
# WinGet manifests consist of multiple YAML files in a specific directory structure

latest=$(git describe --abbrev=0 --tags 2>/dev/null || echo "v0.30.0")
version=${latest#v}  # Remove 'v' prefix
manifest_dir="manifests/e/Equationzhao/G/${version}"

# Create directory structure
mkdir -p "$manifest_dir"

# Generate checksums for Windows binaries
hash64=""
hash386=""
hasharm64=""

if [ -f "../build/g-windows-amd64.exe" ]; then
    hash64=$(shasum -a 256 ../build/g-windows-amd64.exe | cut -d ' ' -f 1)
fi

if [ -f "../build/g-windows-386.exe" ]; then
    hash386=$(shasum -a 256 ../build/g-windows-386.exe | cut -d ' ' -f 1)
fi

if [ -f "../build/g-windows-arm64.exe" ]; then
    hasharm64=$(shasum -a 256 ../build/g-windows-arm64.exe | cut -d ' ' -f 1)
fi

# Generate the version manifest (Equationzhao.G.yaml)
cat > "$manifest_dir/Equationzhao.G.yaml" << EOF
# Created using winget-releaser v2.0.0
# yaml-language-server: \$schema=https://aka.ms/winget-manifest.version.1.6.0.schema.json

PackageIdentifier: Equationzhao.G
PackageVersion: ${version}
DefaultLocale: en-US
ManifestType: version
ManifestVersion: 1.6.0
EOF

# Generate the installer manifest (Equationzhao.G.installer.yaml)
cat > "$manifest_dir/Equationzhao.G.installer.yaml" << EOF
# Created using winget-releaser v2.0.0
# yaml-language-server: \$schema=https://aka.ms/winget-manifest.installer.1.6.0.schema.json

PackageIdentifier: Equationzhao.G
PackageVersion: ${version}
InstallerType: portable
Commands:
- g
ReleaseDate: $(date -u +%Y-%m-%d)
Installers:
EOF

# Add x64 installer if hash is available
if [ -n "$hash64" ]; then
cat >> "$manifest_dir/Equationzhao.G.installer.yaml" << EOF
- Architecture: x64
  InstallerUrl: https://github.com/Equationzhao/g/releases/download/v${version}/g-windows-amd64.exe
  InstallerSha256: ${hash64}
EOF
fi

# Add x86 installer if hash is available  
if [ -n "$hash386" ]; then
cat >> "$manifest_dir/Equationzhao.G.installer.yaml" << EOF
- Architecture: x86
  InstallerUrl: https://github.com/Equationzhao/g/releases/download/v${version}/g-windows-386.exe
  InstallerSha256: ${hash386}
EOF
fi

# Add arm64 installer if hash is available
if [ -n "$hasharm64" ]; then
cat >> "$manifest_dir/Equationzhao.G.installer.yaml" << EOF
- Architecture: arm64
  InstallerUrl: https://github.com/Equationzhao/g/releases/download/v${version}/g-windows-arm64.exe
  InstallerSha256: ${hasharm64}
EOF
fi

# Complete the installer manifest
cat >> "$manifest_dir/Equationzhao.G.installer.yaml" << EOF
ManifestType: installer
ManifestVersion: 1.6.0
EOF

# Generate the locale manifest (Equationzhao.G.locale.en-US.yaml)
cat > "$manifest_dir/Equationzhao.G.locale.en-US.yaml" << EOF
# Created using winget-releaser v2.0.0
# yaml-language-server: \$schema=https://aka.ms/winget-manifest.defaultLocale.1.6.0.schema.json

PackageIdentifier: Equationzhao.G
PackageVersion: ${version}
PackageLocale: en-US
Publisher: Equationzhao
PublisherUrl: https://github.com/Equationzhao
PublisherSupportUrl: https://github.com/Equationzhao/g/issues
Author: Equationzhao
PackageName: G
PackageUrl: https://g.equationzhao.space
License: MIT
LicenseUrl: https://github.com/Equationzhao/g/blob/master/LICENSE
Copyright: Copyright (c) 2024 Equationzhao
ShortDescription: A feature-rich, customizable, and cross-platform ls alternative
Description: |-
  G is a feature-rich, customizable, and cross-platform ls alternative.
  Experience enhanced visuals with type-specific icons, various layout options, and git status integration.
  
  Features:
  - Type-specific icons for files and directories
  - Multiple layout options (oneline, grid, across, zero, comma, table, json, markdown, tree, recurse)
  - Git status integration
  - Cross-platform support (Linux, Windows, macOS)
  - Advanced sorting options
  - Customizable themes
  - MIME type detection
  - Size formatting options
Moniker: g
Tags:
- cli
- command-line
- directory
- file-manager
- ls
- tool
- utilities
ReleaseNotes: See https://github.com/Equationzhao/g/releases/tag/v${version}
ReleaseNotesUrl: https://github.com/Equationzhao/g/releases/tag/v${version}
Documentations:
- DocumentLabel: README
  DocumentUrl: https://github.com/Equationzhao/g/blob/master/README.md
ManifestType: defaultLocale
ManifestVersion: 1.6.0
EOF

echo "WinGet manifest generated in: $manifest_dir"
echo "Files created:"
ls -la "$manifest_dir"

echo ""
echo "To submit to winget-pkgs repository:"
echo "1. Fork https://github.com/microsoft/winget-pkgs"
echo "2. Copy the manifests/ directory to your fork"
echo "3. Create a pull request"
echo ""
echo "For testing locally:"
echo "winget install --manifest $manifest_dir"