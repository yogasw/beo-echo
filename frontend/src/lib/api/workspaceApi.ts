import { apiClient } from './apiClient';
import type { Workspace } from '$lib/types/Workspace';

/**
 * Workspace API service for managing workspaces
 */
export const workspaceApi = {
    /**
     * Get all workspaces (for system owners)
     */
    getAllWorkspaces: async (): Promise<Workspace[]> => {
        const response = await apiClient.get<{
            success: boolean;
            data: Workspace[];
        }>('/workspaces/all');
        return response.data.data;
    },
    
    /**
     * Toggle auto-invite status for a workspace
     * @param id Workspace ID
     * @param enabled Whether to enable or disable auto-invite
     */
    toggleAutoInvite: async (id: string, enabled: boolean): Promise<Workspace> => {
        // First, get current configuration to preserve domains and role
        const currentConfig = await apiClient.get<{
            success: boolean;
            data: {
                enabled: boolean;
                domains: string[];
                role: string;
                workspace_id: string;
                workspace_name: string;
            }
        }>(`/workspaces/${id}/auto-invite`);
        
        // Update only the enabled status
        const response = await apiClient.put<{
            success: boolean;
            data: Workspace;
        }>(`/workspaces/${id}/auto-invite`, {
            enabled,
            domains: currentConfig.data.data.domains,
            role: currentConfig.data.data.role
        });
        return response.data.data;
    },
    
    /**
     * Get a single workspace by ID
     * @param id Workspace ID
     */
    getWorkspace: async (id: string): Promise<Workspace> => {
        const response = await apiClient.get<{
            success: boolean;
            data: Workspace;
        }>(`/workspaces/${id}`);
        return response.data.data;
    },
    
    /**
     * Create a new workspace
     * @param data Workspace data
     */
    createWorkspace: async (data: { name: string }): Promise<Workspace> => {
        const response = await apiClient.post<{
            success: boolean;
            data: Workspace;
        }>('/workspaces', data);
        return response.data.data;
    },
    
    /**
     * Update a workspace
     * @param id Workspace ID
     * @param data Updated workspace data
     */
    updateWorkspace: async (id: string, data: { name: string }): Promise<Workspace> => {
        const response = await apiClient.put<{
            success: boolean;
            data: Workspace;
        }>(`/workspaces/${id}`, data);
        return response.data.data;
    },
    
    /**
     * Delete a workspace
     * @param id Workspace ID
     */
    deleteWorkspace: async (id: string): Promise<void> => {
        await apiClient.delete<{
            success: boolean;
        }>(`/workspaces/${id}`);
    },
    
    /**
     * Get workspace members
     * @param id Workspace ID
     */
    getWorkspaceMembers: async (id: string): Promise<any[]> => {
        const response = await apiClient.get<{
            success: boolean;
            data: any[];
        }>(`/workspaces/${id}/users`);
        return response.data.data;
    },
    
    /**
     * Add a user to a workspace
     * @param id Workspace ID
     * @param data User data to add
     */
    addUserToWorkspace: async (id: string, data: { email: string, role: string }): Promise<any> => {
        const response = await apiClient.post<{
            success: boolean;
            data: any;
        }>(`/workspaces/${id}/users`, data);
        return response.data.data;
    },
    
    /**
     * Remove a user from a workspace
     * @param workspaceId Workspace ID
     * @param userId User ID
     */
    removeUserFromWorkspace: async (workspaceId: string, userId: string): Promise<void> => {
        await apiClient.delete<{
            success: boolean;
        }>(`/workspaces/${workspaceId}/users/${userId}`);
    },
    
    /**
     * Update a user's role in a workspace
     * @param workspaceId Workspace ID
     * @param userId User ID
     * @param role New role
     */
    updateUserRole: async (workspaceId: string, userId: string, role: string): Promise<any> => {
        const response = await apiClient.put<{
            success: boolean;
            data: any;
        }>(`/workspaces/${workspaceId}/users/${userId}/role`, { role });
        return response.data.data;
    },
    
    /**
     * Add an existing member to a workspace by email
     * @param workspaceId Workspace ID
     * @param data Member data (email and role)
     */
    addMember: async (workspaceId: string, data: { email: string, role: string }): Promise<any> => {
        const response = await apiClient.post<{
            success: boolean;
            data: any;
        }>(`/workspaces/${workspaceId}/members`, data);
        return response.data.data;
    }
};
