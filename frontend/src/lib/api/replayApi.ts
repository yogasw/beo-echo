
import { apiClient } from './apiClient';
import type { 
	Replay, 
	CreateReplayRequest, 
	ListReplaysResponse, 
	ExecuteReplayResponse,
	ListReplayLogsResponse
} from '$lib/types/Replay';

export class ReplayApi {
	/**
	 * List all replays for a project
	 */
	static async listReplays(workspaceId: string, projectId: string): Promise<ListReplaysResponse> {
		const response = await apiClient.get(`/workspaces/${workspaceId}/projects/${projectId}/replays`);
		return response.data;
	}

	/**
	 * Create a new replay
	 */
	static async createReplay(
		workspaceId: string, 
		projectId: string, 
		replayData: CreateReplayRequest
	): Promise<{ replay: Replay; message: string }> {
		const response = await apiClient.post(
			`/workspaces/${workspaceId}/projects/${projectId}/replays`,
			replayData
		);
		return response.data;
	}

	/**
	 * Get a specific replay
	 */
	static async getReplay(
		workspaceId: string, 
		projectId: string, 
		replayId: string
	): Promise<{ replay: Replay }> {
		const response = await apiClient.get(
			`/workspaces/${workspaceId}/projects/${projectId}/replays/${replayId}`
		);
		return response.data;
	}

	/**
	 * Execute a replay
	 */
	static async executeReplay(
		workspaceId: string, 
		projectId: string, 
		replayId: string
	): Promise<ExecuteReplayResponse> {
		const response = await apiClient.post(
			`/workspaces/${workspaceId}/projects/${projectId}/replays/${replayId}/execute`
		);
		return response.data;
	}

	/**
	 * Delete a replay
	 */
	static async deleteReplay(
		workspaceId: string, 
		projectId: string, 
		replayId: string
	): Promise<{ message: string }> {
		const response = await apiClient.delete(
			`/workspaces/${workspaceId}/projects/${projectId}/replays/${replayId}`
		);
		return response.data;
	}

	/**
	 * Get replay execution logs
	 */
	static async getReplayLogs(
		workspaceId: string, 
		projectId: string, 
		replayId?: string
	): Promise<ListReplayLogsResponse> {
		const url = `/workspaces/${workspaceId}/projects/${projectId}/replays/${replayId}/logs`;
		const response = await apiClient.get(url);
		return response.data;
	}
}

// Export the class directly for static method access
export const replayApi = ReplayApi;
