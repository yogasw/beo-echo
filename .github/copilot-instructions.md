# Beo Echo Project Guide

## Project Overview
This project is a Beo Echo API mocking service with a Golang backend and Svelte frontend. It includes features for creating mock APIs, forwarding requests, and managing API behaviors, similar to tools like Beeceptor and Mockoon.

### Key Features
- **Mock API Creation**: Define custom API endpoints with configurable responses
- **Request Forwarding**: Forward requests to actual backend services when needed
- **Response Templating**: Create dynamic responses with templates and variables
- **Request Logging**: Comprehensive logging of all requests and responses
- **User Management**: Multi-user support with workspace isolation
- **Dark/Light Mode**: Fully responsive UI with themeing support

## Project Structure
```
/beo-echo/
├── .github/           # GitHub Actions and workflows
├── .vscode/           # VSCode settings
├── backend/           # Golang Backend (BE)
├── docs/              # Documentation files
└── frontend/          # Svelte Frontend (FE)
```

## Technology Stack

### Backend (BE)
- **Language**: Go (1.21+)
- **Framework**: Gin Web Framework
- **ORM**: GORM
- **Database**: SQLite (configurable to other databases)
- **Authentication**: JWT + OAuth (Google)
- **Features**:
  - RESTful API endpoints
  - Mock API management
  - Request logging and filtering
  - Multi-user authentication and authorization
  - Workspace isolation
  - Real-time proxy forwarding

### Frontend (FE)
- **Framework**: SvelteKit
- **Styling**: Tailwind CSS with custom theming
- **Language**: TypeScript
- **State Management**: Svelte stores
- **API Client**: Custom fetch wrappers with type safety
- **Features**:
  - Responsive UI with mobile support
  - Dark/Light mode theming (default: dark)
  - Interactive mock API configuration
  - Request logging visualization with filtering
  - Live request monitoring
  - Request/response body formatting

## UI/UX Design Guidelines

### Color Palette

#### Dark Mode (Default)
- **Background**: 
  - Primary: `bg-gray-800` - Dark grayish blue
  - Secondary: `bg-gray-700` - Lighter dark blue
  - Tertiary accents: `bg-gray-750`, `bg-gray-850`, `bg-gray-900/50` - Various shades for depth
  - Borders: `border-gray-700` - Subtle separators

- **Text**:
  - Primary: `text-white` - Clean white text
  - Secondary: `text-gray-300`, `text-gray-400` - Muted text
  - Links/Accents: `text-blue-400`, `text-blue-500` - Blue highlights

- **Buttons & Interactive Elements**:
  - Primary actions: `bg-blue-600` → `hover:bg-blue-700` - Deep blue with hover state
  - Destructive actions: `bg-red-600` → `hover:bg-red-700` - Deep red
  - Success indicators: `bg-green-600`, `text-green-400` - Green elements
  - Warning indicators: `bg-yellow-600`, `text-yellow-400` - Yellow highlights

#### Light Mode
- **Background**:
  - Primary: `bg-white` or `bg-gray-50`
  - Secondary: `bg-gray-100`
  - Tertiary: `bg-gray-200`
  
- **Text**:
  - Primary: `text-gray-800`
  - Secondary: `text-gray-600`
  - Muted: `text-gray-500`

- **Accent Colors**: Keep consistent with dark mode for brand recognition

### Typography
- **Font Families**:
  - UI: System font stack, sans-serif
  - Code/Monospace: `font-mono` for JSON, request data
  
- **Font Sizes**:
  - Headers: `text-xl` (with `font-bold`)
  - Subheaders: `text-sm` (with `font-semibold`)
  - Body text: `text-sm`
  - Small text: `text-xs`

### Component Styling

#### Containers
- Use rounded corners: `rounded-md`, `rounded-lg`
- Apply subtle shadows: `shadow-md`
- Use border separation: `border border-gray-700`
- For overlay modals: `fixed inset-0 bg-black bg-opacity-50`

