<script lang="ts">
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { theme, toggleTheme } from '$lib/stores/theme';

	// For backwards compatibility with old login system
	import { setLocalStorage } from '$lib/utils/localStorage';
	import featureToggles, { FeatureFlags, getFeatureToggle } from '$lib/stores/featureToggles';

	let email = '';
	let password = '';
	let name = '';
	let error = '';
	let isLogin = true;
	let loading = false;
	let showPassword = false;
	let hasThirdPartyLogin = false; // This will be set based on API response

	onMount(async () => {
		if ($isAuthenticated) {
			goto('/');
			window.location.reload();
		}
		// TODO: When API is ready, fetch third party login availability
		// hasThirdPartyLogin = await auth.checkThirdPartyLoginAvailability();
	});

	async function handleLogin() {
		loading = true;
		error = '';

		try {
			if (isLogin) {
				// Login flow
				if (email && password) {
					if (browser) {
						// New auth system
						await auth.login(email, password);
						await goto('/');
						window.location.reload();
					}
				} else {
					error = 'Please enter both email and password';
				}
			} else {
				// Registration flow
				if (name && email && password) {
					await auth.register(name, email, password);
					await goto('/');
				} else {
					error = 'Please fill all fields';
				}
			}
		} catch (err: any) {
			error = err?.message || 'Authentication failed. Please try again.';
		} finally {
			loading = false;
		}
	}

	// Toggle between login and registration
	function toggleAuthMode() {
		isLogin = !isLogin;
		error = '';
	}

	// Toggle password visibility
	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	// Handler for Google login
	async function handleGoogleLogin() {
		loading = true;
		error = '';
		try {
			// TODO: Implement Google login when API is ready
			// await auth.loginWithGoogle();
			// await goto('/');
		} catch (err: any) {
			error = err?.message || 'Google authentication failed. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="min-h-screen w-full flex items-center justify-center theme-bg-tertiary relative">
	<!-- Theme Toggle Button -->
	<button
		type="button"
		on:click={toggleTheme}
		class="absolute top-4 right-4 w-10 h-10 flex items-center justify-center rounded-full theme-bg-secondary hover:bg-opacity-80 transition-colors border theme-border"
		aria-label="Toggle theme"
	>
		<i class="fas {$theme === 'dark' ? 'fa-sun' : 'fa-moon'} text-base theme-text-primary"></i>
	</button>

	<div class="w-full max-w-md p-8">
		<div class="text-center mb-8">
			<h1 class="text-4xl font-bold theme-text-primary mb-2">Beo Echo</h1>
			<p class="theme-text-secondary">Login or register to your workspace</p>
		</div>

		<div class="theme-bg-primary rounded-lg theme-shadow p-8">
			<form class="space-y-6" on:submit|preventDefault={handleLogin}>
				<div class="space-y-4">
					{#if !isLogin}
						<!-- Name field for registration -->
						<div>
							<label for="name" class="block text-sm font-medium theme-text-secondary mb-1"
								>Name</label
							>
							<div class="relative">
								<div class="absolute inset-y-0 left-0 pl-3 flex items-center">
									<span class="theme-text-muted">
										<i class="fas fa-user"></i>
									</span>
								</div>
								<input
									id="name"
									name="name"
									type="text"
									required={!isLogin}
									bind:value={name}
									disabled={loading}
									class="w-full pl-10 px-4 py-3 theme-bg-secondary theme-border border rounded-lg theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
									placeholder="Full name"
								/>
							</div>
						</div>
					{/if}

					<!-- Email field -->
					<div>
						<label for="email" class="block text-sm font-medium theme-text-secondary mb-1"
							>Email</label
						>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center">
								<span class="theme-text-muted">
									<i class="fas fa-envelope"></i>
								</span>
							</div>
							<input
								id="email"
								name="email"
								type="email"
								required
								bind:value={email}
								disabled={loading}
								class="w-full pl-10 px-4 py-3 theme-bg-secondary theme-border border rounded-lg theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
								placeholder="Email address"
							/>
						</div>
					</div>

					<!-- Password field with visibility toggle -->
					<div>
						<label for="password" class="block text-sm font-medium theme-text-secondary mb-1"
							>Password</label
						>
						<div class="relative">
							<div class="absolute inset-y-0 left-0 pl-3 flex items-center">
								<span class="theme-text-muted">
									<i class="fas fa-lock"></i>
								</span>
							</div>
							<input
								id="password"
								name="password"
								type={showPassword ? 'text' : 'password'}
								required
								bind:value={password}
								disabled={loading}
								class="w-full pl-10 pr-12 py-3 theme-bg-secondary theme-border border rounded-lg theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
								placeholder="Password"
							/>
							<button
								type="button"
								class="absolute inset-y-0 right-0 pr-3 flex items-center"
								on:click={togglePasswordVisibility}
								aria-label={showPassword ? 'Hide password' : 'Show password'}
							>
								<span class="theme-text-muted hover:text-indigo-500">
									<i class="fas {showPassword ? 'fa-eye-slash' : 'fa-eye'}"></i>
								</span>
							</button>
						</div>
					</div>
				</div>

				{#if error}
					<div
						class="text-red-600 dark:text-red-400 text-sm text-center bg-red-500/10 p-3 rounded-lg border border-red-600/30 dark:border-red-400/30"
					>
						{error}
					</div>
				{/if}

				<button
					type="submit"
					disabled={loading}
					class="w-full py-3 px-4 font-medium rounded-lg transition-colors duration-200 
					dark:bg-blue-600 dark:hover:bg-blue-700 dark:text-white
					bg-blue-500 hover:bg-blue-600 text-white
					focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50"
				>
					{#if loading}
						<span class="flex items-center justify-center">
							<svg
								class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
								xmlns="http://www.w3.org/2000/svg"
								fill="none"
								viewBox="0 0 24 24"
							>
								<circle
									class="opacity-25"
									cx="12"
									cy="12"
									r="10"
									stroke="currentColor"
									stroke-width="4"
								></circle>
								<path
									class="opacity-75"
									fill="currentColor"
									d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
								></path>
							</svg>
							{isLogin ? 'Signing in...' : 'Creating account...'}
						</span>
					{:else}
						{isLogin ? 'Sign In' : 'Create Account'}
					{/if}
				</button>

				<!-- Toggle auth mode -->
				{#if getFeatureToggle(FeatureFlags.FEATURE_REGISTER_EMAIL_ENABLED)}
					<div class="text-center pt-2">
						<button
							type="button"
							on:click={toggleAuthMode}
							class="text-blue-400 hover:text-blue-300 text-sm"
						>
							{isLogin ? 'Need an account? Register' : 'Already have an account? Sign in'}
						</button>
					</div>
				{/if}
			</form>
		</div>

		<!-- Divider -->
		<div class="relative my-6">
			<div class="absolute inset-0 flex items-center">
				<div class="w-full border-t theme-border"></div>
			</div>
			<div class="relative flex justify-center text-sm">
				<span class="px-4 theme-bg-tertiary theme-text-secondary">OR</span>
			</div>
		</div>

		<!-- Social Login Buttons -->
		<div class="space-y-3">
			<button
				type="button"
				on:click={handleGoogleLogin}
				class="w-full flex items-center justify-center gap-2 p-3 rounded-lg border transition-colors duration-200
                     dark:bg-gray-700 dark:hover:bg-gray-600 dark:border-gray-600 dark:text-white
                     bg-white hover:bg-gray-50 border-gray-200 text-gray-700
                     focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
			>
				<img src="https://authjs.dev/img/providers/google.svg" alt="Google" class="w-5 h-5" />
				<span>Continue with Google</span>
			</button>
		</div>
	</div>
</div>
