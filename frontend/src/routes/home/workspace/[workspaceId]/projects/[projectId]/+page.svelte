<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { selectProject } from '$lib/utils/recentProjectUtils';
	import { toast } from '$lib/stores/toast';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';
	import BeoEchoLoader from '$lib/components/common/BeoEchoLoader.svelte';

	let isLoading = true;
	let error: string | null = null;

	onMount(async () => {
		try {
			const projectId = $page.params.projectId;
			const workspaceId = $page.params.workspaceId;
			if (!projectId || !workspaceId) {
				throw new Error('Project ID and Workspace ID are required');
			}

			// Use the utility function to select project and handle workspace switching
			await selectProject(projectId, workspaceId);
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
			<BeoEchoLoader message="Loading project...." size="lg" animated={true} />
		</div>
	</div>
{:else if error}
	<div class="flex items-center justify-center min-h-screen">
		<ErrorDisplay message={error} type="error" retryable={false} />
	</div>
{/if}
