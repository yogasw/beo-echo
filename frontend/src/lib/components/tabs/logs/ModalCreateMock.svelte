<script lang="ts">
  import { fade } from 'svelte/transition';
  import { addEndpoint, addResponse, type RequestLog } from '$lib/api/BeoApi';
  import { toast } from '$lib/stores/toast';
  import { theme } from '$lib/stores/theme';
  import HttpMethodDropdown from '$lib/components/common/HttpMethodDropdown.svelte';
  import StatusCodeInput from '$lib/components/common/StatusCodeInput.svelte';
  import StepModal from '$lib/components/common/StepModal.svelte';
  
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
  let currentStep = 1;
  let validationErrors: Record<string, string> = {};

  // Initialize values from log when opened
  $: if (log && isOpen) {
    method = log.method;
    path = log.path;
    statusCode = log.response_status;
    body = log.response_body;
    
    // Try to format headers as JSON if they aren't already
    try {
      const parsedHeaders = JSON.parse(log.response_headers);
      // Remove Content-Length header if it exists
      delete parsedHeaders['Content-Length'];
      headers = JSON.stringify(parsedHeaders, null, 2);
    } catch (e) {
      headers = log.response_headers;
    }
    
    // Initialize documentation with a template
    documentation = `## ${method} ${path}\n\nAuto-created from unmatched request log.\n\n### Description\nThis endpoint was automatically generated from an unmatched request. Please update the documentation to describe its purpose and behavior.`;
    
    // Reset form state
    currentStep = 1;
    validationErrors = {};
    error = null;
  }

  // Validation functions
  function validateStep1(): boolean {
    const errors: Record<string, string> = {};
    
    if (!path.trim()) {
      errors.path = 'Path is required';
    } else if (!path.startsWith('/')) {
      errors.path = 'Path must start with /';
    }
    
    if (!method) {
      errors.method = 'HTTP method is required';
    }
    
    validationErrors = errors;
    return Object.keys(errors).length === 0;
  }

  function validateStep2(): boolean {
    const errors: Record<string, string> = {};
    
    if (statusCode < 100 || statusCode > 599) {
      errors.statusCode = 'Status code must be between 100 and 599';
    }
    
    // Validate headers JSON format
    if (headers.trim()) {
      try {
        JSON.parse(headers);
      } catch (e) {
        errors.headers = 'Headers must be valid JSON format';
      }
    }
    
    // Validate body JSON format if it looks like JSON
    if (body.trim() && (body.trim().startsWith('{') || body.trim().startsWith('['))) {
      try {
        JSON.parse(body);
      } catch (e) {
        errors.body = 'Response body appears to be JSON but is not valid';
      }
    }
    
    validationErrors = errors;
    return Object.keys(errors).length === 0;
  }

  // Step navigation
  function nextStep() {
    if (currentStep === 1 && validateStep1()) {
      currentStep = 2;
    } else if (currentStep === 2 && validateStep2()) {
      currentStep = 3;
    }
  }

  function prevStep() {
    if (currentStep > 1) {
      currentStep--;
    }
  }

  // Format JSON nicely for display
  function formatJson(jsonStr: string): string {
    try {
      return JSON.stringify(JSON.parse(jsonStr), null, 2);
    } catch (e) {
      return jsonStr;
    }
  }

  // Auto-format body as JSON if it's valid
  function handleBodyInput() {
    if (body.trim() && (body.trim().startsWith('{') || body.trim().startsWith('['))) {
      try {
        const formatted = formatJson(body);
        body = formatted;
      } catch (e) {
        // Ignore formatting errors while typing
      }
    }
  }

  // Generate sample response based on method
  function generateSampleResponse() {
    const samples = {
      GET: '{\n  "data": [\n    {\n      "id": 1,\n      "name": "Sample Item"\n    }\n  ],\n  "total": 1\n}',
      POST: '{\n  "id": 1,\n  "message": "Resource created successfully",\n  "data": {\n    "name": "New Item"\n  }\n}',
      PUT: '{\n  "id": 1,\n  "message": "Resource updated successfully"\n}',
      DELETE: '{\n  "message": "Resource deleted successfully"\n}',
      PATCH: '{\n  "id": 1,\n  "message": "Resource partially updated"\n}'
    };
    
    body = samples[method as keyof typeof samples] || samples.GET;
  }

  // Handle form submission
  async function handleSubmit() {
    // Final validation
    if (!validateStep1() || !validateStep2()) {
      toast.error('Please fix the validation errors before submitting');
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

        toast.success(`Mock endpoint created: ${method} ${path}`);
        onSuccess();
        onClose();
      } else {
        throw new Error('Failed to create endpoint');
      }
    } catch (err) {
      console.error('Error creating mock:', err);
      const errorMessage = err instanceof Error ? err.message : 'An unknown error occurred';
      error = errorMessage;
      toast.error(`Failed to create mock: ${errorMessage}`);
    } finally {
      isSubmitting = false;
    }
  }

  // Check if user has made changes
  function hasUnsavedChanges(): boolean {
    return path.trim() !== '' || 
           method !== 'GET' || 
           statusCode !== 200 || 
           body.trim() !== '' || 
           headers !== '{}' || 
           documentation.trim() !== '';
  }
