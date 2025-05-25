<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ReplayResponseFooter from './ReplayResponseFooter.svelte'; // Import the new footer component
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import ReplayBar from './ReplayBar.svelte';

	// Tab management
	interface Tab {
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

		const tabIndex = tabs.findIndex((tab) => tab.id === tabId);
		tabs = tabs.filter((tab) => tab.id !== tabId);

		// Switch to adjacent tab
		if (activeTabId === tabId) {
			if (tabIndex > 0) {
				activeTabId = tabs[tabIndex - 1].id;
			} else {
				activeTabId = tabs[0].id;
			}
		}
		// Update activeTabContent based on the new activeTabId
		const newActiveTab = tabs.find((t) => t.id === activeTabId);
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
		const tab = tabs.find((t) => t.id === tabId);
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
	<ReplayBar 
		activeTabId={activeTabId}
		switchTab={switchTab}
		closeTab={closeTab}
		createNewTab={createNewTab}
		tabs={tabs}
	/>
	<!-- Main content -->
	<main class="flex-grow p-4 space-y-4 overflow-y-auto">
		<!-- Title and actions -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-2">
				<i class="fas fa-folder-open text-orange-500 text-xl"></i>
				<h1 class="text-lg font-semibold theme-text-primary">
					{tabs.find((t) => t.id === activeTabId)?.name || 'New Request'}
				</h1>
			</div>
			<div class="flex items-center space-x-2">
				<button
					class={ThemeUtils.secondaryButton('flex items-center space-x-2 px-3 py-1.5 rounded-md text-sm border theme-border hover:shadow-md transition-all duration-200')}
					title="Save current request configuration"
					aria-label="Save request"
				>
					<i class="fas fa-save text-sm"></i>
					<span>Save</span>
					<i class="fas fa-chevron-down text-sm ml-1"></i>
				</button>
			</div>
		</div>

		<!-- Request builder -->
		<div class={ThemeUtils.themeBgSecondary('flex items-center border theme-border rounded-lg shadow-sm')}>
			<div class="relative">
				<select
					bind:value={activeTabContent.method}
					class={ThemeUtils.themeBgSecondary(
						'appearance-none font-semibold px-4 py-2.5 rounded-l-lg focus:outline-none focus:ring-2 focus:ring-blue-500 border-r theme-border pr-8 theme-text-primary transition-all duration-200'
					)}
					title="Select HTTP method"
					aria-label="HTTP method selector"
				>
					<option>GET</option>
					<option>POST</option>
					<option>PUT</option>
					<option>DELETE</option>
					<option>PATCH</option>
				</select>
				<i
					class="fas fa-chevron-down absolute right-2 top-1/2 transform -translate-y-1/2 theme-text-muted pointer-events-none"
				></i>
			</div>
			<input
				bind:value={activeTabContent.url}
				class={ThemeUtils.themeBgSecondary(
					'flex-grow p-2.5 focus:outline-none focus:ring-2 focus:ring-blue-500 theme-text-secondary placeholder-gray-400 dark:placeholder-gray-500 transition-all duration-200'
				)}
				placeholder="Enter URL or describe the request"
				type="text"
				title="Request URL input"
				aria-label="Enter request URL"
			/>
			<button 
				class={ThemeUtils.primaryButton('px-4 py-2.5 rounded-r-lg space-x-1 shadow-sm hover:shadow-md transition-all duration-200')}
				title="Send HTTP request"
				aria-label="Send request"
			>
				<span>Send</span>
				<i class="fas fa-paper-plane text-sm"></i>
			</button>
		</div>

		<!-- Tab navigation -->
		<div class="border-b theme-border">
			<div class="flex space-x-4 text-sm" role="tablist" aria-label="Request configuration tabs">
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'params'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure query parameters"
					aria-label="Parameters tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'params'}
					on:click={() => setActiveSection('params')}
				>
					Params
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'auth'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request authorization"
					aria-label="Authorization tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'auth'}
					on:click={() => setActiveSection('auth')}
				>
					Authorization
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'headers'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request headers"
					aria-label="Headers tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'headers'}
					on:click={() => setActiveSection('headers')}
				>
					Headers
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'body'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request body"
					aria-label="Body tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'body'}
					on:click={() => setActiveSection('body')}
				>
					Body
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'scripts'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure pre-request scripts and tests"
					aria-label="Scripts tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'scripts'}
					on:click={() => setActiveSection('scripts')}
				>
					Scripts
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTabContent.activeSection === 'settings'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request settings"
					aria-label="Settings tab"
					role="tab"
					aria-selected={activeTabContent.activeSection === 'settings'}
					on:click={() => setActiveSection('settings')}
				>
					Settings
				</button>
			</div>
		</div>

		<!-- Dynamic content based on active section -->
		{#if activeTabContent.activeSection === 'params'}
			<!-- Parameters section -->
			<div role="tabpanel" aria-labelledby="params-tab">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-list-ul text-orange-500 mr-2"></i>
						Query Parameters
					</h2>
					<div class="flex items-center space-x-2">
						<button 
							class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
							title="Open bulk edit mode for parameters"
							aria-label="Bulk edit parameters"
						>
							<i class="fas fa-edit text-xs mr-1"></i>
							Bulk Edit
						</button>
					</div>
				</div>
				<div class="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg border theme-border">
					<table class="w-full text-sm">
						<thead class="bg-gray-50 dark:bg-gray-750">
							<tr class="text-left theme-text-muted">
								<th class="p-3 font-medium w-12">
									<input
										class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
										type="checkbox"
										title="Select all parameters"
										aria-label="Select all parameters"
									/>
								</th>
								<th class="p-3 font-medium w-1/3">Key</th>
								<th class="p-3 font-medium w-1/3">Value</th>
								<th class="p-3 font-medium w-1/3">Description</th>
								<th class="p-3 font-medium w-12">Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr class="border-t theme-border hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors duration-150">
								<td class="p-3">
									<input
										class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
										type="checkbox"
										title="Include this parameter"
										aria-label="Enable parameter"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="Parameter key"
										type="text"
										title="Parameter key name"
										aria-label="Parameter key"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="Parameter value"
										type="text"
										title="Parameter value"
										aria-label="Parameter value"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="Optional description"
										type="text"
										title="Parameter description"
										aria-label="Parameter description"
									/>
								</td>
								<td class="p-3 text-center">
									<button
										class="theme-text-muted hover:text-red-500 p-1 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-all duration-200"
										title="Delete this parameter"
										aria-label="Delete parameter"
									>
										<i class="fas fa-trash text-sm"></i>
									</button>
								</td>
							</tr>
						</tbody>
					</table>
					<div class="p-3 border-t theme-border bg-gray-50 dark:bg-gray-750">
						<button 
							class="text-sm text-blue-500 hover:text-blue-400 font-medium flex items-center transition-colors duration-200"
							title="Add new parameter"
							aria-label="Add new parameter row"
						>
							<i class="fas fa-plus text-xs mr-1"></i>
							Add Parameter
						</button>
					</div>
				</div>
			</div>
		{:else if activeTabContent.activeSection === 'auth'}
			<!-- Authorization section -->
			<div role="tabpanel" aria-labelledby="auth-tab" class="space-y-4">
				<div class="flex items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-lock text-green-500 mr-2"></i>
						Authorization
					</h2>
				</div>
				<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4">
					<div class="space-y-4">
						<div>
							<label for="auth-type" class="block text-sm font-medium theme-text-secondary mb-2">
								Authorization Type
							</label>
							<select
								id="auth-type"
								class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary transition-all duration-200"
								title="Select authorization method"
								aria-label="Choose authorization type"
							>
								<option>No Auth</option>
								<option>Bearer Token</option>
								<option>Basic Auth</option>
								<option>API Key</option>
								<option>OAuth 2.0</option>
							</select>
						</div>
						
						<!-- Token input placeholder (shown when Bearer Token is selected) -->
						<div class="space-y-3">
							<div>
								<label for="bearer-token" class="block text-sm font-medium theme-text-secondary mb-2">
									Token
								</label>
								<input
									id="bearer-token"
									type="password"
									class="w-full theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
									placeholder="Enter your bearer token"
									title="Bearer token for authorization"
									aria-label="Bearer token input"
								/>
							</div>
							<p class="text-xs theme-text-muted">
								<i class="fas fa-info-circle mr-1"></i>
								This token will be included in the Authorization header as "Bearer &#123;token&#125;"
							</p>
						</div>
					</div>
				</div>
			</div>
		{:else if activeTabContent.activeSection === 'headers'}
			<!-- Headers section -->
			<div role="tabpanel" aria-labelledby="headers-tab">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-tags text-blue-500 mr-2"></i>
						Headers
					</h2>
					<button 
						class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
						title="Open bulk edit mode for headers"
						aria-label="Bulk edit headers"
					>
						<i class="fas fa-edit text-xs mr-1"></i>
						Bulk Edit
					</button>
				</div>
				<div class="overflow-x-auto bg-white dark:bg-gray-800 rounded-lg border theme-border">
					<table class="w-full text-sm">
						<thead class="bg-gray-50 dark:bg-gray-750">
							<tr class="text-left theme-text-muted">
								<th class="p-3 font-medium w-12">
									<input
										class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
										type="checkbox"
										title="Select all headers"
										aria-label="Select all headers"
									/>
								</th>
								<th class="p-3 font-medium w-1/3">Key</th>
								<th class="p-3 font-medium w-1/3">Value</th>
								<th class="p-3 font-medium w-1/3">Description</th>
								<th class="p-3 font-medium w-12">Actions</th>
							</tr>
						</thead>
						<tbody>
							<tr class="border-t theme-border hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors duration-150">
								<td class="p-3">
									<input
										class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded"
										type="checkbox"
										title="Include this header"
										aria-label="Enable header"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="Content-Type"
										type="text"
										title="Header key name"
										aria-label="Header key"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="application/json"
										type="text"
										title="Header value"
										aria-label="Header value"
									/>
								</td>
								<td class="p-2">
									<input
										class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border border-transparent focus:border-blue-500 theme-text-secondary transition-all duration-200"
										placeholder="Optional description"
										type="text"
										title="Header description"
										aria-label="Header description"
									/>
								</td>
								<td class="p-3 text-center">
									<button
										class="theme-text-muted hover:text-red-500 p-1 rounded hover:bg-red-50 dark:hover:bg-red-900/20 transition-all duration-200"
										title="Delete this header"
										aria-label="Delete header"
									>
										<i class="fas fa-trash text-sm"></i>
									</button>
								</td>
							</tr>
						</tbody>
					</table>
					<div class="p-3 border-t theme-border bg-gray-50 dark:bg-gray-750">
						<button 
							class="text-sm text-blue-500 hover:text-blue-400 font-medium flex items-center transition-colors duration-200"
							title="Add new header"
							aria-label="Add new header row"
						>
							<i class="fas fa-plus text-xs mr-1"></i>
							Add Header
						</button>
					</div>
				</div>
			</div>
		{:else if activeTabContent.activeSection === 'body'}
			<!-- Body section -->
			<div role="tabpanel" aria-labelledby="body-tab" class="space-y-4">
				<div class="flex items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-file-alt text-purple-500 mr-2"></i>
						Request Body
					</h2>
				</div>
				<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4">
					<div class="space-y-4">
						<fieldset>
							<legend class="sr-only">Body type selection</legend>
							<div class="flex flex-wrap gap-4 text-sm">
								<label class="flex items-center cursor-pointer">
									<input 
										type="radio" 
										name="bodyType" 
										value="none" 
										class="mr-2 text-blue-500 focus:ring-blue-500" 
										checked 
										title="No request body"
										aria-label="No body"
									/>
									<span class="theme-text-secondary">None</span>
								</label>
								<label class="flex items-center cursor-pointer">
									<input 
										type="radio" 
										name="bodyType" 
										value="raw" 
										class="mr-2 text-blue-500 focus:ring-blue-500"
										title="Raw text or JSON body"
										aria-label="Raw body"
									/>
									<span class="theme-text-secondary">Raw</span>
								</label>
								<label class="flex items-center cursor-pointer">
									<input 
										type="radio" 
										name="bodyType" 
										value="form" 
										class="mr-2 text-blue-500 focus:ring-blue-500"
										title="Form data body"
										aria-label="Form data body"
									/>
									<span class="theme-text-secondary">Form Data</span>
								</label>
								<label class="flex items-center cursor-pointer">
									<input 
										type="radio" 
										name="bodyType" 
										value="urlencoded" 
										class="mr-2 text-blue-500 focus:ring-blue-500"
										title="URL encoded form body"
										aria-label="URL encoded body"
									/>
									<span class="theme-text-secondary">x-www-form-urlencoded</span>
								</label>
							</div>
						</fieldset>
						
						<div>
							<label for="request-body" class="block text-sm font-medium theme-text-secondary mb-2">
								Body Content
							</label>
							<textarea
								id="request-body"
								class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary font-mono text-sm transition-all duration-200"
								placeholder="Enter request body content..."
								title="Request body content"
								aria-label="Request body textarea"
							></textarea>
							<div class="flex justify-between items-center mt-2">
								<p class="text-xs theme-text-muted">
									<i class="fas fa-info-circle mr-1"></i>
									Supports JSON, XML, text, and other formats
								</p>
								<button 
									class="text-xs text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200"
									title="Format and beautify JSON content"
									aria-label="Beautify JSON"
								>
									<i class="fas fa-magic mr-1"></i>
									Beautify
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		{:else if activeTabContent.activeSection === 'scripts'}
			<!-- Scripts section -->
			<div role="tabpanel" aria-labelledby="scripts-tab" class="space-y-4">
				<div class="flex justify-between items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-code text-purple-500 mr-2"></i>
						Scripts
					</h2>
					<div class="flex items-center space-x-2">
						<button 
							class="text-sm text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
							title="View script documentation and examples"
							aria-label="View script documentation"
						>
							<i class="fas fa-question-circle text-xs mr-1"></i>
							Help
						</button>
					</div>
				</div>
				<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
					<div class="border-b theme-border">
						<nav class="flex space-x-0" role="tablist" aria-label="Script type navigation">
							<button 
								class="flex-1 py-3 px-4 text-sm font-medium transition-all duration-200 border-r theme-border bg-blue-600 text-white shadow-sm hover:bg-blue-700"
								title="Edit pre-request script that runs before sending the request"
								aria-label="Pre-request Script tab"
								role="tab"
								aria-selected="true"
							>
								<i class="fas fa-play-circle mr-2"></i>
								Pre-request Script
							</button>
							<button 
								class="flex-1 py-3 px-4 text-sm font-medium transition-all duration-200 theme-bg-secondary theme-text-secondary hover:bg-gray-100 dark:hover:bg-gray-600 hover:shadow-sm"
								title="Edit test scripts that run after receiving the response"
								aria-label="Tests tab"
								role="tab"
								aria-selected="false"
							>
								<i class="fas fa-check-circle mr-2"></i>
								Tests
							</button>
						</nav>
					</div>
					<div class="p-4 space-y-4">
						<div>
							<label for="script-editor" class="flex items-center text-sm font-medium theme-text-secondary mb-2">
								<i class="fas fa-edit text-purple-400 mr-2"></i>
								JavaScript Code Editor
							</label>
							<div class="relative">
								<textarea
									id="script-editor"
									class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary font-mono text-sm transition-all duration-200 resize-y"
									placeholder={`// Pre-request Script Example:
// Set environment variables
pm.environment.set('timestamp', Date.now());

// Add dynamic headers
pm.request.headers.add({
    key: 'X-Request-ID',
    value: Math.random().toString(36).substr(2, 9)
});

console.log('Pre-request script executed');`}
									title="JavaScript code editor for request scripts"
									aria-label="Script editor textarea"
								></textarea>
								<div class="absolute top-2 right-2 flex space-x-1">
									<button 
										class="p-1 text-xs theme-text-muted hover:theme-text-primary bg-white dark:bg-gray-700 rounded shadow-sm transition-colors duration-200"
										title="Format and beautify JavaScript code"
										aria-label="Format code"
									>
										<i class="fas fa-indent"></i>
									</button>
									<button 
										class="p-1 text-xs theme-text-muted hover:theme-text-primary bg-white dark:bg-gray-700 rounded shadow-sm transition-colors duration-200"
										title="Clear editor content"
										aria-label="Clear editor"
									>
										<i class="fas fa-eraser"></i>
									</button>
								</div>
							</div>
							<div class="flex justify-between items-start mt-2">
								<div class="space-y-1">
									<p class="text-xs theme-text-muted">
										<i class="fas fa-info-circle mr-1"></i>
										Write JavaScript to execute before sending the request
									</p>
									<p class="text-xs theme-text-muted">
										Available objects: <code class="bg-gray-100 dark:bg-gray-700 px-1 rounded text-blue-600 dark:text-blue-400">pm.request</code>, 
										<code class="bg-gray-100 dark:bg-gray-700 px-1 rounded text-blue-600 dark:text-blue-400">pm.environment</code>, 
										<code class="bg-gray-100 dark:bg-gray-700 px-1 rounded text-blue-600 dark:text-blue-400">pm.globals</code>
									</p>
								</div>
								<button 
									class="text-xs text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200 flex items-center"
									title="View available script APIs and examples"
									aria-label="View script API reference"
								>
									<i class="fas fa-book mr-1"></i>
									API Reference
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
		{:else if activeTabContent.activeSection === 'settings'}
			<!-- Settings section -->
			<div role="tabpanel" aria-labelledby="settings-tab" class="space-y-4">
				<div class="flex items-center mb-4">
					<h2 class="text-sm font-semibold theme-text-primary flex items-center">
						<i class="fas fa-cog text-gray-500 mr-2"></i>
						Request Settings
					</h2>
				</div>
				<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
					<div class="divide-y theme-border">
						<!-- Redirect Settings -->
						<div class="p-4 space-y-3">
							<h3 class="text-sm font-medium theme-text-primary flex items-center mb-3">
								<i class="fas fa-external-link-alt text-blue-500 mr-2"></i>
								Redirect Behavior
							</h3>
							<div class="space-y-3">
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="follow-redirects" class="text-sm font-medium theme-text-primary cursor-pointer">
											Follow redirects
										</label>
										<span class="text-xs theme-text-muted">Automatically follow HTTP redirect responses (3xx)</span>
									</div>
									<div class="flex items-center">
										<input
											id="follow-redirects"
											type="checkbox"
											class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
											title="Automatically follow HTTP redirects"
											aria-label="Enable following redirects"
											checked
										/>
									</div>
								</div>
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="max-redirects" class="text-sm font-medium theme-text-primary">
											Maximum redirects
										</label>
										<span class="text-xs theme-text-muted">Limit the number of redirects to follow</span>
									</div>
									<div class="flex items-center">
										<input
											id="max-redirects"
											type="number"
											min="0"
											max="20"
											value="5"
											class="w-16 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
											title="Maximum number of redirects to follow"
											aria-label="Maximum redirects"
										/>
									</div>
								</div>
							</div>
						</div>

						<!-- Timeout Settings -->
						<div class="p-4 space-y-3">
							<h3 class="text-sm font-medium theme-text-primary flex items-center mb-3">
								<i class="fas fa-clock text-yellow-500 mr-2"></i>
								Timeout Configuration
							</h3>
							<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
								<div class="flex flex-col">
									<label for="request-timeout" class="text-sm font-medium theme-text-primary">
										Request timeout (ms)
									</label>
									<span class="text-xs theme-text-muted">Maximum time to wait for a response</span>
								</div>
								<div class="flex items-center">
									<input
										id="request-timeout"
										type="number"
										min="100"
										max="300000"
										value="30000"
										step="1000"
										class="w-24 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
										title="Request timeout in milliseconds"
										aria-label="Request timeout"
									/>
								</div>
							</div>
						</div>

						<!-- SSL/TLS Settings -->
						<div class="p-4 space-y-3">
							<h3 class="text-sm font-medium theme-text-primary flex items-center mb-3">
								<i class="fas fa-shield-alt text-green-500 mr-2"></i>
								SSL/TLS Options
							</h3>
							<div class="space-y-3">
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="verify-ssl" class="text-sm font-medium theme-text-primary cursor-pointer">
											Verify SSL certificates
										</label>
										<span class="text-xs theme-text-muted">Validate SSL certificates for HTTPS requests</span>
									</div>
									<div class="flex items-center">
										<input
											id="verify-ssl"
											type="checkbox"
											class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
											title="Verify SSL certificates"
											aria-label="Enable SSL certificate verification"
											checked
										/>
									</div>
								</div>
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="ignore-ssl-errors" class="text-sm font-medium theme-text-primary cursor-pointer">
											Ignore SSL errors
										</label>
										<span class="text-xs theme-text-muted">Continue with requests even if SSL verification fails</span>
									</div>
									<div class="flex items-center">
										<input
											id="ignore-ssl-errors"
											type="checkbox"
											class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
											title="Ignore SSL certificate errors"
											aria-label="Enable ignoring SSL errors"
										/>
									</div>
								</div>
							</div>
						</div>

						<!-- Advanced Settings -->
						<div class="p-4 space-y-3">
							<h3 class="text-sm font-medium theme-text-primary flex items-center mb-3">
								<i class="fas fa-cogs text-purple-500 mr-2"></i>
								Advanced Options
							</h3>
							<div class="space-y-3">
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="encoding" class="text-sm font-medium theme-text-primary">
											Response encoding
										</label>
										<span class="text-xs theme-text-muted">Character encoding for response interpretation</span>
									</div>
									<div class="flex items-center">
										<select
											id="encoding"
											class="theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
											title="Select response encoding"
											aria-label="Response encoding"
										>
											<option value="utf-8">UTF-8</option>
											<option value="ascii">ASCII</option>
											<option value="iso-8859-1">ISO-8859-1</option>
											<option value="windows-1252">Windows-1252</option>
										</select>
									</div>
								</div>
								<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
									<div class="flex flex-col">
										<label for="send-cookies" class="text-sm font-medium theme-text-primary cursor-pointer">
											Send cookies
										</label>
										<span class="text-xs theme-text-muted">Include cookies in the request</span>
									</div>
									<div class="flex items-center">
										<input
											id="send-cookies"
											type="checkbox"
											class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
											title="Include cookies in requests"
											aria-label="Enable sending cookies"
											checked
										/>
									</div>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}
	</main>

	<ReplayResponseFooter
		on:toggleExpand={handleFooterToggleExpand}
		on:showHistory={handleFooterShowHistory}
	/>
</div>
