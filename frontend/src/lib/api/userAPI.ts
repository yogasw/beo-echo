// User API methods for retrieving user information
import type { User } from '$lib/types/User';
import { BASE_URL_API } from './mockoonApi';

// Fetch full user profile including ownership status
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
    
    return {
      id: data.data.id,
      email: data.data.email,
      name: data.data.name,
      isOwner: data.data.is_owner || false
    };
  } catch (error) {
    console.error('Error fetching user profile:', error);
    throw error;
  }
}
