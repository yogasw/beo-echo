<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let params: Array<{ key: string; value: string; description: string; enabled: boolean }> = [
		{ key: '', value: '', description: '', enabled: true }
	];

	const dispatch = createEventDispatcher();

	function addParam() {
		params = [...params, { key: '', value: '', description: '', enabled: true }];
		dispatch('paramsChange', { params });
	}

	function removeParam(index: number) {
		params = params.filter((_, i) => i !== index);
		dispatch('paramsChange', { params });
	}

	function updateParam(index: number, field: string, value: string | boolean) {
		params[index] = { ...params[index], [field]: value };
		dispatch('paramsChange', { params });
	}

	function toggleAllParams(enabled: boolean) {
		params = params.map(param => ({ ...param, enabled }));
		dispatch('paramsChange', { params });
	}

	function openBulkEdit() {
		// TODO: Implement bulk edit functionality
		console.log('Open bulk edit for params');
	}

	$: allEnabled = params.every(param => param.enabled);
	$: someEnabled = params.some(param => param.enabled);
</script>

<!-- Parameters section -->
<div role="tabpanel" aria-labelledby="params-tab">
	<div class="flex justify-between items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-list-ul text-orange-500 mr-2"></i>
			Query Parameters
		</h2>
		<div class="flex items-center space-x-2">
			<button 
				class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
				title="Open bulk edit mode for parameters"
				aria-label="Bulk edit parameters"
				on:click={openBulkEdit}
			>
				<i class="fas fa-edit text-xs mr-1"></i>
				Bulk Edit
			</button>
		</div>
	</div>
	
	<div class="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg border theme-border">
		<table class="w-full text-sm">
			<thead class="bg-gray-50 dark:bg-gray-750">
				<tr class="text-left theme-text-muted">
					<th class="p-3 font-medium w-12">
						<input
							class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
							type="checkbox"
							title="Select all parameters"
							aria-label="Select all parameters"
							checked={allEnabled}
							indeterminate={someEnabled && !allEnabled}
							on:change={(e) => toggleAllParams(e.currentTarget.checked)}
						/>
					</th>
					<th class="p-3 font-medium w-1/3">Key</th>
					<th class="p-3 font-medium w-1/3">Value</th>
					<th class="p-3 font-medium w-1/3">Description</th>
					<th class="p-3 font-medium w-12">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each params as param, index (index)}
					<tr class="border-t theme-border hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors duration-150">
						<td class="p-3">
							<input
								class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
								type="checkbox"
								title="Include this parameter"
								aria-label="Enable parameter"
								checked={param.enabled}
								on:change={(e) => updateParam(index, 'enabled', e.currentTarget.checked)}
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="param_name"
								type="text"
								title="Parameter key name"
								aria-label="Parameter key"
								value={param.key}
								on:input={(e) => updateParam(index, 'key', e.currentTarget.value)}
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="param_value"
								type="text"
								title="Parameter value"
								aria-label="Parameter value"
								value={param.value}
								on:input={(e) => updateParam(index, 'value', e.currentTarget.value)}
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="Optional description"
								type="text"
								title="Parameter description"
								aria-label="Parameter description"
								value={param.description}
								on:input={(e) => updateParam(index, 'description', e.currentTarget.value)}
							/>
						</td>
						<td class="p-3">
							<button
								class="text-red-400 hover:text-red-300 p-1 rounded hover:bg-red-400/10 transition-all duration-200"
								title="Remove parameter"
								aria-label="Remove this parameter"
								on:click={() => removeParam(index)}
							>
								<i class="fas fa-trash-alt text-xs"></i>
							</button>
						</td>
					</tr>
				{/each}
				
				<!-- Add new parameter row -->
				<tr class="border-t theme-border bg-gray-25 dark:bg-gray-825">
					<td colspan="5" class="p-3">
						<button
							class="w-full text-left text-sm theme-text-muted hover:theme-text-primary flex items-center transition-colors duration-200"
							title="Add new parameter"
							aria-label="Add new parameter"
							on:click={addParam}
						>
							<i class="fas fa-plus text-xs mr-2"></i>
							Add parameter
						</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>
	
	<div class="mt-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-md border-l-4 border-blue-400">
		<p class="text-sm theme-text-secondary flex items-start">
			<i class="fas fa-info-circle text-blue-400 mr-2 mt-0.5 flex-shrink-0"></i>
			<span>
				Query parameters will be automatically encoded and appended to the request URL. 
				Enable/disable individual parameters using the checkboxes.
			</span>
		</p>
	</div>
</div>
