<script lang="ts">
	import { showSaveButton, saveInprogress, saveButtonHandler, getEndpointsUpdateList, resetEndpointsList } from '$lib/stores/saveButton';
	import { toast } from '$lib/stores/toast';
	import { updateEndpoint as apiUpdateEndpoint } from '$lib/api/mockoonApi';

	async function handleSave() {
		saveInprogress.set(true);
		const updatesList = getEndpointsUpdateList();
		
		try {
			// Process each endpoint update
			for (const update of updatesList) {
				console.log('Saving update:', update);
				
				// Extract API-compatible update data from our update object
				const apiData: any = {};
				
				// Skip projectId and endpointId as they're used for the API path
				for (const [key, value] of Object.entries(update)) {
					if (key !== 'projectId' && key !== 'endpointId') {
						if (key === 'endpoint') {
							// If we're saving the entire endpoint, spread its properties
							Object.assign(apiData, value);
						} else {
							apiData[key] = value;
						}
					}
				}
				
				// Call the API to update the endpoint
				await apiUpdateEndpoint(update.projectId, update.endpointId, apiData);
			}
			
			// Reset the list after successful save
			resetEndpointsList();
			toast.success('Changes saved successfully');
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
