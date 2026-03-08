<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import {
		replays,
		filteredReplays,
		replayFilter,
		replayActions,
		selectedReplay
	} from '$lib/stores/replay';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import type { Replay } from '$lib/types/Replay';
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';
	import { getHttpMethodColor } from '$lib/utils/badgeUtils';

	let {
		isPanelCollapsed = false
	}: {
		isPanelCollapsed?: boolean;
	} = $props();

	const dispatch = createEventDispatcher();

	let sortOrder: 'asc' | 'desc' = $state('asc');
	let showAddDropdown = $state(false);
	let searchTerm = $state('');
	let openMenuId: string | null = $state(null);
	let pendingFolderParent: any = $state(null);

	function handleToggleCollapse() {
		dispatch('toggleCollapse', { animate: true });
	}

	// Sort replays
	let sortedReplays = $derived(
		[...$filteredReplays].sort((a, b) => {
			// Always put folders first in the list
			if (a.itemType !== b.itemType) {
				return a.itemType === 'folder' ? -1 : 1;
			}
			const comparison = (a.name || '').localeCompare(b.name || '');
			return sortOrder === 'asc' ? comparison : -comparison;
		})
	);

	// Update store when filters change
	$effect(() => {
		console.log('Updating replayFilter with searchTerm:', searchTerm);
		replayFilter.set({
			searchTerm: searchTerm,
			protocol: ''
		});
	});

	async function handleDelete(item: any) {
		if (!$selectedWorkspace || !$selectedProject) return;

		if (!confirm(`Are you sure you want to delete the ${item.itemType} "${item.name}"?`)) {
			return;
		}

		try {
			replayActions.setLoading('delete', true);
			// TODO: Api doesn't have delete folder yet, so let's just assume we delete replay for now
			// or conditionally call delete if it's a replay.
			if (item.itemType === 'folder') {
				toast.error('Deleting folders is not yet supported in this view.');
				replayActions.setLoading('delete', false);
				return;
			}
			await replayApi.deleteReplay($selectedWorkspace.id, $selectedProject.id, item.id);
			replayActions.removeReplay(item.id);
			toast.success(`${item.itemType === 'folder' ? 'Folder' : 'Replay'} deleted successfully`);
		} catch (err: any) {
			toast.error(err);
		} finally {
			replayActions.setLoading('delete', false);
		}

		// Close menu after action
		openMenuId = null;
	}

	function handleSelectReplay(item: any) {
		if (item.itemType === 'folder') {
			return;
		}

		console.log('Selecting replay:', item);
		selectedReplay.set(item);
		dispatch('edit', item); // Dispatch an 'edit' event with the selected replay
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

	function handleCreateFolder(item: any) {
		openMenuId = null;
		pendingFolderParent = item;
		dispatch('add', { type: 'folder', parentReplay: item });
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
	{#if isPanelCollapsed}
		<!-- Collapsed state: Show expand button centered vertically -->
		<div class="flex flex-col items-center justify-center h-full">
			<button
				onclick={handleToggleCollapse}
				class="theme-text-primary hover:text-blue-500 px-2 py-2 rounded hover:bg-blue-500/10 transition-all duration-200 border border-gray-700/50 hover:border-blue-500/50"
				title="Expand panel"
				aria-label="Expand panel"
			>
				<i class="fas fa-angle-double-right text-lg"></i>
			</button>
		</div>
	{:else}
		<!-- Expanded state: Show full content -->
		<!-- Header Bar -->
		<div class="flex items-center justify-between p-3 theme-bg-secondary theme-border border-b">
			<!-- Search Section -->
			<div class="flex items-center space-x-3 flex-1">
				<!-- Search Input -->
				<div class="relative flex-1">
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
					onclick={toggleSort}
					class="p-2 theme-text-secondary hover:theme-text-primary transition-colors"
					title="Toggle sort order"
					aria-label="Toggle sort order"
				>
					<i class="fas fa-sort text-sm"></i>
				</button>

				<div class="add-dropdown-container relative">
					<button
						onclick={handleAdd}
						title="Add new replay"
						aria-label="Add new replay"
						class="flex items-center space-x-1 px-3 py-1.5 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm transition-colors"
					>
						<i class="fas fa-plus text-xs"></i>
						<i
							class="fas fa-chevron-down text-xs {showAddDropdown
								? 'rotate-180'
								: ''} transition-transform"
						></i>
					</button>

					<!-- Add Options Dropdown -->
					{#if showAddDropdown}
						<div
							class="absolute top-full right-0 mt-1 w-48 theme-bg-primary theme-border border rounded-md shadow-lg z-20"
						>
							<div class="py-1">
								<button
									onclick={handleAddHttp}
									class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
								>
									<i class="fas fa-globe text-blue-400"></i>
									<span>HTTP Replay</span>
								</button>
								<button
									onclick={handleAddFolder}
									class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
								>
									<i class="fas fa-folder text-yellow-400"></i>
									<span>Folder</span>
								</button>
							</div>
						</div>
					{/if}
				</div>

				<!-- Collapse button -->
				<button
					onclick={handleToggleCollapse}
					class="theme-text-primary hover:text-blue-500 px-2 py-1 rounded hover:bg-blue-500/10 transition-all duration-200 border border-gray-700/50 hover:border-blue-500/50 ml-2"
					title="Collapse panel"
					aria-label="Collapse panel"
				>
					<i class="fas fa-angle-double-left text-sm"></i>
				</button>
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
				<div class="flex flex-col">
					{#each sortedReplays as item (item.id)}
						<div
							class="group flex relative transition-colors cursor-pointer border-l-[3px] px-2 {$selectedReplay?.id ===
							item.id
								? 'bg-white/5 border-[#3b82f6]'
								: 'border-transparent hover:bg-white/5'}"
							onclick={() => handleSelectReplay(item)}
							role="button"
							tabindex="0"
							onkeydown={(e) => e.key === 'Enter' && handleSelectReplay(item)}
						>
							<div class="flex items-center w-full h-[32px] gap-3">
								<!-- Fixed-width badge col so names always align -->
								<div class="w-[50px] flex-shrink-0 flex justify-end">
									{#if item.itemType === 'folder'}
										<i class="fas fa-folder text-yellow-400"></i>
									{:else}
										<HttpMethodBadge method={item.method} size="xs" />
									{/if}
								</div>
								<!-- Name -->
								<span
									class="text-[13px] {$selectedReplay?.id === item.id
										? 'text-gray-100'
										: 'text-gray-300 dark:text-gray-300'} truncate"
								>
									{item.name || 'Unnamed Request'}
								</span>
							</div>

							<!-- Action Buttons -->
							<div
								class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity bg-gray-900/90 pl-2 rounded-l"
							>
								<!-- Three-dot menu -->
								<div class="menu-container relative">
									<button
										onclick={(e) => toggleMenu(item.id, e)}
										class="p-1.5 text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors {openMenuId ===
										item.id
											? 'opacity-100 text-gray-200'
											: ''}"
										title="More options"
										aria-label="More options"
										aria-expanded={openMenuId === item.id}
									>
										<i class="fas fa-ellipsis-h text-sm"></i>
									</button>

									<!-- Dropdown Menu -->
									{#if openMenuId === item.id}
										<div
											class="absolute top-full right-0 mt-1 w-48 theme-bg-primary theme-border border rounded-md shadow-lg z-20"
										>
											<div class="py-1">
												<button
													onclick={(e) => {
														e.stopPropagation();
														handleDelete(item);
													}}
													class="w-full text-left px-4 py-2 text-sm text-red-400 hover:theme-bg-secondary transition-colors flex items-center space-x-2"
												>
													<i class="fas fa-trash"></i>
													<span>Delete</span>
												</button>
												{#if item.itemType === 'folder'}
													<button
														onclick={(e) => {
															e.stopPropagation();
															handleCreateFolder(item);
														}}
														class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
													>
														<i class="fas fa-folder-plus text-yellow-400"></i>
														<span>Create Folder</span>
													</button>
												{/if}
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
	{/if}
</div>
