import { recentProjects } from '$lib/stores/recentProjects';
import type { Project } from '$lib/api/BeoApi';

/**
 * Utility function to add a project to the recent projects list
 * This should be called whenever a user accesses or creates a project
 */
export function addProjectToRecentList(project: Project, workspaceName?: string) {
  recentProjects.addProject({
    id: project.id,
    name: project.name,
    alias: project.alias,
    mode: project.mode,
    workspaceName
  });
}

/**
 * Hook for components to easily add projects to recent list
 */
export function useRecentProjects() {
  return {
    addProject: addProjectToRecentList,
    removeProject: (projectId: string) => recentProjects.removeProject(projectId),
    store: recentProjects
  };
}
