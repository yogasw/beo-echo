<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replays, selectedReplay, replayActions, replayLoading } from '$lib/stores/replay';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import { getReplayPanelWidth, setReplayPanelWidth } from '$lib/utils/localStorage';
	import { 
		getReplayEditorState, 
		setReplayEditorState, 
		clearReplayEditorState,
		createDefaultReplayEditorState,
		createDefaultTab,
		createDefaultTabContent,
		updateTabContentInStorage,
		updateActiveSection
	} from '$lib/utils/replayEditorStorage';

	import ReplayList from './ReplayList.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import ReplayEditor from './ReplayEditor.svelte';
	import type { Tab } from './types';
	import type { ExecuteReplayResponse } from '$lib/types/Replay';

	let isLoading = true;
	let error: string | null = null;
	let activeView: 'list' | 'editor' | 'execution' | 'logs' = 'list';
	let executionResult: ExecuteReplayResponse | null = null;

	// Panel width
	let panelWidth: number; // Initialized in onMount
	let isResizing = false;
	let startX = 0;
	let startWidth = 0;

	// Collapse/Expand state
	let isPanelCollapsed = false;
	let savedPanelWidth = 33; // Store the width before collapsing
	const collapsedWidth = 4; // Width when collapsed (just enough for the toggle button)
	let isAnimating = false; // Flag to enable animation only for toggle button

	// State for ReplayEditor - will be loaded from localStorage
	let editorTabs: Tab[] = [];
	let editorActiveTabId = '';
	let editorActiveTabContent = {
		method: 'GET',
		url: '',
		activeSection: 'params'
	};

	// Flag to track if we should save state changes
	let shouldSaveState = false; // Start as false, enable after load

	// Load editor state from localStorage
	function loadEditorState() {
		// Ensure we have workspace and project IDs before loading state
		if (!$selectedWorkspace?.id || !$selectedProject?.id) {
			console.log('⚠️ Missing workspace or project ID, using default state');
			const defaultState = createDefaultReplayEditorState('default', 'default');
			editorTabs = [...defaultState.tabs];
			editorActiveTabId = defaultState.activeTabId;
			editorActiveTabContent = { ...defaultState.activeTabContent };
			activeView = defaultState.activeView;
			return;
		}

		console.log('📂 Loading editor state from localStorage for workspace:', $selectedWorkspace.id, 'project:', $selectedProject.id);
		
		const savedState = getReplayEditorState($selectedWorkspace.id, $selectedProject.id);
		
		if (savedState && savedState.tabs && savedState.tabs.length > 0) {
			console.log('✅ Found saved state:', savedState);
			
			// Restore all state
			editorTabs = [...savedState.tabs]; // Clone to ensure reactivity
			editorActiveTabId = savedState.activeTabId;
			editorActiveTabContent = { ...savedState.activeTabContent };
			activeView = savedState.activeView;
			
			// Validate that activeTabId exists in tabs
			const activeTabExists = editorTabs.some(tab => tab.id === editorActiveTabId);
			if (!activeTabExists && editorTabs.length > 0) {
				console.warn('⚠️ Active tab not found, using first tab');
				editorActiveTabId = editorTabs[0].id;
				editorActiveTabContent = {
					method: editorTabs[0].method,
					url: editorTabs[0].url,
					activeSection: 'params'
				};
			}
			
			console.log('✅ Editor state restored successfully');
		} else {
			console.log('📝 No saved state found, using default state');
			
			// Use default state if no saved state
			const defaultState = createDefaultReplayEditorState($selectedWorkspace.id, $selectedProject.id);
			editorTabs = [...defaultState.tabs];
			editorActiveTabId = defaultState.activeTabId;
			editorActiveTabContent = { ...defaultState.activeTabContent };
			activeView = defaultState.activeView;
		}
		
		// Log final state
		console.log('📊 Final loaded state:', {
			tabs: editorTabs.length,
			activeTabId: editorActiveTabId,
			activeView: activeView
		});
	}

	// Save current editor state to localStorage
	function saveEditorState() {
		if (!shouldSaveState) {
			console.log('⏸️ Auto-save disabled, skipping save');
			return;
		}
		
		// Don't save if we don't have valid state
		if (!editorTabs || editorTabs.length === 0 || !editorActiveTabId) {
			console.log('⚠️ Invalid state, skipping save');
			return;
		}

		// Ensure we have workspace and project IDs
		if (!$selectedWorkspace?.id || !$selectedProject?.id) {
			console.log('⚠️ Missing workspace or project ID, skipping save');
			return;
		}
		
		const state = {
			tabs: editorTabs,
			activeTabId: editorActiveTabId,
			activeTabContent: editorActiveTabContent,
			activeView: activeView,
			selectedReplayId: $selectedReplay?.id,
			projectId: $selectedProject.id,
			workspaceId: $selectedWorkspace.id
		};
		
		console.log('💾 Saving editor state:', {
			tabsCount: state.tabs.length,
			activeTabId: state.activeTabId,
			activeView: state.activeView
		});
		
		setReplayEditorState(state);
	}

	// Auto-save with debouncing
	let saveTimeout: NodeJS.Timeout | null = null;
	
	function debouncedSave() {
		if (saveTimeout) {
			clearTimeout(saveTimeout);
		}
		
		saveTimeout = setTimeout(() => {
			saveEditorState();
		}, 300); // Save after 300ms of no changes
	}

	// Watch for changes and trigger save
	$: if (shouldSaveState && editorTabs && editorTabs.length > 0) {
		console.log('🔄 Tabs changed, scheduling save...');
		debouncedSave();
	}

	$: if (shouldSaveState && editorActiveTabId) {
		console.log('🔄 Active tab changed, scheduling save...');
		debouncedSave();
	}

	$: if (shouldSaveState && activeView) {
		console.log('🔄 Active view changed, scheduling save...');
		debouncedSave();
	}

	// Handlers for events from ReplayEditor
	function handleTabsChange(event: CustomEvent) {
		console.log('🔄 Tabs changed from editor:', event.detail);
		
		// Update basic tab properties but preserve content structure
		const newTabs = event.detail.tabs.map(newTab => {
			// Find existing tab to preserve content
			const existingTab = editorTabs.find(tab => tab.id === newTab.id);
			
			return {
				...newTab,
				// Preserve existing content or create default if new tab
				content: existingTab?.content || createDefaultTabContent()
			};
		});
		
		editorTabs = [...newTabs]; // Clone to ensure reactivity
		editorActiveTabId = event.detail.activeTabId;
		editorActiveTabContent = { ...event.detail.activeTabContent };
	}

	function handleActiveSectionChange(event: CustomEvent) {
		console.log('🔄 Active section changed:', event.detail);
		
		if (editorActiveTabContent) {
			editorActiveTabContent.activeSection = event.detail.activeSection;
			// Trigger reactivity
			editorActiveTabContent = { ...editorActiveTabContent };
		}
		
	// Also save to localStorage
	if ($selectedWorkspace?.id && $selectedProject?.id) {
		updateActiveSection($selectedWorkspace.id, $selectedProject.id, editorActiveTabId, event.detail.activeSection);
	}
	}

	function handleTabContentChange(event: CustomEvent) {
		// Update the active tab content
		if (editorActiveTabContent) {
			editorActiveTabContent = { ...editorActiveTabContent, ...event.detail };
		}
		
		// Update the corresponding tab and its content
		const activeTab = editorTabs.find(tab => tab.id === editorActiveTabId);
		if (activeTab) {
			// Update basic tab properties
			if (event.detail.method) activeTab.method = event.detail.method;
			if (event.detail.url) activeTab.url = event.detail.url;
			
			// Update tab content
			if (!activeTab.content) {
				activeTab.content = createDefaultTabContent();
			}
			
			// Handle specific field mappings for TabContent structure
			const updatedContent = { ...activeTab.content };
			
			// Map payload to body.content
			if (event.detail.payload !== undefined) {
				updatedContent.body = {
					...updatedContent.body,
					content: event.detail.payload
				};
			}
			
			// Map other fields directly
			if (event.detail.method) updatedContent.method = event.detail.method;
			if (event.detail.url) updatedContent.url = event.detail.url;
			if (event.detail.activeSection) updatedContent.activeSection = event.detail.activeSection;
			if (event.detail.params) updatedContent.params = event.detail.params;
			if (event.detail.headers) updatedContent.headers = event.detail.headers;
			if (event.detail.auth) updatedContent.auth = event.detail.auth;
			if (event.detail.settings) updatedContent.settings = event.detail.settings;
			
			activeTab.content = updatedContent;
			activeTab.isUnsaved = true; // Mark as unsaved when content changes
			
			// Trigger reactivity
			editorTabs = [...editorTabs];
					// Also save to localStorage immediately for important changes
		if ($selectedWorkspace?.id && $selectedProject?.id) {
			updateTabContentInStorage($selectedWorkspace.id, $selectedProject.id, editorActiveTabId, event.detail);
		}
		}
	}

	// Load replays when component mounts or project changes
	$: if ($selectedWorkspace && $selectedProject) {
		loadReplays();
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
		
		// Create a new tab with full content structure
		const newTab = createDefaultTab();
		
		// Add to existing tabs or replace if only one empty tab exists
		if (editorTabs.length === 1 && editorTabs[0].isUnsaved && !editorTabs[0].url) {
			editorTabs = [newTab];
		} else {
			editorTabs = [...editorTabs, newTab];
		}
		
		// Set active tab and content properly
		editorActiveTabId = newTab.id;
		editorActiveTabContent = newTab.content || {
			method: 'GET',
			url: '',
			activeSection: 'params'
		};
		
		// Save the new tab to localStorage
		if ($selectedWorkspace?.id && $selectedProject?.id) {
			setReplayEditorState({
				tabs: editorTabs,
				activeTabId: editorActiveTabId,
				activeTabContent: editorActiveTabContent,
				activeView: activeView,
				projectId: $selectedProject.id,
				workspaceId: $selectedWorkspace.id
			});
		}
	}

	function handleEditReplay(event: CustomEvent) {
		const replay = event.detail;
		selectedReplay.set(replay);
		activeView = 'editor';
		
		// Check if this replay is already open in a tab
		const existingTab = editorTabs.find(tab => 
			tab.id === replay.id || 
			(tab.method === replay.method && tab.url === replay.url && !tab.isUnsaved)
		);
		
		if (existingTab) {
			// Switch to existing tab
			editorActiveTabId = existingTab.id;
			editorActiveTabContent = {
				method: existingTab.method,
				url: existingTab.url,
				activeSection: existingTab.content?.activeSection || 'params'
			};
		} else {
			// Create new tab for this replay with full content
			const newTab: Tab = {
				id: replay.id || `tab-${Date.now()}`,
				name: replay.name || 'Edit Request',
				method: replay.method || 'GET',
				url: replay.url || '',
				isUnsaved: false,
				content: {
					...createDefaultTabContent(),
					method: replay.method || 'GET',
					url: replay.url || '',
					// TODO: Populate from replay data if available
					// params: replay.params || [],
					// headers: replay.headers || [],
					// body: replay.body || { type: 'none', content: '' },
				}
			};
			
			editorTabs = [...editorTabs, newTab];
			editorActiveTabId = newTab.id;
			editorActiveTabContent = {
				method: newTab.method,
				url: newTab.url,
				activeSection: 'params'
			};
		}
	}

	function handleExecuteReplay(event: CustomEvent) {
		const replay = event.detail;
		selectedReplay.set(replay);
		activeView = 'execution';
	}

	function handleViewLogs(event: CustomEvent) {
		const replay = event.detail;
		if (replay) {
			selectedReplay.set(replay);
		}
		activeView = 'logs';
	}

	function handleBackToList() {
		selectedReplay.set(null);
		activeView = 'list';
		// Keep tabs open so user can continue working if they go back to editor
	}

	// Function to manually clear all editor state
	function clearAllEditorState() {
		console.log('🗑️ Clearing all editor state...');
		
		if ($selectedWorkspace?.id && $selectedProject?.id) {
			clearReplayEditorState($selectedWorkspace.id, $selectedProject.id);
			const defaultState = createDefaultReplayEditorState($selectedWorkspace.id, $selectedProject.id);
			
			shouldSaveState = false;
			editorTabs = [...defaultState.tabs];
			editorActiveTabId = defaultState.activeTabId;
			editorActiveTabContent = { ...defaultState.activeTabContent };
			activeView = 'list';
			selectedReplay.set(null);
			
			setTimeout(() => {
				shouldSaveState = true;
			}, 100);
			
			toast.success('Editor state cleared');
			console.log('✅ Editor state cleared and reset to default');
		} else {
			console.log('⚠️ Cannot clear state - missing workspace or project ID');
			toast.error('Cannot clear state - no workspace or project selected');
		}
	}

	// Debug function to test localStorage
	function testLocalStorage() {
		console.log('🧪 Testing localStorage...');
		
		if (!$selectedWorkspace?.id || !$selectedProject?.id) {
			toast.error('Cannot test localStorage - no workspace or project selected');
			return;
		}
		
		const testState = {
			tabs: [{ id: 'test', name: 'Test Tab', method: 'GET', url: 'http://test.com', isUnsaved: false }],
			activeTabId: 'test',
			activeTabContent: { method: 'GET', url: 'http://test.com', activeSection: 'params' },
			activeView: 'editor' as const,
			projectId: $selectedProject.id,
			workspaceId: $selectedWorkspace.id
		};
		
		// Test save
		setReplayEditorState(testState);
		
		// Test load
		const loaded = getReplayEditorState($selectedWorkspace.id, $selectedProject.id);
		console.log('Test result - saved:', testState, 'loaded:', loaded);
		
		// Check if they match
		const matches = JSON.stringify(testState) === JSON.stringify(loaded);
		console.log('LocalStorage test:', matches ? '✅ PASSED' : '❌ FAILED');
		
		toast.success(`LocalStorage test: ${matches ? 'PASSED' : 'FAILED'}`);
	}

	function handleReplayCreated() {
		// Mark current tab as saved
		const activeTab = editorTabs.find(tab => tab.id === editorActiveTabId);
		if (activeTab) {
			activeTab.isUnsaved = false;
			editorTabs = [...editorTabs]; // Trigger reactivity
		}
		
		loadReplays(); // Refresh the list
		toast.success('Replay created successfully');
	}

	function handleReplayUpdated() {
		// Mark current tab as saved
		const activeTab = editorTabs.find(tab => tab.id === editorActiveTabId);
		if (activeTab) {
			activeTab.isUnsaved = false;
			editorTabs = [...editorTabs]; // Trigger reactivity
		}
		
		loadReplays(); // Refresh the list
		toast.success('Replay updated successfully');
	}

	async function executeReplay(replayData: any) {
		if (!$selectedWorkspace || !$selectedProject) {
			toast.error('No workspace or project selected');
			return;
		}

		try {
			replayActions.setLoading('execute', true);
			
			// Prepare request payload from editor data
			const payload = {
				protocol: 'http', // Default to http
				method: replayData.method || 'GET',
				url: replayData.url || '',
				headers: replayData.headers || {},
				body: replayData.body || '',
				query: replayData.query || {}
			};

			// Execute the replay request
			const result = await replayApi.executeReplayRequest(
				$selectedWorkspace.id, 
				$selectedProject.id, 
				payload
			);
			
			executionResult = result;
			console.log('Execution result:', executionResult);
			toast.success('Request executed successfully');
			
			// You can optionally update UI to show the result or navigate to a result view
			activeView = 'execution';
		} catch (err: any) {
			toast.error(err.message || 'Failed to execute request');
			executionResult = null
		} finally {
			replayActions.setLoading('execute', false);
		}
	}

	function handleSendRequest(event: CustomEvent) {
		const requestData = event.detail;
		executeReplay(requestData);
	}

	onMount(async () => {
		console.log('🚀 ReplayManager mounted');

		// Load panel width
		panelWidth = getReplayPanelWidth();

		// Check if panel was collapsed (width is at or near collapsed width)
		if (panelWidth <= collapsedWidth + 1) {
			isPanelCollapsed = true;
		} else {
			savedPanelWidth = panelWidth;
		}

		// Load editor state from localStorage FIRST
		loadEditorState();
		
		// Load replays if workspace and project are selected
		if ($selectedWorkspace && $selectedProject) {
			console.log('📋 Loading replays for workspace and project');
			await loadReplays();
		}
		
		// Add beforeunload listener to save state before page closes/refreshes
		const handleBeforeUnload = (event: BeforeUnloadEvent) => {
			console.log('📤 Page unloading, saving state...');
			saveEditorState();
		};
		
		window.addEventListener('beforeunload', handleBeforeUnload);
		
		// Enable auto-saving after everything is loaded
		setTimeout(() => {
			shouldSaveState = true;
			console.log('▶️ Auto-save enabled');
			
			// Force an immediate save to ensure current state is persisted
			debouncedSave();
		}, 500);
	});

	onDestroy(() => {
		console.log('🛑 ReplayManager destroyed');
		
		// Clear any pending save timeout
		if (saveTimeout) {
			clearTimeout(saveTimeout);
		}
		
		// Force save current state before destroy
		if (shouldSaveState && $selectedWorkspace?.id && $selectedProject?.id) {
			console.log('💾 Force saving state before destroy');
			// Use direct save, not debounced
			const state = {
				tabs: editorTabs,
				activeTabId: editorActiveTabId,
				activeTabContent: editorActiveTabContent,
				activeView: activeView,
				selectedReplayId: $selectedReplay?.id,
				projectId: $selectedProject.id,
				workspaceId: $selectedWorkspace.id
			};
			setReplayEditorState(state);
		}
		
		replayActions.clearAll();
		
		// Clean up resize listeners if component is destroyed while resizing
		if (isResizing) {
			document.removeEventListener('mousemove', handleResize);
			document.removeEventListener('mouseup', stopResize);
		}
		
		// Remove beforeunload listener
		const handleBeforeUnload = (event: BeforeUnloadEvent) => {
			saveEditorState();
		};
		window.removeEventListener('beforeunload', handleBeforeUnload);
	});

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
		// Assuming the ReplayManager's parent container takes full width for percentage calculation
		const containerWidth = (event.target as HTMLElement)?.closest('.flex-1.flex.p-4.h-full')?.clientWidth || window.innerWidth;
		const newWidth = startWidth + (deltaX / containerWidth) * 100;

		// Set minimum width threshold (slightly above collapsed to prevent accidental collapse while dragging)
		const minDragWidth = 15; // Minimum width when dragging (15%)

		// Constrain between minDragWidth and 70%
		panelWidth = Math.min(Math.max(newWidth, minDragWidth), 70);

		// Auto-collapse if dragged close to minimum
		if (panelWidth <= minDragWidth + 2) {
			// Snap to collapsed width
			panelWidth = collapsedWidth;
			isPanelCollapsed = true;
		} else {
			isPanelCollapsed = false;
		}
	}

	function stopResize() {
		isResizing = false;
		document.removeEventListener('mousemove', handleResize);
		document.removeEventListener('mouseup', stopResize);
		document.body.style.cursor = '';
		document.body.style.userSelect = '';

		setReplayPanelWidth(panelWidth);
	}

	function togglePanelCollapse(event?: CustomEvent) {
		// Enable animation for toggle button clicks
		isAnimating = true;

		if (isPanelCollapsed) {
			// Expand: restore to saved width
			panelWidth = savedPanelWidth;
			isPanelCollapsed = false;
		} else {
			// Collapse: save current width and set to minimal
			savedPanelWidth = panelWidth;
			panelWidth = collapsedWidth;
			isPanelCollapsed = true;
		}

		// Save the panel width to localStorage
		setReplayPanelWidth(panelWidth);

		// Disable animation after transition completes
		setTimeout(() => {
			isAnimating = false;
		}, 300);
	}
