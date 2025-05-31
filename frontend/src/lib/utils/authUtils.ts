import type { User } from '$lib/types/User';
import { getApiBaseUrl } from './desktopConfig';

// Define API base URL with desktop mode support
export const BASE_URL_API = getApiBaseUrl();

/**
 * Fetch full user profile including ownership status
 */
export async function fetchUserProfile(token: string): Promise<User> {
  try {
    const response = await fetch(`${BASE_URL_API}/auth/me`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
      }
    });

    if (!response.ok) {
      throw new Error('Failed to fetch user profile');
    }

    const data = await response.json();

    if (!data.success || !data.data) {
      throw new Error(data.message || 'Failed to fetch user profile data');
    }

    return data.data as User;
  } catch (error) {
    console.error('Error fetching user profile:', error);
    throw error;
  }
}
