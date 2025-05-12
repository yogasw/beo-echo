import { browser } from '$app/environment';
import { writable, derived } from 'svelte/store';
import { goto } from '$app/navigation';
import type { User } from '$lib/types/User';
import { BASE_URL_API } from '$lib/api/mockoonApi';

// Types
interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Initial state
const initialState: AuthState = {
  token: browser ? localStorage.getItem('auth_token') : null,
  user: null,
  isAuthenticated: false,
  isLoading: false,
  error: null
};

// Create the store
const authStore = writable<AuthState>(initialState);

// Derived store for checking if user is authenticated
export const isAuthenticated = derived(authStore, $authStore => $authStore.isAuthenticated);

// Derived store for getting the current user
export const currentUser = derived(authStore, $authStore => $authStore.user);

// Helper to extract user data from JWT token
function parseJwt(token: string) {
  try {
    return JSON.parse(atob(token.split('.')[1]));
  } catch (e) {
    return null;
  }
}

// Authentication actions
export const auth = {
  // Initialize authentication state from token (if exists)
  initialize: async () => {
    authStore.update(state => ({ ...state, isLoading: true }));
    
    const currentState = { ...initialState };
    authStore.subscribe(state => {
      currentState.token = state.token;
    })();
    
    if (currentState.token) {
      // Parse token to get basic user details and check expiration
      const payload = parseJwt(currentState.token);
      if (payload) {
        // Check if token is expired
        const expiryDate = new Date(payload.exp * 1000);
        const now = new Date();
        
        if (expiryDate > now) {
          // Token is still valid, get initial basic user info
          const basicUser: User = {
            id: payload.user_id,
            email: payload.email,
            name: payload.name,
            isOwner: false // Will be updated from API
          };
          
          // Update with basic info first
          authStore.update(state => ({
            ...state,
            user: basicUser,
            isAuthenticated: true,
            isLoading: true
          }));
          
          // Then fetch complete profile including owner status
          try {
            import('$lib/api/userAPI').then(async ({ fetchUserProfile }) => {
              try {
                const fullUser = await fetchUserProfile(currentState.token!);
                
                authStore.update(state => ({
                  ...state,
                  user: fullUser,
                  isLoading: false
                }));
              } catch (error) {
                console.error('Failed to fetch user profile:', error);
                authStore.update(state => ({ ...state, isLoading: false }));
              }
            });
          } catch (error) {
            console.error('Failed to import userAPI:', error);
            authStore.update(state => ({ ...state, isLoading: false }));
          }
        } else {
          // Token expired, remove it
          if (browser) {
            localStorage.removeItem('auth_token');
          }
          authStore.update(state => ({
            ...initialState,
            token: null,
            isLoading: false
          }));
        }
      } else {
        authStore.update(state => ({ ...state, isLoading: false }));
      }
    } else {
      authStore.update(state => ({ ...state, isLoading: false }));
    }
  },

  // Login
  login: async (email: string, password: string) => {
    authStore.update(state => ({ ...state, isLoading: true, error: null }));
    
    try {
      const response = await fetch(`${BASE_URL_API}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ email, password })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Failed to login');
      }
      
      // Save token to local storage
      if (browser && data.token) {
        localStorage.setItem('auth_token', data.token);
      }
      
      // Update store with user data
      authStore.update(state => ({
        ...state,
        token: data.token,
        user: data.user,
        isAuthenticated: true,
        isLoading: false
      }));
      
      return data.user;
    } catch (error: any) {
      authStore.update(state => ({
        ...state,
        isLoading: false,
        error: error.message || 'Login failed'
      }));
      
      throw error;
    }
  },

  // Register
  register: async (name: string, email: string, password: string) => {
    authStore.update(state => ({ ...state, isLoading: true, error: null }));
    
    try {
      const response = await fetch(`${BASE_URL_API}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name, email, password })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Registration failed');
      }
      
      // Save token to local storage
      if (browser && data.token) {
        localStorage.setItem('auth_token', data.token);
      }
      
      // Update store with user data
      authStore.update(state => ({
        ...state,
        token: data.token,
        user: data.user,
        isAuthenticated: true,
        isLoading: false
      }));
      
      return data.user;
    } catch (error: any) {
      authStore.update(state => ({
        ...state,
        isLoading: false,
        error: error.message || 'Registration failed'
      }));
      
      throw error;
    }
  },

  // Logout
  logout: () => {
    // Remove token from local storage
    if (browser) {
      localStorage.removeItem('auth_token');
    }
    
    // Reset auth store
    authStore.set({
      ...initialState,
      token: null
    });
    
    // Redirect to login page
    goto('/login');
  },

  // Get the auth token
  getToken: () => {
    let token: string | null = null;
    authStore.subscribe(state => {
      token = state.token;
    })();
    return token;
  },

  // Check if user is an owner
  isOwner: (): boolean => {
    let isOwner = false;
    authStore.subscribe(state => {
      isOwner = !!state.user?.isOwner;
    })();
    return isOwner;
  }
};

// Initialize auth state on load
if (browser) {
  auth.initialize();
}

export default authStore;
