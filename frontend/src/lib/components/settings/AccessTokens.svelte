<script lang="ts">
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast';
	import {
		listApiTokens,
		createApiToken,
		revokeApiToken,
		type ApiToken,
		type CreatedApiToken
	} from '$lib/api/tokensApi';

	let tokens: ApiToken[] = [];
	let isLoading = true;
	let isCreating = false;

	// New token form
	let showCreateForm = false;
	let newTokenName = '';
	let newTokenExpiresDays = 0; // 0 = never

	// The plaintext of a freshly created token, shown once.
	let createdToken: CreatedApiToken | null = null;

	onMount(loadTokens);

	async function loadTokens() {
		isLoading = true;
		try {
			tokens = await listApiTokens();
		} catch (error) {
			toast.error('Failed to load access tokens');
			console.error('Load tokens error:', error);
		} finally {
			isLoading = false;
		}
	}

	async function createToken() {
		if (!newTokenName.trim()) {
			toast.error('Token name cannot be empty');
			return;
		}
		isCreating = true;
		try {
			createdToken = await createApiToken(newTokenName.trim(), newTokenExpiresDays);
			toast.success('Token created — copy it now, it will not be shown again');
			newTokenName = '';
			newTokenExpiresDays = 0;
			showCreateForm = false;
			await loadTokens();
		} catch (error) {
			toast.error('Failed to create token');
			console.error('Create token error:', error);
		} finally {
			isCreating = false;
		}
	}

	async function revokeToken(token: ApiToken) {
		if (!confirm(`Revoke token "${token.name}"? Any client using it will stop working.`)) {
			return;
		}
		try {
			await revokeApiToken(token.id);
			toast.success('Token revoked');
			await loadTokens();
		} catch (error) {
			toast.error('Failed to revoke token');
			console.error('Revoke token error:', error);
		}
	}

	async function copyToken() {
		if (!createdToken) return;
		try {
			await navigator.clipboard.writeText(createdToken.token);
			toast.success('Token copied to clipboard');
		} catch {
			toast.error('Could not copy — select and copy manually');
		}
	}

	function dismissCreatedToken() {
		createdToken = null;
	}

	function formatDate(value: string | null): string {
		if (!value) return 'Never';
		return new Date(value).toLocaleString();
	}
</script>

<div class="space-y-6">
	<div class="flex items-start justify-between gap-4">
		<div>
			<h3 class="text-lg font-semibold theme-text-primary">Access Tokens</h3>
			<p class="theme-text-secondary text-sm mt-1">
				Personal access tokens authenticate the MCP server and CLIs as you. Use one as a
				<code class="theme-bg-primary px-1 rounded">Bearer</code> token in your MCP client config.
			</p>
		</div>
		{#if !showCreateForm}
			<button
				on:click={() => (showCreateForm = true)}
				class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center gap-2 whitespace-nowrap"
				title="Generate a new access token"
				aria-label="Generate a new access token"
			>
				<i class="fas fa-plus"></i>
				<span>New Token</span>
			</button>
		{/if}
	</div>

	<!-- Freshly created token (shown once) -->
	{#if createdToken}
		<div class="p-4 rounded-lg border border-green-300 dark:border-green-700 bg-green-50 dark:bg-green-900/30 space-y-2">
			<div class="flex items-center gap-2 text-green-800 dark:text-green-200 text-sm font-medium">
				<i class="fas fa-check-circle"></i>
				<span>Token created — copy it now. You will not be able to see it again.</span>
			</div>
			<div class="flex items-center gap-2">
				<code class="flex-1 break-all theme-bg-primary theme-text-primary p-2 rounded border theme-border text-xs">
					{createdToken.token}
				</code>
				<button
					on:click={copyToken}
					class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded-md text-sm flex items-center gap-1"
					title="Copy token to clipboard"
					aria-label="Copy token to clipboard"
				>
					<i class="fas fa-copy"></i>
				</button>
				<button
					on:click={dismissCreatedToken}
					class="theme-bg-secondary theme-text-primary px-3 py-2 rounded-md text-sm"
					title="Dismiss"
					aria-label="Dismiss created token"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
		</div>
	{/if}

	<!-- Create form -->
	{#if showCreateForm}
		<form on:submit|preventDefault={createToken} class="p-4 theme-bg-primary rounded-lg space-y-4">
			<div class="space-y-2">
				<label for="tokenName" class="block text-sm font-medium theme-text-primary">Token name</label>
				<input
					type="text"
					id="tokenName"
					bind:value={newTokenName}
					placeholder="e.g. MCP on my laptop"
					class="block w-full p-2 theme-bg-primary theme-text-primary border theme-border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
				/>
			</div>
			<div class="space-y-2">
				<label for="tokenExpiry" class="block text-sm font-medium theme-text-primary">Expires</label>
				<select
					id="tokenExpiry"
					bind:value={newTokenExpiresDays}
					class="block w-full p-2 theme-bg-primary theme-text-primary border theme-border rounded-md"
				>
					<option value={0}>Never</option>
					<option value={30}>30 days</option>
					<option value={90}>90 days</option>
					<option value={365}>1 year</option>
				</select>
			</div>
			<div class="flex justify-end gap-3">
				<button
					type="button"
					on:click={() => {
						showCreateForm = false;
						newTokenName = '';
					}}
					class="theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
					title="Cancel"
					aria-label="Cancel token creation"
				>
					Cancel
				</button>
				<button
					type="submit"
					disabled={isCreating}
					class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center gap-2"
					title="Generate token"
					aria-label="Generate token"
				>
					{#if isCreating}
						<i class="fas fa-spinner fa-spin"></i>
						<span>Generating...</span>
					{:else}
						<i class="fas fa-key"></i>
						<span>Generate</span>
					{/if}
				</button>
			</div>
		</form>
	{/if}

	<!-- Token list -->
	{#if isLoading}
		<div class="flex items-center gap-2 theme-text-secondary text-sm p-4">
			<i class="fas fa-spinner fa-spin"></i>
			<span>Loading tokens...</span>
		</div>
	{:else if tokens.length === 0}
		<div class="text-center theme-text-muted text-sm p-6 theme-bg-primary rounded-lg">
			<i class="fas fa-key text-2xl mb-2 block opacity-50"></i>
			No access tokens yet.
		</div>
	{:else}
		<div class="space-y-2">
			{#each tokens as token (token.id)}
				<div class="flex items-center justify-between gap-4 p-3 theme-bg-primary rounded-lg border theme-border">
					<div class="min-w-0">
						<div class="flex items-center gap-2">
							<span class="font-medium theme-text-primary truncate">{token.name}</span>
							{#if token.source === 'oauth'}
								<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-100">
									OAuth
								</span>
							{/if}
						</div>
						<div class="text-xs theme-text-muted mt-0.5 flex flex-wrap gap-x-3">
							<span><code>{token.prefix}…</code></span>
							<span>Created {formatDate(token.created_at)}</span>
							<span>Last used {formatDate(token.last_used_at)}</span>
							<span>Expires {formatDate(token.expires_at)}</span>
						</div>
					</div>
					<button
						on:click={() => revokeToken(token)}
						class="text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 px-3 py-2 rounded-md text-sm flex items-center gap-1 whitespace-nowrap"
						title="Revoke this token"
						aria-label={`Revoke token ${token.name}`}
					>
						<i class="fas fa-trash"></i>
						<span>Revoke</span>
					</button>
				</div>
			{/each}
		</div>
	{/if}
</div>