#### Cards & Panels
- Container: `bg-gray-800 border border-gray-700 rounded-md shadow-md overflow-hidden`
- Headers: `flex justify-between items-center p-3 bg-gray-750`
- Content areas: `p-4` with appropriate spacing between elements
- Use `transition:fade` for smooth section expansion/collapse

#### Form Elements
- Inputs:
  ```
  block w-full p-3 ps-10 text-sm rounded-lg bg-gray-800 
  border border-gray-700 text-white focus:ring-blue-500 
  focus:border-blue-500 placeholder-gray-400
  ```
- Include iconography in fields with `absolute inset-y-0 start-0 flex items-center ps-3`
- Toggles:
  ```
  w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
  peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
  rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
  after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
  after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
  after:transition-all peer-checked:bg-blue-600
  ```

#### Buttons
- **Primary Action**:
  ```
  bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
  text-sm flex items-center
  ```

- **Secondary Action**:
  ```
  bg-gray-700 hover:bg-gray-600 text-white py-2 px-4 rounded 
  flex items-center
  ```

- **Destructive Action**:
  ```
  bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded 
  flex items-center
  ```

- **Utility/Small**:
  ```
  text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 px-2 py-1 rounded
  ```

- **Badge/Pill**:
  ```
  px-3 py-1 rounded-full bg-gray-900/50 text-xs font-medium
  ```

#### Status Indicators
- Method badges with color coding:
  - GET: `bg-green-600 text-white`
  - POST: `bg-blue-600 text-white`
  - PUT: `bg-yellow-600 text-white`
  - DELETE: `bg-red-600 text-white`
  - Default: `bg-gray-600 text-white`

- Status code indicators follow HTTP conventions:
  - 2xx: `bg-green-600 text-white`
  - 3xx: `bg-blue-600 text-white`
  - 4xx: `bg-yellow-600 text-white`
  - 5xx: `bg-red-600 text-white`

- Live status:
  ```
  <span class="relative flex h-3 w-3 mr-2">
    <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
    <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
  </span>
  ```

#### Notifications
- Toast format:
  ```
  fixed top-6 right-6 bg-gray-700 text-white px-4 py-2 rounded 
  shadow-lg z-50 flex items-center
  ```
- Include appropriate icons: `fas fa-check-circle text-green-400` or `fas fa-exclamation-circle text-red-400`
- Use Svelte transitions for smooth appearance: `transition:fade={{ duration: 200 }}`

#### Icons
- Use Font Awesome icons consistently
- Common uses:
  - Section headers: `fas fa-list-alt`, `fas fa-cogs`, `fas fa-info-circle`
  - Actions: `fas fa-save`, `fas fa-trash-alt`, `fas fa-sync`, `fas fa-copy`
  - Status: `fas fa-check-circle`, `fas fa-exclamation-triangle`, `fas fa-exclamation-circle`
  - Navigation: `fas fa-chevron-up`, `fas fa-chevron-down`

### Responsive Design

#### Mobile Considerations
- Use `grid grid-cols-1 md:grid-cols-2 gap-4` for responsive grids
- Ensure adequate touch targets for mobile (min 44px)
- Stack elements vertically on smaller screens
- Ensure text remains readable at all screen sizes

#### Layout Patterns
- Use Flexbox for alignment: `flex justify-between items-center`
- Space elements consistently: `space-y-4`, `space-x-3`
- For scrollable areas: `overflow-auto max-h-64`

### Interaction Patterns
- Expandable sections with toggle functionality
- Tabbed interfaces for complex data (request/response)
- Form validation with clear feedback
- Toast notifications for actions
- Modal confirmations for destructive actions

## Light/Dark Mode Implementation

### Theme Architecture Overview

The project uses a dual-approach theme system:

1. **Tailwind Dark Mode**: Using Tailwind's built-in dark mode with the class strategy
2. **Theme Utility Functions**: Wrapper utility functions for consistent theme application

### Theme Store Implementation

- Use the Svelte store for theme state management:

