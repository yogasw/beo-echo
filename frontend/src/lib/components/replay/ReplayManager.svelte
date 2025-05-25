<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replays, selectedReplay, replayActions } from '$lib/stores/replay';
	import { toast } from '$lib/stores/toast';
	import { replayApi } from '$lib/api/replayApi';
	
	import ReplayList from './ReplayList.svelte';
	import ReplayEditor from './ReplayEditor.svelte';
	import ReplayExecution from './ReplayExecution.svelte';
	import ReplayLogs from './ReplayLogs.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	let isLoading = true;
	let error: string | null = null;
	let activeView: 'list' | 'editor' | 'execution' | 'logs' = 'list';

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
	}

	function handleEditReplay(replay: any) {
		selectedReplay.set(replay.detail);
		activeView = 'editor';
	}

	function handleExecuteReplay(replay: any) {
		selectedReplay.set(replay.detail);
		activeView = 'execution';
	}

	function handleViewLogs(replay: any) {
		if (replay?.detail) {
			selectedReplay.set(replay.detail);
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
	<!-- Header -->
	<div class="theme-bg-secondary border-b theme-border px-4 py-3">
		<div class="flex items-center justify-between">
			<div class="flex items-center space-x-4">
				{#if activeView !== 'list'}
					<button
						on:click={handleBackToList}
						class="flex items-center text-sm theme-text-secondary hover:theme-text-primary transition-colors"
					>
						<i class="fas fa-arrow-left mr-2"></i>
						Back to Replays
					</button>
				{/if}
				
				<div class="flex items-center space-x-2">
					<i class="fas fa-play-circle text-blue-400"></i>
					<h1 class="text-lg font-semibold theme-text-primary">
						{#if activeView === 'list'}
							API Replays
						{:else if activeView === 'editor'}
							{$selectedReplay ? 'Edit Replay' : 'Create Replay'}
						{:else if activeView === 'execution'}
							Execute Replay
						{:else if activeView === 'logs'}
							Execution Logs
						{/if}
					</h1>
				</div>
			</div>

			<div class="flex items-center space-x-2">
				{#if activeView === 'list'}
					<button
						on:click={handleCreateNew}
						class="px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors"
					>
						<i class="fas fa-plus mr-1"></i>
						New Replay
					</button>
				{/if}
			</div>
		</div>
	</div>

	<!-- Content Area -->
	<div class="flex-1 overflow-hidden">
		{#if !$selectedWorkspace || !$selectedProject}
			<div class="flex items-center justify-center h-full">
				<div class="text-center theme-text-secondary">
					<i class="fas fa-project-diagram text-4xl mb-4 opacity-50"></i>
					<p>Please select a workspace and project to manage replays</p>
				</div>
			</div>
		{:else if isLoading && activeView === 'list'}
			<div class="p-4">
				<SkeletonLoader type="list" count={5} />
			</div>
		{:else if error && activeView === 'list'}
			<div class="p-4">
				<ErrorDisplay 
					message={error} 
					type="error" 
					retryable={true}
					onRetry={loadReplays}
				/>
			</div>
		{:else if activeView === 'list'}
			<ReplayList 
				on:edit={handleEditReplay}
				on:execute={handleExecuteReplay}
				on:logs={handleViewLogs}
				on:refresh={loadReplays}
			/>
		{:else if activeView === 'editor'}
			<ReplayEditor 
				replay={$selectedReplay}
				on:created={handleReplayCreated}
				on:updated={handleReplayUpdated}
				on:cancel={handleBackToList}
			/>
		{:else if activeView === 'execution' && $selectedReplay}
			<ReplayExecution 
				replay={$selectedReplay}
				on:close={handleBackToList}
				on:executed={() => {/* Handle execution completion if needed */}}
			/>
		{:else if activeView === 'logs' && $selectedReplay}
			<ReplayLogs 
				replay={$selectedReplay}
				on:close={handleBackToList}
			/>
		{/if}
	</div>
</div>
