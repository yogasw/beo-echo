// since there's no dynamic data here, we can prerender
// it so that it gets served as a static asset in production
const LANDING_MODE = import.meta.env.VITE_LANDING_MODE === 'true';
const STATIC_MODE = import.meta.env.VITE_STATIC_MODE === 'true';

export const prerender = true;
export const ssr = LANDING_MODE || STATIC_MODE; // SSR/SSG for landing & static builds

export function load() {
  return {
    landingMode: LANDING_MODE || STATIC_MODE
  };
}
