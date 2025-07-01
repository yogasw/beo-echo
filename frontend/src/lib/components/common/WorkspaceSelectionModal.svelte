<script lang="ts">
	import { onMount } from 'svelte';
	import { getWorkspaces, addProject } from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';
	import { goto } from '$app/navigation';
	import { setCurrentWorkspaceId } from '$lib/utils/localStorage';
	import type { Workspace } from '$lib/types/User';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	export let isOpen = false;
	export let projectName = '';
	export let projectAlias = '';
	export let onClose: () => void;

	let workspaces: Workspace[] = [];
	let isLoading = false;
	let isCreating = false;
	let error: string | null = null;
	let selectedWorkspaceId = '';

	// Load workspaces when modal opens
	$: {
		if (isOpen && workspaces.length === 0) {
			loadWorkspaces();
		}
	}

	async function loadWorkspaces() {
		if (!projectName.trim() || !projectAlias.trim()) {
			error = 'Project name and alias are required';
			return;
		}

		isLoading = true;
		error = null;
		
		try {
			workspaces = await getWorkspaces();
			
			// Auto-select first workspace if only one exists
			if (workspaces.length === 1) {
				selectedWorkspaceId = workspaces[0].id;
			}
		} catch (err) {
			console.error('Failed to load workspaces:', err);
			error = 'Failed to load workspaces';
			toast.error(err);
		} finally {
			isLoading = false;
		}
	}

	async function createProjectInWorkspace() {
		if (!selectedWorkspaceId) {
			toast.error('Please select a workspace');
			return;
		}

		if (!projectName.trim() || !projectAlias.trim()) {
			toast.error('Project name and alias are required');
			return;
		}

		isCreating = true;
		try {
			// Set the current workspace ID for the API call
			setCurrentWorkspaceId(selectedWorkspaceId);

			const newProject = await addProject(projectName.trim(), projectAlias.trim());
			
			toast.success(`Project "${projectName}" created successfully!`);
			
			// Navigate to the new project
			await goto(`/home/workspace/${selectedWorkspaceId}/projects/${newProject.id}`);
			
			onClose();
		} catch (err) {
			console.error('Failed to create project:', err);
			toast.error(err);
		} finally {
			isCreating = false;
		}
	}

	function handleWorkspaceSelect(workspaceId: string) {
		selectedWorkspaceId = workspaceId;
	}

	function handleBackdropClick(event: MouseEvent) {
		if (event.target === event.currentTarget) {
			onClose();
		}
	}

	// Close modal on Escape key
	function handleKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') {
			onClose();
		}
	}

	onMount(() => {
		if (isOpen) {
			document.addEventListener('keydown', handleKeydown);
			return () => {
				document.removeEventListener('keydown', handleKeydown);
			};
		}
	});
</script>

{#if isOpen}
	<!-- Modal Backdrop -->
	<div 
		class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center p-4"
		on:click={handleBackdropClick}
		on:keydown={handleKeydown}
		role="dialog"
		aria-labelledby="workspace-modal-title"
		aria-modal="true"
		tabindex="-1"
	>
		<!-- Modal Content -->
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full max-h-[80vh] overflow-hidden">
			<!-- Modal Header -->
			<div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
				<h3 id="workspace-modal-title" class="text-lg font-semibold text-gray-900 dark:text-white">
					<i class="fas fa-folder-plus text-blue-600 dark:text-blue-400 mr-2"></i>
					Select Workspace
				</h3>
				<button
					on:click={onClose}
					class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 p-1 rounded-lg hover:bg-gray-100 dark:hover:bg-gray-700"
					title="Close modal"
					aria-label="Close workspace selection modal"
				>
					<i class="fas fa-times text-lg"></i>
				</button>
			</div>

			<!-- Modal Body -->
			<div class="p-4">
				<!-- Project Info -->
				<div class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-700">
					<p class="text-sm text-gray-700 dark:text-gray-300 mb-1">
						<strong>Project:</strong> {projectName}
					</p>
					<p class="text-sm text-gray-700 dark:text-gray-300">
						<strong>Alias:</strong> <span class="font-mono text-blue-600 dark:text-blue-400">{projectAlias}</span>
					</p>
				</div>

				<!-- Loading State -->
				{#if isLoading}
					<div class="space-y-3">
						<p class="text-sm text-gray-600 dark:text-gray-400 mb-3">Loading your workspaces...</p>
						<SkeletonLoader type="list" count={3} />
					</div>
				{:else if error}
					<!-- Error State -->
					<ErrorDisplay 
						message={error} 
						type="error" 
						retryable={true}
						onRetry={loadWorkspaces}
					/>
				{:else if workspaces.length === 0}
					<!-- No Workspaces -->
					<div class="text-center py-6">
						<i class="fas fa-folder-open text-gray-400 text-3xl mb-3"></i>
						<p class="text-gray-500 dark:text-gray-400 mb-2">No workspaces found</p>
						<p class="text-sm text-gray-400">You need to create a workspace first</p>
					</div>
				{:else}
					<!-- Workspace Selection -->
					<div class="space-y-2">
						<p class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3">
							Choose a workspace for your new project:
						</p>
						
						{#each workspaces as workspace (workspace.id)}
							<button
								on:click={() => handleWorkspaceSelect(workspace.id)}
								class="w-full p-3 text-left rounded-lg border transition-all duration-200 {selectedWorkspaceId === workspace.id 
									? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 ring-2 ring-blue-500 ring-opacity-20' 
									: 'border-gray-200 dark:border-gray-600 hover:border-blue-300 dark:hover:border-blue-500 bg-white dark:bg-gray-700 hover:bg-blue-50 dark:hover:bg-blue-900/10'}"
								title="Select {workspace.name} workspace"
								aria-label="Select {workspace.name} workspace"
							>
								<div class="flex items-center justify-between">
									<div class="flex items-center">
										<div class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white text-sm font-medium mr-3">
											{workspace.name.charAt(0).toUpperCase()}
										</div>
										<div>
											<p class="font-medium text-gray-900 dark:text-white">{workspace.name}</p>
											<p class="text-xs text-gray-500 dark:text-gray-400">
												Workspace
											</p>
										</div>
									</div>
									{#if selectedWorkspaceId === workspace.id}
										<i class="fas fa-check-circle text-blue-500"></i>
									{/if}
								</div>
							</button>
						{/each}
					</div>
				{/if}
			</div>

			<!-- Modal Footer -->
			{#if !isLoading && !error && workspaces.length > 0}
				<div class="flex items-center justify-end gap-3 p-4 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-750">
					<button
						on:click={onClose}
						class="px-4 py-2 text-sm font-medium text-gray-700 dark:text-gray-300 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors"
						title="Cancel project creation"
						aria-label="Cancel project creation"
					>
						Cancel
					</button>
					<button
						on:click={createProjectInWorkspace}
						disabled={!selectedWorkspaceId || isCreating}
						class="px-4 py-2 text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed rounded-lg transition-colors flex items-center"
						title={!selectedWorkspaceId ? "Select a workspace first" : "Create project in selected workspace"}
						aria-label={!selectedWorkspaceId ? "Select a workspace first" : "Create project in selected workspace"}
					>
						{#if isCreating}
							<i class="fas fa-spinner fa-spin mr-2"></i>
							Creating...
						{:else}
							<i class="fas fa-plus mr-2"></i>
							Create Project
						{/if}
					</button>
				</div>
			{/if}
		</div>
	</div>
{/if}

<style>
	/* Ensure modal appears above other content */
	:global(body.modal-open) {
		overflow: hidden;
	}
</style>
