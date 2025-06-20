<script lang="ts">
	import type { Endpoint, Response } from '$lib/api/BeoApi';
	import { addResponse, deleteResponse } from '$lib/api/BeoApi';
	import { updateEndpoint } from '$lib/stores/saveButton';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	let selectedValue: string = '';
	export let selectedEndpoint: Endpoint | null;
	export let selectedResponse: Response | null;

	let flaggedResponseId: string | null = null;

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

		// Sanitize the note text to remove HTML and collapse whitespace
		noteText = sanitizeText(noteText);

		// Use different truncation length depending on whether this is for the selected value or dropdown items
		const maxLength = isForSelectedValue ? 40 : 25;

		// Format display with truncation
		return `Response ${index + 1} ${statusText} ${truncateText(noteText, maxLength)}`;
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
			// Find the index of the selected response in the endpoint's responses
			const index = selectedEndpoint.responses.findIndex((r) => r.id === selectedResponse?.id);
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
		toast.success('Response deleted successfully');
	}
</script>

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
						{#each selectedEndpoint.responses as response, index}
							<li class="flex items-center">
								<button
									type="button"
									class="w-full text-left px-4 py-2 {ThemeUtils.themeHover()} cursor-pointer"
									on:click={() => {
										selectResponse(index, response);
									}}
									title="Select this response"
									aria-label="Select this response"
								>
									Response {index + 1} ({response.status_code}) {truncateText(
										sanitizeText(response?.note)
									)}
								</button>
								<!-- Copy & Delete actions -->
								<div class="flex items-center space-x-1 ml-2">
									<button
										class="p-1 rounded hover:bg-gray-700"
										title="Flag"
										aria-label="Flag"
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
										class="p-1 rounded hover:bg-gray-700"
										title="Copy"
										aria-label="Copy"
										on:click|stopPropagation={() => {
											console.log('Copy', response);
										}}
									>
										<i class="fas fa-copy text-gray-300"></i>
									</button>
									<button
										class="p-1 rounded hover:bg-gray-700"
										title="Delete"
										aria-label="Delete"
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
