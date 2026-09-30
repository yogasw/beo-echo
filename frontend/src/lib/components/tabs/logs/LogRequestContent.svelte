<script lang="ts">
	import type { RequestLog } from '$lib/api/BeoApi';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import HeadersEditor from '../../common/HeadersEditor.svelte';
	import HighlightText from '$lib/components/common/HighlightText.svelte';
	import { highlightTermsFor, type ParsedSearch } from '$lib/utils/logSearch';

	export let log: RequestLog;
	export let copyToClipboard: (text: string, label: string) => Promise<void>;
	export let parseJson: (jsonString: string) => any;
	export let search: ParsedSearch | null = null;
	let hideHeader: boolean;
</script>

<div>
	<!-- General info -->
	<div class="mb-4 bg-gray-100 dark:bg-gray-850 rounded-md p-3">
		<h3 class="text-sm font-semibold theme-text-secondary mb-2">General</h3>
		<div class="grid grid-cols-2 gap-2 text-sm">
			<div>
				<span class="theme-text-muted">Request URL:</span>
				<span class="theme-text-primary font-mono"><HighlightText text={log.path} terms={highlightTermsFor(search, 'path')} /></span>
			</div>
			<div>
				<span class="theme-text-muted">Method:</span>
				<HttpMethodBadge method={log.method} size="sm" />
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
						copyToClipboard(JSON.stringify(parseJson(log.request_headers), null, 2), 'Headers')}
					aria-label="Copy request headers to clipboard"
					title="Copy request headers to clipboard"
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() => (hideHeader = !hideHeader)}
					aria-label={hideHeader ? "Show request headers" : "Hide request headers"}
					title={hideHeader ? "Show request headers" : "Hide request headers"}
				>
					<i class="fas {hideHeader ? 'fa-eye' : 'fa-eye-slash'} mr-1"></i> 
					{hideHeader ? 'Show' : 'Hide'}
				</button>
			</div>
		</div>

		{#if !hideHeader}
			<HeadersEditor
				headers={log.request_headers}
				editable={false}
				title="Request Headers"
				highlightTerms={highlightTermsFor(search, 'request_headers')}
			/>
		{/if}
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
							copyToClipboard(JSON.stringify(parseJson(log.request_body), null, 2), 'Body')}
						aria-label="Copy request body to clipboard"
						title="Copy request body to clipboard"
					>
						<i class="fas fa-copy mr-1"></i> Copy
					</button>
				</div>
			</div>
			<pre
				class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-64"><HighlightText
					text={JSON.stringify(parseJson(log.request_body), null, 2)}
					terms={highlightTermsFor(search, 'request_body')}
				/></pre>
		</div>
	{/if}

	<!-- Query parameters if exists -->
	{#if log.query_params}
		<div class="mt-4">
			<div class="flex justify-between items-center mb-2">
				<h3 class="text-sm font-semibold theme-text-secondary">Query Parameters</h3>
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() => copyToClipboard(log.query_params, 'Query parameters')}
					aria-label="Copy query parameters to clipboard"
					title="Copy query parameters to clipboard"
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
			</div>
			<pre
				class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-32"><HighlightText
					text={log.query_params}
					terms={highlightTermsFor(search, 'query')}
				/></pre>
		</div>
	{/if}
</div>
