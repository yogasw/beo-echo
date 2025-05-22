<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { theme } from '$lib/stores/theme';
	import { createEventDispatcher } from 'svelte';
	
	export let headers: string;
	export let editable: boolean = true;
	export let title: string = 'Headers';
	
	const dispatch = createEventDispatcher<{
		change: string;
	}>();
	
	// Local state
	let localHeaders: Array<{key: string, value: string, isNew?: boolean}> = [];
	let isEditing = false;
	let hasChanges = false;
	
	// Parse JSON headers string into an array of key-value objects
	$: {
		try {
			if (!isEditing) {
				const headersObj = headers ? JSON.parse(headers) : {};
				localHeaders = Object.entries(headersObj).map(([key, value]) => ({ key, value }));
				hasChanges = false;
			}
		} catch (error) {
			console.error('Error parsing headers:', error);
			if (!isEditing) {
				localHeaders = [];
				hasChanges = false;
			}
		}
	}
	
	// Add a new header entry
	function addHeader() {
		localHeaders = [...localHeaders, { key: '', value: '', isNew: true }];
		hasChanges = true;
		isEditing = true;
	}
	
	// Remove a header entry
	function removeHeader(index: number) {
		localHeaders = localHeaders.filter((_, i) => i !== index);
		hasChanges = true;
	}
	
	// Handle header key change
	function handleKeyChange(index: number, newKey: string) {
		localHeaders[index].key = newKey;
		hasChanges = true;
		localHeaders = [...localHeaders]; // Trigger reactivity
	}
	
	// Handle header value change
	function handleValueChange(index: number, newValue: string) {
		localHeaders[index].value = newValue;
		hasChanges = true;
		localHeaders = [...localHeaders]; // Trigger reactivity
	}
	
	// Save changes to headers
	function saveChanges() {
		// Filter out empty headers
		const filteredHeaders = localHeaders.filter(h => h.key.trim() !== '');
		// Convert array back to object
		const headersObj = filteredHeaders.reduce((obj, { key, value }) => {
			obj[key.trim()] = value;
			return obj;
		}, {} as Record<string, string>);
		
		// Convert to JSON string and dispatch event
		const newHeadersJson = JSON.stringify(headersObj);
		headers = newHeadersJson;
		dispatch('change', newHeadersJson);
		isEditing = false;
		hasChanges = false;
	}
	
	// Cancel editing
	function cancelEdit() {
		// Reset to the original headers
		try {
			const headersObj = headers ? JSON.parse(headers) : {};
			localHeaders = Object.entries(headersObj).map(([key, value]) => ({ key, value }));
		} catch (error) {
			localHeaders = [];
		}
		isEditing = false;
		hasChanges = false;
	}
	
	// Enable edit mode
	function startEditing() {
		isEditing = true;
	}
</script>

