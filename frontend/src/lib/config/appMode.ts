// Build-time app modes. These are replaced by Vite at build time and DO NOT
// change in the browser (unlike the per-page LANDING_MODE hydration trick).
//
//   STATIC_MODE  → pure static site with NO backend (GitHub Pages deploy).
//                  Backend-only features (login, cloud, create project,
//                  dashboard, MCP token generation) are hidden.
//   LANDING_MODE → prerendered landing that may still hydrate into the full
//                  app when a backend is present.
export const STATIC_MODE = import.meta.env.VITE_STATIC_MODE === 'true';
export const LANDING_MODE = import.meta.env.VITE_LANDING_MODE === 'true';
