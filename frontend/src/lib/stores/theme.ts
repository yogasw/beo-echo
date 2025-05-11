import { browser } from '$app/environment';
import { writable } from 'svelte/store';

type Theme = 'dark' | 'light';

// Check if we're in the browser and get stored theme or default to 'dark'
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
