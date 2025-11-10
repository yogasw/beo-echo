<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { ActionFilter } from '$lib/types/Action';

	export let filters: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[] = [];

	const dispatch = createEventDispatcher();

	// Add new filter
	function addFilter() {
		filters = [
			...filters,
			{
				type: 'method',
				key: '',
				operator: 'equals',
				value: ''
			}
		];
		dispatch('change', filters);
	}

	// Remove filter
	function removeFilter(index: number) {
		filters = filters.filter((_, i) => i !== index);
		dispatch('change', filters);
	}

	// Update filter
	function updateFilter(index: number, field: string, value: any) {
		filters[index] = { ...filters[index], [field]: value };
		filters = [...filters]; // Trigger reactivity
		dispatch('change', filters);
	}

	// Show/hide key field based on filter type
	function needsKey(type: string): boolean {
		return type === 'header' || type === 'query';
	}
</script>

<div class="space-y-3">
	{#if filters.length === 0}
		<div class="text-center p-6 theme-bg-secondary rounded-lg border theme-border border-dashed">
			<i class="fas fa-filter text-3xl theme-text-secondary mb-2"></i>
			<p class="text-sm theme-text-secondary">No filters added yet</p>
			<p class="text-xs theme-text-secondary mt-1">
				Filters allow you to conditionally execute actions based on request/response properties
			</p>
		</div>
	{:else}
		{#each filters as filter, index}
			<div class="p-4 theme-bg-secondary rounded-lg border theme-border space-y-3">
				<div class="flex justify-between items-start">
					<span class="text-sm font-medium theme-text-primary">Filter #{index + 1}</span>
					<button
						type="button"
						class="text-red-500 hover:text-red-400 text-sm"
						on:click={() => removeFilter(index)}
						title="Remove filter"
						aria-label="Remove filter {index + 1}"
					>
						<i class="fas fa-trash-alt"></i>
					</button>
				</div>

				<div class="grid grid-cols-1 md:grid-cols-3 gap-3">
					<!-- Filter Type -->
					<div>
						<label for="filter-type-{index}" class="block text-xs font-medium theme-text-secondary mb-1">
							Type
						</label>
						<select
							id="filter-type-{index}"
							value={filter.type}
							on:change={(e) => updateFilter(index, 'type', e.currentTarget.value)}
							class="block w-full p-2 text-sm rounded theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							aria-label="Filter type for filter {index + 1}"
						>
							<option value="method">Method</option>
							<option value="path">Path</option>
							<option value="header">Header</option>
							<option value="query">Query Param</option>
							<option value="status_code">Status Code</option>
						</select>
					</div>

					<!-- Operator -->
					<div>
						<label for="filter-operator-{index}" class="block text-xs font-medium theme-text-secondary mb-1">
							Operator
						</label>
						<select
							id="filter-operator-{index}"
							value={filter.operator}
							on:change={(e) => updateFilter(index, 'operator', e.currentTarget.value)}
							class="block w-full p-2 text-sm rounded theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							aria-label="Operator for filter {index + 1}"
						>
							<option value="equals">Equals</option>
							<option value="contains">Contains</option>
							<option value="starts_with">Starts With</option>
							<option value="ends_with">Ends With</option>
							<option value="regex">Regex</option>
						</select>
					</div>

					<!-- Value -->
					<div>
						<label for="filter-value-{index}" class="block text-xs font-medium theme-text-secondary mb-1">
							Value
						</label>
						<input
							id="filter-value-{index}"
							type="text"
							value={filter.value}
							on:input={(e) => updateFilter(index, 'value', e.currentTarget.value)}
							class="block w-full p-2 text-sm rounded theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
							placeholder={filter.type === 'method' ? 'GET, POST' : 'Value to match'}
							aria-label="Value for filter {index + 1}"
						/>
					</div>
				</div>

				<!-- Key field (for header and query types) -->
				{#if needsKey(filter.type)}
					<div>
						<label for="filter-key-{index}" class="block text-xs font-medium theme-text-secondary mb-1">
							{filter.type === 'header' ? 'Header Name' : 'Query Parameter Name'}
						</label>
						<input
							id="filter-key-{index}"
							type="text"
							value={filter.key || ''}
							on:input={(e) => updateFilter(index, 'key', e.currentTarget.value)}
							class="block w-full p-2 text-sm rounded theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
							placeholder={filter.type === 'header' ? 'e.g., Content-Type' : 'e.g., user_id'}
							aria-label="{filter.type === 'header' ? 'Header' : 'Query parameter'} name for filter {index + 1}"
						/>
					</div>
				{/if}

				<!-- Filter Example -->
				<div class="text-xs theme-text-secondary bg-blue-50/70 dark:bg-gray-900/50 border border-blue-200/50 dark:border-gray-700/50 p-2.5 rounded">
					<i class="fas fa-info-circle mr-1 text-blue-500"></i>
					Example: This action will run when
					<span class="theme-text-primary font-medium">
						{#if filter.type === 'method'}
							request method {filter.operator} "{filter.value || '...'}"
						{:else if filter.type === 'path'}
							request path {filter.operator} "{filter.value || '...'}"
						{:else if filter.type === 'header'}
							header "{filter.key || '...'}" {filter.operator} "{filter.value || '...'}"
						{:else if filter.type === 'query'}
							query param "{filter.key || '...'}" {filter.operator} "{filter.value || '...'}"
						{:else if filter.type === 'status_code'}
							response status code {filter.operator} "{filter.value || '...'}"
						{/if}
					</span>
				</div>
			</div>
		{/each}
	{/if}

	<!-- Add Filter Button -->
	<button
		type="button"
		class="w-full py-3 px-4 border-2 border-dashed theme-border hover:border-blue-500 rounded-lg text-sm theme-text-secondary hover:theme-text-primary transition-colors flex items-center justify-center"
		on:click={addFilter}
		title="Add new filter"
		aria-label="Add new filter"
	>
		<i class="fas fa-plus mr-2"></i>
		Add Filter
	</button>

	{#if filters.length > 0}
		<div class="text-xs theme-text-secondary bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-900/50 p-3 rounded">
			<i class="fas fa-info-circle mr-1 text-blue-500"></i>
			<strong>Note:</strong> Multiple filters use OR logic - the action runs if ANY filter matches.
			If no filters are set, the action runs for all requests.
		</div>
	{/if}
</div>
