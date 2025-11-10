<script lang="ts">
	import type { Action } from '$lib/types/Action';
	import { toast } from '$lib/stores/toast';
	import { getActionTypeInfo } from '$lib/utils/actionTypeUtils';

	export let actions: Action[];
	export let onReorder: (actionId: string, executionPoint: string, groupIndex: number) => void;
	export let onActionClick: (actionId: string) => void;
	export let activeActionId: string | null = null;

	let draggedActionId: string | null = null;
	let draggedExecutionPoint: string | null = null;
	let draggedGroupIndex: number | null = null;

	// Group actions by execution point
	$: beforeRequestActions = actions
		.map((action, index) => ({ action, originalIndex: index }))
		.filter(({ action }) => action.execution_point === 'before_request');

	$: afterRequestActions = actions
		.map((action, index) => ({ action, originalIndex: index }))
		.filter(({ action }) => action.execution_point === 'after_request');

	function handleDragStart(action: Action, groupIndex: number) {
		draggedActionId = action.id;
		draggedExecutionPoint = action.execution_point;
		draggedGroupIndex = groupIndex;
	}

	function handleDragOver(event: DragEvent) {
		event.preventDefault();
	}

	function handleDrop(event: DragEvent, dropAction: Action, groupDropIndex: number) {
		event.preventDefault();

		if (draggedActionId === null || draggedExecutionPoint === null || draggedGroupIndex === null) {
			return;
		}

		// Prevent dropping across execution point boundaries
		if (draggedExecutionPoint !== dropAction.execution_point) {
			toast.error('Cannot move actions between Before Request and After Request groups');
			draggedActionId = null;
			draggedExecutionPoint = null;
			draggedGroupIndex = null;
			return;
		}

		// If dropping at the same position, do nothing
		if (draggedGroupIndex === groupDropIndex) {
			draggedActionId = null;
			draggedExecutionPoint = null;
			draggedGroupIndex = null;
			return;
		}

		onReorder(draggedActionId, draggedExecutionPoint, groupDropIndex);
		draggedActionId = null;
		draggedExecutionPoint = null;
		draggedGroupIndex = null;
	}

	function handleDragEnd() {
		draggedActionId = null;
		draggedExecutionPoint = null;
		draggedGroupIndex = null;
	}

	function handleClick(actionId: string) {
		onActionClick(actionId);
	}
</script>

