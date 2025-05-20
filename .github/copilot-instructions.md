# Beo Echo Project Guide

## Project Overview
This project is a Beo Echo API mocking service with a Golang backend and Svelte frontend. It includes features for creating mock APIs, forwarding requests, and managing API behaviors, similar to tools like Beeceptor and Mockoon.

## Project Structure
```
/beo-echo/
├── .github/           # GitHub Actions and workflows
├── .vscode/           # VSCode settings
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
      after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
      after:bg-white after:rounded-full after:h-5 after:w-5 
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