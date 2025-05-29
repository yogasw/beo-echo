# BeoEcho Desktop Application

A **cross-platform native desktop application** built with **Wails v2** that packages the complete BeoEcho mock API service (Svelte 5 frontend + Go backend) into a single executable for **Windows**, **macOS**, and **Linux**.

## ✨ Key Features

- 🖥️ **Native Desktop Experience** - True native application with system integration
- 🚀 **Zero Configuration** - No server setup required, everything bundled
- 🔄 **Hot Reload Development** - Live reload during development
- 📱 **Cross-Platform** - Single codebase for all major operating systems  
- 🗄️ **Embedded Database** - SQLite database stored in user directory
- 📝 **Comprehensive Logging** - Detailed file-based logging for debugging
- 🔧 **Launch Method Independent** - Works from Applications folder, Finder, or terminal

## 🎯 What This Application Does

BeoEcho Desktop transforms the web-based BeoEcho mock API service into a standalone desktop application that:

- **Creates Mock APIs** - Define custom API endpoints with configurable responses
- **Request Forwarding** - Forward requests to actual backend services when needed  
- **Response Templating** - Dynamic responses with templates and variables
- **Request Logging** - Comprehensive logging of all requests and responses
- **Multi-User Support** - User management with workspace isolation
- **Dark/Light Mode** - Fully responsive UI with theming support

## 📋 Prerequisites

