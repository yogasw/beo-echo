<script lang="ts">
	import { onMount, createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';
	import type { Replay, ReplayLog } from '$lib/types/Replay';
	import { replayApi } from '$lib/api/replayApi';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
	import StatusCodeBadge from '$lib/components/common/StatusCodeBadge.svelte';

	export let replay: Replay;

	const dispatch = createEventDispatcher<{
		close: void;
	}>();

	let logs: ReplayLog[] = [];
	let isLoading = true;
	let error: string | null = null;
	let expandedLog: string | null = null;
	let searchTerm = '';
	let statusFilter = '';
	let sortOrder: 'asc' | 'desc' = 'desc';

	$: filteredLogs = logs
		.filter((log) => {
			const matchesSearch = !searchTerm || 
				log.id.toLowerCase().includes(searchTerm.toLowerCase()) ||
				log.error_message?.toLowerCase().includes(searchTerm.toLowerCase());
			
			const matchesStatus = !statusFilter || 
				log.status_code.toString() === statusFilter;
			
			return matchesSearch && matchesStatus;
		})
		.sort((a, b) => {
			const dateA = new Date(a.executed_at).getTime();
			const dateB = new Date(b.executed_at).getTime();
			return sortOrder === 'desc' ? dateB - dateA : dateA - dateB;
		});

	onMount(async () => {
		await loadLogs();
	});

	async function loadLogs() {
		if (!$selectedWorkspace || !$selectedProject) {
			error = 'Please select a workspace and project';
			isLoading = false;
			return;
		}

		try {
			isLoading = true;
			error = null;
			logs = await replayApi.getReplayLogs($selectedWorkspace.id, $selectedProject.id, replay.id);
		} catch (err) {
			error = err instanceof Error ? err.message : 'Failed to load replay logs';
			toast.error(err);
		} finally {
			isLoading = false;
		}
	}

	function toggleLogExpansion(logId: string) {
		expandedLog = expandedLog === logId ? null : logId;
	}

	function formatDuration(ms: number): string {
		if (ms < 1000) return `${ms}ms`;
		return `${(ms / 1000).toFixed(2)}s`;
	}

	function formatTimestamp(timestamp: string): string {
		return new Date(timestamp).toLocaleString();
	}

	function getStatusColor(status: number): string {
		if (status >= 200 && status < 300) return 'text-green-400';
		if (status >= 300 && status < 400) return 'text-blue-400';
		if (status >= 400 && status < 500) return 'text-yellow-400';
		if (status >= 500) return 'text-red-400';
		return 'text-gray-400';
	}

	function exportLogs() {
		const dataStr = JSON.stringify(filteredLogs, null, 2);
		const dataBlob = new Blob([dataStr], { type: 'application/json' });
		const url = URL.createObjectURL(dataBlob);
		const link = document.createElement('a');
		link.href = url;
		link.download = `replay-logs-${replay.name}-${new Date().toISOString().split('T')[0]}.json`;
		link.click();
		URL.revokeObjectURL(url);
		toast.success('Logs exported successfully');
	}

	// Get unique status codes for filter dropdown
	$: uniqueStatusCodes = [...new Set(logs.map(log => log.status_code))].sort((a, b) => a - b);
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" transition:fade>
	<div class="bg-gray-800 rounded-lg shadow-lg w-full max-w-6xl max-h-[90vh] overflow-hidden">
		<!-- Header -->
		<div class="flex justify-between items-center p-4 border-b border-gray-700">
			<div class="flex items-center space-x-4">
				<h2 class="text-lg font-semibold text-white">Replay Execution Logs</h2>
				<div class="flex items-center space-x-2">
					<HttpMethodBadge method={replay.method} />
					<span class="text-gray-300 font-mono text-sm">{replay.name}</span>
				</div>
			</div>
			<button
				on:click={() => dispatch('close')}
				class="text-gray-400 hover:text-white transition-colors"
			>
				<i class="fas fa-times"></i>
			</button>
		</div>

		<!-- Filters and Controls -->
		<div class="p-4 border-b border-gray-700 bg-gray-750">
			<div class="flex flex-wrap items-center justify-between gap-4">
				<div class="flex items-center space-x-4 flex-1">
					<!-- Search -->
					<div class="relative flex-1 max-w-md">
						<i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
						<input
							type="text"
							placeholder="Search logs..."
							bind:value={searchTerm}
							class="w-full pl-10 pr-4 py-2 bg-gray-800 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
						/>
					</div>

					<!-- Status Filter -->
					<select
						bind:value={statusFilter}
						class="bg-gray-800 border border-gray-600 text-white rounded-md px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
					>
						<option value="">All Status Codes</option>
						{#each uniqueStatusCodes as statusCode}
							<option value={statusCode}>{statusCode}</option>
						{/each}
					</select>

					<!-- Sort Order -->
					<button
						on:click={() => (sortOrder = sortOrder === 'desc' ? 'asc' : 'desc')}
						class="bg-gray-700 hover:bg-gray-600 text-white px-3 py-2 rounded-md flex items-center transition-colors"
						title="Toggle sort order"
					>
						<i class="fas fa-sort-amount-{sortOrder === 'desc' ? 'down' : 'up'} mr-2"></i>
						{sortOrder === 'desc' ? 'Newest' : 'Oldest'}
					</button>
				</div>

				<div class="flex items-center space-x-2">
					<!-- Export Button -->
					<button
						on:click={exportLogs}
						disabled={filteredLogs.length === 0}
						class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-600 text-white px-3 py-2 rounded-md flex items-center transition-colors"
					>
						<i class="fas fa-download mr-2"></i>
						Export
					</button>

					<!-- Refresh Button -->
					<button
						on:click={loadLogs}
						disabled={isLoading}
						class="bg-gray-700 hover:bg-gray-600 disabled:bg-gray-600 text-white px-3 py-2 rounded-md flex items-center transition-colors"
					>
						<i class="fas fa-sync {isLoading ? 'fa-spin' : ''} mr-2"></i>
						Refresh
					</button>
				</div>
			</div>

			<!-- Stats -->
			<div class="mt-3 flex items-center space-x-6 text-sm text-gray-400">
				<span>Total: {logs.length} executions</span>
				<span>Filtered: {filteredLogs.length} results</span>
				{#if logs.length > 0}
					<span>Success Rate: {Math.round((logs.filter(l => l.status_code >= 200 && l.status_code < 400).length / logs.length) * 100)}%</span>
				{/if}
			</div>
		</div>

		<div class="flex-1 overflow-hidden">
			{#if isLoading}
				<div class="p-4">
					<SkeletonLoader type="list" count={8} />
				</div>
			{:else if error}
				<div class="p-4">
					<ErrorDisplay 
						message={error} 
						type="error" 
						retryable={true}
						onRetry={loadLogs}
					/>
				</div>
			{:else if filteredLogs.length === 0}
				<div class="p-8 text-center">
					<i class="fas fa-history text-4xl text-gray-600 mb-4"></i>
					<h3 class="text-lg text-gray-400 mb-2">No Execution Logs</h3>
					<p class="text-gray-500">
						{logs.length === 0 
							? 'This replay hasn\'t been executed yet.' 
							: 'No logs match your current filters.'}
					</p>
				</div>
			{:else}
				<div class="max-h-[60vh] overflow-y-auto">
					<div class="space-y-2 p-4">
						{#each filteredLogs as log (log.id)}
							<div class="bg-gray-750 border border-gray-700 rounded-md overflow-hidden" transition:fade={{ duration: 200 }}>
								<!-- Log Header -->
								<button
									on:click={() => toggleLogExpansion(log.id)}
									class="w-full p-3 flex items-center justify-between hover:bg-gray-700 transition-colors"
								>
									<div class="flex items-center space-x-4 flex-1">
										<StatusCodeBadge statusCode={log.status_code} />
										<span class="text-sm text-gray-300 font-mono">{formatTimestamp(log.executed_at)}</span>
										<span class="text-sm text-white">{formatDuration(log.duration_ms)}</span>
										<span class="text-xs text-gray-400">{log.response_size} bytes</span>
										{#if log.error_message}
											<i class="fas fa-exclamation-triangle text-red-400" title="Error occurred"></i>
										{/if}
									</div>
									<i class="fas fa-chevron-{expandedLog === log.id ? 'up' : 'down'} text-gray-400"></i>
								</button>

								<!-- Expanded Log Details -->
								{#if expandedLog === log.id}
									<div class="border-t border-gray-700 p-4 space-y-4" transition:fade={{ duration: 200 }}>
										<!-- Request Details -->
										<div>
											<h4 class="text-sm font-semibold text-white mb-2 flex items-center">
												<i class="fas fa-arrow-up mr-2 text-blue-400"></i>
												Request
											</h4>
											<div class="bg-gray-900 rounded p-3 text-xs font-mono">
												<div class="text-blue-400">{replay.method} {replay.url}</div>
												{#if replay.headers && Object.keys(replay.headers).length > 0}
													<div class="mt-2 text-gray-400">Headers:</div>
													{#each Object.entries(replay.headers) as [key, value]}
														<div class="text-gray-300">{key}: {value}</div>
													{/each}
												{/if}
												{#if replay.body}
													<div class="mt-2 text-gray-400">Body:</div>
													<div class="text-white">{replay.body}</div>
												{/if}
											</div>
										</div>

										<!-- Response Details -->
										<div>
											<h4 class="text-sm font-semibold text-white mb-2 flex items-center">
												<i class="fas fa-arrow-down mr-2 text-green-400"></i>
												Response
											</h4>
											<div class="space-y-2">
												{#if log.response_headers && Object.keys(log.response_headers).length > 0}
													<div>
														<div class="text-xs text-gray-400 mb-1">Headers:</div>
														<div class="bg-gray-900 rounded p-2 max-h-32 overflow-auto">
															{#each Object.entries(log.response_headers) as [key, value]}
																<div class="text-xs font-mono">
																	<span class="text-blue-400">{key}:</span>
																	<span class="text-white ml-2">{Array.isArray(value) ? value.join(', ') : value}</span>
																</div>
															{/each}
														</div>
													</div>
												{/if}

												{#if log.response_body}
													<div>
														<div class="text-xs text-gray-400 mb-1">Body:</div>
														<MonacoEditor
															value={log.response_body}
															language="json"
															height={150}
															readOnly={true}
															theme="vs-dark"
														/>
													</div>
												{/if}
											</div>
										</div>

										<!-- Error Details -->
										{#if log.error_message}
											<div>
												<h4 class="text-sm font-semibold text-red-400 mb-2 flex items-center">
													<i class="fas fa-exclamation-triangle mr-2"></i>
													Error
												</h4>
												<div class="bg-red-900/20 border border-red-600 rounded p-3">
													<pre class="text-xs text-red-300 font-mono whitespace-pre-wrap">{log.error_message}</pre>
												</div>
											</div>
										{/if}
									</div>
								{/if}
							</div>
						{/each}
					</div>
				</div>
			{/if}
		</div>

		<!-- Footer -->
		<div class="flex justify-end space-x-3 p-4 border-t border-gray-700">
			<button
				on:click={() => dispatch('close')}
				class="bg-gray-700 hover:bg-gray-600 text-white px-4 py-2 rounded text-sm"
			>
				Close
			</button>
		</div>
	</div>
</div>
