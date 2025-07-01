/**
 * Version utility for Beo Echo
 * Provides version information from environment variables or fallback
 */

/**
 * Get the current application version
 * During build time, this will be injected from the VERSION file
 * In development, it falls back to a default value
 */
export function getVersion(): string {
  // This will be replaced by Vite during build time
  // @ts-ignore - Vite will define this during build
  return import.meta.env.VITE_APP_VERSION || '2.3.2-dev';
}

/**
 * Get version with prefix
 */
export function getVersionWithPrefix(prefix: string = 'v'): string {
  return `${prefix}${getVersion()}`;
}

/**
 * Check if this is a development build
 */
export function isDevelopmentVersion(): boolean {
  return getVersion().includes('-dev');
}

/**
 * Get version info object with additional metadata
 */
export function getVersionInfo() {
  const version = getVersion();
  const isDev = isDevelopmentVersion();
  
  return {
    version,
    fullVersion: getVersionWithPrefix('v'),
    isDevelopment: isDev,
    buildTime: import.meta.env.VITE_BUILD_TIME || new Date().toISOString(),
  };
}
