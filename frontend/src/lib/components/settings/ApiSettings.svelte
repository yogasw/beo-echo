<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	
	// API key state
	let apiKey = 'sk_test_51NxXBCLKYtpKtDnyBwvfEnGX8kTGvSbZy2S54PXX';
	let isCopied = false;
	let isRegenerating = false;
	
	// Copy API key
	function copyApiKey() {
		if (navigator.clipboard) {
			navigator.clipboard.writeText(apiKey);
			isCopied = true;
			
			// Reset copy status after 3 seconds
			setTimeout(() => {
				isCopied = false;
			}, 3000);
		}
	}
	
	// Regenerate API key
	async function regenerateApiKey() {
		isRegenerating = true;
		
		try {
			// Simulate API call
			await new Promise(r => setTimeout(r, 800));
			
			// Mock new API key
			apiKey = 'sk_test_' + Math.random().toString(36).substring(2, 15) + 
				Math.random().toString(36).substring(2, 15);
			
			toast.success('API key regenerated successfully');
		} catch (error) {
			toast.error('Failed to regenerate API key');
		} finally {
			isRegenerating = false;
		}
	}
</script>

<div class="space-y-6">
	<div class="theme-bg-primary border theme-border rounded-lg p-4">
		<h3 class="font-medium theme-text-primary mb-2">Your API Key</h3>
		<p class="text-sm theme-text-secondary mb-4">
			This key provides programmatic access to the Beo Echo API. Keep it secret!
		</p>
		
		<div class="flex items-stretch">
			<div class="flex-grow p-3 font-mono text-sm theme-bg-secondary rounded-l-md border-t border-b border-l theme-border theme-text-secondary overflow-hidden whitespace-nowrap overflow-ellipsis">
				{apiKey}
			</div>
			<button 
				on:click={copyApiKey}
				class="p-3 text-sm theme-bg-secondary theme-text-primary hover:bg-gray-200 dark:hover:bg-gray-600 border-t border-b theme-border flex items-center"
			>
				{#if isCopied}
					<i class="fas fa-check text-green-500 mr-1"></i>
					<span>Copied</span>
				{:else}
					<i class="fas fa-copy mr-1"></i>
					<span>Copy</span>
				{/if}
			</button>
			<button 
				on:click={regenerateApiKey}
				class="p-3 text-sm bg-red-500 hover:bg-red-600 text-white rounded-r-md border border-red-500 flex items-center"
				disabled={isRegenerating}
			>
				{#if isRegenerating}
					<i class="fas fa-spinner fa-spin mr-1"></i>
				{:else}
					<i class="fas fa-sync-alt mr-1"></i>
				{/if}
				<span>Regenerate</span>
			</button>
		</div>
		
		<p class="mt-3 text-xs theme-text-muted">
			<i class="fas fa-exclamation-triangle text-yellow-500 mr-1"></i>
			Regenerating your API key will revoke the old one. Any applications using this key will need to be updated.
		</p>
	</div>
	
	<div class="bg-blue-500/10 dark:bg-blue-600/10 p-4 rounded-lg border border-blue-200 dark:border-blue-900">
		<h3 class="flex items-center theme-text-primary font-medium mb-2">
			<i class="fas fa-info-circle text-blue-500 mr-2"></i>
			<span>API Key Usage</span>
		</h3>
		<p class="text-sm theme-text-secondary mb-3">
			Use this API key in your requests by including it in the header:
		</p>
		<div class="font-mono text-sm p-3 theme-bg-primary rounded border theme-border overflow-x-auto">
			<code class="theme-text-primary">Authorization: Bearer {apiKey}</code>
		</div>
	</div>
</div>
