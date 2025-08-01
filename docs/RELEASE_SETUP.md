# Release Setup Guide

This document explains how to create releases with platform-specific builds, code signing, and automated publishing.

## Prerequisites

1. **GitHub Repository**: Ensure your project is hosted on GitHub
2. **macOS Developer Account**: Required for code signing certificates (optional)
3. **CHANGELOG.md**: Required file in repository root for release validation

## Required GitHub Secrets (Optional for Code Signing)

Configure the following secrets in your GitHub repository settings (`Settings > Secrets and variables > Actions`) if you want signed macOS binaries:

### Code Signing Secrets
- `MACOS_CERTIFICATE`: Base64-encoded .p12 certificate file
- `MACOS_CERTIFICATE_PWD`: Password for the .p12 certificate
- `MACOS_CERTIFICATE_NAME`: Name of the certificate (Developer ID Application: Your Name)
- `KEYCHAIN_PASSWORD`: Password for temporary keychain

### Notarization Secrets (Required for Release Tags Only)
- `APPLE_ID`: Your Apple ID email
- `APPLE_PASSWORD`: App-specific password for your Apple ID
- `APPLE_TEAM_ID`: Your Apple Developer Team ID

## Setting up Code Signing

### Export Certificate from Xcode/Keychain

**Option A: From Xcode (Recommended)**
1. Open Xcode > Settings > Accounts
2. Select your Apple ID > Manage Certificates
3. Right-click your "Developer ID Application" certificate > Export Certificate
4. Save as `.p12` format and set a strong password

**Option B: From Keychain Access**
1. Open Keychain Access on macOS
2. Find your "Developer ID Application" certificate
3. Right-click and select "Export"
4. Choose `.p12` format and set a strong password

**Convert to Base64:**
```bash
# Convert certificate to base64
base64 -i YourCertificate.p12 | pbcopy
```

**Add to GitHub Secrets:**
- `MACOS_CERTIFICATE`: The base64 string from clipboard
- `MACOS_CERTIFICATE_PWD`: The password you set during export

## Release Process

### Step 1: Update CHANGELOG.md

Before creating a release, you **must** add an entry to `CHANGELOG.md`:

```markdown
## v1.0.0 - 2024-01-15

### Added
- New feature description
- Another new feature

### Changed
- Modified behavior description

### Fixed
- Bug fix description
```

Required format: `## v1.0.0` (without brackets)

### Step 2: Create and Push Tag

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Step 3: Automated Release Process

GitHub Actions will automatically:

1. **Run Tests**: Ensure all tests pass
2. **Build Platform Binaries**:
   - Linux: amd64, arm64 (built on ubuntu-latest)
   - macOS: amd64, arm64 (built on macos-latest)
3. **Sign macOS Binaries**: Using codesign (if certificates configured)
4. **Notarize Binaries**: For release tags only, using xcrun notarytool
4. **Validate Changelog**: Ensure version exists in CHANGELOG.md
5. **Create Release**: Package binaries and upload to GitHub
6. **Generate Checksums**: SHA256 checksums for all binaries

### Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon) - optionally signed

## Configuration Files

- `.github/workflows/ci.yml`: Combined CI and release workflow with platform-specific builds
- `Makefile`: Build targets for different platforms using `dist/$os-$arch/outline` structure
- `CHANGELOG.md`: Required for release validation

## How It Works

1. **Platform-Native Builds**: Each platform builds its own binaries to avoid CGO cross-compilation issues
2. **Code Signing**: macOS binaries are signed using codesign with temporary keychain
3. **Conditional Notarization**: Only release builds (tags) are notarized to save time on regular CI
3. **Changelog Validation**: Release fails if version tag not found in CHANGELOG.md
4. **Artifact Packaging**: Binaries are packaged as tar.gz archives with checksums
5. **GitHub Release**: Archives are automatically attached to the release

## Troubleshooting

### Release Fails: "CHANGELOG.md not found"
- Create a `CHANGELOG.md` file in your repository root
- Add an entry for your version tag

### Release Fails: "Version not found in CHANGELOG.md"  
- Ensure your version tag (e.g., `v1.0.0`) appears as a heading in CHANGELOG.md
- Required format: `## v1.0.0` (without brackets)

### Code Signing Issues (Optional)
- Ensure certificate is valid and not expired
- Check that certificate has "Developer ID Application" type
- Verify `MACOS_CERTIFICATE` and `MACOS_CERTIFICATE_PWD` secrets are set correctly

### Build Issues
- Ensure Go version matches what's specified in workflows (currently 1.24.5)
- Check that all dependencies support the target platforms
- CGO issues: The workflow uses platform-native builds to avoid cross-compilation problems

## Testing Locally

Test builds locally to debug issues:

```bash
# Build for current platform
make build

# Build for specific platform (uses Makefile)
GOOS=linux GOARCH=amd64 make build
GOOS=darwin GOARCH=arm64 make build

# Run tests
make test
```

## Quick Release Checklist

1. ✅ Update `CHANGELOG.md` with new version
2. ✅ Commit and push changes
3. ✅ Create and push version tag: `git tag v1.0.0 && git push origin v1.0.0`
4. ✅ Watch GitHub Actions for build status
5. ✅ Verify release artifacts are uploaded

## Binary Structure

Binaries are now built in a structured format:
- `dist/linux-amd64/outline`
- `dist/linux-arm64/outline`
- `dist/darwin-amd64/outline` 
- `dist/darwin-arm64/outline`

This allows the installer to rename the binary to just `outline` while preserving notarization.

## Notarization Notes

- **Releases only**: Notarization only runs for git tags starting with `v` to save time
- **Stapled tickets**: Notarization tickets are stapled to binaries for offline verification
- **Path independent**: Stapled binaries can be renamed and moved without losing notarization

## Security Notes

- Code signing is optional - releases work without certificates
- Keep certificates and passwords secure if using signing
- Monitor certificate expiration dates
- Notarization requires valid Apple Developer account