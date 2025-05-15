<script lang="ts">
	import type { Endpoint, Response } from '$lib/api/BeoApi';
	import { addResponse } from '$lib/api/BeoApi';
	import type { MockoonResponse, MockoonRoute } from '$lib/types/Config';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	let selectedValue: string = '';
	export let selectedEndpoint: Endpoint | null;
	export let selectedResponse: Response | null;

	const toggleDropdown = (): void => {
		const dropdown = document.getElementById('dropdownMenu');
		if (dropdown) {
			dropdown.classList.toggle('hidden');
		}
	};

	const selectResponse = (index: number, value: Response): void => {
		selectedValue = `Response ${index + 1} (${value.status_code}) ${value.documentation}`;
		const selectedElement = document.getElementById('selectedValue');
		if (selectedElement) {
			selectedElement.innerText = selectedValue;
		}
		toggleDropdown();
		selectedResponse = value;
	};

	const handleAddResponse = async (): Promise<void> => {
		if (selectedEndpoint) {
			try {
				// Sample response values as specified
				const statusCode = 200;
				const body = "{\"message\":\"Hello World\"}";
				const headers = "{\"Content-Type\":\"application/json\"}";
				
				const newResponse = await addResponse(
					selectedEndpoint.project_id, 
					selectedEndpoint.id,
					statusCode,
					body,
					headers,
					"",
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

	$: {
		if (selectedResponse) {
			selectedValue = `Response 1 (${selectedResponse.status_code}) ${selectedResponse.documentation}`;
		} else {
			selectedValue = 'No Response';
		}
	}
</script>

<div class="flex items-center justify-between {ThemeUtils.themeBgPrimary()} {ThemeUtils.themeTextSecondary()} py-2">
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
				on:click={() => { toggleDropdown() }}
			>
				<span id="selectedValue">{selectedValue}</span>
				<i class="fas fa-chevron-down"></i>
			</button>
			<div id="dropdownMenu" class="absolute mt-1 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextSecondary()} rounded shadow-lg w-full hidden z-50 max-h-60 overflow-y-auto">
				<ul class="text-sm">
					{#if selectedEndpoint?.responses && selectedEndpoint.responses.length > 0}
						{#each selectedEndpoint.responses as response, index}
							<li>
								<button type="button" class="w-full text-left px-4 py-2 {ThemeUtils.themeHover()} cursor-pointer"
												on:click={() => { selectResponse(index, response) }}>
									Response {index + 1} ({response.status_code}) {response?.documentation}
								</button>
							</li>
						{/each}
					{:else}
						<li class="px-4 py-2 {ThemeUtils.themeTextMuted()}">No Response</li>
					{/if}
				</ul>
			</div>
		</div>
	</div>
</div>
