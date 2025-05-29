import adapter from '@sveltejs/adapter-static';

import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

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
			pages: process.env.WAILS_BUILD ? '../desktop/assets' : 'build',
			assets: process.env.WAILS_BUILD ? '../desktop/assets' : 'build',
			fallback: 'index.html', // Add fallback for SPA behavior
			precompress: false,
			strict: false
		})
	}
	
};

export default config;


