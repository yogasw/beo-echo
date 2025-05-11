<script lang="ts">
	import type { Response, Endpoint, Project } from '$lib/api/mockoonApi';
	import { updateEndpoint, resetEndpointsList, updateResponse } from '$lib/stores/saveButton';
	import StatusBodyTab from './routes/StatusBodyTab.svelte';
	import HeadersTab from './routes/HeadersTab.svelte';
	import RulesTab from './routes/RulesTab.svelte';
	import CallbacksTab from './routes/CallbacksTab.svelte';
	import RoutesList from '$lib/components/tabs/routes/RoutesList.svelte';
	import DropdownResponse from '$lib/components/tabs/routes/DropdownResponse.svelte';

	export let selectedProject: Project;
	export let endpoints: Endpoint[];
	export let activeContentTab = 'Status & Body';
	let activeConfigName = selectedProject?.name || ''; // Store the active config name

	let selectedEndpoint: Endpoint | null = null;
	let selectedResponse: Response | null = null;
	let filterText: string = ''; // Variable to store filter input

	$: filteredEndpoints = endpoints.filter((endpoint) => {
		if (!filterText.trim()) return true;
		const filterParts = filterText.toLowerCase().split(' ');
		return filterParts.every(
			(part) =>
				endpoint.path.toLowerCase().includes(part) || endpoint.method.toLowerCase().includes(part)
		);
	});

	function selectRoute(route: Endpoint) {
		console.log('Route selected:', route);
		selectedEndpoint = route;
		// Reset endpoints update list when changing endpoints
		resetEndpointsList();
		// selectedResponse = route.responses[0] || null; // Select the first response by default
	}

	function handleRouteStatusChange(route: Endpoint) {
		console.log('Route status changed:', route);
		const index = endpoints.findIndex((r) => r.path === route.path && r.method === route.method);
		if (index !== -1) {
			endpoints[index] = {
				...route
			};
			endpoints = endpoints; // Trigger reactivity
		}
	}

	function handleAddEndpoint(newEndpoint: Endpoint) {
		console.log('Endpoint added:', newEndpoint);
		endpoints = [...endpoints, newEndpoint];
		selectRoute(newEndpoint);
	}
</script>

<div class="flex flex-1 h-full">
	<RoutesList
		{activeConfigName}
		{selectedEndpoint}
		bind:filterText
		{filteredEndpoints}
		{selectRoute}
		{handleRouteStatusChange}
		{handleAddEndpoint}
		project={selectedProject}
	/>

	<!-- Details Section -->
	<div class="w-2/3 bg-gray-800 p-4 flex flex-col overflow-hidden">
		<div class="mb-4">
			<label for="endpoint-method" class="block text-sm font-bold mb-2">Endpoint</label>
			<div class="flex flex-col md:flex-row items-center space-y-2 md:space-y-0 md:space-x-2">
				<select
					id="endpoint-method"
					class="w-full md:w-1/6 rounded bg-gray-700 px-4 py-2 text-white"
					value={selectedEndpoint?.method.toUpperCase()}
				>
					<option value="GET">GET</option>
					<option value="POST">POST</option>
					<option value="PUT">PUT</option>
					<option value="DELETE">DELETE</option>
					<option value="PATCH">PATCH</option>
				</select>
				<span class="text-gray-400 hidden md:block">{selectedProject.url}</span>
				<input
					type="text"
					class="w-full md:flex-1 rounded bg-gray-700 px-4 py-2 text-white"
					value={selectedEndpoint?.path}
					on:blur={(e) => {
						if (selectedEndpoint) {
							let target = e?.target as HTMLInputElement;
							selectedEndpoint = updateEndpoint('path', target?.value || '', selectedEndpoint);
							console.log('Updated endpoint:', selectedEndpoint);
						}
					}}
				/>
				<button
					class="text-gray-400 hover:text-blue-500 disabled:text-gray-600"
					disabled={!selectedEndpoint || selectedEndpoint?.method !== 'GET'}
					aria-label="Open endpoint in a new tab"
					on:click={() => {
						let url = `${selectedProject.url}${selectedEndpoint?.path ? selectedEndpoint.path : ''}`;
						// Open the URL in a new tab
						window.open(url, '_blank');
					}}
				>
					<i class="fas fa-external-link-alt"></i>
				</button>
			</div>
			<span class="text-gray-400 block md:hidden mt-2"></span>
		</div>
		<label for="endpoint-documentation" class="block text-sm font-bold mb-2">
			Documentation for this routes
		</label>
		<textarea
			id="endpoint-documentation"
			class="w-full rounded bg-gray-700 px-4 py-2 text-white"
			rows="3"
			placeholder="Provide a brief description or documentation for this endpoint"
			on:blur={(e) => {
				if (selectedEndpoint) {
					let target = e?.target as HTMLTextAreaElement;
					selectedEndpoint = updateEndpoint('documentation', target?.value || '', selectedEndpoint);
					console.log('Updated endpoint:', selectedEndpoint);
				}
			}}>{selectedEndpoint?.documentation || ''}</textarea
		>

		<DropdownResponse bind:selectedEndpoint bind:selectedResponse />

		<div class="flex space-x-2 mb-4">
			{#each ['Status & Body', 'Headers', 'Rules', 'Callbacks'] as tab}
				<button
					class="text-white py-2 px-4 rounded"
					class:bg-blue-500={tab === activeContentTab}
					class:bg-gray-700={tab !== activeContentTab}
					on:click={() => (activeContentTab = tab)}
				>
					{tab}
				</button>
			{/each}
		</div>

		<div class="flex-1 overflow-auto h-full">
			<div class="max-w-full overflow-x-auto h-full">
				<div class="min-w-0 h-full">
					<div class="h-full flex flex-col">
						{#if selectedEndpoint}
							{#if activeContentTab === 'Status & Body'}
								<StatusBodyTab
									responseBody={selectedResponse?.body || ''}
									statusCode={selectedResponse?.status_code || 200}
									onStatusCodeChange={(val) => {
										if (selectedResponse) {
											console.log('Status code changed:', val);
											selectedResponse = updateResponse(
												'status_code',
												val,
												selectedEndpoint,
												selectedResponse
											);
										}
									}}
									onSaveButtonClick={(content) => {
										console.log('Save button clicked with content:', content);
										if (selectedResponse) {
											selectedResponse = updateResponse(
												'body',
												content,
												selectedEndpoint,
												selectedResponse
											);
										}
									}}
								/>
							{:else if activeContentTab === 'Headers'}
								<HeadersTab headers={selectedResponse?.headers || '{}'} />
							{:else if activeContentTab === 'Rules'}
								<RulesTab
									rules={selectedResponse?.rules || []}
									rulesOperator={selectedResponse?.rulesOperator}
								/>
							{:else if activeContentTab === 'Callbacks'}
								<CallbacksTab callbacks={[]} />
							{/if}
						{:else}
							<div class="text-gray-400">Select a route to view details.</div>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