```typescript
// src/lib/stores/theme.ts
import { browser } from '$app/environment';
import { writable } from 'svelte/store';

type Theme = 'dark' | 'light';

// Get stored theme or default to 'dark'
const userTheme = browser && localStorage.getItem('theme');
const initialTheme = userTheme ? userTheme as Theme : 'dark';

// Create the theme store with initial value
export const theme = writable<Theme>(initialTheme);

// Update localStorage and document class whenever theme changes
theme.subscribe((value) => {
  if (browser) {
    localStorage.setItem('theme', value);
    
    // Update the document class for Tailwind dark mode
    if (value === 'dark') {
      document.documentElement.classList.add('dark');
    } else {
      document.documentElement.classList.remove('dark');
    }
  }
});

// Toggle theme function
export function toggleTheme() {
  theme.update(currentTheme => currentTheme === 'dark' ? 'light' : 'dark');
}
```

### Theme Utilities

Create comprehensive theme utility functions that provide consistent styling:

```typescript
// src/lib/utils/themeUtils.ts
// Apply theme-specific styles consistently across components

export function themeBgPrimary(additionalClasses = '') {
  return `bg-white dark:bg-gray-800 ${additionalClasses}`;
}

export function themeBgSecondary(additionalClasses = '') {
  return `bg-gray-100 dark:bg-gray-700 ${additionalClasses}`;
}

export function themeTextPrimary(additionalClasses = '') {
  return `text-gray-800 dark:text-white ${additionalClasses}`;
}

export function themeTextSecondary(additionalClasses = '') {
  return `text-gray-600 dark:text-gray-300 ${additionalClasses}`;
}

// Component-level utilities for consistent styling
export function card(additionalClasses = '') {
  return `bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 
    rounded-md shadow-md overflow-hidden ${additionalClasses}`;
}

export function primaryButton(additionalClasses = '') {
  return `bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
    text-sm flex items-center ${additionalClasses}`;
}
```

### Theme Toggle Component

Create a reusable theme toggle component:

```svelte
<!-- src/lib/components/ThemeToggle.svelte -->
<script lang="ts">
  import { theme, toggleTheme } from '$lib/stores/theme';
  
  // Optional: Customize toggle appearance
  export let showLabel = false;
  export let size = 'default'; // 'small', 'default', 'large'
  
  $: sizeClass = size === 'small' ? 'w-8 h-4' : 
                 size === 'large' ? 'w-14 h-7' : 'w-11 h-6';
</script>

<div class="flex items-center">
  {#if showLabel}
    <span class="mr-2 text-sm theme-text-secondary">
      {$theme === 'dark' ? 'Dark' : 'Light'}
    </span>
  {/if}
  
  <label class="inline-flex items-center cursor-pointer">
    <input 
      type="checkbox" 
      checked={$theme === 'dark'} 
      on:change={toggleTheme} 
      class="sr-only peer"
    />
    <div class="{sizeClass} bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
      rounded-full peer peer-focus:outline-none peer-focus:ring-2 
      peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 
      peer-checked:after:translate-x-full peer-checked:after:border-white 
      after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
      after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
      after:transition-all dark:border-gray-600"
    >
    </div>
    
    {#if showLabel}
      <i class="ml-2 {$theme === 'dark' ? 'fas fa-moon text-blue-400' : 'fas fa-sun text-yellow-400'}"></i>
    {/if}
  </label>
</div>
```

### Global Utility Classes in app.css

Add these utility classes to your app.css for easier theme class application:

```css
/* Global theme utility classes */
.theme-bg-primary {
  @apply bg-white dark:bg-gray-800;
}

.theme-bg-secondary {
  @apply bg-gray-100 dark:bg-gray-700;
}

.theme-bg-tertiary {
  @apply bg-gray-50 dark:bg-gray-900;
}

.theme-text-primary {
  @apply text-gray-800 dark:text-white;
}

.theme-text-secondary {
  @apply text-gray-600 dark:text-gray-300;
}

.theme-text-muted {
  @apply text-gray-500 dark:text-gray-400;
}

.theme-border {
  @apply border-gray-200 dark:border-gray-700;
}

.theme-border-subtle {
  @apply border-gray-100 dark:border-gray-800;
}
```

