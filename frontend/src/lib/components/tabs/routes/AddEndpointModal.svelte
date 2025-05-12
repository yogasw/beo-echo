<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Project } from '$lib/api/mockoonApi';
  import { addEndpoint } from '$lib/api/mockoonApi';
  import { currentWorkspace } from '$lib/stores/workspace';

  export let isOpen = false;
  export let project: Project;

  const dispatch = createEventDispatcher();

  let method = 'GET';
  let path = '/';
  let isLoading = false;
  let error = '';

  function closeModal() {
    isOpen = false;
    dispatch('close');
  }  async function handleSubmit() {
    if (!path) {
      error = 'Path is required';
      return;
    }

    // Ensure path starts with /
    if (!path.startsWith('/')) {
      path = '/' + path;
    }

    if (!$currentWorkspace) {
      error = 'No workspace selected';
      return;
    }

    isLoading = true;
    error = '';

    try {
      console.log('Creating endpoint:', { workspaceId: $currentWorkspace.id, projectId: project.id, method, path });
      const newEndpoint = await addEndpoint($currentWorkspace.id, project.id, method, path);
      console.log('Endpoint created:', newEndpoint);
      
      // Create event with the new endpoint
      const customEvent = new CustomEvent('endpointCreated', {
        detail: newEndpoint
      });
      dispatch('endpointCreated', newEndpoint);
      
      // Reset form and close modal
      method = 'GET';
      path = '/';
      closeModal();
    } catch (err) {
      console.error('Failed to create endpoint:', err);
      error = 'Failed to create endpoint. Please try again.';
    } finally {
      isLoading = false;
    }
  }

  // Available HTTP methods
  const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS', 'HEAD'];
</script>

{#if isOpen}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
    <div class="bg-gray-800 p-6 rounded-lg shadow-lg w-full max-w-md">
      <div class="flex justify-between items-center mb-4">
        <h2 class="text-xl font-bold text-white">Add New Endpoint</h2>
        <button 
          class="text-white hover:text-gray-300" 
          on:click={closeModal}
        >
          <i class="fas fa-times"></i>
        </button>
      </div>

      {#if error}
        <div class="bg-red-500 text-white p-2 rounded mb-4">
          {error}
        </div>
      {/if}

      <form on:submit|preventDefault={handleSubmit}>
        <div class="mb-4">
          <label for="method" class="block text-sm font-medium text-gray-300 mb-1">
            Method
          </label>
          <select
            id="method"
            bind:value={method}
            class="w-full bg-gray-700 border border-gray-600 rounded py-2 px-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {#each httpMethods as httpMethod}
              <option value={httpMethod}>{httpMethod}</option>
            {/each}
          </select>
        </div>

        <div class="mb-4">
          <label for="path" class="block text-sm font-medium text-gray-300 mb-1">
            Path
          </label>
          <input
            type="text"
            id="path"
            bind:value={path}
            class="w-full bg-gray-700 border border-gray-600 rounded py-2 px-3 text-white focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="/api/resource"
          />
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-300 mb-1">
            Full URL Preview
          </label>
          <div class="bg-gray-700 border border-gray-600 rounded py-2 px-3 text-gray-300">
            {project.url}{path.startsWith('/') ? path : '/' + path}
          </div>
        </div>

        <div class="flex justify-end space-x-2">
          <button
            type="button"
            class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
            on:click={closeModal}
            disabled={isLoading}
          >
            Cancel
          </button>
          <button
            type="submit"
            class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 flex items-center"
            disabled={isLoading}
          >
            {#if isLoading}
              <span class="mr-2">
                <i class="fas fa-circle-notch fa-spin"></i>
              </span>
              Adding...
            {:else}
              Add Endpoint
            {/if}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
