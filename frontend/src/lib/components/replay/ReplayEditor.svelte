<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { onMount } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import ReplayResponseFooter from './ReplayResponseFooter.svelte'; // Import the new footer component
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import ReplayBar from './ReplayBar.svelte';
	import { replayActions, replayLoading } from '$lib/stores/replay';

	// Import modular tab components
	import ParamsTab from './tabs/ParamsTab.svelte';
	import AuthorizationTab from './tabs/AuthorizationTab.svelte';
	import HeadersTab from './tabs/HeadersTab.svelte';
	import FolderOverviewTab from './tabs/FolderOverviewTab.svelte';
	// ScriptTab is not currently used - will be implemented in the future
	// import ScriptTab from './tabs/ScriptTab.svelte';
	import SettingsTab from './tabs/SettingsTab.svelte';
	import ReplayBody from './tabs/ReplayBody.svelte';
	import { createTabActions } from './tabs/tabActions';
	import type { ExecuteReplayResponse, Replay } from '$lib/types/Replay';
	import type { Tab } from './types';

	export let tabs: Tab[] = [];
	export let activeTabId = '';

	$: activeTab = tabs.find((t) => t.id === activeTabId);



	const dispatch = createEventDispatcher();

	// Function to flag a tab as unsaved whenever its content changes.
	function markUnsaved() {
		if (activeTab && !activeTab.isUnsaved) {
			activeTab.isUnsaved = true;
			tabs = [...tabs]; // trigger Svelte reactivity
		}
	}

	function triggerUpdate() {
		tabs = [...tabs];
	}

	// Initialize activeTab.content if it's new or hasn't parsed the original replay yet.
	$: if (activeTab && activeTab.replay && !(activeTab.content as any)?._isParsed) {
		const replay = activeTab.replay;
		const config = parseConfig(replay.config);
		const meta = parseMetadata(replay.metadata);
		
		activeTab.content = {
			...activeTab.content,
			method: replay.method || 'GET',
			url: replay.url || '',
			activeSection: activeTab.content?.activeSection || 'params',
			headers: parseHeaders(replay.headers),
			auth: config.parsedAuth || { type: 'none', config: {} },
			settings: config.parsedSettings || {},
			params: parseUrlParams(replay.url || '', meta.parsedParams || []),
			body: {
				...activeTab.content?.body,
				type: meta.parsedBodyType || 'none',
				content: replay.payload || ''
			}
		} as any;
		(activeTab.content as any)._isParsed = true;
		triggerUpdate();
	}

	// Add props for execution status and results
	export let executionResult: ExecuteReplayResponse | null = null;

	// Footer state
	let isFooterExpanded = false;
	let replayResponseFooter;

	$: if (executionResult) {
		isFooterExpanded = true;
	}

	// Flag to control when content changes should be automatically dispatched
	let shouldAutoDispatch = false;

	// Initialize auto-dispatch after component is fully loaded
	onMount(() => {
		shouldAutoDispatch = true;
	});

	function closeOtherTabs(tabIdToKeep: string) {
		tabs = tabs.filter(t => t.id === tabIdToKeep);
		activeTabId = tabIdToKeep;
		dispatch('tabschange', { tabs, activeTabId });
	}

	function closeAllTabs() {
		tabs = [];
		activeTabId = '';
		dispatch('tabschange', { tabs: [], activeTabId: '' });
	}

	function duplicateTab(tabId: string) {
		const tabToCopy = tabs.find(t => t.id === tabId);
		if (!tabToCopy) return;
		const newTab = JSON.parse(JSON.stringify(tabToCopy));
		newTab.id = `tab-${Date.now()}`;
		newTab.isUnsaved = true;
		tabs = [...tabs, newTab];
		activeTabId = newTab.id;
		dispatch('tabschange', { tabs, activeTabId });
	}

	function switchTab(tabId: string) {
		activeTabId = tabId;
		dispatch('tabschange', { tabs, activeTabId });
	}

	function createNewTab() {
		const newTabId = `tab-${Date.now()}`;
		const newTab: Tab = {
			id: newTabId,
			isUnsaved: true,
			itemType: 'request',
			content: {
				method: 'GET',
				url: '',
				activeSection: 'params',
				headers: [{ key: '', value: '', description: '', enabled: true }],
				params: [{ key: '', value: '', description: '', enabled: true }],
				auth: { type: 'none', config: {} },
				settings: {},
				body: { type: 'none', content: '' }
			} as any
		};
		tabs = [...tabs, newTab];
		activeTabId = newTabId;
		dispatch('tabschange', { tabs, activeTabId });
	}

	function closeTab(tabId: string) {
		const tabIndex = tabs.findIndex((tab) => tab.id === tabId);
		tabs = tabs.filter((tab) => tab.id !== tabId);

		if (tabs.length === 0) {
			activeTabId = '';
			dispatch('tabschange', { tabs: [], activeTabId: '' });
			return;
		}

		if (activeTabId === tabId) {
			if (tabIndex > 0) {
				activeTabId = tabs[tabIndex - 1].id;
			} else {
				activeTabId = tabs[0].id;
			}
		}
		dispatch('tabschange', { tabs, activeTabId });
	}

	function setActiveSection(section: string) {
		if (activeTab?.content) {
			activeTab.content.activeSection = section;
			tabs = [...tabs];
		}
		dispatch('activeSectionChange', { activeSection: section });
	}

	// Event handlers for tab components
	function handleParamsChange(event: CustomEvent) {
		if (activeTab?.content) {
			activeTab.content.params = event.detail.params;
			activeTab.content.url = getUrlFromParams(activeTab.content.url, event.detail.params);
			markUnsaved();
		}
	}

	function handleAuthChange(event: CustomEvent) {
		if (activeTab?.content) {
			activeTab.content.auth = {
				type: event.detail.authType,
				config: event.detail.authConfig
			};
			markUnsaved();
		}
	}

	function handleHeadersChange(event: CustomEvent) {
		if (activeTab?.content) {
			activeTab.content.headers = event.detail.headers;
			markUnsaved();
		}
	}

	function handleSettingsChange(event: CustomEvent) {
		if (activeTab?.content) {
			activeTab.content.settings = event.detail.settings;
			markUnsaved();
		}
	}

	function handleBodyContentChange(event: CustomEvent) {
		if (activeTab?.content) {
			if (!activeTab.content.body) {
				activeTab.content.body = { type: 'none', content: '' } as any;
			}
			if (event.detail.payload !== undefined) {
				activeTab.content.body.content = event.detail.payload;
			}
			if (event.detail.bodyType !== undefined) {
				activeTab.content.body.type = event.detail.bodyType;
			}
			markUnsaved();
		}
	}

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

	async function onSaveRequest() {
		if (!activeTab || !$selectedWorkspace || !$selectedProject) {
			toast.error('No workspace or project selected');
			return;
		}

		try {
			replayActions.setLoading('save', true);

			const payload: any = {
				name: activeTab.replay?.name || activeTab.folder?.name || 'New Request',
				protocol: activeTab.replay?.protocol || 'http',
				method: activeTab.content?.method || 'GET',
				url: activeTab.content?.url || '',
				headers: [],
				payload: activeTab.content?.body?.content || '',
				config: {
					auth: activeTab.content?.auth || { type: 'none', config: {} },
					settings: activeTab.content?.settings || {}
				}
			};

			const validParams = (activeTab.content?.params || []).filter((p: any) => p.key);
			
			// Parse existing metadata to preserve other properties (like bodyType)
			let metaObj: any = {};
			try {
				if (activeTab.replay?.metadata) {
					metaObj = typeof activeTab.replay.metadata === 'string' 
						? JSON.parse(activeTab.replay.metadata) 
						: { ...(activeTab.replay.metadata as any) };
				}
			} catch(e) {
				console.warn('Failed to parse metadata when saving', e);
			}

			if (validParams.length > 0) {
				metaObj.params = validParams;
			}
			metaObj.bodyType = activeTab.content?.body?.type || 'none';
			
			// Always send metadata, even if it's empty, so we don't nullify database fields 
			// if they had bodyType
			payload.metadata = metaObj;

			// Add folder ID if available on the original replay data or current tab
			if (activeTab.replay?.folder_id) {
				payload.folder_id = activeTab.replay.folder_id;
			}

			// Process headers as array with description (disabled = deleted)
			if (activeTab.content?.headers) {
				payload.headers = activeTab.content.headers
					.filter((h: any) => h.enabled && h.key)
					.map((h: any) => ({
						key: h.key,
						value: h.value || '',
						description: h.description || ''
					}));
			}

			if (activeTab.id && !activeTab.id.startsWith('tab-')) {
				// Update existing
				const res = await replayApi.updateReplay(
					$selectedWorkspace.id,
					$selectedProject.id,
					activeTab.id,
					payload
				);
				
				// Mark current tab as saved
				const tabIndex = tabs.findIndex(t => t.id === activeTab.id);
				if (tabIndex !== -1) {
					tabs[tabIndex].isUnsaved = false;
					tabs[tabIndex].replay = {
						...res.replay
					};
					tabs = [...tabs];
				}
				
				dispatch('updated', res.replay);
			} else {
				// Create new
				const res = await replayApi.createReplay(
					$selectedWorkspace.id,
					$selectedProject.id,
					payload
				);
				
				// Update active tab with new ID
				const tabIndex = tabs.findIndex(t => t.id === activeTabId);
				if (tabIndex !== -1) {
					tabs[tabIndex].id = res.replay.id;
					tabs[tabIndex].isUnsaved = false;
					tabs[tabIndex].replay = res.replay;
					activeTabId = res.replay.id;
					tabs = [...tabs];
				}
				
				dispatch('created', res.replay);
			}
		} catch (err: any) {
			toast.error(err.message || 'Failed to save request');
		} finally {
			replayActions.setLoading('save', false);
		}
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
			method: activeTab?.content?.method || 'GET',
			url: activeTab?.content?.url || '',
			headers: {},
			query: {},
			payload: activeTab?.content?.body?.content || ''
		};

		// Process headers
		if (activeTab?.content?.headers) {
			activeTab.content.headers.forEach((header: any) => {
				if (header.enabled && header.key && header.value) {
					requestData.headers[header.key] = header.value;
				}
			});
		}

		// Process query parameters
		if (activeTab?.content?.params) {
			activeTab.content.params.forEach((param: any) => {
				if (param.enabled && param.key && param.value) {
					requestData.query[param.key] = param.value;
				}
			});
		}

		// Process auth
		if (activeTab?.content?.auth?.type !== 'none') {
			const authConfig = activeTab?.content?.auth?.config;

			switch (activeTab?.content?.auth?.type) {
				case 'basic':
					// Add Basic Auth header
					if (authConfig?.username) {
						const credentials = btoa(`${authConfig?.username}:${authConfig?.password || ''}`);
						requestData.headers['Authorization'] = `Basic ${credentials}`;
					}
					break;
				case 'bearer':
					// Add Bearer token header
					if (authConfig?.token) {
						requestData.headers['Authorization'] = `Bearer ${authConfig?.token}`;
					}
					break;
				case 'apiKey':
					// Add API key as header or query param
					if (authConfig?.key && authConfig?.value) {
						if (authConfig?.in === 'header') {
							requestData.headers[authConfig.key] = authConfig.value;
						} else if (authConfig?.in === 'query') {
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
	function parseHeaders(
		headersJson?: string
	): Array<{ key: string; value: string; description: string; enabled: boolean }> {
		if (!headersJson) {
			return [{ key: '', value: '', description: '', enabled: true }];
		}

		try {
			const headers = JSON.parse(headersJson);
			if (Array.isArray(headers)) {
				if (headers.length > 0) {
					return headers.map((h: any) => ({
						key: h.key || '',
						value: h.value || '',
						description: h.description || '',
						enabled: true
					}));
				}
				return [{ key: '', value: '', description: '', enabled: true }];
			}
			// Handle object format: Record<string, string> from backend
			if (typeof headers === 'object' && headers !== null) {
				const entries = Object.entries(headers);
				if (entries.length > 0) {
					return entries.map(([key, value]) => ({
						key,
						value: String(value),
						description: '',
						enabled: true
					}));
				}
			}
			return [{ key: '', value: '', description: '', enabled: true }];
		} catch (e) {
			console.error('Failed to parse headers JSON:', e);
			return [{ key: '', value: '', description: '', enabled: true }];
		}
	}

	function parseConfig(configJson?: string): {
		parsedAuth?: { type: string; config: any };
		parsedSettings?: any;
	} {
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

	function parseMetadata(metadataJson?: string): {
		parsedParams?: Array<{ key: string; value: string; description: string; enabled: boolean }>;
		parsedBodyType?: string;
	} {
		if (!metadataJson) {
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }],
				parsedBodyType: 'none'
			};
		}

		try {
			const metadata = JSON.parse(metadataJson);
			return {
				parsedParams: metadata.params || [{ key: '', value: '', description: '', enabled: true }],
				parsedBodyType: metadata.bodyType || 'none'
			};
		} catch (e) {
			console.error('Failed to parse metadata JSON:', e);
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }],
				parsedBodyType: 'none'
			};
		}
	}

	function parseUrlParams(url: string, existingParams?: Array<{ key: string; value: string; description: string; enabled: boolean }>): Array<{ key: string; value: string; description: string; enabled: boolean }> {
		try {
			const [, queryString] = (url || '').split('?');
			const paramsList: Array<{ key: string; value: string; description: string; enabled: boolean }> = [];
			
			const existingKeyMap = new Map();
			if (existingParams) {
				existingParams.forEach(p => {
					if (p.key) existingKeyMap.set(p.key, p);
				});
			}

			const usedKeys = new Set<string>();

			if (queryString) {
				const urlParams = new URLSearchParams(queryString);
				for (const [key, value] of urlParams.entries()) {
					const existing = existingKeyMap.get(key);
					paramsList.push({ 
						key, 
						value, 
						description: existing ? existing.description : '', 
						enabled: true 
					});
					usedKeys.add(key);
				}
			}

			// Always keep existing unused parameters (keep their description), just mark them as disabled if they're not in the URL
			if (existingParams) {
				existingParams.forEach(p => {
					if (p.key && !usedKeys.has(p.key)) {
						paramsList.push({ ...p, enabled: false });
					}
				});
			}

			paramsList.push({ key: '', value: '', description: '', enabled: true });
			return paramsList;
		} catch (e) {
			return [{ key: '', value: '', description: '', enabled: true }];
		}
	}

	function getUrlFromParams(baseUrl: string, params: Array<{ key: string; value: string; enabled: boolean }>): string {
		try {
			const [base] = (baseUrl || '').split('?');
			const searchParams = new URLSearchParams();
			let hasParams = false;
			params.forEach(p => {
				if (p.enabled && p.key.trim() !== '') {
					searchParams.append(p.key.trim(), p.value);
					hasParams = true;
				}
			});
			return hasParams ? `${base}?${searchParams.toString()}` : base;
		} catch (e) {
			return baseUrl || '';
		}
	}


	let isEditingTitle = false;
	let editTitleValue = '';
	let titleInputRef: HTMLInputElement;

	function startEditTitle() {
		const tab = tabs.find(t => t.id === activeTabId);
		if (tab) {
			editTitleValue = tab.replay?.name || '';
			isEditingTitle = true;
			setTimeout(() => {
				titleInputRef?.focus();
				titleInputRef?.select();
			}, 0);
		}
	}

	async function saveTitle() {
		if (!isEditingTitle) return;
		isEditingTitle = false;
		
		const tabIndex = tabs.findIndex(t => t.id === activeTabId);
		if (tabIndex === -1) return;
		
		const oldName = tabs[tabIndex].replay?.name || tabs[tabIndex].folder?.name;
		const newName = editTitleValue.trim();
		
		if (newName && newName !== oldName) {
			if (tabs[tabIndex].replay) {
				tabs[tabIndex].replay!.name = newName;
			}
			if (activeTab) {
				if (activeTab.itemType === 'folder' && activeTab.folder) {
					activeTab.folder.name = newName;
				}
			}
			tabs = [...tabs]; // trigger Svelte reactivity
			
			// Save using API
			if (tabs[tabIndex].itemType === 'folder' && tabs[tabIndex].folder?.id) {
				try {
					await replayApi.updateFolder($selectedWorkspace?.id as string, $selectedProject?.id as string, tabs[tabIndex].folder.id, { name: newName });
					dispatch('updated', { ...tabs[tabIndex].folder, name: newName });
				} catch (err) {
					console.error('Failed to update folder name:', err);
					toast.error('Failed to update folder name');
				}
			} else if (tabs[tabIndex].itemType !== 'folder') {
				// We update request name right away, then it's technically unsaved 
				// or we just save it immediately using onSaveRequest logic if it exists
				if (activeTab?.id && !activeTab.id.startsWith('tab-')) {
					// Mark tab as unsaved so user knows they need to click Save button,
					// or we can auto save via API. Replay list auto-save title:
					try {
						await replayApi.updateReplay($selectedWorkspace?.id as string, $selectedProject?.id as string, activeTab.id, { name: newName });
						dispatch('updated', { id: activeTab.id, name: newName });
					} catch (err) {
						tabs[tabIndex].isUnsaved = true; // wait for explicit save
					}
				} else {
					tabs[tabIndex].isUnsaved = true;
				}
			}
		}
	}
	
	function handleTitleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			saveTitle();
		} else if (e.key === 'Escape') {
			isEditingTitle = false;
		}
	}