### Best Practices for Theme Implementation

1. **Consistent Theme Application**:
   - Always use theme utility classes or functions for elements that should adapt
   - Use `theme-text-primary` instead of direct colors like `text-gray-800`
   - Apply both light and dark variants for backgrounds: `bg-white dark:bg-gray-800`

2. **Component Theming**:
   - For complex components, use the ThemeUtils functions:
     ```svelte
     <div class={ThemeUtils.card()}>
       <button class={ThemeUtils.primaryButton()}>Action</button>
     </div>
     ```

3. **Dynamic UI Elements**:
   - For UI elements that need to dynamically change with theme:
     ```svelte
     <script>
       import { theme } from '$lib/stores/theme';
     </script>
     
     <div class="flex items-center">
       <i class="fas {$theme === 'dark' ? 'fa-moon text-blue-400' : 'fa-sun text-yellow-400'}"></i>
       <span>{$theme === 'dark' ? 'Dark Mode' : 'Light Mode'}</span>
     </div>
     ```

4. **Color Constants for Special Elements**:
   - Define special-purpose colors that maintain consistency across themes:
     ```typescript
     const successTextClass = (isDark = false) => 
       isDark ? 'text-green-400' : 'text-green-600';
     const warningBgClass = (isDark = false) => 
       isDark ? 'bg-yellow-400/20' : 'bg-yellow-100';
     ```

5. **Test Both Themes**:
   - Visually verify all components in both themes during development
   - Pay special attention to contrast ratios for accessibility
   - Check borders and separators which often need theme-specific adjustments

## Development Guidelines

### Code Organization
- Keep components modular and focused on a single responsibility
- Use Svelte stores for shared state
- Follow TypeScript best practices with proper typing

### API Integration
- All API calls should be centralized in the `$lib/api` directory
- Handle loading states and error conditions consistently


## Frontend Development Guidelines

The frontend is built with SvelteKit and follows a specific structure:

```
frontend/                  # JavaScript/TypeScript frontend with Svelte and Tailwind
├── src/                   # Source code directory
│   ├── lib/               # Library code (components, utils, etc.)
│   │   ├── api/           # API client code for backend communication
│   │   ├── components/    # Reusable UI components
│   │   │   ├── common/    # General UI components (buttons, cards, etc.)
│   │   │   ├── layout/    # Layout components (headers, footers, etc.)
│   │   │   └── specific/  # Feature-specific components
│   │   ├── images/        # Image assets used in components
│   │   ├── services/      # Frontend service layers
│   │   ├── stores/        # Svelte stores for state management
│   │   ├── styles/        # Global and shared styles
│   │   ├── types/         # TypeScript type definitions
│   │   └── utils/         # Utility functions and helpers
│   ├── routes/            # SvelteKit routes (pages)
│   │   ├── home/          # Home page routes
│   │   ├── login/         # Authentication routes
│   │   └── ...            # Other feature routes
│   ├── app.css            # Global CSS
│   └── app.html           # HTML template
├── static/                # Static assets served as-is
│   ├── favicon.png        # Site favicon
│   └── robots.txt         # Robots crawling instructions
└── build/                 # Compiled output (generated)
```

### Frontend Development Best Practices

1. **Component Organization**
   - Keep components small and focused on a single responsibility
   - Place shared components in the appropriate lib/components subdirectory
   - Use props validation for all component inputs
   - Document complex components with JSDoc comments

2. **API Integration**
   - All API calls should be centralized in the `$lib/api` directory
   - Use TypeScript interfaces matching backend models
   - Implement proper error handling with user-friendly messages
   - Use loading states for all async operations

3. **State Management**
   - Use Svelte stores for global state management
   - Document store responsibilities and usage
   - Follow the pattern in existing stores for consistent implementation
   - Consider component-local state for isolated features

4. **Performance Considerations**
   - Lazy-load routes and heavy components
   - Use proper Svelte lifecycle methods
   - Avoid unnecessary re-renders
   - Use efficient event handlers with proper cleanup

