<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replays, selectedReplay, replayActions } from '$lib/stores/replay';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';

	import ReplayList from './ReplayList.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	let isLoading = true;
	let error: string | null = null;
	let activeView: 'list' | 'editor' | 'execution' | 'logs' = 'list';

	// Tab management
	interface Tab {
		id: string;
		name: string;
		method: string;
		url: string;
		isUnsaved: boolean;
	}

	let tabs: Tab[] = [
		{
			id: 'tab-1',
			name: 'New Request',
			method: 'GET',
			url: '',
			isUnsaved: true
		}
	];
	let activeTabId = 'tab-1';

	// Active tab content state
	let activeTabContent = {
		method: 'GET',
		url: '',
		activeSection: 'params' // params, headers, body, auth, scripts, settings
	};

	function createNewTab() {
		const newTabId = `tab-${Date.now()}`;
		const newTab: Tab = {
			id: newTabId,
			name: 'New Request',
			method: 'GET',
			url: '',
			isUnsaved: true
		};
		tabs = [...tabs, newTab];
		activeTabId = newTabId;
		activeTabContent = {
			method: 'GET',
			url: '',
			activeSection: 'params'
		};
	}

	function closeTab(tabId: string) {
		if (tabs.length === 1) {
			// Don't close the last tab, just reset it
			tabs[0] = {
				id: 'tab-1',
				name: 'New Request',
				method: 'GET',
				url: '',
				isUnsaved: true
			};
			activeTabId = 'tab-1';
			return;
		}

		const tabIndex = tabs.findIndex(tab => tab.id === tabId);
		tabs = tabs.filter(tab => tab.id !== tabId);
		
		// Switch to adjacent tab
		if (activeTabId === tabId) {
			if (tabIndex > 0) {
				activeTabId = tabs[tabIndex - 1].id;
			} else {
				activeTabId = tabs[0].id;
			}
		}
	}

	function switchTab(tabId: string) {
		activeTabId = tabId;
		const tab = tabs.find(t => t.id === tabId);
		if (tab) {
			activeTabContent = {
				method: tab.method,
				url: tab.url,
				activeSection: 'params'
			};
		}
	}

	function setActiveSection(section: string) {
		activeTabContent.activeSection = section;
	}

	// Load replays when component mounts or project changes
	$: if ($selectedWorkspace && $selectedProject) {
		loadReplays();
	}

	// Clear replay data when project changes
	$: if ($selectedProject) {
		replayActions.clearAll();
		activeView = 'list';
	}

	async function loadReplays() {
		if (!$selectedWorkspace || !$selectedProject) return;

		try {
			isLoading = true;
			error = null;
			replayActions.setLoading('list', true);

			const response = await replayApi.listReplays($selectedWorkspace.id, $selectedProject.id);
			replays.set(response.replays);
		} catch (err: any) {
			error = err.message || 'Failed to load replays';
			toast.error(err);
		} finally {
			isLoading = false;
			replayActions.setLoading('list', false);
		}
	}

	function handleCreateNew() {
		selectedReplay.set(null);
		activeView = 'editor';
	}

	function handleEditReplay(replay: any) {
		selectedReplay.set(replay.detail);
		activeView = 'editor';
	}

	function handleExecuteReplay(replay: any) {
		selectedReplay.set(replay.detail);
		activeView = 'execution';
	}

	function handleViewLogs(replay: any) {
		if (replay?.detail) {
			selectedReplay.set(replay.detail);
		}
		activeView = 'logs';
	}

	function handleBackToList() {
		selectedReplay.set(null);
		activeView = 'list';
	}

	function handleReplayCreated() {
		activeView = 'list';
		loadReplays(); // Refresh the list
	}

	function handleReplayUpdated() {
		activeView = 'list';
		loadReplays(); // Refresh the list
	}

	onMount(() => {
		if ($selectedWorkspace && $selectedProject) {
			loadReplays();
		}
	});

	onDestroy(() => {
		replayActions.clearAll();
	});
</script>

