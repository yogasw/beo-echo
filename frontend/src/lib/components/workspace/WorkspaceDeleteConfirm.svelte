<script lang="ts">
	import { toast } from '$lib/stores/toast';
	import { workspaces, selectedWorkspace } from '$lib/stores/workspace';
	import type { Workspace } from '$lib/types/workspace';
	
	export let show = false;
	export let workspace: Workspace | null = null;
	export let onClose: () => void;
	
	let isDeleting = false;
	
	async function handleDelete() {
		if (!workspace) return;
		
		isDeleting = true;
		
		try {
			// This would be replaced with an actual API call in a real implementation
			// const response = await fetch(`/api/workspaces/${workspace.id}`, {
			//   method: 'DELETE',
			//   headers: { 'Content-Type': 'application/json' }
			// });
			
			// if (!response.ok) throw new Error('Failed to delete workspace');
			
			// Simulate API delay
			await new Promise(resolve => setTimeout(resolve, 500));
			
			// Remove from store
			workspaces.update(ws => ws.filter(w => w.id !== workspace.id));
			
			// Reset selected workspace if we deleted the current one
			if ($selectedWorkspace && $selectedWorkspace.id === workspace.id) {
				selectedWorkspace.set(null);
			}
			
			toast.success(`Workspace "${workspace.name}" has been deleted`);
			onClose();
		} catch (error) {
			console.error('Error deleting workspace:', error);
			toast.error('Failed to delete workspace. Please try again.');
		} finally {
			isDeleting = false;
		}
	}
</script>

{#if show}
	<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
		<div class="theme-bg-secondary rounded-lg shadow-lg max-w-md w-full mx-4">
			<div class="p-6">
				<h3 class="theme-text-primary text-lg font-semibold mb-4">Delete Workspace</h3>
				
				<p class="theme-text-secondary mb-6">
					Are you sure you want to delete the workspace "{workspace?.name}"? This action cannot be undone.
				</p>
				
				<div class="flex justify-end space-x-3">
					<button
						class="px-4 py-2 rounded bg-gray-700 text-white hover:bg-gray-600 transition-colors"
						on:click={onClose}
						disabled={isDeleting}
					>
						Cancel
					</button>
					<button
						class="px-4 py-2 rounded bg-red-600 text-white hover:bg-red-700 transition-colors flex items-center"
						on:click={handleDelete}
						disabled={isDeleting}
					>
						{#if isDeleting}
							<i class="fas fa-spinner fa-spin mr-2"></i>
						{:else}
							<i class="fas fa-trash-alt mr-2"></i>
						{/if}
						Delete
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
