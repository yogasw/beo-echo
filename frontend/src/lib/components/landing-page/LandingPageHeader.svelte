<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { isAuthenticated, auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	const dispatch = createEventDispatcher();

	// Props
	export let showNavigation = true; // Show navigation links (features, modes, pricing)
	export let showBackButton = false; // Show back button (for login page)
	export let showUserMenu = true; // Show user menu when authenticated

	async function handleLogin() {
		await goto('/login');
	}

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
<header class="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-50">
	<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
		<div class="flex items-center justify-between h-16">
			<!-- Left Side: Logo + Navigation Menu -->
			<div class="flex items-center space-x-8">
				<!-- Back Button (for login page) -->
				{#if showBackButton}
					<button
						on:click={handleBackClick}
						class="p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
						title="Go back to landing page"
						aria-label="Go back to landing page"
					>
						<i class="fas fa-arrow-left text-gray-600 dark:text-gray-300"></i>
					</button>
				{/if}

				<!-- Logo and Brand -->
				<div class="flex items-center">
					<div class="flex-shrink-0 flex items-center">
						<!-- Logo using favicon.svg -->
						<img src="/favicon.svg" alt="Beo Echo" class="w-8 h-8 mr-3" />
						<h1 class="text-xl font-bold text-gray-800 dark:text-white">Beo Echo</h1>
					</div>
				</div>

				<!-- Navigation Links -->
				{#if showNavigation}
					<nav class="hidden md:flex items-center space-x-6">
						<a 
							href="#features" 
							class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							title="View features"
							aria-label="View features"
						>
							Features
						</a>
						<a 
							href="#modes" 
							class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							title="View operating modes"
							aria-label="View operating modes"
						>
							Modes
						</a>
						<a 
							href="#pricing" 
							class="text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
							title="View pricing plans"
							aria-label="View pricing plans"
						>
							Pricing
						</a>
					</nav>
				{/if}
			</div>

			<!-- Right Side: Actions -->
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
					<img src="https://img.shields.io/github/stars/yogasw/beo-echo?style=social" alt="GitHub stars" class="h-5" />
				</a>

				<!-- Theme Toggle -->
				<button
					on:click={toggleTheme}
					class="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors"
					title={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
					aria-label={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
				>
					<i class="fas {$theme === 'dark' ? 'fa-sun text-yellow-400' : 'fa-moon text-gray-600'} text-sm"></i>
				</button>

				{#if $isAuthenticated && showUserMenu}
					<!-- User Menu -->
					<div class="flex items-center space-x-2">
						<button
							on:click={goToDashboard}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
							title="Go to dashboard"
							aria-label="Go to dashboard"
						>
							<i class="fas fa-tachometer-alt mr-2"></i>
							Dashboard
						</button>
						
						<!-- User Avatar & Dropdown -->
						<div class="relative">
							<button
								class="flex items-center p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
								title="User menu"
								aria-label="User menu"
							>
								<div class="w-8 h-8 bg-gradient-to-br from-purple-500 to-pink-500 rounded-full flex items-center justify-center">
									<i class="fas fa-user text-white text-xs"></i>
								</div>
								<i class="fas fa-chevron-down text-gray-400 ml-2 text-xs"></i>
							</button>
							
							<!-- Dropdown Menu (hidden by default, can be toggled with JavaScript) -->
							<div class="hidden absolute right-0 mt-2 w-48 bg-white dark:bg-gray-800 rounded-lg shadow-lg border border-gray-200 dark:border-gray-700 py-1">
								<a 
									href="/home" 
									class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
									title="Go to dashboard"
									aria-label="Go to dashboard"
								>
									<i class="fas fa-home mr-2"></i>
									Dashboard
								</a>
								<a 
									href="/settings" 
									class="block px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
									title="Go to settings"
									aria-label="Go to settings"
								>
									<i class="fas fa-cog mr-2"></i>
									Settings
								</a>
								<hr class="border-gray-200 dark:border-gray-700 my-1">
								<button
									on:click={handleLogout}
									class="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700"
									title="Logout from account"
									aria-label="Logout from account"
								>
									<i class="fas fa-sign-out-alt mr-2"></i>
									Logout
								</button>
							</div>
						</div>
					</div>
				{:else if showUserMenu}
					<div class="flex items-center space-x-2">
						<button
							on:click={handleLogin}
							class="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-all duration-200"
							title="Sign in to your account"
							aria-label="Sign in to your account"
						>
							Sign in
						</button>
						<button
							on:click={handleLogin}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
							title="Get started with Beo Echo"
							aria-label="Get started with Beo Echo"
						>
							<i class="fas fa-rocket mr-2"></i>
							Get Started
						</button>
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
		</div>
	</div>
</header>