<div class="flex flex-col h-full theme-bg-primary">
	<!-- Main Content Area -->
	<div class="flex-1 grid grid-cols-3 gap-4 p-4 h-full">
		<!-- Left: Replay List -->
		<div class="flex flex-col h-full space-y-4">
			<button
				on:click={handleCreateNew}
				class="w-full px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors flex-shrink-0"
			>
				<i class="fas fa-plus mr-1"></i>
				New Replay
			</button>

			<div class="flex-1 min-h-0">
				{#if !$selectedWorkspace || !$selectedProject}
					<div class="flex items-center justify-center h-full">
						<div class="text-center theme-text-secondary">
							<i class="fas fa-project-diagram text-4xl mb-4 opacity-50"></i>
							<p>Please select a workspace and project to manage replays</p>
						</div>
					</div>
				{:else if isLoading}
					<SkeletonLoader type="list" count={5} />
				{:else if error}
					<ErrorDisplay message={error} type="error" retryable={true} onRetry={loadReplays} />
				{:else}
					<ReplayList
						on:edit={handleEditReplay}
						on:execute={handleExecuteReplay}
						on:logs={handleViewLogs}
						on:refresh={loadReplays}
					/>
				{/if}
			</div>
		</div>
		<!-- Center & Right: Replay Content (Postman-like Interface) -->
		<div
			class="col-span-2 flex flex-col h-full bg-gray-800 dark:bg-gray-800 border border-gray-700 dark:border-gray-700 rounded-lg overflow-hidden"
		>
			<!-- Postman-like Request Interface -->
			<div class="flex flex-col h-full">
				<!-- Header with tabs/actions -->
				<header class="bg-gray-900 dark:bg-gray-900 border-b border-gray-700 dark:border-gray-700">
					<div class="flex items-center justify-between px-4 py-2 text-sm">
						<div class="flex items-center space-x-2">
							<button
								class="flex items-center space-x-1 px-2 py-1 hover:bg-gray-700 dark:hover:bg-gray-700 rounded-md"
							>
								<i class="fas fa-play text-sm"></i>
								<span>Replay</span>
							</button>
							
							<!-- Tab Navigation -->
							<div class="flex items-center space-x-1">
								{#each tabs as tab (tab.id)}
									<div class="flex items-center">
										<button
											class="flex items-center space-x-1 px-2 py-1 {activeTabId === tab.id ? 'bg-gray-700' : 'hover:bg-gray-700'} dark:hover:bg-gray-700 rounded-l-md"
											on:click={() => switchTab(tab.id)}
										>
											<span class="text-{tab.method === 'GET' ? 'green' : tab.method === 'POST' ? 'blue' : tab.method === 'PUT' ? 'yellow' : tab.method === 'DELETE' ? 'red' : 'gray'}-500 font-semibold">{tab.method}</span>
											<span class="max-w-32 truncate">{tab.name}</span>
											{#if tab.isUnsaved}
												<span class="w-2 h-2 bg-orange-500 rounded-full"></span>
											{/if}
										</button>
										<button 
											aria-label="Close tab"
											class="p-1 hover:bg-gray-700 dark:hover:bg-gray-700 rounded-r-md"
											on:click={() => closeTab(tab.id)}
										>
											<i class="fas fa-times text-sm"></i>
										</button>
									</div>
								{/each}
								
								<!-- Add new tab button -->
								<button 
									class="p-1 hover:bg-gray-700 dark:hover:bg-gray-700 rounded-md"
									aria-label="Add new tab"
									on:click={createNewTab}
								>
									<i class="fas fa-plus text-sm"></i>
								</button>
							</div>
						</div>
						
					</div>
				</header>

				<!-- Main content -->
				<main class="flex-grow p-4 space-y-4">
					<!-- Title and actions -->
					<div class="flex items-center justify-between">
						<div class="flex items-center space-x-2">
							<i class="fas fa-folder-open text-orange-500 text-xl"></i>
							<h1 class="text-lg font-semibold text-white dark:text-white">
								{tabs.find(t => t.id === activeTabId)?.name || 'New Request'}
							</h1>
						</div>
						<div class="flex items-center space-x-2">
							<button
								class="flex items-center space-x-1 px-3 py-1.5 bg-gray-700 dark:bg-gray-700 hover:bg-gray-600 dark:hover:bg-gray-600 rounded-md text-sm"
							>
								<i class="fas fa-save text-sm"></i>
								<span>Save</span>
								<i class="fas fa-chevron-down text-sm"></i>
							</button>
							
						</div>
					</div>

					<!-- Request builder -->
					<div
						class="flex items-center bg-gray-900 dark:bg-gray-900 border border-gray-700 dark:border-gray-700 rounded-lg"
					>
						<div class="relative">
							<select
								bind:value={activeTabContent.method}
								class="appearance-none bg-gray-900 dark:bg-gray-900 text-green-500 font-semibold px-4 py-2.5 rounded-l-lg focus:outline-none border-r border-gray-700 dark:border-gray-700 pr-8"
							>
								<option>GET</option>
								<option>POST</option>
								<option>PUT</option>
								<option>DELETE</option>
								<option>PATCH</option>
							</select>
							<i
								class="fas fa-chevron-down absolute right-2 top-1/2 transform -translate-y-1/2 text-gray-400 pointer-events-none"
							></i>
						</div>
						<input
							bind:value={activeTabContent.url}
							class="flex-grow bg-gray-900 dark:bg-gray-900 p-2.5 focus:outline-none text-gray-300 dark:text-gray-300 placeholder-gray-500 dark:placeholder-gray-500"
							placeholder="Enter URL or describe the request"
							type="text"
						/>
						<button
							class="bg-blue-600 hover:bg-blue-500 text-white px-4 py-2.5 rounded-r-lg flex items-center space-x-1"
						>
							<span>Send</span>
							<i class="fas fa-paper-plane text-sm"></i>
						</button>
					</div>

					<!-- Tab navigation -->
					<div class="border-b border-gray-700 dark:border-gray-700">
						<nav class="flex space-x-4 text-sm">
							<button 
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'params' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('params')}
							>
								Params
							</button>
							<button
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'auth' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('auth')}
							>
								Authorization
							</button>
							<button
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'headers' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('headers')}
							>
								Headers
							</button>
							<button
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'body' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('body')}
							>
								Body
							</button>
							<button
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'scripts' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('scripts')}
							>
								Scripts
							</button>
							<button
								class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'settings' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:text-gray-100 dark:hover:text-gray-100'}"
								on:click={() => setActiveSection('settings')}
							>
								Settings
							</button>
						</nav>
					</div>

					<!-- Dynamic content based on active section -->
					{#if activeTabContent.activeSection === 'params'}
						<!-- Parameters section -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-sm font-medium text-gray-300 dark:text-gray-300">Query Params</p>
								<div class="flex items-center space-x-2">
									<button class="text-sm text-blue-400 hover:underline">
										<i class="fas fa-ellipsis-h text-base align-middle"></i> Bulk Edit
									</button>
								</div>
							</div>
							<div class="overflow-x-auto">
								<table class="w-full text-sm">
									<thead>
										<tr class="text-left text-gray-400 dark:text-gray-400">
											<th class="p-2 font-normal w-1/12"></th>
											<th class="p-2 font-normal w-1/3">Key</th>
											<th class="p-2 font-normal w-1/3">Value</th>
											<th class="p-2 font-normal w-1/3">Description</th>
											<th class="p-2 font-normal w-auto"></th>
										</tr>
									</thead>
									<tbody>
										<tr class="border-t border-gray-700 dark:border-gray-700">
											<td class="p-2">
												<input
													class="form-checkbox h-4 w-4 bg-gray-800 dark:bg-gray-800 border-gray-600 dark:border-gray-600 text-blue-500 focus:ring-blue-500 rounded"
													type="checkbox"
												/>
											</td>
											<td class="p-1">
												<input
													class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300"
													placeholder="Key"
													type="text"
												/>
											</td>
											<td class="p-1">
												<input
													class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300"
													placeholder="Value"
													type="text"
												/>
											</td>
											<td class="p-1">
												<input
													class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300"
													placeholder="Description"
													type="text"
												/>
											</td>
											<td class="p-2">
												<button
													class="text-gray-500 dark:text-gray-500 hover:text-gray-300 dark:hover:text-gray-300"
												>
													<i class="fas fa-trash text-base"></i>
												</button>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
					{:else if activeTabContent.activeSection === 'auth'}
						<!-- Authorization section -->
						<div class="space-y-4">
							<div>
								<label class="block text-sm font-medium text-gray-300 dark:text-gray-300 mb-2">Authorization Type</label>
								<select class="w-full bg-gray-800 dark:bg-gray-800 p-2 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-gray-700 dark:border-gray-700 text-gray-300 dark:text-gray-300">
									<option>No Auth</option>
									<option>Bearer Token</option>
									<option>Basic Auth</option>
									<option>API Key</option>
								</select>
							</div>
						</div>
					{:else if activeTabContent.activeSection === 'headers'}
						<!-- Headers section -->
						<div>
							<div class="flex justify-between items-center mb-2">
								<p class="text-sm font-medium text-gray-300 dark:text-gray-300">Headers</p>
								<button class="text-sm text-blue-400 hover:underline">
									<i class="fas fa-ellipsis-h text-base align-middle"></i> Bulk Edit
								</button>
							</div>
							<div class="overflow-x-auto">
								<table class="w-full text-sm">
									<thead>
										<tr class="text-left text-gray-400 dark:text-gray-400">
											<th class="p-2 font-normal w-1/12"></th>
											<th class="p-2 font-normal w-1/3">Key</th>
											<th class="p-2 font-normal w-1/3">Value</th>
											<th class="p-2 font-normal w-1/3">Description</th>
											<th class="p-2 font-normal w-auto"></th>
										</tr>
									</thead>
									<tbody>
										<tr class="border-t border-gray-700 dark:border-gray-700">
											<td class="p-2">
												<input class="form-checkbox h-4 w-4 bg-gray-800 dark:bg-gray-800 border-gray-600 dark:border-gray-600 text-blue-500 focus:ring-blue-500 rounded" type="checkbox"/>
											</td>
											<td class="p-1">
												<input class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300" placeholder="Content-Type" type="text"/>
											</td>
											<td class="p-1">
												<input class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300" placeholder="application/json" type="text"/>
											</td>
											<td class="p-1">
												<input class="w-full bg-gray-800 dark:bg-gray-800 p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 text-gray-300 dark:text-gray-300" placeholder="Description" type="text"/>
											</td>
											<td class="p-2">
												<button class="text-gray-500 dark:text-gray-500 hover:text-gray-300 dark:hover:text-gray-300">
													<i class="fas fa-trash text-base"></i>
												</button>
											</td>
										</tr>
									</tbody>
								</table>
							</div>
						</div>
					{:else if activeTabContent.activeSection === 'body'}
						<!-- Body section -->
						<div class="space-y-4">
							<div class="flex space-x-4 text-sm">
								<label class="flex items-center">
									<input type="radio" name="bodyType" value="none" class="mr-2" checked>
									<span>None</span>
								</label>
								<label class="flex items-center">
									<input type="radio" name="bodyType" value="raw" class="mr-2">
									<span>Raw</span>
								</label>
								<label class="flex items-center">
									<input type="radio" name="bodyType" value="form" class="mr-2">
									<span>Form Data</span>
								</label>
							</div>
							<textarea 
								class="w-full h-64 bg-gray-800 dark:bg-gray-800 p-3 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-gray-700 dark:border-gray-700 text-gray-300 dark:text-gray-300"
								placeholder="Enter request body..."
							></textarea>
						</div>
					{:else if activeTabContent.activeSection === 'scripts'}
						<!-- Scripts section -->
						<div class="space-y-4">
							<div class="flex space-x-4 text-sm">
								<button class="py-2 px-4 bg-gray-700 dark:bg-gray-700 rounded-md">Pre-request Script</button>
								<button class="py-2 px-4 hover:bg-gray-700 dark:hover:bg-gray-700 rounded-md">Tests</button>
							</div>
							<textarea 
								class="w-full h-64 bg-gray-800 dark:bg-gray-800 p-3 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-gray-700 dark:border-gray-700 text-gray-300 dark:text-gray-300 font-mono"
								placeholder="// Add your script here..."
							></textarea>
						</div>
					{:else if activeTabContent.activeSection === 'settings'}
						<!-- Settings section -->
						<div class="space-y-4">
							<div class="flex items-center justify-between">
								<label class="text-sm font-medium text-gray-300 dark:text-gray-300">Follow redirects</label>
								<input type="checkbox" class="form-checkbox h-4 w-4 bg-gray-800 dark:bg-gray-800 border-gray-600 dark:border-gray-600 text-blue-500 focus:ring-blue-500 rounded">
							</div>
							<div class="flex items-center justify-between">
								<label class="text-sm font-medium text-gray-300 dark:text-gray-300">Automatically follow redirects</label>
								<input type="checkbox" class="form-checkbox h-4 w-4 bg-gray-800 dark:bg-gray-800 border-gray-600 dark:border-gray-600 text-blue-500 focus:ring-blue-500 rounded">
							</div>
						</div>
					{/if}
				</main>

				<!-- Response footer -->
				<footer class="bg-gray-900 dark:bg-gray-900 border-t border-gray-700 dark:border-gray-700">
					<div
						class="flex items-center justify-between p-2 border-b border-gray-700 dark:border-gray-700"
					>
						<div class="flex items-center space-x-2">
							<span class="text-sm font-medium text-gray-300 dark:text-gray-300">Response</span>
							<button
								class="flex items-center space-x-1 text-sm text-gray-400 dark:text-gray-400 hover:text-gray-200 dark:hover:text-gray-200"
							>
								<i class="fas fa-history text-base"></i>
								<span>History</span>
							</button>
						</div>
						<button
							class="text-gray-400 dark:text-gray-400 hover:text-gray-200 dark:hover:text-gray-200"
						>
							<i class="fas fa-chevron-up"></i>
						</button>
					</div>
					<div class="flex flex-col items-center justify-center h-64 text-center p-4">
						<div
							class="h-32 w-32 mb-4 opacity-75 bg-gray-700 dark:bg-gray-700 rounded-lg flex items-center justify-center"
						>
							<i class="fas fa-rocket text-4xl text-gray-500 dark:text-gray-500"></i>
						</div>
						<p class="text-gray-400 dark:text-gray-400">
							Enter the URL and click Send to get a response
						</p>
					</div>
				</footer>
			</div>
		</div>
	</div>
</div>
