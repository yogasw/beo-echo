<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { ssoStore, type GoogleOAuthConfig } from '$lib/stores/ssoStore';
	import { ssoApi } from '$lib/api/ssoApi';
	import GoogleOAuthConfigModal from './GoogleOAuthConfigModal.svelte';
	
	export let visible = false;

	let showGoogleConfig = false;
	let loading = false;
	let error: string | null = null;

	$: googleConfigured = !!$ssoStore.googleConfig;
	$: googleEnabled = $ssoStore.isGoogleEnabled;

	onMount(async () => {
		loading = true;
		try {
			const { config, enabled } = await ssoApi.getGoogleConfig();
			ssoStore.setGoogleConfig(config);
			ssoStore.setGoogleEnabled(enabled);
		} catch (err) {
			error = 'Failed to load SSO configuration';
		} finally {
			loading = false;
		}
	});

	async function handleAutoCreateAccountsChange(enabled: boolean) {
		try {
			await ssoApi.updateSystemConfig('FEATURE_OAUTH_AUTO_REGISTER', enabled);
			ssoStore.setAutoCreateAccounts(enabled);
		} catch (err) {
			error = 'Failed to update auto-create accounts setting';
		}
	}

	async function handleAllowLocalAuthChange(enabled: boolean) {
		try {
			await ssoApi.updateSystemConfig('FEATURE_LOCAL_AUTH', enabled);
			ssoStore.setAllowLocalAuth(enabled);
		} catch (err) {
			error = 'Failed to update local authentication setting';
		}
	}

	async function handleDefaultSignInMethodChange(method: string) {
		try {
			await ssoApi.updateSystemConfig('DEFAULT_SIGN_IN_METHOD', method);
			ssoStore.setDefaultSignInMethod(method as any);
		} catch (err) {
			error = 'Failed to update default sign-in method';
		}
	}
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<div class="space-y-5">
			{#if error}
				<div class="bg-red-900/20 border border-red-700 text-red-500 p-3 rounded-md text-sm">
					{error}
				</div>
			{/if}

			<div>
				<div class="flex items-center justify-between mb-3">
					<h3 class="theme-text-primary font-medium">SSO Providers</h3>
					<div>
						<button class="px-3 py-2 theme-bg-secondary hover:bg-gray-600 theme-text-primary rounded-md text-sm flex items-center gap-2">
							<i class="fas fa-plus"></i>
							<span>Add Provider</span>
						</button>
					</div>
				</div>
				
				<div class="space-y-4">
					<!-- Google SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-white rounded flex items-center justify-center">
									<i class="fab fa-google text-gray-800"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">Google</h4>
									{#if loading}
										<p class="theme-text-muted text-xs">Loading...</p>
									{:else if googleConfigured}
										<p class="flex items-center gap-2">
											<span class="text-xs {googleEnabled ? 'text-green-400' : 'text-yellow-400'}">
												{googleEnabled ? 'Enabled' : 'Disabled'}
											</span>
											{#if googleEnabled}
												<span class="relative flex h-2 w-2">
													<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
													<span class="relative inline-flex rounded-full h-2 w-2 bg-green-500"></span>
												</span>
											{/if}
										</p>
									{:else}
										<p class="theme-text-muted text-xs">Not configured</p>
									{/if}
								</div>
							</div>
							<div>
								<button 
									class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs"
									on:click={() => showGoogleConfig = true}
									disabled={loading}>
									{googleConfigured ? 'Edit' : 'Configure'}
								</button>
							</div>
						</div>
					</div>

					<GoogleOAuthConfigModal
						visible={showGoogleConfig}
						onClose={() => showGoogleConfig = false}
					/>
					
					<!-- GitHub SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-gray-900 rounded flex items-center justify-center">
									<i class="fab fa-github text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">GitHub</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
					
					<!-- Microsoft SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
									<i class="fab fa-microsoft text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">Microsoft</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
					
					<!-- SAML SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-blue-600 rounded flex items-center justify-center">
									<i class="fas fa-id-card text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">SAML</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
			
			<div class="pt-4 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-3">SSO Settings</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Default Sign-in Method</p>
							<p class="theme-text-muted text-xs">Select which sign-in method to display by default</p>
						</div>
						<select 
							class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary"
							bind:value={$ssoStore.defaultSignInMethod}
							on:change={(e) => handleDefaultSignInMethodChange(e.currentTarget.value)}
							disabled={loading}>
							<option value="local">Local Authentication</option>
							<option value="google" disabled={!googleConfigured}>Google</option>
							<option value="github" disabled>GitHub</option>
							<option value="microsoft" disabled>Microsoft</option>
							<option value="saml" disabled>SAML</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Allow Local Authentication</p>
							<p class="theme-text-muted text-xs">Enable login with username and password</p>
						</div>
						<label class="relative inline-flex items-center cursor-pointer">
							<input 
								type="checkbox" 
								class="sr-only peer" 
								bind:checked={$ssoStore.allowLocalAuth}
								on:change={(e) => handleAllowLocalAuthChange(e.currentTarget.checked)}
								disabled={loading}>
							<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
							peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
							rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
							after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
							after:transition-all peer-checked:bg-blue-600"></div>
						</label>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Auto-create User Accounts</p>
							<p class="theme-text-muted text-xs">Automatically create accounts for new SSO users</p>
						</div>
						<label class="relative inline-flex items-center cursor-pointer">
							<input 
								type="checkbox" 
								class="sr-only peer"
								bind:checked={$ssoStore.autoCreateAccounts}
								on:change={(e) => handleAutoCreateAccountsChange(e.currentTarget.checked)}
								disabled={loading}>
							<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
							peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
							rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
							after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
							after:transition-all peer-checked:bg-blue-600"></div>
						</label>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
