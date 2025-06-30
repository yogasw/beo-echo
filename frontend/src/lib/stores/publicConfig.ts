import { writable } from 'svelte/store';
import { getPublicConfig, type PublicConfigResponse } from '$lib/api/BeoApi';

// Create a writable store for public configuration
export const publicConfig = writable<PublicConfigResponse | null>(null);

// Flag to track if config has been loaded
export const isPublicConfigLoaded = writable<boolean>(false);

// Function to load public config (only loads once)
export async function loadPublicConfig(): Promise<PublicConfigResponse | null> {
	let isLoaded = false;
	
	// Check if already loaded
	isPublicConfigLoaded.subscribe((loaded) => {
		isLoaded = loaded;
	})();
	
	if (isLoaded) {
		// Already loaded, return current config
		let currentConfig: PublicConfigResponse | null = null;
		publicConfig.subscribe((config) => {
			currentConfig = config;
		})();
		return currentConfig;
	}
	
	// Not loaded yet, load from API
	try {
		const config = await getPublicConfig();
		publicConfig.set(config);
		isPublicConfigLoaded.set(true);
		return config;
	} catch (error) {
		console.error('Failed to load public config:', error);
		// Set default values on error
		const defaultConfig: PublicConfigResponse = {
			is_authenticated: false,
			landing_enabled: true,
			mock_url_format: 'subdomain'
		};
		publicConfig.set(defaultConfig);
		isPublicConfigLoaded.set(true);
		return defaultConfig;
	}
}
