// This file ensures the page loads correctly with query parameters in a static build
export const prerender = true;
export const ssr = false;

// Make the page statically generate but allow for client-side behavior with URL parameters
export function load() {
  return {};
}
