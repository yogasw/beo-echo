import { recentProjects } from '$lib/stores/recentProjects';
import { currentWorkspace } from '$lib/stores/workspace';
import type { Project } from '$lib/api/BeoApi';
import { get } from 'svelte/store';

/**
 * Add a project to the recent projects list with comprehensive error handling
 * @param project - The project to add from BeoApi
 * @param workspaceName - Optional workspace name override
 * @throws {Error} If project data is invalid
 */
export function addProjectToRecent(project: Project, workspaceName?: string): void {
  try {
    // Validate project data
    if (!project || !project.id || !project.name || !project.alias) {
      throw new Error('Invalid project data: missing required fields (id, name, alias)');
    }

    const workspace = get(currentWorkspace);
    const finalWorkspaceName = workspaceName || workspace?.name || 'Unknown Workspace';
    
    recentProjects.addProject({
      id: project.id,
      name: project.name,
      alias: project.alias,
      workspaceName: finalWorkspaceName,
      mode: project.mode || 'mock' // Default mode if not specified
    });

    console.log(`Project "${project.name}" added to recent projects`);
  } catch (error) {
    console.error('Error adding project to recent list:', error);
    throw error; // Re-throw for caller to handle if needed
  }
}
