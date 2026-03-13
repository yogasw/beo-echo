<script lang="ts">
	import type { ExecuteReplayResponse } from '$lib/types/Replay';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';

	let { executionResult }: { executionResult: ExecuteReplayResponse } = $props();

	let format = $state('json'); // 'json', 'xml', 'html', 'javascript', 'plaintext'
	let isPreview = $state(false);

	// Try to auto-detect content type when executionResult changes
	$effect(() => {
		if (executionResult?.response_headers) {
			const contentType = Object.keys(executionResult.response_headers)
				.find(k => k.toLowerCase() === 'content-type');
			if (contentType) {
				const val = executionResult.response_headers[contentType].toLowerCase();
				if (val.includes('application/json')) {
					format = 'json';
				} else if (val.includes('text/html')) {
					format = 'html';
				} else if (val.includes('text/xml') || val.includes('application/xml')) {
					format = 'xml';
				} else if (val.includes('javascript')) {
					format = 'javascript';
				} else if (val.includes('text/plain')) {
					format = 'plaintext';
				}
			}
		}
	});

	let formattedBody = $derived.by(() => {
		if (!executionResult?.response_body) return '';
		
		const body = executionResult.response_body;
		if (format === 'json') {
			try {
				const ob = typeof body === 'string' ? JSON.parse(body) : body;
				return JSON.stringify(ob, null, 2);
			} catch(e) {
				return typeof body === 'string' ? body : JSON.stringify(body, null, 2);
			}
		}
		return typeof body === 'string' ? body : JSON.stringify(body, null, 2);
	});

	let previewUrl = $derived.by(() => {
		if (!isPreview || !executionResult?.response_body) return '';
		const body = executionResult.response_body;
		const content = typeof body === 'string' ? body : JSON.stringify(body);
		const blob = new Blob([content], { type: 'text/html' });
		return URL.createObjectURL(blob);
	});
</script>

{#if executionResult.error}
	<div class="h-full w-full p-4 overflow-auto">
		<div class="p-4 bg-red-100 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md text-red-800 dark:text-red-300">
			<h3 class="font-semibold">Error</h3>
			<p>{executionResult.error}</p>
		</div>
	</div>
{:else if executionResult.response_body}
	<div class="flex flex-col h-full w-full">
		<!-- Toolbar -->
		<div class="flex flex-row items-center gap-2 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-900 px-4 py-2">
			<!-- Format Dropdown -->
			<div class="relative flex items-center bg-gray-200 dark:bg-gray-800 hover:bg-gray-300 dark:hover:bg-gray-700 rounded transition-colors duration-150">
				<select
					bind:value={format}
					class="appearance-none bg-transparent text-gray-700 dark:text-gray-300 py-1 pl-3 pr-8 rounded text-sm focus:outline-none cursor-pointer"
				>
					<option value="json">JSON</option>
					<option value="xml">XML</option>
					<option value="html">HTML</option>
					<option value="javascript">JavaScript</option>
					<option value="plaintext">Raw</option>
				</select>
				<div class="pointer-events-none absolute right-0 flex items-center pr-2 text-gray-500">
					<i class="fas fa-chevron-down text-[10px]"></i>
				</div>
			</div>

			<!-- Preview Button -->
			<button
				onclick={() => {
					isPreview = !isPreview;
				}}
				class="flex items-center gap-2 px-3 py-1 rounded text-sm transition-colors {isPreview ? 'bg-orange-100 dark:bg-orange-900/20 text-orange-600 dark:text-orange-400' : 'bg-transparent text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-800 hover:text-gray-800 dark:hover:text-gray-200'}"
			>
				<i class="fas fa-play text-xs {isPreview ? 'text-orange-500 text-opacity-80' : ''}"></i>
				<span>Preview</span>
			</button>
		</div>

		<!-- Content Area -->
		<div class="flex-grow relative w-full h-full overflow-hidden">
			{#if isPreview}
				<iframe
					src={previewUrl}
					class="w-full h-full bg-white border-0"
					title="HTML Preview"
					sandbox="allow-scripts"
				></iframe>
			{:else}
				<MonacoEditor
					value={formattedBody}
					language={format}
					readOnly={true}
				/>
			{/if}
		</div>
	</div>
{:else}
	<div class="h-full w-full p-4 flex items-center justify-center">
		<div class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300">
			No response body
		</div>
	</div>
{/if}
