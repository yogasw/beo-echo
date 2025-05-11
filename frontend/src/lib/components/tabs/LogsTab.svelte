<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { createLogStream, getLogs, type Project, type RequestLog } from "$lib/api/mockoonApi";

  export let selectedProject: Project;

  let logs: RequestLog[] = [];
  let isLoading = true;
  let error: string | null = null;
  let searchTerm = '';
  let total = 0;
  let page = 1;
  const pageSize = 100;
  let eventSource: EventSource | null = null;
  let autoScroll = true;
  // Map to track expanded logs
  let expandedLogs: Record<string, boolean> = {};
  
  // Function to toggle log expansion
  function toggleLogExpansion(logId: string) {
    expandedLogs[logId] = !expandedLogs[logId];
    expandedLogs = expandedLogs; // Force Svelte reactivity update
  }
  
  $: filteredLogs = searchTerm 
    ? logs.filter(log => 
        log.path.toLowerCase().includes(searchTerm.toLowerCase()) ||
        log.method.toLowerCase().includes(searchTerm.toLowerCase()) ||
        log.request_body.toLowerCase().includes(searchTerm.toLowerCase()) ||
        log.response_body.toLowerCase().includes(searchTerm.toLowerCase())
      )
    : logs;
  
  // Convert JSON string to object for display
  function parseJson(jsonString: string): any {
    try {
      return JSON.parse(jsonString);
    } catch (e) {
      return jsonString;
    }
  }
  
  // Format timestamp for display
  function formatDate(dateString: string | Date): string {
    try {
      const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
      return date.toLocaleString();
    } catch (e) {
      return String(dateString);
    }
  }
  
  async function loadInitialLogs() {
    try {
      isLoading = true;
      const result = await getLogs(1, pageSize, selectedProject.id);
      logs = result.logs;
      total = result.total;
      isLoading = false;
    } catch (err) {
      console.error('Failed to load logs:', err);
      error = 'Failed to load logs: ' + (err instanceof Error ? err.message : String(err));
      isLoading = false;
    }
  }
  
  function setupLogStream() {
    // Close any existing connection
    if (eventSource) {
      eventSource.close();
    }
    
    console.log('Setting up log stream for project:', selectedProject.id);
    
    // Create new connection
    eventSource = createLogStream(selectedProject.id, pageSize);
    
    // Setup event handlers
    eventSource.addEventListener('log', (event) => {
      try {
        console.log('Log event received:', event.data);
        const newLog = JSON.parse(event.data);
        
        // Check if log already exists to prevent duplicates
        if (!logs.some(log => log.id === newLog.id)) {
          // Add to beginning of array (newest first) and force Svelte reactivity
          logs = [newLog, ...logs].slice(0, 1000); // Limit to 1000 logs to prevent browser slowdown
          console.log('Added new log, total logs:', logs.length);
          
          // Auto-scroll to top if enabled
          if (autoScroll) {
            window.scrollTo(0, 0);
          }
        }
      } catch (err) {
        console.error('Error processing log event:', err, event.data);
      }
    });
    
    // Direct message event (fallback)
    eventSource.onmessage = (event) => {
      console.log('Generic message received:', event.data);
      try {
        const newLog = JSON.parse(event.data);
        if (newLog && newLog.id && !logs.some(log => log.id === newLog.id)) {
          logs = [newLog, ...logs].slice(0, 1000);
        }
      } catch (err) {
        console.error('Error processing generic message:', err);
      }
    };
    
    eventSource.addEventListener('ping', (event) => {
      // Keep connection alive, no action needed
      console.log('Ping received from server:', event.data);
      isConnected = true;
      reconnectAttempts = 0; // Reset reconnect counter on successful ping
    });
    
    eventSource.onopen = () => {
      console.log('Log stream connection established');
      isConnected = true;
      reconnectAttempts = 0; // Reset reconnect counter on connection
    };
    
    eventSource.onerror = (err) => {
      console.error('EventSource error:', err);
      isConnected = false;
      
      // Implement smart reconnection strategy with backoff
      if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
        reconnectAttempts++;
        const delay = RECONNECT_DELAY_MS * reconnectAttempts; // Increase delay with each attempt
        console.log(`Attempting to reconnect log stream (attempt ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS}) in ${delay}ms...`);
        
        setTimeout(() => {
          setupLogStream();
        }, delay);
      } else {
        console.error('Max reconnection attempts reached. Please refresh manually.');
      }
    };
  }
  
  // Track connection status for UI feedback
  let isConnected = false;
  let reconnectAttempts = 0;
  const MAX_RECONNECT_ATTEMPTS = 5;
  const RECONNECT_DELAY_MS = 3000;

  // Initialize on component mount
  onMount(() => {
    loadInitialLogs().then(() => {
      setupLogStream();
    });
  });
  
  // Clean up on component destroy
  onDestroy(() => {
    if (eventSource) {
      eventSource.close();
    }
  });
