import adapter from '@sveltejs/adapter-static';

import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

// Check build modes
const isDesktopMode = process.env.VITE_DESKTOP_MODE === 'true';
const isLandingMode = process.env.VITE_LANDING_MODE === 'true';
const isStaticMode = process.env.VITE_STATIC_MODE === 'true';
// Both landing and static builds are prerendered SSG with no SPA fallback.
const isSSG = isLandingMode || isStaticMode;
// Base path for GitHub Pages project sites (e.g. "/beo-echo"). Empty for root
// (custom domain / user or org page). Set via the BASE_PATH env in CI.
const basePath = process.env.BASE_PATH ?? '';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	// compilerOptions: { // Add this section
	// 	runes: true
	// },

	kit: {
		// adapter-auto only supports some environments, see https://svelte.dev/docs/kit/adapter-auto for a list.
		// If your environment is not supported, or you settled on a specific environment, switch out the adapter.
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		paths: {
			base: basePath,
			// SSG on GitHub Pages: use absolute, base-prefixed asset URLs so they
			// resolve correctly whether the page is served at /beo-echo or
			// /beo-echo/ (relative './asset' breaks on the no-trailing-slash URL).
			// Desktop keeps relative paths for file:// loading.
			relative: isDesktopMode
		},

		adapter: adapter({
			// default options are shown. On some platforms
			// these options are set automatically — see below
			pages: isDesktopMode ? '../desktop/frontend' : 'build',
			assets: isDesktopMode ? '../desktop/frontend' : 'build',
			fallback: isSSG ? (isStaticMode ? '404.html' : null) : 'index.html', // static: 404 fallback; landing SSG: none; SPA: index.html
			precompress: false,
			strict: false
		}),

		// Configure prerendering for landing page SSG
		prerender: {
			entries: isSSG ? ['/', '/guide', '/login', '*'] : [], // Prerender landing + guide (+ login for BE landing) for SSG builds
			handleHttpError: 'warn',
			handleMissingId: 'warn',
			handleEntryGeneratorMismatch: 'warn',
			handleUnseenRoutes: ({ id }) => {
				// Ignore acceptable unlinked routes
				if (['/demo', '/demo/badges', '/demo/components', '/demo/toggle-switch', '/home', '/login'].includes(id)) {
					return;
				}
				// Default behavior for everything else
				console.warn(`Unseen route: ${id}`);
			}
		}
	}
	
};

export default config;