
import type { User } from '$lib/types/User';
import { BASE_URL_API } from '$lib/utils/authUtils';
import { auth } from '$lib/stores/auth';

interface UpdateUserPayload {
  name?: string;
  email?: string;
  isEnabled?: boolean;
}

/**
 * API service for user management
 */
export const UserAPI = {
  /**
   * Update user profile
   * @param userId User ID to update
   * @param data Update payload
   * @returns Updated user data
   */
  async updateProfile(userId: string, data: UpdateUserPayload): Promise<User> {
    const token = auth.getToken();
    
    const response = await fetch(`${BASE_URL_API}/users/${userId}`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(data)
    });
    
    if (!response.ok) {
      throw new Error(`Failed to update user: ${response.statusText}`);
    }
    
    const updatedUser = await response.json();
    
    // Update the store with the new user data
    auth.updateUserData(updatedUser);
    
    return updatedUser;
  },
  
  /**
   * Toggle user enabled status
   * @param userId User ID to update
   * @param isEnabled New enabled status
   * @returns Updated user data
   */
  async toggleUserEnabled(userId: string, isEnabled: boolean): Promise<User> {
    return this.updateProfile(userId, { isEnabled });
  }
};
