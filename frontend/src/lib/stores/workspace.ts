import { writable, derived } from 'svelte/store';
import type { Workspace } from '$lib/types/User';
import authStore from './auth';
import { createWorkspace, getWorkspaces } from '$lib/api/BeoApi';
import { setCurrentWorkspaceId } from '$lib/utils/localStorage';

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
export const workspaceStore = writable<WorkspaceState>(initialState);

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

// Create a selectedWorkspace store (alias for currentWorkspace for better semantics)
export const selectedWorkspace = currentWorkspace;

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
      const workspaces = await getWorkspaces();
      if (!workspaces) {
        console.error('No workspaces found');
      }

      // Update store with workspace data
      workspaceStore.update(state => ({
        ...state,
        workspaces: workspaces || [],
        isLoading: false
      }));

      // If there are workspaces and no current one is set, set the first one as current
      if (workspaces?.length > 0) {
        workspaceStore.update(state => {
          if (!state.currentWorkspace) {
            return {
              ...state,
              currentWorkspace: workspaces[0]
            };
          }
          return state;
        });
      }

      return workspaces;
    } catch (error: any) {
      workspaceStore.update(state => ({
        ...state,
        isLoading: false,
        error: error.message || 'Failed to load workspaces'
      }));
      console.error('Failed to load workspaces:', error);
    }
  },

  // Create a new workspace
  create: async (name: string) => {
    workspaceStore.update(state => ({ ...state, isLoading: true, error: null }));

    try {
      const workspace = await createWorkspace(name);
      if (!workspace) {
        throw new Error(workspace || 'Failed to create workspace');
      }

      // Update store with new workspace
      workspaceStore.update(state => ({
        ...state,
        workspaces: [...state.workspaces, workspace],
        currentWorkspace: workspace,
        isLoading: false
      }));

      return workspace;
    } catch (error: any) {
      workspaceStore.update(state => ({
        ...state,
        isLoading: false,
        error: error?.message || 'Failed to create workspace'
      }));

      throw error;
    }
  },

  // Set the current workspace
  setCurrent: (workspaceId: string) => {
    console.log('Setting current workspace:', workspaceId);
    setCurrentWorkspaceId(workspaceId);
    workspaceStore.update(state => {
      const workspace = state.workspaces.find(w => w.id === workspaceId);
      return {
        ...state,
        currentWorkspace: workspace || null
      };
    });
  },

  switchWorkspace: (workspace: Workspace) => {
    console.log('Switching to workspace:', workspace);
    setCurrentWorkspaceId(workspace.id);
    workspaceStore.update(state => ({
      ...state,
      currentWorkspace: workspace
    }));
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

};

export default workspaceStore;
