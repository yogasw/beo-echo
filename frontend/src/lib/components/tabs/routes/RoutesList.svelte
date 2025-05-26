<script lang="ts">
	import { deleteEndpoint, updateEndpoint, type Endpoint, type RequestLog } from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';
	import { onMount, onDestroy } from 'svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { currentWorkspace } from '$lib/stores/workspace';
	import ModalCreateMock from '../logs/ModalCreateMock.svelte';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { getRoutesPanelWidth, setRoutesPanelWidth } from '$lib/utils/localStorage';

	export let selectedEndpoint: Endpoint | null;
	export let activeConfigName: string;
	export let filterText: string;
	export let filteredEndpoints: Endpoint[];
	export let selectRoute: (route: Endpoint) => void;
	export let handleRouteStatusChange: (route: Endpoint) => void;
	export let handleAddEndpoint: (endpoint: Endpoint) => void;
	export let panelWidth: number = 33; // Panel width as percentage (33% = w-1/3) - will be initialized in onMount
	let defaultRequestLog: RequestLog = {
		id: '',
		project_id: $selectedProject?.id || '',
		method: 'GET',
		path: '/',
		query_params: '',
		request_headers: `{}`,
		request_body: '',
		response_status: 200,
		response_body: '',
		response_headers: '{"Content-Type": "application/json"}',
		latency_ms: 0,
		execution_mode: 'mock',
		matched: false,
		bookmark: false,
		created_at: new Date(),
	};

	let showAddEndpointModal = false;
	let activeMenuEndpointId: string | null = null;
	// Reference to the endpoints container for scrolling
	let endpointsContainer: HTMLDivElement;
	
	// Resizable panel variables
	let isResizing = false;
	let startX = 0;
	let startWidth = panelWidth;

	function onEndpointCreated(event: CustomEvent<Endpoint>) {
		handleAddEndpoint(event.detail);
		showAddEndpointModal = false;
		
		// Scroll to the bottom of the endpoints list after a short delay to ensure rendering is complete
		setTimeout(() => {
			if (endpointsContainer) {
				endpointsContainer.scrollTop = endpointsContainer.scrollHeight;
			}
		}, 150);
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

		switch (action) {
			case 'enable':
			case 'disable':
				if (!$currentWorkspace) {
					toast.error('No workspace selected');
					return;
				}

				endpoint.enabled = action === 'enable';
				handleRouteStatusChange(endpoint);
				updateEndpoint(endpoint.project_id, endpoint.id, {
					enabled: endpoint.enabled
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
				deleteEndpoint(endpoint.project_id, endpoint.id)
					.then(() => {
						toast.success('Endpoint successfully deleted!');
						handleRouteStatusChange(endpoint);
					})
					.catch((error) => {
						toast.error(`Failed to delete endpoint: ${error.message}`);
					});
				break;
		}
	}

	// Close menu when clicking outside
	function handleClickOutside(event: MouseEvent) {
		if (activeMenuEndpointId) {
			activeMenuEndpointId = null;
		}
	}

	// Resize functions
	function startResize(event: MouseEvent) {
		isResizing = true;
		startX = event.clientX;
		startWidth = panelWidth;
		
		document.addEventListener('mousemove', handleResize);
		document.addEventListener('mouseup', stopResize);
		document.body.style.cursor = 'col-resize';
		document.body.style.userSelect = 'none';
	}

	function handleResize(event: MouseEvent) {
		if (!isResizing) return;
		
		const deltaX = event.clientX - startX;
		const containerWidth = window.innerWidth;
		const newWidth = startWidth + (deltaX / containerWidth) * 100;
		
		// Constrain between 20% and 60%
		panelWidth = Math.min(Math.max(newWidth, 20), 60);
	}

	function stopResize() {
		isResizing = false;
		document.removeEventListener('mousemove', handleResize);
		document.removeEventListener('mouseup', stopResize);
		document.body.style.cursor = '';
		document.body.style.userSelect = '';
		
		// Save panel width to localStorage when resize is complete
		setRoutesPanelWidth(panelWidth);
	}

	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		
		// Initialize panel width from localStorage
		panelWidth = getRoutesPanelWidth();
	});

	onDestroy(() => {
		document.removeEventListener('click', handleClickOutside);
	});
	console.log("projectId", selectedEndpoint?.project_id);
</script>

