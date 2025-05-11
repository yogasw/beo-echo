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
        }
      }
    },
  },
  plugins: [],
} 