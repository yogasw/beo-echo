<script lang="ts">
	import StatusBodyTab from './routes/StatusBodyTab.svelte';
	import HeadersTab from './routes/HeadersTab.svelte';
	import RulesTab from './routes/RulesTab.svelte';
	import CallbacksTab from './routes/CallbacksTab.svelte';
	import type { MockoonResponse, MockoonRoute } from '$lib/types/Config';
	import RoutesList from '$lib/components/tabs/routes/RoutesList.svelte';
	import DropdownResponse from '$lib/components/tabs/routes/DropdownResponse.svelte';
	import type { ConfigResponse } from '$lib/api/mockoonApi';

	export let selectedConfig: ConfigResponse;
	export let routes: MockoonRoute[];
	export let activeContentTab = 'Status & Body';
	let activeConfigName = selectedConfig.name; // Store the active config name

	let selectedRoute: MockoonRoute | null = null;
	let selectedResponse: MockoonResponse | null = null;
	let filterText: string = ''; // Variable to store filter input

	$: filteredRoutes = routes.filter(route => {
		if (!filterText.trim()) return true;
		const filterParts = filterText.toLowerCase().split(' ');
		return filterParts.every(part =>
			route.endpoint.toLowerCase().includes(part) ||
			route.method.toLowerCase().includes(part)
		);
	});

	function selectRoute(route: MockoonRoute) {
		console.log('Route selected:', route);
		selectedRoute = route;
		selectedResponse = route.responses[0] || null; // Select the first response by default
	}

	function handleRouteStatusChange(route: MockoonRoute) {
		console.log('Route status changed:', route);
		const index = routes.findIndex(r => r.endpoint === route.endpoint && r.method === route.method);
		if (index !== -1) {
			routes[index] = {
				...route,
				status: route.status === 'enabled' ? 'disabled' : 'enabled'
			};
			routes = routes; // Trigger reactivity
		}
	}
</script>

<div class="flex flex-1 h-full">
	<RoutesList
		{activeConfigName}
		{selectedRoute}
		bind:filterText
		{filteredRoutes}
		{selectRoute}
		{handleRouteStatusChange}
	/>

	<!-- Details Section -->
	<div class="w-2/3 bg-gray-800 p-4 flex flex-col overflow-hidden">
		<div class="mb-4">
			<label class="block text-sm font-bold mb-2">Endpoint</label>
			<div class="flex flex-col md:flex-row items-center space-y-2 md:space-y-0 md:space-x-2">
				<select class="w-full md:w-1/6 rounded bg-gray-700 px-4 py-2 text-white"
								value="{selectedRoute?.method.toUpperCase()}">
					<option value="GET">GET</option>
					<option value="POST">POST</option>
					<option value="PUT">PUT</option>
					<option value="DELETE">DELETE</option>
					<option value="PATCH">PATCH</option>
				</select>
				<span class="text-gray-400 hidden md:block">{selectedConfig.url}/</span>
				<input type="text" class="w-full md:flex-1 rounded bg-gray-700 px-4 py-2 text-white"
							 value="{selectedRoute?.endpoint}">
				<button
					class="text-gray-400 hover:text-blue-500 disabled:text-gray-600"
					disabled={!selectedRoute || selectedRoute?.method !== 'GET'}
					on:click={() =>{
						let url = `${selectedConfig.url}/${selectedRoute?.endpoint ? selectedRoute.endpoint : ''}`;
						// Open the URL in a new tab
						window.open(url, '_blank');
					}}>
					<i class="fas fa-external-link-alt"></i>
				</button>
			</div>
			<span class="text-gray-400 block md:hidden mt-2"></span>
		</div>
		<label class="block text-sm font-bold mb-2">Documentation for this routes</label>
		<textarea
			class="w-full rounded bg-gray-700 px-4 py-2 text-white" rows="3"
			placeholder="Provide a brief description or documentation for this endpoint">{selectedRoute?.documentation || ''}</textarea>

		<DropdownResponse bind:selectedRoute bind:selectedResponse />

		<div class="flex space-x-2 mb-4">
			{#each ["Status & Body", "Headers", "Rules", "Callbacks"] as tab}
				<button
					class="text-white py-2 px-4 rounded"
					class:bg-blue-500={tab === activeContentTab}
					class:bg-gray-700={tab !== activeContentTab}
					on:click={() => activeContentTab = tab}
				>
					{tab}
				</button>
			{/each}
		</div>

		<div class="flex-1 overflow-auto h-full">
			<div class="max-w-full overflow-x-auto h-full">
				<div class="min-w-0 h-full">
					<div class="h-full flex flex-col">
						{#if selectedRoute}
							{#if activeContentTab === 'Status & Body'}
								<StatusBodyTab
									responseBody={selectedResponse?.body || ''}
									statusCode={selectedResponse?.statusCode || 200}
									onBodyChange={(val) => {
											if (selectedResponse) {
												selectedResponse.body = val;
											}
									}}
									onStatusCodeChange={(val) => {
										if (selectedResponse) {
											selectedResponse.statusCode = val;
										}
									}}
								/>
							{:else if activeContentTab === 'Headers'}
								<HeadersTab headers={selectedResponse?.headers} />
							{:else if activeContentTab === 'Rules'}
								<RulesTab rules={selectedResponse?.rules} rulesOperator={selectedResponse?.rulesOperator} />
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
