<script lang="ts">
  import { fade } from 'svelte/transition';
  import { addEndpoint, addResponse, type RequestLog } from '$lib/api/BeoApi';
  
  export let isOpen: boolean;
  export let log: RequestLog | null = null;
  export let projectId: string;
  export let onClose: () => void;
  export let onSuccess: () => void;

  let path = '';
  let method = 'GET';
  let statusCode = 200;
  let body = '';
  let headers = '{}';
  let documentation = '';
  let isSubmitting = false;
  let error: string | null = null;

  // HTTP methods for dropdown
  const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'OPTIONS', 'HEAD'];

  // Initialize values from log when opened
  $: if (log && isOpen) {
    method = log.method;
    path = log.path;
    statusCode = log.response_status;
    body = log.response_body;
    
    // Try to format headers as JSON if they aren't already
    try {
      const parsedHeaders = JSON.parse(log.response_headers);
      headers = JSON.stringify(parsedHeaders, null, 2);
    } catch (e) {
      headers = log.response_headers;
    }
    
    // Initialize documentation with a template
    documentation = `## ${method} ${path}\n\nAuto-created from unmatched request.`;
  }

  // Format JSON nicely for display
  function formatJson(jsonStr: string): string {
    try {
      return JSON.stringify(JSON.parse(jsonStr), null, 2);
    } catch (e) {
      return jsonStr;
    }
  }

  // Handle form submission
  async function handleSubmit() {
    if (!path || !method) {
      error = 'Path and method are required';
      return;
    }

    try {
      isSubmitting = true;
      error = null;

      // First create the endpoint
      const endpoint = await addEndpoint(projectId, method, path);

      if (endpoint && endpoint.id) {
        // Then create the response for this endpoint
        await addResponse(
          projectId,
          endpoint.id,
          statusCode,
          body,
          headers,
          documentation
        );

        onSuccess();
        onClose();
      } else {
        throw new Error('Failed to create endpoint');
      }
    } catch (err) {
      console.error('Error creating mock:', err);
      error = err instanceof Error ? err.message : 'An unknown error occurred';
    } finally {
      isSubmitting = false;
    }
  }

  // Close the modal when clicking outside or on close button
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      onClose();
    }
  }
</script>

