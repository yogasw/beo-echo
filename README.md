
# Beo Echo

> ⚠️ **AI-Generated Code Notice**: A significant portion of this codebase was generated using GitHub Copilot AI agent. While thoroughly tested, users should be aware that any bugs or issues arising from AI-generated code are not the responsibility of the original creator. Please report issues and contribute fixes as needed.

Beo Echo is a powerful API mocking service with a robust Golang backend and sleek Svelte frontend. It serves as a comprehensive control panel for creating, managing, and monitoring mock APIs and request forwarding, similar to popular tools like Beeceptor and Mockoon.

## 🌟 Project Overview

Beo Echo provides a user-friendly interface for:
- Creating and configuring mock APIs
- Forwarding requests to existing endpoints
- Real-time request logging and inspection
- API behavior management with multiple response types

## 🚀 How to Run

### Using Docker (Recommended)

Run Beo Echo using the pre-built Docker image:

```bash
docker run -d \
  -p 8080:80 \
  -v $(pwd)/beo-echo-config/:/app/configs/ \
  ghcr.io/yogasw/beo-echo:latest
```

> **Note**: The Docker image is built for x86 architecture. Running on ARM-based systems may result in errors as ARM builds are not currently provided due to extended build times.

After running the container, access the application at [http://localhost:8080](http://localhost:8080)

### Local Development

For development purposes, you can run the backend and frontend separately:

## 🏗️ Architecture

The project follows a clean architecture with:

```
/beo-echo/
├── .github/           # GitHub Actions and workflows
├── .vscode/           # VSCode settings
├── backend/           # Golang Backend
│   └── ...            # Go files for the backend service
└── frontend/          # Svelte Frontend
    └── ...            # Svelte and TypeScript files
```

### Backend Stack
- **Language**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: SQLite (configurable)

### Frontend Stack
- **Framework**: SvelteKit
- **Styling**: Tailwind CSS
- **Language**: TypeScript
- **State Management**: Svelte stores

---

## 🚀 Key Features

### API Request Handling
Beo Echo supports multiple execution modes:

- **Mock Mode**: Create and serve mock API responses
- **Forwarder Mode**: Forward requests to existing endpoints
- **Hybrid Mode**: Conditionally mock or forward based on rules

### Advanced Request Logging
Track and analyze all incoming requests with detailed logs including:
- Request method, path, headers, and body
- Response status, time, and content
- Execution mode indicators

### User Authentication
Secure access to your mock APIs with:
- JWT-based authentication
- Role-based access control
- Workspace isolation

### Theming Support
Choose your preferred visual experience:
- Dark mode (default)
- Light mode
- Consistent design across all components

---

## 🖼️ UI Screenshots

### 🔐 Login Page
![Login Page](./images/login.png)

### 📊 Logs Tab
![Logs Tab](./images/logs.png)

### 🛣️ Routes Tab
![Routes Tab](./images/routes.png)

---

## 📦 Getting Started

1. Clone the repository:
```bash
git clone https://github.com/yogasw/beo-echo.git
cd beo-echo
```

2. Start the backend:
```bash
cd backend
go run main.go
```

3. Start the frontend:
```bash
cd frontend
npm install
npm run dev
```

4. Access the application at [http://localhost:5173](http://localhost:5173)

### Using Caddy Server

Beo Echo now uses Caddy for simplified deployment instead of Traefik:

```bash
caddy run
```

The Caddyfile in the project root handles all the necessary routing.

---

## 🗄️ Database Configuration

By default, Beo Echo uses a SQLite database stored in the configured data directory. You can override the database connection by setting the `DATABASE_URL` environment variable to use PostgreSQL instead:

### Example using PostgreSQL:

```bash
DATABASE_URL=postgresql://username:password@host:port/database go run main.go
```

When a valid PostgreSQL connection string is provided in the `DATABASE_URL` environment variable, the system automatically uses PostgreSQL instead of SQLite.

---

## ⚠️ Disclaimer

- This is an active development project
- Feel free to fork and adapt it to suit your needs
- If you encounter issues or bugs, contributions are welcome

---

## 🎨 UI/UX Design

### Color Palette

#### Dark Mode (Default)
- **Background**: Dark grayish blue (`bg-gray-800`) with lighter accents (`bg-gray-700`, `bg-gray-750`)
- **Text**: Clean white (`text-white`) with muted variants (`text-gray-300`, `text-gray-400`)
- **Accent Colors**: Blue highlights (`text-blue-400`, `bg-blue-600`), status indicators in green, red, and yellow

#### Light Mode
- **Background**: Clean white (`bg-white`) with subtle gray accents (`bg-gray-50`, `bg-gray-100`)
- **Text**: Dark gray (`text-gray-800`) with lighter variants
- **Accent Colors**: Consistent with dark mode for brand recognition

### Component Styling

- **Cards & Panels**: Rounded containers with subtle shadows and borders
- **Form Elements**: Consistent styling with clear focus states
- **Buttons**: Color-coded by action type (primary, secondary, destructive)
- **Status Indicators**: Color-coded method badges (GET: green, POST: blue, etc.)

### Responsive Design
- Responsive grid layouts that adapt to different screen sizes
- Mobile-optimized touch targets and navigation
- Properly stacked elements on smaller screens

## 🌓 Theme Implementation

Beo Echo uses a dual-approach theme system:

1. **Tailwind Dark Mode**: Class-based strategy for toggling themes
2. **Theme Utility Functions**: Consistent application of theme styles

### Theme Toggle Example

```svelte
<ThemeToggle showLabel={true} size="default" />
```

This provides an elegant interface for users to switch between dark and light modes, with the selected preference persisted across sessions.

---

---

## 👨‍💻 Development Guidelines

### Code Organization
- Keep components modular and focused on a single responsibility
- Use Svelte stores for shared state
- Follow TypeScript best practices with proper typing

### API Integration
- All API calls should be centralized in the `$lib/api` directory
- Handle loading states and error conditions consistently

### Accessibility
- Ensure keyboard navigation works properly
- Use proper ARIA attributes for interactive elements
- Maintain sufficient color contrast ratios
- Test with screen readers

## 🧩 Component Library Reference

Based on the existing components, follow these established patterns:

### Page Structure
1. Main container with `w-full bg-gray-800 p-4 relative`
2. Header section with title, subtitle, and action buttons
3. Content area with appropriate spacing between sections
4. Optional notification system for feedback

### Common Elements
- **Section Headers**: Include an icon in a colored background with title
- **Expandable Sections**: Toggle visibility with chevron indicators
- **Search Fields**: Include icon prefix and clear, descriptive placeholder text
- **Status Indicators**: Use color-coding for status representation
- **Data Displays**: Use structured layouts with proper labeling

### Mobile Responsiveness
- Ensure touch targets are sufficiently large (min 44px)
- Stack elements vertically on small screens
- Use responsive grids: `grid-cols-1 md:grid-cols-2`
- Test on various screen sizes
