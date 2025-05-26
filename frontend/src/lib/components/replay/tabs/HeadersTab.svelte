<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let headers: Array<{ key: string; value: string; description: string; enabled: boolean }> = [
		{ key: '', value: '', description: '', enabled: true }
	];

	const dispatch = createEventDispatcher();

	function addHeader() {
		headers = [...headers, { key: '', value: '', description: '', enabled: true }];
		dispatch('headersChange', { headers });
	}

	function removeHeader(index: number) {
		headers = headers.filter((_, i) => i !== index);
		dispatch('headersChange', { headers });
	}

	function updateHeader(index: number, field: string, value: string | boolean) {
		headers[index] = { ...headers[index], [field]: value };
		dispatch('headersChange', { headers });
	}

	function toggleAllHeaders(enabled: boolean) {
		headers = headers.map(header => ({ ...header, enabled }));
		dispatch('headersChange', { headers });
	}

	function openBulkEdit() {
		// TODO: Implement bulk edit functionality
		console.log('Open bulk edit for headers');
	}

	function addPresetHeader(key: string, value: string) {
		headers = [...headers, { key, value, description: '', enabled: true }];
		dispatch('headersChange', { headers });
	}

	// Common header presets
	const headerPresets = [
		{ key: 'Content-Type', value: 'application/json', description: 'JSON content type' },
		{ key: 'Accept', value: 'application/json', description: 'Accept JSON responses' },
		{ key: 'User-Agent', value: 'Beo-Echo/1.0', description: 'Custom user agent' },
		{ key: 'X-Requested-With', value: 'XMLHttpRequest', description: 'AJAX request identifier' },
		{ key: 'Cache-Control', value: 'no-cache', description: 'Disable caching' },
		{ key: 'Accept-Encoding', value: 'gzip, deflate', description: 'Supported encodings' }
	];

	$: allEnabled = headers.every(header => header.enabled);
	$: someEnabled = headers.some(header => header.enabled);
	$: enabledHeadersCount = headers.filter(header => header.enabled && header.key).length;
</script>

<!-- Headers section -->
<div role="tabpanel" aria-labelledby="headers-tab">
	<div class="flex justify-between items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-tags text-blue-500 mr-2"></i>
			Headers
			{#if enabledHeadersCount > 0}
				<span class="ml-2 px-2 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 text-xs rounded-full">
					{enabledHeadersCount}
				</span>
			{/if}
		</h2>
		<div class="flex items-center space-x-2">
			<button 
				class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
				title="Open bulk edit mode for headers"
				aria-label="Bulk edit headers"
				on:click={openBulkEdit}
			>
				<i class="fas fa-edit text-xs mr-1"></i>
				Bulk Edit
			</button>
		</div>
	</div>

	<!-- Header Presets -->
	<div class="mb-4 p-3 bg-gray-50 dark:bg-gray-900 rounded-lg border theme-border">
		<h3 class="text-xs font-medium theme-text-muted mb-2 uppercase tracking-wide">
			<i class="fas fa-magic mr-1"></i>
			Quick Add
		</h3>
		<div class="flex flex-wrap gap-2">
			{#each headerPresets as preset}
				<button
					class="text-xs bg-blue-100 dark:bg-blue-900/30 text-blue-600 dark:text-blue-400 px-2 py-1 rounded hover:bg-blue-200 dark:hover:bg-blue-900/50 transition-colors duration-200"
					title="Add {preset.key} header with value {preset.value}"
					aria-label="Add {preset.key} header"
					on:click={() => addPresetHeader(preset.key, preset.value)}
				>
					{preset.key}
				</button>
			{/each}
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
							title="Select all headers"
							aria-label="Select all headers"
							checked={allEnabled}
							indeterminate={someEnabled && !allEnabled}
							on:change={(e) => toggleAllHeaders(e.currentTarget.checked)}
						/>
					</th>
					<th class="p-3 font-medium w-1/3">Key</th>
					<th class="p-3 font-medium w-1/3">Value</th>
					<th class="p-3 font-medium w-1/3">Description</th>
					<th class="p-3 font-medium w-12">Actions</th>
				</tr>
			</thead>
			<tbody>
				{#each headers as header, index (index)}
					<tr class="border-t theme-border hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors duration-150">
						<td class="p-3">
							<input
								class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
								type="checkbox"
								title="Include this header"
								aria-label="Enable header"
								checked={header.enabled}
								on:change={(e) => updateHeader(index, 'enabled', e.currentTarget.checked)}
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="Content-Type"
								type="text"
								title="Header key name"
								aria-label="Header key"
								value={header.key}
								on:input={(e) => updateHeader(index, 'key', e.currentTarget.value)}
								list="header-suggestions"
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="application/json"
								type="text"
								title="Header value"
								aria-label="Header value"
								value={header.value}
								on:input={(e) => updateHeader(index, 'value', e.currentTarget.value)}
							/>
						</td>
						<td class="p-2">
							<input
								class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
								placeholder="Optional description"
								type="text"
								title="Header description"
								aria-label="Header description"
								value={header.description}
								on:input={(e) => updateHeader(index, 'description', e.currentTarget.value)}
							/>
						</td>
						<td class="p-3">
							<button
								class="text-red-400 hover:text-red-300 p-1 rounded hover:bg-red-400/10 transition-all duration-200"
								title="Remove header"
								aria-label="Remove this header"
								on:click={() => removeHeader(index)}
							>
								<i class="fas fa-trash-alt text-xs"></i>
							</button>
						</td>
					</tr>
				{/each}
				
				<!-- Add new header row -->
				<tr class="border-t theme-border bg-gray-25 dark:bg-gray-825">
					<td colspan="5" class="p-3">
						<button
							class="w-full text-left text-sm theme-text-muted hover:theme-text-primary flex items-center transition-colors duration-200"
							title="Add new header"
							aria-label="Add new header"
							on:click={addHeader}
						>
							<i class="fas fa-plus text-xs mr-2"></i>
							Add header
						</button>
					</td>
				</tr>
			</tbody>
		</table>
	</div>

	<!-- Header suggestions datalist -->
	<datalist id="header-suggestions">
		<option value="Accept">
		<option value="Accept-Encoding">
		<option value="Accept-Language">
		<option value="Authorization">
		<option value="Cache-Control">
		<option value="Connection">
		<option value="Content-Length">
		<option value="Content-Type">
		<option value="Cookie">
		<option value="Host">
		<option value="Origin">
		<option value="Referer">
		<option value="User-Agent">
		<option value="X-Requested-With">
		<option value="X-CSRF-Token">
		<option value="X-API-Key">
		<option value="X-Custom-Header">
	</datalist>

	<!-- Info and Tips -->
	<div class="mt-4 space-y-3">
		<div class="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-md border-l-4 border-blue-400">
			<p class="text-sm theme-text-secondary flex items-start">
				<i class="fas fa-info-circle text-blue-400 mr-2 mt-0.5 flex-shrink-0"></i>
				<span>
					Headers are case-insensitive. Use the quick add buttons for common headers or type custom header names.
					Disabled headers will not be sent with the request.
				</span>
			</p>
		</div>
		
		{#if enabledHeadersCount > 0}
			<div class="p-3 bg-green-50 dark:bg-green-900/20 rounded-md border-l-4 border-green-400">
				<p class="text-sm theme-text-secondary flex items-center">
					<i class="fas fa-check-circle text-green-400 mr-2"></i>
					{enabledHeadersCount} header{enabledHeadersCount === 1 ? '' : 's'} will be sent with the request
				</p>
			</div>
		{/if}
	</div>
</div>
