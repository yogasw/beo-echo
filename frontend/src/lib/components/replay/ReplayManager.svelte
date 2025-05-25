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
	import ReplayEditor from './ReplayEditor.svelte';
	import type { Tab } from './ReplayEditor.svelte';

	let isLoading = true;
	let error: string | null = null;
	let activeView: 'list' | 'editor' | 'execution' | 'logs' = 'list';

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
						on:add={handleCreateNew}
						on:edit={handleEditReplay}
						on:execute={handleExecuteReplay}
						on:logs={handleViewLogs}
						on:refresh={loadReplays}
					/>
				{/if}
			</div>
		</div>
		<!-- Center & Right: Replay Content -->
		<div
			class="col-span-2 flex flex-col h-full bg-gray-800 dark:bg-gray-800 border border-gray-700 dark:border-gray-700 rounded-lg overflow-hidden"
		>
			{#if activeView === 'editor' || activeView === 'execution' || activeView === 'logs'}
				<ReplayEditor 
					bind:tabs={editorTabs} 
					bind:activeTabId={editorActiveTabId} 
					bind:activeTabContent={editorActiveTabContent}
					on:tabschange={handleTabsChange}
					on:activeSectionChange={handleActiveSectionChange}
					on:tabContentChange={handleTabContentChange}
				/>
			{:else if activeView === 'list' && ($selectedWorkspace && $selectedProject)}
				<div class="flex flex-col items-center justify-center h-full text-center p-4 theme-text-secondary">
					<i class="fas fa-mouse-pointer text-5xl mb-4 opacity-50"></i>
					<p class="text-lg">Select a replay from the list to view or edit,</p>
					<p class="text-lg">or click "New Replay" to create one.</p>
				</div>
			{:else}
				<div class="flex flex-col items-center justify-center h-full text-center p-4 theme-text-secondary">
					<i class="fas fa-info-circle text-5xl mb-4 opacity-50"></i>
					<p class="text-lg">Select a workspace and project to get started.</p>
				</div>
			{/if}
		</div>
	</div>
</div>
