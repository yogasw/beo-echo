<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import ReplayResponseFooter from './ReplayResponseFooter.svelte'; // Import the new footer component
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import ReplayBar from './ReplayBar.svelte';
	
	// Import modular tab components
	import ParamsTab from './tabs/ParamsTab.svelte';
	import AuthorizationTab from './tabs/AuthorizationTab.svelte';
	import HeadersTab from './tabs/HeadersTab.svelte';
	import ScriptTab from './tabs/ScriptTab.svelte';
	import SettingsTab from './tabs/SettingsTab.svelte';
	import ReplayBody from './tabs/ReplayBody.svelte';


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
		activeSection: 'params', // params, headers, body, auth, scripts, settings
		// Data for each tab component
		params: [{ key: '', value: '', description: '', enabled: true }] as Param[],
		headers: [{ key: '', value: '', description: '', enabled: true }] as Header[],
		auth: { type: 'none', config: {} } as AuthConfig,
		scripts: { preRequestScript: '', testScript: '' } as ScriptConfig,
		settings: {
			timeout: 30000,
			followRedirects: true,
			maxRedirects: 5,
			verifySsl: true,
			ignoreSslErrors: false,
			encoding: 'utf-8',
			sendCookies: true,
			storeCookies: true,
			keepAlive: true,
			userAgent: 'Beo-Echo/1.0',
			retryOnFailure: false,
			retryCount: 3,
			retryDelay: 1000
		} as SettingsConfig
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
			activeSection: 'params',
			params: [{ key: '', value: '', description: '', enabled: true }],
			headers: [{ key: '', value: '', description: '', enabled: true }],
			auth: { type: 'none', config: {} },
			scripts: { preRequestScript: '', testScript: '' },
			settings: {
				timeout: 30000,
				followRedirects: true,
				maxRedirects: 5,
				verifySsl: true,
				ignoreSslErrors: false,
				encoding: 'utf-8',
				sendCookies: true,
				storeCookies: true,
				keepAlive: true,
				userAgent: 'Beo-Echo/1.0',
				retryOnFailure: false,
				retryCount: 3,
				retryDelay: 1000
			}
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
				activeSection: 'params',
				params: [{ key: '', value: '', description: '', enabled: true }],
				headers: [{ key: '', value: '', description: '', enabled: true }],
				auth: { type: 'none', config: {} },
				scripts: { preRequestScript: '', testScript: '' },
				settings: {
					timeout: 30000,
					followRedirects: true,
					maxRedirects: 5,
					verifySsl: true,
					ignoreSslErrors: false,
					encoding: 'utf-8',
					sendCookies: true,
					storeCookies: true,
					keepAlive: true,
					userAgent: 'Beo-Echo/1.0',
					retryOnFailure: false,
					retryCount: 3,
					retryDelay: 1000
				}
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
				activeSection: 'params',
				params: [{ key: '', value: '', description: '', enabled: true }],
				headers: [{ key: '', value: '', description: '', enabled: true }],
				auth: { type: 'none', config: {} },
				scripts: { preRequestScript: '', testScript: '' },
				settings: {
					timeout: 30000,
					followRedirects: true,
					maxRedirects: 5,
					verifySsl: true,
					ignoreSslErrors: false,
					encoding: 'utf-8',
					sendCookies: true,
					storeCookies: true,
					keepAlive: true,
					userAgent: 'Beo-Echo/1.0',
					retryOnFailure: false,
					retryCount: 3,
					retryDelay: 1000
				}
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
				activeSection: 'params',
				params: [{ key: '', value: '', description: '', enabled: true }],
				headers: [{ key: '', value: '', description: '', enabled: true }],
				auth: { type: 'none', config: {} },
				scripts: { preRequestScript: '', testScript: '' },
				settings: {
					timeout: 30000,
					followRedirects: true,
					maxRedirects: 5,
					verifySsl: true,
					ignoreSslErrors: false,
					encoding: 'utf-8',
					sendCookies: true,
					storeCookies: true,
					keepAlive: true,
					userAgent: 'Beo-Echo/1.0',
					retryOnFailure: false,
					retryCount: 3,
					retryDelay: 1000
				}
			};
		}
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function setActiveSection(section: string) {
		activeTabContent.activeSection = section;
		dispatch('activeSectionChange', { activeSection: section });
	}

	// Event handlers for tab components
	function handleParamsChange(event: CustomEvent) {
		activeTabContent.params = event.detail.params;
		dispatch('tabContentChange', activeTabContent);
	}

	function handleAuthChange(event: CustomEvent) {
		activeTabContent.auth = {
			type: event.detail.authType,
			config: event.detail.authConfig
		};
		dispatch('tabContentChange', activeTabContent);
	}

	function handleHeadersChange(event: CustomEvent) {
		activeTabContent.headers = event.detail.headers;
		dispatch('tabContentChange', activeTabContent);
	}

	function handleScriptChange(event: CustomEvent) {
		activeTabContent.scripts = {
			preRequestScript: event.detail.preRequestScript,
			testScript: event.detail.testScript
		};
		dispatch('tabContentChange', activeTabContent);
	}

	function handleSettingsChange(event: CustomEvent) {
		activeTabContent.settings = event.detail.settings;
		dispatch('tabContentChange', activeTabContent);
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

	function onExcuteRequest() {
		console.log('Execute request with current configuration:', activeTabContent);
		// Here you would typically send the request using fetch or another HTTP client
		// For now, just log the configuration
		dispatch('executeRequest', { request: activeTabContent });
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
				on:click={onExcuteRequest}
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
			<ParamsTab 
				params={activeTabContent?.params}
				on:paramsChange={handleParamsChange}
			/>
		{:else if activeTabContent.activeSection === 'auth'}
			<AuthorizationTab 
				authType={activeTabContent?.auth?.type}
				authConfig={activeTabContent?.auth?.config}
				on:authChange={handleAuthChange}
			/>
		{:else if activeTabContent.activeSection === 'headers'}
			<HeadersTab 
				headers={activeTabContent?.headers}
				on:headersChange={handleHeadersChange}
			/>
		{:else if activeTabContent.activeSection === 'body'}
			<ReplayBody/>
		{:else if activeTabContent.activeSection === 'scripts'}
			<ScriptTab 
				preRequestScript={activeTabContent?.scripts?.preRequestScript}
				testScript={activeTabContent?.scripts?.testScript}
				on:scriptChange={handleScriptChange}
			/>
		{:else if activeTabContent.activeSection === 'settings'}
			<SettingsTab 
				settings={activeTabContent?.settings}
				on:settingsChange={handleSettingsChange}
			/>
		{/if}
	</main>

	<ReplayResponseFooter
		on:toggleExpand={handleFooterToggleExpand}
		on:showHistory={handleFooterShowHistory}
	/>
</div>