<div class="w-full {!isEditing ? 'rounded-md overflow-hidden shadow-sm' : ThemeUtils.card('w-full')}">
	<div class="p-0">
		<!-- View Mode -->
		{#if !isEditing}
			{#if localHeaders && localHeaders.length > 0}
				<div class="max-h-[350px] overflow-y-auto">
					<table class="w-full border-collapse">
						<thead class="sticky top-0 z-10">
							<tr class="bg-gray-100 dark:bg-gray-700 border-b dark:border-gray-600">
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}">Header Name</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}">Value</th>
							</tr>
						</thead>
						<tbody class="divide-y dark:divide-gray-700">
							{#each localHeaders as header, index}
								<tr class="{index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-3 align-top">
										<span class="font-medium whitespace-nowrap text-blue-600 dark:text-blue-400">{header.key}</span>
									</td>
									<td class="px-4 py-3 break-all {ThemeUtils.themeTextSecondary()}">{header.value}</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
				
				{#if editable}
					<div class="flex justify-end mt-2">
						<button 
							on:click={startEditing} 
							class="{ThemeUtils.utilityButton('py-1 px-2 text-xs')}"
							type="button"
							aria-label="Edit headers"
						>
							<i class="fas fa-edit text-xs mr-1"></i>
							<span>Edit</span>
						</button>
					</div>
				{/if}
			{:else}
				<div class="text-center py-6 rounded-lg border border-dashed {ThemeUtils.themeBorder()}">
					<i class="fas fa-info-circle mb-2 text-lg {ThemeUtils.themeTextMuted()}"></i>
					<div class="{ThemeUtils.themeTextMuted()} italic">No headers available.</div>
					
					{#if editable}
						<button 
							on:click={startEditing} 
							class="{ThemeUtils.secondaryButton('py-1.5 px-3 text-xs mt-3')}"
							type="button"
						>
							<i class="fas fa-plus text-xs mr-1"></i>
							<span>Add Headers</span>
						</button>
					{/if}
				</div>
			{/if}
		
		<!-- Edit Mode -->
		{:else}
			<div class="space-y-2">
				{#if localHeaders.length === 0}
					<div class="flex justify-center py-4">
						<button
							on:click={addHeader}
							class="{ThemeUtils.secondaryButton('py-2 px-4 text-sm gap-2')}"
							type="button"
						>
							<i class="fas fa-plus"></i>
							<span>Add First Header</span>
						</button>
					</div>
				{:else}
					<div class="max-h-[350px] overflow-y-auto px-1 py-1 border dark:border-gray-700 rounded-md">
						{#each localHeaders as header, i}
							<div class="flex items-center gap-2 px-3 py-2 {i !== 0 ? 'border-t dark:border-gray-700/50' : ''} {i % 2 === 0 ? '' : 'bg-gray-50/50 dark:bg-gray-750/50'}">
								<!-- Key Input -->
								<div class="flex-1">
									<input 
										type="text"
										bind:value={header.key}
										on:input={() => handleKeyChange(i, header.key)}
										placeholder="Header name"
										class="block w-full p-1.5 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
										aria-label="Header name"
									/>
								</div>
								
								<!-- Value Input -->
								<div class="flex-1">
									<input 
										type="text" 
										bind:value={header.value}
										on:input={() => handleValueChange(i, header.value)}
										placeholder="Value"
										class="block w-full p-1.5 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
										aria-label="Header value"
									/>
								</div>
								
								<!-- Delete Button -->
								<button
									on:click={() => removeHeader(i)}
									title="Remove header"
									class="h-6 w-6 flex items-center justify-center text-red-500 hover:bg-red-100 dark:hover:bg-red-900/30 dark:text-red-400 dark:hover:text-red-300 rounded-full"
									type="button"
									aria-label="Remove header"
								>
									<i class="fas fa-trash-alt text-xs"></i>
								</button>
							</div>
						{/each}
					</div>
					
					<!-- Add New Header Button -->
					<div>
						<button
							on:click={addHeader}
							class="{ThemeUtils.utilityButton('py-1 px-2 text-xs gap-1')}"
							type="button"
							aria-label="Add new header"
						>
							<i class="fas fa-plus text-xs"></i>
							<span>Add Header</span>
						</button>
					</div>
				{/if}
				
				<!-- Action Buttons -->
				<div class="flex justify-end items-center gap-2 pt-4 mt-2 border-t {ThemeUtils.themeBorderLight()}">
					<button
						on:click={cancelEdit}
						class="{ThemeUtils.utilityButton('py-1.5 px-3 text-xs')}"
						type="button"
						aria-label="Cancel editing"
					>
						Cancel
					</button>
					
					<button
						on:click={saveChanges}
						class="{hasChanges ? ThemeUtils.primaryButton('py-1.5 px-3 text-xs') : ThemeUtils.secondaryButton('py-1.5 px-3 text-xs opacity-60')}"
						type="button"
						disabled={!hasChanges}
						aria-label="Save changes"
					>
						<i class="fas fa-save mr-1"></i>
						Save
					</button>
				</div>
			</div>
		{/if}
	</div>
</div>
