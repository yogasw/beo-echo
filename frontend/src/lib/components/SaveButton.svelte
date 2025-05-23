<script lang="ts">
	import { showSaveButton, saveInprogress, saveButtonHandler, getEndpointsUpdateList, resetEndpointsList } from '$lib/stores/saveButton';
	import { toast } from '$lib/stores/toast';
	import { updateEndpoint as apiUpdateEndpoint, updateResponse as apiUpdateResponse, getProjectDetail } from '$lib/api/BeoApi';
	import { currentWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';

	async function handleSave() {
		saveInprogress.set(true);
		const updatesList = getEndpointsUpdateList();
		
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			saveInprogress.set(false);
			return;
		}

		try {
			// Process each update (endpoint or response)
			for (const update of updatesList) {
				console.log('Saving update:', update);
				
				// Extract API-compatible update data from our update object
				const apiData: any = {};
				
				// Check if this is a response update (has responseId) or endpoint update
				const isResponseUpdate = 'responseId' in update;
				
				// Skip projectId, endpointId, and responseId as they're used for the API path
				for (const [key, value] of Object.entries(update)) {
					if (key !== 'projectId' && key !== 'endpointId' && key !== 'responseId') {
						if (key === 'endpoint') {
							// If we're saving the entire endpoint, spread its properties
							Object.assign(apiData, value);
						} else if (key === 'response') {
							// If we're saving the entire response, spread its properties
							Object.assign(apiData, value);
						} else {
							apiData[key] = value;
						}
					}
				}
				
				if (isResponseUpdate) {
					// Call the API to update the response
					await apiUpdateResponse(
						update.projectId, 
						update.endpointId, 
						(update as any).responseId, 
						apiData
					);
				} else {
					// Call the API to update the endpoint
					await apiUpdateEndpoint(update.projectId, update.endpointId, apiData);
				}
			}
			
			// Reset the list after successful save
			resetEndpointsList();
			toast.success('Changes saved successfully');
			
			// Refresh the selected project data to update the endpoints in the store
			if ($selectedProject) {
				const refreshedProject = await getProjectDetail($selectedProject.id);
				selectedProject.set(refreshedProject);
			}
		} catch (error) {
			console.error('Error saving changes:', error);
			toast.error('Failed to save changes');
		} finally {
			saveInprogress.set(false);
		}
	}
</script>

{#if $showSaveButton}
	<div class="fixed bottom-6 right-6 z-50">
		<button
			class="flex items-center justify-center gap-2 bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-full shadow-lg transition-all duration-200 transform hover:scale-105"
			on:click={handleSave}
			aria-label="Save changes"
			disabled={$saveInprogress}
		>
			{#if $saveInprogress}
				<i class="fas fa-spinner fa-spin"></i>
				<span>Saving...</span>
			{:else}
				<i class="fas fa-save"></i>
				<span>
					Save
				</span>
			{/if}
		</button>
	</div>
{/if}
