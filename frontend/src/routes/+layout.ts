// Force client-side routing for the entire app
export const ssr = false;
export const prerender = true;

// This ensures the app works as an SPA in production
export function load() {
  return {};
}
