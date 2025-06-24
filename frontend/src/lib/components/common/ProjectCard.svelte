<script lang="ts">
  import { addProjectToRecent } from '$lib/utils/recentProjectsUtils';
  import type { Project } from '$lib/api/BeoApi';
  import { toast } from '$lib/stores/toast';
  
  // Example project data
  export let project: Project;
  export let onProjectOpen: ((project: Project) => void) | null = null;
  
  async function handleOpenProject() {
    try {
      // Add to recent projects first
      addProjectToRecent(project);
      
      // Show success notification
      toast.success(`Project "${project.name}" opened`);
      
      // Call parent handler if provided
      if (onProjectOpen) {
        onProjectOpen(project);
      }
    } catch (error) {
      toast.error('Failed to open project');
      console.error('Error opening project:', error);
    }
  }
  
  function getModeColor(mode: string): string {
    switch (mode?.toLowerCase()) {
      case 'mock':
        return 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200';
      case 'proxy':
        return 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200';
      case 'forwarder':
        return 'bg-orange-100 dark:bg-orange-900 text-orange-800 dark:text-orange-200';
      default:
        return 'bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200';
    }
  }
</script>

<!-- Project Card with Accessibility -->
<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-4 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors">
  <div class="flex justify-between items-start mb-3">
    <h3 class="text-lg font-semibold text-gray-900 dark:text-white truncate">
      {project.name}
    </h3>
    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getModeColor(project.mode)}">
      {project.mode}
    </span>
  </div>
  
  <div class="space-y-2 mb-4">
    <div class="flex items-center text-sm text-gray-600 dark:text-gray-400">
      <i class="fas fa-link mr-2 text-blue-500" aria-hidden="true"></i>
      <span class="truncate" title={project.url}>
        {project.alias}.beo-echo.dev
      </span>
    </div>
    
    <div class="flex items-center text-sm text-gray-600 dark:text-gray-400">
      <i class="fas fa-circle mr-2 text-green-500" aria-hidden="true"></i>
      <span>Status: {project.status || 'stopped'}</span>
    </div>
  </div>
  
  <!-- Action Buttons with MANDATORY Accessibility -->
  <div class="flex space-x-2">
    <button
      class="flex-1 bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md text-sm font-medium transition-colors flex items-center justify-center"
      on:click={handleOpenProject}
      title="Open project {project.name} and add to recent list"
      aria-label="Open project {project.name} and add to recent list"
    >
      <i class="fas fa-play mr-2" aria-hidden="true"></i>
      Open Project
    </button>
    
    <button
      class="bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-500 text-gray-700 dark:text-gray-200 py-2 px-3 rounded-md text-sm transition-colors"
      title="View project settings for {project.name}"
      aria-label="View project settings for {project.name}"
    >
      <i class="fas fa-cog" aria-hidden="true"></i>
    </button>
    
    <button
      class="bg-gray-100 dark:bg-gray-600 hover:bg-gray-200 dark:hover:bg-gray-500 text-gray-700 dark:text-gray-200 py-2 px-3 rounded-md text-sm transition-colors"
      title="Copy project URL {project.alias}.beo-echo.dev to clipboard"
      aria-label="Copy project URL {project.alias}.beo-echo.dev to clipboard"
    >
      <i class="fas fa-copy" aria-hidden="true"></i>
    </button>
  </div>
</div>
