<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { replays, filteredReplays, replayFilter, replayActions } from '$lib/stores/replay';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import { HTTP_METHODS, PROTOCOLS } from '$lib/types/Replay';
	import type { Replay } from '$lib/types/Replay';

	const dispatch = createEventDispatcher();

	let searchTerm = '';
	let methodFilter = '';
	let protocolFilter = '';

	// Update store when filters change
	$: replayFilter.set({
		searchTerm,
		method: methodFilter,
		protocol: protocolFilter
	});

	async function handleDelete(replay: Replay) {
		if (!$selectedWorkspace || !$selectedProject) return;

		if (!confirm(`Are you sure you want to delete the replay "${replay.alias}"?`)) {
			return;
		}

		try {
			replayActions.setLoading('delete', true);
			await replayApi.deleteReplay($selectedWorkspace.id, $selectedProject.id, replay.id);
			replayActions.removeReplay(replay.id);
			toast.success('Replay deleted successfully');
		} catch (err: any) {
			toast.error(err);
		} finally {
			replayActions.setLoading('delete', false);
		}
	}

	function handleEdit(replay: Replay) {
		dispatch('edit', replay);
	}

	function handleExecute(replay: Replay) {
		dispatch('execute', replay);
	}

	function handleLogs(replay: Replay) {
		dispatch('logs', replay);
	}

	function handleRefresh() {
		dispatch('refresh');
	}

	function clearFilters() {
		searchTerm = '';
		methodFilter = '';
		protocolFilter = '';
	}

	function getMethodColor(method: string): string {
		switch (method.toUpperCase()) {
			case 'GET': return 'bg-green-600 text-white';
			case 'POST': return 'bg-blue-600 text-white';
			case 'PUT': return 'bg-yellow-600 text-white';
			case 'PATCH': return 'bg-orange-600 text-white';
			case 'DELETE': return 'bg-red-600 text-white';
			case 'HEAD': return 'bg-purple-600 text-white';
			case 'OPTIONS': return 'bg-gray-600 text-white';
			default: return 'bg-gray-600 text-white';
		}
	}

	function formatDate(dateString: string): string {
		return new Date(dateString).toLocaleDateString('en-US', {
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}
</script>

<div class="flex flex-col h-full">
	<!-- Search and Filter Bar -->
	<div class="theme-bg-secondary border-b theme-border p-4">
		<div class="flex flex-col lg:flex-row gap-4">
			<!-- Search Input -->
			<div class="flex-1">
				<div class="relative">
					<div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
						<i class="fas fa-search theme-text-secondary"></i>
					</div>
					<input
						type="text"
						bind:value={searchTerm}
						placeholder="Search replays by alias or URL..."
						class="block w-full pl-10 pr-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
					/>
				</div>
			</div>

			<!-- Filters -->
			<div class="flex gap-3">
				<select
					bind:value={methodFilter}
					class="px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">All Methods</option>
					{#each HTTP_METHODS as method}
						<option value={method}>{method}</option>
					{/each}
				</select>

				<select
					bind:value={protocolFilter}
					class="px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary focus:outline-none focus:ring-2 focus:ring-blue-500"
				>
					<option value="">All Protocols</option>
					{#each PROTOCOLS as protocol}
						<option value={protocol}>{protocol.toUpperCase()}</option>
					{/each}
				</select>

				<button
					on:click={clearFilters}
					class="px-3 py-2 text-sm theme-text-secondary hover:theme-text-primary transition-colors"
					title="Clear filters"
				>
					<i class="fas fa-times"></i>
				</button>

				<button
					on:click={handleRefresh}
					class="px-3 py-2 text-sm theme-text-secondary hover:theme-text-primary transition-colors"
					title="Refresh list"
				>
					<i class="fas fa-sync-alt"></i>
				</button>
			</div>
		</div>

		<!-- Results Summary -->
		{#if searchTerm || methodFilter || protocolFilter}
			<div class="mt-3 text-sm theme-text-secondary">
				Showing {$filteredReplays.length} of {$replays.length} replays
				{#if searchTerm || methodFilter || protocolFilter}
					<button
						on:click={clearFilters}
						class="ml-2 text-blue-400 hover:text-blue-300 underline"
					>
						Clear filters
					</button>
				{/if}
			</div>
		{/if}
	</div>

	<!-- Replay List -->
	<div class="flex-1 overflow-auto">
		{#if $filteredReplays.length === 0}
			<div class="flex items-center justify-center h-full">
				<div class="text-center theme-text-secondary">
					{#if $replays.length === 0}
						<i class="fas fa-play-circle text-4xl mb-4 opacity-50"></i>
						<h3 class="text-lg font-medium mb-2">No replays yet</h3>
						<p>Create your first API replay to get started</p>
					{:else}
						<i class="fas fa-search text-4xl mb-4 opacity-50"></i>
						<h3 class="text-lg font-medium mb-2">No matching replays</h3>
						<p>Try adjusting your search criteria</p>
					{/if}
				</div>
			</div>
		{:else}
			<div class="divide-y theme-border">
				{#each $filteredReplays as replay (replay.id)}
					<div class="p-4 hover:theme-bg-secondary transition-colors group">
						<div class="flex items-center justify-between">
							<!-- Replay Info -->
							<div class="flex-1 min-w-0">
								<div class="flex items-center space-x-3 mb-2">
									<!-- Method Badge -->
									<span class="px-2 py-1 text-xs font-medium rounded {getMethodColor(replay.method)}">
										{replay.method}
									</span>
									
									<!-- Protocol Badge -->
									<span class="px-2 py-1 text-xs font-medium rounded bg-gray-600 text-white">
										{replay.protocol.toUpperCase()}
									</span>
									
									<!-- Alias -->
									<h3 class="font-medium theme-text-primary truncate">
										{replay.alias}
									</h3>
								</div>
								
								<!-- URL -->
								<div class="text-sm theme-text-secondary font-mono truncate mb-1">
									{replay.url}
								</div>
								
								<!-- Metadata -->
								<div class="text-xs theme-text-secondary">
									Created {formatDate(replay.created_at)}
									{#if replay.updated_at !== replay.created_at}
										• Updated {formatDate(replay.updated_at)}
									{/if}
								</div>
							</div>

							<!-- Actions -->
							<div class="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity">
								<button
									on:click={() => handleExecute(replay)}
									class="p-2 text-green-400 hover:text-green-300 transition-colors"
									title="Execute replay"
								>
									<i class="fas fa-play"></i>
								</button>
								
								<button
									on:click={() => handleLogs(replay)}
									class="p-2 text-blue-400 hover:text-blue-300 transition-colors"
									title="View execution logs"
								>
									<i class="fas fa-history"></i>
								</button>
								
								<button
									on:click={() => handleEdit(replay)}
									class="p-2 theme-text-secondary hover:theme-text-primary transition-colors"
									title="Edit replay"
								>
									<i class="fas fa-edit"></i>
								</button>
								
								<button
									on:click={() => handleDelete(replay)}
									class="p-2 text-red-400 hover:text-red-300 transition-colors"
									title="Delete replay"
								>
									<i class="fas fa-trash"></i>
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
