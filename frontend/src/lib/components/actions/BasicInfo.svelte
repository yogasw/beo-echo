<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let name: string = '';
	export let executionPoint: 'before_request' | 'after_request' = 'after_request';
	export let enabled: boolean = true;

	const dispatch = createEventDispatcher();

	// Dispatch changes to parent
	$: dispatch('change', { name, executionPoint, enabled });
</script>

<div class="space-y-4">
	<!-- Action Name -->
	<div>
		<label for="action-name" class="block text-sm font-medium theme-text-primary mb-2">
			Action Name <span class="text-red-500">*</span>
		</label>
		<input
			id="action-name"
			type="text"
			bind:value={name}
			class="block w-full p-3 text-sm rounded-lg theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
			placeholder="e.g., Replace API URL"
			required
			aria-required="true"
		/>
		<p class="mt-1 text-xs theme-text-secondary">
			A descriptive name for this action
		</p>
	</div>

	<!-- Execution Point -->
	<div>
		<label for="execution-point" class="block text-sm font-medium theme-text-primary mb-2">
			When to Execute <span class="text-red-500">*</span>
		</label>
		<select
			id="execution-point"
			bind:value={executionPoint}
			class="block w-full p-3 text-sm rounded-lg theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
			required
			aria-required="true"
		>
			<option value="after_request">After Request (Modify Response)</option>
			<option value="before_request">Before Request (Modify Request)</option>
		</select>
		<p class="mt-1 text-xs theme-text-secondary">
			{executionPoint === 'after_request'
				? 'Executes after receiving the response, allowing you to modify response data'
				: 'Executes before forwarding the request, allowing you to modify request data'}
		</p>
	</div>

	<!-- Enabled Status -->
	<div>
		<label for="enabled" class="block text-sm font-medium theme-text-primary mb-2">
			Status
		</label>
		<label class="relative inline-flex items-center cursor-pointer">
			<input
				id="enabled"
				type="checkbox"
				bind:checked={enabled}
				class="sr-only peer"
				aria-label="Enable or disable this action"
			/>
			<div
				class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
			></div>
			<span class="ml-3 text-sm font-medium theme-text-primary">
				{enabled ? 'Enabled' : 'Disabled'}
			</span>
		</label>
		<p class="mt-1 text-xs theme-text-secondary">
			{enabled ? 'This action will run when conditions are met' : 'This action is currently disabled'}
		</p>
	</div>
</div>
