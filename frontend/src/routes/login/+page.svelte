<script lang="ts">
	import { goto } from '$app/navigation';
	import { browser } from '$app/environment';
	import { setLocalStorage } from '$lib/utils/localStorage';
	import { fetchConfigsStore } from '$lib/stores/configurations';
	import { isAuthenticated } from '$lib/stores/authentication';
	import { login } from '$lib/api/mockoonApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	let username = '';
	let password = '';
	let error = '';

	async function handleLogin() {
		try {
			if (username && password) {
				if (browser) {
					setLocalStorage('username', username);
					setLocalStorage('password', password);
					await login({
						username,
						password
					}).then(async d => {
						isAuthenticated.set(true);
						await goto('/home');
					}).catch(e => {
						error = 'Login failed. Please try again.';
					});
				}
			} else {
				error = 'Please enter both username and password';
			}
		} catch (err) {
			error = 'Login failed. Please try again.';
		}
	}
</script>

<div class="min-h-screen w-full flex items-center justify-center theme-bg-tertiary">
	<div class="w-full max-w-md p-8">
		<div class="text-center mb-8">
			<h1 class="text-4xl font-bold theme-text-primary mb-2">Mockoon Control Panel</h1>
			<p class="theme-text-secondary">Sign in to manage your mock servers</p>
		</div>
		<div class="theme-bg-primary rounded-lg theme-shadow p-8">
			<form class="space-y-6" on:submit|preventDefault={handleLogin}>
				<div class="space-y-4">
					<div>
						<label for="username" class="block text-sm font-medium theme-text-secondary mb-1">Username</label>
						<input
							id="username"
							name="username"
							type="text"
							required
							bind:value={username}
							class="w-full px-4 py-3 theme-bg-secondary theme-border border rounded-lg theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
							placeholder="Enter your username"
						/>
					</div>
					<div>
						<label for="password" class="block text-sm font-medium theme-text-secondary mb-1">Password</label>
						<input
							id="password"
							name="password"
							type="password"
							required
							bind:value={password}
							class="w-full px-4 py-3 theme-bg-secondary theme-border border rounded-lg theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors"
							placeholder="Enter your password"
						/>
					</div>
				</div>

				{#if error}
					<div class="text-red-600 dark:text-red-400 text-sm text-center bg-red-500/10 p-3 rounded-lg border border-red-600/30 dark:border-red-400/30">{error}</div>
				{/if}

				<button
					type="submit"
					class="{ThemeUtils.primaryButton('w-full py-3 px-4 font-medium rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500')}"
				>
					Sign in
				</button>
			</form>
		</div>
	</div>
</div>
