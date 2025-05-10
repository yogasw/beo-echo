<script lang="ts">
	import type { Endpoint } from '$lib/api/mockoonApi';
	import type { MockoonRoute } from '$lib/types/Config';
	export let selectedEndpoint: Endpoint | null;
	export let activeConfigName: string;
	export let filterText: string;
	export let filteredEndpoints: Endpoint[];
	export let selectRoute: (route: Endpoint) => void;
	export let handleRouteStatusChange: (route: Endpoint) => void;
</script>

<!-- Routes Section -->
<div class="w-1/3 bg-gray-800 p-4 flex flex-col">
	<div class="bg-gray-700 p-4 rounded mb-4">
		<div class="flex items-center">
			<i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
			<span class="text-xl font-bold text-blue-500">Editing Configuration:</span>
		</div>
		<span class="text-xl font-bold text-blue-500">{activeConfigName}</span>
	</div>
	<div class="flex items-center bg-gray-700 p-2 rounded mb-4">
		<i class="fas fa-search text-white text-lg mr-2"></i>
		<input
			type="text"
			id="route-search"
			placeholder="Search Path or Method"
			class="w-full bg-gray-700 text-white py-1 px-2 rounded text-sm"
			bind:value={filterText}
		/>
	</div>

	<div class="flex-1 overflow-y-auto hide-scrollbar">
		<div class="space-y-4 pr-2 py-2">
			{#each filteredEndpoints as endpoint}
				<div
					class="flex items-center justify-between bg-gray-700 p-4 rounded cursor-pointer {selectedEndpoint ===
					endpoint
						? 'border-2 border-blue-500'
						: ''}"
					on:click={() => selectRoute(endpoint)}
					on:keydown={(e) => e.key === 'Enter' && selectRoute(endpoint)}
					tabindex="0"
					role="button"
				>
					<span class="text-sm font-bold truncate">
						<strong>{endpoint.method}</strong>
						{endpoint.path.length > 30 ? endpoint.path.slice(0, 30) + '...' : endpoint.path}
					</span>
					<button
						class="text-white py-1 px-2 rounded flex items-center"
						class:bg-green-500={endpoint.enabled == true}
						class:bg-red-500={endpoint.enabled == false}
						on:click|stopPropagation={() => {
							endpoint.enabled = !endpoint.enabled;
							// Call the function to handle route status change';
							handleRouteStatusChange(endpoint);
						}}
					>
						{#if endpoint.enabled === false}
							<i class="fas fa-toggle-off mr-1"></i> Disable
						{:else}
							<i class="fas fa-toggle-on mr-1"></i> Enable
						{/if}
					</button>
				</div>
			{/each}
		</div>
	</div>
</div>
