// This file ensures the page loads correctly with query parameters in a static build
const LANDING_MODE = process.env.VITE_LANDING_MODE === 'true';

export const prerender = true;
export const ssr = LANDING_MODE; // Enable SSR for landing mode

// Make the page statically generate but allow for client-side behavior with URL parameters
export function load() {
  return {
    landingMode: LANDING_MODE
  };
}
