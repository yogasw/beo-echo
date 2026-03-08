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
		itemType?: 'request' | 'folder';
		folder?: any;
		folder_id?: string;

		// Raw fields from replayData
		headers?: string;
		config?: string;
		metadata?: string;
		payload?: string;

		// Parsed fields for editing
		parsedParams?: Array<{ key: string; value: string; description: string; enabled: boolean }>;
		parsedHeaders?: Array<{ key: string; value: string; description: string; enabled: boolean }>;
		parsedAuth?: { type: string; config: any };
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
			itemType: (replayData as any).itemType === 'folder' ? 'folder' : 'request',
			folder: (replayData as any).itemType === 'folder' ? replayData : undefined,

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
		
		// If editing an existing item that has a folder_id, keep it in activeTabContent:
		if (replayData.folder_id) {
			activeTabContent.folder_id = replayData.folder_id;
		} else {
			// Try finding the tab by ID to inherit any initialized folder_id (from Create Request action)
			const currentTab = tabs.find((t) => t.id === replayData.id || t.id === activeTabId);
			if (currentTab?.folder_id) {
				activeTabContent.folder_id = currentTab.folder_id;
			}
		}

		lastDispatchedContent = JSON.stringify(getStructuredContentForStorage());
		lastDispatchedTabId = replayData.id || activeTabId || '';
	} else {
		// If replayData is null, we might be on a brand new tab, attempt to load folder_id from the tab array
		const currentTab = tabs.find((t) => t.id === activeTabId);
		if (currentTab?.folder_id && activeTabContent) {
			activeTabContent.folder_id = currentTab.folder_id;
		}
	}

	// Add props for execution status and results
	export let executionResult: ExecuteReplayResponse | null = null;

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
		const emptyHeadersJson = JSON.stringify([
			{ key: '', value: '', description: '', enabled: true }
		]);
		const emptyConfigJson = JSON.stringify({
			auth: { type: 'none', config: {} },
			settings: {}
		});
		const emptyMetadataJson = JSON.stringify({
			params: [{ key: '', value: '', description: '', enabled: true }]
		});

		// Set activeTabContent with raw JSON fields and parsed values
		activeTabContent = {
			method: 'GET',
			url: '',
			protocol: 'http',
			activeSection: 'params',
			itemType: 'request',
			folder: undefined,

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
		// Reset tracker when resetting tab
		lastDispatchedContent = JSON.stringify(getStructuredContentForStorage());
	}

	// Flag to control when content changes should be automatically dispatched
	let shouldAutoDispatch = false;

	// Initialize auto-dispatch after component is fully loaded
	onMount(() => {
		shouldAutoDispatch = true;
	});

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
		
		shouldAutoDispatch = false; // Disable during reset
		resetActiveTabContent();
		shouldAutoDispatch = true; // Re-enable after reset
		
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	function closeTab(tabId: string) {
		const tabIndex = tabs.findIndex((tab) => tab.id === tabId);
		tabs = tabs.filter((tab) => tab.id !== tabId);

		if (tabs.length === 0) {
			activeTabId = '';
			shouldAutoDispatch = false;
			resetActiveTabContent();
			shouldAutoDispatch = true;
			dispatch('tabschange', { tabs: [], activeTabId: '', activeTabContent: null });
			return;
		}

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
			shouldAutoDispatch = false; // Disable during content update
			
			// Create default empty strings for JSON fields
			const emptyHeadersJson = JSON.stringify([
				{ key: '', value: '', description: '', enabled: true }
			]);
			const emptyConfigJson = JSON.stringify({
				auth: { type: 'none', config: {} },
				settings: {}
			});
			const emptyMetadataJson = JSON.stringify({
				params: [{ key: '', value: '', description: '', enabled: true }]
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
				itemType: newActiveTab.itemType || 'request',
				folder: newActiveTab.folder,

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
			
			lastDispatchedContent = JSON.stringify(getStructuredContentForStorage());
			shouldAutoDispatch = true; // Re-enable after content update
		}
		dispatch('tabschange', { tabs, activeTabId, activeTabContent });
	}

	const tabActions = createTabActions({
		getTabs: () => tabs,
		setTabs: (t) => { tabs = t; },
		getActiveTabId: () => activeTabId,
		setActiveTabId: (id) => { activeTabId = id; },
		getActiveTabContent: () => activeTabContent,
		dispatchTabsChange: () => dispatch('tabschange', { tabs, activeTabId, activeTabContent }),
		triggerResetActiveTabContent: () => {
			shouldAutoDispatch = false;
			resetActiveTabContent();
			shouldAutoDispatch = true;
			dispatch('tabschange', { tabs: [], activeTabId: '', activeTabContent: null });
		}
	});

	function closeOtherTabs(tabIdToKeep: string) {
		tabActions.closeOtherTabs(tabIdToKeep);
	}

	function closeAllTabs() {
		tabActions.closeAllTabs();
	}

	function duplicateTab(tabId: string) {
		tabActions.duplicateTab(tabId);
	}

	function switchTab(tabId: string) {
		activeTabId = tabId;
		const tab = tabs.find((t) => t.id === tabId);
		if (tab) {
			shouldAutoDispatch = false; // Disable during content update
			
			// Create default empty strings for JSON fields
			const emptyHeadersJson = JSON.stringify([
				{ key: '', value: '', description: '', enabled: true }
			]);
			const emptyConfigJson = JSON.stringify({
				auth: { type: 'none', config: {} },
				settings: {}
			});
			const emptyMetadataJson = JSON.stringify({
				params: [{ key: '', value: '', description: '', enabled: true }]
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
				itemType: tab.itemType || 'request',
				folder: tab.folder,

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
		
		lastDispatchedContent = JSON.stringify(getStructuredContentForStorage());
		shouldAutoDispatch = true; // Re-enable after content update
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
				params: event.detail.params
			};
			activeTabContent.metadata = JSON.stringify(metadata);
		}

		dispatch('tabContentChange', getStructuredContentForStorage());
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

		dispatch('tabContentChange', getStructuredContentForStorage());
	}

	function handleHeadersChange(event: CustomEvent) {
		// Update parsedHeaders field
		activeTabContent.parsedHeaders = event.detail.headers;

		// Update raw headers JSON string
		activeTabContent.headers = JSON.stringify(event.detail.headers);

		dispatch('tabContentChange', getStructuredContentForStorage());
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

		dispatch('tabContentChange', getStructuredContentForStorage());
	}

	function handleBodyContentChange(event: CustomEvent) {
		// Update payload field
		activeTabContent.payload = event.detail.payload;
		
		dispatch('tabContentChange', getStructuredContentForStorage());
	}

	// Function to convert activeTabContent to structured format for storage
	function getStructuredContentForStorage() {
		return {
			method: activeTabContent.method,
			url: activeTabContent.url,
			activeSection: activeTabContent.activeSection,
			payload: activeTabContent.payload,
			params: activeTabContent.parsedParams || [],
			headers: activeTabContent.parsedHeaders || [],
			auth: activeTabContent.parsedAuth || { type: 'none', config: {} },
			settings: activeTabContent.parsedSettings || {}
		};
	}

	// Track the last dispatched content and its tab ID to prevent spurious changes flag when switching tabs
	let lastDispatchedContent = '';
	let lastDispatchedTabId = '';

	// Propagate changes upwards with structured content
	$: if (activeTabContent && shouldAutoDispatch) {
		const newContent = getStructuredContentForStorage();
		const newContentStr = JSON.stringify(newContent);
		
		// If we just switched to a new tab or loaded fresh data, reset the baseline without dispatching
		if (activeTabId !== lastDispatchedTabId) {
			lastDispatchedTabId = activeTabId;
			lastDispatchedContent = newContentStr;
		} 
		// If it's the exact same tab but content truly changed, dispatch the change
		else if (newContentStr !== lastDispatchedContent) {
			lastDispatchedContent = newContentStr;
			dispatch('tabContentChange', newContent);
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
		if (!$selectedWorkspace || !$selectedProject) {
			toast.error('No workspace or project selected');
			return;
		}

		try {
			replayActions.setLoading('save', true);

			const payload: any = {
				name: activeTabContent.name || 'New Request',
				protocol: activeTabContent.protocol || 'http',
				method: activeTabContent.method || 'GET',
				url: activeTabContent.url || '',
				headers: {},
				payload: activeTabContent.payload || '',
				config: {
					auth: activeTabContent.parsedAuth || { type: 'none', config: {} },
					settings: activeTabContent.parsedSettings || {}
				}
			};

			const validParams = (activeTabContent.parsedParams || []).filter((p: any) => p.key && p.enabled);
			if (validParams.length > 0) {
				payload.metadata = { params: validParams };
			}

			// Add folder ID if available on the original replay data or current tab
			if (activeTabContent.folder_id) {
				payload.folder_id = activeTabContent.folder_id;
			} else if (replayData && replayData.folder_id) {
				payload.folder_id = replayData.folder_id;
			}

			// Process headers as object map
			if (activeTabContent.parsedHeaders) {
				activeTabContent.parsedHeaders.forEach((h: any) => {
					if (h.enabled && h.key) {
						payload.headers[h.key] = h.value;
					}
				});
			}

			if (activeTabContent.id && !activeTabContent.id.startsWith('tab-')) {
				// Update existing
				const res = await replayApi.updateReplay(
					$selectedWorkspace.id,
					$selectedProject.id,
					activeTabContent.id,
					payload
				);
				
				// Mark current tab as saved
				const tabIndex = tabs.findIndex(t => t.id === activeTabContent.id);
				if (tabIndex !== -1) {
					tabs[tabIndex].isUnsaved = false;
					tabs[tabIndex].name = res.replay.name;
					tabs[tabIndex].url = res.replay.url;
					tabs[tabIndex].method = res.replay.method;
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
				activeTabContent.id = res.replay.id;
				
				const tabIndex = tabs.findIndex(t => t.id === activeTabId);
				if (tabIndex !== -1) {
					tabs[tabIndex].id = res.replay.id;
					tabs[tabIndex].isUnsaved = false;
					tabs[tabIndex].name = res.replay.name;
					tabs[tabIndex].url = res.replay.url;
					tabs[tabIndex].method = res.replay.method;
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
			method: activeTabContent.method,
			url: activeTabContent.url,
			headers: {},
			query: {},
			payload: activeTabContent.payload || ''
		};

		// Process headers
		if (activeTabContent.parsedHeaders) {
			activeTabContent.parsedHeaders.forEach((header) => {
				if (header.enabled && header.key && header.value) {
					requestData.headers[header.key] = header.value;
				}
			});
		}

		// Process query parameters
		if (activeTabContent.parsedParams) {
			activeTabContent.parsedParams.forEach((param) => {
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
	function parseHeaders(
		headersJson?: string
	): Array<{ key: string; value: string; description: string; enabled: boolean }> {
		if (!headersJson) {
			return [{ key: '', value: '', description: '', enabled: true }];
		}

		try {
			const headers = JSON.parse(headersJson);
			return Array.isArray(headers)
				? headers
				: [{ key: '', value: '', description: '', enabled: true }];
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
	} {
		if (!metadataJson) {
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }]
			};
		}

		try {
			const metadata = JSON.parse(metadataJson);
			return {
				parsedParams: metadata.params || [{ key: '', value: '', description: '', enabled: true }]
			};
		} catch (e) {
			console.error('Failed to parse metadata JSON:', e);
			return {
				parsedParams: [{ key: '', value: '', description: '', enabled: true }]
			};
		}
	}

	// Initialize activeTabContent if it's null
	$: if (activeTabContent === undefined || activeTabContent === null) {
		resetActiveTabContent();
	}

	// Sync activeTabContent when activeTabId or tabs change
	$: if (activeTabId && tabs.length > 0) {
		const currentTab = tabs.find(tab => tab.id === activeTabId);
		if (currentTab) {
			// Check if we need to initialize or update activeTabContent for this tab
			if (!activeTabContent || activeTabContent.id !== activeTabId) {
				// Create default empty strings for JSON fields
				const emptyHeadersJson = JSON.stringify([
					{ key: '', value: '', description: '', enabled: true }
				]);
				const emptyConfigJson = JSON.stringify({
					auth: { type: 'none', config: {} },
					settings: {}
				});
				const emptyMetadataJson = JSON.stringify({
					params: [{ key: '', value: '', description: '', enabled: true }]
				});

				// Use tab content if available, otherwise defaults
				const tabContent = currentTab.content;
				
				activeTabContent = {
					id: currentTab.id,
					name: currentTab.name || 'New Request',
					method: currentTab.method || 'GET',
					url: currentTab.url || '',
					activeSection: tabContent?.activeSection || 'params',
					protocol: 'http',
					itemType: currentTab.itemType || 'request',
					folder: currentTab.folder,

					// Use proper TabContent structure for data
					headers: JSON.stringify(tabContent?.headers || [{ key: '', value: '', description: '', enabled: true }]),
					config: JSON.stringify({
						auth: tabContent?.auth || { type: 'none', config: {} },
						settings: tabContent?.settings || {}
					}),
					metadata: JSON.stringify({
						params: tabContent?.params || [{ key: '', value: '', description: '', enabled: true }]
					}),
					payload: tabContent?.body?.content || '',

					// Parse JSON fields
					parsedHeaders: tabContent?.headers || [{ key: '', value: '', description: '', enabled: true }],
					...parseConfig(JSON.stringify({
						auth: tabContent?.auth || { type: 'none', config: {} },
						settings: tabContent?.settings || {}
					})),
					...parseMetadata(JSON.stringify({
						params: tabContent?.params || [{ key: '', value: '', description: '', enabled: true }]
					}))
				};
			}
		}
	}
	let isEditingTitle = false;
	let editTitleValue = '';
	let titleInputRef: HTMLInputElement;

	function startEditTitle() {
		const tab = tabs.find(t => t.id === activeTabId);
		if (tab) {
			editTitleValue = tab.name || '';
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
		
		const oldName = tabs[tabIndex].name;
		const newName = editTitleValue.trim();
		
		if (newName && newName !== oldName) {
			// Update local tab array immediately for quick UI update
			tabs[tabIndex].name = newName;
			if (activeTabContent) {
				activeTabContent.name = newName;
				if (activeTabContent.itemType === 'folder' && activeTabContent.folder) {
					activeTabContent.folder.name = newName;
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
				if (activeTabContent.id && !activeTabContent.id.startsWith('tab-')) {
					// Mark tab as unsaved so user knows they need to click Save button,
					// or we can auto save via API. Replay list auto-save title:
					try {
						await replayApi.updateReplay($selectedWorkspace?.id as string, $selectedProject?.id as string, activeTabContent.id, { name: newName });
						dispatch('updated', { id: activeTabContent.id, name: newName });
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
				<i class="fas {activeTabContent?.itemType === 'folder' ? 'fa-folder-open text-orange-500' : 'fa-file-alt text-blue-500'} text-xl"></i>
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
						{tabs.find((t) => t.id === activeTabId)?.name || 'New Request'}
						<i class="fas fa-pencil-alt text-xs opacity-0 group-hover:opacity-50 transition-opacity" on:click={startEditTitle}></i>
					</h1>
				{/if}
			</div>
			<div class="flex items-center space-x-2">
				{#if activeTabContent?.itemType !== 'folder'}
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

		{#if activeTabContent?.itemType === 'folder'}
			<!-- Folder Overview Mode -->
			<div class="flex-grow flex flex-col min-h-0 w-full">
				<FolderOverviewTab
					folder={activeTabContent.folder}
				/>
			</div>
		{:else}
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
			<HeadersTab
				headers={activeTabContent?.parsedHeaders}
				on:headersChange={handleHeadersChange}
			/>
		{:else if activeTabContent.activeSection === 'body'}
			<ReplayBody 
				payload={activeTabContent?.payload} 
				on:change={handleBodyContentChange}
			/>
		{:else if activeTabContent.activeSection === 'settings'}
			<SettingsTab
				settings={activeTabContent?.parsedSettings}
				on:settingsChange={handleSettingsChange}
			/>
		{/if}
		{/if} <!-- End of request vs folder branch -->
	</main>

	{#if activeTabContent?.itemType === 'request'}
		<ReplayResponseFooter
			bind:isExpanded={isFooterExpanded}
			{executionResult}
			on:toggleExpand={handleFooterToggleExpand}
			on:showHistory={handleFooterShowHistory}
		/>
	{/if}
</div>
