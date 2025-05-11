# Mockoon Control Panel Project Guide

## Project Overview
This project is a control panel for Mockoon API mocking service with a Golang backend and Svelte frontend.

## Project Structure
```
/mockoon-control-panel
├── backend/           # Golang Backend (BE)
│   └── ...            # Go files for the backend service
└── frontend/          # Svelte Frontend (FE)
    └── ...            # Svelte and TypeScript files
```

## Technology Stack

### Backend (BE)
- **Language**: Go
- **Framework**: Gin
- **ORM**: GORM
- **Database**: SQLite (configurable)
- **Features**:
  - RESTful API endpoints
  - Mock API management
  - Request logging
  - API authentication

### Frontend (FE)
- **Framework**: SvelteKit
- **Styling**: Tailwind CSS
- **Language**: TypeScript
- **State Management**: Svelte stores
- **Features**:
  - Responsive UI with mobile support
  - Dark/Light mode theming (default: dark)
  - Interactive mock API configuration
  - Request logging visualization

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

### Theme Toggle Component
- Create a theme toggle component that switches between modes
- Store theme preference in local storage
- Apply theme changes immediately without page refresh

### CSS Variables for Theming
- Consider using CSS variables for dynamic theme changes:
  ```css
  :root {
    /* Light theme defaults */
    --bg-primary: #ffffff;
    --bg-secondary: #f3f4f6;
    --text-primary: #1f2937;
    --text-secondary: #4b5563;
  }
  
  .dark {
    --bg-primary: #1f2937;
    --bg-secondary: #374151;
    --text-primary: #f9fafb;
    --text-secondary: #d1d5db;
  }
  ```

### Class-Based Toggle
- Use a class-based approach with Tailwind:
  ```html
  <div class="bg-white dark:bg-gray-800 text-gray-800 dark:text-white">
    <!-- Content -->
  </div>
  ```

## Development Guidelines

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

## Development Workflow

### Getting Started
1. Clone the repository
2. Install dependencies:
   ```
   # Backend
   cd backend
   go mod download
   
   # Frontend
   cd frontend
   npm install
   ```
3. Run development servers:
   ```
   # Backend
   cd backend
   go run main.go
   
   # Frontend
   cd frontend
   npm run dev
   ```

### Development Best Practices
- Follow the established UI patterns from the LogsTab and ConfigurationTab
- Test on mobile devices or with responsive device emulation
- Ensure both light and dark modes look polished
- Maintain consistent spacing and component sizing

## Component Library Reference
Based on the existing LogsTab and ConfigurationTab components, follow these established patterns:

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