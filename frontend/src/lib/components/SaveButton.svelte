<script lang="ts">
  import { toast } from '$lib/stores/toast';
  import { selectedProject } from '$lib/stores/selectedConfig';
  import { updateEndpoint, updateResponse } from '$lib/api/mockoonApi';
  import { changes } from '$lib/stores/changes';
  import { saveButton } from '$lib/stores/saveButton';

  async function handleSave() {
    if (!$selectedProject) {
      toast.error('No project selected');
      return;
    }

    try {
      saveButton.setSaving(true);
      
      // Saving endpoint if there are changes
      if ($changes.hasEndpointChanges && $changes.endpointId && $changes.currentEndpointData) {
        const endpointData = {
          method: $changes.currentEndpointData.method,
          path: $changes.currentEndpointData.path,
          enabled: $changes.currentEndpointData.enabled,
          responseMode: $changes.currentEndpointData.response_mode
        };
        
        await updateEndpoint($selectedProject.id, $changes.endpointId, endpointData);
        toast.success('Endpoint updated successfully');
        changes.clearEndpointChanges();
      } 
      
      // Saving response if there are changes
      if ($changes.hasResponseChanges && $changes.endpointId && $changes.responseId && $changes.currentResponseData) {
        const responseData = {
          statusCode: $changes.currentResponseData.status_code,
          body: $changes.currentResponseData.body,
          headers: $changes.currentResponseData.headers,
          priority: $changes.currentResponseData.priority,
          delayMS: $changes.currentResponseData.delay_ms,
          stream: $changes.currentResponseData.stream,
          enabled: $changes.currentResponseData.enabled
        };
        
        await updateResponse($selectedProject.id, $changes.endpointId, $changes.responseId, responseData);
        toast.success('Response updated successfully');
        changes.clearResponseChanges();
      }
      
    } catch (error) {
      console.error('Failed to save:', error);
      toast.error('Failed to save changes');
    } finally {
      saveButton.setSaving(false);
    }
  }
</script>

{#if $saveButton.visible}
  <div class="fixed bottom-6 right-6 z-50">
    <button 
      class="flex items-center justify-center gap-2 bg-green-500 hover:bg-green-600 text-white px-4 py-3 rounded-full shadow-lg transition-all duration-200 transform hover:scale-105"
      on:click={handleSave}
      aria-label="Save changes"
      disabled={$saveButton.isSaving}
    >
      {#if $saveButton.isSaving}
        <i class="fas fa-spinner fa-spin"></i>
        <span>Saving...</span>
      {:else}
        <i class="fas fa-save"></i>
        <span>
          {#if $saveButton.saveType === 'both'}
            Save All Changes
          {:else if $saveButton.saveType === 'endpoint'}
            Save Endpoint
          {:else}
            Save Response
          {/if}
        </span>
      {/if}
    </button>
  </div>
{/if}
