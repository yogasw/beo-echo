import adapter from '@sveltejs/adapter-static';

import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

// Check build modes
const isDesktopMode = process.env.VITE_DESKTOP_MODE === 'true';
const isLandingMode = process.env.VITE_LANDING_MODE === 'true';

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
		adapter: adapter({
			// default options are shown. On some platforms
			// these options are set automatically — see below
			pages: isDesktopMode ? '../desktop/frontend' : 'build',
			assets: isDesktopMode ? '../desktop/frontend' : 'build',
			fallback: isLandingMode ? null : 'index.html', // No fallback for SSG landing, fallback for SPA
			precompress: false,
			strict: false
		}),
		
		// Configure prerendering for landing page SSG
		prerender: {
			entries: isLandingMode ? ['/', '/login', '*'] : [], // Always prerender root and login in landing mode
			handleHttpError: 'warn',
			handleMissingId: 'warn',
			handleEntryGeneratorMismatch: 'warn'
		}
	}
	
};

export default config;