</script>

<div class="w-full bg-gray-800 p-4">
  {#if !isConnected && reconnectAttempts > 0}
    <div class="p-2 rounded mb-4 flex items-center justify-between">
      <div class="flex items-center">
        <i class="fas fa-exclamation-triangle text-yellow-400 text-lg mr-2"></i>
        <span class="text-white">Live stream disconnected. Using manual refresh only.</span>
      </div>
      <button 
        class="bg-blue-500 hover:bg-blue-600 text-white py-1 px-3 rounded text-sm"
        on:click={() => setupLogStream()}
      >
        <i class="fas fa-sync mr-1"></i> Reconnect Stream
      </button>
    </div>
  {/if}
  <div class="bg-gray-700 p-4 rounded mb-4 flex items-center justify-between">
    <div class="flex items-center">
      <i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
      <span class="text-xl font-bold text-blue-500">
        Logs for: {selectedProject.name}
      </span>
      
      <div class="ml-4 flex items-center">
        <!-- Stream status indicator -->
        <span class="relative flex h-3 w-3">
          {#if isConnected}
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
            <span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
          {:else}
            <span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
          {/if}
        </span>
        <span class="ml-2 text-xs {isConnected ? 'text-green-400' : 'text-red-400'}">
          {isConnected ? 'Live stream' : 'Disconnected'}
        </span>
        {#if !isConnected}
          <button 
            class="ml-2 bg-blue-500 hover:bg-blue-600 text-white py-0.5 px-2 rounded text-xs"
            on:click={() => setupLogStream()}
          >
            Reconnect
          </button>
        {/if}
      </div>
    </div>
    <div class="flex items-center">
      <span class="text-xs text-gray-400 mr-2">Auto-scroll</span>
      <input type="checkbox" bind:checked={autoScroll} class="form-checkbox h-4 w-4 text-blue-500" />
      
      <button 
        class="ml-4 bg-blue-500 hover:bg-blue-600 text-white py-1 px-3 rounded text-xs flex items-center"
        on:click={() => {
          loadInitialLogs();
          if (!isConnected) {
            setupLogStream(); // Also try to reconnect if disconnected
          }
        }}
      >
        <i class="fas fa-sync mr-1"></i> Manual Refresh
      </button>
    </div>
  </div>
  

  
  <div class="flex items-center bg-gray-700 p-2 rounded mb-4">
    <i class="fas fa-search text-white text-lg mr-2"></i>
    <input
      type="text"
      bind:value={searchTerm}
      placeholder="Search Logs (path, method, request/response body)"
      class="w-full bg-gray-700 text-white py-1 px-2 rounded text-sm"
    />
  </div>
  
  {#if isLoading}
    <div class="flex justify-center py-8">
      <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
    </div>
  {:else if error}
    <div class="bg-red-800 p-4 rounded mb-4 text-center">
      <p class="text-white">{error}</p>
      <button 
        on:click={loadInitialLogs} 
        class="mt-2 bg-blue-500 hover:bg-blue-600 text-white py-1 px-4 rounded text-sm"
      >
        Retry
      </button>
    </div>
  {:else if filteredLogs.length === 0}
    <div class="bg-gray-700 p-4 rounded text-center">
      <p class="text-white">No logs found {searchTerm ? 'matching your search criteria' : 'for this project'}</p>
    </div>
  {:else}
    <div class="space-y-4">
      {#each filteredLogs as log (log.id)}
          
      <!-- class:bg-green-900={log.matched} 
          class:bg-red-900={!log.matched} -->
          
        <div 
          class="bg-gray-700 p-4 rounded cursor-pointer" 
          on:click={() => toggleLogExpansion(log.id)}
          on:keydown={(e) => e.key === 'Enter' && toggleLogExpansion(log.id)}
          tabindex="0"
          role="button"
          aria-expanded={!!expandedLogs[log.id]}
        >
          <div class="flex justify-between items-center">
            <div class="flex items-center">
              <span class="text-sm font-bold">
                <span class="mr-2 px-2 py-0.5 rounded {log.method === 'GET' ? 'bg-green-600' : log.method === 'POST' ? 'bg-blue-600' : log.method === 'PUT' ? 'bg-yellow-600' : log.method === 'DELETE' ? 'bg-red-600' : 'bg-gray-600'}">
                  {log.method}
                </span>
                {log.path}
              </span>
              <span class="ml-2 text-xs px-2 py-0.5 rounded {log.response_status < 300 ? 'bg-green-600' : log.response_status < 400 ? 'bg-blue-600' : log.response_status < 500 ? 'bg-yellow-600' : 'bg-red-600'}">
                {log.response_status}
              </span>
              <span class="ml-2 text-xs px-2 py-0.5 rounded bg-purple-600">
                {log.execution_mode}
              </span>
              <span class="ml-2 text-xs px-2 py-0.5 rounded {log.matched ? 'bg-green-600' : 'bg-red-600'}">
                {log.matched ? 'Matched' : 'Unmatched'}
              </span>
            </div>
            <div class="flex items-center">
              <span class="text-xs text-gray-400 mr-2">{formatDate(log.created_at)}</span>
              <span class="text-xs px-2 py-0.5 rounded bg-blue-600">{log.latency_ms}ms</span>
              <i class="fas {expandedLogs[log.id] ? 'fa-chevron-up' : 'fa-chevron-down'} ml-3 text-gray-400"></i>
            </div>
          </div>
          
          {#if expandedLogs[log.id]}
            <div class="grid grid-cols-2 gap-4 mt-4">
              <div>
                <div class="flex justify-between items-center mb-2">
                  <h3 class="text-sm font-semibold">Request</h3>
                  {#if log.query_params}
                    <span class="text-xs text-gray-400">Query: {log.query_params}</span>
                  {/if}
                </div>
                
                <div class="mb-2">
                  <h4 class="text-xs font-semibold text-gray-400">Headers</h4>
                  <pre class="bg-gray-800 p-2 rounded text-xs overflow-auto max-h-32">{JSON.stringify(parseJson(log.request_headers), null, 2)}</pre>
                </div>
                
                {#if log.request_body}
                  <div>
                    <h4 class="text-xs font-semibold text-gray-400">Body</h4>
                    <pre class="bg-gray-800 p-2 rounded text-xs overflow-auto max-h-64">{JSON.stringify(parseJson(log.request_body), null, 2)}</pre>
                  </div>
                {/if}
              </div>
              
              <div>
                <div class="flex justify-between items-center mb-2">
                  <h3 class="text-sm font-semibold">Response</h3>
                </div>
                
                <div class="mb-2">
                  <h4 class="text-xs font-semibold text-gray-400">Headers</h4>
                  <pre class="bg-gray-800 p-2 rounded text-xs overflow-auto max-h-32">{JSON.stringify(parseJson(log.response_headers), null, 2)}</pre>
                </div>
                
                <div>
                  <h4 class="text-xs font-semibold text-gray-400">Body</h4>
                  <pre class="bg-gray-800 p-2 rounded text-xs overflow-auto max-h-64">{JSON.stringify(parseJson(log.response_body), null, 2)}</pre>
                </div>
              </div>
            </div>
          {/if}
        </div>
      {/each}
    </div>
    
    {#if filteredLogs.length < total}
      <div class="mt-4 text-center">
        <span class="text-xs text-gray-400">Showing {filteredLogs.length} of {total} logs</span>
      </div>
    {/if}
  {/if}
</div>