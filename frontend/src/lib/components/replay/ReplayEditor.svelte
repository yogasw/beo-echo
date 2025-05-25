<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ReplayResponseFooter from './ReplayResponseFooter.svelte'; // Import the new footer component
	import * as ThemeUtils from '$lib/utils/themeUtils';

	// Tab management
	export interface Tab {
		id: string;
		name: string;
		method: string;
		url: string;
		isUnsaved: boolean;
	}

	export let tabs: Tab[] = [
		{
			id: 'tab-1',
			name: 'New Request',
			method: 'GET',
			url: '',
			isUnsaved: true
		}
	];
	export let activeTabId = 'tab-1';

	// Active tab content state
	export let activeTabContent = {
		method: 'GET',
		url: '',
		activeSection: 'params' // params, headers, body, auth, scripts, settings
	};

	const dispatch = createEventDispatcher();

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
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
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
			activeTabContent = {
				method: 'GET',
				url: '',
				activeSection: 'params'
			};
			dispatch('tabschange', { tabs, activeTabId, activeTabContent });
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
		// Update activeTabContent based on the new activeTabId
		const newActiveTab = tabs.find(t => t.id === activeTabId);
		if (newActiveTab) {
			activeTabContent = {
				method: newActiveTab.method,
				url: newActiveTab.url,
				activeSection: 'params' // Default or restore previous section
			};
		}
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function switchTab(tabId: string) {
		activeTabId = tabId;
		const tab = tabs.find(t => t.id === tabId);
		if (tab) {
			activeTabContent = {
				method: tab.method,
				url: tab.url,
				activeSection: 'params' // Default or restore previous section
			};
		}
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function setActiveSection(section: string) {
		activeTabContent.activeSection = section;
		dispatch('activeSectionChange', { activeSection: section });
	}

	// Propagate changes upwards
	$: dispatch('tabContentChange', activeTabContent);

	// Event handlers for ReplayResponseFooter
	function handleFooterToggleExpand(event: CustomEvent) {
		console.log('Footer expansion toggled:', event.detail.expanded);
		// Handle footer expansion state if needed in ReplayEditor or pass up to parent
	}

	function handleFooterShowHistory() {
		console.log('Show history clicked');
		// Handle showing history, potentially dispatching another event upwards
	}

</script>

<!-- Postman-like Request Interface -->
<div class="flex flex-col h-full">
	<!-- Header with tabs/actions -->
	<header class={ThemeUtils.themeBgSecondary('border-b theme-border')}>
		<div class="flex items-center justify-between px-4 py-2 text-sm">
			<div class="flex items-center space-x-2">
				<button
					class={ThemeUtils.themeBgAccent('flex items-center space-x-1 px-2 py-1 rounded-md theme-text-primary')}
				>
					<i class="fas fa-play text-sm"></i>
					<span>Replay</span>
				</button>
				
				<!-- Tab Navigation -->
				<div class="flex items-center space-x-1">
					{#each tabs as tab (tab.id)}
						<div class="flex items-center">
							<button
								class="flex items-center space-x-1 px-2 py-1 {activeTabId === tab.id ? ThemeUtils.themeBgAccent() : 'hover:' + ThemeUtils.themeBgAccent()} rounded-l-md theme-text-primary"
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
								class="p-1 hover:{ThemeUtils.themeBgAccent()} rounded-r-md theme-text-primary"
								on:click={() => closeTab(tab.id)}
							>
								<i class="fas fa-times text-sm"></i>
							</button>
						</div>
					{/each}
					
					<!-- Add new tab button -->
					<button 
						class="p-1 hover:{ThemeUtils.themeBgAccent()} rounded-md theme-text-primary"
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
	<main class="flex-grow p-4 space-y-4 overflow-y-auto">
		<!-- Title and actions -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-2">
				<i class="fas fa-folder-open text-orange-500 text-xl"></i>
				<h1 class="text-lg font-semibold theme-text-primary">
					{tabs.find(t => t.id === activeTabId)?.name || 'New Request'}
				</h1>
			</div>
			<div class="flex items-center space-x-2">
				<button
					class={ThemeUtils.themeBgSecondary('flex items-center space-x-1 px-3 py-1.5 hover:bg-gray-600 dark:hover:bg-gray-600 rounded-md text-sm')}
				>
					<i class="fas fa-save text-sm"></i>
					<span>Save</span>
					<i class="fas fa-chevron-down text-sm"></i>
				</button>
				
			</div>
		</div>

		<!-- Request builder -->
		<div
			class={ThemeUtils.themeBgSecondary('flex items-center border theme-border rounded-lg')}
		>
			<div class="relative">
				<select
					bind:value={activeTabContent.method}
					class={ThemeUtils.themeBgSecondary('appearance-none text-green-500 font-semibold px-4 py-2.5 rounded-l-lg focus:outline-none border-r theme-border pr-8')}
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
				class={ThemeUtils.themeBgSecondary('flex-grow p-2.5 focus:outline-none theme-text-secondary placeholder-gray-500 dark:placeholder-gray-500')}
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
		<div class="border-b theme-border">
			<nav class="flex space-x-4 text-sm">
				<button 
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'params' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
					on:click={() => setActiveSection('params')}
				>
					Params
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'auth' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
					on:click={() => setActiveSection('auth')}
				>
					Authorization
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'headers' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
					on:click={() => setActiveSection('headers')}
				>
					Headers
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'body' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
					on:click={() => setActiveSection('body')}
				>
					Body
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'scripts' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
					on:click={() => setActiveSection('scripts')}
				>
					Scripts
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'settings' ? 'border-orange-600 text-orange-600' : 'border-transparent hover:theme-text-primary'}"
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
					<p class="text-sm font-medium theme-text-secondary">Query Params</p>
					<div class="flex items-center space-x-2">
						<button class="text-sm text-blue-400 hover:underline">
							<i class="fas fa-ellipsis-h text-base align-middle"></i> Bulk Edit
						</button>
					</div>
				</div>
				<div class="overflow-x-auto">
					<table class="w-full text-sm">
						<thead>
							<tr class="text-left theme-text-muted">
								<th class="p-2 font-normal w-1/12"></th>
								<th class="p-2 font-normal w-1/3">Key</th>
								<th class="p-2 font-normal w-1/3">Value</th>
								<th class="p-2 font-normal w-1/3">Description</th>
								<th class="p-2 font-normal w-auto"></th>
							</tr>
						</thead>
						<tbody>
							<tr class="border-t theme-border">
								<td class="p-2">
									<input
										class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
										type="checkbox"
									/>
								</td>
								<td class="p-1">
									<input
										class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary"
										placeholder="Key"
										type="text"
									/>
								</td>
								<td class="p-1">
									<input
										class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary"
										placeholder="Value"
										type="text"
									/>
								</td>
								<td class="p-1">
									<input
										class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary"
										placeholder="Description"
										type="text"
									/>
								</td>
								<td class="p-2">
									<button
										class="theme-text-muted hover:text-gray-300 dark:hover:text-gray-300"
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
					<label class="block text-sm font-medium theme-text-secondary mb-2">Authorization Type</label>
					<select class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary">
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
					<p class="text-sm font-medium theme-text-secondary">Headers</p>
					<button class="text-sm text-blue-400 hover:underline">
						<i class="fas fa-ellipsis-h text-base align-middle"></i> Bulk Edit
					</button>
				</div>
				<div class="overflow-x-auto">
					<table class="w-full text-sm">
						<thead>
							<tr class="text-left theme-text-muted">
								<th class="p-2 font-normal w-1/12"></th>
								<th class="p-2 font-normal w-1/3">Key</th>
								<th class="p-2 font-normal w-1/3">Value</th>
								<th class="p-2 font-normal w-1/3">Description</th>
								<th class="p-2 font-normal w-auto"></th>
							</tr>
						</thead>
						<tbody>
							<tr class="border-t theme-border">
								<td class="p-2">
									<input class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded" type="checkbox"/>
								</td>
								<td class="p-1">
									<input class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary" placeholder="Content-Type" type="text"/>
								</td>
								<td class="p-1">
									<input class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary" placeholder="application/json" type="text"/>
								</td>
								<td class="p-1">
									<input class="w-full theme-bg-secondary p-1 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary" placeholder="Description" type="text"/>
								</td>
								<td class="p-2">
									<button class="theme-text-muted hover:text-gray-300 dark:hover:text-gray-300">
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
					class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary"
					placeholder="Enter request body..."
				></textarea>
			</div>
		{:else if activeTabContent.activeSection === 'scripts'}
			<!-- Scripts section -->
			<div class="space-y-4">
				<div class="flex space-x-4 text-sm">
					<button class="py-2 px-4 theme-bg-accent rounded-md">Pre-request Script</button>
					<button class="py-2 px-4 hover:theme-bg-accent rounded-md">Tests</button>
				</div>
				<textarea 
					class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-1 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary font-mono"
					placeholder="// Add your script here..."
				></textarea>
			</div>
		{:else if activeTabContent.activeSection === 'settings'}
			<!-- Settings section -->
			<div class="space-y-4">
				<div class="flex items-center justify-between">
					<label class="text-sm font-medium theme-text-secondary">Follow redirects</label>
					<input type="checkbox" class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded">
				</div>
				<div class="flex items-center justify-between">
					<label class="text-sm font-medium theme-text-secondary">Automatically follow redirects</label>
					<input type="checkbox" class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded">
				</div>
			</div>
		{/if}
	</main>

	<ReplayResponseFooter 
		on:toggleExpand={handleFooterToggleExpand}
		on:showHistory={handleFooterShowHistory}
	/>
</div>
