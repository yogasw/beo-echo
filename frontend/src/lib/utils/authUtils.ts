import type { User } from '$lib/types/User';

// Define API base URL
export const BASE_URL_API = import.meta.env.VITE_API_BASE_URL || 'http://localhost:3600/api';

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
