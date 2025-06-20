<script lang="ts">
	import type { Endpoint, Response } from '$lib/api/BeoApi';
	import { addResponse, deleteResponse, duplicateResponse, updateResponse, reorderResponses } from '$lib/api/BeoApi';
	import { updateEndpoint } from '$lib/stores/saveButton';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	let selectedValue: string = '';
	export let selectedEndpoint: Endpoint | null;
	export let selectedResponse: Response | null;

	let flaggedResponseId: string | null = null;
	
	// Drag and drop state
	let draggedIndex: number | null = null;
	let dragOverIndex: number | null = null;
	let isDragging = false;

	const toggleDropdown = (): void => {
		const dropdown = document.getElementById('dropdownMenu');
		if (dropdown) {
			dropdown.classList.toggle('hidden');
		}
	};

	// Helper function to truncate text with ellipsis
	function truncateText(text: string, maxLength: number = 25): string {
		if (!text) return '';
		return text.length > maxLength ? text.substring(0, maxLength) + '...' : text;
	}

	// Helper function to sanitize text by removing HTML and collapsing whitespace
	function sanitizeText(text: string): string {
		if (!text) return '';
		// Remove HTML tags
		const withoutHtml = text.replace(/<[^>]*>/g, ' ');
		// Collapse whitespace (newlines, tabs, multiple spaces)
		return withoutHtml.replace(/\s+/g, ' ').trim();
	}

	// Format response label for display
	function formatResponseLabel(
		index: number,
		response: Response,
		isForSelectedValue: boolean = false
	): string {
		const statusText = `(${response.status_code})`;
		let noteText = response.note ? ` ${response.note}` : '';
		const enabledText = response.enabled ? '' : ' [DISABLED]';

		// Sanitize the note text to remove HTML and collapse whitespace
		noteText = sanitizeText(noteText);

		// Use different truncation length depending on whether this is for the selected value or dropdown items
		const maxLength = isForSelectedValue ? 40 : 25;

		// Format display with truncation
		const baseText = `Response ${index + 1} ${statusText} ${truncateText(noteText, maxLength)}`;
		return baseText + enabledText;
	}

	const selectResponse = (index: number, value: Response): void => {
		toggleDropdown();
		selectedResponse = value;
	};

	const handleAddResponse = async (): Promise<void> => {
		if (selectedEndpoint) {
			try {
				// Sample response values as specified
				const statusCode = 200;
				const body = '{"message":"Hello World"}';
				const headers = '{"Content-Type":"application/json"}';

				const newResponse = await addResponse(
					selectedEndpoint.project_id,
					selectedEndpoint.id,
					statusCode,
					body,
					headers,
					''
				);

				// Update the endpoint responses array
				if (selectedEndpoint.responses) {
					selectedEndpoint.responses = [...selectedEndpoint.responses, newResponse];
				} else {
					selectedEndpoint.responses = [newResponse];
				}

				// Auto-select the newly created response
				selectResponse(selectedEndpoint.responses.length - 1, newResponse);
			} catch (error) {
				console.error('Failed to add response:', error);
				// You may want to show a toast or notification here
			}
		}
	};

	// Toggle flag for a response
	function toggleFlag(responseId: string) {
		flaggedResponseId = responseId === flaggedResponseId ? null : responseId;
	}

	// Handle response mode change
	const handleResponseModeChange = async (
		mode: 'static' | 'random' | 'round_robin'
	): Promise<void> => {
		if (!selectedEndpoint) {
			toast.error('No endpoint selected');
			return;
		}
		selectedEndpoint = updateEndpoint('response_mode', mode, selectedEndpoint);
	};

	$: {
		if (selectedResponse && selectedEndpoint?.responses) {
			// Sort responses by priority before finding index
			const sortedResponses = selectedEndpoint.responses.sort((a, b) => a.priority - b.priority);
			const index = sortedResponses.findIndex((r) => r.id === selectedResponse?.id);
			if (index !== -1) {
				selectedValue = formatResponseLabel(index, selectedResponse, true);
			} else {
				selectedValue = 'Response not found';
			}
		} else {
			selectedValue = 'No Response';
		}
	}

	async function handleDeleteResponse(response: Response): Promise<void> {
		if (!selectedEndpoint || !response) {
			toast.error('No endpoint or response selected');
			return;
		}

		let deleted = await deleteResponse(
			selectedEndpoint?.project_id,
			response.endpoint_id,
			response.id
		);

		if (!deleted) {
			toast.error('Failed to delete response');
			return;
		}

		// select other response if available
		let otherResponses = selectedEndpoint?.responses?.filter((r) => r.id !== response.id);
		if (otherResponses) {
			selectedResponse = otherResponses[0] || null;
			if (selectedEndpoint?.responses) {
				selectedEndpoint.responses = otherResponses || [];
			}
		}
	}

	async function handleDuplicateResponse(response: Response): Promise<void> {
		if (!selectedEndpoint || !response) {
			toast.error('No endpoint or response selected');
			return;
		}

		try {
			const duplicatedResponse = await duplicateResponse(
				selectedEndpoint.project_id,
				response.endpoint_id,
				response.id
			);

			// Add the duplicated response to the endpoint's responses array
			if (selectedEndpoint.responses) {
				selectedEndpoint.responses = [...selectedEndpoint.responses, duplicatedResponse];
			} else {
				selectedEndpoint.responses = [duplicatedResponse];
			}

			// Auto-select the newly duplicated response
			selectResponse(selectedEndpoint.responses.length - 1, duplicatedResponse);
			
			toast.success('Response duplicated successfully');
		} catch (error) {
			console.error('Failed to duplicate response:', error);
			toast.error('Failed to duplicate response');
		}
	}

	async function handleToggleResponseEnabled(response: Response): Promise<void> {
		if (!selectedEndpoint || !response) {
			toast.error('No endpoint or response selected');
			return;
		}

		try {
			const updatedResponse = await updateResponse(
				selectedEndpoint.project_id,
				response.endpoint_id,
				response.id,
				{ enabled: !response.enabled }
			);

			// Update the response in the endpoint's responses array
			if (selectedEndpoint.responses) {
				const index = selectedEndpoint.responses.findIndex((r) => r.id === response.id);
				if (index !== -1) {
					selectedEndpoint.responses[index] = updatedResponse;
					// If this is the currently selected response, update it too
					if (selectedResponse?.id === response.id) {
						selectedResponse = updatedResponse;
					}
				}
			}

			const statusMessage = updatedResponse.enabled ? 'enabled' : 'disabled';
		} catch (error) {
			console.error('Failed to toggle response enabled state:', error);
			toast.error(error);
		}
	}

	// Drag and drop handlers
	function handleDragStart(event: DragEvent, index: number): void {
		if (!event.dataTransfer) return;
		
		draggedIndex = index;
		isDragging = true;
		event.dataTransfer.effectAllowed = 'move';
		event.dataTransfer.setData('text/html', ''); // Required for Firefox
		
		// Add visual feedback
		if (event.target instanceof HTMLElement) {
			event.target.style.opacity = '0.5';
		}
	}

	function handleDragEnd(event: DragEvent): void {
		if (event.target instanceof HTMLElement) {
			event.target.style.opacity = '1';
		}
		draggedIndex = null;
		dragOverIndex = null;
		isDragging = false;
	}

	function handleDragOver(event: DragEvent, index: number): void {
		event.preventDefault();
		if (draggedIndex === null || draggedIndex === index) return;
		
		dragOverIndex = index;
		
		if (event.dataTransfer) {
			event.dataTransfer.dropEffect = 'move';
		}
	}

	function handleDragLeave(): void {
		dragOverIndex = null;
	}

	async function handleDrop(event: DragEvent, dropIndex: number): Promise<void> {
		event.preventDefault();
		
		if (draggedIndex === null || draggedIndex === dropIndex || !selectedEndpoint?.responses) {
			return;
		}

		try {
			// Sort responses by priority first to get the correct order
			const sortedResponses = [...selectedEndpoint.responses].sort((a, b) => a.priority - b.priority);
			const draggedResponse = sortedResponses[draggedIndex];
			
			// Remove from old position
			sortedResponses.splice(draggedIndex, 1);
			
			// Insert at new position
			sortedResponses.splice(dropIndex, 0, draggedResponse);
			
			// Create order array with response IDs
			const responseOrder = sortedResponses.map(r => r.id);
			
			// Call API to update order
			const updatedResponses = await reorderResponses(
				selectedEndpoint.project_id,
				selectedEndpoint.id,
				responseOrder
			);
			
			// Update local state
			selectedEndpoint.responses = updatedResponses;
			
			// If the selected response was moved, keep it selected
			if (selectedResponse && draggedResponse.id === selectedResponse.id) {
				// Find the new index in the updated responses
				const newSortedResponses = updatedResponses.sort((a, b) => a.priority - b.priority);
				const newIndex = newSortedResponses.findIndex(r => r.id === draggedResponse.id);
				if (newIndex !== -1) {
					selectedResponse = newSortedResponses[newIndex];
				}
			}
			
			toast.success('Response order updated successfully');
		} catch (error) {
			console.error('Failed to reorder responses:', error);
			toast.error('Failed to update response order');
		} finally {
			draggedIndex = null;
			dragOverIndex = null;
			isDragging = false;
		}
	}
