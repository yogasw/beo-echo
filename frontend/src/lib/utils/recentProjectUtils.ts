import { goto } from '$app/navigation';
import { selectedProject } from '$lib/stores/selectedConfig';
import { workspaces, currentWorkspace } from '$lib/stores/workspace';
import { recentProjects, type RecentProject } from '$lib/stores/recentProjects';
import { getProjectDetail } from '$lib/api/BeoApi';
import { toast } from '$lib/stores/toast';
import type { Project } from '$lib/api/BeoApi';
import { get } from 'svelte/store';
import { logStatus } from '$lib/stores/logStatus';
import { initializeLogsStream } from '$lib/services/logsService';
import { activeTab } from '$lib/stores/activeTab';
import { triggerScrollToProject } from '$lib/stores/scrollToProject';

/**
 * Add a project to the recent projects list with comprehensive error handling
 * @param project - The project to add from BeoApi
 * @param workspaceName - Optional workspace name override
 * @throws {Error} If project data is invalid
 */
export function addProjectToRecent(recentProject: RecentProject): void {
    try {
        recentProjects.addProject(recentProject);
        console.log(`Project added to recent projects`, recentProject);
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
        console.log(`Selecting project with ID: ${projectId}, workspace ID: ${workspaceId}`);
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
                // remove from recent projects if it exists
                recentProjects.removeProject(projectId);
                // back to home
                await goto('/home', { replaceState: true });
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

        // Create a RecentProject object
        const recentProject: RecentProject = {
            id: project.id,
            name: project.name,
            alias: project.alias,
            workspaceId: project.workspace_id,
            workspaceName: get(currentWorkspace)?.name || 'Unknown Workspace',
            lastUsed: new Date().toISOString(),
            url: project.url || '',
        };

        // Update recent projects
        addProjectToRecent(recentProject);

        logStatus.reset();
        initializeLogsStream(project.id, 100, true);

        // Navigate to home if requested
        if (navigateToHome) {
            await goto('/home', { replaceState: true });
        }

        // Switch to logs tab after selecting recent project for better UX
        // This allows users to immediately see request logs when switching projects
        setTimeout(() => {
            activeTab.set('logs');
            // Trigger scroll to project after tab is set
            triggerScrollToProject(project.id);
        }, 200); // Ensure tab switch happens after next tick

        return project;

    } catch (error: any) {
        const errorMessage = error.message || 'Failed to select project';
        toast.error(errorMessage);
        console.error('Error selecting project:', error);
        await goto('/home', { replaceState: true });
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
        workspaceId: project.workspace_id || workspace?.id || '',
        mode: project.mode || 'mock',
        lastUsed: new Date().toISOString(),
        url: project.url || '',
    };
}
