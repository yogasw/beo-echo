<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { isFirstOpenPage } from '$lib/stores/isFirstOpen';
	import LandingPage from '$lib/components/landing-page/LandingPage.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import { publicConfig, loadPublicConfig } from '$lib/stores/publicConfig';
	import { toast } from '$lib/stores/toast';

	// State management
	let isLoading = true;
	let error: Error | null = null;
	let showLandingPage = false; // Default to false until config is loaded

	// Load public configuration
	async function loadConfig() {
		try {
			isLoading = true;
			error = null;
			
			// Load config using store (will only hit API once)
			const config = await loadPublicConfig();
			
			// Handle redirection based on configuration and authentication
			if (config?.landing_enabled) {
				showLandingPage = true;
			} else {
				// Landing page disabled - redirect based on auth status
				if ($isAuthenticated) {
					await goto('/home');
				} else {
					await goto('/login');
				}
			}
		} catch (err) {
			console.error('Failed to load public config:', err);
			error = err as Error;
			toast.error(err);
			// Default fallback - show landing page on error
			showLandingPage = true;
		} finally {
			isLoading = false;
		}
	}

	onMount(async () => {
		console.log('onMount: page - loading public configuration');
		
		// Load public configuration first
		await loadConfig();
		
		// Handle first open behavior
		if ($isFirstOpenPage) {
			isFirstOpenPage.set(false);
		}
	});
</script>

{#if isLoading}
	<!-- Show loading skeleton while fetching configuration -->
	<div class="min-h-screen flex items-center justify-center theme-bg-primary">
		<div class="text-center">
			<SkeletonLoader type="card" count={1} />
			<p class="mt-4 theme-text-secondary">Loading configuration...</p>
		</div>
	</div>
{:else if error}
	<!-- Show error with retry option -->
	<div class="min-h-screen flex items-center justify-center theme-bg-primary">
		<ErrorDisplay 
			message={error.message} 
			type="error" 
			retryable={true}
			onRetry={loadConfig}
		/>
	</div>
{:else if showLandingPage}
	<!-- Show landing page if enabled -->
	<LandingPage />
{/if}
