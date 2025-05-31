# BeoEcho Desktop Application

This document describes how to build and develop the desktop version of BeoEcho using Wails v2.

## Overview

BeoEcho Desktop is a cross-platform desktop application that packages the existing web frontend and Go backend into a native desktop app using Wails v2. The desktop version provides:

- **Native desktop experience** with system integration
- **Embedded backend server** - no need to run separate backend
- **Local SQLite database** - data stored locally on user's machine
- **Cross-platform support** - Windows, macOS, and Linux
- **Auto-startup** - backend starts automatically with the app
- **Desktop-specific features** - system tray, native dialogs, etc.

## Architecture

```
┌─────────────────────────────────────┐
│           Desktop App               │
├─────────────────────────────────────┤
│  Frontend (Svelte + Wails Bridge)   │  ← Web UI with desktop adaptations
├─────────────────────────────────────┤
│         Wails Runtime               │  ← Native app framework
├─────────────────────────────────────┤
│    Embedded Go Backend Server      │  ← Same backend as web version
├─────────────────────────────────────┤
│      Local SQLite Database         │  ← Data stored in user directory
└─────────────────────────────────────┘
```

## Prerequisites

Before you can build the desktop version, ensure you have:

### Required Software

1. **Go 1.21+** - Backend development
   ```bash
   # Check Go version
   go version
   ```

2. **Node.js 18+** - Frontend development
   ```bash
   # Check Node.js version
   node --version
   npm --version
   ```

3. **Wails v2** - Desktop app framework
   ```bash
   # Install Wails (will be installed automatically by scripts)
   go install github.com/wailsapp/wails/v2/cmd/wails@v2.8.0
   ```

### Platform-Specific Requirements

#### macOS
- **Xcode Command Line Tools**
  ```bash
  xcode-select --install
  ```

#### Windows
- **WebView2 Runtime** (usually pre-installed on Windows 10/11)
- **GCC compiler** (via TDM-GCC or similar)

#### Linux
- **WebKitGTK** development libraries
  ```bash
  # Ubuntu/Debian
  sudo apt-get install libwebkit2gtk-4.0-dev
  
  # Fedora
  sudo dnf install webkit2gtk3-devel
  
  # Arch Linux
  sudo pacman -S webkit2gtk
  ```

## Quick Start

### Method 1: Using the Desktop Script (Recommended)

The easiest way to get started is using the provided script:

```bash
# Setup development environment (one-time setup)
./desktop.sh setup

# Start development mode with hot reload
./desktop.sh dev

# Build for current platform
./desktop.sh build

# Build for all platforms
./desktop.sh build-all

# Package for distribution
./desktop.sh package
```

### Method 2: Using Make

Alternatively, you can use the Makefile:

```bash
# Setup environment
make -f Makefile.desktop setup-desktop

# Development mode
make -f Makefile.desktop dev-desktop

# Build
make -f Makefile.desktop build-desktop

# Clean
make -f Makefile.desktop clean-desktop
```

### Method 3: Direct Wails Commands

For advanced users who want direct control:

```bash
# Setup dependencies
cd frontend && npm install && cd ..
cd desktop && go mod tidy && cd ..

# Development mode
wails dev

# Build
cd frontend && npm run build:desktop && cd ..
wails build
```

## Development Workflow

### 1. Initial Setup

```bash
# Clone the repository (if not already done)
git clone <repository-url>
cd beo-echo

# Setup desktop development environment
./desktop.sh setup
```

This will:
- Install Wails if not present
- Install frontend dependencies
- Setup Go modules for desktop
- Verify all requirements

### 2. Development Mode

```bash
# Start development with hot reload
./desktop.sh dev
```

This starts the Wails development server with:
- **Hot reload** for frontend changes
- **Automatic rebuild** for backend changes
- **Live debugging** with Chrome DevTools
- **Native desktop window** for testing

### 3. Building for Production

```bash
# Build for current platform
./desktop.sh build

# Or build for all platforms
./desktop.sh build-all
```

Build artifacts will be in `build/bin/`:
- **macOS**: `BeoEcho.app`
- **Windows**: `BeoEcho.exe`
- **Linux**: `BeoEcho`

## Project Structure

```
/
├── wails.json              # Wails configuration
├── desktop.sh              # Development script
├── Makefile.desktop        # Build automation
├── desktop/                # Desktop-specific code
│   ├── main.go            # Wails app entry point
│   ├── go.mod             # Desktop Go module
│   └── app.json           # App metadata
├── frontend/              # Existing Svelte frontend
│   ├── src/lib/utils/
│   │   ├── wails-bridge.js     # Wails JavaScript bridge
│   │   └── desktopConfig.ts    # Desktop configuration
│   └── src/lib/components/desktop/
│       ├── DesktopMenuBar.svelte   # Desktop menu bar
│       └── BackendStatus.svelte    # Backend status indicator
└── backend/               # Existing Go backend (unchanged)
```

## Configuration

### Wails Configuration (`wails.json`)

Key configurations for the desktop app:

```json
{
  "name": "BeoEcho Desktop",
  "frontend": {
    "dir": "./frontend",
    "build": "npm run build:desktop"
  },
  "backend": {
    "dir": "./desktop"
  },
  "info": {
    "productName": "BeoEcho Desktop",
    "companyName": "BeoEcho Team"
  }
}
```

### Frontend Desktop Configuration

