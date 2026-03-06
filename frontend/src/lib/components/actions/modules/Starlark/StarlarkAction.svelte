<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import DebugPanel from './DebugPanel.svelte';
	import type { StarlarkConfig } from '$lib/types/Action';

	export let config: StarlarkConfig | null = null;
	export let executionPoint: 'before_request' | 'after_request' = 'after_request';

	const dispatch = createEventDispatcher();

	let showDebugPanel = false;

	// Initialize script with default template
	let script =
		config?.script ||
		`# Modify request or response using Starlark (Python-like syntax)
# Available: request, response (if after_request), json module
# Use print() to debug

# Example: Add a custom header
if request:
	request["headers"]["X-Custom-Header"] = "MyValue"

${executionPoint === 'after_request' ? `# Example: Modify response body
if response:
	body = json.decode(response["body"])
	body["modified"] = True
	body["starlark"] = "processed"
	response["body"] = json.encode(body)
` : ''}
`;

	let editor: any;

	// Update parent when script changes
	function handleScriptChange(event: CustomEvent) {
		script = event.detail;
		const newConfig: StarlarkConfig = { script };
		dispatch('change', newConfig);
	}

	// Update template when execution point changes
	$: {
		if (!config) {
			// Only update template for new actions
			if (executionPoint === 'after_request' && !script.includes('response')) {
				script += `\n\n# Example: Modify response
if response:
	print("Response status:", response["status_code"])
`;
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
		const initialConfig: StarlarkConfig = { script };
		dispatch('change', initialConfig);
	});
</script>

<div class="space-y-4">
	<!-- Starlark Editor -->
	<div>
		<div class="flex items-center justify-between mb-2">
			<label for="script-editor" class="block text-sm font-medium theme-text-primary">
				Starlark Code <span class="text-red-500">*</span>
			</label>
			<div class="flex items-center gap-2">
				<button
					type="button"
					on:click={() => showDebugPanel = true}
					class="text-xs px-3 py-1.5 rounded bg-green-600 hover:bg-green-700 text-white transition-colors flex items-center gap-1.5"
					title="Test script with sample log data"
					aria-label="Test Starlark code"
				>
					<i class="fas fa-play text-xs"></i>
					Test Script
				</button>
				<button
					type="button"
					on:click={formatCode}
					class="text-xs px-3 py-1.5 rounded bg-gray-700 hover:bg-gray-600 text-white transition-colors flex items-center gap-1.5"
					title="Format code"
					aria-label="Format Starlark code"
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
				language="python"
				on:change={handleScriptChange}
			/>
		</div>

		<p id="script-help" class="mt-2 text-xs theme-text-secondary">
			Write Starlark (Python-like) code to modify the request or response
		</p>
	</div>

	<!-- Context Information -->
	<div class="p-4 bg-blue-50/70 dark:bg-gray-900/50 rounded-lg border border-blue-200/50 dark:border-blue-900/50">
		<div class="flex items-start gap-2 mb-3">
			<i class="fas fa-info-circle text-blue-500 mt-0.5"></i>
			<div>
				<h4 class="text-sm font-semibold theme-text-primary">Available Context</h4>
				<p class="text-xs theme-text-secondary mt-1">
					The following objects are available in your Starlark code:
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
					<div><code class="text-purple-600 dark:text-purple-400">query</code> - Query parameters dict</div>
					<div><code class="text-purple-600 dark:text-purple-400">headers</code> - Request headers dict</div>
					<div><code class="text-purple-600 dark:text-purple-400">body</code> - Request body (string)</div>
				</div>
			</div>

			<!-- Response Object -->
			{#if executionPoint === 'after_request'}
				<div class="bg-white/50 dark:bg-gray-800/50 rounded p-2">
					<div class="font-semibold theme-text-primary mb-1 font-mono">response</div>
					<div class="space-y-1 theme-text-secondary pl-2">
						<div><code class="text-purple-600 dark:text-purple-400">status_code</code> - HTTP status code</div>
						<div><code class="text-purple-600 dark:text-purple-400">headers</code> - Response headers dict</div>
						<div><code class="text-purple-600 dark:text-purple-400">body</code> - Response body (string)</div>
					</div>
				</div>
			{/if}

			<!-- Built-in Modules -->
			<div class="bg-white/50 dark:bg-gray-800/50 rounded p-2">
				<div class="font-semibold theme-text-primary mb-1 font-mono">Built-in Modules</div>
				<div class="space-y-1 theme-text-secondary pl-2">
					<div><code class="text-green-600 dark:text-green-400">json</code> - JSON encoding/decoding (json.encode, json.decode)</div>
					<div><code class="text-green-600 dark:text-green-400">print()</code> - Print for debugging (visible in test mode)</div>
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
				<div class="font-semibold theme-text-primary mb-1">Secure Sandbox Environment</div>
				<div class="theme-text-secondary space-y-1">
					<div>Starlark is a Python-like language designed for safe, deterministic execution.</div>
					<div><strong>Features:</strong> Supports functions, loops, conditionals, and built-in data types (list, dict, string).</div>
					<div><strong>Built-ins:</strong> JSON module for parsing/encoding, http module for requests (sandboxed).</div>
					<div><strong>Security:</strong> No file system access, no arbitrary imports. Scripts timeout after 5 seconds.</div>
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