{#if isOpen}
  <div
    class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
    transition:fade={{ duration: 200 }}
    on:click={handleBackdropClick}
    role="dialog"
    aria-modal="true"
  >
    <div 
      class="bg-gray-800 rounded-lg w-full max-w-3xl max-h-[90vh] overflow-y-auto shadow-xl"
      on:click|stopPropagation={() => {}}
    >
      <!-- Header -->
      <div class="p-4 border-b border-gray-700 flex justify-between items-center">
        <h2 class="text-xl font-semibold text-white">Create Mock Endpoint</h2>
        <button
          class="text-gray-400 hover:text-white"
          on:click={onClose}
          aria-label="Close"
        >
          <i class="fas fa-times"></i>
        </button>
      </div>

      <!-- Body -->
      <div class="p-4">
        {#if error}
          <div class="mb-4 p-3 bg-red-900/30 border border-red-700 rounded text-white">
            <div class="flex items-center">
              <i class="fas fa-exclamation-triangle text-yellow-400 mr-2"></i>
              <span>{error}</span>
            </div>
          </div>
        {/if}

        <form on:submit|preventDefault={handleSubmit}>
          <!-- Endpoint section -->
          <div class="mb-5">
            <h3 class="text-gray-300 font-medium mb-3 pb-1 border-b border-gray-700">Endpoint Configuration</h3>
            
            <!-- Method dropdown -->
            <div class="mb-4">
              <label for="method" class="block text-sm font-medium text-gray-300 mb-1">
                HTTP Method
              </label>
              <div class="relative">
                <select
                  id="method"
                  bind:value={method}
                  class="bg-gray-700 border border-gray-600 text-white rounded-md block w-full py-2 pl-3 pr-10"
                >
                  {#each methods as httpMethod}
                    <option value={httpMethod}>{httpMethod}</option>
                  {/each}
                </select>
              </div>
            </div>
            
            <!-- Path input -->
            <div class="mb-4">
              <label for="path" class="block text-sm font-medium text-gray-300 mb-1">
                Path
              </label>
              <input
                type="text"
                id="path"
                bind:value={path}
                class="bg-gray-700 border border-gray-600 text-white rounded-md block w-full py-2 px-3"
                placeholder="/api/resources"
              />
            </div>
          </div>
          
          <!-- Response section -->
          <div class="mb-5">
            <h3 class="text-gray-300 font-medium mb-3 pb-1 border-b border-gray-700">Response Configuration</h3>
            
            <!-- Status code -->
            <div class="mb-4">
              <label for="statusCode" class="block text-sm font-medium text-gray-300 mb-1">
                Status Code
              </label>
              <input
                type="number"
                id="statusCode"
                bind:value={statusCode}
                min="100"
                max="599"
                class="bg-gray-700 border border-gray-600 text-white rounded-md block w-full py-2 px-3"
              />
            </div>
            
            <!-- Response body -->
            <div class="mb-4">
              <label for="body" class="block text-sm font-medium text-gray-300 mb-1">
                Response Body
              </label>
              <div class="relative">
                <button 
                  type="button"
                  class="absolute right-2 top-2 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs px-2 py-1"
                  on:click={() => body = formatJson(body)}
                  title="Format JSON"
                >
                  <i class="fas fa-code"></i> Format
                </button>
                <textarea
                  id="body"
                  bind:value={body}
                  class="bg-gray-700 border border-gray-600 text-white font-mono rounded-md block w-full py-2 px-3"
                  rows="7"
                ></textarea>
              </div>
            </div>
            
            <!-- Headers -->
            <div class="mb-4">
              <label for="headers" class="block text-sm font-medium text-gray-300 mb-1">
                Response Headers (JSON format)
              </label>
              <div class="relative">
                <button 
                  type="button"
                  class="absolute right-2 top-2 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs px-2 py-1"
                  on:click={() => headers = formatJson(headers)}
                  title="Format JSON"
                >
                  <i class="fas fa-code"></i> Format
                </button>
                <textarea
                  id="headers"
                  bind:value={headers}
                  class="bg-gray-700 border border-gray-600 text-white font-mono rounded-md block w-full py-2 px-3"
                  rows="4"
                  placeholder="Headers"
                ></textarea>
              </div>
            </div>
          </div>
          
          <!-- Documentation section -->
          <div class="mb-5">
            <h3 class="text-gray-300 font-medium mb-3 pb-1 border-b border-gray-700">Documentation</h3>
            
            <div class="mb-4">
              <label for="documentation" class="block text-sm font-medium text-gray-300 mb-1">
                Documentation (Markdown)
              </label>
              <textarea
                id="documentation"
                bind:value={documentation}
                class="bg-gray-700 border border-gray-600 text-white rounded-md block w-full py-2 px-3"
                rows="5"
                placeholder="## Endpoint Description
                
Describe your endpoint here using Markdown."
              ></textarea>
            </div>
          </div>
        </form>
      </div>

      <!-- Footer -->
      <div class="p-4 border-t border-gray-700 flex justify-end space-x-3">
        <button
          type="button"
          class="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md text-sm"
          on:click={onClose}
          disabled={isSubmitting}
        >
          Cancel
        </button>
        <button
          type="button"
          class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center"
          on:click={handleSubmit}
          disabled={isSubmitting}
        >
          {#if isSubmitting}
            <span class="animate-spin h-4 w-4 border-t-2 border-b-2 border-white rounded-full mr-2"></span>
            Creating...
          {:else}
            <i class="fas fa-save mr-1"></i> Create Mock
          {/if}
        </button>
      </div>
    </div>
  </div>
{/if}
