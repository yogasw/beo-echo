<script lang="ts">
	import { fade } from 'svelte/transition';
	import type { Project, RequestLog } from '$lib/api/BeoApi';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import LogItem from './LogItem.svelte';

	export let filteredLogs: RequestLog[] = [];
	export let expandedLogs: Record<string, boolean> = {};
	export let activeTabs: Record<string, 'request' | 'response'> = {};
	export let selectedProject: Project;
	export let logsConnectionStatus: any;
	export let toggleLogExpansion: (logId: string) => void;
	export let switchTab: (logId: string, tab: 'request' | 'response') => void;
	export let copyToClipboard: (text: string, label: string) => Promise<void>;
	export let parseJson: (jsonString: string) => any;
	export let formatDate: (dateString: string | Date) => string;
	export let bookmarkLog: (log: RequestLog) => Promise<void>;
	export let createMockFromLog: (log: RequestLog) => void;
	export let replayLog: (log: RequestLog) => void;
</script>

{#if filteredLogs.length === 0}
	<div class="theme-bg-secondary border theme-border p-6 rounded-lg flex flex-col items-center justify-center text-center">
		<!-- Empty logs state -->
		<div class="bg-blue-600/10 dark:bg-blue-600/20 p-5 rounded-full mb-5">
			<i class="fas fa-satellite-dish text-blue-500 text-4xl"></i>
		</div>
		<h3 class="text-xl font-semibold theme-text-primary mb-2">Waiting for requests</h3>
		<p class="theme-text-secondary mb-5 max-w-lg">
			Your endpoint is ready! Send your first HTTP request to start populating your logs.
		</p>
		
		<!-- Example request section -->
		<div class="w-full max-w-2xl mb-6">
			<div class="flex items-center justify-between mb-2">
				<h4 class="text-sm font-medium theme-text-primary">Example cURL request</h4>
				<button
					class={ThemeUtils.utilityButton()}
					on:click|stopPropagation={() => copyToClipboard(`curl -X GET "${selectedProject.url}"`, 'Example request')}
				>
					<i class="fas fa-copy mr-1"></i> Copy
				</button>
			</div>
			<div class="theme-bg-tertiary border theme-border rounded-md p-3">
				<pre class="font-mono text-xs theme-text-secondary overflow-x-auto">curl -X GET {selectedProject.url}</pre>
			</div>
		</div>
		
		<div class="w-full max-w-2xl grid gap-4 md:grid-cols-2">
			<!-- Live status -->
			<div class="theme-bg-tertiary border theme-border rounded-md p-4 flex items-center">
				<span class="relative flex h-3 w-3 mr-3">
					{#if logsConnectionStatus.isConnected}
						<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
						<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
					{:else}
						<span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
					{/if}
				</span>
				<div>
					<h4 class="text-sm font-medium theme-text-primary">Live connection</h4>
					<p class="text-xs theme-text-muted">
						{logsConnectionStatus.isConnected ? 'Connected and ready for requests' : 'Currently offline'}
					</p>
				</div>
			</div>
			
			<!-- Create mock CTA -->
			<div class="theme-bg-tertiary border theme-border rounded-md p-4">
				<h4 class="text-sm font-medium theme-text-primary mb-1 flex items-center">
					<i class="fas fa-magic text-blue-500 mr-2"></i>
					Need mock responses?
				</h4>
				<p class="text-xs theme-text-muted mb-2">
					Create mock endpoints from the Configuration tab
				</p>
				<a href="/configuration" class="text-xs text-blue-500 hover:text-blue-400 underline">
					Go to Configuration
				</a>
			</div>
		</div>
	</div>
{:else}
	<div class="space-y-4">
		{#each filteredLogs as log (log.id)}
			<LogItem 
				{log}
				{selectedProject}
				isExpanded={!!expandedLogs[log.id]}
				activeTab={activeTabs[log.id] || 'request'}
				{toggleLogExpansion}
				{switchTab}
				{copyToClipboard}
				{parseJson}
				{formatDate}
				{bookmarkLog}
				{createMockFromLog}
				{replayLog}
			/>
		{/each}
	</div>

	{#if filteredLogs.length < logsConnectionStatus.total}
		<div class="mt-4 text-center">
			<span class="text-xs text-gray-400">Showing {filteredLogs.length} of {logsConnectionStatus.total} logs</span>
		</div>
	{/if}
{/if}
