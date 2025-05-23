<script lang="ts">
	import type { RequestLog } from '$lib/api/BeoApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import HeadersEditor from '../../common/HeadersEditor.svelte';

	export let log: RequestLog;
	export let copyToClipboard: (text: string, label: string) => Promise<void>;
	export let parseJson: (jsonString: string) => any;
</script>

<div>
	<!-- General info -->
	<div class="mb-4 bg-gray-100 dark:bg-gray-850 rounded-md p-3">
		<h3 class="text-sm font-semibold theme-text-secondary mb-2">General</h3>
		<div class="grid grid-cols-2 gap-2 text-sm">
			<div>
				<span class="theme-text-muted">Request URL:</span>
				<span class="theme-text-primary font-mono">{log.path}</span>
			</div>
			<div>
				<span class="theme-text-muted">Method:</span>
				<span class="theme-text-primary font-mono">{log.method}</span>
			</div>
		</div>
	</div>

	<!-- Headers with copy button -->
	<div class="mb-4">
		<div class="flex justify-between items-center mb-2">
			<h3 class="text-sm font-semibold theme-text-secondary">Headers</h3>
			<div class="flex space-x-2">
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() =>
						copyToClipboard(
							JSON.stringify(parseJson(log.request_headers), null, 2),
							'Headers'
						)}
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() =>
						copyToClipboard(
							JSON.stringify(parseJson(log.request_headers)),
							'Headers (minified)'
						)}
				>
					<i class="fas fa-compress-alt mr-1"></i> Minify
				</button>
			</div>
		</div>
		
		<HeadersEditor 
			headers={log.request_headers} 
			editable={false} 
			title="Request Headers" 
		/>
	</div>

	<!-- Request body if exists -->
	{#if log.request_body}
		<div>
			<div class="flex justify-between items-center mb-2">
				<h3 class="text-sm font-semibold theme-text-secondary">Body</h3>
				<div class="flex space-x-2">
					<button
						class={ThemeUtils.utilityButton()}
						on:click|stopPropagation={() =>
							copyToClipboard(
								JSON.stringify(parseJson(log.request_body), null, 2),
								'Body'
							)}
					>
						<i class="fas fa-copy mr-1"></i> Copy
					</button>
					<button
						class={ThemeUtils.utilityButton()}
						on:click|stopPropagation={() =>
							copyToClipboard(
								JSON.stringify(parseJson(log.request_body)),
								'Body (minified)'
							)}
					>
						<i class="fas fa-compress-alt mr-1"></i> Minify
					</button>
				</div>
			</div>
			<pre
				class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-64">{JSON.stringify(
					parseJson(log.request_body),
					null,
					2
				)}</pre>
		</div>
	{/if}

	<!-- Query parameters if exists -->
	{#if log.query_params}
		<div class="mt-4">
			<div class="flex justify-between items-center mb-2">
				<h3 class="text-sm font-semibold theme-text-secondary">Query Parameters</h3>
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() =>
						copyToClipboard(log.query_params, 'Query parameters')}
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
			</div>
			<pre
				class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-32">{log.query_params}</pre>
		</div>
	{/if}
</div>
