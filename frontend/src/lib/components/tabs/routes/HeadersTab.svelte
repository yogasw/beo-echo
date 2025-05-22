<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { theme } from '$lib/stores/theme';
	import { createEventDispatcher } from 'svelte';
	import { tick } from 'svelte';
	
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
	// References for input elements to enable focusing
	let valueInputs: HTMLInputElement[] = [];
	let keyInputs: HTMLInputElement[] = [];	
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
	async function addHeader() {
		const newIndex = localHeaders.length;
		localHeaders = [...localHeaders, { key: '', value: '', isNew: true }];
		hasChanges = true;
		isEditing = true;
		
		// Wait for DOM update, then focus on the new input
		await tick();
		if (keyInputs[newIndex]) {
			keyInputs[newIndex].focus();
		}
	}
	
	// Remove a header entry
	function removeHeader(index: number) {
		localHeaders = localHeaders.filter((_, i) => i !== index);
		hasChanges = true;
	}
	
	// Duplicate a header entry
	function duplicateHeader(index: number) {
		const header = localHeaders[index];
		localHeaders = [...localHeaders.slice(0, index + 1), { key: '', value: '', isNew: true }, ...localHeaders.slice(index + 1)];
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
	
	// Handle key press events
	async function handleKeyDown(event: KeyboardEvent, index: number) {
		// If Enter key is pressed in the last row's value input, add a new header
		if (event.key === 'Enter' && index === localHeaders.length - 1) {
			event.preventDefault(); // Prevent form submission
			await addHeader();
		}
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

<div class="w-full rounded-md overflow-hidden shadow-sm">
	<div class="p-0">
		<div class="bg-gray-100 dark:bg-gray-700 border-b dark:border-gray-600 flex justify-between items-center">
			<div class="flex">
				<div class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}">Header Name</div>
				<div class="text-left pr-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}">Value</div>
			</div>
			
			{#if editable}
				<div class="flex pr-4">
					<button 
						on:click={isEditing ? saveChanges : startEditing} 
						class="{isEditing && !hasChanges ? ThemeUtils.secondaryButton('py-1 px-2 text-xs opacity-60') : isEditing ? ThemeUtils.primaryButton('py-1 px-2 text-xs') : ThemeUtils.utilityButton('py-1 px-2 text-xs')}"
						type="button"
						disabled={isEditing && !hasChanges}
						aria-label={isEditing ? "Save headers" : "Edit headers"}
					>
						{#if isEditing}
							<i class="fas fa-save text-xs mr-1"></i>
							<span>Save</span>
						{:else}
							<i class="fas fa-edit text-xs mr-1"></i>
							<span>Edit</span>
						{/if}
					</button>
					
					{#if isEditing}
						<button 
							on:click={cancelEdit} 
							class="{ThemeUtils.utilityButton('py-1 px-2 text-xs ml-1')}"
							type="button"
							aria-label="Cancel editing"
						>
							Cancel
						</button>
					{/if}
				</div>
			{/if}
		</div>
		
		<!-- View Mode -->
		{#if !isEditing}
			{#if localHeaders && localHeaders.length > 0}
				<div class="max-h-[350px] overflow-y-auto">
					<table class="w-full border-collapse">
						<tbody class="divide-y dark:divide-gray-700">
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
				
			{:else}
				<div class="text-center py-6 rounded-lg border border-dashed {ThemeUtils.themeBorder()}">
					<i class="fas fa-info-circle mb-2 text-lg {ThemeUtils.themeTextMuted()}"></i>
					<div class="{ThemeUtils.themeTextMuted()} italic">No headers available.</div>
				</div>
			{/if}
		
		<!-- Edit Mode -->
		{:else}
			{#if localHeaders.length === 0}
				<div class="text-center py-6">
					<button
						on:click={addHeader}
						class="{ThemeUtils.secondaryButton('py-1.5 px-3 text-xs')}"
						type="button"
					>
						<i class="fas fa-plus text-xs mr-1"></i>
						<span>Add First Header</span>
					</button>
				</div>
			{:else}
				<div class="max-h-[350px] overflow-y-auto">
					<table class="w-full border-collapse">
						<tbody class="divide-y dark:divide-gray-700">
							{#each localHeaders as header, i}
								<tr class="{i % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-2 align-top" style="width: 50%">
										<input 
											type="text"
											bind:value={header.key}
											on:input={() => handleKeyChange(i, header.key)}
											placeholder="Header name"
											class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
											aria-label="Header name"
											bind:this={keyInputs[i]}
										/>
									</td>
									<td class="px-4 py-2 align-top relative">
										<div class="flex items-center gap-2 w-full">
											<div class="flex-1">
												<input 
													type="text" 
													bind:value={header.value}
													on:input={() => handleValueChange(i, header.value)}
													placeholder="Value"
													class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
													aria-label="Header value"
													bind:this={valueInputs[i]}
													on:keydown={(e) => handleKeyDown(e, i)}
												/>
											</div>
											{#if i === localHeaders.length - 1}
												<button
													on:click={() => addHeader()}
													title="Add new header"
													class="h-6 w-6 flex-shrink-0 flex items-center justify-center text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-900/30 dark:text-blue-400 dark:hover:text-blue-300 rounded-full"
													type="button"
													aria-label="Add new header"
												>
													<i class="fas fa-plus text-xs"></i>
												</button>
											{/if}
											<button
												on:click={() => removeHeader(i)}
												title="Remove header"
												class="h-6 w-6 flex-shrink-0 flex items-center justify-center text-red-500 hover:bg-red-100 dark:hover:bg-red-900/30 dark:text-red-400 dark:hover:text-red-300 rounded-full"
												type="button"
												aria-label="Remove header"
											>
												<i class="fas fa-trash-alt text-xs"></i>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		{/if}
	</div>
</div>
