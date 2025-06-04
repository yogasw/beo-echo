import type { Tab, TabContent } from '$lib/components/replay/types';

const REPLAY_EDITOR_STORAGE_PREFIX = 'replay-editor-state';
const WORKSPACE_PROJECT_REGISTRY_KEY = 'replay-editor-workspaces';

export interface ReplayEditorState {
	tabs: Tab[];
	activeTabId: string;
	activeTabContent: {
		method: string;
		url: string;
		activeSection: string;
	};
	activeView: 'list' | 'editor' | 'execution' | 'logs';
	selectedReplayId?: string;
	projectId: string;
	workspaceId: string;
}

// Registry to track which projects exist for cleanup
interface WorkspaceProjectRegistry {
	[workspaceId: string]: string[]; // workspace -> project IDs
}

/**
 * Generate storage key for specific project
 */
function getProjectStorageKey(workspaceId: string, projectId: string): string {
	return `${REPLAY_EDITOR_STORAGE_PREFIX}-${workspaceId}-${projectId}`;
}

/**
 * Get or create workspace-project registry
 */
function getWorkspaceProjectRegistry(): WorkspaceProjectRegistry {
	if (typeof window === 'undefined') return {};
	
	try {
		const stored = localStorage.getItem(WORKSPACE_PROJECT_REGISTRY_KEY);
		return stored ? JSON.parse(stored) : {};
	} catch (error) {
		console.error('Failed to load workspace-project registry:', error);
		return {};
	}
}

/**
 * Update workspace-project registry
 */
function updateWorkspaceProjectRegistry(workspaceId: string, projectId: string): void {
	if (typeof window === 'undefined') return;
	
	try {
		const registry = getWorkspaceProjectRegistry();
		
		if (!registry[workspaceId]) {
			registry[workspaceId] = [];
		}
		
		if (!registry[workspaceId].includes(projectId)) {
			registry[workspaceId].push(projectId);
		}
		
		localStorage.setItem(WORKSPACE_PROJECT_REGISTRY_KEY, JSON.stringify(registry));
	} catch (error) {
		console.error('Failed to update workspace-project registry:', error);
	}
}

/**
 * Remove project from registry and clean up its storage
 */
function removeProjectFromRegistry(workspaceId: string, projectId: string): void {
	if (typeof window === 'undefined') return;
	
	try {
		const registry = getWorkspaceProjectRegistry();
		
		if (registry[workspaceId]) {
			registry[workspaceId] = registry[workspaceId].filter(id => id !== projectId);
			
			// Remove workspace entry if no projects left
			if (registry[workspaceId].length === 0) {
				delete registry[workspaceId];
			}
		}
		
		// Clean up the project's storage
		const projectKey = getProjectStorageKey(workspaceId, projectId);
		localStorage.removeItem(projectKey);
		
		// Update registry
		localStorage.setItem(WORKSPACE_PROJECT_REGISTRY_KEY, JSON.stringify(registry));
	} catch (error) {
		console.error('Failed to remove project from registry:', error);
	}
}

/**
 * Clean up all projects for a workspace
 */
function cleanupWorkspaceProjects(workspaceId: string): void {
	if (typeof window === 'undefined') return;
	
	try {
		const registry = getWorkspaceProjectRegistry();
		
		if (registry[workspaceId]) {
			// Remove all project storage for this workspace
			registry[workspaceId].forEach(projectId => {
				const projectKey = getProjectStorageKey(workspaceId, projectId);
				localStorage.removeItem(projectKey);
			});
			
			// Remove workspace from registry
			delete registry[workspaceId];
			localStorage.setItem(WORKSPACE_PROJECT_REGISTRY_KEY, JSON.stringify(registry));
		}
	} catch (error) {
		console.error('Failed to cleanup workspace projects:', error);
	}
}

/**
 * Create default tab content
 */
export function createDefaultTabContent(): TabContent {
	return {
		method: 'GET',
		url: '',
		params: [],
		headers: [
			{ key: 'User-Agent', value: 'BeoEcho/1.0', description: 'User agent header', enabled: true }
		],
		body: {
			type: 'none',
			content: '',
			formData: [],
			urlEncoded: []
		},
		auth: {
			type: 'none',
			config: {}
		},
		scripts: {
			preRequestScript: '',
			testScript: ''
		},
		settings: {
			timeout: 5000,
			followRedirects: true,
			maxRedirects: 5,
			verifySsl: true,
			ignoreSslErrors: false,
			encoding: 'utf-8',
			sendCookies: true,
			storeCookies: true,
			keepAlive: true,
			userAgent: 'BeoEcho/1.0',
			retryOnFailure: false,
			retryCount: 3,
			retryDelay: 1000
		},
		activeSection: 'params'
	};
}

