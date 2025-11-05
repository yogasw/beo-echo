import { apiClient } from './apiClient';
import type {
	Action,
	ActionTypeInfo,
	CreateActionRequest,
	UpdateActionRequest,
	ListActionsResponse,
	ActionResponse,
	ActionMessageResponse
} from '$lib/types/Action';

export class ActionsApi {
	/**
	 * Get all available action types
	 */
	static async getActionTypes(): Promise<{ success: boolean; data: ActionTypeInfo[] }> {
		const response = await apiClient.get('/action-types');
		return response.data;
	}
	/**
	 * List all actions for a project
	 */
	static async listActions(workspaceId: string, projectId: string): Promise<ListActionsResponse> {
		const response = await apiClient.get(`/workspaces/${workspaceId}/projects/${projectId}/actions`);
		return response.data;
	}

	/**
	 * Create a new action
	 */
	static async createAction(
		workspaceId: string,
		projectId: string,
		actionData: CreateActionRequest
	): Promise<ActionResponse> {
		const response = await apiClient.post(
			`/workspaces/${workspaceId}/projects/${projectId}/actions`,
			actionData
		);
		return response.data;
	}

	/**
	 * Get a specific action
	 */
	static async getAction(
		workspaceId: string,
		projectId: string,
		actionId: string
	): Promise<ActionResponse> {
		const response = await apiClient.get(
			`/workspaces/${workspaceId}/projects/${projectId}/actions/${actionId}`
		);
		return response.data;
	}

	/**
	 * Update an existing action
	 */
	static async updateAction(
		workspaceId: string,
		projectId: string,
		actionId: string,
		actionData: UpdateActionRequest
	): Promise<ActionResponse> {
		const response = await apiClient.put(
			`/workspaces/${workspaceId}/projects/${projectId}/actions/${actionId}`,
			actionData
		);
		return response.data;
	}

	/**
	 * Toggle action enabled/disabled status
	 */
	static async toggleAction(
		workspaceId: string,
		projectId: string,
		actionId: string
	): Promise<ActionResponse> {
		const response = await apiClient.post(
			`/workspaces/${workspaceId}/projects/${projectId}/actions/${actionId}/toggle`
		);
		return response.data;
	}

	/**
	 * Delete an action
	 */
	static async deleteAction(
		workspaceId: string,
		projectId: string,
		actionId: string
	): Promise<ActionMessageResponse> {
		const response = await apiClient.delete(
			`/workspaces/${workspaceId}/projects/${projectId}/actions/${actionId}`
		);
		return response.data;
	}
}

// Export the class directly for static method access
export const actionsApi = ActionsApi;
