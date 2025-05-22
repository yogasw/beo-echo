<script lang="ts">
	import type { Response, Endpoint, Project } from '$lib/api/BeoApi';
	import { updateEndpoint, resetEndpointsList, updateResponse } from '$lib/stores/saveButton';
	import StatusBodyTab from './routes/StatusBodyTab.svelte';
	import HeadersTab from './routes/HeadersTab.svelte';
	import RulesTab from './routes/RulesTab.svelte';
	import CallbacksTab from './routes/CallbacksTab.svelte';
	import ProxyTab from './routes/ProxyTab.svelte';
	import NotesTab from './routes/NotesTab.svelte';
	import RoutesList from '$lib/components/tabs/routes/RoutesList.svelte';
	import DropdownResponse from '$lib/components/tabs/routes/DropdownResponse.svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { selectedProject } from '$lib/stores/selectedConfig';

	export let activeContentTab = 'Status & Body';
	let activeConfigName = $selectedProject?.name || ''; // Store the active config name
	let endpoints: Endpoint[] = $selectedProject?.endpoints || []; // Store the list of endpoints
	let selectedEndpoint: Endpoint | null = null;
	let selectedResponse: Response | null = null;
	let filterText: string = ''; // Variable to store filter input
	
	// Update endpoints and activeConfigName when selectedProject changes
	$: {
		if ($selectedProject) {
			endpoints = $selectedProject.endpoints || [];
			activeConfigName = $selectedProject.name || '';
			selectedEndpoint = null; // Reset selected endpoint when project changes
			selectedResponse = null; // Reset selected response when project changes
		}
		console.log('Selected selectedProject:', $selectedProject);
	}

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

	function handleProxyChange(updatedEndpoint: Endpoint) {
		// Update the endpoint in the list
		const index = endpoints.findIndex((e) => e.id === updatedEndpoint.id);
		if (index !== -1) {
			endpoints[index] = updatedEndpoint;
			endpoints = [...endpoints]; // Trigger reactivity
		}
		selectedEndpoint = updatedEndpoint;
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
	/>

	<!-- Details Section -->
	<div class="w-2/3 {ThemeUtils.themeBgPrimary()} p-4 flex flex-col overflow-hidden">
		<div class="mb-4">
			<label
				for="endpoint-method"
				class="block text-sm font-bold mb-2 {ThemeUtils.themeTextPrimary()}">Endpoint</label
			>
			<div class="flex flex-col md:flex-row items-center space-y-2 md:space-y-0 md:space-x-2">
				<select
					id="endpoint-method"
					class="w-full md:w-1/6 rounded {ThemeUtils.themeBgSecondary()} px-4 py-2 {ThemeUtils.themeTextPrimary()}"
					value={selectedEndpoint?.method.toUpperCase()}
				>
					<option value="GET">GET</option>
					<option value="POST">POST</option>
					<option value="PUT">PUT</option>
					<option value="DELETE">DELETE</option>
					<option value="PATCH">PATCH</option>
				</select>
				<div class="flex items-center relative">
					<div class="hidden md:flex items-center rounded-md border border-blue-400/30 bg-blue-500/10 px-3 py-1 shadow-sm hover:border-blue-500/40 hover:bg-blue-500/15 transition-colors">
						<i class="fas fa-globe-americas text-blue-400 mr-1.5 text-[10px]"></i>
						<span class="text-blue-400 font-medium text-xs">API HOST</span>
						<div class="group relative ml-1.5">
							<i class="fas fa-info-circle text-blue-400 hover:text-blue-300 cursor-pointer transition-colors"></i>
							<div class="absolute z-10 hidden group-hover:block bg-gray-800 text-white text-xs rounded-md p-3 left-1/2 -translate-x-1/2 mt-2 w-auto min-w-[200px] max-w-xs whitespace-normal break-all shadow-lg border border-gray-700 transition-all duration-200 ease-in-out">
								<!-- Triangle pointer -->
								<div class="absolute -top-2 left-1/2 -translate-x-1/2 w-0 h-0 border-l-[6px] border-l-transparent border-r-[6px] border-r-transparent border-b-[6px] border-b-gray-800"></div>
								<div class="flex flex-col space-y-2">
									<div class="flex items-start">
										<i class="fas fa-link text-blue-400 mr-2 mt-0.5"></i>
										<span>{$selectedProject?.url || "No API host defined"}</span>
									</div>
									<div class="text-[10px] text-gray-400 italic">This is the base URL for all endpoints in this project</div>
								</div>
							</div>
						</div>
					</div>
				</div>
				<input
					type="text"
					class="w-full md:flex-1 rounded {ThemeUtils.themeBgSecondary()} px-4 py-2 {ThemeUtils.themeTextPrimary()}"
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
					class="{ThemeUtils.themeTextMuted()} hover:text-blue-500 disabled:{ThemeUtils.themeTextMuted(
						'opacity-50'
					)}"
					disabled={!selectedEndpoint || selectedEndpoint?.method !== 'GET'}
					aria-label="Open endpoint in a new tab"
					on:click={() => {
						let url = `${$selectedProject?.url || ""}${selectedEndpoint?.path ? selectedEndpoint.path : ''}`;
						// Open the URL in a new tab
						window.open(url, '_blank');
					}}
				>
					<i class="fas fa-external-link-alt"></i>
				</button>
			</div>
			<span class="{ThemeUtils.themeTextMuted()} block md:hidden mt-2"></span>
		</div>
		<label
			for="endpoint-documentation"
			class="block text-sm font-bold mb-2 {ThemeUtils.themeTextPrimary()}"
		>
			Documentation for this routes
		</label>
		<textarea
			id="endpoint-documentation"
			class="w-full rounded {ThemeUtils.themeBgSecondary()} px-4 py-2 {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()}"
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

		{#if selectedEndpoint}
			<div class="mb-4 mt-4">
				<!-- ProxyTab moved here, above DropdownResponse -->
				<ProxyTab endpoint={selectedEndpoint} onChange={handleProxyChange} />
			</div>
		{/if}

		<!-- Response section only shown when proxy is disabled -->
		{#if !selectedEndpoint?.use_proxy}
			<DropdownResponse bind:selectedEndpoint bind:selectedResponse />

			<div class="flex space-x-2 mb-4">
				{#each ['Status & Body', 'Headers', 'Rules', 'Callbacks', 'Notes'] as tab}
					{#if tab === activeContentTab}
						<button
							class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded"
							on:click={() => (activeContentTab = tab)}
						>
							{tab}
						</button>
					{:else}
						<button
							class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} py-2 px-4 rounded hover:bg-opacity-80"
							on:click={() => (activeContentTab = tab)}
						>
							{tab}
						</button>
					{/if}
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
									<RulesTab rules={selectedResponse?.rules || []} rulesOperator="AND" />
								{:else if activeContentTab === 'Callbacks'}
									<CallbacksTab callbacks={[]} />
								{:else if activeContentTab === 'Notes'}
									<NotesTab 
										notes={selectedResponse?.note || ''} 
										onSaveNotes={(notes) => {
											if (selectedResponse) {
												console.log('Notes saved:', notes);
												// Ensure we're not exceeding the backend character limit
												const trimmedNotes = notes.substring(0, 500);
												selectedResponse = updateResponse(
													'note',
													trimmedNotes,
													selectedEndpoint,
													selectedResponse
												);
											}
										}}
									/>
								{/if}
							{:else}
								<div class={ThemeUtils.themeTextMuted()}>Select a route to view details.</div>
							{/if}
						</div>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
