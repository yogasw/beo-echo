<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replays, selectedReplay, replayActions } from '$lib/stores/replay';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	import { getReplayPanelWidth, setReplayPanelWidth } from '$lib/utils/localStorage';

	import ReplayList from './ReplayList.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import ReplayEditor from './ReplayEditor.svelte';

	let isLoading = true;
	let error: string | null = null;
	let activeView: 'list' | 'editor' | 'execution' | 'logs' = 'list';
	let isExecuting = false;
	let executionResult: any = null;

	// Panel width
	let panelWidth: number; // Initialized in onMount
	let isResizing = false;
	let startX = 0;
	let startWidth = 0;

	// State for ReplayEditor, to be passed as props
	let editorTabs: Tab[] = [
		{
			id: 'tab-1',
			name: 'New Request',
			method: 'GET',
			url: '',
			isUnsaved: true
		}
	];
	let editorActiveTabId = 'tab-1';
	let editorActiveTabContent = {
		method: 'GET',
		url: '',
		activeSection: 'params'
	};

	// Handlers for events from ReplayEditor
	function handleTabsChange(event: CustomEvent) {
		editorTabs = event.detail.tabs;
		editorActiveTabId = event.detail.activeTabId;
		editorActiveTabContent = event.detail.activeTabContent;
	}

	function handleActiveSectionChange(event: CustomEvent) {
		if (editorActiveTabContent) {
			editorActiveTabContent.activeSection = event.detail.activeSection;
		}
	}

	function handleTabContentChange(event: CustomEvent) {
		const activeTab = editorTabs.find(tab => tab.id === editorActiveTabId);
		if (activeTab) {
			activeTab.method = event.detail.method;
			activeTab.url = event.detail.url;
			// Potentially mark as unsaved, etc.
		}
		if (editorActiveTabContent) {
			editorActiveTabContent = {...editorActiveTabContent, ...event.detail};
		}
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
		// Reset editor state for a new replay
		editorTabs = [
			{
				id: `tab-${Date.now()}`,
				name: 'New Request',
				method: 'GET',
				url: '',
				isUnsaved: true
			}
		];
		editorActiveTabId = editorTabs[0].id;
		editorActiveTabContent = {
			method: 'GET',
			url: '',
			activeSection: 'params'
		};
	}

	function handleEditReplay(event: CustomEvent) {
		const replay = event.detail;
		selectedReplay.set(replay);
		activeView = 'editor';
		// Populate editor state from the selected replay
		const replayData = replay; // Assuming replay has the necessary data
		editorTabs = [
			{
				id: replayData.id || `tab-${Date.now()}`,
				name: replayData.name || 'Edit Request',
				method: replayData.request?.method || 'GET',
				url: replayData.request?.url || '',
				isUnsaved: false // Or determine based on actual state
			}
		];
		editorActiveTabId = editorTabs[0].id;
		editorActiveTabContent = {
			method: replayData.request?.method || 'GET',
			url: replayData.request?.url || '',
			activeSection: 'params' // Or restore last active section for this replay
		};
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
	}

	function handleReplayCreated() {
		activeView = 'list';
		loadReplays(); // Refresh the list
	}

	function handleReplayUpdated() {
		activeView = 'list';
		loadReplays(); // Refresh the list
	}

	async function executeReplay(replayData: any) {
		if (!$selectedWorkspace || !$selectedProject) {
			toast.error('No workspace or project selected');
			return;
		}

		try {
			isExecuting = true;
			executionResult = null;
			
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
			toast.success('Request executed successfully');
			
			// You can optionally update UI to show the result or navigate to a result view
			activeView = 'execution';
		} catch (err: any) {
			toast.error(err.message || 'Failed to execute request');
			executionResult = {
				error: err.message || 'An unknown error occurred'
			};
		} finally {
			isExecuting = false;
		}
	}

	function handleSendRequest(event: CustomEvent) {
		const requestData = event.detail;
		executeReplay(requestData);
	}

	onMount(() => {
		if ($selectedWorkspace && $selectedProject) {
			loadReplays();
		}
		panelWidth = getReplayPanelWidth();
		// document.addEventListener('click', handleClickOutside); // If needed for closing menus
	});

	onDestroy(() => {
		replayActions.clearAll();
		// document.removeEventListener('click', handleClickOutside); // If added in onMount
		if (isResizing) { // Clean up listeners if component is destroyed while resizing
			document.removeEventListener('mousemove', handleResize);
			document.removeEventListener('mouseup', stopResize);
		}
	});

	// function handleClickOutside(event: MouseEvent) {
	// 	// Logic to close any dropdowns or modals if necessary
	// }

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

		// Constrain between 20% and 60% (adjust as needed)
		panelWidth = Math.min(Math.max(newWidth, 15), 70);
	}

	function stopResize() {
		isResizing = false;
		document.removeEventListener('mousemove', handleResize);
		document.removeEventListener('mouseup', stopResize);
		document.body.style.cursor = '';
		document.body.style.userSelect = '';

		setReplayPanelWidth(panelWidth);
	}
</script>

<div class="flex flex-col h-full bg-gray-50 dark:bg-gray-900">
	<!-- Main Content Area -->
	<div class="flex-1 flex p-4 h-full overflow-hidden bg-gray-50 dark:bg-gray-900"> 
		<!-- Left: Replay List (Resizable) -->
		<div 
			class="flex flex-col h-full space-y-4 relative mr-4"  
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
						on:add={handleCreateNew}
						on:edit={handleEditReplay}
						on:execute={handleExecuteReplay}
						on:logs={handleViewLogs}
						on:refresh={loadReplays}
					/>
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
					bind:activeTabContent={editorActiveTabContent}
					isExecuting={isExecuting}
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
