
import { writable, derived } from 'svelte/store';
import type { Replay, ReplayLog, ReplayExecutionResult } from '$lib/types/Replay';

// Main replay list store
export const replays = writable<Replay[]>([]);

// Currently selected replay
export const selectedReplay = writable<Replay | null>(null);

// Replay execution results/logs
export const replayLogs = writable<ReplayLog[]>([]);

// Replay execution state
export const replayExecution = writable<{
	isExecuting: boolean;
	lastResult: ReplayExecutionResult | null;
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
}>({
	list: false,
	create: false,
	execute: false,
	delete: false,
	logs: false
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
	[replays, replayFilter],
	([$replays, $filter]) => {
		return $replays.filter(replay => {
			// Search by alias or URL
			const matchesSearch = !$filter.searchTerm || 
				replay.alias?.toLowerCase().includes($filter.searchTerm?.toLowerCase()) ||
				replay.url?.toLowerCase().includes($filter.searchTerm?.toLowerCase());
			// Filter by protocol
			const matchesProtocol = !$filter.protocol || replay.protocol === $filter.protocol;

			return matchesSearch && matchesProtocol;
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
	setLoading: (action: 'list' | 'create' | 'execute' | 'delete' | 'logs', isLoading: boolean) => {
		replayLoading.update(state => ({
			...state,
			[action]: isLoading
		}));
	},

	// Add a new replay to the list
	addReplay: (replay: Replay) => {
		replays.update(list => [...list, replay]);
	},

	// Update an existing replay
	updateReplay: (updatedReplay: Replay) => {
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

	// Set execution state
	setExecuting: (isExecuting: boolean) => {
		replayExecution.update(state => ({
			...state,
			isExecuting
		}));
	},

	// Set last execution result
	setLastResult: (result: ReplayExecutionResult | null) => {
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
		selectedReplay.set(null);
		replayLogs.set([]);
		replayExecution.set({ isExecuting: false, lastResult: null });
		replayFilter.set({ searchTerm: '', protocol: '' });
	}
};
