<script lang="ts">
	import type { Action, ReplaceTextConfig } from '$lib/types/Action';

	export let action: Action;
	export let onEdit: () => void;
	export let onDelete: () => void;
	export let onToggle: () => void;

	let config: ReplaceTextConfig | null = null;

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

	// Get target label for replace_text
	function getTargetLabel(target: string): string {
		const labels: Record<string, string> = {
			request_body: 'Request Body',
			response_body: 'Response Body',
			request_header: 'Request Header',
			response_header: 'Response Header'
		};
		return labels[target] || target;
	}
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
						<span class="text-xs px-2 py-1 rounded bg-gray-600 text-white" title="Execution priority">
							Priority: {action.priority}
						</span>
					{/if}
				</div>
				<div class="flex items-center space-x-2 text-sm theme-text-secondary">
					<i class="fas fa-bolt text-amber-500"></i>
					<span>Replace Text</span>
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

	<!-- Config Details -->
	{#if config}
		<div class="mt-3 p-3 bg-gray-900/50 rounded border theme-border">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
				<div>
					<span class="theme-text-secondary">Target:</span>
					<span class="theme-text-primary font-medium ml-2">{getTargetLabel(config.target)}</span>
					{#if config.header_key}
						<span class="theme-text-secondary ml-2">({config.header_key})</span>
					{/if}
				</div>
				<div>
					<span class="theme-text-secondary">Mode:</span>
					<span class="theme-text-primary font-medium ml-2">
						{config.use_regex ? 'Regex' : 'Text'}
					</span>
				</div>
				<div class="md:col-span-2">
					<span class="theme-text-secondary">Pattern:</span>
					<code
						class="ml-2 px-2 py-1 bg-gray-800 rounded text-xs theme-text-primary font-mono"
					>
						{config.pattern}
					</code>
				</div>
				<div class="md:col-span-2">
					<span class="theme-text-secondary">Replacement:</span>
					<code
						class="ml-2 px-2 py-1 bg-gray-800 rounded text-xs theme-text-primary font-mono"
					>
						{config.replacement}
					</code>
				</div>
			</div>
		</div>
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
					<span class="text-xs px-2 py-1 rounded bg-blue-900/50 border border-blue-700 theme-text-primary">
						{filter.type}: {filter.operator} "{filter.value}"
					</span>
				{/each}
			</div>
		</div>
	{/if}
</div>