The frontend automatically detects desktop mode and adapts:

```typescript
// Check if running in desktop mode
import { isDesktopMode } from '$lib/utils/desktopConfig';

if (isDesktopMode()) {
  // Desktop-specific behavior
  // API calls go to embedded backend
  // Show desktop-specific UI elements
}
```

### Backend Configuration

The desktop backend:
- **Auto-configures** data directories in user's home folder (`~/.beoecho/`)
- **Starts embedded server** on localhost:3600
- **Uses SQLite database** for local storage
- **Same API endpoints** as web version

## Desktop-Specific Features

### 1. Desktop Menu Bar

Shows application information and quick actions:
- App version and platform info
- Links to documentation and GitHub
- About dialog
- Quit application

### 2. Backend Status Indicator

Monitors the embedded backend server:
- Real-time health check
- Connection status
- Port information
- Error indicators

### 3. Native Integration

- **System tray** support (future feature)
- **Native file dialogs** for imports/exports
- **Desktop notifications** for important events
- **Auto-startup** option (future feature)

### 4. Offline Operation

- **No internet required** for core functionality
- **Local data storage** in user directory
- **Embedded assets** for complete offline operation

## Building for Distribution

### Single Platform Build

```bash
# Build for current platform
./desktop.sh build

# Platform-specific builds
./desktop.sh build     # Current platform
make -f Makefile.desktop build-desktop-windows
make -f Makefile.desktop build-desktop-mac
make -f Makefile.desktop build-desktop-linux
```

### Cross-Platform Build

```bash
# Build for all platforms at once
./desktop.sh build-all
```

This creates binaries for:
- **macOS**: Intel (amd64) and Apple Silicon (arm64)
- **Windows**: 64-bit (amd64)
- **Linux**: 64-bit (amd64)

### Application Packaging

```bash
# Package for distribution
./desktop.sh package
```

This creates a `dist/` folder with:
- **Packaged application** ready for distribution
- **Platform-appropriate format** (.app, .exe, binary)

## Testing

### Running Tests

```bash
# Run all tests (backend + frontend)
./desktop.sh test

# Backend tests only
cd backend && go test ./...

# Frontend tests only (if available)
cd frontend && npm test
```

### Manual Testing

1. **Start development mode**: `./desktop.sh dev`
2. **Test desktop features**:
   - Menu bar functionality
   - Backend status monitoring
   - API connectivity
   - Data persistence
3. **Test cross-platform** builds on target platforms

## Troubleshooting

### Common Issues

#### 1. Wails Installation Failed
```bash
# Manual Wails installation
go install github.com/wailsapp/wails/v2/cmd/wails@v2.8.0

# Verify installation
wails version
```

#### 2. Frontend Build Errors
```bash
# Clean and reinstall frontend dependencies
cd frontend
rm -rf node_modules package-lock.json
npm install
npm run build:desktop
```

#### 3. Backend Connection Issues
- Check if port 3600 is available
- Verify backend health endpoint: `curl http://localhost:3600/api/health`
- Check desktop logs for backend startup errors

#### 4. Cross-Platform Build Issues
- Ensure target platform dependencies are installed
- Use platform-specific build commands if cross-compilation fails
- Check Wails documentation for platform-specific requirements

### Debug Mode

Enable debug logging:

```bash
# Development with debug info
WAILS_DEBUG=1 ./desktop.sh dev

# Build with debug info
WAILS_DEBUG=1 ./desktop.sh build
```

### Getting Help

1. **Check logs** in the terminal output
2. **Review Wails documentation**: https://wails.io/docs
3. **Check GitHub issues** for similar problems
4. **Use browser DevTools** in development mode (F12)

## Deployment

### macOS
- **Code signing** required for distribution outside Mac App Store
- **Notarization** required for macOS 10.15+
- Consider **Mac App Store** distribution

### Windows
- **Code signing** recommended for trust
- **NSIS installer** can be generated automatically
- Consider **Microsoft Store** distribution

### Linux
- **AppImage** format for universal distribution
- **Snap packages** for Ubuntu users
- **Flatpak** for other distributions

## Advanced Configuration

### Custom Build Flags

```bash
# Build with custom flags
wails build -ldflags="-X main.version=1.0.0" -tags=production

# Optimize binary size
wails build -upx -upxflags="-9"

# Obfuscate binary
wails build -garbleargs="-literals -tiny"
```

### Environment Variables

- `WAILS_DEBUG=1` - Enable debug mode
- `VITE_DESKTOP_MODE=true` - Frontend desktop mode
- `GO_ENV=production` - Backend production mode

## Future Enhancements

Planned desktop-specific features:

1. **System Tray Integration**
   - Minimize to tray
   - Quick actions from tray menu
   - Status indicators

2. **Auto-Update System**
   - Automatic update checking
   - Background downloads
   - User-prompted updates

3. **Enhanced Native Integration**
   - File associations for mock files
   - Protocol handler for beoecho:// URLs
   - Desktop drag-and-drop support

4. **Performance Optimizations**
   - Faster startup time
   - Reduced memory usage
   - Background processing

## Contributing

When contributing to desktop functionality:

1. **Test on multiple platforms** before submitting PRs
2. **Follow desktop UI guidelines** for each platform
3. **Update documentation** for new features
4. **Add tests** for desktop-specific functionality

## License

Same license as the main BeoEcho project.
