<script lang="ts">
	import { onMount } from 'svelte';
	import { actionsApi } from '$lib/api/actionsApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { toast } from '$lib/stores/toast';
	import type { Action } from '$lib/types/Action';

	import ActionItem from './ActionItem.svelte';
	import ActionWizard from './ActionWizard.svelte';
	import ActionsSidebar from './ActionsSidebar.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	let actions: Action[] = [];
	let isLoading = true;
	let error: string | null = null;
	let showEditor = false;
	let editingAction: Action | null = null;
	let draggedIndex: number | null = null;
	let draggedExecutionPoint: string | null = null;
	let activeActionId: string | null = null;

	// Store references to action elements for scrolling
	let actionElements: Record<string, HTMLElement> = {};

	// Group actions by execution point
	$: beforeRequestActions = actions
		.map((action, index) => ({ action, originalIndex: index }))
		.filter(({ action }) => action.execution_point === 'before_request');

	$: afterRequestActions = actions
		.map((action, index) => ({ action, originalIndex: index }))
		.filter(({ action }) => action.execution_point === 'after_request');

	// Load actions when project changes
	$: if ($selectedProject && $selectedWorkspace) {
		loadActions();
	}

	async function loadActions() {
		if (!$selectedWorkspace || !$selectedProject) return;

		try {
			isLoading = true;
			error = null;
			const response = await actionsApi.listActions($selectedWorkspace.id, $selectedProject.id);
			actions = response.data || [];
		} catch (err: any) {
			error = err.message || 'Failed to load actions';
			toast.error(err);
		} finally {
			isLoading = false;
		}
	}

	function handleCreateAction() {
		editingAction = null;
		showEditor = true;
	}

	function handleEditAction(action: Action) {
		editingAction = action;
		showEditor = true;
	}

	async function handleDeleteAction(actionId: string) {
		if (!$selectedWorkspace || !$selectedProject) return;
		if (!confirm('Are you sure you want to delete this action?')) return;

		try {
			await actionsApi.deleteAction($selectedWorkspace.id, $selectedProject.id, actionId);
			toast.success('Action deleted successfully');
			await loadActions();
		} catch (err: any) {
			toast.error(err);
		}
	}

	async function handleToggleAction(action: Action) {
		if (!$selectedWorkspace || !$selectedProject) return;

		try {
			await actionsApi.toggleAction($selectedWorkspace.id, $selectedProject.id, action.id);
			toast.success(`Action ${action.enabled ? 'disabled' : 'enabled'} successfully`);
			await loadActions();
		} catch (err: any) {
			toast.error(err);
		}
	}

	function handleEditorCancel() {
		showEditor = false;
		editingAction = null;
	}

	function handleEditorSave() {
		showEditor = false;
		editingAction = null;
		loadActions();
	}

	// Drag & Drop handlers
	function handleDragStart(action: Action, groupIndex: number) {
		draggedIndex = groupIndex;
		draggedExecutionPoint = action.execution_point;
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
	}

	async function handleDrop(event: DragEvent, dropAction: Action, groupDropIndex: number) {
		event.preventDefault();

		if (draggedIndex === null || draggedExecutionPoint === null) {
			return;
		}

		// Prevent dropping across execution point boundaries
		if (draggedExecutionPoint !== dropAction.execution_point) {
			toast.error('Cannot move actions between Before Request and After Request groups');
			draggedIndex = null;
			draggedExecutionPoint = null;
			return;
		}

		// If dropping at the same position, do nothing
		if (draggedIndex === groupDropIndex) {
			draggedIndex = null;
			draggedExecutionPoint = null;
			return;
		}

		if (!$selectedWorkspace || !$selectedProject) return;

		// Get the actions in the same execution point group
		const sameGroupActions =
			dropAction.execution_point === 'before_request' ? beforeRequestActions : afterRequestActions;

		const draggedAction = sameGroupActions[draggedIndex].action;

		// Reset drag state immediately
		const executionPoint = draggedExecutionPoint;
		draggedIndex = null;
		draggedExecutionPoint = null;

		// Update priority via backend (using 1-based priority)
		// Convert 0-based index to 1-based priority
		try {
			await actionsApi.updateActionPriority(
				$selectedWorkspace.id,
				$selectedProject.id,
				draggedAction.id,
				groupDropIndex + 1 // Convert to 1-based priority
			);

			// Reload to get the correct order from backend
			await loadActions();
			toast.success('Action priority updated successfully');
		} catch (err: any) {
			toast.error(err);
			// Reload to restore correct order
			await loadActions();
		}
	}

	function handleDragEnd() {
		draggedIndex = null;
		draggedExecutionPoint = null;
	}

	// Handle reorder from sidebar - with execution point validation
	async function handleSidebarReorder(
		actionId: string,
		executionPoint: string,
		groupIndex: number
	) {
		if (!$selectedWorkspace || !$selectedProject) return;

		try {
			// Convert 0-based index to 1-based priority
			await actionsApi.updateActionPriority(
				$selectedWorkspace.id,
				$selectedProject.id,
				actionId,
				groupIndex + 1 // Convert to 1-based priority
			);

			// Reload to get the correct order from backend
			await loadActions();
			toast.success('Action priority updated successfully');
		} catch (err: any) {
			toast.error(err);
			// Reload to restore correct order
			await loadActions();
		}
	}

	// Handle click from sidebar - scroll to action
	function handleSidebarActionClick(actionId: string) {
		activeActionId = actionId;
		const element = actionElements[actionId];
		if (element) {
			element.scrollIntoView({ behavior: 'smooth', block: 'center' });
			// Flash highlight effect
			setTimeout(() => {
				activeActionId = null;
			}, 2000);
		}
	}

	onMount(() => {
		loadActions();
	});
