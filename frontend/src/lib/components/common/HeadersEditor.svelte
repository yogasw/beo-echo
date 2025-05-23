<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { createEventDispatcher, onMount, onDestroy } from 'svelte';
	import { tick } from 'svelte';
	
	export let headers: string;
	export let editable: boolean = true;
	export const title: string = 'Headers'; // Changed to const since it's not used internally
	export let maxContentHeight: string = ''; // Optional explicit height
	export let onSave: ((headers: string) => void) | undefined = undefined; // Callback when headers are saved
	
	const dispatch = createEventDispatcher<{
		change: string;
	}>();
	
	// Local state
	let localHeaders: Array<{key: string, value: string, isNew?: boolean}> = [];
	let isEditing = false;
	let hasChanges = false;
	// References for input elements to enable focusing
	let valueInputs: HTMLTextAreaElement[] = [];
	let keyInputs: HTMLInputElement[] = [];
	// Container references for height calculation
	let containerElement: HTMLElement;
	let tableContainer: HTMLElement;
	let calculatedHeight = '300px'; // Default height
	
	// Parse JSON headers string into an array of key-value objects
	$: {
		try {
			if (!isEditing) {
				const headersObj = headers ? JSON.parse(headers) : {};
				localHeaders = Object.entries(headersObj).map(([key, value]) => ({ 
					key, 
					value: typeof value === 'string' ? value : String(value) 
				}));
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
	
	// Calculate the height based on the parent container
	function calculateHeight() {
		if (!containerElement || !tableContainer) return;
		
		// If explicit height is provided, use it
		if (maxContentHeight) {
			calculatedHeight = maxContentHeight;
			return;
		}
		
		// Get the position of the container relative to the viewport
		const containerRect = containerElement.getBoundingClientRect();
		// Determine how much space is available (viewport height minus container top position minus some padding)
		const availableHeight = window.innerHeight - containerRect.top - 75;
		// Ensure the height is at least 100px but not more than the available space
		calculatedHeight = `${Math.max(100, Math.min(350, availableHeight))}px`;
	}
	
	// Set up resize observer to recalculate height when container size changes
	let resizeObserver: ResizeObserver;
	
	onMount(() => {
		// Initial height calculation
		calculateHeight();
		
		// Set up resize observer
		resizeObserver = new ResizeObserver(() => {
			calculateHeight();
		});
		
		if (containerElement) {
			resizeObserver.observe(containerElement);
		}
		
		// Also listen to window resize events
		window.addEventListener('resize', calculateHeight);
		
		// Initialize textarea heights for existing values
		if (isEditing) {
			setTimeout(() => {
				valueInputs.forEach(textarea => {
					if (textarea) {
						autoResizeTextarea(textarea);
					}
				});
			}, 0);
		}
	});
	
	onDestroy(() => {
		// Clean up resize observer
		if (resizeObserver) {
			resizeObserver.disconnect();
		}
		window.removeEventListener('resize', calculateHeight);
	});
	
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
	
	// Add a function to auto-resize textareas based on content
	function autoResizeTextarea(textarea: HTMLTextAreaElement) {
		if (!textarea) return;
		// Reset height to auto to get the correct scrollHeight
		textarea.style.height = 'auto';
		// Set to the scrollHeight to show all content
		textarea.style.height = textarea.scrollHeight + 'px';
	}
	
	// Handle value changes and resize the textarea
	function handleValueChange(index: number, newValue: string) {
		localHeaders[index].value = newValue;
		hasChanges = true;
		localHeaders = [...localHeaders]; // Trigger reactivity
		
		// Resize the textarea after the value changes
		setTimeout(() => {
			if (valueInputs[index]) {
				autoResizeTextarea(valueInputs[index]);
			}
		}, 0);
	}
	
	// Handle key press events
	async function handleKeyDown(event: KeyboardEvent, index: number) {
		// Enter key handling is now directly in the textarea element
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
		
		// Call onSave if provided
		if (onSave) {
			onSave(newHeadersJson);
		}
		
		isEditing = false;
		hasChanges = false;
	}
	
	// Cancel editing
	function cancelEdit() {
		// Reset to the original headers
		try {
			const headersObj = headers ? JSON.parse(headers) : {};
			localHeaders = Object.entries(headersObj).map(([key, value]) => ({ 
				key, 
				value: typeof value === 'string' ? value : String(value) 
			}));
		} catch (error) {
			localHeaders = [];
		}
		isEditing = false;
		hasChanges = false;
	}
	
	// Enable edit mode
	function startEditing() {
		isEditing = true;
		
		// Initialize textarea heights after switching to edit mode
		setTimeout(() => {
			valueInputs.forEach(textarea => {
				if (textarea) {
					autoResizeTextarea(textarea);
				}
			});
		}, 0);
	}
</script>

<div class="w-full rounded-md overflow-hidden shadow-sm" bind:this={containerElement}>
	<div class="p-0">
		<div class="bg-gray-100 dark:bg-gray-700 border-b dark:border-gray-600 flex justify-between items-center">
			<div class="flex w-full">									<div class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-1/2">Header Name</div>
								<div class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}" style="flex: 1">Value</div>
			</div>
			
			{#if editable}
				<div class="flex pr-4 shrink-0">
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
				<div bind:this={tableContainer} class="overflow-y-auto" style="max-height: {calculatedHeight}">
					<table class="w-full border-collapse table-fixed">
						<tbody class="divide-y dark:divide-gray-700">
							{#each localHeaders as header, index}
								<tr class="{index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-2 align-top w-1/2" style="min-width: 150px">
										<span class="font-medium whitespace-nowrap text-blue-600 dark:text-blue-400 text-xs">{header.key}</span>
									</td>
									<td class="px-4 py-2 align-top" style="min-width: 200px">
										<div class="break-all whitespace-pre-wrap text-xs {ThemeUtils.themeTextSecondary()}" style="max-width: 100%; overflow-x: auto">{header.value}</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
				
			{:else}
				<div class="text-center py-6 rounded-lg border border-dashed {ThemeUtils.themeBorder()}">
					<i class="fas fa-info-circle mb-2 text-lg {ThemeUtils.themeTextMuted()}"></i>
					<div class="{ThemeUtils.themeTextMuted()} italic text-xs">No headers available.</div>
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
				<div bind:this={tableContainer} class="overflow-y-auto" style="max-height: {calculatedHeight}">
					<table class="w-full border-collapse table-fixed">
						<tbody class="divide-y dark:divide-gray-700">
							{#each localHeaders as header, i}
								<tr class="{i % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-2 align-top w-1/2" style="min-width: 150px">
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
									<td class="px-4 py-2 align-top relative" style="min-width: 200px">
										<div class="flex items-center gap-2 w-full">
											<div class="flex-1">
												<textarea 
													bind:value={header.value}
													on:input={(e) => {
														handleValueChange(i, header.value);
														// Auto-adjust height based on content
														const textarea = e.target as HTMLTextAreaElement;
														if (textarea) {
															textarea.style.height = 'auto';
															textarea.style.height = textarea.scrollHeight + 'px';
														}
													}}
													placeholder="Value"
													class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500 resize-none"
													style="min-height: 32px; height: auto; overflow: hidden;"
													aria-label="Header value"
													bind:this={valueInputs[i]}
													on:keydown={(e) => {
														if (e.key === 'Enter') {
															e.preventDefault();
															if (i === localHeaders.length - 1) {
																addHeader();
															}
														}
														handleKeyDown(e, i);
													}}
												></textarea>
											</div>
											<div class="flex shrink-0">
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
