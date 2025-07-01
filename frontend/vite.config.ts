import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';
import { readFileSync, existsSync } from 'fs';
import { resolve, dirname } from 'path';
import { fileURLToPath } from 'url';

// Get current directory for ES modules
const __filename = fileURLToPath(import.meta.url);
const __dirname = dirname(__filename);

// Read version from VERSION file
function getVersion(): string {
	try {
		// In Docker build, VERSION is copied to /app/VERSION
		// In local development, VERSION is in parent directory
		const possiblePaths = [
			resolve(__dirname, '..', 'VERSION'),        // Local development
			resolve('/app', 'VERSION'),                 // Docker build
			resolve(__dirname, '..', '..', 'VERSION')   // Alternative path
		];
		
		for (const versionPath of possiblePaths) {
			if (existsSync(versionPath)) {
				try {
					const version = readFileSync(versionPath, 'utf-8').trim();
					console.log(`✅ Version loaded from ${versionPath}: ${version}`);
					return version;
				} catch (error) {
					console.warn(`Failed to read ${versionPath}:`, error);
					continue;
				}
			}
		}
		
		throw new Error('VERSION file not found in any expected location');
	} catch (error) {
		console.warn('Could not read VERSION file, using fallback version:', error);
		return '2.3.2-dev';
	}
}

export default defineConfig({
	plugins: [
		sveltekit()
	],
	server: {
		allowedHosts: ['local.yogasw.my.id', 'localhost']
	},
	define: {
		// Inject version and build time during build
		'import.meta.env.VITE_APP_VERSION': JSON.stringify(getVersion()),
		'import.meta.env.VITE_BUILD_TIME': JSON.stringify(new Date().toISOString()),
	}
});