<div class="w-64 theme-bg-secondary border-l theme-border flex flex-col h-full overflow-hidden">
	<!-- Sidebar Header -->
	<div class="px-4 py-3 border-b theme-border">
		<div class="flex items-center gap-2">
			<i class="fas fa-list-ol text-blue-500 dark:text-blue-400"></i>
			<h3 class="text-sm font-semibold theme-text-primary">Actions Flow</h3>
		</div>
		<p class="text-xs theme-text-secondary mt-1">Drag to reorder • Click to jump</p>
	</div>

	<!-- Actions List -->
	<div class="flex-1 overflow-y-auto p-3 space-y-3">
		<!-- Before Request Section -->
		{#if beforeRequestActions.length > 0}
			<div>
				<!-- Section Header -->
				<div class="flex items-center gap-2 mb-2.5 px-2 py-1.5 bg-blue-50 dark:bg-blue-900/20 rounded-md border border-blue-200 dark:border-blue-800">
					<div class="flex items-center justify-center w-5 h-5 rounded bg-blue-500 dark:bg-blue-600 flex-shrink-0">
						<i class="fas fa-arrow-right text-white" style="font-size: 9px;"></i>
					</div>
					<span class="text-xs font-bold theme-text-primary uppercase tracking-wide flex-1">
						Before Request
					</span>
					<span class="text-xs font-semibold text-blue-600 dark:text-blue-400 bg-blue-100 dark:bg-blue-900/40 px-1.5 py-0.5 rounded">
						{beforeRequestActions.length}
					</span>
				</div>

				<!-- Action Items -->
				<div class="space-y-1">
					{#each beforeRequestActions as { action, originalIndex }, idx (action.id)}
					{@const typeInfo = getActionTypeInfo(action.type)}
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							draggable="true"
							on:dragstart={() => handleDragStart(action, idx)}
							on:dragover={handleDragOver}
							on:drop={(e) => handleDrop(e, action, idx)}
							on:dragend={handleDragEnd}
							on:click={() => handleClick(action.id)}
							on:keydown={(e) => e.key === 'Enter' && handleClick(action.id)}
							class="group px-2 py-2 rounded cursor-pointer transition-all duration-150 {activeActionId === action.id
								? 'bg-blue-100 dark:bg-blue-900/40 border-l-2 border-blue-500'
								: 'hover:bg-gray-200 dark:hover:bg-gray-700/50 border-l-2 border-transparent'}"
							class:opacity-50={draggedActionId === action.id && draggedExecutionPoint === 'before_request'}
							title="{action.name || action.type} - Click to jump"
						>
							<div class="flex items-center gap-2 min-w-0">
								<i class="{typeInfo.iconClass} {typeInfo.icon} text-xs {typeInfo.color} flex-shrink-0"></i>
								<span class="text-xs font-medium theme-text-primary truncate flex-1">
									{action.name || action.type}
								</span>
								{#if !action.enabled}
									<i class="fas fa-pause-circle text-xs text-gray-400 flex-shrink-0" title="Disabled"></i>
								{/if}
							</div>
						</div>

						<!-- Mini Arrow Between Actions -->
						<!-- {#if idx < beforeRequestActions.length - 1}
							<div class="flex justify-center py-0.5" aria-hidden="true">
								<i class="fas fa-chevron-down text-blue-400 dark:text-blue-500" style="font-size: 8px;"></i>
							</div>
						{/if} -->
					{/each}
				</div>
			</div>
		{/if}

		<!-- Separator between Before and After -->
		{#if beforeRequestActions.length > 0 && afterRequestActions.length > 0}
			<div class="flex items-center justify-center gap-2 py-2 px-2" aria-hidden="true">
				<!-- Left line -->
				<div class="flex-1 h-0.5 bg-gradient-to-r from-blue-400 to-purple-500 dark:from-blue-500 dark:to-purple-600"></div>
				<!-- Server Icon Circle -->
				<div class="flex items-center justify-center w-6 h-6 rounded-full bg-purple-500 dark:bg-purple-600 shadow-sm flex-shrink-0">
					<i class="fas fa-server text-white" style="font-size: 9px;"></i>
				</div>
				<!-- Right line -->
				<div class="flex-1 h-0.5 bg-gradient-to-r from-purple-500 to-green-400 dark:from-purple-600 dark:to-green-500"></div>
			</div>
		{/if}

		<!-- After Request Section -->
		{#if afterRequestActions.length > 0}
			<div>
				<!-- Section Header -->
				<div class="flex items-center gap-2 mb-2.5 px-2 py-1.5 bg-green-50 dark:bg-green-900/20 rounded-md border border-green-200 dark:border-green-800">
					<div class="flex items-center justify-center w-5 h-5 rounded bg-green-500 dark:bg-green-600 flex-shrink-0">
						<i class="fas fa-arrow-left text-white" style="font-size: 9px;"></i>
					</div>
					<span class="text-xs font-bold theme-text-primary uppercase tracking-wide flex-1">
						After Request
					</span>
					<span class="text-xs font-semibold text-green-600 dark:text-green-400 bg-green-100 dark:bg-green-900/40 px-1.5 py-0.5 rounded">
						{afterRequestActions.length}
					</span>
				</div>

				<!-- Action Items -->
				<div class="space-y-1">
					{#each afterRequestActions as { action, originalIndex }, idx (action.id)}
						{@const typeInfo = getActionTypeInfo(action.type)}
						<!-- svelte-ignore a11y-no-static-element-interactions -->
						<div
							draggable="true"
							on:dragstart={() => handleDragStart(action, idx)}
							on:dragover={handleDragOver}
							on:drop={(e) => handleDrop(e, action, idx)}
							on:dragend={handleDragEnd}
							on:click={() => handleClick(action.id)}
							on:keydown={(e) => e.key === 'Enter' && handleClick(action.id)}
							class="group px-2 py-2 rounded cursor-pointer transition-all duration-150 {activeActionId === action.id
								? 'bg-green-100 dark:bg-green-900/40 border-l-2 border-green-500'
								: 'hover:bg-gray-200 dark:hover:bg-gray-700/50 border-l-2 border-transparent'}"
							class:opacity-50={draggedActionId === action.id && draggedExecutionPoint === 'after_request'}
							title="{action.name || action.type} - Click to jump"
						>
							<div class="flex items-center gap-2 min-w-0">
								<i class="{typeInfo.iconClass} {typeInfo.icon} text-xs {typeInfo.color} flex-shrink-0"></i>
								<span class="text-xs font-medium theme-text-primary truncate flex-1">
									{action.name || action.type}
								</span>
								{#if !action.enabled}
									<i class="fas fa-pause-circle text-xs text-gray-400 flex-shrink-0" title="Disabled"></i>
								{/if}
							</div>
						</div>

						<!-- Mini Arrow Between Actions -->
						<!-- {#if idx < afterRequestActions.length - 1}
							<div class="flex justify-center py-0.5" aria-hidden="true">
								<i class="fas fa-chevron-down text-green-400 dark:text-green-500" style="font-size: 8px;"></i>
							</div>
						{/if} -->
					{/each}
				</div>
			</div>
		{/if}

		<!-- Empty State -->
		{#if actions.length === 0}
			<div class="flex flex-col items-center justify-center py-8 px-4 text-center">
				<i class="fas fa-inbox text-3xl theme-text-secondary opacity-50 mb-2"></i>
				<p class="text-xs theme-text-secondary">No actions yet</p>
			</div>
		{/if}
	</div>
</div>
