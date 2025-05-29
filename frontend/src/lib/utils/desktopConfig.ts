// Desktop mode detection and API configuration
import { browser } from '$app/environment';

// Check if running in desktop mode (Wails)
export const isDesktopMode = () => {
  return browser && typeof window !== 'undefined' && window.wails !== undefined;
};

// Get appropriate API base URL based on environment
export const getApiBaseUrl = () => {
  // In desktop mode, backend server runs on embedded port
  if (isDesktopMode()) {
    // Desktop mode - backend will be embedded and accessible via localhost
    // Wails will handle the backend server startup
    return 'http://localhost:3600/api';
  }
  
  // Web mode - use environment variable or default
  return import.meta.env.VITE_API_BASE_URL || 'http://localhost:3600/api';
};

// Desktop-specific configuration
export const desktopConfig = {
  // Window configuration for desktop app
  windowTitle: 'BeoEcho Desktop',
  minWidth: 800,
  minHeight: 600,
  defaultWidth: 1200,
  defaultHeight: 800,
  
  // App behavior configuration
  hideWindowOnClose: false,
  autoStartBackend: true,
  
  // Desktop-specific features
  enableSystemTray: true,
  enableAutoUpdate: false, // Future feature
  enableOfflineMode: true,
  
  // Data storage paths (handled by backend)
  userDataPath: '~/.beoecho',
  configPath: '~/.beoecho/configs',
  logPath: '~/.beoecho/logs'
};

// Export for use in components
export default {
  isDesktopMode,
  getApiBaseUrl,
  desktopConfig
};
