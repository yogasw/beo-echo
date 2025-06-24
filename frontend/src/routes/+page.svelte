<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { isFirstOpenPage } from '$lib/stores/isFirstOpen';
	import LandingPage from '$lib/components/landing-page/LandingPage.svelte';

	// For authenticated users, we can either show the landing page or redirect
	// This allows both authenticated and unauthenticated users to see the landing page
	let showLandingPage = true;

	onMount(async () => {
		console.log('onMount: page - landing page');
		
		// Handle first open behavior - optionally redirect authenticated users
		if ($isFirstOpenPage) {
			isFirstOpenPage.set(false);
			// Uncomment the line below if you want to redirect authenticated users to /home automatically
			// if ($isAuthenticated) {
			//     await goto('/home');
			// }
		}
	});
</script>

{#if showLandingPage}
	<LandingPage />
{/if}
