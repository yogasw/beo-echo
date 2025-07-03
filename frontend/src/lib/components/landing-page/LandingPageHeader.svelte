<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { isAuthenticated, auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	const dispatch = createEventDispatcher();

	// Props
	export let showUserMenu = true; // Show user menu when authenticated

	async function handleLogout() {
		await auth.logout();
		dispatch('logout');
	}

	async function goToDashboard() {
		await goto('/home');
	}

	function handleBackClick() {
		dispatch('back');
	}
</script>

<!-- Header/Navigation Bar -->
<header
	class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50"
>
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Left Side: Logo + Navigation Menu -->
			<div class="flex items-center space-x-8">
				<!-- Logo and Brand -->
				<div class="flex items-center">
					<div class="flex-shrink-0 flex items-center">
						<!-- Clickable Logo -->
						<a
							href="/"
							class="flex items-center hover:opacity-80 transition-opacity no-underline hover:no-underline"
							title="Go to home page"
							aria-label="Go to Beo Echo home page"
						>
							<img src="/favicon.svg" alt="Beo Echo" class="w-8 h-8 mr-3" />
							<h1 class="text-xl font-bold text-gray-800 dark:text-white">Beo Echo</h1>
						</a>
					</div>
				</div>

				<!-- Navigation Links -->
				<nav class="hidden md:flex items-center space-x-6">
					<a
						href="/#features"
						class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
						title="View features"
						aria-label="View features"
					>
						Features
					</a>
					<a
						href="/#modes"
						class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
						title="View operating modes"
						aria-label="View operating modes"
					>
						Modes
					</a>
					<a
						href="/#pricing"
						class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
						title="View pricing plans"
						aria-label="View pricing plans"
					>
						Pricing
					</a>
				</nav>
			</div>

			<!-- Right Side: Actions -->
			<div class="flex items-center space-x-4">
				<!-- Main Actions Group -->
				<div class="flex items-center space-x-3">
					<!-- GitHub Stars Badge -->
					<a
						href="https://github.com/yogasw/beo-echo"
						target="_blank"
						rel="noopener noreferrer"
						class="hidden sm:block"
						title="Star us on GitHub"
						aria-label="Star Beo Echo on GitHub"
					>
						<img
							src="https://img.shields.io/github/stars/yogasw/beo-echo?style=social"
							alt="GitHub stars"
							class="h-5"
						/>
					</a>

					{#if $isAuthenticated && showUserMenu}
						<!-- User Menu -->
						<div class="flex items-center space-x-2">
							<button
								on:click={goToDashboard}
								class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
								title="Move to dashboard"
								aria-label="Move to dashboard"
							>
								<i class="fas fa-tachometer-alt mr-2"></i>
								Move to Dashboard
							</button>
						</div>
					{:else}
						<!-- Login/Signup Buttons for non-authenticated users -->
						<div class="flex items-center space-x-2">
							<a
								href="/login"
								class="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-all duration-200"
								title="Sign in to your account"
								aria-label="Sign in to your account"
							>
								Sign in
							</a>
							<a
								href="/login"
								class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
								title="Get started with Beo Echo"
								aria-label="Get started with Beo Echo"
							>
								<i class="fas fa-rocket mr-2"></i>
								Get Started
							</a>
						</div>
					{/if}

					<!-- Mobile Menu Toggle -->
					<button
						class="md:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
						title="Open mobile menu"
						aria-label="Open mobile menu"
					>
						<i class="fas fa-bars text-gray-600 dark:text-gray-300"></i>
					</button>
				</div>

				<!-- Theme Toggle - Separated on the right -->
				<button
					type="button"
					on:click={toggleTheme}
					class="w-10 h-10 flex items-center justify-center rounded-full theme-bg-secondary hover:bg-opacity-80 transition-colors border theme-border"
					title={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
					aria-label={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
				>
					<i class="fas {$theme === 'dark' ? 'fa-sun' : 'fa-moon'} text-base theme-text-primary"
					></i>
				</button>
			</div>
		</div>
	</div>
</header>
