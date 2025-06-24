<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated, auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import LandingPageHeader from '$lib/components/landing-page/LandingPageHeader.svelte';
	import LandingPageFooter from '$lib/components/landing-page/LandingPageFooter.svelte';
	import LandingContent from '$lib/components/landing-page/LandingContent.svelte';

	onMount(async () => {
		// Initialize auth
		await auth.initialize();
	});

	async function handleLogout() {
		await auth.logout();
		// Refresh the page to update the landing content
		window.location.reload();
	}

	async function handleBackToLanding() {
		await goto('/');
	}
</script>

<!-- Landing Page Layout -->
<div class="min-h-screen bg-white dark:bg-gray-900 transition-colors duration-200">
	<!-- Header -->
	<LandingPageHeader 
		showNavigation={true} 
		showBackButton={false} 
		showUserMenu={true}
		on:logout={handleLogout}
		on:back={handleBackToLanding}
	/>

	<!-- Main Content -->
	<LandingContent />

	<!-- Footer -->
	<LandingPageFooter />
</div>