5. **Testing**
   - Write unit tests for critical components and utilities
   - Implement integration tests for complex page interactions
   - Ensure accessibility testing (a11y) is part of the process

### Key Frontend Technologies
- **SvelteKit**: For routing and server-side rendering
- **Tailwind CSS**: For styling with utility classes
- **TypeScript**: For type safety and better developer experience
- **Font Awesome**: For consistent iconography


----------------------------

# Simplified Backend Architecture Guide

## Dependency Inversion Principle (DIP)

The backend follows a clean architecture using the Dependency Inversion Principle:

> High-level modules should not depend on low-level modules. Both should depend on abstractions.
> Abstractions should not depend on details. Details should depend on abstractions.

## Layer Structure

```
┌────────────────┐
│    Handler     │  HTTP interface, uses Service interfaces
├────────────────┤
│    Service     │  Business logic, defines & uses Repository interfaces
├────────────────┤
│   Repository   │  Data access, implements interfaces defined by Service
└────────────────┘
```

## Directory Structure

```
backend/                   # Golang Backend
├── cmd/                   # Application entry points
│   ├── root.go            # Root command definitions
│   └── server.go          # Server startup command
├── logs/                  # Log output directory
├── src/                   # Core application code
│   ├── auth/              # Authentication module
│   │   ├── handler/       # HTTP handlers for auth routes
│   │   └── services/      # Authentication business logic
│   ├── database/          # Database operations
│   │   ├── models.go      # Central data model definitions
│   │   ├── db.go          # Database connection management
│   │   └── *_repo.go      # Repository implementations
│   ├── middlewares/       # HTTP middleware components
│   ├── mocks/             # Mock API management
│   │   ├── handler/       # HTTP handlers for mock endpoints
│   │   ├── services/      # Mock routing business logic
│   │   └── repositories/  # Data access layer
│   ├── logs/              # Request logging module
│   │   ├── handlers.go    # HTTP handlers for logs
│   │   └── services.go    # Log processing logic
│   ├── types/             # Common type definitions
│   └── utils/             # Utility functions
├── uploads/               # File upload storage
├── configs/               # Configuration files
└── Makefile               # Build commands and automation

## Data Model Reference

All database models are defined in `backend/src/database/models.go`. This file serves as the **central reference** for all data models in the application. When working with data in the backend:

1. **Always reference existing models**: Before creating new data-related functions, always check the models defined in `models.go` to understand the existing schema and relationships.

2. **Field documentation requirements**: When adding new fields to existing models or creating new models, each field must include proper documentation comments explaining:
   - Purpose of the field
   - Format requirements if applicable
   - Relationship to other models
   - Default values if applicable
   - Any validation requirements

3. **Key Data Models**:
   - `SystemConfig`: System-wide configuration settings
   - `Project`: API mock project definition with endpoints and settings
   - `MockEndpoint`: Specific API endpoint configurations 
   - `MockResponse`: Response configuration for endpoints
   - `RequestLog`: Detailed request/response logging
   - `User`: User account information
   - `Workspace`: Shared workspace for projects
   - `UserWorkspace`: User membership in workspaces

4. **Model reference example**:
   ```go
   // ProjectMode defines operation mode of mock system per project
   type ProjectMode string
   
   const (
     ModeMock      ProjectMode = "mock"      // Serves predefined mock responses only
     ModeProxy     ProjectMode = "proxy"     // Uses mocks when available, otherwise forwards requests
     ModeForwarder ProjectMode = "forwarder" // Always forwards all requests to target endpoint
     ModeDisabled  ProjectMode = "disabled"  // Endpoint inactive - no responses served
   )
   ```

5. **Field documentation example**:
   ```go
   type Project struct {
     ID            string         `gorm:"type:string;primaryKey" json:"id"`
     Name          string         `gorm:"type:string" json:"name"`
     WorkspaceID   string         `gorm:"type:string;index" json:"workspace_id"`       // Foreign key to the associated workspace
     Mode          ProjectMode    `gorm:"type:string;default:'mock'" json:"mode"`      // default: mock
     Status        string         `gorm:"type:string;default:'running'" json:"status"` // running, stopped, error
     ActiveProxyID *string        `gorm:"type:string" json:"active_proxy_id"`          // Current active proxy configuration
     Alias         string         `gorm:"type:string;uniqueIndex;not null" json:"alias"` // Subdomain or alias for the project
     URL           string         `json:"url"`                                           // URL for the project UI display
   }
   ```

## Module Organization

Each feature module follows a simplified structure with concrete implementations:

1. **Handler Layer** (`modules/handler.go`):
   - Processes HTTP requests and produces responses
   - Validates and transforms input data 
   - Uses the concrete service implementation
   - Converts domain errors to appropriate HTTP status codes
   - One file per module, example: `workspaces/handler.go`

2. **Service Layer** (`modules/service.go`):
   - Contains core business logic and domain rules
   - Defines repository interfaces that it needs
   - Implements concrete service implementation
   - Orchestrates operations across multiple repositories
   - One file per module, example: `workspaces/service.go`

3. **Repository Interfaces**:
   - Defined in the service file alongside the service implementation
   - Specify data access requirements for each service
   - Clean separation from database implementation details

4. **Repository Implementations** (`database/repositories/`):
   - Implement the interfaces defined in service files
   - Handle database operations and queries
   - Convert between domain and database models
   - All database operations are centralized here
   - Files should be named `*_repo.go`

## Key Principles

1. **Simplified Structure**
   - Concrete service implementations directly in the module directory
   - Repository interfaces defined in the same file as the service that uses them
   - Clear dependency flow from handler → service → repository

2. **Layer Responsibilities**
   - **Handlers**: Handle external requests and input validation
   - **Services**: Implement business logic and define data needs
   - **Repository Interfaces**: Define data access requirements (in service.go)
   - **Repository Implementations**: Implement database access (in database/)

## Implementation Guidelines

### Step 1: Define Service with Repository Interface

```go
// users/service.go
package users

