import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

export default defineConfig({
	plugins: [
		sveltekit()
	],
	server: {
		allowedHosts: ['local.yogasw.my.id', 'localhost']
	}
});
