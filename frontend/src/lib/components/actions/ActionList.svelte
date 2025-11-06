<script lang="ts">
	import { onMount } from 'svelte';
	import { actionsApi } from '$lib/api/actionsApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { toast } from '$lib/stores/toast';
	import type { Action, ActionTypeInfo } from '$lib/types/Action';

	import ActionItem from './ActionItem.svelte';
	import ActionWizard from './ActionWizard.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	let actions: Action[] = [];
	let actionTypes: ActionTypeInfo[] = [];
	let isLoading = true;
	let error: string | null = null;
	let showEditor = false;
	let editingAction: Action | null = null;
	let draggedIndex: number | null = null;

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
	function handleDragStart(index: number) {
		draggedIndex = index;
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
	}

	async function handleDrop(event: DragEvent, dropIndex: number) {
		event.preventDefault();

		if (draggedIndex === null || draggedIndex === dropIndex) {
			draggedIndex = null;
			return;
		}

		// Reorder actions array
		const reorderedActions = [...actions];
		const [draggedAction] = reorderedActions.splice(draggedIndex, 1);
		reorderedActions.splice(dropIndex, 0, draggedAction);

		// Update UI immediately
		actions = reorderedActions;
		draggedIndex = null;

		// Update priorities in backend
		try {
			for (let i = 0; i < reorderedActions.length; i++) {
				const action = reorderedActions[i];

				if (!$selectedWorkspace || !$selectedProject) return;
				await actionsApi.updateAction($selectedWorkspace.id, $selectedProject.id, action.id, {
					priority: i
				});
			}
			toast.success('Actions reordered successfully');
		} catch (err: any) {
			toast.error(err);
			// Reload to get correct order
			await loadActions();
		}
	}

	function handleDragEnd() {
		draggedIndex = null;
	}

	onMount(async () => {
		loadActions();
		// Load action types for category grouping
		try {
			const response = await actionsApi.getActionTypes();
			actionTypes = response.data;
		} catch (err) {
			// Silently fail if action types can't be loaded
			console.error('Failed to load action types:', err);
		}
	});
</script>

{#if showEditor}
	<!-- Full-page Wizard View -->
	<ActionWizard action={editingAction} onCancel={handleEditorCancel} onSave={handleEditorSave} />
{:else}
	<!-- List View -->
	<div class="w-full theme-bg-primary p-4 relative">
		<!-- Header -->
		<div class="mb-6">
			<div class="flex justify-between items-center mb-4">
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

		<!-- Content -->
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
				{@const groupedActions = actions.reduce<Record<string, Action[]>>((acc, action) => {
					const actionType = actionTypes.find((t) => t.id === action.type);
					const category = actionType?.category || 'Other';
					if (!acc[category]) acc[category] = [];
					acc[category].push(action);
					return acc;
				}, {})}

				{#each Object.entries(groupedActions) as [category, categoryActions]}
					<div class="mb-6">
						<div class="flex items-center mb-3">
							<div class="flex-1 h-px theme-border"></div>
							<h3 class="px-4 text-sm font-semibold theme-text-secondary uppercase tracking-wide">
								{category}
							</h3>
							<div class="flex-1 h-px theme-border"></div>
						</div>
						<div class="grid grid-cols-1 gap-4">
							{#each categoryActions as action (action.id)}
								<div
									role="button"
									tabindex="0"
									draggable="true"
									on:dragstart={() => handleDragStart(actions.indexOf(action))}
									on:dragover={handleDragOver}
									on:drop={(e) => handleDrop(e, actions.indexOf(action))}
									on:dragend={handleDragEnd}
									on:keydown={(e) => {
										if (e.key === 'Enter' || e.key === ' ') {
											e.preventDefault();
										}
									}}
									class="transition-opacity cursor-move"
									class:opacity-50={draggedIndex === actions.indexOf(action)}
									title="Drag to reorder action"
									aria-label="Drag to reorder action: {action.name || action.type}"
								>
									<ActionItem
										{action}
										onEdit={() => handleEditAction(action)}
										onDelete={() => handleDeleteAction(action.id)}
										onToggle={() => handleToggleAction(action)}
									/>
								</div>
							{/each}
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
{/if}
