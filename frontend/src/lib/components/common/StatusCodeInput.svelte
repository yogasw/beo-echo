<script lang="ts">
  import { createEventDispatcher, onMount, afterUpdate } from 'svelte';
  import StatusCodeBadge from './StatusCodeBadge.svelte';
  
  const dispatch = createEventDispatcher<{
    change: { value: number; statusCode: StatusCode }
  }>();

  export let value: number = 200;
  export let placeholder: string = 'Enter status code';
  export let disabled: boolean = false;
  export let error: string = '';
  export let label: string = 'Status Code';
  export let showLabel: boolean = true;
  export let showQuickSelect: boolean = true;
  export let className: string = '';

  interface StatusCode {
    code: number;
    name: string;
    description: string;
    category: string;
    color: string;
    bgColor: string;
  }

  const statusCodes: StatusCode[] = [
    // 2xx Success
    { code: 200, name: 'OK', description: 'Request successful', category: 'Success', color: 'text-green-600', bgColor: 'bg-green-600' },
    { code: 201, name: 'Created', description: 'Resource created successfully', category: 'Success', color: 'text-green-600', bgColor: 'bg-green-600' },
    { code: 202, name: 'Accepted', description: 'Request accepted for processing', category: 'Success', color: 'text-green-600', bgColor: 'bg-green-600' },
    { code: 204, name: 'No Content', description: 'Success with no response body', category: 'Success', color: 'text-green-600', bgColor: 'bg-green-600' },
    { code: 206, name: 'Partial Content', description: 'Partial content served', category: 'Success', color: 'text-green-600', bgColor: 'bg-green-600' },
    
    // 3xx Redirection  
    { code: 301, name: 'Moved Permanently', description: 'Resource moved permanently', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 302, name: 'Found', description: 'Resource temporarily moved', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 304, name: 'Not Modified', description: 'Resource not modified', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 307, name: 'Temporary Redirect', description: 'Temporary redirect', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 308, name: 'Permanent Redirect', description: 'Permanent redirect', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    
    // 4xx Client Error
    { code: 400, name: 'Bad Request', description: 'Invalid request syntax', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 401, name: 'Unauthorized', description: 'Authentication required', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 403, name: 'Forbidden', description: 'Access denied', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 404, name: 'Not Found', description: 'Resource not found', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 405, name: 'Method Not Allowed', description: 'HTTP method not allowed', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 406, name: 'Not Acceptable', description: 'Response format not accepted', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 408, name: 'Request Timeout', description: 'Request timeout', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 409, name: 'Conflict', description: 'Request conflicts with current state', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 410, name: 'Gone', description: 'Resource no longer available', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 412, name: 'Precondition Failed', description: 'Precondition failed', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 413, name: 'Payload Too Large', description: 'Request payload too large', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 415, name: 'Unsupported Media Type', description: 'Media type not supported', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 422, name: 'Unprocessable Entity', description: 'Request validation failed', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 429, name: 'Too Many Requests', description: 'Rate limit exceeded', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    
    // 5xx Server Error
    { code: 500, name: 'Internal Server Error', description: 'Server encountered an error', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 501, name: 'Not Implemented', description: 'Server does not support functionality', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 502, name: 'Bad Gateway', description: 'Invalid response from upstream', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 503, name: 'Service Unavailable', description: 'Server temporarily unavailable', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 504, name: 'Gateway Timeout', description: 'Upstream server timeout', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 505, name: 'HTTP Version Not Supported', description: 'HTTP version not supported', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' }
  ];

  // Function to get appropriate color for status code based on its range
  function getStatusCodeColor(code: number): { color: string; bgColor: string; category: string } {
    if (code >= 200 && code < 300) {
      return { color: 'text-green-600', bgColor: 'bg-green-600', category: 'Success' };
    } else if (code >= 300 && code < 400) {
      return { color: 'text-blue-600', bgColor: 'bg-blue-600', category: 'Redirection' };
    } else if (code >= 400 && code < 500) {
      return { color: 'text-yellow-600', bgColor: 'bg-yellow-600', category: 'Client Error' };
    } else if (code >= 500 && code < 600) {
      return { color: 'text-red-600', bgColor: 'bg-red-600', category: 'Server Error' };
    } else {
      return { color: 'text-gray-600', bgColor: 'bg-gray-600', category: 'Custom' };
    }
  }

  // Quick select options (most commonly used)
  const quickSelectCodes = [200, 201, 204, 400, 401, 403, 404, 500];

  let isOpen = false;
  let searchTerm = '';
  let inputElement: HTMLInputElement;
  let dropdownElement: HTMLDivElement;
  let selectedIndex = -1;
  let statusBadgeElement: HTMLElement;

  $: selectedStatus = statusCodes.find(s => s.code === value) || (() => {
    const colors = getStatusCodeColor(value);
    return {
      code: value,
      name: 'Custom',
      description: `Custom ${colors.category.toLowerCase()} status code`,
      category: colors.category,
      color: colors.color,
      bgColor: colors.bgColor
    };
  })();
  
  // Enhanced filtering: supports multi-criteria search (numeric + text)
  $: filteredCodes = statusCodes.filter(status => {
    const search = searchTerm.toLowerCase().trim();
    
    if (!search) return true;
    
    // Split search terms by spaces for multi-criteria search
    const searchTerms = search.split(/\s+/).filter(term => term.length > 0);
    
    // All terms must match (AND logic)
    return searchTerms.every(term => {
      // Check if term is numeric
      if (/^\d/.test(term)) {
        // Support searching for ranges like "20" matches "200-299"
        if (term.length <= 2) {
          return status.code.toString().startsWith(term);
        } else {
          // Exact or partial code match
          return status.code.toString().includes(term);
        }
      } else {
        // Text-based search in name, description, and category
        return status.name.toLowerCase().includes(term) ||
               status.description.toLowerCase().includes(term) ||
               status.category.toLowerCase().includes(term);
      }
    });
  });

  // Calculate dynamic padding based on status badge width
  let dynamicPadding = '56px'; // Default fallback

  function selectStatusCode(status: StatusCode) {
    value = status.code;
    searchTerm = '';
    isOpen = false;
    selectedIndex = -1;
    dispatch('change', { value: status.code, statusCode: status });
  }

  function handleInputKeydown(event: KeyboardEvent) {
    if (!isOpen && (event.key === 'Enter' || event.key === 'ArrowDown' || event.key === 'ArrowUp')) {
      isOpen = true;
      selectedIndex = filteredCodes.findIndex(s => s.code === value);
      return;
    }

    if (!isOpen) return;

    switch (event.key) {
      case 'ArrowDown':
        event.preventDefault();
        selectedIndex = Math.min(selectedIndex + 1, filteredCodes.length - 1);
        break;
      case 'ArrowUp':
        event.preventDefault();
        selectedIndex = Math.max(selectedIndex - 1, -1);
        break;
      case 'Enter':
        event.preventDefault();
        if (selectedIndex >= 0 && filteredCodes[selectedIndex]) {
          selectStatusCode(filteredCodes[selectedIndex]);
        } else if (searchTerm.trim()) {
          const searchValue = searchTerm.trim();
          // Try to create custom status code
          if (!isNaN(Number(searchValue))) {
            const customCode = Number(searchValue);
            if (customCode > 0) {
              const colors = getStatusCodeColor(customCode);
              const customStatus = {
                code: customCode,
                name: 'Custom',
                description: `Custom ${colors.category.toLowerCase()} status code`,
                category: colors.category,
                color: colors.color,
                bgColor: colors.bgColor
              };
              selectStatusCode(customStatus);
            }
          }
        }
        break;
      case 'Escape':
        isOpen = false;
        selectedIndex = -1;
        searchTerm = '';
        inputElement.blur();
        break;
    }
  }

  function handleInputChange(event: Event) {
    const target = event.target as HTMLInputElement;
    searchTerm = target.value;
    
    // If the search term is a valid status code, auto-apply it immediately
    const newValue = Number(target.value);
    if (!isNaN(newValue) && newValue > 0) {
      value = newValue;
      const status = statusCodes.find(s => s.code === newValue) || (() => {
        const colors = getStatusCodeColor(newValue);
        return {
          code: newValue,
          name: 'Custom',
          description: `Custom ${colors.category.toLowerCase()} status code`,
          category: colors.category,
          color: colors.color,
          bgColor: colors.bgColor
        };
      })();
      dispatch('change', { value: newValue, statusCode: status });
    }
  }

  function handleInputFocus() {
    isOpen = true;
    searchTerm = '';
  }

  function handleInputBlur() {
    // Delay to allow click on dropdown items
    setTimeout(() => {
      if (!dropdownElement?.contains(document.activeElement)) {
        isOpen = false;
        searchTerm = '';
        selectedIndex = -1;
      }
    }, 150);
  }

  function handleClickOutside(event: MouseEvent) {
    if (!dropdownElement?.contains(event.target as Node) && 
        !inputElement?.contains(event.target as Node)) {
      isOpen = false;
      searchTerm = '';
      selectedIndex = -1;
    }
  }

  function updatePadding() {
    if (statusBadgeElement) {
      const badgeWidth = statusBadgeElement.offsetWidth;
      const containerPadding = 12; // pl-3 = 12px
      const gapBetweenBadgeAndText = 12; // 12px gap
      const totalPadding = containerPadding + badgeWidth + gapBetweenBadgeAndText;
      dynamicPadding = `${totalPadding}px`;
      
      console.log(`Status: ${selectedStatus.code}, Badge width: ${badgeWidth}px, Total padding: ${totalPadding}px`);
    }
  }

  onMount(() => {
    updatePadding();
  });

  afterUpdate(() => {
    updatePadding();
  });
</script>

<svelte:window on:click={handleClickOutside} />

<div class="relative {className}">
  {#if showLabel}
    <label for="status-code" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
      <i class="fas fa-hashtag mr-2"></i>{label}
    </label>
  {/if}

  <!-- Quick Select Buttons -->
  {#if showQuickSelect}
    <div class="grid grid-cols-4 gap-2 mb-3">
      {#each quickSelectCodes as code}
        {@const status = statusCodes.find(s => s.code === code)}
        {#if status}
          <button
            type="button"
            class="p-2 rounded-lg border text-left transition-all hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 text-xs"
            class:bg-blue-600={value === code}
            class:border-blue-500={value === code}
            class:text-white={value === code}
            class:bg-white={value !== code}
            class:dark:bg-gray-700={value !== code}
            class:border-gray-300={value !== code}
            class:dark:border-gray-600={value !== code}
            class:text-gray-800={value !== code}
            class:dark:text-gray-300={value !== code}
            class:hover:border-blue-400={value !== code}
            on:click={() => selectStatusCode(status)}
          >
            <div class="font-bold">{status.code}</div>
            <div class="opacity-80 text-xs">{status.name}</div>
          </button>
        {/if}
      {/each}
    </div>
  {/if}
  
  <div class="relative">
    <div class="relative">
      <input
        bind:this={inputElement}
        bind:value={searchTerm}
        id="status-code"
        type="text"
        class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg block w-full py-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors focus:outline-none"
        class:border-red-500={error}
        class:dark:border-red-500={error}
        style="padding-left: {dynamicPadding}"
        placeholder={isOpen ? 'Type to search (e.g., "20 ok", "404") or enter any custom code...' : (value ? value.toString() : placeholder)}
        {disabled}
        on:keydown={handleInputKeydown}
        on:focus={handleInputFocus}
        on:blur={handleInputBlur}
        on:input={handleInputChange}
        autocomplete="off"
      />
      
      <!-- Status Badge -->
      <div class="absolute inset-y-0 left-0 flex items-center pl-3">
        <span bind:this={statusBadgeElement}>
          <StatusCodeBadge statusCode={selectedStatus.code} size="sm" />
        </span>
      </div>
      
      <!-- Dropdown Arrow -->
      <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
        <i class="fas fa-chevron-down text-gray-400 transition-transform duration-200" 
           class:rotate-180={isOpen}></i>
      </div>
    </div>

    <!-- Dropdown Menu -->
    {#if isOpen}
      <div 
        bind:this={dropdownElement}
        class="absolute z-50 w-full mt-1 bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg shadow-lg max-h-64 overflow-auto"
      >
        <!-- Group by category -->
        {#each ['Success', 'Redirection', 'Client Error', 'Server Error'] as category}
          {@const categoryItems = filteredCodes.filter(code => code.category === category)}
          {#if categoryItems.length > 0}
            <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800 border-b border-gray-200 dark:border-gray-600">
              {category}
            </div>
            {#each categoryItems as status, index}
              {@const globalIndex = filteredCodes.indexOf(status)}
              <button
                type="button"
                class="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-600 focus:bg-gray-100 dark:focus:bg-gray-600 focus:outline-none transition-colors border-b border-gray-100 dark:border-gray-600 last:border-b-0"
                class:bg-blue-50={globalIndex === selectedIndex}
                class:dark:bg-blue-900={globalIndex === selectedIndex && globalIndex >= 0}
                on:click={() => selectStatusCode(status)}
              >
                <div class="flex items-center justify-between">
                  <div class="flex items-center">
                    <div class="mr-3">
                      <StatusCodeBadge statusCode={status.code} size="sm" />
                    </div>
                    <div>
                      <div class="font-medium text-gray-900 dark:text-white">{status.code} {status.name}</div>
                      <div class="text-xs text-gray-500 dark:text-gray-400">{status.description}</div>
                    </div>
                  </div>
                  {#if status.code === value}
                    <i class="fas fa-check text-blue-500"></i>
                  {/if}
                </div>
              </button>
            {/each}
          {/if}
        {/each}
        
        {#if filteredCodes.length === 0 && searchTerm.trim() && isNaN(Number(searchTerm))}
          <div class="px-4 py-3 text-gray-500 dark:text-gray-400 text-sm">
            <i class="fas fa-search mr-2"></i>No matching status codes found
            <div class="mt-2 text-xs">
              <div class="mb-2 font-medium">Search suggestions:</div>
              <ul class="list-disc list-inside space-y-1">
                <li>Status code numbers (e.g., "20" for 2xx codes, "404")</li>
                <li>Status names (e.g., "success", "error", "not found")</li>
                <li>Categories (e.g., "client error", "server error")</li>
                <li>Combined search (e.g., "20 ok", "40 error", "50 server")</li>
              </ul>
            </div>
          </div>
        {:else if filteredCodes.length === 0 && !searchTerm.trim()}
          <div class="px-4 py-3 text-gray-500 dark:text-gray-400 text-sm">
            <i class="fas fa-list mr-2"></i>Start typing to search status codes
          </div>
        {/if}
        
        <!-- Custom status code option -->
        {#if searchTerm.trim() && !isNaN(Number(searchTerm)) && Number(searchTerm) > 0}
          <div class="border-t border-gray-200 dark:border-gray-600">
            <div class="px-3 py-2 text-xs font-semibold text-gray-500 dark:text-gray-400 bg-gray-50 dark:bg-gray-800">
              {#if filteredCodes.some(s => s.code === Number(searchTerm))}
                Predefined Status Code
              {:else}
                Custom Status Code
              {/if}
            </div>
            <button
              type="button"
              class="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-600 focus:bg-gray-100 dark:focus:bg-gray-600 focus:outline-none transition-colors"
              on:click={() => {
                const customCode = Number(searchTerm);
                const existingStatus = statusCodes.find(s => s.code === customCode);
                if (existingStatus) {
                  selectStatusCode(existingStatus);
                } else {
                  const colors = getStatusCodeColor(customCode);
                  selectStatusCode({
                    code: customCode,
                    name: 'Custom',
                    description: `Custom ${colors.category.toLowerCase()} status code`,
                    category: colors.category,
                    color: colors.color,
                    bgColor: colors.bgColor
                  });
                }
              }}
            >
              <div class="flex items-center">
                <div class="mr-3">
                  <StatusCodeBadge statusCode={Number(searchTerm)} size="sm" />
                </div>
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">
                    {#if filteredCodes.some(s => s.code === Number(searchTerm))}
                      Use "{searchTerm}"
                    {:else}
                      Use "{searchTerm}" (Custom)
                    {/if}
                  </div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">
                    {#if filteredCodes.some(s => s.code === Number(searchTerm))}
                      {@const existingStatus = statusCodes.find(s => s.code === Number(searchTerm))}
                      {existingStatus?.name} - {existingStatus?.description}
                    {:else}
                      Custom {getStatusCodeColor(Number(searchTerm)).category.toLowerCase()} status code
                    {/if}
                  </div>
                </div>
                <div class="ml-auto">
                  {#if filteredCodes.some(s => s.code === Number(searchTerm))}
                    <i class="fas fa-check text-green-400"></i>
                  {:else}
                    <i class="fas fa-check text-blue-400"></i>
                  {/if}
                </div>
              </div>
            </button>
            
            <!-- Show relevant range info for guidance only for standard HTTP status codes -->
            {#if Number(searchTerm) >= 100 && Number(searchTerm) <= 599}
              {@const codeRange = Math.floor(Number(searchTerm) / 100) * 100}
              <div class="px-4 py-2 bg-blue-50 dark:bg-blue-900/20 text-blue-700 dark:text-blue-300 text-xs">
                <i class="fas fa-info-circle mr-1"></i>
                <strong>{codeRange}xx range:</strong>
                {#if codeRange === 100}
                  Informational responses (rarely used in APIs).
                {:else if codeRange === 200}
                  Success responses - indicate the request was successfully received, understood, and processed.
                {:else if codeRange === 300}
                  Redirection responses - indicate further action needs to be taken to complete the request.
                {:else if codeRange === 400}
                  Client error responses - indicate the request contains bad syntax or cannot be fulfilled.
                {:else if codeRange === 500}
                  Server error responses - indicate the server failed to fulfill a valid request.
                {/if}
              </div>
            {:else if Number(searchTerm) > 599}
              <div class="px-4 py-2 bg-gray-50 dark:bg-gray-900/20 text-gray-700 dark:text-gray-300 text-xs">
                <i class="fas fa-info-circle mr-1"></i>
                <strong>Custom status code:</strong> This is outside the standard HTTP status code range (100-599) but will work for custom API responses.
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  {#if error}
    <p class="text-red-500 dark:text-red-400 text-xs mt-1">
      <i class="fas fa-exclamation-circle mr-1"></i>{error}
    </p>
  {:else}
    <div class="text-gray-500 dark:text-gray-500 text-xs mt-1 space-y-1">
      <p>
        <i class="fas fa-info-circle mr-1 opacity-75"></i>
        {selectedStatus.name} - {selectedStatus.description}
        {#if selectedStatus.category === 'Custom'}
          • Valid custom code
        {/if}
      </p>
      {#if !isOpen && showQuickSelect}
        <p class="opacity-75">
          <i class="fas fa-lightbulb mr-1"></i>
          Tip: Use quick select buttons above, or type to search/create any custom status code
        </p>
      {/if}
    </div>
  {/if}
</div>

<style>
  .rotate-180 {
    transform: rotate(180deg);
  }
</style>
