/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  darkMode: 'class',
  theme: {
    extend: {
      colors: {
        gray: {
          750: '#2d3748', // Custom shade for headers
          850: '#1a202c',  // Custom shade for depth
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
        theme: 'var(--border-color)'
      }
    },
  },
  plugins: [],
} 