</script>

<StepModal
  bind:isOpen
  title="Create Mock Endpoint"
  currentStep={currentStep}
  totalSteps={3}
  stepLabels={['Endpoint', 'Response', 'Documentation']}
  onClose={onClose}
  onNext={nextStep}
  onPrevious={prevStep}
  onSubmit={handleSubmit}
  canGoNext={currentStep === 1 ? validateStep1() : currentStep === 2 ? validateStep2() : true}
  canGoPrevious={currentStep > 1}
  canSubmit={validateStep1() && validateStep2()}
  {isSubmitting}
  {error}
  hasUnsavedChanges={hasUnsavedChanges}
  nextButtonText="Continue"
  previousButtonText="Back"
  submitButtonText="Create Endpoint"
>
  <!-- Step 1: Endpoint Configuration -->
  {#if currentStep === 1}
    <div class="p-6" transition:fade={{ duration: 200 }}>
      <div class="mb-6">
        <div class="flex items-center mb-4">
          <div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center mr-3">
            <i class="fas fa-route text-white text-sm"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 dark:text-white">Endpoint Configuration</h3>
        </div>
        <p class="text-gray-600 dark:text-gray-400 text-sm">Define the HTTP method and path for your mock endpoint.</p>
      </div>
      
      <!-- Method and Path in a Grid -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <!-- Method Selection -->
        <div class="md:col-span-1">
          <HttpMethodDropdown 
            bind:value={method}
            error={validationErrors.method}
            on:change={(e: CustomEvent) => method = e.detail.value}
          />
        </div>
        
        <!-- Path Input -->
        <div class="md:col-span-2">
          <label for="path" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
            <i class="fas fa-link mr-2"></i>Path
          </label>
          <div class="relative">
            <input
              type="text"
              id="path"
              bind:value={path}
              class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg block w-full py-3 pl-4 pr-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors focus:outline-none"
              class:border-red-500={validationErrors.path}
              class:dark:border-red-500={validationErrors.path}
              placeholder="/api/users/:id"
            />
          </div>
          {#if validationErrors.path}
            <p class="text-red-500 dark:text-red-400 text-xs mt-1">{validationErrors.path}</p>
          {:else}
            <p class="text-gray-500 dark:text-gray-500 text-xs mt-1">Example: /api/users, /api/orders/:id</p>
          {/if}
        </div>
      </div>
      
      <!-- Preview Card -->
      {#if path && method}
        <div class="mt-6 p-4 bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-750 dark:to-gray-800 rounded-lg border border-gray-200 dark:border-gray-600 shadow-sm">
          <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-3 flex items-center">
            <i class="fas fa-eye mr-2 text-blue-600 dark:text-blue-400"></i>Endpoint Preview
          </h4>
          <div class="flex items-center p-3 bg-white dark:bg-gray-700 rounded-md border border-gray-100 dark:border-gray-600">
            <span class="px-3 py-1 rounded-full text-xs font-medium mr-3
              {method === 'GET' ? 'bg-green-600 text-white' : 
               method === 'POST' ? 'bg-blue-600 text-white' : 
               method === 'PUT' ? 'bg-yellow-600 text-white' : 
               method === 'DELETE' ? 'bg-red-600 text-white' : 
               method === 'PATCH' ? 'bg-purple-600 text-white' : 
               'bg-gray-600 text-white'}">
              {method}
            </span>
            <span class="text-gray-900 dark:text-white font-mono">{path}</span>
          </div>
        </div>
      {/if}
    </div>
  {/if}
  
  <!-- Step 2: Response Configuration -->
  {#if currentStep === 2}
    <div class="p-6" transition:fade={{ duration: 200 }}>
      <div class="mb-6">
        <div class="flex items-center mb-4">
          <div class="w-8 h-8 bg-green-600 rounded-full flex items-center justify-center mr-3">
            <i class="fas fa-reply text-white text-sm"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 dark:text-white">Response Configuration</h3>
        </div>
        <p class="text-gray-600 dark:text-gray-400 text-sm">Configure how your mock endpoint should respond to requests.</p>
      </div>
      
      <!-- Status Code -->
      <div class="mb-6">
        <StatusCodeInput 
          bind:value={statusCode}
          error={validationErrors.statusCode}
          on:change={(e: CustomEvent) => statusCode = e.detail.value}
        />
      </div>
      
      <!-- Response Body -->
      <div class="mb-6">
        <div class="flex justify-between items-center mb-2">
          <label for="body" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
            <i class="fas fa-file-code mr-2"></i>Response Body
          </label>
          <div class="flex space-x-2">
            <button 
              type="button"
              class="bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white rounded-md text-xs px-3 py-1 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500"
              on:click={generateSampleResponse}
              title="Generate sample response"
            >
              <i class="fas fa-magic mr-1"></i>Sample
            </button>
            <button 
              type="button"
              class="bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white rounded-md text-xs px-3 py-1 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500"
              on:click={() => body = formatJson(body)}
              title="Format JSON"
            >
              <i class="fas fa-code mr-1"></i>Format
            </button>
          </div>
        </div>
        <div class="relative">
          <textarea
            id="body"
            bind:value={body}
            on:blur={handleBodyInput}
            class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white font-mono rounded-lg block w-full py-3 px-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors resize-none focus:outline-none"
            class:border-red-500={validationErrors.body}
            class:dark:border-red-500={validationErrors.body}
            rows="8"
            placeholder='&#123;&#10;  "message": "Hello from mock endpoint",&#10;  "data": &#123;&#10;    "id": 1,&#10;    "name": "Sample"&#10;  &#125;&#10;&#125;'
          ></textarea>
          {#if body}
            <div class="absolute top-2 right-2 text-xs text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-800 px-2 py-1 rounded">
              {body.length} chars
            </div>
          {/if}
        </div>
        {#if validationErrors.body}
          <p class="text-red-500 dark:text-red-400 text-xs mt-1">{validationErrors.body}</p>
        {/if}
      </div>
      
      <!-- Headers -->
      <div class="mb-6">
        <div class="flex justify-between items-center mb-2">
          <label for="headers" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
            <i class="fas fa-tags mr-2"></i>Response Headers
          </label>
          <button 
            type="button"
            class="bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white rounded-md text-xs px-3 py-1 transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500"
            on:click={() => headers = formatJson(headers)}
            title="Format JSON"
          >
            <i class="fas fa-code mr-1"></i>Format
          </button>
        </div>
        <textarea
          id="headers"
          bind:value={headers}
          class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white font-mono rounded-lg block w-full py-3 px-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors resize-none focus:outline-none"
          class:border-red-500={validationErrors.headers}
          class:dark:border-red-500={validationErrors.headers}
          rows="5"
          placeholder='&#123;&#10;  "Content-Type": "application/json",&#10;  "Cache-Control": "no-cache"&#10;&#125;'
        ></textarea>
        {#if validationErrors.headers}
          <p class="text-red-500 dark:text-red-400 text-xs mt-1">{validationErrors.headers}</p>
        {:else}
          <p class="text-gray-500 dark:text-gray-500 text-xs mt-1">JSON format - Common headers will be added automatically</p>
        {/if}
      </div>
    </div>
  {/if}  <!-- Step 3: Documentation -->
  {#if currentStep === 3}
    <div class="p-6" transition:fade={{ duration: 200 }}>
      <div class="mb-6">
        <div class="flex items-center mb-4">
          <div class="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center mr-3">
            <i class="fas fa-book text-white text-sm"></i>
          </div>
          <h3 class="text-xl font-semibold text-gray-800 dark:text-white">Documentation & Review</h3>
        </div>
        <p class="text-gray-600 dark:text-gray-400 text-sm">Add documentation and review your mock endpoint configuration.</p>
      </div>
      
      <!-- Documentation -->
      <div class="mb-6">
        <label for="documentation" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
          <i class="fas fa-edit mr-2"></i>Documentation (Markdown)
        </label>
        <textarea
          id="documentation"
          bind:value={documentation}
          class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg block w-full py-3 px-4 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors resize-none focus:outline-none"
          rows="6"
          placeholder="## Endpoint Description

Describe your endpoint here using Markdown.

### Parameters
- `id` (path): Resource identifier

### Response
Returns a JSON object with the requested resource."
        ></textarea>
        <p class="text-gray-500 dark:text-gray-500 text-xs mt-1">Use Markdown to document your endpoint's purpose, parameters, and expected responses</p>
      </div>
      
      <!-- Configuration Summary -->
      <div class="bg-gradient-to-br from-gray-50 to-gray-100 dark:from-gray-750 dark:to-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-600 shadow-sm">
        <h4 class="text-gray-800 dark:text-white font-semibold mb-4 flex items-center">
          <i class="fas fa-eye text-blue-600 dark:text-blue-400 mr-2"></i>Configuration Summary
        </h4>
        <div class="grid grid-cols-1 gap-3">
          <div class="bg-white dark:bg-gray-800 rounded-md p-3 border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400 text-sm font-medium">Endpoint:</span>
              <div class="flex items-center">
                <span class="px-2 py-1 rounded text-xs font-medium mr-2
                  {method === 'GET' ? 'bg-green-600 text-white' : 
                   method === 'POST' ? 'bg-blue-600 text-white' : 
                   method === 'PUT' ? 'bg-yellow-600 text-white' : 
                   method === 'DELETE' ? 'bg-red-600 text-white' : 
                   method === 'PATCH' ? 'bg-purple-600 text-white' : 
                   'bg-gray-600 text-white'}">
                  {method}
                </span>
                <span class="text-gray-900 dark:text-white font-mono text-sm">{path}</span>
              </div>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-md p-3 border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400 text-sm font-medium">Status Code:</span>
              <span class="text-gray-900 dark:text-white font-mono">{statusCode}</span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-md p-3 border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400 text-sm font-medium">Response Size:</span>
              <span class="text-gray-900 dark:text-white">{body.length} characters</span>
            </div>
          </div>
          <div class="bg-white dark:bg-gray-800 rounded-md p-3 border border-gray-200 dark:border-gray-700">
            <div class="flex items-center justify-between">
              <span class="text-gray-600 dark:text-gray-400 text-sm font-medium">Headers:</span>
              <span class="text-gray-900 dark:text-white">{Object.keys(JSON.parse(headers || '{}')).length} defined</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  {/if}
</StepModal>
