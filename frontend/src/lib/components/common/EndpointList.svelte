<!-- Endpoint List Component -->
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import EndpointListItem from './EndpointListItem.svelte';
  import SkeletonLoader from './SkeletonLoader.svelte';
  import ErrorDisplay from './ErrorDisplay.svelte';

  export let endpoints: Array<{
    id: string;
    method: string;
    path: string;
    name?: string;
    description?: string;
    status: 'active' | 'inactive' | 'error';
    statusCode?: number;
    responseTime?: number;
    lastUsed?: Date;
    isDefault?: boolean;
  }> = [];
  
  export let loading = false;
  export let error: string | null = null;
  export let emptyMessage = 'No endpoints found';
  export let showStatusCode = true;
  export let showMetrics = true;
  export let showActions = true;
  export let compact = false;
  export let searchQuery = '';
  export let filterStatus: 'all' | 'active' | 'inactive' | 'error' = 'all';
  export let filterMethod: string = 'all';
  export let className = '';

  const dispatch = createEventDispatcher<{
    edit: { endpoint: any };
    delete: { endpoint: any };
    toggle: { endpoint: any };
    duplicate: { endpoint: any };
    test: { endpoint: any };
    create: void;
    refresh: void;
  }>();

  // Filter endpoints based on search query, status, and method
  $: filteredEndpoints = endpoints.filter(endpoint => {
    // Search filter
    const matchesSearch = searchQuery === '' || 
      endpoint.path.toLowerCase().includes(searchQuery.toLowerCase()) ||
      endpoint.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      endpoint.description?.toLowerCase().includes(searchQuery.toLowerCase()) ||
      endpoint.method.toLowerCase().includes(searchQuery.toLowerCase());

    // Status filter
    const matchesStatus = filterStatus === 'all' || endpoint.status === filterStatus;

    // Method filter
    const matchesMethod = filterMethod === 'all' || endpoint.method === filterMethod;

    return matchesSearch && matchesStatus && matchesMethod;
  });

  // Group endpoints by status for better organization
  $: groupedEndpoints = {
    active: filteredEndpoints.filter(e => e.status === 'active'),
    inactive: filteredEndpoints.filter(e => e.status === 'inactive'),
    error: filteredEndpoints.filter(e => e.status === 'error')
  };

  function handleEndpointAction(event: CustomEvent, action: string) {
    dispatch(action as any, event.detail);
  }

  function handleCreate() {
    dispatch('create');
  }

  function handleRefresh() {
    dispatch('refresh');
  }

  // Get unique methods for filter dropdown
  $: uniqueMethods = [...new Set(endpoints.map(e => e.method))].sort();
</script>

<div class="space-y-4 {className}">
  <!-- Header with Filters and Actions -->
  <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
    <div class="flex items-center space-x-4">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        <i class="fas fa-list-alt mr-2"></i>
        Endpoints
        {#if filteredEndpoints.length > 0}
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">
            ({filteredEndpoints.length})
          </span>
        {/if}
      </h2>
    </div>

    <!-- Filters and Actions -->
    <div class="flex items-center space-x-3">
      <!-- Search -->
      <div class="relative">
        <input
          type="text"
          placeholder="Search endpoints..."
          bind:value={searchQuery}
          class="pl-10 pr-4 py-2 text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
        />
        <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
      </div>

      <!-- Status Filter -->
      <select
        bind:value={filterStatus}
        class="text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
      >
        <option value="all">All Status</option>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
        <option value="error">Error</option>
      </select>

      <!-- Method Filter -->
      <select
        bind:value={filterMethod}
        class="text-sm bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-lg px-3 py-2 focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
      >
        <option value="all">All Methods</option>
        {#each uniqueMethods as method}
          <option value={method}>{method}</option>
        {/each}
      </select>

      <!-- Refresh Button -->
      <button
        on:click={handleRefresh}
        class="p-2 text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors rounded-lg border border-gray-300 dark:border-gray-600"
        title="Refresh endpoints"
        aria-label="Refresh endpoints list"
      >
        <i class="fas fa-sync text-sm"></i>
      </button>

      <!-- Create Button -->
      <button
        on:click={handleCreate}
        class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white text-sm font-medium rounded-lg transition-colors flex items-center"
        aria-label="Create new endpoint"
        title="Create new endpoint"
      >
        <i class="fas fa-plus mr-2"></i>
        New Endpoint
      </button>
    </div>
  </div>

  <!-- Loading State -->
  {#if loading}
    <div class="space-y-3">
      <SkeletonLoader type="list" count={5} />
    </div>
  
  <!-- Error State -->
  {:else if error}
    <ErrorDisplay 
      message={error} 
      type="error" 
      retryable={true}
      onRetry={handleRefresh}
    />
  
  <!-- Empty State -->
  {:else if filteredEndpoints.length === 0}
    <div class="text-center py-12">
      <div class="mx-auto w-24 h-24 bg-gray-100 dark:bg-gray-800 rounded-full flex items-center justify-center mb-4">
        <i class="fas fa-route text-3xl text-gray-400"></i>
      </div>
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">
        {searchQuery || filterStatus !== 'all' || filterMethod !== 'all' ? 'No matching endpoints' : emptyMessage}
      </h3>
      <p class="text-gray-500 dark:text-gray-400 mb-6">
        {#if searchQuery || filterStatus !== 'all' || filterMethod !== 'all'}
          Try adjusting your search criteria or filters.
        {:else}
          Get started by creating your first API endpoint.
        {/if}
      </p>
      {#if !searchQuery && filterStatus === 'all' && filterMethod === 'all'}
        <button
          on:click={handleCreate}
          class="px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white font-medium rounded-lg transition-colors"
          aria-label="Create your first endpoint"
          title="Create your first endpoint"
        >
          <i class="fas fa-plus mr-2"></i>
          Create Endpoint
        </button>
      {/if}
    </div>
  
  <!-- Endpoint List -->
  {:else}
    <div class="space-y-3">
      {#each filteredEndpoints as endpoint (endpoint.id)}
        <EndpointListItem
          {endpoint}
          {showStatusCode}
          {showMetrics}
          {showActions}
          {compact}
          on:edit={e => handleEndpointAction(e, 'edit')}
          on:delete={e => handleEndpointAction(e, 'delete')}
          on:toggle={e => handleEndpointAction(e, 'toggle')}
          on:duplicate={e => handleEndpointAction(e, 'duplicate')}
          on:test={e => handleEndpointAction(e, 'test')}
        />
      {/each}
    </div>

    <!-- Summary Stats -->
    {#if !compact && endpoints.length > 0}
      <div class="bg-gray-50 dark:bg-gray-800/50 rounded-lg p-4 border border-gray-200 dark:border-gray-700">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 text-center">
          <div>
            <div class="text-2xl font-bold text-green-600 dark:text-green-400">
              {groupedEndpoints.active.length}
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">Active</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-gray-600 dark:text-gray-400">
              {groupedEndpoints.inactive.length}
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">Inactive</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-red-600 dark:text-red-400">
              {groupedEndpoints.error.length}
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">Errors</div>
          </div>
          <div>
            <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">
              {endpoints.length}
            </div>
            <div class="text-xs text-gray-500 dark:text-gray-400">Total</div>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</div>
