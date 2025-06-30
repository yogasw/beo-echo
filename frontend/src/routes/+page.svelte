<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { isFirstOpenPage } from '$lib/stores/isFirstOpen';
	import LandingPage from '$lib/components/landing-page/LandingPage.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import { getPublicConfig, type PublicConfigResponse } from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';

	// State management
	let publicConfig: PublicConfigResponse | null = null;
	let isLoading = true;
	let error: Error | null = null;
	let showLandingPage = false; // Default to false until config is loaded

	// Load public configuration
	async function loadPublicConfig() {
		try {
			isLoading = true;
			error = null;
			publicConfig = await getPublicConfig();
			
			// Handle redirection based on configuration and authentication
			if (publicConfig.landing_enabled) {
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
		await loadPublicConfig();
		
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
			onRetry={loadPublicConfig}
		/>
	</div>
{:else if showLandingPage}
	<!-- Show landing page if enabled -->
	<LandingPage />
{/if}
