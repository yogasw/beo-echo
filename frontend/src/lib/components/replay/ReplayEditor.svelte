<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ReplayResponseFooter from './ReplayResponseFooter.svelte'; // Import the new footer component
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import ReplayBar from './ReplayBar.svelte';
	import { replayActions, replayLoading } from '$lib/stores/replay';

	// Import modular tab components
	import ParamsTab from './tabs/ParamsTab.svelte';
	import AuthorizationTab from './tabs/AuthorizationTab.svelte';
	import HeadersTab from './tabs/HeadersTab.svelte';
	// ScriptTab is not currently used - will be implemented in the future
	// import ScriptTab from './tabs/ScriptTab.svelte';
	import SettingsTab from './tabs/SettingsTab.svelte';
	import ReplayBody from './tabs/ReplayBody.svelte';
	import type { Replay } from '$lib/types/Replay';
	import type { Tab } from './types';

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
	
	// Original replay data from API
	export let replayData: Replay | null = null;

	// Active tab content state - UI state and parsed fields separate from raw Replay data
	let activeTabContent: {
		// Basic fields for editing
		id?: string;
		name?: string;
		project_id?: string;
		protocol?: string;
		method: string;
		url: string;
		
		// UI state
		activeSection: string;
		
		// Raw fields from replayData
		headers?: string;
		config?: string;
		metadata?: string;
		payload?: string;
		
		// Parsed fields for editing
		parsedParams?: Array<{key: string; value: string; description: string; enabled: boolean}>;
		parsedHeaders?: Array<{key: string; value: string; description: string; enabled: boolean}>;
		parsedAuth?: {type: string; config: any};
		parsedBody?: string;
		parsedSettings?: any;
	};
	
	// Watch for replayData changes and update activeTabContent
	$: if (replayData) {
		// Initialize with basic fields from replayData
		activeTabContent = {
			id: replayData.id,
			name: replayData.name,
			project_id: replayData.project_id,
			protocol: replayData.protocol || 'http',
			method: replayData.method || 'GET',
			url: replayData.url || '',
			activeSection: 'params',
			
			// Store raw JSON fields
			headers: replayData.headers,
			config: replayData.config,
			metadata: replayData.metadata,
			payload: replayData.payload,
			
			// Parse the JSON fields with default fallbacks
			parsedHeaders: parseHeaders(replayData.headers),
			...parseConfig(replayData.config),
			...parseMetadata(replayData.metadata)
		};
	}
	
	// Add props for execution status and results
	export let isExecuting = false;
	export let executionResult = null;
	
	// Update footer expansion when execution result changes
	$: if (executionResult) {
		isFooterExpanded = true;
	}
	const dispatch = createEventDispatcher();

	// Footer state
	let isFooterExpanded = false;
	let replayResponseFooter; // Reference to the footer component

	// Update footer expansion when execution result changes
	$: if (executionResult) {
		isFooterExpanded = true;
	}

	// Active tab content reset function
	function resetActiveTabContent() {
		// Create default empty strings for JSON fields
		const emptyHeadersJson = JSON.stringify([{ key: '', value: '', description: '', enabled: true }]);
		const emptyConfigJson = JSON.stringify({
			auth: { type: 'none', config: {} },
			settings: {}
		});
		const emptyMetadataJson = JSON.stringify({
			params: [{ key: '', value: '', description: '', enabled: true }],
		});
		
		// Set activeTabContent with raw JSON fields and parsed values
		activeTabContent = {
			method: 'GET',
			url: '',
			protocol: 'http',
			activeSection: 'params',
			
			// Store raw JSON fields
			headers: emptyHeadersJson,
			config: emptyConfigJson,
			metadata: emptyMetadataJson,
			payload: '',
			
			// Parse the JSON fields with default fallbacks
			parsedHeaders: parseHeaders(emptyHeadersJson),
			...parseConfig(emptyConfigJson),
			...parseMetadata(emptyMetadataJson)
		};
	}

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
		resetActiveTabContent();
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
			resetActiveTabContent();
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
			// Create default empty strings for JSON fields
			const emptyHeadersJson = JSON.stringify([{ key: '', value: '', description: '', enabled: true }]);
			const emptyConfigJson = JSON.stringify({
				auth: { type: 'none', config: {} },
				settings: {}
			});
			const emptyMetadataJson = JSON.stringify({
				params: [{ key: '', value: '', description: '', enabled: true }],
			});
			
			// Find the replay data for this tab if it exists
			const tabReplayData = replayData && replayData.id === newActiveTab.id ? replayData : null;
			
			// Start with basic fields from tab or default values
			activeTabContent = {
				id: newActiveTab.id,
				name: newActiveTab.name || 'New Request',
				method: newActiveTab.method || 'GET',
				url: newActiveTab.url || '',
				activeSection: 'params',
				protocol: 'http',
				
				// Use replay data fields if available, otherwise defaults
				headers: tabReplayData?.headers || emptyHeadersJson,
				config: tabReplayData?.config || emptyConfigJson,
				metadata: tabReplayData?.metadata || emptyMetadataJson,
				payload: tabReplayData?.payload || '',
				
				// Parse JSON fields
				parsedHeaders: parseHeaders(tabReplayData?.headers || emptyHeadersJson),
				...parseConfig(tabReplayData?.config || emptyConfigJson),
				...parseMetadata(tabReplayData?.metadata || emptyMetadataJson)
			};
		}
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function switchTab(tabId: string) {
		activeTabId = tabId;
		const tab = tabs.find((t) => t.id === tabId);
		if (tab) {
			// Create default empty strings for JSON fields
			const emptyHeadersJson = JSON.stringify([{ key: '', value: '', description: '', enabled: true }]);
			const emptyConfigJson = JSON.stringify({
				auth: { type: 'none', config: {} },
				settings: {}
			});
			const emptyMetadataJson = JSON.stringify({
				params: [{ key: '', value: '', description: '', enabled: true }],
			});
			
			// Find the replay data for this tab if it exists
			const tabReplayData = replayData && replayData.id === tab.id ? replayData : null;
			
			// Start with basic fields from tab or default values
			activeTabContent = {
				id: tab.id,
				name: tab.name || 'New Request',
				method: tab.method || 'GET',
				url: tab.url || '',
				activeSection: 'params',
				protocol: 'http',
				
				// Use replay data fields if available, otherwise defaults
				headers: tabReplayData?.headers || emptyHeadersJson,
				config: tabReplayData?.config || emptyConfigJson,
				metadata: tabReplayData?.metadata || emptyMetadataJson,
				payload: tabReplayData?.payload || '',
				
				// Parse JSON fields
				parsedHeaders: parseHeaders(tabReplayData?.headers || emptyHeadersJson),
				...parseConfig(tabReplayData?.config || emptyConfigJson),
				...parseMetadata(tabReplayData?.metadata || emptyMetadataJson)
			};
		}
		
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function setActiveSection(section: string) {
		// to display the "Coming Soon" message
		activeTabContent.activeSection = section;
		
		dispatch('activeSectionChange', { activeSection: section });
	}

	// Event handlers for tab components
	function handleParamsChange(event: CustomEvent) {
		// Update parsedParams field
		activeTabContent.parsedParams = event.detail.params;
		
		// Update raw metadata JSON string to reflect the change
		try {
			const metadata = activeTabContent.metadata ? JSON.parse(activeTabContent.metadata) : {};
			metadata.params = event.detail.params;
			activeTabContent.metadata = JSON.stringify(metadata);
		} catch (e) {
			console.error('Failed to update metadata JSON with params change:', e);
			// Create a new metadata object if parsing failed
			const metadata = {
				params: event.detail.params,
			};
			activeTabContent.metadata = JSON.stringify(metadata);
		}
		
		dispatch('tabContentChange', activeTabContent);
	}

	function handleAuthChange(event: CustomEvent) {
		// Update parsedAuth field
		activeTabContent.parsedAuth = {
			type: event.detail.authType,
			config: event.detail.authConfig
		};
		
		// Update raw config JSON string to reflect the change
		try {
			const config = activeTabContent.config ? JSON.parse(activeTabContent.config) : {};
			config.auth = activeTabContent.parsedAuth;
			activeTabContent.config = JSON.stringify(config);
		} catch (e) {
			console.error('Failed to update config JSON with auth change:', e);
			// Create a new config object if parsing failed
			const config = {
				auth: activeTabContent.parsedAuth,
				settings: activeTabContent.parsedSettings || {}
			};
			activeTabContent.config = JSON.stringify(config);
		}
		
		dispatch('tabContentChange', activeTabContent);
	}

	function handleHeadersChange(event: CustomEvent) {
		// Update parsedHeaders field
		activeTabContent.parsedHeaders = event.detail.headers;
		
		// Update raw headers JSON string
		activeTabContent.headers = JSON.stringify(event.detail.headers);
		
		dispatch('tabContentChange', activeTabContent);
	}

	function handleSettingsChange(event: CustomEvent) {
		// Update parsedSettings field
		activeTabContent.parsedSettings = event.detail.settings;
		
		// Update raw config JSON string to reflect the change
		try {
			const config = activeTabContent.config ? JSON.parse(activeTabContent.config) : {};
			config.settings = event.detail.settings;
			activeTabContent.config = JSON.stringify(config);
		} catch (e) {
			console.error('Failed to update config JSON with settings change:', e);
			// Create a new config object if parsing failed
			const config = {
				auth: activeTabContent.parsedAuth || { type: 'none', config: {} },
				settings: event.detail.settings
			};
			activeTabContent.config = JSON.stringify(config);
		}
		
		dispatch('tabContentChange', activeTabContent);
	}

	// Propagate changes upwards
	$: dispatch('tabContentChange', activeTabContent);

	// Event handlers for ReplayResponseFooter
	function handleFooterToggleExpand(event: CustomEvent) {
		isFooterExpanded = event.detail.expanded;
		console.log('Footer expansion toggled:', event.detail.expanded);
		// Handle footer expansion state if needed in ReplayEditor or pass up to parent
	}

	function handleFooterShowHistory() {
		console.log('Show history clicked');
		// Handle showing history, potentially dispatching another event upwards
	}

	function onExcuteRequest() {
		// Prepare request data
		const requestData: {
			method: string;
			url: string;
			headers: Record<string, string>;
			query: Record<string, string>;
			payload: string;
		} = {
			method: activeTabContent.method,
			url: activeTabContent.url,
			headers: {},
			query: {},
			payload: activeTabContent.payload || ''
		};

		// Process headers
		if (activeTabContent.parsedHeaders) {
			activeTabContent.parsedHeaders.forEach(header => {
				if (header.enabled && header.key && header.value) {
					requestData.headers[header.key] = header.value;
				}
			});
		}

		// Process query parameters
		if (activeTabContent.parsedParams) {
			activeTabContent.parsedParams.forEach(param => {
				if (param.enabled && param.key && param.value) {
					requestData.query[param.key] = param.value;
				}
			});
		}

		// Process auth
		if (activeTabContent.parsedAuth?.type !== 'none') {
			const authConfig = activeTabContent.parsedAuth?.config;
			
			switch (activeTabContent.parsedAuth?.type) {
				case 'basic':
					// Add Basic Auth header
					if (authConfig.username) {
						const credentials = btoa(`${authConfig.username}:${authConfig.password || ''}`);
						requestData.headers['Authorization'] = `Basic ${credentials}`;
					}
					break;
				case 'bearer':
					// Add Bearer token header
					if (authConfig.token) {
						requestData.headers['Authorization'] = `Bearer ${authConfig.token}`;
					}
					break;
				case 'apiKey':
					// Add API key as header or query param
					if (authConfig.key && authConfig.value) {
						if (authConfig.in === 'header') {
							requestData.headers[authConfig.key] = authConfig.value;
						} else if (authConfig.in === 'query') {
							requestData.query[authConfig.key] = authConfig.value;
						}
					}
					break;
				// Add other auth types as needed
			}
		}

		// Dispatch the send event with request data
		dispatch('send', requestData);
	}

	function onCancelRequest() {
		console.log('Cancel request clicked');
		// Add logic to cancel the request if needed
	}

	// Utility functions to parse the JSON fields from replayData
	function parseHeaders(headersJson?: string): Array<{key: string; value: string; description: string; enabled: boolean}> {
		if (!headersJson) {
			return [{ key: '', value: '', description: '', enabled: true }];
		}
		
		try {
			const headers = JSON.parse(headersJson);
			return Array.isArray(headers) ? headers : [{ key: '', value: '', description: '', enabled: true }];
		} catch (e) {
			console.error('Failed to parse headers JSON:', e);
			return [{ key: '', value: '', description: '', enabled: true }];
		}
	}
	
	function parseConfig(configJson?: string): { parsedAuth?: {type: string; config: any}; parsedSettings?: any } {
		if (!configJson) {
			return { 
				parsedAuth: { type: 'none', config: {} },
				parsedSettings: {}
			};
		}
		
		try {
			const config = JSON.parse(configJson);
			return {
				parsedAuth: config.auth || { type: 'none', config: {} },
				parsedSettings: config.settings || {}
			};
		} catch (e) {
			console.error('Failed to parse config JSON:', e);
			return { 
				parsedAuth: { type: 'none', config: {} },
				parsedSettings: {}
			};
		}
	}
	
	function parseMetadata(metadataJson?: string): { parsedParams?: Array<{key: string; value: string; description: string; enabled: boolean}>; } {
		if (!metadataJson) {
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }],
			};
		}
		
		try {
			const metadata = JSON.parse(metadataJson);
			return {
				parsedParams: metadata.params || [{ key: '', value: '', description: '', enabled: true }],
			};
		} catch (e) {
			console.error('Failed to parse metadata JSON:', e);
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }],
			};
		}
	}

	// Initialize activeTabContent if it's null
	$: if (activeTabContent === undefined || activeTabContent === null) {
		resetActiveTabContent();
	}
