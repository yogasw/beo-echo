import { writable } from 'svelte/store';
import { browser } from '$app/environment';

export interface RecentProject {
  id: string;
  name: string;
  alias: string;
  lastUsed: string; // ISO date string
  workspaceName?: string;
  workspaceId?: string;
  mode?: string;
}

const STORAGE_KEY = 'beo_echo_recent_projects';
const MAX_RECENT_PROJECTS = 5;

// Initialize store with data from localStorage
function createRecentProjectsStore() {
  let initialProjects: RecentProject[] = [];
  
  if (browser) {
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        initialProjects = JSON.parse(stored);
        // Sort by lastUsed date (most recent first)
        initialProjects.sort((a, b) => new Date(b.lastUsed).getTime() - new Date(a.lastUsed).getTime());
      }
    } catch (error) {
      console.error('Error loading recent projects from localStorage:', error);
    }
  }

  const { subscribe, set, update } = writable<RecentProject[]>(initialProjects);

  return {
    subscribe,
    
    // Add or update a project in recent list
    addProject: (project: Omit<RecentProject, 'lastUsed'>) => {
      update(projects => {
        const now = new Date().toISOString();
        const existingIndex = projects.findIndex(p => p.id === project.id);
        
        const recentProject: RecentProject = {
          ...project,
          lastUsed: now
        };
        
        let updatedProjects: RecentProject[];
        
        if (existingIndex >= 0) {
          // Update existing project and move to top
          updatedProjects = [
            recentProject,
            ...projects.filter(p => p.id !== project.id)
          ];
        } else {
          // Add new project to top
          updatedProjects = [recentProject, ...projects];
        }
        
        // Keep only the most recent projects
        const limitedProjects = updatedProjects.slice(0, MAX_RECENT_PROJECTS);
        
        // Save to localStorage
        if (browser) {
          try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(limitedProjects));
          } catch (error) {
            console.error('Error saving recent projects to localStorage:', error);
          }
        }
        
        return limitedProjects;
      });
    },
    
    // Remove a project from recent list
    removeProject: (projectId: string) => {
      update(projects => {
        const filteredProjects = projects.filter(p => p.id !== projectId);
        
        // Save to localStorage
        if (browser) {
          try {
            localStorage.setItem(STORAGE_KEY, JSON.stringify(filteredProjects));
          } catch (error) {
            console.error('Error saving recent projects to localStorage:', error);
          }
        }
        
        return filteredProjects;
      });
    },
    
    // Clear all recent projects (if needed)
    clear: () => {
      if (browser) {
        try {
          localStorage.removeItem(STORAGE_KEY);
        } catch (error) {
          console.error('Error clearing recent projects from localStorage:', error);
        }
      }
      set([]);
    }
  };
}

export const recentProjects = createRecentProjectsStore();