import (
    "context"
    "beo-echo/backend/src/database"
)

// UserRepository defines data access requirements for user operations
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*database.User, error)
    Create(ctx context.Context, user *database.User) error
}

// UserService implements user business operations
type UserService struct {
    repo UserRepository
}

// NewUserService creates a new user service
func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(ctx context.Context, id string) (*database.User, error) {
    return s.repo.FindByID(ctx, id)
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, user *database.User) error {
    return s.repo.Create(ctx, user)
}
```

### Step 2: Create Handler

```go
// users/handler.go
package users

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
    service *UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(service *UserService) *UserHandler {
    return &UserHandler{service: service}
}

// GetUser handles GET /users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
    userID := c.Param("id")
    user, err := h.service.GetUser(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, user)
}

// CreateUser handles POST /users
func (h *UserHandler) CreateUser(c *gin.Context) {
    // Implementation details...
}
```

### Step 3: Implement Repository

```go
// database/repositories/user_repo.go
package repositories

import (
    "context"
    "gorm.io/gorm"
    
    "beo-echo/backend/src/database"
    "beo-echo/backend/src/users"  // Import the package, not a subpackage
)

// userRepository implements users.UserRepository
type userRepository struct {
    db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) users.UserRepository {
    return &userRepository{db: db}
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*database.User, error) {
    // Implementation details...
    return nil, nil
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *database.User) error {
    // Implementation details...
    return nil
}
```

### Step 4: Register Dependencies and Routes

All routes must be registered in the central `src/server.go` file in the `SetupRouter()` function. Do not create separate route registration functions in feature modules:

```go
// In src/server.go
func SetupRouter() *gin.Engine {
    // Create Gin router
    router := gin.Default()
    
    // Setup middleware
    // ...
    
    // API routes group
    apiGroup := router.Group("/api")
    apiGroup.Use(middlewares.JWTAuthMiddleware())
    {
        // Get database
        db := database.DB
        
        // Create repository
        userRepo := repositories.NewUserRepository(db)
        
        // Create service with repository
        userService := users.NewUserService(userRepo)
        
        // Create handler with service
        userHandler := users.NewUserHandler(userService)
        
        // Register routes directly in server.go
        apiGroup.GET("/users/:id", userHandler.GetUser)
        apiGroup.POST("/users", userHandler.CreateUser)
        
        // For complex feature modules, group them together
        usersGroup := apiGroup.Group("/users")
        {
            usersGroup.GET("/profile", userHandler.GetProfile)
            usersGroup.PUT("/settings", userHandler.UpdateSettings)
        }
    }
    
    return router
}
```

## Key Benefits

1. **Testability**: Mock repositories easily for unit testing
2. **Modularity**: Change implementations without affecting services
3. **Clean Architecture**: Clear separation of concerns
4. **Consistent Patterns**: Standard approach for all features

## Testing Requirements
- Every new feature must include unit tests
- Use testify for assertions and mocking
- Use mockery for generating mocks
- Test coverage should be maintained or improved
- Run all tests after any changes

### Test Structure
```go
func TestServiceMethod(t *testing.T) {
    // Given - Setup test dependencies and expectations
    mockRepo := mocks.NewMockRepository(t)
    mockRepo.EXPECT().Method(mock.Anything, "input").Return(expectedValue, nil)
    
    svc := NewService(mockRepo)
    
    // When - Call the method under test
    result, err := svc.Method(context.Background(), "input")
    
    // Then - Verify results and behaviors
    assert.NoError(t, err)
    assert.Equal(t, expectedValue, result)
    mockRepo.AssertExpectations(t) // Verify all expected calls were made
}
```

## Code Quality Requirements

### 1. Testing
- Run tests before committing: `go test ./...`
- Write tests for new features
- Update tests when modifying existing features
- Use table-driven tests where appropriate

### 2. Error Handling
- Use custom error types
- Return meaningful error messages
- Handle all error cases
- Log errors appropriately

### 3. Logging
- Use zerolog package for all logging operations 
- Always include context in every function
- Follow structured logging pattern
- Log levels must be appropriate for the message
   
Example:
```go
// Initialize logger in main or init
logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
ctx = logger.WithContext(ctx)

