import { browser } from '$app/environment';
import { writable } from 'svelte/store';

/**
 * Theme utility functions that help with applying the right theme classes
 * Following project design guidelines for consistent styling
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
  return `bg-gray-50 dark:bg-gray-900/50 ${additionalClasses}`;
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
  return `block w-full p-3 ps-10 text-sm rounded-lg bg-gray-800 
    border border-gray-700 text-white focus:ring-blue-500 
    focus:border-blue-500 placeholder-gray-400 ${additionalClasses}`;
}

export function methodBadge(method = 'GET', additionalClasses = '') {
  let baseClass = 'px-2 py-1 rounded-md text-xs font-medium';
  switch(method.toUpperCase()) {
    case 'GET':
      return `${baseClass} bg-green-600 text-white ${additionalClasses}`;
    case 'POST':
      return `${baseClass} bg-blue-600 text-white ${additionalClasses}`;
    case 'PUT':
      return `${baseClass} bg-yellow-600 text-white ${additionalClasses}`;
    case 'PATCH':
      return `${baseClass} bg-orange-600 text-white ${additionalClasses}`;
    case 'DELETE':
      return `${baseClass} bg-red-600 text-white ${additionalClasses}`;
    default:
      return `${baseClass} bg-gray-600 text-white ${additionalClasses}`;
  }
}

export function statusBadge(statusCode = 200, additionalClasses = '') {
  let baseClass = 'px-2 py-1 rounded-md text-xs font-medium';
  if (statusCode >= 200 && statusCode < 300) {
    return `${baseClass} bg-green-600 text-white ${additionalClasses}`;
  } else if (statusCode >= 300 && statusCode < 400) {
    return `${baseClass} bg-blue-600 text-white ${additionalClasses}`;
  } else if (statusCode >= 400 && statusCode < 500) {
    return `${baseClass} bg-yellow-600 text-white ${additionalClasses}`;
  } else if (statusCode >= 500) {
    return `${baseClass} bg-red-600 text-white ${additionalClasses}`;
  } else {
    return `${baseClass} bg-gray-600 text-white ${additionalClasses}`;
  }
}

export function primaryButton(additionalClasses = '') {
  return `bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
    text-sm flex items-center ${additionalClasses}`;
}

export function secondaryButton(additionalClasses = '') {
  return `bg-gray-700 hover:bg-gray-600 text-white py-2 px-4 rounded 
    flex items-center ${additionalClasses}`;
}

export function destructiveButton(additionalClasses = '') {
  return `bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded 
    flex items-center ${additionalClasses}`;
}

export function utilityButton(additionalClasses = '') {
  return `text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 px-2 py-1 rounded ${additionalClasses}`;
}

export function badge(type = 'default', additionalClasses = '') {
  let baseClass = 'px-3 py-1 rounded-full text-xs font-medium';
  switch(type) {
    case 'success':
      return `${baseClass} bg-green-600 text-white ${additionalClasses}`;
    case 'info':
      return `${baseClass} bg-blue-600 text-white ${additionalClasses}`;
    case 'warning':
      return `${baseClass} bg-yellow-600 text-white ${additionalClasses}`;
    case 'danger':
      return `${baseClass} bg-red-600 text-white ${additionalClasses}`;
    default:
      return `${baseClass} bg-gray-900/50 text-xs font-medium ${additionalClasses}`;
  }
}

export function card(additionalClasses = '') {
  return `bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-md overflow-hidden 
    rounded-md shadow-md overflow-hidden ${additionalClasses}`;
}

export function cardHeader(additionalClasses = '') {
  return `flex justify-between items-center p-3 bg-gray-100 dark:bg-gray-750 ${additionalClasses}`;
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

export function notificationToast(type = 'default', additionalClasses = '') {
  const baseClass = `fixed top-6 right-6 bg-gray-700 text-white px-4 py-2 rounded shadow-lg z-50 flex items-center`;
  
  switch(type) {
    case 'success':
      return `${baseClass} border-l-4 border-green-500 ${additionalClasses}`;
    case 'error':
      return `${baseClass} border-l-4 border-red-500 ${additionalClasses}`;
    case 'warning':
      return `${baseClass} border-l-4 border-yellow-500 ${additionalClasses}`;
    case 'info':
    default:
      return `${baseClass} border-l-4 border-blue-500 ${additionalClasses}`;
  }
}

export function headerSection(additionalClasses = '') {
  return `flex items-center p-3 ${additionalClasses}`;
}

export function contentSection(additionalClasses = '') {
  return `p-4 ${additionalClasses}`;
}

export function liveStatus(isActive = true, additionalClasses = '') {
  return `
    <span class="relative flex h-3 w-3 mr-2 ${additionalClasses}">
      ${isActive ? '<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>' : ''}
      <span class="relative inline-flex rounded-full h-3 w-3 ${isActive ? 'bg-green-500' : 'bg-gray-500'}"></span>
    </span>
  `;
}

export function projectStatusBadge(status: string, additionalClasses = '') {
  const baseClass = 'px-2 py-1 rounded-md text-xs font-medium flex items-center';
  
  switch(status) {
    case 'running':
      return `${baseClass} bg-green-600/20 text-green-400 dark:bg-green-900/30 dark:text-green-400 ${additionalClasses}`;
    case 'stopped':
      return `${baseClass} bg-gray-600/20 text-gray-400 dark:bg-gray-900/30 dark:text-gray-400 ${additionalClasses}`;
    case 'error':
      return `${baseClass} bg-red-600/20 text-red-400 dark:bg-red-900/30 dark:text-red-400 ${additionalClasses}`;
    default:
      return `${baseClass} bg-gray-600/20 text-gray-400 dark:bg-gray-900/30 dark:text-gray-400 ${additionalClasses}`;
  }
}

export function projectStatusIndicator(status: string, additionalClasses = '') {
  const baseClass = 'relative flex h-3 w-3 mr-2';
  
  switch(status) {
    case 'running':
      return `
        <span class="${baseClass} ${additionalClasses}">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
        </span>
      `;
    case 'stopped':
      return `
        <span class="${baseClass} ${additionalClasses}">
          <span class="relative inline-flex rounded-full h-3 w-3 bg-gray-500"></span>
        </span>
      `;
    case 'error':
      return `
        <span class="${baseClass} ${additionalClasses}">
          <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-red-400 opacity-75"></span>
          <span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
        </span>
      `;
    default:
      return `
        <span class="${baseClass} ${additionalClasses}">
          <span class="relative inline-flex rounded-full h-3 w-3 bg-gray-500"></span>
        </span>
      `;
  }
}
