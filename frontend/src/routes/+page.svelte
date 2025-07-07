<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { isFirstOpenPage } from '$lib/stores/isFirstOpen';
	import LandingPage from '$lib/components/landing-page/LandingPage.svelte';
	import BeoEchoLoader from '$lib/components/common/BeoEchoLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import { publicConfig, loadPublicConfig } from '$lib/stores/publicConfig';
	import { toast } from '$lib/stores/toast';
	import { browser } from '$app/environment';

	// Check if we're in landing mode (build time environment variable)
	let LANDING_MODE = import.meta.env.VITE_LANDING_MODE === 'true';
	if (browser) {
		LANDING_MODE = false; // Ensure this is false in browser context
	}

	// State management
	let isLoading = !LANDING_MODE; // Don't show loading if in landing mode
	let error: Error | null = null;
	let showLandingPage = LANDING_MODE; // Default to true if in landing mode

	// Load public configuration
	async function loadConfig() {
		try {
			isLoading = true;
			error = null;

			// If in landing mode, always show landing page
			if (LANDING_MODE) {
				showLandingPage = true;
				return;
			}

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

		// If in landing mode, skip API calls and show landing immediately
		if (LANDING_MODE) {
			isLoading = false;
			return;
		}

		// Load public configuration first
		await loadConfig();

		// Handle first open behavior
		if ($isFirstOpenPage) {
			isFirstOpenPage.set(false);
		}
	});
</script>

<!-- Landing page always rendered if in landing mode (for SSG) -->
{#if LANDING_MODE}
	<LandingPage />
{:else}
	<!-- Dynamic content for non-landing mode -->
	{#if isLoading}
		<!-- Show beautiful Beo Echo loader while fetching configuration -->
		<div class="min-h-screen flex items-center justify-center theme-bg-primary">
			<BeoEchoLoader 
				message="Loading configuration..." 
				size="lg"
				animated={true}
			/>
		</div>
	{:else if error}
		<!-- Show error with retry option -->
		<div class="min-h-screen flex items-center justify-center theme-bg-primary">
			<ErrorDisplay message={error.message} type="error" retryable={true} onRetry={loadConfig} />
		</div>
	{:else if showLandingPage}
		<!-- Show landing page if enabled -->
		<LandingPage />
	{/if}
{/if}