// In functions, always accept and use context
func (s *service) DoSomething(ctx context.Context, input Input) error {
    log := zerolog.Ctx(ctx)
    
    log.Info().
        Str("input_id", input.ID).
        Str("user", input.UserID).
        Msg("processing request")
        
    // ... function logic ...
    
    if err != nil {
        log.Error().
            Err(err).
            Str("input_id", input.ID).
            Msg("failed to process request")
        return err
    }
    
    return nil
}
```

Log Level Guidelines:
- `Trace()`: For very detailed debugging
- `Debug()`: For development debugging
- `Info()`: For tracking normal operations
- `Warn()`: For handled issues that might need attention
- `Error()`: For unhandled errors that need immediate attention
- `Fatal()`: For errors that prevent application startup

Context Guidelines:
- Every function must accept `context.Context` as first parameter
- Pass context through all layers (handler → service → repository)
- Add relevant request information to context (user ID, request ID, etc.)
- Use context timeouts for external operations

Security Guidelines:
- Never log sensitive information (passwords, tokens, personal data)
- Mask or truncate potentially sensitive IDs in logs
- Use appropriate log levels to avoid exposing detailed errors in production

### 4. Documentation
- Document all exported functions and types
- Include examples in documentation
- Keep README.md up to date
- Document configuration options

### 5. Code Style
- Follow Go best practices
- Use consistent naming
- Keep functions focused and small
- Use meaningful variable names

## Development Workflow

1. Create new feature module with proper structure
2. Implement interfaces and types
3. Generate mocks using mockery
4. Implement business logic with tests
5. Add HTTP handlers with tests
6. Run all tests before committing:
   ```bash
   cd backend && go test ./...
   ```

### Running Tests
The project includes predefined VS Code tasks to facilitate testing:

```bash
# Run all backend tests
cd backend && go test ./...

# Run specific package tests
cd backend && go test ./src/database

# With code coverage report
cd backend && go test -cover ./...
```

You can also use the VS Code task "Run all Tests in Backend" to execute all tests in the backend directory.
