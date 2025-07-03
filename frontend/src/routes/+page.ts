// since there's no dynamic data here, we can prerender
// it so that it gets served as a static asset in production
const LANDING_MODE = process.env.VITE_LANDING_MODE === 'true';

export const prerender = true;
export const ssr = LANDING_MODE; // Enable SSR for landing mode

export function load() {
  return {
    landingMode: LANDING_MODE
  };
}
