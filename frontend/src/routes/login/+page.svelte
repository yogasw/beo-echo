<script lang="ts">
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';
	import { auth, isAuthenticated } from '$lib/stores/auth';
	import Button from '$lib/components/ui/Button.svelte';
	import { BASE_URL_API } from '$lib/utils/authUtils';
	import LandingPageHeader from '$lib/components/landing-page/LandingPageHeader.svelte';
	import LandingPageFooter from '$lib/components/landing-page/LandingPageFooter.svelte';
	import { publicConfig, loadPublicConfig } from '$lib/stores/publicConfig';

	// Check if we're in landing mode (build time environment variable)
	const LANDING_MODE = import.meta.env.VITE_LANDING_MODE === 'true';

	let email = '';
	let password = '';
	let error = '';
	let loading = false;
	let showPassword = false;
	let configLoading = !LANDING_MODE; // Don't show loading if in landing mode

	// Load public configuration
	async function loadConfig() {
		try {
			configLoading = true;
			
			// If in landing mode, skip API calls and use defaults
			if (LANDING_MODE) {
				configLoading = false;
				return;
			}
			
			await loadPublicConfig();
		} catch (err) {
			console.error('Failed to load public config:', err);
		} finally {
			configLoading = false;
		}
	}

	onMount(async () => {
		// If in landing mode, skip API calls and show form immediately
		if (LANDING_MODE) {
			configLoading = false;
			return;
		}
		
		// Load public configuration first
		await loadConfig();
		
		if ($isAuthenticated) {
			goto('/home', { replaceState: true });
			// window.location.reload();
		}
		// TODO: When API is ready, fetch third party login availability
		// hasThirdPartyLogin = await auth.checkThirdPartyLoginAvailability();
	});

	async function handleLogin() {
		loading = true;
		error = '';

		try {
			if (email && password) {
				if (browser) {
					// New auth system
					await auth.login(email, password);
					await goto('/home', { replaceState: true });
					// window.location.reload();
				}
			} else {
				error = 'Please enter both email and password';
			}
		} catch (err: any) {
			error = err?.message || 'Authentication failed. Please try again.';
		} finally {
			loading = false;
		}
	}

	// Toggle password visibility
	function togglePasswordVisibility() {
		showPassword = !showPassword;
	}

	// Check for OAuth response on page load
	onMount(() => {
		if ($isAuthenticated) {
			goto('/home', { replaceState: true });
			return;
		}

		// Using setTimeout to ensure this runs after the component is fully mounted and
		// the URL parameters are fully available in the static build
		setTimeout(() => {
			const params = new URLSearchParams(window.location.search);
			const oauthError = params.get('error');
			const errorMessage = params.get('message');
			const token = params.get('token');
			const success = params.get('success');
			const user = params.get('user');
			const sso = params.get('sso');
			
			// Preserve any returnUrl or other important params
			const returnUrl = params.get('returnUrl');
			const preservedParams = new URLSearchParams();
			if (returnUrl) {
				preservedParams.set('returnUrl', returnUrl);
			}

			if (oauthError) {
				console.log('OAuth error detected:', oauthError, errorMessage);
				error = errorMessage || 'Authentication failed. Please contact an administrator.';
			} else if (success && token && user) {
				// Store the token and additional info
				localStorage.setItem('sso_provider', sso || '');
				auth.setToken(token);
				
				// If we have a returnUrl, go there, otherwise go to home
				if (returnUrl) {
					goto(returnUrl, { replaceState: true });
				} else {
					goto('/home', { replaceState: true });
				}
				return;
			}

			// Clean up sensitive params but preserve others
			const cleanUrl = preservedParams.toString() 
				? `/login?${preservedParams.toString()}`
				: '/login';
			
			if (oauthError || success) {
				window.history.replaceState({}, '', cleanUrl);
			}
		}, 100); // Small delay to ensure URL is fully processed
	});

	// Handler for Google login
	function handleGoogleLogin() {
		console.log('Google login clicked');
		loading = true;
		error = '';
		try {
			// Ensure we use the full absolute URL for redirects in production
			const currentURL = window.location.origin + '/login';
			// Append any existing query parameters except sensitive ones
			const params = new URLSearchParams(window.location.search);
			params.delete('error');
			params.delete('message');
			params.delete('token');
			params.delete('success');
			params.delete('user');
			params.delete('sso');
			
			const queryString = params.toString() ? `?${params.toString()}` : '';
			const redirectUri = `${currentURL}${queryString}`;
			
			console.log('Redirecting to Google OAuth with redirect URI:', redirectUri);
			window.location.href = `${BASE_URL_API}/oauth/google/login?redirect_uri=${encodeURIComponent(redirectUri)}`;
		} catch (err: any) {
			loading = false;
			error = err?.message || 'Google authentication failed. Please try again.';
		}
	}
</script>

<svelte:head>
	<title>Login - Beo Echo API Mock Service</title>
	<meta name="description" content="Login to access your Beo Echo API mock dashboard" />
</svelte:head>

<div class="min-h-screen flex flex-col theme-bg-tertiary">
	<!-- Header - always show in landing mode or when landing enabled -->
	{#if LANDING_MODE || (!configLoading && $publicConfig?.landing_enabled)}
		<LandingPageHeader 
			showUserMenu={false}
			on:back={() => goto('/')}
		/>
	{/if}

	<!-- Main Content -->
	<div class="flex-1 flex items-center justify-center relative">

		<div class="w-full max-w-md p-8">
			<div class="text-center mb-8">
				<h1 class="text-4xl font-bold theme-text-primary mb-2">Beo Echo</h1>
				<p class="theme-text-secondary">Login or register to your workspace</p>
			</div>

			<div class="theme-bg-primary rounded-lg theme-shadow p-8">
				<form class="space-y-6" on:submit|preventDefault={handleLogin}>
					<div class="space-y-4">
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
									title={showPassword ? 'Hide password' : 'Show password'}
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

					<Button
						type="submit"
						disabled={loading || !email || !password}
						loading={loading}
						fullWidth
					>
						Sign In
					</Button>
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
				<Button
					variant="outline"
					on:click={handleGoogleLogin}
					disabled={loading}
					loading={loading}
					fullWidth
				>
					<img src="/images/providers/google.svg" alt="Google" class="w-5 h-5" />
					<span>Continue with Google</span>
				</Button>
			</div>
		</div>
	</div>

	<!-- Footer - always show in landing mode or when landing enabled -->
	{#if LANDING_MODE || (!configLoading && $publicConfig?.landing_enabled)}
		<LandingPageFooter />
	{/if}
</div>
