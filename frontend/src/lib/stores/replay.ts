
import { writable, derived } from 'svelte/store';
import type { Replay, ReplayListItem, ReplayLog, ExecuteReplayResponse } from '$lib/types/Replay';

// Main replay list store — lightweight items from list API
export const replays = writable<ReplayListItem[]>([]);

// Replay folders
export const replayFolders = writable<import('$lib/types/Replay').ReplayFolder[]>([]);

// Combined type for UI
export type ReplayItem = (ReplayListItem & { itemType: 'replay' }) | (import('$lib/types/Replay').ReplayFolder & { itemType: 'folder' });

// Currently selected replay
export const selectedReplay = writable<Replay | null>(null);

// Replay execution results/logs
export const replayLogs = writable<ReplayLog[]>([]);

// Replay execution state
export const replayExecution = writable<{
	isExecuting: boolean;
	lastResult: ExecuteReplayResponse | null;
}>({
	isExecuting: false,
	lastResult: null
});

// Loading states
export const replayLoading = writable<{
	list: boolean;
	create: boolean;
	execute: boolean;
	delete: boolean;
	logs: boolean;
	save: boolean;
}>({
	list: false,
	create: false,
	execute: false,
	delete: false,
	logs: false,
	save: false
});

// Replay search and filter
export const replayFilter = writable<{
	searchTerm: string;
	protocol: string;
}>({
	searchTerm: '',
	protocol: ''
});

// Filtered replays based on search and filter criteria
export const filteredReplays = derived(
	[replays, replayFolders, replayFilter],
	([$replays, $folders, $filter]) => {
		const combined: ReplayItem[] = [
			...$folders.map(f => ({ ...f, itemType: 'folder' as const })),
			...$replays.map(r => ({ ...r, itemType: 'replay' as const }))
		];

		return combined.filter(item => {
			// Search by name only (url not available in list items)
			const matchesSearch = !$filter.searchTerm || 
				item.name?.toLowerCase().includes($filter.searchTerm?.toLowerCase());

			return matchesSearch;
		});
	}
);

// Replay statistics
export const replayStats = derived(
	[replays, replayLogs],
	([$replays, $logs]) => {
		const totalReplays = $replays.length;
		const totalExecutions = $logs.length;
		
		// Calculate average response time
		const avgResponseTime = $logs.length > 0 
			? $logs.reduce((sum, log) => sum + log.latency_ms, 0) / $logs.length 
			: 0;

		// Success rate (2xx status codes)
		const successfulExecutions = $logs.filter(log => 
			log.status_code >= 200 && log.status_code < 300
		).length;
		const successRate = totalExecutions > 0 
			? (successfulExecutions / totalExecutions) * 100 
			: 0;

		return {
			totalReplays,
			totalExecutions,
			avgResponseTime: Math.round(avgResponseTime),
			successRate: Math.round(successRate * 100) / 100
		};
	}
);

// Helper functions
export const replayActions = {
	// Set loading state for specific action
	setLoading: (action: 'list' | 'create' | 'execute' | 'delete' | 'logs' | 'save', isLoading: boolean) => {
		replayLoading.update(state => ({
			...state,
			[action]: isLoading
		}));
	},

	// Add a new replay to the list (after create, list is refreshed — this is for optimistic updates)
	addReplay: (replay: ReplayListItem) => {
		replays.update(list => [...list, replay]);
	},

	// Update an existing replay in the list
	updateReplay: (updatedReplay: ReplayListItem) => {
		replays.update(list => 
			list.map(replay => 
				replay.id === updatedReplay.id ? updatedReplay : replay
			)
		);
	},

	// Remove a replay from the list
	removeReplay: (replayId: string) => {
		replays.update(list => list.filter(replay => replay.id !== replayId));
		
		// Clear selected replay if it was the one being deleted
		selectedReplay.update(selected => 
			selected?.id === replayId ? null : selected
		);
	},

	// Update an existing replay folder
	updateFolder: (updatedFolder: import('$lib/types/Replay').ReplayFolder) => {
		replayFolders.update(list => 
			list.map(folder => 
				folder.id === updatedFolder.id ? updatedFolder : folder
			)
		);
	},

	// Optimistically move an item (replay or folder) to a new parent folder
	// Can be used before the API call finishes to make the UI feel responsive
	moveItem: (itemId: string, itemType: 'replay' | 'folder', newParentId: string | null) => {
		const targetParentId = newParentId === null ? undefined : newParentId;
		if (itemType === 'replay') {
			replays.update(list => list.map(r => r.id === itemId ? { ...r, folder_id: targetParentId } : r));
		} else if (itemType === 'folder') {
			replayFolders.update(list => list.map(f => f.id === itemId ? { ...f, parent_id: targetParentId } : f));
		}
	},

	// Set execution state
	setExecuting: (isExecuting: boolean) => {
		replayExecution.update(state => ({
			...state,
			isExecuting
		}));
	},

	// Set last execution result
	setLastResult: (result: ExecuteReplayResponse | null) => {
		replayExecution.update(state => ({
			...state,
			lastResult: result
		}));
	},

	// Add execution log
	addExecutionLog: (log: ReplayLog) => {
		replayLogs.update(logs => [log, ...logs]);
	},

	// Clear all state (useful for project switching)
	clearAll: () => {
		replays.set([]);
		replayFolders.set([]);
		selectedReplay.set(null);
		replayLogs.set([]);
		replayExecution.set({ isExecuting: false, lastResult: null });
		replayFilter.set({ searchTerm: '', protocol: '' });
	}
};
