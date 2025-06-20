<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { replays, filteredReplays, replayFilter, replayActions, selectedReplay } from '$lib/stores/replay';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import type { Replay } from '$lib/types/Replay';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';

	const dispatch = createEventDispatcher();

	let sortOrder: 'asc' | 'desc' = 'asc';
	let showAddDropdown = false;
	let searchTerm = '';
	let openMenuId: string | null = null;

	// Sort replays
	$: sortedReplays = $filteredReplays.sort((a, b) => {
		const comparison = a.name.localeCompare(b.name);
		return sortOrder === 'asc' ? comparison : -comparison;
	});

	// Update store when filters change
	$: {
		console.log('Updating replayFilter with searchTerm:', searchTerm);
		replayFilter.set({
			searchTerm: searchTerm,
			protocol: ''
		});
	}

	async function handleDelete(replay: Replay) {
		if (!$selectedWorkspace || !$selectedProject) return;

		if (!confirm(`Are you sure you want to delete the replay "${replay.name}"?`)) {
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
		
		// Close menu after action
		openMenuId = null;
	}

	function handleSelectReplay(replay: Replay) {
		console.log('Selecting replay:', replay);
		selectedReplay.set(replay);
		dispatch('edit', replay); // Dispatch an 'edit' event with the selected replay
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

	function handleCreateFolder(replay: Replay) {
		dispatch('add', { type: 'folder', parentReplay: replay });
		openMenuId = null;
	}

	function toggleMenu(replayId: string, event: Event) {
		event.stopPropagation();
		openMenuId = openMenuId === replayId ? null : replayId;
	}

	function closeMenu() {
		openMenuId = null;
	}

	function toggleSort() {
		sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
	}
</script>

<div 
	class="flex flex-col h-full theme-bg-primary theme-border border rounded-lg shadow-md overflow-hidden"
>
	<!-- Header Bar -->
	<div class="flex items-center justify-between p-3 theme-bg-secondary theme-border border-b">
		<!-- Search Section -->
		<div class="flex items-center space-x-3">
			<!-- Search Input -->
			<div class="relative">
				<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
					<i class="fas fa-search theme-text-muted text-sm"></i>
				</div>
				<input
					type="text"
					bind:value={searchTerm}
					placeholder="Search replays..."
					class="block w-full p-2 ps-10 text-sm rounded-lg theme-bg-primary theme-border border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400 dark:placeholder-gray-400"
				/>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="flex items-center space-x-2">
			<button
				on:click={toggleSort}
				class="p-2 theme-text-secondary hover:theme-text-primary transition-colors"
				title="Toggle sort order"
				aria-label="Toggle sort order"
			>
				<i class="fas fa-sort text-sm"></i>
			</button>
			
			<div class="add-dropdown-container relative">
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
					<div class="absolute top-full right-0 mt-1 w-48 theme-bg-primary theme-border border rounded-md shadow-lg z-20">
						<div class="py-1">
							<button
								on:click={handleAddHttp}
								class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
							>
								<i class="fas fa-globe text-blue-400"></i>
								<span>HTTP Replay</span>
							</button>
							<button
								on:click={handleAddFolder}
								class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
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
				<div class="text-center theme-text-muted">
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
			<div class="divide-y theme-border">
				{#each sortedReplays as replay (replay.id)}
					<div 
						class="group hover:theme-bg-secondary transition-colors cursor-pointer {$selectedReplay?.id === replay.id ? 'bg-blue-600/20 border-l-4 border-l-blue-500' : ''}"
						on:click={() => handleSelectReplay(replay)}
						role="button"
						tabindex="0"
						on:keydown={(e) => e.key === 'Enter' && handleSelectReplay(replay)}
					>
						<div class="flex items-center justify-between px-4 py-3">
							<!-- Method Icon and Info -->
							<div class="flex items-center space-x-3 flex-1 min-w-0">
								
								<HttpMethodBadge method={replay.method} />
								
								<div class="flex-1 min-w-0">
									<h4 class="text-sm font-medium {$selectedReplay?.id === replay.id ? 'text-blue-300' : 'theme-text-primary'} truncate">
										{replay.name}
									</h4>
									<p class="text-xs theme-text-muted truncate">
										{replay.url}
									</p>
								</div>
							</div>

							<!-- Three-dot menu -->
							<div class="menu-container relative">
								<button
									on:click={(e) => toggleMenu(replay.id, e)}
									class="p-1.5 theme-text-muted hover:theme-text-primary transition-colors opacity-0 group-hover:opacity-100 {openMenuId === replay.id ? 'opacity-100' : ''}"
									title="More options"
									aria-label="More options"
									aria-expanded={openMenuId === replay.id}
								>
									<i class="fas fa-ellipsis-v text-xs"></i>
								</button>

								<!-- Dropdown Menu -->
								{#if openMenuId === replay.id}
									<div class="absolute top-full right-0 mt-1 w-48 theme-bg-primary theme-border border rounded-md shadow-lg z-20">
										<div class="py-1">
											<button
												on:click={() => handleDelete(replay)}
												class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
											>
												<i class="fas fa-trash text-red-400"></i>
												<span>Delete</span>
											</button>
											<button
												on:click={() => handleCreateFolder(replay)}
												class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
											>
												<i class="fas fa-folder-plus text-yellow-400"></i>
												<span>Create Folder</span>
											</button>
										</div>
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>