/**
 * Create default tab with full content
 */
export function createDefaultTab(id?: string): Tab {
	const tabId = id || `tab-${Date.now()}`;
	
	return {
		id: tabId,
		name: 'New Request',
		method: 'GET',
		url: '',
		isUnsaved: true,
		content: createDefaultTabContent()
	};
}

/**
 * Get replay editor state from localStorage for specific project
 */
export function getReplayEditorState(workspaceId: string, projectId: string): ReplayEditorState | null {
	if (typeof window === 'undefined') return null;
	
	try {
		const storageKey = getProjectStorageKey(workspaceId, projectId);
		const stored = localStorage.getItem(storageKey);
		
		if (!stored) {
			return null;
		}
		
		const state = JSON.parse(stored);
		
		// Validate the state structure
		if (!state.tabs || !Array.isArray(state.tabs) || !state.activeTabId) {
			console.warn('❌ Invalid state structure:', state);
			return null;
		}
		
		// Validate and fix each tab's content
		state.tabs = state.tabs.map((tab: Tab) => validateAndFixTabContent(tab));
		
		// Ensure projectId and workspaceId are set
		state.projectId = projectId;
		state.workspaceId = workspaceId;
		
		return state;
	} catch (error) {
		console.warn('❌ Failed to load replay editor state from localStorage:', error);
		return null;
	}
}

/**
 * Save replay editor state to localStorage for specific project
 */
export function setReplayEditorState(state: ReplayEditorState): void {
	if (typeof window === 'undefined') return;
	
	try {
		const storageKey = getProjectStorageKey(state.workspaceId, state.projectId);
		const serialized = JSON.stringify(state);
		localStorage.setItem(storageKey, serialized);
		
		// Update registry to track this project
		updateWorkspaceProjectRegistry(state.workspaceId, state.projectId);
	} catch (error) {
		console.warn('❌ Failed to save replay editor state to localStorage:', error);
	}
}

/**
 * Clear replay editor state from localStorage for specific project
 */
export function clearReplayEditorState(workspaceId: string, projectId: string): void {
	if (typeof window === 'undefined') return;
	
	try {
		removeProjectFromRegistry(workspaceId, projectId);
	} catch (error) {
		console.warn('Failed to clear replay editor state from localStorage:', error);
	}
}

/**
 * Create default state for new replay editor
 */
export function createDefaultReplayEditorState(workspaceId: string, projectId: string): ReplayEditorState {
	const defaultTab = createDefaultTab();
	
	return {
		tabs: [defaultTab],
		activeTabId: defaultTab.id,
		activeTabContent: {
			method: 'GET',
			url: '',
			activeSection: 'params'
		},
		activeView: 'list',
		projectId,
		workspaceId
	};
}

/**
 * Update specific tab content in the stored state
 */
export function updateTabContentInStorage(workspaceId: string, projectId: string, tabId: string, content: any): void {
	const currentState = getReplayEditorState(workspaceId, projectId);
	if (!currentState) return;
	
	const tabIndex = currentState.tabs.findIndex(tab => tab.id === tabId);
	if (tabIndex === -1) return;
	
	const currentTab = currentState.tabs[tabIndex];
	
	// Ensure tab has content structure
	if (!currentTab.content) {
		currentTab.content = createDefaultTabContent();
	}
	
	// Handle specific field mappings
	const updatedContent = { ...currentTab.content };
	
	// Map payload to body.content
	if (content.payload !== undefined) {
		updatedContent.body = {
			...updatedContent.body,
			content: content.payload
		};
	}
	
	// Map other fields directly
	if (content.method !== undefined) updatedContent.method = content.method;
	if (content.url !== undefined) updatedContent.url = content.url;
	if (content.activeSection !== undefined) updatedContent.activeSection = content.activeSection;
	if (content.params !== undefined) updatedContent.params = content.params;
	if (content.headers !== undefined) updatedContent.headers = content.headers;
	if (content.auth !== undefined) updatedContent.auth = content.auth;
	if (content.settings !== undefined) updatedContent.settings = content.settings;
	
	currentTab.content = updatedContent;
	
	// Update basic tab properties if they changed
	if (content.method) currentTab.method = content.method;
	if (content.url) currentTab.url = content.url;
	
	// Mark as unsaved if content changed
	currentTab.isUnsaved = true;
	
	setReplayEditorState(currentState);
}

