<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replayApi } from '$lib/api/replayApi';
	import { toast } from '$lib/stores/toast';
	import { replays, replayFolders } from '$lib/stores/replay';
	import FolderOverviewTab from './tabs/FolderOverviewTab.svelte';

	let {
		folder
	}: {
		folder: any;
	} = $props();

	const dispatch = createEventDispatcher();

	// --- Editable title ---
	let isEditingTitle = $state(false);
	let titleValue = $state('');
	let isSavingTitle = $state(false);
	let titleInputEl: HTMLInputElement | undefined = $state(undefined);

	$effect(() => {
		titleValue = folder?.name || '';
		isEditingTitle = false;
	});

	function startEditTitle() {
		isEditingTitle = true;
		titleValue = folder?.name || '';
		setTimeout(() => titleInputEl?.focus(), 0);
	}

	function cancelEditTitle() {
		isEditingTitle = false;
		titleValue = folder?.name || '';
	}

	async function saveTitle() {
		const trimmed = titleValue.trim();
		if (!trimmed || trimmed === folder?.name) {
			cancelEditTitle();
			return;
		}
		if (!$selectedWorkspace || !$selectedProject) return;

		try {
			isSavingTitle = true;
			await replayApi.updateFolder($selectedWorkspace.id, $selectedProject.id, folder.id, {
				name: trimmed
			});
			folder = { ...folder, name: trimmed };
			isEditingTitle = false;
			toast.success('Folder renamed');
			dispatch('folderUpdated', { ...folder, name: trimmed });
		} catch (err: any) {
			toast.error(err);
		} finally {
			isSavingTitle = false;
		}
	}

	function handleTitleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') saveTitle();
		if (e.key === 'Escape') cancelEditTitle();
	}

	// --- Child count ---
	let childCount = $derived.by(() => {
		const childReplays = $replays.filter((r: any) => r.folder_id === folder?.id).length;
		const childFolders = $replayFolders.filter((f: any) => f.parent_id === folder?.id).length;
		return childReplays + childFolders;
	});

	// --- Breadcrumb ---
	let breadcrumbParts = $derived.by(() => {
		const parts: string[] = [];
		let current = folder;
		const visited = new Set<string>();
		while (current?.parent_id && !visited.has(current.id)) {
			visited.add(current.id);
			const parent = $replayFolders.find((f: any) => f.id === current.parent_id);
			if (parent) {
				parts.unshift(parent.name);
				current = parent;
			} else {
				break;
			}
		}
		return parts;
	});
</script>

<div class="flex flex-col h-full overflow-hidden theme-bg-primary">
	<!-- ── Tab bar (matches ReplayEditor style) ── -->
	<div class="flex-shrink-0 theme-border border-b theme-bg-secondary">
		<!-- Breadcrumb row -->
		{#if breadcrumbParts.length > 0}
			<nav aria-label="Folder breadcrumb" class="flex items-center gap-1 text-xs theme-text-muted px-4 pt-2">
				{#each breadcrumbParts as part, i}
					<span class="theme-text-muted opacity-75">{part}</span>
					<i class="fas fa-chevron-right text-[9px] opacity-50"></i>
				{/each}
				<span class="theme-text-primary font-medium">{folder?.name}</span>
			</nav>
		{/if}

		<!-- Title + tab row -->
		<div class="flex items-center justify-between px-4 pt-2 pb-0">
			<!-- Left: folder icon + editable title -->
			<div class="flex items-center gap-2.5 group min-w-0 flex-1 mr-4">
				<i class="far fa-folder-open text-yellow-400 flex-shrink-0"></i>

				{#if isEditingTitle}
					<div class="flex items-center gap-1.5 flex-1 min-w-0">
						<input
							bind:this={titleInputEl}
							type="text"
							bind:value={titleValue}
							onkeydown={handleTitleKeydown}
							onblur={saveTitle}
							class="flex-1 min-w-0 theme-bg-primary theme-border border border-blue-500 rounded px-2 py-0.5 text-sm font-semibold theme-text-primary focus:outline-none focus:ring-1 focus:ring-blue-500"
							title="Folder name"
							aria-label="Edit folder name"
							aria-required="true"
						/>
						<button
							onclick={saveTitle}
							disabled={isSavingTitle}
							class="flex-shrink-0 px-2 py-1 bg-blue-600 hover:bg-blue-700 disabled:opacity-50 text-white rounded text-xs transition-colors"
							title="Save folder name"
							aria-label="Save folder name"
						>
							{#if isSavingTitle}
								<i class="fas fa-spinner fa-spin text-[10px]"></i>
							{:else}
								<i class="fas fa-check text-[10px]"></i>
							{/if}
						</button>
						<button
							onclick={cancelEditTitle}
							class="flex-shrink-0 px-2 py-1 theme-bg-secondary hover:theme-bg-hover theme-text-primary rounded text-xs transition-colors theme-border border"
							title="Cancel rename"
							aria-label="Cancel renaming folder"
						>
							<i class="fas fa-times text-[10px]"></i>
						</button>
					</div>
				{:else}
					<h1
						class="text-sm font-semibold theme-text-primary truncate flex-1 min-w-0 leading-tight"
						title={folder?.name}
					>
						{folder?.name || 'Unnamed Folder'}
					</h1>
					<button
						onclick={startEditTitle}
						class="flex-shrink-0 opacity-0 group-hover:opacity-100 p-1 theme-text-muted hover:theme-text-primary hover:bg-gray-500/10 rounded transition-all"
						title="Rename folder"
						aria-label="Rename this folder"
					>
						<i class="fas fa-pencil-alt text-[10px]"></i>
					</button>
				{/if}
			</div>

			<!-- Right: item count badge -->
			<span
				class="flex-shrink-0 flex items-center gap-1 text-xs theme-text-muted"
				title="{childCount} {childCount === 1 ? 'item' : 'items'} in folder"
				aria-label="{childCount} items"
			>
				<i class="fas fa-layer-group text-[10px]"></i>
				{childCount}
				{childCount === 1 ? 'item' : 'items'}
			</span>
		</div>

		<!-- Active tab indicator -->
		<div class="flex items-center px-4 mt-2">
			<div
				class="flex items-center gap-1.5 px-3 py-2 text-xs font-medium theme-text-primary border-b-2 border-[#ff6600]"
				role="tab"
				aria-selected="true"
			>
				<i class="fas fa-file-alt text-[10px]"></i>
				Documentation
			</div>
		</div>
	</div>

	<!-- ── Content ── -->
	<div class="flex-1 overflow-auto p-5 min-h-0">
		<FolderOverviewTab {folder} />
	</div>
</div>
