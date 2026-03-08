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
	import HttpMethodBadge from '$lib/components/common/HttpMethodBadge.svelte';

	let {
		isPanelCollapsed = false
	}: {
		isPanelCollapsed?: boolean;
	} = $props();

	const dispatch = createEventDispatcher();

	let sortOrder: 'asc' | 'desc' = $state('asc');
	let showAddDropdown = $state(false);
	let searchTerm = $state('');
	let collapsedFolders = $state(new Set<string>());
	let inlineFolderParent: string | null | undefined = $state(undefined);
	let newFolderName = $state('');

	let draggedItem: any | null = $state(null);
	let dragOverItemId: string | null = $state(null);
	let dragOverPosition: 'top' | 'bottom' | 'inside' | 'root' | null = $state(null);
	let dragCounter = $state(0);

	let contextMenu = $state<{ isOpen: boolean; x: number; y: number; item: any | null }>({
		isOpen: false,
		x: 0,
		y: 0,
		item: null
	});

	function handleContextMenu(event: MouseEvent, item: any) {
		event.preventDefault();
		contextMenu = {
			isOpen: true,
			x: event.clientX,
			y: event.clientY,
			item
		};
	}

	function closeContextMenu() {
		contextMenu.isOpen = false;
	}

	function toggleFolderCollapse(e: MouseEvent, folderId: string) {
		e.stopPropagation();
		const newCollapsed = new Set(collapsedFolders);
		if (newCollapsed.has(folderId)) {
			newCollapsed.delete(folderId);
		} else {
			newCollapsed.add(folderId);
		}
		collapsedFolders = newCollapsed;
	}

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

	// Build a hierarchical ordered list that computes depth and handles collapse state
	let displayItems = $derived.by(() => {
		const items = sortedReplays;
		const result: any[] = [];
		const folders = items.filter(i => i.itemType === 'folder');
		const replaysList = items.filter(i => i.itemType === 'replay');

		// Defend against circular references
		const processedFolderIds = new Set<string>();
		
		// Setup known folders to handle orphaned items
		const knownFolderIds = new Set(folders.map(f => f.id));

		function addChildren(parentId: string | null, depth: number, isVisible: boolean) {
			// Add folders first
			const childFolders = folders.filter(f => {
				const fPid = f.parent_id ? String(f.parent_id).trim() : null;
				// If parentId is null (we are currently at root), also include orphaned folders whose parent no longer exists
				if (parentId === null) {
					return fPid === null || !knownFolderIds.has(fPid);
				}
				return fPid === parentId;
			});
			
			for (const folder of childFolders) {
				if (processedFolderIds.has(folder.id)) continue;
				processedFolderIds.add(folder.id);
				
				if (isVisible) {
					result.push({ ...folder, depth, isVisible });
				}
				
				// Calculate children recursively
				const isExpanded = !collapsedFolders.has(folder.id);
				addChildren(folder.id, depth + 1, isVisible && isExpanded);
			}

			// Add replays
			const childReplays = replaysList.filter(r => {
				const rFid = r.folder_id ? String(r.folder_id).trim() : null;
				// If parentId is null (we are currently at root), also include orphaned replays whose folder no longer exists
				if (parentId === null) {
					return rFid === null || !knownFolderIds.has(rFid);
				}
				return rFid === parentId;
			});
			for (const replay of childReplays) {
				if (isVisible) {
					result.push({ ...replay, depth, isVisible });
				}
			}
		}

		addChildren(null, 0, true);
		return result;
	});

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

		let message = `Are you sure you want to delete the ${item.itemType} "${item.name}"?`;
		if (item.itemType === 'folder') {
			message = `Are you sure you want to delete the folder "${item.name}"? All nested folders and replays inside it will also be deleted. This action cannot be undone.`;
		}

		if (!confirm(message)) {
			return;
		}

		try {
			replayActions.setLoading('delete', true);
			
			if (item.itemType === 'folder') {
				await replayApi.deleteFolder($selectedWorkspace.id, $selectedProject.id, item.id);
				dispatch('refresh');
				toast.success('Folder deleted successfully');
			} else {
				await replayApi.deleteReplay($selectedWorkspace.id, $selectedProject.id, item.id);
				replayActions.removeReplay(item.id);
				toast.success('Replay deleted successfully');
			}
		} catch (err: any) {
			toast.error(err);
		} finally {
			replayActions.setLoading('delete', false);
		}

		// Close menu after action
		closeContextMenu();
	}

	function handleSelectReplay(item: any) {
		if (item.itemType === 'folder') {
			// Select folder to show documentation view
			selectedReplay.set(item);
			dispatch('edit', item);
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
		inlineFolderParent = null; // null means root level
		newFolderName = '';
		showAddDropdown = false;
	}

	function handleCreateFolder(item: any) {
		closeContextMenu();
		inlineFolderParent = item.id;
		newFolderName = '';
		// Ensure the parent folder is expanded
		const newCollapsed = new Set(collapsedFolders);
		newCollapsed.delete(item.id);
		collapsedFolders = newCollapsed;
	}

	function handleCreateHttpInFolder(item: any) {
		closeContextMenu();
		dispatch('add', { type: 'http', parentReplay: { id: item.id } });
	}

	function handleDuplicate(item: any) {
		closeContextMenu();
		dispatch('duplicate', item);
	}

	function submitInlineFolder() {
		if (newFolderName.trim()) {
			const parentItem = inlineFolderParent === null ? null : { id: inlineFolderParent };
			dispatch('add', { type: 'folder', name: newFolderName.trim(), parentReplay: parentItem });
		}
		inlineFolderParent = undefined;
		newFolderName = '';
	}

	function handleDragStart(e: DragEvent, item: any) {
		draggedItem = item;
		if (e.dataTransfer) {
			e.dataTransfer.effectAllowed = 'move';
			e.dataTransfer.setData('text/plain', item.id);
			// Optional: Custom ghost image could be set here
		}
	}

	function handleDragOver(e: DragEvent, targetItem: any) {
		e.preventDefault();
		e.stopPropagation();
		
		if (!draggedItem || draggedItem.id === targetItem.id) return;
		
		// Determine position for visual feedback
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		const y = e.clientY - rect.top;
		
		dragOverItemId = targetItem.id;
		
		if (targetItem.itemType === 'folder') {
			// Droppable inside folder
			dragOverPosition = 'inside';
		} else {
			// If dropping on a replay, determine if it's top or bottom (for reordering, which isn't fully implemented backend yet, so just fallback to same folder)
			// For now, if we drop on a replay, we just move it to that replay's folder
			dragOverPosition = 'inside';
		}

		// Prevent folder dropping into itself
		if (draggedItem.itemType === 'folder' && draggedItem.id === targetItem.id) {
			dragOverItemId = null;
			dragOverPosition = null;
			return;
		}

		if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
	}

	function handleDragEnter(e: DragEvent, targetItem: any) {
		e.preventDefault();
		dragCounter++;
	}

	function handleDragLeave(e: DragEvent, targetItem: any) {
		dragCounter--;
		if (dragCounter === 0) {
			// dragOverItemId = null;
			// dragOverPosition = null;
		}
		// A more reliable way is just to nullify if we leave the main bounds, 
		// but since we only care about drop, we can be a bit loose with leave or just use CSS pointer-events.
		// To fix flickering, we won't clear dragOverItemId on leave, we'll clear it on drop or dragend.
	}
	
	function handleDragEnd(e: DragEvent) {
		draggedItem = null;
		dragOverItemId = null;
		dragOverPosition = null;
		dragCounter = 0;
	}

	async function handleDrop(e: DragEvent, targetItem: any) {
		e.preventDefault();
		e.stopPropagation();
		dragOverItemId = null;
		dragOverPosition = null;
		dragCounter = 0;
		
		if (!draggedItem || draggedItem.id === targetItem.id) return;
		
		const itemId = draggedItem.id;
		const itemType = draggedItem.itemType;
		// If dropped on a folder, target is the folder. If dropped on a replay, target is the replay's folder.
		const targetId = targetItem.itemType === 'folder' ? targetItem.id : targetItem.folder_id || targetItem.parent_id;
		
		// Prevent dropping folder into itself
		if (itemType === 'folder' && itemId === targetId) return;

		if (!$selectedWorkspace || !$selectedProject) return;

		// Optimistic UI update
		replayActions.moveItem(itemId, itemType, targetId);

		try {
			if (itemType === 'replay') {
				await replayApi.updateReplay($selectedWorkspace.id, $selectedProject.id, itemId, {
					folder_id: targetId,
					update_folder_id: true
				});
			} else {
				await replayApi.updateFolder($selectedWorkspace.id, $selectedProject.id, itemId, {
					parent_id: targetId,
					update_parent_id: true
				});
			}
			toast.success(`Moved ${itemType} to folder`);
		} catch (error: any) {
			toast.error(`Failed to move: ${error.message || 'Unknown error'}`);
			dispatch('refresh'); // Refresh on error to restore state
		}
		
		draggedItem = null;
	}

	function handleRootDragOver(e: DragEvent) {
		e.preventDefault();
		if (draggedItem) {
			dragOverPosition = 'root';
			if (e.dataTransfer) e.dataTransfer.dropEffect = 'move';
		}
	}

	function handleRootDragLeave(e: DragEvent) {
		const rect = (e.currentTarget as HTMLElement).getBoundingClientRect();
		const { clientX: x, clientY: y } = e;
		if (
			x <= rect.left ||
			x >= rect.right ||
			y <= rect.top ||
			y >= rect.bottom
		) {
			dragOverPosition = null;
		}
	}

	async function handleRootDrop(e: DragEvent) {
		e.preventDefault();
		dragOverPosition = null;
		
		if (!draggedItem) return;
		
		const itemId = draggedItem.id;
		const itemType = draggedItem.itemType;

		if (!$selectedWorkspace || !$selectedProject) return;

		// Optimistic UI update
		replayActions.moveItem(itemId, itemType, null);

		try {
			if (itemType === 'replay') {
				await replayApi.updateReplay($selectedWorkspace.id, $selectedProject.id, itemId, {
					folder_id: null,
					update_folder_id: true
				});
			} else {
				await replayApi.updateFolder($selectedWorkspace.id, $selectedProject.id, itemId, {
					parent_id: null,
					update_parent_id: true
				});
			}
			toast.success(`Moved ${itemType} to root`);
		} catch (error: any) {
			toast.error(`Failed to move: ${error.message || 'Unknown error'}`);
			dispatch('refresh'); // Refresh on error to restore state
		}
		
		draggedItem = null;
	}

	function toggleSort() {
		sortOrder = sortOrder === 'asc' ? 'desc' : 'asc';
	}
</script>

<svelte:window onclick={closeContextMenu} />

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
		<div class="flex-1 overflow-auto"
			 role="group"
			 ondragover={handleRootDragOver}
			 ondragleave={handleRootDragLeave}
			 ondrop={handleRootDrop}
		>
			{#if dragOverPosition === 'root'}
				<div class="bg-blue-500/10 border-2 border-dashed border-blue-500 rounded-lg m-2 p-4 flex items-center justify-center pointer-events-none">
					<span class="text-blue-500 font-medium">Move to Top Level</span>
				</div>
			{/if}

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
					<!-- ROOT INLINE FOLDER -->
					{#if inlineFolderParent === null}
						<div class="group flex items-center transition-all px-2 py-1">
							<div class="flex items-center w-full gap-2 pl-2">
								<i class="fas fa-folder text-yellow-400"></i>
								<!-- svelte-ignore a11y_autofocus -->
								<input
									type="text"
									bind:value={newFolderName}
									class="flex-1 bg-transparent border-b border-blue-500 outline-none text-[13px] theme-text-primary py-1 px-1"
									placeholder="Folder name"
									autofocus
									onkeydown={(e) => {
										if (e.key === 'Enter') submitInlineFolder();
										if (e.key === 'Escape') inlineFolderParent = undefined;
									}}
									onblur={submitInlineFolder}
								/>
							</div>
						</div>
					{/if}

					{#each displayItems as item (item.id)}
						<div
							draggable="true"
							ondragstart={(e) => handleDragStart(e, item)}
							ondragover={(e) => handleDragOver(e, item)}
							ondragenter={(e) => handleDragEnter(e, item)}
							ondragleave={(e) => handleDragLeave(e, item)}
							ondragend={(e) => handleDragEnd(e)}
							ondrop={(e) => handleDrop(e, item)}
							class="group flex transition-all cursor-pointer {draggedItem?.id === item.id ? 'opacity-30' : ''} {dragOverItemId === item.id ? (item.itemType === 'folder' ? 'bg-[#ff6600]/10 ring-1 ring-inset ring-[#ff6600] z-10' : 'border-t-2 border-t-[#ff6600]') : 'border-t-2 border-t-transparent'} {$selectedReplay?.id ===
							item.id && dragOverItemId !== item.id
								? 'bg-blue-500/10 border-l-[3px] border-l-[#ff6600]'
								: 'border-l-[3px] border-l-transparent hover:bg-gray-500/10'} {contextMenu.isOpen && contextMenu.item?.id === item.id ? 'bg-gray-500/10 relative z-40' : 'relative'} px-2"
							style="padding-left: {item.depth * 1.5}rem; padding-right: 0.5rem;"
							onclick={() => handleSelectReplay(item)}
							oncontextmenu={(e) => handleContextMenu(e, item)}
							role="button"
							tabindex="0"
							onkeydown={(e) => e.key === 'Enter' && handleSelectReplay(item)}
						>
							{#if item.depth > 0}
								<!-- Vertical guide line -->
								{#each Array(item.depth) as _, i}
									<div class="absolute top-0 bottom-0 border-l border-gray-700/50 pointer-events-none" style="left: {i * 1.5 + 0.625}rem"></div>
								{/each}
							{/if}
							
							<div class="flex items-center w-full h-[32px] gap-2 {draggedItem ? 'pointer-events-none' : ''}">
								<!-- Icon column -->
								<div class="w-[45px] flex-shrink-0 flex items-center justify-end gap-2">
									{#if item.itemType === 'folder'}
										<button 
											onclick={(e) => toggleFolderCollapse(e, item.id)}
											class="w-4 h-4 flex items-center justify-center theme-text-muted hover:theme-text-primary transition-colors rounded hover:bg-gray-500/10"
											aria-label="Toggle Folder"
										>
											<i class="fas fa-chevron-{collapsedFolders.has(item.id) ? 'right' : 'down'} text-[10px]"></i>
										</button>
										<i class="far fa-folder theme-text-muted text-sm"></i>
									{:else}
										<HttpMethodBadge method={item.method} size="folder-size" />
									{/if}
								</div>
								
								<!-- Name -->
								<span
									class="text-[13px] {$selectedReplay?.id === item.id
										? 'theme-text-primary font-medium'
										: 'theme-text-secondary'} truncate"
								>
									{item.name || 'Unnamed Request'}
								</span>
							</div>

							<!-- Action Buttons -->
							<div
								class="absolute right-2 top-1/2 -translate-y-1/2 flex items-center gap-1 transition-opacity pl-2 opacity-0 group-hover:opacity-100 z-30 {draggedItem ? 'pointer-events-none' : ''}"
							>
								<!-- Three-dot menu -->
								<div class="menu-container relative">
									<button
										onclick={(e) => {
											e.stopPropagation();
											handleContextMenu(e, item)
										}}
										class="p-1.5 theme-text-muted hover:theme-text-primary transition-colors"
										title="More options"
										aria-label="More options"
									>
										<i class="fas fa-ellipsis-h text-sm"></i>
									</button>
								</div>
							</div>
						</div>
						
						<!-- NESTED INLINE FOLDER -->
						{#if inlineFolderParent === item.id}
							<div class="group flex items-center transition-all px-2 py-1" style="padding-left: {(item.depth + 1) * 1.5}rem; padding-right: 0.5rem;">
								{#each Array(item.depth + 1) as _, i}
									<div class="absolute top-0 bottom-0 border-l border-gray-700/50 pointer-events-none" style="left: {i * 1.5 + 0.625}rem"></div>
								{/each}
								<div class="flex flex-row items-center w-full gap-2 pl-6 relative z-10">
									<i class="fas fa-folder text-yellow-400"></i>
									<!-- svelte-ignore a11y_autofocus -->
									<input
										type="text"
										bind:value={newFolderName}
										class="flex-1 bg-transparent border-b border-blue-500 outline-none text-[13px] theme-text-primary py-1 px-1"
										placeholder="Folder name"
										autofocus
										onkeydown={(e) => {
											if (e.key === 'Enter') submitInlineFolder();
											if (e.key === 'Escape') inlineFolderParent = undefined;
										}}
										onblur={submitInlineFolder}
									/>
								</div>
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		</div>
	{/if}

	<!-- Global Context Menu Overlay -->
	{#if contextMenu.isOpen && contextMenu.item}
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div
			class="fixed theme-bg-primary theme-border border rounded-md shadow-xl z-[100] overflow-hidden min-w-[160px]"
			style="top: {contextMenu.y}px; left: {contextMenu.x}px;"
			onclick={(e) => e.stopPropagation()}
			oncontextmenu={(e) => e.stopPropagation()}
		>
			<div class="py-1">
				<button
					onclick={(e) => {
						e.stopPropagation();
						handleDelete(contextMenu.item);
					}}
					class="w-full text-left px-4 py-2 text-sm text-red-500 hover:theme-bg-secondary transition-colors flex items-center space-x-2"
				>
					<i class="fas fa-trash"></i>
					<span>Delete</span>
				</button>
				{#if contextMenu.item.itemType === 'folder'}
					<button
						onclick={(e) => {
							e.stopPropagation();
							handleCreateHttpInFolder(contextMenu.item);
						}}
						class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
					>
						<i class="fas fa-globe text-blue-400"></i>
						<span>Create Request</span>
					</button>
					<button
						onclick={(e) => {
							e.stopPropagation();
							handleCreateFolder(contextMenu.item);
						}}
						class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
					>
						<i class="fas fa-folder-plus text-yellow-400"></i>
						<span>Create Folder</span>
					</button>
				{/if}
				{#if contextMenu.item.itemType === 'replay'}
					<button
						onclick={(e) => {
							e.stopPropagation();
							handleDuplicate(contextMenu.item);
						}}
						class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:theme-bg-secondary transition-colors flex items-center space-x-2"
					>
						<i class="fas fa-copy text-blue-400"></i>
						<span>Duplicate</span>
					</button>
				{/if}
			</div>
		</div>
	{/if}
</div>
