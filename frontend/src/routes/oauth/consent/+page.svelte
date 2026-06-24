<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { auth, isAuthenticated, currentUser } from '$lib/stores/auth';
	import { BASE_URL_API } from '$lib/utils/authUtils';

	let requestId = '';
	let isWorking = false;
	let error = '';
	let ready = false;

	onMount(() => {
		if (!browser) return;

		requestId = $page.url.searchParams.get('request_id') ?? '';

		if (!requestId) {
			error = 'Missing authorization request. Start the connection from your MCP client again.';
			ready = true;
			return;
		}

		// The consent decision must be made by a logged-in user. If not
		// authenticated, bounce to login and come back to this exact URL.
		if (!$isAuthenticated || !auth.getToken()) {
			const returnTo = `/oauth/consent?request_id=${encodeURIComponent(requestId)}`;
			goto(`/login?returnUrl=${encodeURIComponent(returnTo)}`, { replaceState: true });
			return;
		}

		ready = true;
	});

	async function decide(approve: boolean) {
		if (!requestId || isWorking) return;
		isWorking = true;
		error = '';

		try {
			const res = await fetch(`${BASE_URL_API}/oauth/mcp/approve`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${auth.getToken()}`
				},
				body: JSON.stringify({ request_id: requestId, approve })
			});

			const data = await res.json();

			if (!res.ok || !data.success) {
				throw new Error(data.message || 'Failed to process the authorization request');
			}

			// Hand control back to the MCP client at its redirect URL (carries the
			// authorization code on approve, or an error on deny).
			if (data.redirect_url) {
				window.location.href = data.redirect_url;
			} else {
				error = 'The server did not return a redirect URL.';
				isWorking = false;
			}
		} catch (e: any) {
			error = e.message || 'Something went wrong';
			isWorking = false;
		}
	}
</script>

<div class="min-h-screen flex items-center justify-center theme-bg-secondary p-4">
	<div class="w-full max-w-md theme-bg-primary rounded-lg shadow-lg p-6 space-y-6">
		<div class="text-center">
			<div class="w-14 h-14 mx-auto rounded-full bg-blue-500 flex items-center justify-center text-white text-2xl mb-3">
				<i class="fas fa-plug"></i>
			</div>
			<h1 class="text-xl font-semibold theme-text-primary">Authorize MCP access</h1>
			<p class="theme-text-secondary text-sm mt-2">
				An MCP client wants to access Beo Echo on your behalf.
			</p>
		</div>

		{#if error}
			<div class="p-3 rounded-md bg-red-50 dark:bg-red-900/30 border border-red-300 dark:border-red-700 text-red-800 dark:text-red-200 text-sm">
				<i class="fas fa-exclamation-circle mr-1"></i>
				{error}
			</div>
		{/if}

		{#if ready && !error}
			<div class="theme-bg-secondary rounded-md p-4 text-sm theme-text-secondary space-y-1">
				<div class="flex justify-between">
					<span>Signed in as</span>
					<span class="theme-text-primary font-medium">{$currentUser?.email ?? '—'}</span>
				</div>
				<div class="flex justify-between">
					<span>Access granted</span>
					<span class="theme-text-primary font-medium">Your workspaces &amp; projects</span>
				</div>
			</div>

			<p class="text-xs theme-text-muted">
				Approving issues an access token to this client. You can revoke it any time from
				Profile → Access Tokens.
			</p>

			<div class="flex gap-3">
				<button
					on:click={() => decide(false)}
					disabled={isWorking}
					class="flex-1 theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
					title="Deny access"
					aria-label="Deny access"
				>
					Deny
				</button>
				<button
					on:click={() => decide(true)}
					disabled={isWorking}
					class="flex-1 bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center justify-center gap-2"
					title="Approve access"
					aria-label="Approve access"
				>
					{#if isWorking}
						<i class="fas fa-spinner fa-spin"></i>
						<span>Working...</span>
					{:else}
						<i class="fas fa-check"></i>
						<span>Approve</span>
					{/if}
				</button>
			</div>
		{:else if !error}
			<div class="flex items-center justify-center gap-2 theme-text-secondary text-sm py-4">
				<i class="fas fa-spinner fa-spin"></i>
				<span>Loading...</span>
			</div>
		{/if}
	</div>
</div>
