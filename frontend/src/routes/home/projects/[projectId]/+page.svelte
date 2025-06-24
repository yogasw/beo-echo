<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { selectProject } from '$lib/utils/recentProjectUtils';
  import { toast } from '$lib/stores/toast';
  import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
  import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

  let isLoading = true;
  let error: string | null = null;

  onMount(async () => {
    try {
      const projectId = $page.params.projectId;
      
      if (!projectId) {
        throw new Error('Project ID is required');
      }

      // Use the utility function to select project and handle workspace switching
      await selectProject(projectId);
      
    } catch (err: any) {
      error = err.message || 'Failed to load project';
      toast.error(err);
      console.error('Error loading project:', err);
    } finally {
      isLoading = false;
    }
  });
</script>

{#if isLoading}
  <div class="flex items-center justify-center min-h-screen">
    <div class="text-center">
      <SkeletonLoader type="custom" />
      <p class="mt-4 text-sm text-gray-600 dark:text-gray-400">Loading project...</p>
    </div>
  </div>
{:else if error}
  <div class="flex items-center justify-center min-h-screen">
    <ErrorDisplay 
      message={error} 
      type="error" 
      retryable={false}
    />
  </div>
{/if}
