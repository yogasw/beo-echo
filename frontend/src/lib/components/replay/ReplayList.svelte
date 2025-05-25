<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { replays, filteredReplays, replayFilter, replayActions } from '$lib/stores/replay';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import { HTTP_METHODS } from '$lib/types/Replay';
	import type { Replay } from '$lib/types/Replay';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
	import StatusCodeBadge from '$lib/components/common/StatusCodeBadge.svelte';

	const dispatch = createEventDispatcher();

	let searchTerm = '';
	let sortOrder: 'asc' | 'desc' = 'asc';
	let showDropdown = false;
	let showAddDropdown = false;

	// Sort replays
	$: sortedReplays = $filteredReplays.sort((a, b) => {
		const comparison = a.alias.localeCompare(b.alias);
		return sortOrder === 'asc' ? comparison : -comparison;
	});

	// Update store when filters change
	$: replayFilter.set({
		searchTerm,
		method: '',
		protocol: ''
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

	function handleAdd() {
		showAddDropdown = !showAddDropdown;
	}

	function handleAddHttp() {
		dispatch('add', { type: 'http' });
		showAddDropdown = false;
	}

	function handleAddFolder() {
		dispatch('add', { type: 'folder' });
		showAddDropdown = false;
	}

	function toggleSort() {
		sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
	}
</script>

<div class="flex flex-col h-full theme-bg-primary border border-gray-700 rounded-lg shadow-md overflow-hidden">
	<!-- Header Bar -->
	<div class="flex items-center justify-between p-3 bg-gray-750 border-b border-gray-700">
		<!-- Filter Section -->
		<div class="flex items-center space-x-3">
			<div class="relative">
				<button 
					on:click={() => showDropdown = !showDropdown}
					class="flex items-center space-x-2 px-3 py-1.5 bg-gray-700 hover:bg-gray-600 text-white rounded-md text-sm transition-colors"
				>
					<i class="fas fa-search text-xs"></i>
					<span>Filter</span>
					<i class="fas fa-chevron-down text-xs {showDropdown ? 'rotate-180' : ''} transition-transform"></i>
				</button>

				<!-- Dropdown Filter Panel -->
				{#if showDropdown}
					<div class="absolute top-full left-0 mt-1 w-72 bg-gray-800 border border-gray-700 rounded-md shadow-lg z-10">
						<div class="p-3 space-y-3">
							<div>
								<label class="block text-xs font-medium text-gray-300 mb-1">Search</label>
								<input
									type="text"
									bind:value={searchTerm}
									placeholder="Search by alias or URL..."
									class="w-full px-3 py-2 bg-gray-700 border border-gray-600 rounded text-white text-sm placeholder-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
								/>
							</div>
							<div class="flex justify-end">
								<button
									on:click={() => { searchTerm = ''; showDropdown = false; }}
									class="px-3 py-1 bg-blue-600 hover:bg-blue-700 text-white text-xs rounded"
								>
									Apply
								</button>
							</div>
						</div>
					</div>
				{/if}
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="flex items-center space-x-2">
			<button
				on:click={toggleSort}
				class="p-2 text-gray-300 hover:text-white transition-colors"
				title="Toggle sort order"
				aria-label="Toggle sort order"
			>
				<i class="fas fa-sort text-sm"></i>
			</button>
			
			<div class="relative">
				<button
					on:click={handleAdd}
					title="Add new replay"
					aria-label="Add new replay"
					class="flex items-center space-x-1 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm transition-colors"
				>
					<i class="fas fa-plus text-xs"></i>
					<i class="fas fa-chevron-down text-xs {showAddDropdown ? 'rotate-180' : ''} transition-transform"></i>
				</button>

				<!-- Add Options Dropdown -->
				{#if showAddDropdown}
					<div class="absolute top-full right-0 mt-1 w-48 bg-gray-800 border border-gray-700 rounded-md shadow-lg z-20">
						<div class="py-1">
							<button
								on:click={handleAddHttp}
								class="w-full text-left px-4 py-2 text-sm text-white hover:bg-gray-700 transition-colors flex items-center space-x-2"
							>
								<i class="fas fa-globe text-blue-400"></i>
								<span>HTTP Replay</span>
							</button>
							<button
								on:click={handleAddFolder}
								class="w-full text-left px-4 py-2 text-sm text-white hover:bg-gray-700 transition-colors flex items-center space-x-2"
							>
								<i class="fas fa-folder text-yellow-400"></i>
								<span>Folder</span>
							</button>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>

	<!-- Content Area -->
	<div class="flex-1 overflow-auto">
		{#if $filteredReplays.length === 0}
			<div class="flex items-center justify-center h-full">
				<div class="text-center text-gray-400">
					{#if $replays.length === 0}
						<i class="fas fa-play-circle text-4xl mb-4 opacity-50"></i>
						<h3 class="text-lg font-medium mb-2">No replays yet</h3>
						<p class="text-sm">Create your first API replay to get started</p>
					{:else}
						<i class="fas fa-search text-4xl mb-4 opacity-50"></i>
						<h3 class="text-lg font-medium mb-2">No matching replays</h3>
						<p class="text-sm">Try adjusting your search criteria</p>
					{/if}
				</div>
			</div>
		{:else}
			<div class="divide-y divide-gray-700">
				{#each sortedReplays as replay (replay.id)}
					<div class="group hover:bg-gray-750 transition-colors">
						<div class="flex items-center justify-between px-4 py-1">
							<!-- Method Icon and Info -->
							<div class="flex items-center space-x-3 flex-1 min-w-0">
								
								<HttpMethodBadge method={replay.method} />
								
								<div class="flex-1 min-w-0">
									<h4 class="text-sm font-medium text-white truncate">
										{replay.alias}
									</h4>
									<p class="text-xs text-gray-400 truncate">
										{replay.url}
									</p>
								</div>
							</div>

							<!-- Actions -->
							<div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
								<button
									on:click={() => handleEdit(replay)}
									class="p-1.5 text-gray-400 hover:text-white transition-colors"
									title="Edit"
									aria-label="Edit replay"
								>
									<i class="fas fa-edit text-xs"></i>
								</button>
								
								<button
									on:click={() => handleDelete(replay)}
									class="p-1.5 text-gray-400 hover:text-red-400 transition-colors"
									title="Delete"
									aria-label="Delete replay"
								>
									<i class="fas fa-trash text-xs"></i>
								</button>
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Click outside to close dropdown -->
{#if showDropdown}
	<div class="fixed inset-0 z-0" on:click={() => showDropdown = false}></div>
{/if}

<!-- Click outside to close add dropdown -->
{#if showAddDropdown}
	<div class="fixed inset-0 z-10" on:click={() => showAddDropdown = false}></div>
{/if}