</script>

<div class="flex flex-col h-full bg-gray-50 dark:bg-gray-900">
	<!-- Main Content Area -->
	<div class="flex-1 flex p-4 h-full overflow-hidden bg-gray-50 dark:bg-gray-900"> 
		<!-- Left: Replay List (Resizable) -->
		<div
			class="flex flex-col h-full space-y-4 relative mr-4 {isAnimating ? 'panel-transition' : ''}"
			style="width: {panelWidth}%;"
		>
			<div class="flex-1 min-h-0 overflow-y-auto pr-1 hide-scrollbar"> 
				{#if !$selectedWorkspace || !$selectedProject}
					<div class="flex items-center justify-center h-full">
						<div class="text-center theme-text-secondary p-4">
							<i class="fas fa-project-diagram text-4xl mb-4 opacity-50"></i>
							<p>Please select a workspace and project to manage replays</p>
						</div>
					</div>
				{:else if isLoading}
					<div class="p-4">
						<SkeletonLoader type="list" count={5} />
					</div>
				{:else if error}
					<div class="p-4">
						<ErrorDisplay message={error} type="error" retryable={true} onRetry={loadReplays} />
					</div>
				{:else}
					<ReplayList
						{isPanelCollapsed}
						on:add={handleCreateNew}
						on:edit={handleEditReplay}
						on:execute={handleExecuteReplay}
						on:logs={handleViewLogs}
						on:refresh={loadReplays}
						on:toggleCollapse={togglePanelCollapse}
					/>
					
					<!-- Debug Panel -->
					<div class="p-4 border-t border-gray-200 dark:border-gray-700 space-y-2">
						<div class="text-xs text-gray-500 dark:text-gray-400 mb-2">
							<div>Tabs: {editorTabs.length}</div>
							<div>Active: {editorActiveTabId}</div>
							<div>View: {activeView}</div>
							<div>Auto-save: {shouldSaveState ? '✅' : '❌'}</div>
						</div>
						
						<button
							class="w-full bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded text-xs flex items-center justify-center"
							title="Test localStorage functionality"
							aria-label="Test localStorage"
							on:click={testLocalStorage}
						>
							<i class="fas fa-flask mr-2"></i>
							Test LocalStorage
						</button>
						
						<button
							class="w-full bg-gray-600 hover:bg-gray-700 text-white py-2 px-4 rounded text-xs flex items-center justify-center"
							title="Clear all editor tabs and reset state"
							aria-label="Clear editor state"
							on:click={clearAllEditorState}
						>
							<i class="fas fa-trash-alt mr-2"></i>
							Clear Editor State
						</button>
					</div>
				{/if}
			</div>
			<!-- Resizable handle -->
			<div
				class="absolute top-0 right-0 bottom-0 w-1.5 cursor-col-resize group z-10 hover:bg-blue-500/20 dark:hover:bg-blue-400/20 transition-colors duration-200"
				role="button"
				tabindex="0"
				on:mousedown|preventDefault={startResize}
				aria-label="Resize panel"
				title="Drag to resize panel"
			>
				<div class="w-full h-full bg-transparent group-hover:bg-blue-500/30 dark:group-hover:bg-blue-400/30 transition-colors duration-150"></div>
			</div>
		</div>

		<!-- Right: Replay Content -->
		<div
			class="flex-1 flex flex-col h-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden shadow-sm"
		>
			{#if activeView === 'editor' || activeView === 'execution' || activeView === 'logs'}
				<ReplayEditor 
					bind:tabs={editorTabs} 
					bind:activeTabId={editorActiveTabId} 
					executionResult={executionResult}
					on:tabschange={handleTabsChange}
					on:activeSectionChange={handleActiveSectionChange}
					on:tabContentChange={handleTabContentChange}
					on:back={handleBackToList}
					on:created={handleReplayCreated}
					on:updated={handleReplayUpdated}
					on:send={handleSendRequest}
				/>
			{:else if activeView === 'list' && ($selectedWorkspace && $selectedProject)}
				<div class="flex flex-col items-center justify-center h-full text-center p-8 text-gray-500 dark:text-gray-400">
					<div class="bg-gray-100 dark:bg-gray-700 rounded-full p-6 mb-4">
						<i class="fas fa-mouse-pointer text-4xl text-gray-400 dark:text-gray-500"></i>
					</div>
					<h3 class="text-lg font-medium mb-2 text-gray-700 dark:text-gray-300">Select a Replay</h3>
					<p class="text-sm max-w-xs leading-relaxed">
						Choose a replay from the list to view or edit, or click "New Replay" to create one.
					</p>
				</div>
			{:else}
				<div class="flex flex-col items-center justify-center h-full text-center p-8 text-gray-500 dark:text-gray-400">
					<div class="bg-gray-100 dark:bg-gray-700 rounded-full p-6 mb-4">
						<i class="fas fa-info-circle text-4xl text-gray-400 dark:text-gray-500"></i>
					</div>
					<h3 class="text-lg font-medium mb-2 text-gray-700 dark:text-gray-300">Get Started</h3>
					<p class="text-sm max-w-xs leading-relaxed">
						Select a workspace and project to manage your API replays.
					</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.panel-transition {
		transition: width 0.3s ease-in-out;
	}
</style>
