<script lang="ts">
	import type { Endpoint, Response } from '$lib/api/mockoonApi';
	import type { MockoonResponse, MockoonRoute } from '$lib/types/Config';

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

	$: {
		if (selectedResponse) {
			selectedValue = `Response 1 (${selectedResponse.status_code}) ${selectedResponse.documentation}`;
		} else {
			selectedValue = 'No Response';
		}
	}
</script>

<div class="flex items-center justify-between bg-gray-800 text-gray-300 py-5">
	<div class="flex items-center w-full">
		<button class="bg-blue-500 text-white rounded px-2 py-1 text-sm font-medium">
			<i class="fas fa-plus"></i>
		</button>
		<div class="relative ml-4 w-full">
			<button
				class="text-sm font-medium bg-gray-700 text-gray-300 rounded px-2 py-1 flex items-center justify-between w-full"
				on:click={() => { toggleDropdown() }}
			>
				<span id="selectedValue">{selectedValue}</span>
				<i class="fas fa-chevron-down"></i>
			</button>
			<div id="dropdownMenu" class="absolute mt-1 bg-gray-700 text-gray-300 rounded shadow-lg w-full hidden">
				<ul class="text-sm">
					{#if selectedEndpoint?.responses && selectedEndpoint.responses.length > 0}
						{#each selectedEndpoint.responses as response, index}
							<li>
								<button type="button" class="w-full text-left px-4 py-2 hover:bg-gray-600 cursor-pointer"
												on:click={() => { selectResponse(index, response) }}>
									Response {index + 1} ({response.status_code}) {response?.documentation}
								</button>
							</li>
						{/each}
					{:else}
						<li class="px-4 py-2 text-gray-500">No Response</li>
					{/if}
				</ul>
			</div>
		</div>
	</div>
</div>
