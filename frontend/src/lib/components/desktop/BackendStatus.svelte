<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  import { isDesktopMode } from '$lib/utils/desktopConfig';
  import { BASE_URL_API } from '$lib/utils/authUtils';
  
  let backendStatus: 'starting' | 'running' | 'error' | 'stopped' = 'starting';
  let backendPort = 3600;
  let lastChecked: Date | null = null;
  let healthCheckInterval: NodeJS.Timeout | null = null;
  
  onMount(() => {
    if (isDesktopMode()) {
      startHealthCheck();
    }
  });
  
  onDestroy(() => {
    if (healthCheckInterval) {
      clearInterval(healthCheckInterval);
    }
  });
  
  const checkBackendHealth = async () => {
    try {
      const response = await fetch(`${BASE_URL_API}/health`, {
        method: 'GET',
        headers: { 'Content-Type': 'application/json' }
      });
      
      if (response.ok) {
        backendStatus = 'running';
        lastChecked = new Date();
      } else {
        backendStatus = 'error';
      }
    } catch (error) {
      console.log('Backend health check failed:', error);
      backendStatus = 'error';
    }
  };
  
  const startHealthCheck = () => {
    // Initial check
    checkBackendHealth();
    
    // Set up periodic health check every 30 seconds
    healthCheckInterval = setInterval(checkBackendHealth, 30000);
  };
  
  const getStatusColor = (status: string) => {
    switch (status) {
      case 'running':
        return 'text-green-400';
      case 'starting':
        return 'text-yellow-400';
      case 'error':
      case 'stopped':
        return 'text-red-400';
      default:
        return 'text-gray-400';
    }
  };
  
  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'running':
        return 'fas fa-check-circle';
      case 'starting':
        return 'fas fa-sync fa-spin';
      case 'error':
        return 'fas fa-exclamation-circle';
      case 'stopped':
        return 'fas fa-stop-circle';
      default:
        return 'fas fa-question-circle';
    }
  };
  
  const formatLastChecked = (date: Date | null) => {
    if (!date) return 'Never';
    return date.toLocaleTimeString('en-US', { 
      hour12: false, 
      hour: '2-digit', 
      minute: '2-digit',
      second: '2-digit'
    });
  };
</script>

<!-- Desktop Backend Status - Only shown in desktop mode -->
{#if browser && isDesktopMode()}
  <div class="bg-gray-800 border border-gray-700 rounded-md p-3 mb-4">
    <div class="flex items-center justify-between">
      <!-- Status Information -->
      <div class="flex items-center space-x-3">
        <div class="flex items-center space-x-2">
          <i class="{getStatusIcon(backendStatus)} {getStatusColor(backendStatus)}"></i>
          <span class="text-sm font-medium text-white">Backend Service</span>
        </div>
        
        <div class="flex items-center space-x-4 text-xs text-gray-400">
          <span class="capitalize {getStatusColor(backendStatus)}">
            {backendStatus}
          </span>
          <span>Port: {backendPort}</span>
          <span>Last checked: {formatLastChecked(lastChecked)}</span>
        </div>
      </div>
      
      <!-- Actions -->
      <div class="flex items-center space-x-2">
        <button
          class="px-2 py-1 bg-gray-700 hover:bg-gray-600 text-white rounded text-xs transition-colors"
          title="Check backend health"
          aria-label="Check backend service health"
          on:click={checkBackendHealth}
        >
          <i class="fas fa-sync mr-1"></i>
          Check
        </button>
        
        {#if backendStatus === 'error'}
          <span class="px-2 py-1 bg-red-900/20 text-red-300 rounded text-xs">
            <i class="fas fa-exclamation-triangle mr-1"></i>
            Backend Offline
          </span>
        {/if}
      </div>
    </div>
    
    <!-- Additional Info for Desktop -->
    <div class="mt-2 text-xs text-gray-500">
      <div class="flex items-center space-x-4">
        <span>
          <i class="fas fa-desktop mr-1"></i>
          Desktop Mode
        </span>
        <span>
          <i class="fas fa-server mr-1"></i>
          Embedded Backend
        </span>
        <span>
          <i class="fas fa-database mr-1"></i>
          Local SQLite Database
        </span>
      </div>
    </div>
  </div>
{/if}
