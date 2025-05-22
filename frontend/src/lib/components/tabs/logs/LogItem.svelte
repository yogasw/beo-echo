<script lang="ts">
	import { fade } from 'svelte/transition';
	import type { RequestLog } from '$lib/api/BeoApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import LogRequestContent from './LogRequestContent.svelte';
	import LogResponseContent from './LogResponseContent.svelte';

	export let log: RequestLog;
	export let isExpanded: boolean = false;
	export let activeTab: 'request' | 'response' = 'request';
	export let toggleLogExpansion: (logId: string) => void;
	export let switchTab: (logId: string, tab: 'request' | 'response') => void;
	export let copyToClipboard: (text: string, label: string) => Promise<void>;
	export let parseJson: (jsonString: string) => any;
	export let formatDate: (dateString: string | Date) => string;
	export let bookmarkLog: (log: RequestLog) => Promise<void>;
	export let createMockFromLog: (log: RequestLog) => void;
</script>

<div class={ThemeUtils.card('overflow-hidden')}>
	<!-- Log header - clickable to expand/collapse -->
	<div class="flex flex-col">
		<div
			class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer"
			on:click={() => toggleLogExpansion(log.id)}
			on:keydown={(e) => e.key === 'Enter' && toggleLogExpansion(log.id)}
			tabindex="0"
			role="button"
			aria-expanded={isExpanded}
		>
			<div class="flex items-center space-x-2">
				<!-- Method badge -->
				<span
					class="px-2 py-0.5 text-sm font-mono rounded {log.method === 'GET'
						? 'bg-green-600 text-white'
						: log.method === 'POST'
							? 'bg-blue-600 text-white'
							: log.method === 'PUT'
								? 'bg-yellow-600 text-white'
								: log.method === 'DELETE'
									? 'bg-red-600 text-white'
									: 'bg-gray-600 text-white'}"
				>
					{log.method}
				</span>

				<!-- Path with truncation -->
				<span class="font-mono text-sm theme-text-primary truncate max-w-sm">
					{log.path}
				</span>

				<!-- Status code -->
				<span
					class="px-2 py-0.5 text-xs font-mono rounded {log.response_status < 300
						? 'bg-green-600 text-white'
						: log.response_status < 400
							? 'bg-blue-600 text-white'
							: log.response_status < 500
								? 'bg-yellow-600 text-white'
								: 'bg-red-600 text-white'}"
				>
					{log.response_status}
				</span>

				<!-- Match status badge -->
				<span
					class="px-2 py-0.5 text-xs font-mono rounded {log.matched
						? 'bg-green-600 text-white'
						: 'bg-red-600 text-white'}"
				>
					{log.matched ? 'Matched' : 'Unmatched'}
				</span>
				
				<!-- Execution Mode badge (proxy/forwarder) -->
				{#if log.execution_mode === 'proxy' || log.execution_mode === 'forwarder'}
					<span class="px-2 py-0.5 text-xs font-mono rounded bg-purple-600 text-white">
						{log.execution_mode === 'proxy' ? 'Proxy' : 'Forwarder'}
					</span>
				{/if}
				
			</div>

			<div class="flex items-center space-x-3">
				<span class="text-xs theme-text-muted">{formatDate(log.created_at)}</span>
				<span class="px-2 py-0.5 text-xs bg-blue-600 rounded text-white"
					>{log.latency_ms}ms</span
				>
				<i
					class="fas {isExpanded ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>
		</div>

		<!-- Action buttons row -->
		<div class="flex justify-end px-3 py-1 border-t theme-border-light">
			<!-- Bookmark button -->
			<button
				class="{log.bookmark ? 'bg-yellow-600 hover:bg-yellow-700' : 'bg-gray-600 hover:bg-gray-700'} text-white py-1 px-3 rounded text-xs flex items-center transition-all duration-200 transform hover:scale-105 mr-2"
				on:click|stopPropagation={() => bookmarkLog(log)}
				title={log.bookmark ? "Remove from bookmarks" : "Add to bookmarks"}
			>
				<i class="fas {log.bookmark ? 'fa-bookmark' : 'fa-bookmark'} mr-1"></i> 
				{log.bookmark ? 'Bookmarked' : 'Bookmark'}
			</button>
			
			<!-- Create Mock button - only for unmatched requests -->
			{#if !log.matched}
				<button
					class="bg-emerald-600 hover:bg-emerald-700 text-white py-1 px-3 rounded text-xs flex items-center transition-all duration-200 transform hover:scale-105"
					on:click|stopPropagation={() => createMockFromLog(log)}
					title="Create a new mock endpoint using the data from this request"
				>
					<i class="fas fa-magic mr-1"></i> Create Mock from this Request
				</button>
			{/if}
		</div>
	</div>

	<!-- Expanded details -->
	{#if isExpanded}
		<div transition:fade={{ duration: 150 }} class="border-t theme-border px-4 py-3">
			<!-- Tab navigation -->
			<div class="flex mb-4 border-b theme-border">
				<button
					class="px-4 py-2 font-medium text-sm {activeTab === 'request'
						? 'text-blue-500 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400'
						: 'theme-text-muted hover:text-gray-600 dark:hover:text-gray-300'}"
					on:click|stopPropagation={() => switchTab(log.id, 'request')}
				>
					Request
				</button>
				<button
					class="px-4 py-2 font-medium text-sm {activeTab === 'response'
						? 'text-blue-500 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400'
						: 'theme-text-muted hover:text-gray-600 dark:hover:text-gray-300'}"
					on:click|stopPropagation={() => switchTab(log.id, 'response')}
				>
					Response
				</button>
			</div>

			<!-- Request content -->
			{#if activeTab === 'request'}
				<LogRequestContent {log} {copyToClipboard} {parseJson} />
			{:else}
				<LogResponseContent {log} {copyToClipboard} {parseJson} />
			{/if}
		</div>
	{/if}
</div>
