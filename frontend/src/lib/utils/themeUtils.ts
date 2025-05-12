import { browser } from '$app/environment';
import { writable } from 'svelte/store';

/**
 * Theme utility functions that help with applying the right theme classes
 */

// Apply primary background class based on theme
export function themeBgPrimary(additionalClasses = '') {
  return `bg-white dark:bg-gray-800 ${additionalClasses}`;
}

// Apply secondary background class based on theme
export function themeBgSecondary(additionalClasses = '') {
  return `bg-gray-100 dark:bg-gray-700 ${additionalClasses}`;
}

// Apply tertiary background class based on theme
export function themeBgTertiary(additionalClasses = '') {
  return `bg-gray-50 dark:bg-gray-900 ${additionalClasses}`;
}

// Apply accent background class based on theme
export function themeBgAccent(additionalClasses = '') {
  return `bg-gray-200 dark:bg-gray-600 ${additionalClasses}`;
}

// Apply primary text color based on theme
export function themeTextPrimary(additionalClasses = '') {
  return `text-gray-800 dark:text-white ${additionalClasses}`;
}

// Apply secondary text color based on theme
export function themeTextSecondary(additionalClasses = '') {
  return `text-gray-600 dark:text-gray-300 ${additionalClasses}`;
}

// Apply muted text color based on theme
export function themeTextMuted(additionalClasses = '') {
  return `text-gray-500 dark:text-gray-400 ${additionalClasses}`;
}

// Apply border color based on theme
export function themeBorder(additionalClasses = '') {
  return `border-gray-200 dark:border-gray-700 ${additionalClasses}`;
}

// Apply light/semi-transparent border color for subtle separators
export function themeBorderLight(additionalClasses = '') {
  return `border-gray-200/70 dark:border-gray-700/70 ${additionalClasses}`;
}

// Apply very subtle/thin semi-transparent border
export function themeBorderSubtleLight(additionalClasses = '') {
  return `border-gray-200/30 dark:border-gray-700/30 ${additionalClasses}`;
}

// Apply subtle border color based on theme (less prominent)
export function themeBorderSubtle(additionalClasses = '') {
  return `border-gray-100 dark:border-gray-800 ${additionalClasses}`;
}

// Apply hover effect based on theme
export function themeHover(additionalClasses = '') {
  return `hover:bg-gray-200 dark:hover:bg-gray-600 ${additionalClasses}`;
}

// Apply shadow based on theme
export function themeShadow(additionalClasses = '') {
  return `shadow-md dark:shadow-gray-900/50 ${additionalClasses}`;
}

// Full component style helpers
export function inputField(additionalClasses = '') {
  return `block w-full p-3 ps-10 text-sm rounded-lg bg-gray-100 dark:bg-gray-800 
    border border-gray-200 dark:border-gray-700 text-gray-800 dark:text-white 
    focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400 ${additionalClasses}`;
}

export function primaryButton(additionalClasses = '') {
  return `bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
    text-sm flex items-center ${additionalClasses}`;
}

export function secondaryButton(additionalClasses = '') {
  return `bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-white hover:bg-gray-300 
    dark:hover:bg-gray-600 py-2 px-4 rounded flex items-center ${additionalClasses}`;
}

export function destructiveButton(additionalClasses = '') {
  return `bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded 
    flex items-center ${additionalClasses}`;
}

export function utilityButton(additionalClasses = '') {
  return `text-xs bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 
    hover:bg-gray-300 dark:hover:bg-gray-600 px-2 py-1 rounded ${additionalClasses}`;
}

export function card(additionalClasses = '') {
  return `bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 
    rounded-md shadow-md overflow-hidden ${additionalClasses}`;
}

export function cardHeader(additionalClasses = '') {
  return `flex justify-between items-center p-3 bg-gray-100 dark:bg-gray-750 ${additionalClasses}`;
}

export function methodBadge(method: string) {
  switch(method.toUpperCase()) {
    case 'GET': 
      return 'bg-green-600 text-white';
    case 'POST':
      return 'bg-blue-600 text-white';
    case 'PUT':
      return 'bg-yellow-600 text-white';
    case 'DELETE':
      return 'bg-red-600 text-white';
    default:
      return 'bg-gray-600 text-white';
  }
}

export function statusCodeBadge(statusCode: number) {
  if (statusCode >= 200 && statusCode < 300) {
    return 'bg-green-600 text-white';
  } else if (statusCode >= 300 && statusCode < 400) {
    return 'bg-blue-600 text-white';
  } else if (statusCode >= 400 && statusCode < 500) {
    return 'bg-yellow-600 text-white';
  } else if (statusCode >= 500) {
    return 'bg-red-600 text-white';
  } else {
    return 'bg-gray-600 text-white';
  }
}
