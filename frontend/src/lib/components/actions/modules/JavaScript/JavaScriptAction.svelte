<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import DebugPanel from './DebugPanel.svelte';
	import type { JavaScriptConfig } from '$lib/types/Action';

	export let config: JavaScriptConfig | null = null;
	export let executionPoint: 'before_request' | 'after_request' = 'after_request';

	const dispatch = createEventDispatcher();

	let showDebugPanel = false;

	// Initialize script with default template
	let script =
		config?.script ||
		`// Modify request or response
// Available variables: request, response (if after_request)
// Use console.log() to debug

// Example: Add a custom header
if (request) {
	request.headers["X-Custom-Header"] = "MyValue";
}

${executionPoint === 'after_request' ? `// Example: Modify response body
if (response) {
	try {
		var body = JSON.parse(response.body);
		body.modified = true;
		response.body = JSON.stringify(body);
	} catch(e) {
		console.log("Response is not JSON");
	}
}` : ''}
`;

	let editor: any;

	// Update parent when script changes
	function handleScriptChange(event: CustomEvent) {
		script = event.detail;
		const newConfig: JavaScriptConfig = { script };
		dispatch('change', newConfig);
	}

	// Update template when execution point changes
	$: {
		if (!config) {
			// Only update template for new actions
			if (executionPoint === 'after_request' && !script.includes('response')) {
				script += `\n\n// Example: Modify response
if (response) {
	console.log("Response status:", response.status_code);
}`;
			}
		}
	}

	// Format code
	function formatCode() {
		if (editor) {
			editor.format();
		}
	}

	onMount(() => {
		// Dispatch initial config
		const initialConfig: JavaScriptConfig = { script };
		dispatch('change', initialConfig);
	});
</script>

<div class="space-y-4">
	<!-- JavaScript Editor -->
	<div>
		<div class="flex items-center justify-between mb-2">
			<label for="script-editor" class="block text-sm font-medium theme-text-primary">
				JavaScript Code <span class="text-red-500">*</span>
			</label>
			<div class="flex items-center gap-2">
				<button
					type="button"
					on:click={() => showDebugPanel = true}
					class="text-xs px-3 py-1.5 rounded bg-green-600 hover:bg-green-700 text-white transition-colors flex items-center gap-1.5"
					title="Test script with sample log data"
					aria-label="Test JavaScript code"
				>
					<i class="fas fa-play text-xs"></i>
					Test Script
				</button>
				<button
					type="button"
					on:click={formatCode}
					class="text-xs px-3 py-1.5 rounded bg-gray-700 hover:bg-gray-600 text-white transition-colors flex items-center gap-1.5"
					title="Format code"
					aria-label="Format JavaScript code"
				>
					<i class="fas fa-magic text-xs"></i>
					Format
				</button>
			</div>
		</div>

		<div class="border theme-border rounded-lg overflow-hidden" style="height: 350px;">
			<MonacoEditor
				bind:this={editor}
				bind:value={script}
				language="javascript"
				on:change={handleScriptChange}
			/>
		</div>

		<p id="script-help" class="mt-2 text-xs theme-text-secondary">
			Write JavaScript code to modify the request or response
		</p>
	</div>

	<!-- Context Information -->
	<div class="p-4 bg-blue-50/70 dark:bg-gray-900/50 rounded-lg border border-blue-200/50 dark:border-blue-900/50">
		<div class="flex items-start gap-2 mb-3">
			<i class="fas fa-info-circle text-blue-500 mt-0.5"></i>
			<div>
				<h4 class="text-sm font-semibold theme-text-primary">Available Context</h4>
				<p class="text-xs theme-text-secondary mt-1">
					The following objects are available in your JavaScript code:
				</p>
			</div>
		</div>

		<div class="space-y-3 text-xs">
			<!-- Request Object -->
			<div class="bg-white/50 dark:bg-gray-800/50 rounded p-2">
				<div class="font-semibold theme-text-primary mb-1 font-mono">request</div>
				<div class="space-y-1 theme-text-secondary pl-2">
					<div><code class="text-purple-600 dark:text-purple-400">method</code> - HTTP method (GET, POST, etc.)</div>
					<div><code class="text-purple-600 dark:text-purple-400">path</code> - Request path</div>
					<div><code class="text-purple-600 dark:text-purple-400">query</code> - Query parameters object</div>
					<div><code class="text-purple-600 dark:text-purple-400">headers</code> - Request headers object</div>
					<div><code class="text-purple-600 dark:text-purple-400">body</code> - Request body (string)</div>
				</div>
			</div>

			<!-- Response Object -->
			{#if executionPoint === 'after_request'}
				<div class="bg-white/50 dark:bg-gray-800/50 rounded p-2">
					<div class="font-semibold theme-text-primary mb-1 font-mono">response</div>
					<div class="space-y-1 theme-text-secondary pl-2">
						<div><code class="text-purple-600 dark:text-purple-400">status_code</code> - HTTP status code</div>
						<div><code class="text-purple-600 dark:text-purple-400">headers</code> - Response headers object</div>
						<div><code class="text-purple-600 dark:text-purple-400">body</code> - Response body (string)</div>
					</div>
				</div>
			{/if}

			<!-- Console -->
			<div class="bg-white/50 dark:bg-gray-800/50 rounded p-2">
				<div class="font-semibold theme-text-primary mb-1 font-mono">console.log()</div>
				<div class="theme-text-secondary pl-2">
					Use console.log() for debugging. Logs will appear in test mode.
				</div>
			</div>
		</div>
	</div>

	<!-- Execution Point Notice -->
	<div class="p-3 bg-amber-50/70 dark:bg-amber-900/10 rounded border border-amber-200/50 dark:border-amber-800/50">
		<div class="flex items-start gap-2">
			<i class="fas fa-exclamation-triangle text-amber-500 dark:text-amber-400 text-sm mt-0.5"></i>
			<div class="text-xs">
				<div class="font-semibold theme-text-primary mb-1">Execution Point: {executionPoint === 'before_request' ? 'Before Request' : 'After Request'}</div>
				<div class="theme-text-secondary">
					{#if executionPoint === 'before_request'}
						Your code will run <strong>before</strong> the request is sent. Only the <code class="text-purple-600 dark:text-purple-400">request</code> object is available.
					{:else}
						Your code will run <strong>after</strong> receiving the response. Both <code class="text-purple-600 dark:text-purple-400">request</code> and <code class="text-purple-600 dark:text-purple-400">response</code> objects are available.
					{/if}
				</div>
			</div>
		</div>
	</div>

	<!-- Security Notice -->
	<div class="p-3 bg-gray-50/70 dark:bg-gray-900/50 rounded border theme-border-subtle">
		<div class="flex items-start gap-2">
			<i class="fas fa-shield-alt text-gray-500 dark:text-gray-400 text-sm mt-0.5"></i>
			<div class="text-xs">
				<div class="font-semibold theme-text-primary mb-1">Sandboxed Environment</div>
				<div class="theme-text-secondary space-y-1">
					<div>JavaScript code runs in a secure sandbox with ECMAScript 5.1 support.</div>
					<div><strong>Limitations:</strong> No async/await, arrow functions (=&gt;), let/const, template literals, or other ES6+ features.</div>
					<div>Network requests and file system access are disabled. Scripts timeout after 5 seconds.</div>
				</div>
			</div>
		</div>
	</div>
</div>

<!-- Debug Panel -->
<DebugPanel
	bind:visible={showDebugPanel}
	{script}
	{executionPoint}
/>