</script>

<!-- Postman-like Request Interface -->
<div class="flex flex-col h-full">
	<ReplayBar {activeTabId} {switchTab} {closeTab} {createNewTab} {tabs} />
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
					class={ThemeUtils.secondaryButton(
						'flex items-center space-x-2 px-3 py-1.5 rounded-md text-sm border theme-border hover:shadow-md transition-all duration-200'
					)}
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
		<div
			class={ThemeUtils.themeBgSecondary(
				'flex items-center border theme-border rounded-lg shadow-sm'
			)}
		>
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
				class={$replayLoading.execute
					? 'bg-gray-600 hover:bg-gray-700 text-white px-4 py-2.5 rounded-r-lg flex items-center space-x-1 shadow-sm hover:shadow-md transition-all duration-200'
					: ThemeUtils.primaryButton(
							'px-4 py-2.5 rounded-r-lg space-x-1 shadow-sm hover:shadow-md transition-all duration-200'
						)}
				title={$replayLoading.execute ? 'Cancel HTTP request' : 'Send HTTP request'}
				aria-label={$replayLoading.execute ? 'Cancel request' : 'Send request'}
				on:click={$replayLoading.execute ? onCancelRequest : onExcuteRequest}
			>
				{#if $replayLoading.execute}
					<span>Cancel</span>
					<i class="fas fa-times text-sm"></i>
				{:else}
					<span>Send</span>
					<i class="fas fa-paper-plane text-sm"></i>
				{/if}
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
				<!-- <button
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
				</button> -->
			</div>
		</div>

		<!-- Dynamic content based on active section -->
		{#if activeTabContent.activeSection === 'params'}
			<ParamsTab params={activeTabContent?.parsedParams} on:paramsChange={handleParamsChange} />
		{:else if activeTabContent.activeSection === 'auth'}
			<AuthorizationTab
				authType={activeTabContent?.parsedAuth?.type}
				authConfig={activeTabContent?.parsedAuth?.config}
				on:authChange={handleAuthChange}
			/>
		{:else if activeTabContent.activeSection === 'headers'}
			<HeadersTab headers={activeTabContent?.parsedHeaders} on:headersChange={handleHeadersChange} />
		{:else if activeTabContent.activeSection === 'body'}
			<ReplayBody payload={activeTabContent?.payload} />
		{:else if activeTabContent.activeSection === 'settings'}
			<SettingsTab settings={activeTabContent?.parsedSettings} on:settingsChange={handleSettingsChange} />
		{/if}
	</main>

	<ReplayResponseFooter
		bind:isExpanded={isFooterExpanded}
		{executionResult}
		on:toggleExpand={handleFooterToggleExpand}
		on:showHistory={handleFooterShowHistory}
	/>
</div>
