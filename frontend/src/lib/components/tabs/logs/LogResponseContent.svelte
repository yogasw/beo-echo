<script lang="ts">
	import type { RequestLog } from '$lib/api/BeoApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import HeadersTab from '../../common/HeadersEditor.svelte';
	import StatusCodeBadge from '$lib/components/common/StatusCodeBadge.svelte';

	export let log: RequestLog;
	export let copyToClipboard: (text: string, label: string) => Promise<void>;
	export let parseJson: (jsonString: string) => any;
	let hideHeader: boolean = false;
</script>

<div>
	<!-- General info -->
	<div class="mb-4 bg-gray-100 dark:bg-gray-850 rounded-md p-3">
		<h3 class="text-sm font-semibold theme-text-secondary mb-2">General</h3>
		<div class="grid grid-cols-2 gap-2 text-sm">
			<div>
				<span class="theme-text-muted">Status Code:</span>
				<StatusCodeBadge
					statusCode={log.response_status}
					size="sm"
					showDescription={true}
					className="inline-block"
				/>
			</div>
			<div>
				<span class="theme-text-muted">Execution Mode:</span>
				<span
					class="{log.execution_mode === 'proxy' || log.execution_mode === 'forwarder'
						? 'text-purple-600 dark:text-purple-400'
						: 'theme-text-primary'} font-mono"
				>
					{#if log.execution_mode === 'proxy'}
						Proxy (Forwarded Request)
					{:else if log.execution_mode === 'forwarder'}
						Forwarder
					{:else}
						{log.execution_mode}
					{/if}
				</span>
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
						copyToClipboard(JSON.stringify(parseJson(log.response_headers), null, 2), 'Headers')}
					aria-label="Copy response headers to clipboard"
					title="Copy response headers to clipboard"
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>

				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() => (hideHeader = !hideHeader)}
					aria-label={hideHeader ? 'Show response headers' : 'Hide response headers'}
					title={hideHeader ? 'Show response headers' : 'Hide response headers'}
				>
					<i class="fas {hideHeader ? 'fa-eye' : 'fa-eye-slash'} mr-1"></i>
					{hideHeader ? 'Show' : 'Hide'}
				</button>
			</div>
		</div>
		{#if !hideHeader}
			<HeadersTab headers={log.response_headers} editable={false} title="Response Headers" />
		{/if}
	</div>

	<!-- Response body -->
	<div>
		<div class="flex justify-between items-center mb-2">
			<h3 class="text-sm font-semibold text-gray-300">Body</h3>
			<div class="flex space-x-2">
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() =>
						copyToClipboard(JSON.stringify(parseJson(log.response_body), null, 2), 'Body')}
					aria-label="Copy response body to clipboard"
					title="Copy response body to clipboard"
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
			</div>
		</div>

		<!-- Special handling for endpoint not found error -->
		{#if log.response_status >= 400 && parseJson(log.response_body)?.error === true}
			<div
				class="bg-red-100/30 dark:bg-red-900/30 border border-red-300 dark:border-red-700 p-3 rounded-md"
			>
				<div class="flex items-center">
					<i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-2"></i>
					<span class="text-sm theme-text-primary">
						{parseJson(log.response_body)?.message || 'Error'}
					</span>
				</div>
			</div>
		{:else}
			<pre
				class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-64">{JSON.stringify(
					parseJson(log.response_body),
					null,
					2
				)}</pre>
		{/if}
	</div>
</div>
