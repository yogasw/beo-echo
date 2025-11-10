<script lang="ts">
	import type { Action, ReplaceTextConfig, JavaScriptConfig } from '$lib/types/Action';
	import ReplacetextItem from './modules/ReplaceText/ReplacetextItem.svelte';
	import JavaScriptItem from './modules/JavaScript/JavaScriptItem.svelte';
	import { getActionTypeInfo } from '$lib/utils/actionTypeUtils';

	export let action: Action;
	export let onEdit: () => void;
	export let onDelete: () => void;
	export let onToggle: () => void;

	let config: ReplaceTextConfig | JavaScriptConfig | null = null;

	// Parse config
	$: {
		try {
			config = JSON.parse(action.config);
		} catch {
			config = null;
		}
	}

	// Get execution point label
	function getExecutionPointLabel(point: string): string {
		return point === 'before_request' ? 'Before Request' : 'After Request';
	}

	// Get execution point color
	function getExecutionPointColor(point: string): string {
		return point === 'before_request' ? 'bg-purple-600' : 'bg-green-600';
	}

	$: actionTypeInfo = getActionTypeInfo(action.type);
</script>

<div
	class="theme-bg-secondary border theme-border rounded-lg p-4 hover:shadow-lg transition-shadow"
	class:opacity-50={!action.enabled}
>
	<div class="flex items-start justify-between mb-3">
		<div class="flex items-center space-x-3 flex-1">
			<!-- Toggle Switch -->
			<label class="relative inline-flex items-center cursor-pointer" title="Toggle action">
				<input
					type="checkbox"
					checked={action.enabled}
					on:change={onToggle}
					class="sr-only peer"
					aria-label="Toggle action enabled/disabled"
				/>
				<div
					class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
				></div>
			</label>

			<!-- Action Name & Type -->
			<div class="flex-1">
				<div class="flex items-center space-x-2 mb-1">
					<h3 class="text-lg font-semibold theme-text-primary">{action.name}</h3>
					<span
						class="text-xs px-2 py-1 rounded {getExecutionPointColor(
							action.execution_point
						)} text-white"
					>
						{getExecutionPointLabel(action.execution_point)}
					</span>
					{#if action.priority > 0}
						<span
							class="text-xs px-2 py-1 rounded bg-gray-600 text-white"
							title="Execution priority"
							aria-label="Execution priority: {action.priority}"
						>
							Priority: {action.priority}
						</span>
					{/if}
				</div>
				<div class="flex items-center space-x-2 text-sm theme-text-secondary">
					<i class="{actionTypeInfo.iconClass} {actionTypeInfo.icon} {actionTypeInfo.color}"></i>
					<span>{actionTypeInfo.label}</span>
				</div>
			</div>
		</div>

		<!-- Action Buttons -->
		<div class="flex space-x-2">
			<button
				class="text-xs bg-gray-700 hover:bg-gray-600 text-gray-300 px-3 py-1 rounded"
				on:click={onEdit}
				title="Edit action"
				aria-label="Edit action"
			>
				<i class="fas fa-edit mr-1"></i>
				Edit
			</button>
			<button
				class="text-xs bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded"
				on:click={onDelete}
				title="Delete action"
				aria-label="Delete action"
			>
				<i class="fas fa-trash-alt mr-1"></i>
				Delete
			</button>
		</div>
	</div>

	<!-- Config Details - Dynamic based on action type -->
	{#if config}
		{#if action.type === 'replace_text'}
			<ReplacetextItem config={config as ReplaceTextConfig} />
		{:else if action.type === 'run_javascript'}
			<JavaScriptItem config={config as JavaScriptConfig} />
		{:else}
			<!-- Other action types: no preview (empty) -->
			<div
				class="mt-2 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-md border border-gray-200 dark:border-gray-700"
			>
				<p class="text-xs theme-text-secondary italic text-center">
					<i class="fas fa-info-circle mr-1"></i>
					Preview not available for this action type
				</p>
			</div>
		{/if}
	{/if}

	<!-- Filters -->
	{#if action.filters && action.filters.length > 0}
		<div class="mt-3">
			<div class="text-sm theme-text-secondary mb-2">
				<i class="fas fa-filter mr-1"></i>
				Filters ({action.filters.length}):
			</div>
			<div class="flex flex-wrap gap-2">
				{#each action.filters as filter}
					<span
						class="text-xs px-2 py-1 rounded bg-blue-900/50 border border-blue-700 theme-text-primary"
					>
						{filter.type}: {filter.operator} "{filter.value}"
					</span>
				{/each}
			</div>
		</div>
	{/if}
</div>
