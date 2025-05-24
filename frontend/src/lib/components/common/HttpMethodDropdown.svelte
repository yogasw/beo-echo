<script lang="ts">
  import { createEventDispatcher, onMount, afterUpdate } from 'svelte';
  import HttpMethodBadge from './HttpMethodBadge.svelte';
  
  const dispatch = createEventDispatcher<{
    change: { value: string; method: HttpMethod }
  }>();

  export let value: string = 'GET';
  export let placeholder: string = 'Select method';
  export let disabled: boolean = false;
  export let error: string = '';
  export let label: string = 'HTTP Method';
  export let showLabel: boolean = true;
  export let allowCustom: boolean = false;
  export let className: string = '';

  interface HttpMethod {
    value: string;
    label: string;
    description: string;
    color: string;
    bgColor: string;
  }

  const httpMethods: HttpMethod[] = [
    { 
      value: 'GET', 
      label: 'GET', 
      description: 'Retrieve data from server', 
      color: 'text-green-600', 
      bgColor: 'bg-green-600' 
    },
    { 
      value: 'POST', 
      label: 'POST', 
      description: 'Create new resource', 
      color: 'text-blue-600', 
      bgColor: 'bg-blue-600' 
    },
    { 
      value: 'PUT', 
      label: 'PUT', 
      description: 'Update existing resource', 
      color: 'text-yellow-600', 
      bgColor: 'bg-yellow-600' 
    },
    { 
      value: 'DELETE', 
      label: 'DELETE', 
      description: 'Remove resource', 
      color: 'text-red-600', 
      bgColor: 'bg-red-600' 
    },
    { 
      value: 'PATCH', 
      label: 'PATCH', 
      description: 'Partially update resource', 
      color: 'text-purple-600', 
      bgColor: 'bg-purple-600' 
    },
    { 
      value: 'OPTIONS', 
      label: 'OPTIONS', 
      description: 'Get allowed methods', 
      color: 'text-gray-600', 
      bgColor: 'bg-gray-600' 
    },
    { 
      value: 'HEAD', 
      label: 'HEAD', 
      description: 'Get headers only', 
      color: 'text-gray-600', 
      bgColor: 'bg-gray-600' 
    }
  ];

  let isOpen = false;
  let searchTerm = '';
  let inputElement: HTMLInputElement;
  let dropdownElement: HTMLDivElement;
  let selectedIndex = -1;
  let methodBadgeElement: HTMLElement;

  $: selectedMethod = httpMethods.find(m => m.value === value) || httpMethods[0];
  $: filteredMethods = httpMethods.filter(method => {
    const search = searchTerm.toLowerCase().trim();
    
    if (!search) return true;
    
    // Split search terms by spaces for multi-criteria search
    const searchTerms = search.split(/\s+/).filter(term => term.length > 0);
    
    // All terms must match (AND logic)
    return searchTerms.every(term =>
      method.value.toLowerCase().includes(term) ||
      method.description.toLowerCase().includes(term)
    );
  });

  // Calculate dynamic padding based on method badge width
  // Add extra padding: 12px (container left padding) + badge width + 12px (gap between badge and text)
  let dynamicPadding = '56px'; // Default fallback

  function selectMethod(method: HttpMethod) {
    value = method.value;
    searchTerm = '';
    isOpen = false;
    selectedIndex = -1;
    dispatch('change', { value: method.value, method });
  }

  function handleInputKeydown(event: KeyboardEvent) {
    if (!isOpen && (event.key === 'Enter' || event.key === 'ArrowDown' || event.key === 'ArrowUp')) {
      isOpen = true;
      selectedIndex = filteredMethods.findIndex(m => m.value === value);
      return;
    }

    if (!isOpen) return;

    switch (event.key) {
      case 'ArrowDown':
        event.preventDefault();
        selectedIndex = Math.min(selectedIndex + 1, filteredMethods.length - 1);
        break;
      case 'ArrowUp':
        event.preventDefault();
        selectedIndex = Math.max(selectedIndex - 1, -1);
        break;
      case 'Enter':
        event.preventDefault();
        if (selectedIndex >= 0 && filteredMethods[selectedIndex]) {
          selectMethod(filteredMethods[selectedIndex]);
        } else if (allowCustom && searchTerm.trim()) {
          const customMethod = {
            value: searchTerm.trim().toUpperCase(),
            label: searchTerm.trim().toUpperCase(),
            description: 'Custom HTTP method',
            color: 'text-gray-600',
            bgColor: 'bg-gray-600'
          };
          selectMethod(customMethod);
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
    if (methodBadgeElement) {
      const badgeWidth = methodBadgeElement.offsetWidth;
      const containerPadding = 10; // pl-3 = 10px
      const gapBetweenBadgeAndText = 10; // 10px gap
      const totalPadding = containerPadding + badgeWidth + gapBetweenBadgeAndText;
      dynamicPadding = `${totalPadding}px`;
      
      console.log(`Method: ${selectedMethod.value}, Badge width: ${badgeWidth}px, Total padding: ${totalPadding}px`);
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
    <label for="http-method" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
      <i class="fas fa-code mr-2"></i>{label}
    </label>
  {/if}
  
  <div class="relative">
    <div class="relative">
      <input
        bind:this={inputElement}
        bind:value={searchTerm}
        id="http-method"
        type="text"
        class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white rounded-lg block w-full py-3 pr-10 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors focus:outline-none cursor-pointer"
        class:border-red-500={error}
        class:dark:border-red-500={error}
        style="padding-left: {dynamicPadding}"
        placeholder={isOpen ? 'Type to search... (e.g., "get data", "post create")' : (selectedMethod ? selectedMethod.label : placeholder)}
        {disabled}
        on:keydown={handleInputKeydown}
        on:focus={handleInputFocus}
        on:blur={handleInputBlur}
        autocomplete="off"
      />
      
      <!-- Method Badge -->
      <div class="absolute inset-y-0 left-0 flex items-center pl-3">
        <span bind:this={methodBadgeElement}>
          <HttpMethodBadge method={selectedMethod.value} size="sm" />
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
        {#each filteredMethods as method, index}
          <button
            type="button"
            class="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-600 focus:bg-gray-100 dark:focus:bg-gray-600 focus:outline-none transition-colors border-b border-gray-100 dark:border-gray-600 last:border-b-0"
            class:bg-blue-50={index === selectedIndex}
            class:dark:bg-blue-900={index === selectedIndex && index >= 0}
            on:click={() => selectMethod(method)}
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center">
                <div class="mr-3">
                  <HttpMethodBadge method={method.value} size="sm" />
                </div>
                <div>
                  <div class="font-medium text-gray-900 dark:text-white">{method.label}</div>
                  <div class="text-xs text-gray-500 dark:text-gray-400">{method.description}</div>
                </div>
              </div>
              {#if method.value === value}
                <i class="fas fa-check text-blue-500"></i>
              {/if}
            </div>
          </button>
        {/each}
        
        {#if allowCustom && searchTerm.trim() && !filteredMethods.some(m => m.value.toLowerCase() === searchTerm.toLowerCase())}
          <button
            type="button"
            class="w-full text-left px-4 py-3 hover:bg-gray-100 dark:hover:bg-gray-600 focus:bg-gray-100 dark:focus:bg-gray-600 focus:outline-none transition-colors border-t border-gray-200 dark:border-gray-600"
            on:click={() => selectMethod({
              value: searchTerm.trim().toUpperCase(),
              label: searchTerm.trim().toUpperCase(),
              description: 'Custom HTTP method',
              color: 'text-gray-600',
              bgColor: 'bg-gray-600'
            })}
          >
            <div class="flex items-center">
              <div class="mr-3">
                <HttpMethodBadge method={searchTerm.trim().toUpperCase()} size="sm" />
              </div>
              <div>
                <div class="font-medium text-gray-900 dark:text-white">Create "{searchTerm.trim().toUpperCase()}"</div>
                <div class="text-xs text-gray-500 dark:text-gray-400">Custom HTTP method</div>
              </div>
            </div>
          </button>
        {/if}
        
        {#if filteredMethods.length === 0 && !allowCustom}
          <div class="px-4 py-3 text-gray-500 dark:text-gray-400 text-sm">
            <i class="fas fa-search mr-2"></i>No methods found
          </div>
        {/if}
      </div>
    {/if}
  </div>
  
  {#if error}
    <p class="text-red-500 dark:text-red-400 text-xs mt-1">{error}</p>
  {/if}
</div>

<style>
  .rotate-180 {
    transform: rotate(180deg);
  }
</style>
