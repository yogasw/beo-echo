import UserManagement from "$lib/components/instance/UserManagement.svelte";
import type { User } from "$lib/types/User";
import { apiClient } from "./apiClient";
import { toast } from "$lib/stores/toast";

export const userManagementApi = {
    /**
     * Get all users in the system
     */
    getAllUsers: async (): Promise<User[]> => {
        try {
            const { data } = await apiClient.get('/users');
            return data.data || [];
        } catch (error) {
            toast.error('Failed to fetch users');
            console.error('Failed to fetch users:', error);
            return [];
        }
    },

    /**
     * Get all users in a workspace
     * @param workspaceId The ID of the workspace
     */
    getWorkspaceUsers: async (workspaceId: string): Promise<User[]> => {
        try {
            const { data } = await apiClient.get(`/workspaces/${workspaceId}/users`);
            console.log('API response:', data);

            // Transform backend data to User type
            const users = data.data || [];
            return users.map((user: any) => ({
                id: user.userId || user.user?.id || '',
                name: user.user?.name || 'Unknown User',
                email: user.user?.email || 'unknown@example.com',
                role: mapRoleToFrontendRole(user.role),
                status: 'Active',
                initials: generateInitials(user.user?.name || 'Unknown User')
            }));
        } catch (error) {
            toast.error('Failed to fetch workspace users');
            console.error('Failed to fetch workspace users:', error);
            return [];
        }
    },

    /**
     * Add a user to a workspace
     * @param workspaceId The ID of the workspace
     * @param email The email of the user to add
     * @param role The role to assign to the user
     */
    addWorkspaceUser: async (workspaceId: string, email: string, role: string): Promise<boolean> => {
        try {
            await apiClient.post(`/workspaces/${workspaceId}/users`, {
                email,
                role: role.toLowerCase() // Convert to backend format
            });
            return true;
        } catch (error) {
            toast.error('Failed to add user to workspace');
            console.error('Failed to add user to workspace:', error);
            return false;
        }
    },

    /**
     * Update a user's role in a workspace
     * @param workspaceId The ID of the workspace
     * @param userId The ID of the user
     * @param role The new role to assign
     */
    updateWorkspaceUserRole: async (workspaceId: string, userId: string, role: string): Promise<boolean> => {
        try {
            await apiClient.patch(`/workspaces/${workspaceId}/users/${userId}/role`, {
                role: role.toLowerCase() // Convert to backend format
            });
            return true;
        } catch (error) {
            toast.error('Failed to update user role');
            console.error('Failed to update user role:', error);
            return false;
        }
    },

    /**
     * Remove a user from a workspace
     * @param workspaceId The ID of the workspace
     * @param userId The ID of the user to remove
     */
    removeWorkspaceUser: async (workspaceId: string, userId: string): Promise<boolean> => {
        try {
            await apiClient.delete(`/workspaces/${workspaceId}/users/${userId}`);
            return true;
        } catch (error) {
            toast.error('Failed to remove user from workspace');
            console.error('Failed to remove user from workspace:', error);
            return false;
        }
    },
    /**
     * Update a user's owner status
     * @param userId The ID of the user to update
     * @param is_owner Whether the user should be an owner or not
     */
    updateUser: async (userId: string, {
        is_owner,
        is_active
    }: {
        is_owner: boolean,
        is_active: boolean
    }): Promise<boolean> => {
        try {
            await apiClient.patch(`/users/${userId}`, {
                is_owner,
                is_active
            });
            return true;
        } catch (error) {
            toast.error('Failed to update user');
            console.error('Failed to update user:', error);
            return false;
        }
    },

    /**
     * Delete a user from the system
     * @param userId The ID of the user to delete
     */
    deleteUser: async (userId: string): Promise<boolean> => {
        try {
            await apiClient.delete(`/users/${userId}`);
            return true;
        } catch (error) {
            toast.error('Failed to delete user');
            console.error('Failed to delete user:', error);
            return false;
        }
    }
};

/**
 * Map backend role to frontend role format
 */
function mapRoleToFrontendRole(role: string): 'Admin' | 'User' | 'Viewer' {
    const lowercaseRole = (role || '').toLowerCase();
    if (lowercaseRole === 'admin') return 'Admin';
    if (lowercaseRole === 'viewer') return 'Viewer';
    return 'User'; // Default for 'editor', 'member', or any other role
}

/**
 * Generate initials from name
 */
function generateInitials(name: string): string {
    const nameParts = name.split(' ');
    const initials = nameParts.length > 1
        ? `${nameParts[0][0]}${nameParts[1][0]}`
        : name.substring(0, 2);
    return initials.toUpperCase();
}