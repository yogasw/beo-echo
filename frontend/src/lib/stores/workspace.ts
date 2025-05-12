import { writable, derived } from 'svelte/store';
import type { Workspace } from '$lib/types/User';
import authStore from './auth';

interface WorkspaceState {
  workspaces: Workspace[];
  currentWorkspace: Workspace | null;
  isLoading: boolean;
  error: string | null;
}

// Initial state
const initialState: WorkspaceState = {
  workspaces: [],
  currentWorkspace: null,
  isLoading: false,
  error: null
};

// Create the store
const workspaceStore = writable<WorkspaceState>(initialState);

// API URL
const API_URL = '/mock/api';

// Helper to get the auth token from the auth store
function getAuthToken(): string | null {
  let token: string | null = null;
  
  authStore.subscribe((state) => {
    token = state.token;
  })();
  
  return token;
}

// Derived store for getting all workspaces
export const allWorkspaces = derived(workspaceStore, $state => $state.workspaces);

// Derived store for getting the current workspace
export const currentWorkspace = derived(workspaceStore, $state => $state.currentWorkspace);

// Workspace actions
export const workspaces = {
  // Load all workspaces for the current user
  loadAll: async () => {
    workspaceStore.update(state => ({ ...state, isLoading: true, error: null }));
    
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('Not authenticated');
      }
      
      const response = await fetch(`${API_URL}/workspaces`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Failed to load workspaces');
      }
      
      // Update store with workspace data
      workspaceStore.update(state => ({
        ...state,
        workspaces: data.data || [],
        isLoading: false
      }));
      
      // If there are workspaces and no current one is set, set the first one as current
      if (data.data?.length > 0) {
        workspaceStore.update(state => {
          if (!state.currentWorkspace) {
            return {
              ...state,
              currentWorkspace: data.data[0]
            };
          }
          return state;
        });
      }
      
      return data.data;
    } catch (error) {
      workspaceStore.update(state => ({
        ...state,
        isLoading: false,
        error: error.message || 'Failed to load workspaces'
      }));
      
      throw error;
    }
  },
  
  // Create a new workspace
  create: async (name: string) => {
    workspaceStore.update(state => ({ ...state, isLoading: true, error: null }));
    
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('Not authenticated');
      }
      
      const response = await fetch(`${API_URL}/workspaces`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify({ name })
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Failed to create workspace');
      }
      
      // Update store with new workspace
      workspaceStore.update(state => ({
        ...state,
        workspaces: [...state.workspaces, data.data],
        currentWorkspace: data.data,
        isLoading: false
      }));
      
      return data.data;
    } catch (error) {
      workspaceStore.update(state => ({
        ...state,
        isLoading: false,
        error: error.message || 'Failed to create workspace'
      }));
      
      throw error;
    }
  },
  
  // Set the current workspace
  setCurrent: (workspaceId: string) => {
    workspaceStore.update(state => {
      const workspace = state.workspaces.find(w => w.id === workspaceId);
      return {
        ...state,
        currentWorkspace: workspace || null
      };
    });
  },
  
  // Check user role in a workspace
  checkRole: async (workspaceId: string, userId?: string) => {
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('Not authenticated');
      }
      
      const url = userId 
        ? `${API_URL}/workspaces/${workspaceId}/role?user_id=${userId}`
        : `${API_URL}/workspaces/${workspaceId}/role`;
      
      const response = await fetch(url, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Failed to check role');
      }
      
      return data.data;
    } catch (error) {
      throw error;
    }
  },
  
  // Get projects in a workspace
  getProjects: async (workspaceId: string) => {
    try {
      const token = getAuthToken();
      if (!token) {
        throw new Error('Not authenticated');
      }
      
      const response = await fetch(`${API_URL}/workspaces/${workspaceId}/projects`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });
      
      const data = await response.json();
      
      if (!response.ok) {
        throw new Error(data.message || 'Failed to load workspace projects');
      }
      
      return data.data;
    } catch (error) {
      throw error;
    }
  }
};

export default workspaceStore;
