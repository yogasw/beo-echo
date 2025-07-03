<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { isAuthenticated, auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';

	const dispatch = createEventDispatcher();

	// Props
	export let showUserMenu = true; // Show user menu when authenticated

	// Mobile menu state
	let mobileMenuOpen = false;

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

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}

	// Navigate to section and close mobile menu
	function navigateToSection(sectionId: string) {
		const element = document.getElementById(sectionId);
		if (element) {
			element.scrollIntoView({ behavior: 'smooth' });
		}
		closeMobileMenu();
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
			<div class="flex items-center space-x-3">
				<!-- GitHub Stars Badge - Hidden on mobile -->
				<a
					href="https://github.com/yogasw/beo-echo"
					target="_blank"
					rel="noopener noreferrer"
					class="hidden lg:block"
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
					<div class="flex items-center">
						<button
							on:click={goToDashboard}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
							title="Go to dashboard"
							aria-label="Go to dashboard"
						>
							<i class="fas fa-tachometer-alt mr-2"></i>
							<span class="hidden sm:inline">Dashboard</span>
							<span class="sm:hidden">Go</span>
						</button>
					</div>
				{:else}
					<!-- Login/Signup Buttons for non-authenticated users -->
					<div class="hidden md:flex items-center space-x-2">
						<a
							href="/login"
							class="px-4 py-2 text-sm font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-all duration-200"
							title="Sign in to your account"
							aria-label="Sign in to your account"
						>
							Sign In
						</a>
						<button
							on:click={() => goto('/login')}
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
							title="Get started with Beo Echo"
							aria-label="Get started with Beo Echo"
						>
							<i class="fas fa-rocket mr-2"></i>
							Get Started
						</button>
					</div>
					
					<!-- Mobile CTA Button -->
					<div class="md:hidden">
						<button
							on:click={() => goto('/login')}
							class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors flex items-center"
							title="Get started with Beo Echo"
							aria-label="Get started with Beo Echo"
						>
							<i class="fas fa-rocket mr-1"></i>
							Get Started
						</button>
					</div>
				{/if}

				<!-- Theme Toggle -->
				<button
					type="button"
					on:click={toggleTheme}
					class="w-10 h-10 flex items-center justify-center rounded-full bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors border border-gray-200 dark:border-gray-600"
					title={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
					aria-label={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
				>
					<i class="fas {$theme === 'dark' ? 'fa-sun' : 'fa-moon'} text-base text-gray-600 dark:text-gray-300"
					></i>
				</button>

				<!-- Mobile Menu Toggle -->
				<button
					on:click={toggleMobileMenu}
					class="md:hidden p-2 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
					title={mobileMenuOpen ? 'Close mobile menu' : 'Open mobile menu'}
					aria-label={mobileMenuOpen ? 'Close mobile menu' : 'Open mobile menu'}
					aria-expanded={mobileMenuOpen}
				>
					<i class="fas {mobileMenuOpen ? 'fa-times' : 'fa-bars'} text-gray-600 dark:text-gray-300"></i>
				</button>
			</div>
		</div>

		<!-- Mobile Menu -->
		{#if mobileMenuOpen}
			<div
				class="md:hidden border-t border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800"
				role="menu"
				aria-orientation="vertical"
			>
				<div class="px-2 pt-2 pb-3 space-y-1">
					<!-- Navigation Links -->
					<a
						href="/#features"
						class="block w-full text-left px-3 py-2 text-base font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
						title="View features"
						aria-label="View features"
						role="menuitem"
						on:click={closeMobileMenu}
					>
						<i class="fas fa-star mr-3"></i>
						Features
					</a>
					<a
						href="/#modes"
						class="block w-full text-left px-3 py-2 text-base font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
						title="View operating modes"
						aria-label="View operating modes"
						role="menuitem"
						on:click={closeMobileMenu}
					>
						<i class="fas fa-cogs mr-3"></i>
						Modes
					</a>
					<a
						href="/#pricing"
						class="block w-full text-left px-3 py-2 text-base font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
						title="View pricing plans"
						aria-label="View pricing plans"
						role="menuitem"
						on:click={closeMobileMenu}
					>
						<i class="fas fa-dollar-sign mr-3"></i>
						Pricing
					</a>

					{#if !$isAuthenticated}
						<!-- Authentication Links for Mobile -->
						<div class="border-t border-gray-200 dark:border-gray-700 pt-3 mt-3">
							<a
								href="/login"
								class="block w-full text-left px-3 py-2 text-base font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
								title="Sign in to your account"
								aria-label="Sign in to your account"
								role="menuitem"
								on:click={closeMobileMenu}
							>
								<i class="fas fa-sign-in-alt mr-3"></i>
								Sign In
							</a>
							<button
								on:click={() => { goto('/login'); closeMobileMenu(); }}
								class="block w-full text-left px-3 py-2 mt-2 text-base font-medium bg-blue-600 text-white hover:bg-blue-700 rounded-lg transition-colors"
								title="Get started with Beo Echo"
								aria-label="Get started with Beo Echo"
								role="menuitem"
							>
								<i class="fas fa-rocket mr-3"></i>
								Get Started
							</button>
						</div>
					{/if}

					<!-- GitHub Link for Mobile -->
					<div class="border-t border-gray-200 dark:border-gray-700 pt-3 mt-3">
						<button
							on:click={() => { window.open('https://github.com/yogasw/beo-echo', '_blank'); closeMobileMenu(); }}
							class="block w-full text-left px-3 py-2 text-base font-medium text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg transition-colors"
							title="Star us on GitHub"
							aria-label="Star Beo Echo on GitHub"
							role="menuitem"
						>
							<i class="fab fa-github mr-3"></i>
							GitHub
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</header>