</script>

<style>
	.drag-ghost {
		opacity: 0.5;
	}
	
	.drag-over {
		border-top: 2px solid #3b82f6;
		background-color: rgba(59, 130, 246, 0.1);
	}
	
	.cursor-grab {
		cursor: grab;
	}
	
	.cursor-grab:active,
	.cursor-grabbing {
		cursor: grabbing;
	}
</style>

<div
	class="flex items-center justify-between {ThemeUtils.themeBgPrimary()} {ThemeUtils.themeTextSecondary()} py-2"
>
	<div class="flex items-center w-full">
		<button
			class="bg-blue-500 text-white rounded px-2 py-1 text-sm font-medium"
			on:click={handleAddResponse}
			aria-label="Add Response"
		>
			<i class="fas fa-plus"></i>
		</button>
		<div class="relative ml-4 w-full">
			<button
				class="text-sm font-medium {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextSecondary()} rounded px-2 py-1 flex items-center justify-between w-full"
				on:click={() => {
					toggleDropdown();
				}}
				title="Select response"
				aria-label="Select response"
			>
				<span
					class="truncate mr-2 inline-block max-w-[calc(100%-20px)] whitespace-nowrap overflow-hidden"
					>{selectedValue}</span
				>
				<i class="fas fa-chevron-down flex-shrink-0"></i>
			</button>
			<div
				id="dropdownMenu"
				class="absolute mt-1 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextSecondary()} rounded shadow-lg w-full hidden z-50 max-h-60 overflow-y-auto"
			>
				<ul class="text-sm">
					{#if selectedEndpoint?.responses && selectedEndpoint.responses.length > 0}
						{#each selectedEndpoint.responses.sort((a, b) => a.priority - b.priority) as response, index}
							<li 
								class="flex items-center group transition-all duration-200 {dragOverIndex === index ? 'drag-over' : ''} {isDragging && draggedIndex === index ? 'drag-ghost' : ''}"
								draggable="true"
								on:dragstart={(event) => handleDragStart(event, index)}
								on:dragend={handleDragEnd}
								on:dragover={(event) => handleDragOver(event, index)}
								on:dragleave={handleDragLeave}
								on:drop={(event) => handleDrop(event, index)}
							>
								<!-- Drag handle -->
								<div class="flex items-center px-2 py-2 cursor-grab active:cursor-grabbing opacity-0 group-hover:opacity-100 transition-opacity duration-200" 
								     title="Drag to reorder"
								     aria-label="Drag to reorder response">
									<i class="fas fa-grip-vertical text-gray-400 text-xs"></i>
								</div>
								
								<button
									type="button"
									class="flex-1 text-left px-2 py-2 {ThemeUtils.themeHover()} cursor-pointer {!response.enabled ? 'opacity-50' : ''} transition-colors duration-200"
									on:click={() => {
										selectResponse(index, response);
									}}
									title="Select this response"
									aria-label="Select this response"
								>
									<span class="{!response.enabled ? 'line-through' : ''}">
										Response {index + 1} ({response.status_code}) {truncateText(
											sanitizeText(response?.note)
										)}
									</span>
									{#if !response.enabled}
										<span class="text-xs text-gray-500 ml-2">[DISABLED]</span>
									{/if}
								</button>
								<!-- Copy & Delete actions -->
								<div class="flex items-center space-x-1 mr-2 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
									<button
										class="p-1 rounded hover:bg-gray-700 transition-colors duration-200"
										title="Flag response"
										aria-label="Flag this response"
										on:click|stopPropagation={() => {
											toggleFlag(response.id);
										}}
									>
										<i
											class="fas fa-flag {flaggedResponseId === response.id
												? 'text-red-500'
												: 'text-gray-300'}"
										></i>
									</button>
									<button
										class="p-1 rounded hover:bg-gray-700 transition-colors duration-200"
										title="{response.enabled ? 'Disable response' : 'Enable response'}"
										aria-label="{response.enabled ? 'Disable this response' : 'Enable this response'}"
										on:click|stopPropagation={() => {
											handleToggleResponseEnabled(response);
										}}
									>
										<i
											class="fas {response.enabled ? 'fa-eye text-green-400' : 'fa-eye-slash text-gray-500'}"
										></i>
									</button>
									<button
										class="p-1 rounded hover:bg-gray-700 transition-colors duration-200"
										title="Duplicate response"
										aria-label="Duplicate this response"
										on:click|stopPropagation={() => {
											handleDuplicateResponse(response);
										}}
									>
										<i class="fas fa-copy text-gray-300"></i>
									</button>
									<button
										class="p-1 rounded hover:bg-gray-700 transition-colors duration-200"
										title="Delete response"
										aria-label="Delete this response"
										on:click|stopPropagation={() => {
											handleDeleteResponse(response);
										}}
									>
										<i class="fas fa-trash text-gray-300"></i>
									</button>
								</div>
							</li>
						{/each}
					{:else}
						<li class="px-4 py-2 {ThemeUtils.themeTextMuted()}">No Response</li>
					{/if}
				</ul>
			</div>
		</div>
	</div>
	<!-- Action buttons group -->
	<div class="flex items-center space-x-1 ml-4">
		<button
			class="p-2 rounded-md text-sm transition-all duration-200 {selectedEndpoint?.response_mode ===
			'static'
				? 'bg-blue-600 text-white shadow-md'
				: 'bg-gray-700 dark:bg-gray-600 text-gray-300 hover:bg-gray-600 dark:hover:bg-gray-500 hover:text-white'}"
			title="Static Mode - Return highest priority response"
			aria-label="Static response mode"
			on:click={() => handleResponseModeChange('static')}
			disabled={!selectedEndpoint}
		>
			<i class="fas fa-list text-xs"></i>
		</button>
		<button
			class="p-2 rounded-md text-sm transition-all duration-200 {selectedEndpoint?.response_mode ===
			'random'
				? 'bg-green-600 text-white shadow-md'
				: 'bg-gray-700 dark:bg-gray-600 text-gray-300 hover:bg-gray-600 dark:hover:bg-gray-500 hover:text-white'}"
			title="Random Mode - Return random response from available options"
			aria-label="Random response mode"
			on:click={() => handleResponseModeChange('random')}
			disabled={!selectedEndpoint}
		>
			<i class="fas fa-random text-xs"></i>
		</button>
		<button
			class="p-2 rounded-md text-sm transition-all duration-200 {selectedEndpoint?.response_mode ===
			'round_robin'
				? 'bg-purple-600 text-white shadow-md'
				: 'bg-gray-700 dark:bg-gray-600 text-gray-300 hover:bg-gray-600 dark:hover:bg-gray-500 hover:text-white'}"
			title="Round Robin Mode - Cycle through responses in order"
			aria-label="Round robin response mode"
			on:click={() => handleResponseModeChange('round_robin')}
			disabled={!selectedEndpoint}
		>
			<i class="fas fa-sync-alt text-xs"></i>
		</button>
	</div>
</div>
