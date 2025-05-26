<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let authType: string = 'No Auth';
	export let authConfig: Record<string, any> = {};

	const dispatch = createEventDispatcher();

	function updateAuthType(type: string) {
		authType = type;
		// Reset config when changing auth type
		authConfig = {};
		dispatch('authChange', { authType, authConfig });
	}

	function updateAuthConfig(key: string, value: any) {
		authConfig = { ...authConfig, [key]: value };
		dispatch('authChange', { authType, authConfig });
	}

	// Auth type options
	const authTypes = [
		{ value: 'No Auth', label: 'No Auth', icon: 'fas fa-ban' },
		{ value: 'Bearer Token', label: 'Bearer Token', icon: 'fas fa-key' },
		{ value: 'Basic Auth', label: 'Basic Auth', icon: 'fas fa-user-shield' },
		{ value: 'API Key', label: 'API Key', icon: 'fas fa-code' },
		{ value: 'OAuth 2.0', label: 'OAuth 2.0', icon: 'fas fa-lock' },
		{ value: 'Digest Auth', label: 'Digest Auth', icon: 'fas fa-fingerprint' }
	];

	$: selectedAuthType = authTypes.find(type => type.value === authType) || authTypes[0];
</script>

<!-- Authorization section -->
<div role="tabpanel" aria-labelledby="auth-tab">
	<div class="flex justify-between items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-shield-alt text-green-500 mr-2"></i>
			Authorization
		</h2>
		<div class="flex items-center space-x-2">
			<span class="text-xs theme-text-muted">
				Type: <span class="font-medium">{authType}</span>
			</span>
		</div>
	</div>

	<div class="space-y-4">
		<!-- Auth Type Selector -->
		<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4">
			<label for="auth-type" class="block text-sm font-medium theme-text-secondary mb-2">
				<i class="{selectedAuthType.icon} text-green-400 mr-2"></i>
				Authentication Type
			</label>
			<select
				id="auth-type"
				class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary transition-all duration-200"
				title="Select authentication method"
				aria-label="Authentication type"
				value={authType}
				on:change={(e) => updateAuthType(e.currentTarget.value)}
			>
				{#each authTypes as type}
					<option value={type.value}>{type.label}</option>
				{/each}
			</select>
		</div>

		<!-- Dynamic Auth Configuration -->
		{#if authType === 'Bearer Token'}
			<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 space-y-3">
				<h3 class="text-sm font-medium theme-text-primary flex items-center">
					<i class="fas fa-key text-yellow-400 mr-2"></i>
					Bearer Token Configuration
				</h3>
				<div>
					<label for="bearer-token" class="block text-sm font-medium theme-text-secondary mb-2">
						Token
					</label>
					<input
						id="bearer-token"
						type="password"
						class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
						placeholder="Enter your bearer token"
						title="Bearer token for authorization"
						aria-label="Bearer token input"
						value={authConfig.token || ''}
						on:input={(e) => updateAuthConfig('token', e.currentTarget.value)}
					/>
				</div>
				<p class="text-xs theme-text-muted">
					<i class="fas fa-info-circle mr-1"></i>
					This token will be included in the Authorization header as "Bearer {authConfig.token || '[your_token]'}"
				</p>
			</div>

		{:else if authType === 'Basic Auth'}
			<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 space-y-3">
				<h3 class="text-sm font-medium theme-text-primary flex items-center">
					<i class="fas fa-user-shield text-blue-400 mr-2"></i>
					Basic Authentication Configuration
				</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label for="basic-username" class="block text-sm font-medium theme-text-secondary mb-2">
							Username
						</label>
						<input
							id="basic-username"
							type="text"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
							placeholder="Enter username"
							title="Username for basic authentication"
							aria-label="Username input"
							value={authConfig.username || ''}
							on:input={(e) => updateAuthConfig('username', e.currentTarget.value)}
						/>
					</div>
					<div>
						<label for="basic-password" class="block text-sm font-medium theme-text-secondary mb-2">
							Password
						</label>
						<input
							id="basic-password"
							type="password"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
							placeholder="Enter password"
							title="Password for basic authentication"
							aria-label="Password input"
							value={authConfig.password || ''}
							on:input={(e) => updateAuthConfig('password', e.currentTarget.value)}
						/>
					</div>
				</div>
				<p class="text-xs theme-text-muted">
					<i class="fas fa-info-circle mr-1"></i>
					Credentials will be Base64 encoded and sent in the Authorization header
				</p>
			</div>

		{:else if authType === 'API Key'}
			<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 space-y-3">
				<h3 class="text-sm font-medium theme-text-primary flex items-center">
					<i class="fas fa-code text-purple-400 mr-2"></i>
					API Key Configuration
				</h3>
				<div class="space-y-4">
					<div>
						<label for="api-key-name" class="block text-sm font-medium theme-text-secondary mb-2">
							Key Name
						</label>
						<input
							id="api-key-name"
							type="text"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
							placeholder="X-API-Key"
							title="API key header name"
							aria-label="API key name"
							value={authConfig.keyName || ''}
							on:input={(e) => updateAuthConfig('keyName', e.currentTarget.value)}
						/>
					</div>
					<div>
						<label for="api-key-value" class="block text-sm font-medium theme-text-secondary mb-2">
							Value
						</label>
						<input
							id="api-key-value"
							type="password"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
							placeholder="Enter your API key"
							title="API key value"
							aria-label="API key value"
							value={authConfig.keyValue || ''}
							on:input={(e) => updateAuthConfig('keyValue', e.currentTarget.value)}
						/>
					</div>
					<div>
						<label for="api-key-location" class="block text-sm font-medium theme-text-secondary mb-2">
							Add to
						</label>
						<select
							id="api-key-location"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary transition-all duration-200"
							title="Where to include the API key"
							aria-label="API key location"
							value={authConfig.location || 'Header'}
							on:change={(e) => updateAuthConfig('location', e.currentTarget.value)}
						>
							<option value="Header">Header</option>
							<option value="Query Params">Query Parameters</option>
						</select>
					</div>
				</div>
				<p class="text-xs theme-text-muted">
					<i class="fas fa-info-circle mr-1"></i>
					API key will be added to the request {authConfig.location === 'Query Params' ? 'as a query parameter' : 'as a header'}
				</p>
			</div>

		{:else if authType === 'OAuth 2.0'}
			<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 space-y-3">
				<h3 class="text-sm font-medium theme-text-primary flex items-center">
					<i class="fas fa-lock text-red-400 mr-2"></i>
					OAuth 2.0 Configuration
				</h3>
				<div class="space-y-4">
					<div>
						<label for="oauth-token" class="block text-sm font-medium theme-text-secondary mb-2">
							Access Token
						</label>
						<input
							id="oauth-token"
							type="password"
							class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
							placeholder="Enter access token"
							title="OAuth 2.0 access token"
							aria-label="OAuth access token"
							value={authConfig.accessToken || ''}
							on:input={(e) => updateAuthConfig('accessToken', e.currentTarget.value)}
						/>
					</div>
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label for="oauth-header-prefix" class="block text-sm font-medium theme-text-secondary mb-2">
								Header Prefix
							</label>
							<select
								id="oauth-header-prefix"
								class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary transition-all duration-200"
								title="OAuth token header prefix"
								aria-label="OAuth header prefix"
								value={authConfig.headerPrefix || 'Bearer'}
								on:change={(e) => updateAuthConfig('headerPrefix', e.currentTarget.value)}
							>
								<option value="Bearer">Bearer</option>
								<option value="Token">Token</option>
								<option value="OAuth">OAuth</option>
							</select>
						</div>
						<div class="flex items-end">
							<button
								class="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 px-4 rounded-md text-sm transition-all duration-200"
								title="Configure OAuth 2.0 flow settings"
								aria-label="Configure OAuth settings"
							>
								<i class="fas fa-cog mr-2"></i>
								Configure OAuth Flow
							</button>
						</div>
					</div>
				</div>
				<p class="text-xs theme-text-muted">
					<i class="fas fa-info-circle mr-1"></i>
					For full OAuth 2.0 flow configuration, use the Configure OAuth Flow button
				</p>
			</div>

		{:else if authType === 'No Auth'}
			<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 text-center">
				<div class="py-8">
					<i class="fas fa-ban text-gray-400 text-4xl mb-4"></i>
					<h3 class="text-sm font-medium theme-text-primary mb-2">No Authentication</h3>
					<p class="text-xs theme-text-muted">
						This request will not include any authentication headers
					</p>
				</div>
			</div>
		{/if}

		<!-- Preview Section -->
		{#if authType !== 'No Auth'}
			<div class="bg-gray-50 dark:bg-gray-900 border theme-border rounded-lg p-4">
				<h4 class="text-xs font-medium theme-text-muted mb-2 uppercase tracking-wide">
					<i class="fas fa-eye mr-1"></i>
					Preview
				</h4>
				<div class="text-sm theme-text-secondary font-mono bg-white dark:bg-gray-800 p-3 rounded border theme-border">
					{#if authType === 'Bearer Token' && authConfig.token}
						<span class="text-blue-400">Authorization:</span> Bearer {authConfig.token.substring(0, 10)}***
					{:else if authType === 'Basic Auth' && authConfig.username}
						<span class="text-blue-400">Authorization:</span> Basic {btoa(`${authConfig.username}:${authConfig.password || ''}`).substring(0, 10)}***
					{:else if authType === 'API Key' && authConfig.keyName}
						{#if authConfig.location === 'Query Params'}
							<span class="text-blue-400">Query Parameter:</span> {authConfig.keyName}={authConfig.keyValue?.substring(0, 10) || ''}***
						{:else}
							<span class="text-blue-400">{authConfig.keyName}:</span> {authConfig.keyValue?.substring(0, 10) || ''}***
						{/if}
					{:else if authType === 'OAuth 2.0' && authConfig.accessToken}
						<span class="text-blue-400">Authorization:</span> {authConfig.headerPrefix || 'Bearer'} {authConfig.accessToken.substring(0, 10)}***
					{:else}
						<span class="text-gray-400 italic">Configure authentication details to see preview</span>
					{/if}
				</div>
			</div>
		{/if}
	</div>
</div>
