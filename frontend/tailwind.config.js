/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        gray: {
          750: '#2d3748', // Custom shade for headers
          775: '#252f3e', // Additional shade for dark mode elements
          825: '#1e2633', // Additional shade for dark mode hover states
          850: '#1a202c', // Custom shade for depth
          875: '#161b27', // Additional deeper shade for dark mode
        },
        theme: {
          // These can be accessed as bg-theme-primary, text-theme-secondary, etc.
          primary: 'var(--bg-primary)',
          secondary: 'var(--bg-secondary)',
          tertiary: 'var(--bg-tertiary)',
          accent: 'var(--bg-accent)',
          text: {
            primary: 'var(--text-primary)',
            secondary: 'var(--text-secondary)',
            muted: 'var(--text-muted)'
          },
          border: 'var(--border-color)'
        }
      },
      backgroundColor: {
        theme: {
          primary: 'var(--bg-primary)',
          secondary: 'var(--bg-secondary)',
          tertiary: 'var(--bg-tertiary)',
          accent: 'var(--bg-accent)'
        }
      },
      textColor: {
        theme: {
          primary: 'var(--text-primary)',
          secondary: 'var(--text-secondary)',
          muted: 'var(--text-muted)'
        }
      },
      borderColor: {
        theme: {
          DEFAULT: 'var(--border-color)',
          light: 'var(--border-color-light)',
          dark: 'var(--border-color-dark)'
        }
      }
    },
  },
  plugins: [],
} 