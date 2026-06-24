<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { toast } from '$lib/stores/toast';
	import AccessTokens from './AccessTokens.svelte';

	// The MCP endpoint lives on the same origin as the app, at /mcp.
	let endpoint = '';

	onMount(() => {
		if (browser) {
			endpoint = `${window.location.origin}/mcp`;
		}
	});

	$: desktopConfig = `{
  "mcpServers": {
    "beo-echo": {
      "url": "${endpoint}",
      "headers": { "Authorization": "Bearer YOUR_TOKEN" }
    }
  }
}`;

	async function copy(text: string, label: string) {
		try {
			await navigator.clipboard.writeText(text);
			toast.success(`${label} copied`);
		} catch {
			toast.error('Could not copy — select and copy manually');
		}
	}
</script>

<div class="space-y-6">
	<div>
		<h3 class="text-lg font-semibold theme-text-primary">Connect over MCP</h3>
		<p class="theme-text-secondary text-sm mt-1">
			Drive your mocks, routes, logs, replays, actions, config, and workspaces from AI clients
			(Claude, Cursor, VS Code, Claude Code) over the Model Context Protocol.
		</p>
	</div>

	<!-- Endpoint -->
	<div class="space-y-2">
		<span class="block text-sm font-medium theme-text-primary">MCP endpoint</span>
		<div class="flex items-center gap-2">
			<code class="flex-1 break-all theme-bg-primary theme-text-primary p-2 rounded border theme-border text-sm">
				{endpoint || '…'}
			</code>
			<button
				on:click={() => copy(endpoint, 'Endpoint')}
				class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded-md text-sm flex items-center gap-1"
				title="Copy endpoint URL"
				aria-label="Copy endpoint URL"
			>
				<i class="fas fa-copy"></i>
			</button>
		</div>
		<p class="text-xs theme-text-muted">This is the URL every MCP client needs.</p>
	</div>

	<!-- OAuth (Claude.ai) -->
	<div class="p-4 theme-bg-secondary rounded-lg space-y-2">
		<div class="flex items-center gap-2">
			<i class="fas fa-shield-halved text-purple-400"></i>
			<h4 class="font-medium theme-text-primary">Claude.ai &amp; OAuth-aware clients</h4>
		</div>
		<p class="theme-text-secondary text-sm">
			Paste just the endpoint URL — no token needed. The client discovers OAuth automatically,
			sends you here to log in and approve, then connects.
		</p>
		<ol class="text-sm theme-text-secondary list-decimal list-inside space-y-1">
			<li>In Claude.ai: Settings → Integrations → Add custom MCP server.</li>
			<li>Paste the endpoint URL above.</li>
			<li>Log in if prompted, then click <span class="theme-text-primary font-medium">Approve</span>.</li>
		</ol>
	</div>

	<!-- Static token clients -->
	<div class="p-4 theme-bg-secondary rounded-lg space-y-3">
		<div class="flex items-center gap-2">
			<i class="fas fa-key text-blue-400"></i>
			<h4 class="font-medium theme-text-primary">Claude Desktop, Cursor, VS Code</h4>
		</div>
		<p class="theme-text-secondary text-sm">
			These read a static <code class="theme-bg-primary px-1 rounded">Authorization: Bearer</code>
			header from their config. Generate a token below and drop it in.
		</p>
		<div class="relative">
			<pre class="theme-bg-primary theme-text-primary p-3 rounded border theme-border text-xs overflow-x-auto"><code>{desktopConfig}</code></pre>
			<button
				on:click={() => copy(desktopConfig, 'Config')}
				class="absolute top-2 right-2 bg-blue-600 hover:bg-blue-700 text-white px-2 py-1 rounded text-xs flex items-center gap-1"
				title="Copy config snippet"
				aria-label="Copy config snippet"
			>
				<i class="fas fa-copy"></i>
			</button>
		</div>
	</div>

	<!-- Tokens -->
	<div class="border-t theme-border pt-6">
		<AccessTokens />
	</div>
</div>
