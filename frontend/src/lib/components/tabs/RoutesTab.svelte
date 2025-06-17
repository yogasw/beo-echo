<script lang="ts">
	import type { Response, Endpoint, Project } from '$lib/api/BeoApi';
	import { updateEndpoint, resetEndpointsList, updateResponse } from '$lib/stores/saveButton';
	import StatusBodyTab from './routes/StatusBodyTab.svelte';
	import RulesTab from './routes/RulesTab.svelte';
	import ProxyTab from './routes/ProxyTab.svelte';
	import NotesTab from './routes/NotesTab.svelte';
	import AdvancedSettingsTab from './routes/AdvancedSettingsTab.svelte';
	import RoutesList from '$lib/components/tabs/routes/RoutesList.svelte';
	import DropdownResponse from '$lib/components/tabs/routes/DropdownResponse.svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import HeadersEditor from '../common/HeadersEditor.svelte';
	import HttpMethodDropdown from '../common/HttpMethodDropdown.svelte';

	export let activeContentTab = 'Status & Body';
	let activeConfigName = $selectedProject?.name || ''; // Store the active config name
	let endpoints: Endpoint[] = $selectedProject?.endpoints || []; // Store the list of endpoints
	let selectedEndpoint: Endpoint | null = null;
	let selectedResponse: Response | null = null;
	let filterText: string = ''; // Variable to store filter input
	let localUseProxy: boolean = false; // Local state for proxy status
	let panelWidth: number = 33; // Panel width as percentage (33% = w-1/3)

	// Update endpoints and activeConfigName when selectedProject changes
	$: {
		if ($selectedProject) {
			endpoints = $selectedProject.endpoints || [];
			activeConfigName = $selectedProject.name || '';

			// Automatically select the first endpoint if available
			if (endpoints.length > 0) {
				selectRoute(endpoints[0]);
			} else {
				selectedEndpoint = null; // Reset selected endpoint when no endpoints available
				selectedResponse = null; // Reset selected response when no endpoints available
			}
		}
		console.log('Selected selectedProject:', $selectedProject);
	}

	// Update localUseProxy when selectedEndpoint changes
	$: {
		if (selectedEndpoint) {
			localUseProxy = selectedEndpoint.use_proxy || false;
		}
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

		// Update local proxy state from the selected endpoint
		localUseProxy = route.use_proxy || false;

		// Automatically select the first response if available
		if (route.responses && route.responses.length > 0) {
			selectedResponse = route.responses[0];
		} else {
			selectedResponse = null;
		}
	}

	function handleRouteStatusChange(route: Endpoint) {
		console.log('Route status changed:', route);
		const index = endpoints.findIndex((r) => r.id === route.id);
		if (index !== -1) {
			endpoints[index] = {
				...route
			};
			endpoints = [...endpoints]; // Trigger reactivity with a new array reference

			// Also update in the selectedProject store
			if ($selectedProject && $selectedProject.endpoints) {
				const projectEndpointIndex = $selectedProject.endpoints.findIndex((e) => e.id === route.id);
				if (projectEndpointIndex !== -1) {
					$selectedProject.endpoints[projectEndpointIndex] = route;
					selectedProject.set($selectedProject); // Trigger the store update
				}
			}
		}
	}

	function handleAddEndpoint(newEndpoint: Endpoint) {
		console.log('Endpoint added:', newEndpoint);

		// Add the new endpoint to the endpoints array
		endpoints = [...endpoints, newEndpoint];

		// Also update the selectedProject endpoints to make the change persist
		if ($selectedProject) {
			// If selectedProject.endpoints is undefined, initialize it as an empty array
			if (!$selectedProject.endpoints) {
				$selectedProject.endpoints = [];
			}

			// Add the new endpoint to the selectedProject's endpoints
			$selectedProject.endpoints = [...$selectedProject.endpoints, newEndpoint];

			// Update the selectedProject store to trigger reactivity
			selectedProject.set($selectedProject);
		}

		// Clear any active filter to ensure the new endpoint is visible
		if (filterText) {
			filterText = '';
		}

		setTimeout(() => {
			// Automatically select the new endpoint after a short delay
			selectRoute(newEndpoint);
		}, 100); // Adjust the delay as needed
	}

	function handleProxyChange(updatedEndpoint: Endpoint) {
		// Update local proxy state
		localUseProxy = updatedEndpoint.use_proxy || false;
		console.log('Proxy status updated, localUseProxy:', localUseProxy);
	}

	function handleHeadersSave(headers: string): void {
		console.log('Headers saved:', headers);
		if (selectedResponse) {
			selectedResponse = updateResponse('headers', headers, selectedEndpoint, selectedResponse);
		}
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
		bind:panelWidth
	/>

	<!-- Details Section -->
	<div class="{ThemeUtils.themeBgPrimary()} p-4 flex flex-col overflow-hidden" style="width: {100 - panelWidth}%;">
		<div class="mb-4">
			<label
				for="endpoint-method"
				class="block text-sm font-bold mb-2 {ThemeUtils.themeTextPrimary()}">Endpoint</label
			>
			<div class="flex flex-col md:flex-row items-center space-y-2 md:space-y-0 md:space-x-2">
				<!-- HTTP Method Dropdown - Enhanced with better styling -->
				<HttpMethodDropdown
					showLabel={false}
					showPlaceholder={false}
					className="w-full md:w-[120px]"
					value={selectedEndpoint?.method.toUpperCase()}
					placeholder={''}
					on:change={(e) => {
						if (selectedEndpoint) {
							let endpoint = e.detail.value;
							let target = e?.target as HTMLSelectElement;
							updateEndpoint('method', endpoint, selectedEndpoint);
							console.log('Updated endpoint method:', selectedEndpoint);
						}
					}}
				/>

				<!-- API Host Indicator - Enhanced with better tooltip and responsive design -->
				<div class="flex items-center relative">
					<div
						class="hidden md:flex items-center rounded-md border border-blue-400/30 bg-blue-500/10 px-3 py-2 shadow-sm hover:border-blue-500/40 hover:bg-blue-500/15 transition-colors cursor-pointer group"
					>
						<i class="fas fa-globe-americas text-blue-400 mr-1.5 text-xs"></i>
						<span class="text-blue-400 font-medium text-xs">API HOST</span>
						<i class="fas fa-chevron-down text-blue-400 opacity-50 ml-1 text-[10px]"></i>

						<!-- Improved tooltip -->
						<div
							class="absolute z-20 hidden group-hover:block bg-gray-800 text-white text-xs rounded-md p-4 left-0 mt-[42px] w-auto min-w-[240px] max-w-sm whitespace-normal break-all shadow-lg border border-gray-700 transition-all duration-200 ease-in-out"
						>
							<!-- Triangle pointer -->
							<div
								class="absolute -top-2 left-4 w-0 h-0 border-l-[6px] border-l-transparent border-r-[6px] border-r-transparent border-b-[6px] border-b-gray-800"
							></div>
							<div class="flex flex-col space-y-3">
								<div class="flex items-center">
									<div class="bg-blue-600/20 rounded-full p-1.5 mr-3">
										<i class="fas fa-link text-blue-400"></i>
									</div>
									<div class="flex flex-col">
										<span class="font-semibold mb-1">Base URL</span>
										<span class="font-mono text-blue-300"
											>{$selectedProject?.url || 'No API host defined'}</span
										>
									</div>
								</div>
								<div class="text-xs text-gray-400">
									This is the base URL for all endpoints in this project. It will be prepended to
									all endpoint paths when requests are processed.
								</div>
							</div>
						</div>
					</div>
				</div>

				<!-- Path Input Field - Enhanced with styling and icons -->
				<div class="relative w-full md:flex-1">
					<div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
						<i class="fas fa-route text-gray-400 dark:text-gray-500"></i>
					</div>
					<input
						type="text"
						class="w-full rounded-md {ThemeUtils.themeBgSecondary()} pl-10 pr-10 py-2 {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
						value={selectedEndpoint?.path}
						placeholder="/api/resource/id"
						on:blur={(e) => {
							if (selectedEndpoint) {
								let target = e?.target as HTMLInputElement;
								selectedEndpoint = updateEndpoint('path', target?.value || '', selectedEndpoint);
							}
						}}
					/>

					<!-- Enhanced Open in New Tab Button -->
					<div class="absolute inset-y-0 right-0 flex items-center pr-3">					<button
						class="p-1.5 rounded-md hover:bg-blue-500/10 {ThemeUtils.themeTextMuted()} hover:text-blue-500 transition-colors disabled:opacity-50 disabled:hover:bg-transparent"
						disabled={!selectedEndpoint || selectedEndpoint?.method !== 'GET'}
						aria-label="Open endpoint in a new tab"
						title={selectedEndpoint?.method !== 'GET'
							? 'Only GET endpoints can be opened directly'
							: 'Open endpoint in new tab'}
							on:click={() => {
								let url = `${$selectedProject?.url || ''}${selectedEndpoint?.path ? selectedEndpoint.path : ''}`;
								window.open(url, '_blank');
							}}
						>
							<i class="fas fa-external-link-alt"></i>
						</button>
					</div>
				</div>
			</div>
			<span class="{ThemeUtils.themeTextMuted()} block md:hidden mt-2"></span>
		</div>
		
		<div class="mb-2">
			<label
				for="endpoint-documentation"
				class="block text-sm font-bold mb-2 {ThemeUtils.themeTextPrimary()} flex items-center"
			>
				<i class="fas fa-file-alt mr-2 text-blue-500"></i>
				Documentation
			</label>
			<div class="relative">
				<div class="absolute inset-y-0 left-0 pl-3 pt-3 pointer-events-none">
					<i class="fas fa-pencil-alt {ThemeUtils.themeTextMuted()}"></i>
				</div>
				<textarea
					id="endpoint-documentation"
					class="w-full rounded-md {ThemeUtils.themeBgSecondary()} pl-10 px-4 py-3 {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all"
					rows="3"
					placeholder="Describe what this endpoint does, expected parameters, and example responses..."
					on:blur={(e) => {
						if (selectedEndpoint) {
							let target = e?.target as HTMLTextAreaElement;
							selectedEndpoint = updateEndpoint(
								'documentation',
								target?.value || '',
								selectedEndpoint
							);
						}
					}}>{selectedEndpoint?.documentation || ''}</textarea
				>
				<div class="absolute bottom-3 right-3 {ThemeUtils.themeTextMuted()} text-xs opacity-70">
					<i class="fas fa-info-circle mr-1"></i>
					Good documentation helps team members understand the API
				</div>
			</div>
		</div>

		{#if selectedEndpoint}
			<div class="mb-2">
				<!-- ProxyTab endpoint={selectedEndpoint} onChange={handleProxyChange} /-->
				<ProxyTab
					endpoint={selectedEndpoint}
					onChange={handleProxyChange}
					bind:use_proxy={localUseProxy}
				/>
			</div>
		{/if}

		<!-- Response section only shown when proxy is disabled -->
		{#if !localUseProxy}
			<DropdownResponse bind:selectedEndpoint bind:selectedResponse />

			<!-- Enhanced Tab Navigation -->
			<div class="flex mb-4 border-b {ThemeUtils.themeBorder()} overflow-x-auto no-scrollbar">
				{#each [{ id: 'Status & Body', icon: 'fas fa-code' }, { id: 'Headers', icon: 'fas fa-exchange-alt' }, { id: 'Rules', icon: 'fas fa-filter' }, { id: 'Notes', icon: 'fas fa-sticky-note' }, { id: 'Advanced Settings', icon: 'fas fa-cogs' }] as tab}				<button
					class="relative flex items-center py-3 px-4 font-medium text-sm whitespace-nowrap transition-all duration-200 {tab.id ===
					activeContentTab
						? `${ThemeUtils.themeTextPrimary()} border-b-2 border-blue-500`
						: `${ThemeUtils.themeTextMuted()} hover:${ThemeUtils.themeTextPrimary('opacity-80')}`}"
					on:click={() => (activeContentTab = tab.id)}
					aria-label="Switch to {tab.id} tab"
					title="Switch to {tab.id} tab"
				>
						<i class="{tab.icon} mr-2 {tab.id === activeContentTab ? 'text-blue-500' : ''}"></i>
						{tab.id}
						<!-- Active indicator for current tab -->
						{#if tab.id === activeContentTab}
							<span
								class="absolute bottom-0 left-0 h-0.5 w-full bg-blue-500 transform transition-transform"
							></span>
						{/if}
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
									<HeadersEditor
										onSave={handleHeadersSave}
										headers={selectedResponse?.headers || '{}'}
									/>
								{:else if activeContentTab === 'Rules'}
									<RulesTab
										projectId={$selectedProject?.id || ''}
										rules={selectedResponse?.rules || []}
										endpointId={selectedEndpoint?.id || ''}
										responseId={selectedResponse?.id || ''}
									/>
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
								{:else if activeContentTab === 'Advanced Settings'}
									<AdvancedSettingsTab
										delayMs={selectedResponse?.delay_ms || 0}
										onDelayUpdate={(newDelayMs) => {
											if (selectedResponse) {
												console.log('Delay updated:', newDelayMs);
												selectedResponse = updateResponse(
													'delay_ms',
													newDelayMs,
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