### Required Software
- **Go 1.21+** - [Download](https://golang.org/dl/)
- **Node.js 18+** - [Download](https://nodejs.org/)
- **Wails CLI v2.8.0+** - [Installation Guide](https://wails.io/docs/gettingstarted/installation)

### Platform-Specific Requirements

#### macOS
- **macOS 10.15+** (Catalina or newer)
- **Xcode Command Line Tools** - `xcode-select --install`
- For distribution: **Apple Developer Account** (for code signing)

#### Windows  
- **Windows 10+** (64-bit)
- **WebView2 Runtime** - Usually pre-installed on Windows 11
- **GCC Compiler** - via [TDM-GCC](https://jmeubank.github.io/tdm-gcc/) or [MinGW-w64](https://www.mingw-w64.org/)

#### Linux
- **GTK 3.20+** development libraries
- **WebKitGTK** development libraries
```bash
# Ubuntu/Debian
sudo apt-get install build-essential pkg-config libgtk-3-dev libwebkit2gtk-4.0-dev

# Fedora/RHEL
sudo dnf install gtk3-devel webkit2gtk3-devel

# Arch Linux  
sudo pacman -S gtk3 webkit2gtk
```

### Install Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Verify Installation
```bash
wails doctor
```

## 🚀 Quick Start

### 1. Initial Setup
```bash
# Clone the repository (if not already done)
git clone <repository-url>
cd mockoon-control-panel

# Navigate to desktop directory
cd desktop

# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd ../frontend && npm install && cd ../desktop
```

### 2. Development Mode (Recommended for Development)
```bash
# Start development server with hot reload
wails dev
```

**Development Mode Features:**
- 🔄 **Live Reload** - Frontend and backend changes reload automatically
- 🐛 **Debug Console** - Browser DevTools available for frontend debugging
- 📡 **API Testing** - Backend server accessible at `http://localhost:3600`
- 🔍 **Real-time Logs** - Console logs from both frontend and backend

### 3. Production Build
```bash
# Build for current platform
wails build

# Build with specific options
wails build -clean -upx -s  # Clean, compress, strip symbols
```

### 4. Platform-Specific Builds
```bash
# macOS (Intel)
wails build -platform darwin/amd64

# macOS (Apple Silicon) 
wails build -platform darwin/arm64

# Windows (64-bit)
wails build -platform windows/amd64

# Linux (64-bit)
wails build -platform linux/amd64

# Build for all platforms
wails build -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64
```

## 📁 Project Structure

```
desktop/                                    # Desktop application root
├── main.go                                # Desktop application entry point & Wails integration
├── go.mod                                 # Go module configuration
├── go.sum                                 # Go dependencies lock file
├── wails.json                            # Wails configuration (build settings, platforms)
├── app.json                              # Application metadata (name, version, description)
├── README.md                             # This documentation file
├── build/                                # Build output directory (generated)
│   ├── bin/                             # Compiled executable binaries
│   │   ├── BeoEcho Desktop.app/        # macOS application bundle
│   │   ├── BeoEcho.exe                 # Windows executable
│   │   └── BeoEcho                     # Linux executable  
│   ├── configs/                        # Build-time configuration files
│   └── darwin/                         # Platform-specific build artifacts
├── frontend/                           # Generated frontend assets (embedded)
│   ├── index.html                      # Main frontend entry point
│   ├── demo.html                       # Demo page
│   ├── login.html                      # Authentication page
│   ├── home.html                       # Dashboard page
│   ├── robots.txt                      # Web crawler instructions
│   ├── favicon.png                     # Application icon
│   ├── _app/                          # Compiled Svelte application
│   │   ├── immutable/                 # Hashed static assets
│   │   └── version.json               # Build version info
│   ├── css/                           # Compiled stylesheets
│   ├── fonts/                         # Web fonts and typography
│   └── images/                        # Static image assets
└── uploads/                           # User upload storage (runtime)
```

### Runtime Directory Structure (Desktop Mode)

When running as a desktop app, all data is stored in the user's home directory:

```
~/.beoecho/                            # Application data directory
├── configs/                           # Configuration files
│   ├── db/                           # Database storage
│   │   └── db.sqlite                 # Main SQLite database
│   └── caddy/                        # Reverse proxy configuration
├── uploads/                          # User uploaded files
└── logs/                             # Application logs
    └── desktop-YYYY-MM-DD.log        # Daily log files with detailed debugging
```

## ⚙️ Configuration

### Wails Configuration (`wails.json`)

```json
{
  "name": "BeoEcho Desktop",
  "outputfilename": "BeoEcho",
  "frontend:install": "npm install",
  "frontend:build": "npm run build:desktop",
  "frontend:dev:build": "npm run dev",
  "frontend:dev:install": "npm install",
  "frontend:dev:watcher": "npm run dev",
  "frontend:dev:serverUrl": "auto",
  "build": {
    "platforms": ["darwin/amd64", "darwin/arm64", "windows/amd64", "linux/amd64"]
  }
}
```

### Application Metadata (`app.json`)

```json
{
  "name": "BeoEcho Desktop", 
  "description": "Native desktop application for BeoEcho mock API service",
  "version": "1.0.0",
  "author": "BeoEcho Team"
}
```

### Desktop-Specific Features

#### 🖥️ Native System Integration
- **Menu Bar Integration** - Native application menus for each platform
- **System Tray Support** - Background operation with system tray icon
- **Auto-Launch** - Optional startup with system boot

#### 🔧 Cross-Platform Path Handling  
- **User Directory Storage** - All data stored in `~/.beoecho/` 
- **Launch Method Independence** - Works from any location (Applications, Finder, Terminal)
- **Automatic Directory Creation** - Required folders created on first run
- **Permission Handling** - Proper file system permissions across platforms

#### 📊 Enhanced Logging & Debugging
- **File-Based Logging** - Detailed logs saved to `~/.beoecho/logs/`
- **Multi-Writer Logging** - Console + file output simultaneously  
- **Startup Diagnostics** - Working directory, environment variables, path resolution
- **Error Recovery** - Graceful handling of startup failures

## 🛠️ Development Commands

### Frontend Development (SvelteKit)
```bash
# Build frontend for desktop embedding
cd ../frontend && npm run build:desktop

# Development mode with desktop integration
VITE_DESKTOP_MODE=true npm run dev

# Run frontend tests
npm test

# Type checking
npm run check
```

### Backend Development (Go)
```bash
# Run backend tests
cd ../backend && go test ./...

# Run backend directly (standalone mode)
go run main.go

# Run backend with specific port
go run main.go -p 8080

# Generate mocks for testing
go generate ./...
```

### Desktop Development (Wails)
```bash
# Development mode with hot reload (recommended)
wails dev

# Generate frontend bindings
wails generate module

# Build application (production)
wails build

# Build with optimizations
wails build -clean -upx -s

# Clean build artifacts
rm -rf build/ && wails build
```

### Cross-Platform Building
```bash
# Check available platforms
wails build -help

# Build for specific platform
wails build -platform darwin/arm64

# Build universal macOS binary
wails build -platform darwin/amd64,darwin/arm64

# Build for distribution
wails build -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64 -clean -upx
```

## 🧪 Testing & Quality Assurance

### Automated Testing
```bash
# Run all backend tests
cd ../backend && go test ./... -v

# Run tests with coverage
cd ../backend && go test ./... -cover

# Run frontend tests
cd ../frontend && npm test

# Type checking
cd ../frontend && npm run check
```

### Manual Testing Scenarios

#### 🖱️ Launch Method Testing
```bash
# Test 1: Launch from Applications folder (macOS)
open ~/Applications/"BeoEcho Desktop.app"

# Test 2: Launch from terminal
cd ~/Applications && open "BeoEcho Desktop.app"

# Test 3: Launch from build directory
./build/bin/"BeoEcho Desktop.app"/Contents/MacOS/BeoEcho

# Verify: Check logs show correct paths
tail -f ~/.beoecho/logs/desktop-$(date +%Y-%m-%d).log
```

#### 🌐 API Functionality Testing
```bash
# Wait for app startup (5 seconds), then test APIs
sleep 5

# Test health endpoint
curl -s http://localhost:3600/api/health | jq .

# Test authentication (should return auth error)
curl -s http://localhost:3600/api/api/workspaces

# Check database creation
ls -la ~/.beoecho/configs/db/
```

#### 🔧 Development Workflow Testing
```bash
# Test development mode
wails dev

# Make a change to frontend, verify hot reload
# Make a change to backend, verify restart and reload
```

### Platform-Specific Testing

#### macOS Testing
```bash
# Test both Intel and Apple Silicon builds
file ./build/bin/"BeoEcho Desktop.app"/Contents/MacOS/BeoEcho

# Test code signing (for distribution)
codesign -v ./build/bin/"BeoEcho Desktop.app"

# Test notarization requirements
spctl --assess --verbose ./build/bin/"BeoEcho Desktop.app"
```

#### Windows Testing  
```bash
# Test executable
./build/bin/BeoEcho.exe

# Check dependencies
ldd ./build/bin/BeoEcho.exe  # or equivalent Windows tool
```

#### Linux Testing
```bash
# Test executable
./build/bin/BeoEcho

# Check GTK dependencies
ldd ./build/bin/BeoEcho | grep gtk
```

## 📦 Distribution & Deployment

### Build for Distribution
```bash
# Production build with all optimizations
wails build \
  -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64 \
  -clean \
  -upx \
  -s \
  -ldflags "-X main.Version=$(git describe --tags --always)"
```

### Platform-Specific Distribution

#### 🍎 macOS Distribution
```bash
# Build universal binary
wails build -platform darwin/amd64,darwin/arm64 -clean

# Code signing (requires Apple Developer account)
codesign --deep --force --verify --verbose --sign "Developer ID Application: Your Name" ./build/bin/"BeoEcho Desktop.app"

# Create DMG installer
hdiutil create -volname "BeoEcho Desktop" -srcfolder ./build/bin/"BeoEcho Desktop.app" -ov -format UDZO BeoEcho-Desktop.dmg

# Notarization (for Gatekeeper)
xcrun notarytool submit BeoEcho-Desktop.dmg --apple-id your@email.com --password app-specific-password --wait
```

#### 🪟 Windows Distribution
```bash
# Build Windows executable
wails build -platform windows/amd64 -clean -upx

# Create installer with NSIS or Inno Setup
# Package with WebView2 bootstrapper for older Windows versions
```

#### 🐧 Linux Distribution
```bash
# Build Linux binary  
wails build -platform linux/amd64 -clean -upx

# Create AppImage
# Create .deb package (Ubuntu/Debian)
# Create .rpm package (Fedora/RHEL)
# Create Flatpak (universal Linux distribution)
```

## 🔧 Troubleshooting

### Common Development Issues

#### ❌ "no index.html found" Error
```bash
# Solution: Build frontend first
cd ../frontend && npm run build:desktop
cd ../desktop && wails build
```

#### ❌ Backend Dependencies Issues  
```bash
# Solution: Update and clean Go modules
go clean -modcache
go mod tidy
go mod download
```

#### ❌ Frontend Dependencies Issues
```bash
# Solution: Clean and reinstall Node modules
cd ../frontend
rm -rf node_modules package-lock.json .svelte-kit
npm install
npm run build:desktop
```

#### ❌ Wails CLI Issues
```bash
# Solution: Update Wails and check system
go install github.com/wailsapp/wails/v2/cmd/wails@latest
wails doctor
```

### Build & Runtime Issues

#### 🔄 Clean Build Environment
```bash
# Complete clean rebuild
rm -rf build/
rm -rf ../frontend/build/
rm -rf ../frontend/.svelte-kit/
rm -rf ../frontend/node_modules/

cd ../frontend && npm install && npm run build:desktop
cd ../desktop && wails build
```

#### 🗄️ Database Connection Issues
```bash
# Check database path and permissions
ls -la ~/.beoecho/configs/db/
sqlite3 ~/.beoecho/configs/db/db.sqlite ".tables"

# Check logs for database errors
tail -f ~/.beoecho/logs/desktop-$(date +%Y-%m-%d).log | grep -i database
```

#### 🌐 Network & Port Issues
```bash
# Check if port 3600 is available
lsof -i :3600

# Test API accessibility
curl -v http://localhost:3600/api/health

# Check firewall settings (macOS)
sudo pfctl -sr | grep 3600
```

### Platform-Specific Issues

#### macOS Issues
```bash
# Permission issues
sudo xattr -dr com.apple.quarantine ./build/bin/"BeoEcho Desktop.app"

# Code signing issues
codesign --remove-signature ./build/bin/"BeoEcho Desktop.app"

# Gatekeeper issues
spctl --master-disable  # Temporarily disable (not recommended for production)
```

#### Windows Issues  
```bash
# WebView2 missing
# Download and install WebView2 Runtime from Microsoft

# Antivirus false positives
# Add build directory to antivirus exclusions

# DLL missing errors
# Ensure all dependencies are included in build
```

#### Linux Issues
```bash
# GTK/WebKit missing
sudo apt-get install libgtk-3-0 libwebkit2gtk-4.0-37

# Permission denied
chmod +x ./build/bin/BeoEcho

# Library path issues
export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/usr/local/lib
```

## 📋 Development Workflow

### Daily Development Routine
1. **Start Development Environment**
   ```bash
   cd desktop && wails dev
   ```

2. **Make Changes** 
   - Frontend: Edit files in `../frontend/src/`
   - Backend: Edit files in `../backend/src/`
   - Desktop: Edit `main.go` for desktop-specific features

3. **Test Changes**
   - Changes auto-reload in development mode
   - Use browser DevTools for frontend debugging
   - Check terminal for backend logs

4. **Run Tests**
   ```bash
   cd ../backend && go test ./...
   cd ../frontend && npm test
   ```

### Release Preparation Workflow
1. **Code Quality Checks**
   ```bash
   # Backend tests
   cd ../backend && go test ./... -cover
   
   # Frontend type checking  
   cd ../frontend && npm run check
   
   # Linting
   cd ../frontend && npm run lint
   ```

2. **Build Frontend**
   ```bash
   cd ../frontend && npm run build:desktop
   ```

3. **Build Desktop Application**
   ```bash
   cd ../desktop && wails build -platform darwin/amd64,darwin/arm64,windows/amd64,linux/amd64 -clean
   ```

4. **Test Builds**
   ```bash
   # Test each platform build
   ./build/bin/"BeoEcho Desktop.app"/Contents/MacOS/BeoEcho  # macOS
   ./build/bin/BeoEcho.exe                                   # Windows (via Wine/VM)
   ./build/bin/BeoEcho                                       # Linux (via Docker/VM)
   ```

5. **Distribution Preparation**
   - Code signing (macOS/Windows)
   - Installer creation
   - Documentation updates
   - Release notes preparation

### Debugging Workflow
1. **Check Application Logs**
   ```bash
   tail -f ~/.beoecho/logs/desktop-$(date +%Y-%m-%d).log
   ```

2. **API Debugging**
   ```bash
   # Test API endpoints
   curl -v http://localhost:3600/api/health
   
   # Check backend logs
   grep -i error ~/.beoecho/logs/desktop-$(date +%Y-%m-%d).log
   ```

3. **Frontend Debugging**
   - Use browser DevTools in development mode
   - Check Svelte component state
   - Monitor network requests

4. **Database Debugging**
   ```bash
   # Open database directly
   sqlite3 ~/.beoecho/configs/db/db.sqlite
   .tables
   .schema users
   SELECT * FROM users LIMIT 5;
   ```

## 🔗 Related Documentation

### Core Documentation
- **[Wails Documentation](https://wails.io/docs/)** - Official Wails framework docs
- **[Backend API Documentation](../backend/README.md)** - Go backend service details
- **[Frontend Documentation](../frontend/README.md)** - Svelte frontend setup and development  
- **[Main Project Documentation](../README.md)** - Overall project setup and architecture

### Advanced Guides
- **[Desktop Development Guide](../docs/Desktop_Application_Guide.md)** - Detailed desktop-specific development
- **[API Authentication Guide](../docs/API%20Authentication%20and%20Authorization.md)** - Authentication system details
- **[Multi-User & Workspaces](../docs/Multi%20User%20and%20workspace.md)** - User management and workspaces
- **[Docker & Deployment](../docs/Setup%20Docker%20Traefik%20Cloudflare%20Tunnel.md)** - Server deployment options

### API References  
- **[Mock Rules API](../docs/Mock_Rules_API.md)** - Mock API configuration
- **[Replay API](../docs/Replay_API_Documentation.md)** - Request replay functionality
- **[Refactoring Plan](../docs/REFACTORING_PLAN.md)** - Architecture evolution plans

## 🆘 Support & Community

### Getting Help
- **Desktop Application Issues** - Check this README and [Wails documentation](https://wails.io/docs/)
- **Backend/API Issues** - See [Backend README](../backend/README.md)  
- **Frontend/UI Issues** - See [Frontend README](../frontend/README.md)
- **General Setup** - See [Main Project README](../README.md)

### Reporting Issues
When reporting issues, please include:
1. **Platform Information** - OS version, architecture  
2. **Version Information** - Wails version, Go version, Node.js version
3. **Steps to Reproduce** - Exact commands and actions taken
4. **Log Files** - Contents of `~/.beoecho/logs/desktop-YYYY-MM-DD.log`
5. **Expected vs Actual Behavior** - What should happen vs what actually happens

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes with tests
4. Test across platforms if possible
5. Submit a pull request with detailed description

---

**Built with ❤️ using [Wails v2](https://wails.io/), [Go](https://golang.org/), and [Svelte 5](https://svelte.dev/)**
