<script lang="ts">
	import { onMount } from 'svelte';
	import { workspaces, currentWorkspace, workspaceStore } from '$lib/stores/workspace';
	export let className = '';

	// Local state
	let loading = true;
	let error: string | null = null;
	let modalOpen = false;
	let newWorkspaceName = '';

	// Handle workspace selection
	function selectWorkspace(workspaceId: string) {
		workspaces.setCurrent(workspaceId);
		modalOpen = false;
	}

	// Handle new workspace creation
	async function createWorkspace() {
		if (!newWorkspaceName?.trim()) {
			error = 'Workspace name is required';
			return;
		}

		loading = true;
		error = null;

		try {
			await workspaces.create(newWorkspaceName);
			newWorkspaceName = '';
			modalOpen = false;
		} catch (err) {
			error = err.message || 'Failed to create workspace';
		} finally {
			loading = false;
		}
	}


	// Load workspaces on component mount
	onMount(async () => {
		try {
			await workspaces.loadAll();
		} catch (err) {
			error = 'Failed to load workspaces';
		} finally {
			loading = false;
		}
	});

	// Handle modal toggle
	function toggleModal() {
		modalOpen = !modalOpen;
	}

</script>

<!-- Workspace Manager Component -->
<div class="relative {className}">
	<!-- Company/Workspace Button -->
	<button on:click={toggleModal} class="flex flex-col items-center">
		<div
			class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-green-500 flex items-center justify-center"
		>
			<i class="fas fa-building"></i>
		</div>
		<span class="text-xs mt-1 theme-text-primary">{$currentWorkspace?.name}</span>
	</button>

	<!-- Workspace Modal -->
	{#if modalOpen}
		<div
			class="absolute top-full right-0 mt-2 w-64 theme-bg-primary rounded-md shadow-lg z-40 border theme-border"
		>
			<div class="p-2">
				<!-- Error Message -->
				{#if error}
					<div
						class="mx-3 my-2 p-2 text-xs bg-red-500/10 text-red-400 dark:text-red-300 rounded border border-red-400/30"
					>
						{error}
					</div>
				{/if}

				<!-- Workspaces List -->
				<div class="max-h-48 overflow-y-auto">
					{#if !loading && $workspaceStore.workspaces.length > 0}
						{#each $workspaceStore.workspaces as workspace}
							<button
								on:click={() => selectWorkspace(workspace.id)}
								class="flex items-center w-full px-3 py-2 text-left hover:bg-blue-500/20 rounded-md transition-colors"
							>
								<div class="flex-1">
									<div class="flex items-center gap-2">
										<span class="theme-text-primary">{workspace.name}</span>
										{#if workspace.id === $currentWorkspace?.id}
											<i class="fas fa-check text-blue-500 text-xs"></i>
										{/if}
									</div>
									{#if workspace.role}
										<span class="theme-text-muted text-xs">{workspace.role}</span>
									{/if}
								</div>
							</button>
						{/each}
					{:else if !loading}
						<div class="px-3 py-2 theme-text-muted text-sm">No workspaces found</div>
					{/if}
				</div>

				<!-- Create New Workspace -->
				<div class="mt-2 pt-2 border-t theme-border">
					<div class="px-3 py-2">
						<h4 class="theme-text-secondary text-xs font-medium mb-1">New Workspace</h4>
						<div class="flex items-center gap-2">
							<input
								type="text"
								placeholder="Workspace name"
								bind:value={newWorkspaceName}
								class="flex-1 px-2 py-1 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
							/>
							<button
								on:click={createWorkspace}
								disabled={loading || !newWorkspaceName}
								class="p-1 theme-bg-secondary hover:bg-blue-600 disabled:opacity-50 disabled:hover:bg-gray-700 rounded text-white text-sm"
							>
								<i class="fas fa-plus"></i>
							</button>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>

<!-- Delete Confirmation Modal
{#if showDeleteConfirmation && workspaceToDelete}
	<WorkspaceDeleteConfirmation workspace={workspaceToDelete} onConfirm={handleDeleteConfirmation} />
{/if} -->