<!-- Routes Section -->
<div class="flex flex-col theme-bg-primary relative border-r border-theme" style="width: {panelWidth}%;">
	<div class="pr-2 flex flex-col h-full">
		<div class={ThemeUtils.headerSection('rounded mb-4')}>
			<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
				<i class="fas fa-route text-blue-500 text-xl"></i>
			</div>
			<div>
				<h2 class="text-xl font-bold theme-text-primary">{activeConfigName}</h2>
				<p class="text-sm theme-text-muted">API Endpoint Management</p>
			</div>
		</div>
		<div class="relative mb-6">
			<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
				<i class="fas fa-search theme-text-muted"></i>
			</div>
			<input
				type="text"
				placeholder="Search Path or Method"
				bind:value={filterText}
				class={ThemeUtils.inputField('p-3 ps-10 text-sm rounded-lg')}
			/>
		</div>

		<!-- Add Endpoint Button -->
		<button
			class={ThemeUtils.primaryButton('w-full justify-center mb-4')}
			on:click={() => (showAddEndpointModal = true)}
		>
			<i class="fas fa-plus mr-2"></i> Add Endpoint
		</button>

		<ModalCreateMock
			bind:isOpen={showAddEndpointModal}
			projectId={$selectedProject?.id || ''}
			onClose={() => (showAddEndpointModal = false)}
			onSuccess={() => (showAddEndpointModal = false)}
			on:endpointCreated={onEndpointCreated}
			on:close={() => (showAddEndpointModal = false)}
			log={defaultRequestLog}
		/>

		<div class="flex-1 overflow-y-auto hide-scrollbar" bind:this={endpointsContainer}>
			<div class="space-y-4 py-2">
				{#each filteredEndpoints as endpoint}
					<div
						class={ThemeUtils.themeBgSecondary(`flex items-center justify-between py-2 px-4 rounded cursor-pointer relative group 
							${selectedEndpoint === endpoint ? 'border-2 border-blue-500' : 'theme-border'}`)}
						on:click={() => selectRoute(endpoint)}
						on:keydown={(e) => e.key === 'Enter' && selectRoute(endpoint)}
						tabindex="0"
						role="button"
					>
						{#if !endpoint.enabled}
							<div class="absolute left-0 top-0 bottom-0 w-1 bg-red-500 rounded-bl rounded-tl"></div>
						{/if}
						<span class={ThemeUtils.themeTextPrimary(`flex items-center text-sm font-bold truncate ${!endpoint.enabled ? 'opacity-75' : ''}`)}>
							<span class={ThemeUtils.methodBadge(endpoint.method, 'mr-2')}>
								{endpoint.method}
							</span>
							{endpoint.path.length > 30 ? endpoint.path.slice(0, 30) + '...' : endpoint.path}
							{#if endpoint.use_proxy}
								<span class="ml-2 px-2 py-0.5 text-xs font-medium rounded-full bg-purple-600 text-white">
									Proxy
								</span>
							{/if}
						</span>

						<!-- Three-dot menu button only shown on hover -->
						<div class="relative menu-container">
							<button
								class="theme-text-primary h-8 w-8 flex items-center justify-center rounded hover:bg-gray-600 focus:outline-none opacity-0 group-hover:opacity-100 hover:opacity-100"
								on:click|stopPropagation={(e) => toggleMenu(e, endpoint.id)}
								aria-label="Options menu"
							>
								<div class="flex flex-col space-y-0.5">
									<div class="w-1 h-1 rounded-full theme-bg-accent"></div>
									<div class="w-1 h-1 rounded-full theme-bg-accent"></div>
									<div class="w-1 h-1 rounded-full theme-bg-accent"></div>
								</div>
							</button>

							{#if activeMenuEndpointId === endpoint.id}
								<div
									class={ThemeUtils.card(
										'absolute right-0 top-full mt-2 w-48 theme-bg-primary border z-50'
									)}
								>
									<div class="py-1">
										<button
											class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:bg-gray-700 flex items-center"
											on:click|stopPropagation={(e) =>
												handleMenuAction(e, endpoint.enabled ? 'disable' : 'enable', endpoint)}
										>
											{#if endpoint.enabled}
												<i class="fas fa-toggle-off mr-2 text-red-500"></i> Disable
											{:else}
												<i class="fas fa-toggle-on mr-2 text-green-500"></i> Enable
											{/if}
										</button>
										<button
											class="w-full text-left px-4 py-2 text-sm theme-text-primary hover:bg-gray-700 flex items-center"
											on:click|stopPropagation={(e) => handleMenuAction(e, 'duplicate', endpoint)}
										>
											<i class="fas fa-clone mr-2"></i> Duplicate
										</button>
										<hr class="theme-border my-1" />
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
	
	<!-- Resizable handle -->
	<div 
		class="absolute top-0 right-0 bottom-0 w-1 cursor-col-resize hover:bg-blue-500 transition-colors duration-200 group"
		on:mousedown={startResize}
		title="Drag to resize panel"
	>
		<div class="w-full h-full bg-transparent group-hover:bg-blue-500/30"></div>
	</div>
</div>