</script>

{#if showEditor}
	<!-- Full-page Wizard View -->
	<ActionWizard action={editingAction} onCancel={handleEditorCancel} onSave={handleEditorSave} />
{:else}
	<!-- List View -->
	<div class="w-full h-full theme-bg-primary flex flex-col">
		<!-- Header -->
		<div class="px-4 pt-4 pb-3 border-b theme-border">
			<div class="flex justify-between items-center">
				<div class="flex items-center">
					<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
						<i class="fas fa-cogs text-blue-500 text-xl"></i>
					</div>
					<div>
						<h2 class="text-xl font-bold theme-text-primary">Actions</h2>
						<p class="text-sm theme-text-muted">Automate request and response modifications</p>
					</div>
				</div>

				<div class="flex items-center space-x-3">
					<button
						class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md text-sm flex items-center shadow-sm hover:shadow-md transition-all duration-200"
						on:click={handleCreateAction}
						disabled={!$selectedProject}
						title="Create new action"
						aria-label="Create new action"
					>
						<i class="fas fa-plus mr-2"></i>
						New Action
					</button>
				</div>
			</div>
		</div>

		<!-- Main Content with Sidebar -->
		<div class="flex-1 flex overflow-hidden">
			<!-- Actions List (Main Content) -->
			<div class="flex-1 overflow-auto p-4">
			{#if isLoading}
				<SkeletonLoader type="card" count={3} />
			{:else if error}
				<ErrorDisplay message={error} type="error" retryable={true} onRetry={loadActions} />
			{:else if actions.length === 0}
				<div class="flex flex-col items-center justify-center h-full theme-text-secondary">
					<i class="fas fa-bolt text-6xl mb-4 opacity-50"></i>
					<h3 class="text-xl font-semibold mb-2">No Actions Yet</h3>
					<p class="text-sm mb-4 text-center max-w-md">
						Create your first action to automatically modify requests or responses based on your
						rules.
					</p>
					<button
						class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md text-sm"
						on:click={handleCreateAction}
						title="Create your first action"
						aria-label="Create your first action"
					>
						<i class="fas fa-plus mr-2"></i>
						Create Action
					</button>
				</div>
			{:else}
				<div class="space-y-0">
					<!-- Before Request Actions -->
					{#if beforeRequestActions.length > 0}
						{#each beforeRequestActions as { action, originalIndex }, idx (action.id)}
							<!-- Action Item -->
							<div
								bind:this={actionElements[action.id]}
								role="button"
								tabindex="0"
								draggable="true"
								on:dragstart={() => handleDragStart(action, idx)}
								on:dragover={handleDragOver}
								on:drop={(e) => handleDrop(e, action, idx)}
								on:dragend={handleDragEnd}
								on:keydown={(e) => {
									if (e.key === 'Enter' || e.key === ' ') {
										e.preventDefault();
									}
								}}
								class="transition-all duration-200 cursor-move {activeActionId === action.id
									? 'ring-2 ring-blue-500 dark:ring-blue-400 rounded-lg'
									: ''}"
								class:opacity-50={draggedIndex === idx && draggedExecutionPoint === 'before_request'}
								title="Drag to reorder action within Before Request group"
								aria-label="Drag to reorder action: {action.name || action.type}"
							>
								<ActionItem
									{action}
									onEdit={() => handleEditAction(action)}
									onDelete={() => handleDeleteAction(action.id)}
									onToggle={() => handleToggleAction(action)}
								/>
							</div>

							<!-- Flow Arrow Between Actions -->
							{#if idx < beforeRequestActions.length - 1}
								<div class="flex items-center justify-center py-2" aria-hidden="true">
									<div class="flex flex-col items-center">
										<div class="w-0.5 h-3 bg-gradient-to-b from-blue-400 to-blue-500 dark:from-blue-500 dark:to-blue-600"></div>
										<div class="flex items-center justify-center w-6 h-6 rounded-full bg-blue-500 dark:bg-blue-600 shadow-sm">
											<i class="fas fa-chevron-down text-white text-xs"></i>
										</div>
										<div class="w-0.5 h-3 bg-gradient-to-b from-blue-500 to-blue-400 dark:from-blue-600 dark:to-blue-500"></div>
									</div>
								</div>
							{/if}
						{/each}
					{/if}

					<!-- Main Separator: Before → Server → After -->
					{#if beforeRequestActions.length > 0 && afterRequestActions.length > 0}
						<div class="flex items-center justify-center py-4" aria-hidden="true">
							<div class="flex flex-col items-center gap-1 max-w-md mx-auto">
								<!-- Top connector -->
								<div class="w-0.5 h-4 bg-gradient-to-b from-blue-400 to-purple-500 dark:from-blue-500 dark:to-purple-600"></div>

								<!-- Server Processing Badge -->
								<div class="flex items-center gap-3 px-4 py-2 bg-purple-50 dark:bg-purple-900/20 rounded-lg border border-purple-200 dark:border-purple-800 shadow-sm">
									<div class="flex-1 h-px bg-gradient-to-r from-transparent via-purple-300 dark:via-purple-600 to-purple-500"></div>
									<div class="flex items-center gap-2">
										<div class="flex items-center justify-center w-7 h-7 rounded-full bg-purple-500 dark:bg-purple-600 shadow-md">
											<i class="fas fa-server text-white text-xs"></i>
										</div>
										<span class="text-xs font-semibold text-purple-700 dark:text-purple-300 uppercase tracking-wide whitespace-nowrap">
											Server Processing
										</span>
									</div>
									<div class="flex-1 h-px bg-gradient-to-r from-purple-500 via-purple-300 dark:via-purple-600 to-transparent"></div>
								</div>

								<!-- Bottom connector -->
								<div class="w-0.5 h-4 bg-gradient-to-b from-purple-500 to-green-400 dark:from-purple-600 dark:to-green-500"></div>
							</div>
						</div>
					{/if}

					<!-- After Request Actions -->
					{#if afterRequestActions.length > 0}
						{#each afterRequestActions as { action, originalIndex }, idx (action.id)}
							<!-- Action Item -->
							<div
								bind:this={actionElements[action.id]}
								role="button"
								tabindex="0"
								draggable="true"
								on:dragstart={() => handleDragStart(action, idx)}
								on:dragover={handleDragOver}
								on:drop={(e) => handleDrop(e, action, idx)}
								on:dragend={handleDragEnd}
								on:keydown={(e) => {
									if (e.key === 'Enter' || e.key === ' ') {
										e.preventDefault();
									}
								}}
								class="transition-all duration-200 cursor-move {activeActionId === action.id
									? 'ring-2 ring-green-500 dark:ring-green-400 rounded-lg'
									: ''}"
								class:opacity-50={draggedIndex === idx && draggedExecutionPoint === 'after_request'}
								title="Drag to reorder action within After Request group"
								aria-label="Drag to reorder action: {action.name || action.type}"
							>
								<ActionItem
									{action}
									onEdit={() => handleEditAction(action)}
									onDelete={() => handleDeleteAction(action.id)}
									onToggle={() => handleToggleAction(action)}
								/>
							</div>

							<!-- Flow Arrow Between Actions -->
							{#if idx < afterRequestActions.length - 1}
								<div class="flex items-center justify-center py-2" aria-hidden="true">
									<div class="flex flex-col items-center">
										<div class="w-0.5 h-3 bg-gradient-to-b from-green-400 to-green-500 dark:from-green-500 dark:to-green-600"></div>
										<div class="flex items-center justify-center w-6 h-6 rounded-full bg-green-500 dark:bg-green-600 shadow-sm">
											<i class="fas fa-chevron-down text-white text-xs"></i>
										</div>
										<div class="w-0.5 h-3 bg-gradient-to-b from-green-500 to-green-400 dark:from-green-600 dark:to-green-500"></div>
									</div>
								</div>
							{/if}
						{/each}
					{/if}
				</div>
			{/if}
			</div>

			<!-- Sidebar (Right) -->
			{#if !isLoading && !error && actions.length > 0}
				<ActionsSidebar
					{actions}
					{activeActionId}
					onReorder={handleSidebarReorder}
					onActionClick={handleSidebarActionClick}
				/>
			{/if}
		</div>
	</div>
{/if}
