<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { fade } from 'svelte/transition';
	import type { Replay, ExecuteReplayResponse } from '$lib/types/Replay';
	import { replayApi } from '$lib/api/replayApi';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
	import StatusCodeBadge from '$lib/components/common/StatusCodeBadge.svelte';

	export let replay: Replay;

	const dispatch = createEventDispatcher<{
		close: void;
		executed: ExecuteReplayResponse;
	}>();

	let isExecuting = false;
	let executionResult: ExecuteReplayResponse | null = null;
	let showRequestDetails = false;
	let showResponseDetails = false;

	async function executeReplay() {
		if (!$selectedWorkspace || !$selectedProject) {
			toast.error('Please select a workspace and project');
			return;
		}

		try {
			isExecuting = true;
			executionResult = await replayApi.executeReplay(
				$selectedWorkspace.id,
				$selectedProject.id,
				replay.id
			);
			if (executionResult) {
				dispatch('executed', executionResult);
			}
			toast.success('Replay executed successfully');
		} catch (error) {
			toast.error(error);
		} finally {
			isExecuting = false;
		}
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
</script>

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" transition:fade>
	<div class="bg-gray-800 rounded-lg shadow-lg w-full max-w-4xl max-h-[90vh] overflow-hidden">
		<!-- Header -->
		<div class="flex justify-between items-center p-4 border-b border-gray-700">
			<div class="flex items-center space-x-4">
				<h2 class="text-lg font-semibold text-white">Execute Replay</h2>
				<div class="flex items-center space-x-2">
					<HttpMethodBadge method={replay.method} />
					<span class="text-gray-300 font-mono text-sm">{replay.alias}</span>
				</div>
			</div>
			<button
				on:click={() => dispatch('close')}
				class="text-gray-400 hover:text-white transition-colors"
			>
				<i class="fas fa-times"></i>
			</button>
		</div>

		<div class="p-4 max-h-[80vh] overflow-y-auto">
			<!-- Replay Details -->
			<div class="mb-6 p-4 bg-gray-750 rounded-md">
				<h3 class="text-sm font-semibold text-white mb-3 flex items-center">
					<i class="fas fa-info-circle mr-2 text-blue-400"></i>
					Request Details
				</h3>
				<div class="grid grid-cols-2 gap-4 text-sm">
					<div>
						<span class="text-gray-400">URL:</span>
						<span class="text-white font-mono ml-2">{replay.target_url}</span>
					</div>
					<div>
						<span class="text-gray-400">Method:</span>
						<HttpMethodBadge method={replay.method} className="ml-2" />
					</div>
				</div>

				{#if replay.headers && Object.keys(replay.headers).length > 0}
					<div class="mt-3">
						<button
							class="text-xs text-blue-400 hover:text-blue-300 flex items-center"
							on:click={() => (showRequestDetails = !showRequestDetails)}
						>
							<i class="fas fa-chevron-{showRequestDetails ? 'up' : 'down'} mr-1"></i>
							Headers ({Object.keys(replay.headers).length})
						</button>
						{#if showRequestDetails}
							<div class="mt-2 bg-gray-900 rounded p-2 max-h-32 overflow-auto">
								{#each Object.entries(replay.headers) as [key, value]}
									<div class="text-xs font-mono">
										<span class="text-blue-400">{key}:</span>
										<span class="text-white ml-2">{value}</span>
									</div>
								{/each}
							</div>
						{/if}
					</div>
				{/if}

				{#if replay.payload}
					<div class="mt-3">
						<div class="text-xs text-gray-400 mb-1">Request Body:</div>
						<MonacoEditor
							value={replay.payload}
							language="json"
							readOnly={true}
							theme="vs-dark"
						/>
					</div>
				{/if}
			</div>

			<!-- Execute Button -->
			<div class="mb-6 text-center">
				<button
					on:click={executeReplay}
					disabled={isExecuting}
					class="bg-green-600 hover:bg-green-700 disabled:bg-gray-600 text-white px-6 py-2 rounded-md flex items-center mx-auto transition-colors"
				>
					{#if isExecuting}
						<i class="fas fa-spinner fa-spin mr-2"></i>
						Executing...
					{:else}
						<i class="fas fa-play mr-2"></i>
						Execute Replay
					{/if}
				</button>
			</div>

			<!-- Execution Result -->
			{#if executionResult}
				<div class="space-y-4" transition:fade={{ delay: 200 }}>
					<!-- Result Summary -->
					<div class="p-4 bg-gray-750 rounded-md">
						<h3 class="text-sm font-semibold text-white mb-3 flex items-center">
							<i class="fas fa-chart-bar mr-2 text-green-400"></i>
							Execution Result
						</h3>
						<div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
							<div>
								<span class="text-gray-400">Status:</span>
								<StatusCodeBadge statusCode={executionResult.status_code} className="ml-2" />
							</div>
							<div>
								<span class="text-gray-400">Duration:</span>
								<span class="text-white ml-2">{formatDuration(executionResult.duration_ms)}</span>
							</div>
							<div>
								<span class="text-gray-400">Size:</span>
								<span class="text-white ml-2">{executionResult.response_size} bytes</span>
							</div>
							<div>
								<span class="text-gray-400">Executed:</span>
								<span class="text-white ml-2">{formatTimestamp(executionResult.executed_at)}</span>
							</div>
						</div>
					</div>

					<!-- Response Details -->
					<div class="p-4 bg-gray-750 rounded-md">
						<div class="flex justify-between items-center mb-3">
							<h3 class="text-sm font-semibold text-white flex items-center">
								<i class="fas fa-reply mr-2 text-blue-400"></i>
								Response Details
							</h3>
							<button
								class="text-xs text-blue-400 hover:text-blue-300 flex items-center"
								on:click={() => (showResponseDetails = !showResponseDetails)}
							>
								<i class="fas fa-chevron-{showResponseDetails ? 'up' : 'down'} mr-1"></i>
								{showResponseDetails ? 'Hide' : 'Show'} Details
							</button>
						</div>

						{#if showResponseDetails}
							<div class="space-y-3" transition:fade>
								<!-- Response Headers -->
								{#if executionResult.response_headers && Object.keys(executionResult.response_headers).length > 0}
									<div>
										<div class="text-xs text-gray-400 mb-1">Response Headers:</div>
										<div class="bg-gray-900 rounded p-2 max-h-32 overflow-auto">
											{#each Object.entries(executionResult.response_headers) as [key, value]}
												<div class="text-xs font-mono">
													<span class="text-blue-400">{key}:</span>
													<span class="text-white ml-2">{Array.isArray(value) ? value.join(', ') : value}</span>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								<!-- Response Body -->
								{#if executionResult.response_body}
									<div>
										<div class="text-xs text-gray-400 mb-1">Response Body:</div>
										<MonacoEditor
											value={executionResult.response_body}
											language="json"
											height={200}
											readOnly={true}
											theme="vs-dark"
										/>
									</div>
								{/if}
							</div>
						{/if}
					</div>

					<!-- Error Details -->
					{#if executionResult.error_message}
						<div class="p-4 bg-red-900/20 border border-red-600 rounded-md">
							<h3 class="text-sm font-semibold text-red-400 mb-2 flex items-center">
								<i class="fas fa-exclamation-triangle mr-2"></i>
								Error Details
							</h3>
							<pre class="text-xs text-red-300 font-mono whitespace-pre-wrap">{executionResult.error_message}</pre>
						</div>
					{/if}
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
