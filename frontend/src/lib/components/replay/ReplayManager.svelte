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
	<!-- Main Content Area -->
	<div class="flex-1 grid grid-cols-3 gap-4 p-4">
		<!-- Left: Replay List -->
		<div class="space-y-4">
			<button
				on:click={handleCreateNew}
				class="w-full px-3 py-1.5 text-sm bg-blue-600 hover:bg-blue-700 text-white rounded transition-colors"
			>
				<i class="fas fa-plus mr-1"></i>
				New Replay
			</button>

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
				<ErrorDisplay 
					message={error} 
					type="error" 
					retryable={true}
					onRetry={loadReplays}
				/>
			{:else}
				<ReplayList 
					on:edit={handleEditReplay}
					on:execute={handleExecuteReplay}
					on:logs={handleViewLogs}
					on:refresh={loadReplays}
				/>
			{/if}
		</div>

		<!-- Center: Replay Editor/Input -->
		<div>
			{#if activeView === 'editor'}
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
			{/if}
		</div>

		<!-- Right: Replay Logs/Results -->
		<div>
			{#if activeView === 'logs' && $selectedReplay}
				<ReplayLogs 
					replay={$selectedReplay}
					on:close={handleBackToList}
				/>
			{/if}
		</div>
		<div>
			{#if activeView === 'logs' && $selectedReplay}
				<ReplayLogs 
					replay={$selectedReplay}
					on:close={handleBackToList}
				/>
			{/if}
		</div>
	</div>
</div>