/**
 * Get tab content by ID
 */
export function getTabContentFromStorage(workspaceId: string, projectId: string, tabId: string): TabContent | null {
	const currentState = getReplayEditorState(workspaceId, projectId);
	if (!currentState) return null;
	
	const tab = currentState.tabs.find(tab => tab.id === tabId);
	if (!tab || !tab.content) return null;
	
	return tab.content;
}

/**
 * Update active section for a specific tab
 */
export function updateActiveSection(workspaceId: string, projectId: string, tabId: string, activeSection: string): void {
	const currentState = getReplayEditorState(workspaceId, projectId);
	if (!currentState) return;
	
	const tab = currentState.tabs.find(tab => tab.id === tabId);
	if (!tab || !tab.content) return;
	
	tab.content.activeSection = activeSection;
	
	// Also update the global active tab content if this is the active tab
	if (tabId === currentState.activeTabId) {
		currentState.activeTabContent.activeSection = activeSection;
	}
	
	setReplayEditorState(currentState);
}

/**
 * Add new tab to stored state
 */
export function addTabToStorage(workspaceId: string, projectId: string, tab: Tab): void {
	const currentState = getReplayEditorState(workspaceId, projectId);
	if (!currentState) return;
	
	// Ensure tab has content
	if (!tab.content) {
		tab.content = createDefaultTabContent();
	}
	
	currentState.tabs.push(tab);
	currentState.activeTabId = tab.id;
	
	setReplayEditorState(currentState);
}

/**
 * Remove tab from stored state
 */
export function removeTabFromStorage(workspaceId: string, projectId: string, tabId: string): void {
	const currentState = getReplayEditorState(workspaceId, projectId);
	if (!currentState) return;
	
	const updatedTabs = currentState.tabs.filter(tab => tab.id !== tabId);
	
	// If we removed the active tab, set a new active tab
	if (currentState.activeTabId === tabId && updatedTabs.length > 0) {
		currentState.activeTabId = updatedTabs[0].id;
	}
	
	// If no tabs left, create a default tab
	if (updatedTabs.length === 0) {
		const defaultTab = createDefaultTab();
		updatedTabs.push(defaultTab);
		currentState.activeTabId = defaultTab.id;
	}
	
	currentState.tabs = updatedTabs;
	setReplayEditorState(currentState);
}

/**
 * Validate and fix tab content structure
 */
export function validateAndFixTabContent(tab: Tab): Tab {
	if (!tab.content) {
		tab.content = createDefaultTabContent();
	} else {
		// Ensure all required properties exist
		const defaultContent = createDefaultTabContent();
		tab.content = {
			...defaultContent,
			...tab.content,
			// Ensure nested objects are properly merged
			body: {
				...defaultContent.body,
				...tab.content.body
			},
			auth: {
				...defaultContent.auth,
				...tab.content.auth
			},
			scripts: {
				...defaultContent.scripts,
				...tab.content.scripts
			},
			settings: {
				...defaultContent.settings,
				...tab.content.settings
			}
		};
	}
	
	return tab;
}

/**
 * Public cleanup utilities
 */

/**
 * Clean up storage when a project is deleted
 */
export function cleanupProjectStorage(workspaceId: string, projectId: string): void {
	removeProjectFromRegistry(workspaceId, projectId);
}

/**
 * Clean up storage when switching workspaces
 */
export function cleanupWorkspaceStorage(workspaceId: string): void {
	cleanupWorkspaceProjects(workspaceId);
}

/**
 * Get all stored project IDs for a workspace (useful for debugging or migration)
 */
export function getStoredProjectIds(workspaceId: string): string[] {
	const registry = getWorkspaceProjectRegistry();
	return registry[workspaceId] || [];
}

/**
 * Get total number of stored projects across all workspaces
 */
export function getStorageInfo(): { totalProjects: number; workspaces: string[] } {
	const registry = getWorkspaceProjectRegistry();
	const workspaces = Object.keys(registry);
	const totalProjects = workspaces.reduce((total, workspaceId) => {
		return total + (registry[workspaceId]?.length || 0);
	}, 0);
	
	return {
		totalProjects,
		workspaces
	};
}
