import { goto } from '$app/navigation';
import { selectedProject } from '$lib/stores/selectedConfig';
import { workspaces, currentWorkspace } from '$lib/stores/workspace';
import { recentProjects, type RecentProject } from '$lib/stores/recentProjects';
import { getProjectDetail } from '$lib/api/BeoApi';
import { toast } from '$lib/stores/toast';
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
    const finalWorkspaceId = project.workspace_id || workspace?.id;
    
    recentProjects.addProject({
      id: project.id,
      name: project.name,
      alias: project.alias,
      workspaceName: finalWorkspaceName,
      workspaceId: finalWorkspaceId,
      mode: project.mode || 'mock' // Default mode if not specified
    });

    console.log(`Project "${project.name}" added to recent projects`);
  } catch (error) {
    console.error('Error adding project to recent list:', error);
    throw error; // Re-throw for caller to handle if needed
  }
}

/**
 * Select a project and optionally switch to its workspace
 * @param projectId - The ID of the project to select
 * @param workspaceId - Optional workspace ID if known, will be fetched if not provided
 * @param navigateToHome - Whether to navigate to /home after selection (default: true)
 */
export async function selectProject(
  projectId: string, 
  workspaceId?: string,
  navigateToHome: boolean = true
): Promise<Project | null> {
  try {
    let project: Project;
    
    // If workspace ID is provided, we can set it first
    if (workspaceId) {
      await workspaces.loadAll();
      workspaces.setCurrent(workspaceId);
    }
    
    // Fetch project details
    try {
      project = await getProjectDetail(projectId);
    } catch (error) {
      // If we failed and don't have workspace ID, it might be wrong workspace
      if (!workspaceId) {
        throw new Error('Project not found. It might be in a different workspace.');
      }
      throw error;
    }
    
    // If we didn't have workspace ID initially, switch to the project's workspace
    if (!workspaceId && project.workspace_id) {
      await workspaces.loadAll();
      workspaces.setCurrent(project.workspace_id);
    }
    
    // Set the selected project
    selectedProject.set(project);
    
    // Update recent projects
    addProjectToRecent(project);
    
    // Navigate to home if requested
    if (navigateToHome) {
      await goto('/home', { replaceState: true });
    }
    
    return project;
    
  } catch (error: any) {
    const errorMessage = error.message || 'Failed to select project';
    toast.error(errorMessage);
    console.error('Error selecting project:', error);
    throw error;
  }
}

/**
 * Select a project from recent projects list
 * @param recentProject - The recent project to select
 * @param navigateToHome - Whether to navigate to /home after selection (default: false for home page usage)
 */
export async function selectRecentProject(
  recentProject: RecentProject,
  navigateToHome: boolean = false
): Promise<Project | null> {
  try {
    // Update recent projects first (move to top)
    recentProjects.addProject({
      id: recentProject.id,
      name: recentProject.name,
      alias: recentProject.alias,
      workspaceId: recentProject.workspaceId,
      workspaceName: recentProject.workspaceName,
      mode: recentProject.mode
    });
    
    // Select the project with known workspace
    return await selectProject(recentProject.id, recentProject.workspaceId, navigateToHome);
    
  } catch (error: any) {
    toast.error('Failed to open recent project');
    console.error('Error selecting recent project:', error);
    throw error;
  }
}

/**
 * Create a RecentProject object from a Project for easier handling
 * @param project - The project to convert
 * @param workspaceName - Optional workspace name
 */
export function projectToRecentProject(project: Project, workspaceName?: string): RecentProject {
  const workspace = get(currentWorkspace);
  return {
    id: project.id,
    name: project.name,
    alias: project.alias,
    workspaceName: workspaceName || workspace?.name || 'Unknown Workspace',
    workspaceId: project.workspace_id || workspace?.id,
    mode: project.mode || 'mock',
    lastUsed: new Date().toISOString()
  };
}
