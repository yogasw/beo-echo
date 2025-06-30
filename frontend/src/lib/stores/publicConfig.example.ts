// Example usage of publicConfig store in other components

import { publicConfig, loadPublicConfig } from '$lib/stores/publicConfig';

// In a component:
// 1. Import the store
// 2. Use $publicConfig to access the data reactively
// 3. Call loadPublicConfig() only when needed (it will only hit API once)

// Example in a Svelte component:
/*
<script>
	import { publicConfig, loadPublicConfig } from '$lib/stores/publicConfig';
	
	// The store will automatically update this variable when config changes
	$: config = $publicConfig;
	
	// Use the data
	$: urlFormat = config?.mock_url_format || 'subdomain';
	$: landingEnabled = config?.landing_enabled || false;
	
	// Load config only if needed (will not hit API if already loaded)
	async function ensureConfigLoaded() {
		await loadPublicConfig();
	}
</script>

<div>
	{#if $publicConfig}
		<p>URL Format: {$publicConfig.mock_url_format}</p>
		<p>Landing Enabled: {$publicConfig.landing_enabled}</p>
	{:else}
		<p>Loading config...</p>
	{/if}
</div>
*/