</script>

<!-- Postman-like Request Interface -->
<div class="flex flex-col h-full">
	<ReplayBar {activeTabId} {switchTab} {closeTab} {closeOtherTabs} {closeAllTabs} {duplicateTab} {createNewTab} {tabs} />
	<!-- Main content -->
	<main class="flex-grow p-4 space-y-4 flex flex-col overflow-y-auto">
		<!-- Title and actions -->
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-2">
				<i class="fas {activeTab?.itemType === 'folder' ? 'fa-folder-open text-orange-500' : 'fa-file-alt text-blue-500'} text-xl"></i>
				{#if isEditingTitle}
					<input
						bind:this={titleInputRef}
						bind:value={editTitleValue}
						on:keydown={handleTitleKeydown}
						on:blur={saveTitle}
						class="text-lg font-semibold px-2 py-0.5 border border-blue-500 rounded focus:outline-none theme-bg-primary theme-text-primary"
						type="text"
					/>
				{:else}
					<h1 
						class="text-lg font-semibold theme-text-primary cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-800 px-2 py-0.5 rounded transition-colors group flex items-center gap-2"
						on:dblclick={startEditTitle}
					>
						{tabs.find((t) => t.id === activeTabId)?.replay?.name || tabs.find((t) => t.id === activeTabId)?.folder?.name || 'New Request'}
						<!-- svelte-ignore a11y_click_events_have_key_events -->
						<!-- svelte-ignore a11y_no_static_element_interactions -->
						<i class="fas fa-pencil-alt text-xs opacity-0 group-hover:opacity-50 transition-opacity" on:click={startEditTitle}></i>
					</h1>
				{/if}
			</div>
			<div class="flex items-center space-x-2">
				{#if activeTab?.itemType !== 'folder'}
				<button
					class={ThemeUtils.secondaryButton(
						'flex items-center space-x-2 px-3 py-1.5 rounded-md text-sm border theme-border hover:shadow-md transition-all duration-200'
					)}
					class:opacity-50={$replayLoading.save || !(tabs.find(t => t.id === activeTabId)?.isUnsaved)}
					class:cursor-not-allowed={$replayLoading.save || !(tabs.find(t => t.id === activeTabId)?.isUnsaved)}
					title="Save current request configuration"
					aria-label="Save request"
					on:click={onSaveRequest}
					disabled={$replayLoading.save || !(tabs.find(t => t.id === activeTabId)?.isUnsaved)}
				>
					{#if $replayLoading.save}
						<i class="fas fa-circle-notch fa-spin text-sm"></i>
						<span>Saving...</span>
					{:else}
						<i class="fas fa-save text-sm"></i>
						<span>Save</span>
					{/if}
				</button>
				{/if}
			</div>
		</div>

		{#if activeTab?.itemType === 'folder'}
			<!-- Folder Overview Mode -->
			<div class="flex-grow flex flex-col min-h-0 w-full">
				<FolderOverviewTab
					folder={activeTab.folder}
				/>
			</div>
		{:else if activeTab && activeTab.content}
			<!-- Request builder -->
			<div
				class={ThemeUtils.themeBgSecondary(
					'flex items-center border theme-border rounded-lg shadow-sm'
				)}
			>
			<div class="relative">
				<select
					bind:value={activeTab.content.method}
					on:change={markUnsaved}
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
				bind:value={activeTab.content.url}
				on:input={(e) => {
					markUnsaved();
					if (activeTab?.content) {
						activeTab.content.params = parseUrlParams(e.currentTarget.value, activeTab.content.params);
					}
				}}
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
					class="py-2 px-1 border-b-2 {activeTab.content.activeSection === 'params'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure query parameters"
					aria-label="Parameters tab"
					role="tab"
					aria-selected={activeTab.content.activeSection === 'params'}
					on:click={() => setActiveSection('params')}
				>
					Params
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTab.content.activeSection === 'auth'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request authorization"
					aria-label="Authorization tab"
					role="tab"
					aria-selected={activeTab.content.activeSection === 'auth'}
					on:click={() => setActiveSection('auth')}
				>
					Authorization
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTab.content.activeSection === 'headers'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request headers"
					aria-label="Headers tab"
					role="tab"
					aria-selected={activeTab.content.activeSection === 'headers'}
					on:click={() => setActiveSection('headers')}
				>
					Headers
				</button>
				<button
					class="py-2 px-1 border-b-2 {activeTab.content.activeSection === 'body'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request body"
					aria-label="Body tab"
					role="tab"
					aria-selected={activeTab.content.activeSection === 'body'}
					on:click={() => setActiveSection('body')}
				>
					Body
				</button>
				<!-- <button
					class="py-2 px-1 border-b-2 {activeTab.content.activeSection === 'settings'
						? 'border-orange-600 text-orange-600'
						: 'border-transparent hover:theme-text-primary'} transition-colors duration-200"
					title="Configure request settings"
					aria-label="Settings tab"
					role="tab"
					aria-selected={activeTab.content.activeSection === 'settings'}
					on:click={() => setActiveSection('settings')}
				>
					Settings
				</button> -->
			</div>
		</div>

		<!-- Dynamic content based on active section -->
		{#if activeTab.content.activeSection === 'params'}
			<ParamsTab params={activeTab.content?.params} on:paramsChange={handleParamsChange} />
		{:else if activeTab.content.activeSection === 'auth'}
			<AuthorizationTab
				authType={activeTab.content?.auth?.type}
				authConfig={activeTab.content?.auth?.config}
				on:authChange={handleAuthChange}
			/>
		{:else if activeTab.content.activeSection === 'headers'}
			<HeadersTab
				headers={activeTab.content?.headers}
				on:headersChange={handleHeadersChange}
			/>
		{:else if activeTab.content.activeSection === 'body'}
			<ReplayBody 
				payload={activeTab.content?.body?.content}
				metadata={JSON.stringify({ params: activeTab.content?.params, bodyType: activeTab.content?.body?.type })}
				protocol={activeTab.replay?.protocol || 'http'}
				on:change={handleBodyContentChange}
			/>
		{:else if activeTab.content.activeSection === 'settings'}
			<SettingsTab
				settings={activeTab.content?.settings}
				on:settingsChange={handleSettingsChange}
			/>
		{/if}
		{/if} <!-- End of request vs folder branch -->
	</main>

	{#if activeTab?.itemType === 'request'}
		<ReplayResponseFooter
			bind:isExpanded={isFooterExpanded}
			{executionResult}
			on:toggleExpand={handleFooterToggleExpand}
			on:showHistory={handleFooterShowHistory}
		/>
	{/if}
</div>
