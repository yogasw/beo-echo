# BeoEcho Desktop Application

BeoEcho Desktop is a native desktop application built with Wails v2 that packages the BeoEcho web application (Svelte frontend + Go backend) into a single executable for Windows, macOS, and Linux.

## 📋 Prerequisites

### Required Software
- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Wails CLI v2.8.0+** - [Installation Guide](https://wails.io/docs/gettingstarted/installation)

### Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Verify Installation
```bash
wails doctor
```

## 🚀 Quick Start

### 1. Setup Development Environment
```bash
# Navigate to desktop directory
cd desktop

# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd ../frontend && npm install

# Return to desktop directory
cd ../desktop
```

### 2. Development Mode (with Hot Reload)
```bash
# Start development server with hot reload
wails dev
```

This will:
- Start the backend server
- Launch the desktop application with hot reload
- Watch for changes in both frontend and backend code
- Automatically reload when files change

### 3. Production Build
```bash
# Build for current platform
wails build

# Build for specific platform
wails build -platform darwin/amd64    # macOS Intel
wails build -platform darwin/arm64    # macOS Apple Silicon
wails build -platform windows/amd64   # Windows 64-bit
wails build -platform linux/amd64     # Linux 64-bit

# Build for all platforms
wails build -platform darwin/amd64,darwin/arm64,linux/amd64,windows/amd64
```

Built applications will be in `build/bin/` directory.

## 📁 Project Structure

```
desktop/
├── main.go              # Desktop application entry point
├── go.mod              # Go module configuration
├── go.sum              # Go dependencies
├── wails.json          # Wails configuration
├── assets/             # Embedded static assets
├── build/              # Build output directory
│   └── bin/           # Compiled binaries
└── frontend/          # Generated Wails frontend bindings
    └── wailsjs/       # JavaScript bridge for Wails APIs
```

## ⚙️ Configuration

### Wails Configuration (`wails.json`)
- **Frontend Directory**: `../frontend` - Points to the main frontend source
- **Build Directory**: `../frontend/build` - Where compiled frontend assets are located
- **Backend**: Current directory (`.`) with `main.go` as entry point

### Desktop-Specific Features
- **Native Menus**: Desktop menu bar integration
- **Backend Status**: Real-time backend health monitoring
- **Desktop Environment**: Runs in `~/.beoecho` directory
- **Database**: SQLite stored in user's home directory

## 🛠️ Development Commands

### Frontend Development
```bash
# Build frontend for desktop (from frontend directory)
cd ../frontend
npm run build:desktop

# Development with desktop mode enabled
VITE_DESKTOP_MODE=true npm run dev
```

### Backend Development
```bash
# Run backend tests (from backend directory)
cd ../backend
go test ./...

# Run backend directly
go run main.go
```

### Desktop Development
```bash
# Development mode with hot reload
wails dev

# Generate Wails bindings
wails generate module

# Build application
wails build

# Clean build artifacts
rm -rf build/
```

## 🧪 Testing

### Test Backend Integration
```bash
cd ../backend
go test ./...
```

### Test Frontend Build
```bash
cd ../frontend
npm run build:desktop
ls -la build/  # Verify build artifacts
```

### Test Desktop Application
```bash
# After building
./build/bin/BeoEcho\ Desktop.app/Contents/MacOS/BeoEcho  # macOS
./build/bin/BeoEcho.exe                                   # Windows
./build/bin/BeoEcho                                       # Linux
```

## 📦 Distribution

### Build for All Platforms
```bash
# Single command to build for all supported platforms
wails build -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64
```

### Platform-Specific Builds
```bash
# macOS (Universal Binary)
wails build -platform darwin/amd64,darwin/arm64

# Windows
wails build -platform windows/amd64

# Linux
wails build -platform linux/amd64
```

### Output Locations
- **macOS**: `build/bin/BeoEcho Desktop.app/`
- **Windows**: `build/bin/BeoEcho.exe`
- **Linux**: `build/bin/BeoEcho`

## 🔧 Troubleshooting

### Common Issues

**1. "no index.html found" Error**
```bash
# Ensure frontend is built first
cd ../frontend && npm run build:desktop
cd ../desktop && wails build
```

**2. Backend Dependencies Issues**
```bash
# Update Go dependencies
go mod tidy
go mod download
```

**3. Frontend Dependencies Issues**
```bash
# Clean and reinstall
cd ../frontend
rm -rf node_modules package-lock.json
npm install
```

**4. Wails CLI Issues**
```bash
# Update Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# Check installation
wails doctor
```

### Build Issues

**Clean Build Environment**
```bash
# Clean all build artifacts
rm -rf build/
rm -rf ../frontend/build/
rm -rf ../frontend/.svelte-kit/

# Rebuild everything
cd ../frontend && npm run build:desktop
cd ../desktop && wails build
```

**Dependency Issues**
```bash
# Reset Go modules
go clean -modcache
go mod tidy

# Reset Node modules
cd ../frontend
rm -rf node_modules
npm install
```

## 📋 Development Workflow

### Daily Development
1. Start development mode: `wails dev`
2. Make changes to frontend (`../frontend/src/`) or backend (`../backend/src/`)
3. Changes are automatically reloaded in development mode

### Testing Changes
1. Build frontend: `cd ../frontend && npm run build:desktop`
2. Build desktop: `cd ../desktop && wails build`
3. Test the built application

### Releasing
1. Ensure all tests pass: `cd ../backend && go test ./...`
2. Build frontend: `cd ../frontend && npm run build:desktop`
3. Build for all platforms: `cd ../desktop && wails build -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64`
4. Test binaries on target platforms
5. Package for distribution

## 🔗 Related Documentation

- [Wails Documentation](https://wails.io/docs/)
- [BeoEcho Backend README](../backend/README.md)
- [BeoEcho Frontend Setup](../frontend/README.md)
- [Desktop Development Guide](../docs/Desktop_Application_Guide.md)

## 🆘 Support

For issues related to:
- **Desktop Application**: Check this README and Wails documentation
- **Backend API**: See `../backend/README.md`
- **Frontend UI**: See `../frontend/README.md`
- **General Setup**: See main project `../README.md`
