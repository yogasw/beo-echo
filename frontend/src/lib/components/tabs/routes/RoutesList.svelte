<script lang="ts">
	import { updateEndpoint, type Endpoint, type Project } from '$lib/api/mockoonApi';
	import { toast } from '$lib/stores/toast';
	import type { MockoonRoute } from '$lib/types/Config';
	import AddEndpointModal from './AddEndpointModal.svelte';
	import { onMount, onDestroy } from 'svelte';

	export let selectedEndpoint: Endpoint | null;
	export let activeConfigName: string;
	export let filterText: string;
	export let filteredEndpoints: Endpoint[];
	export let selectRoute: (route: Endpoint) => void;
	export let handleRouteStatusChange: (route: Endpoint) => void;
	export let handleAddEndpoint: (endpoint: Endpoint) => void;
	export let project: Project;
	
	let showAddEndpointModal = false;
	let activeMenuEndpointId: string | null = null;

	function onEndpointCreated(event: CustomEvent<Endpoint>) {
		handleAddEndpoint(event.detail);
		showAddEndpointModal = false;
	}

	function toggleMenu(event: MouseEvent, endpointId: string) {
		event.stopPropagation();
		if (activeMenuEndpointId === endpointId) {
			activeMenuEndpointId = null;
		} else {
			activeMenuEndpointId = endpointId;
		}
	}

	function handleMenuAction(event: MouseEvent, action: string, endpoint: Endpoint) {
		event.stopPropagation();
		activeMenuEndpointId = null;
		
		switch(action) {
			case 'enable':
			case 'disable':
				endpoint.enabled = action === 'enable';
				handleRouteStatusChange(endpoint);
				updateEndpoint(endpoint.project_id, endpoint.id, {
					enabled: endpoint.enabled,
				})
					.then(() => {
						toast.success(`Endpoint successfully ${action}d!`);
					})
					.catch((error) => {
						toast.error(`Failed to ${action} endpoint: ${error.message}`);
					});
				break;
			case 'duplicate':
				// Add your duplicate functionality here
				console.log('Duplicate endpoint', endpoint);
				// Call API or service to duplicate the endpoint
				break;
			case 'delete':
				// Add your delete functionality here
				console.log('Delete endpoint', endpoint);
				// Call API or service to delete the endpoint
				break;
		}
	}

	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (activeMenuEndpointId) {
			activeMenuEndpointId = null;
		}
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
	});

	onDestroy(() => {
		document.removeEventListener('click', handleClickOutside);
	});
</script>

<!-- Routes Section -->
<div class="w-1/3 bg-gray-800 p-4 flex flex-col">
	<div class="bg-gray-700 p-4 rounded mb-4 flex items-center">
		<i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
		<span class="text-xl font-bold text-blue-500">Project: {activeConfigName}</span>
	</div>
	<div class="flex items-center bg-gray-700 p-2 rounded mb-2">
		<i class="fas fa-search text-white text-lg mr-2"></i>
		<input
			type="text"
			id="route-search"
			placeholder="Search Path or Method"
			class="w-full bg-gray-700 text-white py-1 px-2 rounded text-sm"
			bind:value={filterText}
		/>
	</div>
	
	<!-- Add Endpoint Button -->
	<button 
		class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded mb-4 flex items-center justify-center"
		on:click={() => showAddEndpointModal = true}
	>
		<i class="fas fa-plus mr-2"></i> Add Endpoint
	</button>
	
	<AddEndpointModal 
		bind:isOpen={showAddEndpointModal} 
		{project}
		on:endpointCreated={onEndpointCreated}
		on:close={() => showAddEndpointModal = false}
	/>

	<div class="flex-1 overflow-y-auto hide-scrollbar">
		<div class="space-y-4 pr-2 py-2">
			{#each filteredEndpoints as endpoint}
				<div
					class="flex items-center justify-between bg-gray-700 py-2 px-4 rounded cursor-pointer relative group {selectedEndpoint ===
					endpoint
						? 'border-2 border-blue-500'
						: ''} {!endpoint.enabled ? 'disabled-endpoint' : ''}"
					on:click={() => selectRoute(endpoint)}
					on:keydown={(e) => e.key === 'Enter' && selectRoute(endpoint)}
					tabindex="0"
					role="button"
				>
					{#if !endpoint.enabled}
						<div class="absolute left-0 top-0 bottom-0 w-1 bg-red-500"></div>
					{/if}
					<span class="text-sm font-bold truncate">
						<strong>{endpoint.method}</strong>
						{endpoint.path.length > 30 ? endpoint.path.slice(0, 30) + '...' : endpoint.path}
					</span>
					
					<!-- Three-dot menu button only shown on hover -->
					<div class="relative menu-container">
						<button
							class="text-white h-8 w-8 flex items-center justify-center rounded hover:bg-gray-600 focus:outline-none opacity-0 group-hover:opacity-100 hover:opacity-100"
							on:click|stopPropagation={(e) => toggleMenu(e, endpoint.id)}
							aria-label="Options menu"
						>
							<div class="flex flex-col space-y-0.5">
								<div class="w-1 h-1 rounded-full bg-white"></div>
								<div class="w-1 h-1 rounded-full bg-white"></div>
								<div class="w-1 h-1 rounded-full bg-white"></div>
							</div>
						</button>
						
						{#if activeMenuEndpointId === endpoint.id}
							<div class="absolute right-0 top-full mt-2 w-48 bg-gray-800 border border-gray-600 rounded shadow-lg z-50">
								<div class="py-1">
									<button
										class="w-full text-left px-4 py-2 text-sm text-white hover:bg-gray-700 flex items-center"
										on:click|stopPropagation={(e) => handleMenuAction(e, endpoint.enabled ? 'disable' : 'enable', endpoint)}
									>
										{#if endpoint.enabled}
											<i class="fas fa-toggle-off mr-2 text-red-500"></i> Disable
										{:else}
											<i class="fas fa-toggle-on mr-2 text-green-500"></i> Enable
										{/if}
									</button>
									<button
										class="w-full text-left px-4 py-2 text-sm text-white hover:bg-gray-700 flex items-center"
										on:click|stopPropagation={(e) => handleMenuAction(e, 'duplicate', endpoint)}
									>
										<i class="fas fa-clone mr-2"></i> Duplicate
									</button>
									<hr class="border-gray-600 my-1" />
									<button
										class="w-full text-left px-4 py-2 text-sm text-red-400 hover:bg-gray-700 flex items-center"
										on:click|stopPropagation={(e) => handleMenuAction(e, 'delete', endpoint)}
									>
										<i class="fas fa-trash-alt mr-2"></i> Delete
									</button>
								</div>
							</div>
						{/if}
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>
