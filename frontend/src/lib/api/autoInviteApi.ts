import { apiClient } from '$lib/api/apiClient';
import type { AutoInviteConfig } from '$lib/types/autoInviteTypes';

/**
 * Auto-invite API service 
 * Handles retrieving and updating auto-invite configurations for workspaces
 */
export const autoInviteApi = {
    /**
     * Get auto-invite configuration for a workspace
     * @param workspaceId The ID of the workspace to get configuration for
     */
    getConfig: async (workspaceId: string): Promise<AutoInviteConfig> => {
        const response = await apiClient.get<{
            success: boolean;
            data: AutoInviteConfig;
        }>(`/workspaces/${workspaceId}/auto-invite`);

        return response.data.data;
    },

    /**
     * Update auto-invite configuration for a workspace
     * @param workspaceId The ID of the workspace to update
     * @param config The updated configuration
     */
    updateConfig: async (workspaceId: string, config: {
        enabled: boolean;
        domains: string[];
        role: string;
    }): Promise<AutoInviteConfig> => {
        const response = await apiClient.put<{
            success: boolean;
            data: AutoInviteConfig;
        }>(`/workspaces/${workspaceId}/auto-invite`, config);

        return response.data.data;
    }
};
