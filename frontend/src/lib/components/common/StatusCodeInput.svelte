<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
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
    
    // 3xx Redirection  
    { code: 301, name: 'Moved Permanently', description: 'Resource moved permanently', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 302, name: 'Found', description: 'Resource temporarily moved', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    { code: 304, name: 'Not Modified', description: 'Resource not modified', category: 'Redirection', color: 'text-blue-600', bgColor: 'bg-blue-600' },
    
    // 4xx Client Error
    { code: 400, name: 'Bad Request', description: 'Invalid request syntax', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 401, name: 'Unauthorized', description: 'Authentication required', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 403, name: 'Forbidden', description: 'Access denied', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 404, name: 'Not Found', description: 'Resource not found', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 405, name: 'Method Not Allowed', description: 'HTTP method not allowed', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 409, name: 'Conflict', description: 'Request conflicts with current state', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 422, name: 'Unprocessable Entity', description: 'Request validation failed', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    { code: 429, name: 'Too Many Requests', description: 'Rate limit exceeded', category: 'Client Error', color: 'text-yellow-600', bgColor: 'bg-yellow-600' },
    
    // 5xx Server Error
    { code: 500, name: 'Internal Server Error', description: 'Server encountered an error', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 501, name: 'Not Implemented', description: 'Server does not support functionality', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 502, name: 'Bad Gateway', description: 'Invalid response from upstream', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 503, name: 'Service Unavailable', description: 'Server temporarily unavailable', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' },
    { code: 504, name: 'Gateway Timeout', description: 'Upstream server timeout', category: 'Server Error', color: 'text-red-600', bgColor: 'bg-red-600' }
  ];

  // Quick select options (most commonly used)
  const quickSelectCodes = [200, 201, 204, 400, 401, 403, 404, 500];

  let isOpen = false;
  let searchTerm = '';
  let inputElement: HTMLInputElement;
  let dropdownElement: HTMLDivElement;
  let selectedIndex = -1;

  $: selectedStatus = statusCodes.find(s => s.code === value) || 
    { code: value, name: 'Custom', description: 'Custom status code', category: 'Custom', color: 'text-gray-600', bgColor: 'bg-gray-600' };
  
  $: filteredCodes = statusCodes.filter(status => 
    status.code.toString().includes(searchTerm) ||
    status.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
    status.description.toLowerCase().includes(searchTerm.toLowerCase())
  );

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
        } else if (searchTerm.trim() && !isNaN(Number(searchTerm))) {
          const customCode = Number(searchTerm);
          if (customCode >= 100 && customCode <= 599) {
            const customStatus = {
              code: customCode,
              name: 'Custom',
              description: 'Custom status code',
              category: 'Custom',
              color: 'text-gray-600',
              bgColor: 'bg-gray-600'
            };
            selectStatusCode(customStatus);
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
    const newValue = Number(target.value);
    
    if (!isNaN(newValue) && newValue >= 100 && newValue <= 599) {
      value = newValue;
      const status = statusCodes.find(s => s.code === newValue) || {
        code: newValue,
        name: 'Custom',
        description: 'Custom status code',
        category: 'Custom',
        color: 'text-gray-600',
        bgColor: 'bg-gray-600'
      };
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
        type="number"
        min="100"
        max="599"
        class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg block w-full py-3 pl-14 pr-10 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors focus:outline-none"
        class:border-red-500={error}
        class:dark:border-red-500={error}
        placeholder={isOpen ? 'Type to search...' : (value ? value.toString() : placeholder)}
        {disabled}
        on:keydown={handleInputKeydown}
        on:focus={handleInputFocus}
        on:blur={handleInputBlur}
        on:input={handleInputChange}
        autocomplete="off"
      />
      
      <!-- Status Badge -->
      <div class="absolute inset-y-0 left-0 flex items-center pl-3">
        <span class="px-2 py-1 rounded text-xs font-medium text-white {selectedStatus.bgColor}">
          {selectedStatus.code}
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
                    <span class="px-2 py-1 rounded text-xs font-medium text-white mr-3 {status.bgColor}">
                      {status.code}
                    </span>
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
        
        {#if filteredCodes.length === 0}
          <div class="px-4 py-3 text-gray-500 dark:text-gray-400 text-sm">
            <i class="fas fa-search mr-2"></i>No status codes found
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  {#if error}
    <p class="text-red-500 dark:text-red-400 text-xs mt-1">{error}</p>
  {:else}
    <p class="text-gray-500 dark:text-gray-500 text-xs mt-1">
      {selectedStatus.name} - {selectedStatus.description}
    </p>
  {/if}
</div>

<style>
  .rotate-180 {
    transform: rotate(180deg);
  }
</style>
