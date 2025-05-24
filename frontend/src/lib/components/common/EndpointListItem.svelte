<!-- Endpoint List Item Component -->
<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import HttpMethodBadge from './HttpMethodBadge.svelte';
  import StatusCodeBadge from './StatusCodeBadge.svelte';

  export let endpoint: {
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
  };
  
  export let showStatusCode = true;
  export let showMetrics = true;
  export let showActions = true;
  export let compact = false;
  export let className = '';

  const dispatch = createEventDispatcher<{
    edit: { endpoint: typeof endpoint };
    delete: { endpoint: typeof endpoint };
    toggle: { endpoint: typeof endpoint };
    duplicate: { endpoint: typeof endpoint };
    test: { endpoint: typeof endpoint };
  }>();

  function handleEdit() {
    dispatch('edit', { endpoint });
  }

  function handleDelete() {
    dispatch('delete', { endpoint });
  }

  function handleToggle() {
    dispatch('toggle', { endpoint });
  }

  function handleDuplicate() {
    dispatch('duplicate', { endpoint });
  }

  function handleTest() {
    dispatch('test', { endpoint });
  }

  function formatLastUsed(date: Date | undefined): string {
    if (!date) return 'Never';
    
    const now = new Date();
    const diff = now.getTime() - date.getTime();
    const minutes = Math.floor(diff / 60000);
    const hours = Math.floor(diff / 3600000);
    const days = Math.floor(diff / 86400000);
    
    if (minutes < 1) return 'Just now';
    if (minutes < 60) return `${minutes}m ago`;
    if (hours < 24) return `${hours}h ago`;
    if (days < 7) return `${days}d ago`;
    return date.toLocaleDateString();
  }

  function getStatusColor(status: string): string {
    switch (status) {
      case 'active': return 'text-green-400';
      case 'inactive': return 'text-gray-400';
      case 'error': return 'text-red-400';
      default: return 'text-gray-400';
    }
  }

  function getStatusBgColor(status: string): string {
    switch (status) {
      case 'active': return 'bg-green-400/20';
      case 'inactive': return 'bg-gray-400/20';
      case 'error': return 'bg-red-400/20';
      default: return 'bg-gray-400/20';
    }
  }
</script>

<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-sm hover:shadow-md transition-all duration-200 {className}">
  <div class="p-4 {compact ? 'py-3' : ''}">
    <div class="flex items-center justify-between">
      <!-- Left side: Method, Path, and Details -->
      <div class="flex items-center space-x-3 flex-1 min-w-0">
        <!-- HTTP Method Badge -->
        <div class="flex-shrink-0">
          <HttpMethodBadge method={endpoint.method} size={compact ? 'sm' : 'md'} />
        </div>
        
        <!-- Path and Details -->
        <div class="flex-1 min-w-0">
          <div class="flex items-center space-x-2">
            <h3 class="text-sm font-medium text-gray-900 dark:text-white truncate">
              {endpoint.path}
            </h3>
            
            <!-- Status Indicator -->
            <div class="flex items-center space-x-1 flex-shrink-0">
              <div class="w-2 h-2 rounded-full {getStatusBgColor(endpoint.status)}">
                <div class="w-full h-full rounded-full {getStatusColor(endpoint.status)}" 
                     class:animate-pulse={endpoint.status === 'active'}></div>
              </div>
              {#if !compact}
                <span class="text-xs {getStatusColor(endpoint.status)} capitalize">
                  {endpoint.status}
                </span>
              {/if}
            </div>

            <!-- Default Badge -->
            {#if endpoint.isDefault}
              <span class="px-2 py-1 bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-400 text-xs font-medium rounded">
                Default
              </span>
            {/if}
          </div>
          
          <!-- Name and Description -->
          {#if endpoint.name && !compact}
            <p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
              {endpoint.name}
            </p>
          {/if}
          
          {#if endpoint.description && !compact}
            <p class="text-xs text-gray-500 dark:text-gray-500 mt-1 truncate">
              {endpoint.description}
            </p>
          {/if}
        </div>
      </div>
      
      <!-- Right side: Status Code, Metrics, and Actions -->
      <div class="flex items-center space-x-3 flex-shrink-0">
        <!-- Status Code Badge -->
        {#if showStatusCode && endpoint.statusCode}
          <div class="flex-shrink-0">
            <StatusCodeBadge 
              statusCode={endpoint.statusCode} 
              size={compact ? 'sm' : 'md'} 
              showDescription={!compact}
            />
          </div>
        {/if}
        
        <!-- Metrics -->
        {#if showMetrics && !compact}
          <div class="text-xs text-gray-500 dark:text-gray-400 text-right space-y-1">
            {#if endpoint.responseTime}
              <div class="flex items-center">
                <i class="fas fa-clock mr-1"></i>
                {endpoint.responseTime}ms
              </div>
            {/if}
            {#if endpoint.lastUsed}
              <div class="flex items-center">
                <i class="fas fa-history mr-1"></i>
                {formatLastUsed(endpoint.lastUsed)}
              </div>
            {/if}
          </div>
        {/if}
        
        <!-- Actions -->
        {#if showActions}
          <div class="flex items-center space-x-1">
            <!-- Status Toggle -->
            <button
              class="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 transition-colors rounded"
              title={endpoint.status === 'active' ? 'Disable endpoint' : 'Enable endpoint'}
              on:click={handleToggle}
            >
              <i class="fas {endpoint.status === 'active' ? 'fa-pause' : 'fa-play'} text-sm"></i>
            </button>
            
            <!-- Test Endpoint -->
            <button
              class="p-2 text-gray-400 hover:text-blue-600 dark:text-gray-500 dark:hover:text-blue-400 transition-colors rounded"
              title="Test endpoint"
              on:click={handleTest}
            >
              <i class="fas fa-play-circle text-sm"></i>
            </button>
            
            <!-- More Actions Dropdown -->
            <div class="relative group">
              <button
                class="p-2 text-gray-400 hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300 transition-colors rounded"
                title="More actions"
              >
                <i class="fas fa-ellipsis-v text-sm"></i>
              </button>
              
              <!-- Dropdown Menu -->
              <div class="absolute right-0 top-full mt-1 w-48 bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg shadow-lg opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all duration-200 z-10">
                <div class="py-1">
                  <button
                    class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
                    on:click={handleEdit}
                  >
                    <i class="fas fa-edit mr-2"></i>Edit
                  </button>
                  
                  <button
                    class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
                    on:click={handleDuplicate}
                  >
                    <i class="fas fa-copy mr-2"></i>Duplicate
                  </button>
                  
                  <hr class="border-gray-200 dark:border-gray-600 my-1">
                  
                  <button
                    class="w-full text-left px-4 py-2 text-sm text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors"
                    on:click={handleDelete}
                  >
                    <i class="fas fa-trash-alt mr-2"></i>Delete
                  </button>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
</div>
