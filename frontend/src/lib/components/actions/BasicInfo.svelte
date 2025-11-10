<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let name: string = 'New Action';
	export let executionPoint: 'before_request' | 'after_request' = 'after_request';
	export let enabled: boolean = true;

	const dispatch = createEventDispatcher();

	// Dispatch changes to parent
	$: dispatch('change', { name, executionPoint, enabled });
</script>

<div class="flex flex-col gap-4">
	<!-- Line 1: Name and Status -->
	<div class="flex items-end gap-3">
		<!-- Action Name - Takes most space -->
		<div class="flex-1 min-w-0">
			<label for="action-name" class="block text-sm font-medium theme-text-primary mb-1.5">
				Action Name <span class="text-red-500">*</span>
			</label>
			<input
				id="action-name"
				type="text"
				bind:value={name}
				class="block w-full px-3.5 py-2.5 text-sm rounded-lg theme-bg-secondary border border-gray-300 dark:border-gray-600 theme-text-primary placeholder:text-gray-500 dark:placeholder:text-gray-400 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
				placeholder="Enter action name..."
				required
				aria-required="true"
			/>
		</div>

		<!-- Status Toggle - Compact on right -->
		<div class="flex-shrink-0">
			<div class="block text-sm font-medium theme-text-primary mb-1.5">
				Status
			</div>
			<button
				type="button"
				class="px-4 py-2.5 h-[42px] rounded-lg border-2 transition-all duration-150 flex items-center gap-2.5 min-w-[120px] {enabled
					? 'border-green-500 bg-green-500/10 dark:bg-green-500/20'
					: 'border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800'}"
				on:click={() => (enabled = !enabled)}
				title={enabled ? 'Click to disable' : 'Click to enable'}
				aria-label={enabled ? 'Action is enabled' : 'Action is disabled'}
				aria-pressed={enabled}
			>
				<div class="relative w-9 h-5 rounded-full transition-all duration-200 {enabled ? 'bg-green-500' : 'bg-gray-400 dark:bg-gray-600'}">
					<div class="absolute top-0.5 transition-transform duration-200 {enabled ? 'left-[20px]' : 'left-0.5'}">
						<div class="w-4 h-4 rounded-full bg-white shadow-md"></div>
					</div>
				</div>
				<span class="text-sm font-semibold whitespace-nowrap {enabled ? 'text-green-600 dark:text-green-400' : 'theme-text-secondary'}">
					{enabled ? 'Enabled' : 'Disabled'}
				</span>
			</button>
		</div>
	</div>

	<!-- Line 2: When to Execute -->
	<div>
		<div class="block text-sm font-medium theme-text-primary mb-1.5">
			When to Execute <span class="text-red-500">*</span>
		</div>
		<div class="grid grid-cols-2 gap-2.5">
			<button
				type="button"
				class="px-3 py-2 rounded-md border-2 transition-all duration-150 text-left hover:shadow-sm {executionPoint === 'before_request'
					? 'border-blue-500 bg-blue-500/10 dark:bg-blue-500/20 shadow-sm'
					: 'border-gray-300 dark:border-gray-600 hover:border-blue-400 dark:hover:border-blue-500 bg-white dark:bg-gray-800'}"
				on:click={() => (executionPoint = 'before_request')}
				title="Execute before forwarding the request to the server"
				aria-label="Execute before request"
				aria-pressed={executionPoint === 'before_request'}
			>
				<div class="flex items-center gap-1.5 mb-0.5">
					<i class="fas fa-arrow-right text-xs {executionPoint === 'before_request' ? 'text-blue-500' : 'text-gray-500 dark:text-gray-400'}"></i>
					<span class="text-xs font-semibold {executionPoint === 'before_request' ? 'text-blue-600 dark:text-blue-400' : 'theme-text-primary'}">Before Request</span>
				</div>
				<p class="text-xs theme-text-secondary leading-tight">Modify request data</p>
			</button>

			<button
				type="button"
				class="px-3 py-2 rounded-md border-2 transition-all duration-150 text-left hover:shadow-sm {executionPoint === 'after_request'
					? 'border-green-500 bg-green-500/10 dark:bg-green-500/20 shadow-sm'
					: 'border-gray-300 dark:border-gray-600 hover:border-green-400 dark:hover:border-green-500 bg-white dark:bg-gray-800'}"
				on:click={() => (executionPoint = 'after_request')}
				title="Execute after receiving the response from the server"
				aria-label="Execute after request"
				aria-pressed={executionPoint === 'after_request'}
			>
				<div class="flex items-center gap-1.5 mb-0.5">
					<i class="fas fa-arrow-left text-xs {executionPoint === 'after_request' ? 'text-green-500' : 'text-gray-500 dark:text-gray-400'}"></i>
					<span class="text-xs font-semibold {executionPoint === 'after_request' ? 'text-green-600 dark:text-green-400' : 'theme-text-primary'}">After Request</span>
				</div>
				<p class="text-xs theme-text-secondary leading-tight">Modify response data</p>
			</button>
		</div>
	</div>
</div>
