<!-- Demo page for testing StatusCodeInput and HttpMethodDropdown components -->
<script lang="ts">
  import HttpMethodDropdown from '$lib/components/common/HttpMethodDropdown.svelte';
  import StatusCodeInput from '$lib/components/common/StatusCodeInput.svelte';
  
  let selectedMethod = 'GET';
  let selectedStatusCode = 200;
  let selectedStatusData: any = null;

  function handleMethodChange(event: CustomEvent) {
    selectedMethod = event.detail.value;
    console.log('Method changed:', event.detail);
  }

  function handleStatusCodeChange(event: CustomEvent) {
    selectedStatusCode = event.detail.value;
    selectedStatusData = event.detail.statusCode;
    console.log('Status code changed:', event.detail);
  }
</script>

<svelte:head>
  <title>Component Demo - Beo Echo</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
  <div class="max-w-4xl mx-auto px-4">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">
        <i class="fas fa-cogs mr-3"></i>Component Demo
      </h1>
      <p class="text-gray-600 dark:text-gray-400">
        Testing enhanced StatusCodeInput and HttpMethodDropdown components with dynamic padding and multi-criteria search
      </p>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- HTTP Method Dropdown Demo -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          <i class="fas fa-list-alt mr-2"></i>HTTP Method Dropdown
        </h2>
        
        <div class="space-y-4">
          <HttpMethodDropdown 
            bind:value={selectedMethod}
            on:change={handleMethodChange}
            label="Request Method"
            showLabel={true}
          />
          
          <div class="bg-gray-50 dark:bg-gray-700 rounded-md p-4">
            <h3 class="font-medium text-gray-900 dark:text-white mb-2">Selected Value:</h3>
            <p class="text-sm text-gray-600 dark:text-gray-300">{selectedMethod}</p>
          </div>
          
          <div class="text-sm text-gray-500 dark:text-gray-400">
            <h4 class="font-medium mb-2">Features:</h4>
            <ul class="list-disc list-inside space-y-1">
              <li>Dynamic padding based on badge width</li>
              <li>Multi-criteria search (e.g., "po get", "del remove")</li>
              <li>Keyboard navigation (arrows, enter, escape)</li>
              <li>Method-specific styling with appropriate colors</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Status Code Input Demo -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
          <i class="fas fa-hashtag mr-2"></i>Status Code Input
        </h2>
        
        <div class="space-y-4">
          <StatusCodeInput 
            bind:value={selectedStatusCode}
            on:change={handleStatusCodeChange}
            label="Response Status Code"
            showLabel={true}
            showQuickSelect={true}
          />
          
          <div class="bg-gray-50 dark:bg-gray-700 rounded-md p-4">
            <h3 class="font-medium text-gray-900 dark:text-white mb-2">Selected Value:</h3>
            <div class="text-sm text-gray-600 dark:text-gray-300 space-y-1">
              <p><strong>Code:</strong> {selectedStatusCode}</p>
              {#if selectedStatusData}
                <p><strong>Name:</strong> {selectedStatusData.name}</p>
                <p><strong>Description:</strong> {selectedStatusData.description}</p>
                <p><strong>Category:</strong> {selectedStatusData.category}</p>
              {/if}
            </div>
          </div>
          
          <div class="text-sm text-gray-500 dark:text-gray-400">
            <h4 class="font-medium mb-2">Features:</h4>
            <ul class="list-disc list-inside space-y-1">
              <li>Dynamic padding based on status badge width</li>
              <li>Multi-criteria search (e.g., "20 ok", "40 error")</li>
              <li>Quick select buttons for common codes</li>
              <li>Custom status code creation (100-599)</li>
              <li>Automatic category-based coloring</li>
              <li>Range-specific guidance for custom codes</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Test Cases Section -->
    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 p-6">
      <h2 class="text-xl font-semibold text-gray-900 dark:text-white mb-4">
        <i class="fas fa-vial mr-2"></i>Test Cases
      </h2>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <h3 class="font-medium text-gray-900 dark:text-white mb-3">HTTP Method Search Tests:</h3>
          <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
            <li>• Try: "po" (should show POST, PUT)</li>
            <li>• Try: "get read" (should show GET)</li>
            <li>• Try: "del remove" (should show DELETE)</li>
            <li>• Try: "patch update" (should show PATCH)</li>
          </ul>
        </div>
        
        <div>
          <h3 class="font-medium text-gray-900 dark:text-white mb-3">Status Code Search Tests:</h3>
          <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
            <li>• Try: "20" (should show 2xx codes)</li>
            <li>• Try: "20 ok" (should show 200 OK)</li>
            <li>• Try: "40 error" (should show 4xx errors)</li>
            <li>• Try: "418" (should offer custom creation)</li>
            <li>• Try: "not found" (should show 404)</li>
          </ul>
        </div>
      </div>
    </div>

    <!-- Back Navigation -->
    <div class="mt-8 text-center">
      <a 
        href="/demo" 
        class="inline-flex items-center px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md transition-colors"
      >
        <i class="fas fa-arrow-left mr-2"></i>
        Back to Demo Menu
      </a>
    </div>
  </div>
</div>
