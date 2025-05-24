<script lang="ts">
  /**
   * A reusable component for displaying error messages
   * 
   * @component
   * 
   * @prop {string} message - The error message to display
   * @prop {string} type - The type of error: 'error', 'warning', or 'info'
   * @prop {boolean} retryable - Whether the error can be retried
   * @prop {Function} onRetry - Callback function when retry button is clicked
   * @prop {string} className - Additional CSS classes to apply
   */
  export let message: string = 'An error occurred';
  export let type: 'error' | 'warning' | 'info' = 'error';
  export let retryable: boolean = false;
  export let onRetry: () => void = () => {};
  export let className: string = '';

  // Determine icon and colors based on error type
  $: icon = type === 'error' ? 'fas fa-exclamation-circle' : 
            type === 'warning' ? 'fas fa-exclamation-triangle' : 
            'fas fa-info-circle';
  
  $: bgColor = type === 'error' ? 'bg-red-500/20 border-red-500' : 
               type === 'warning' ? 'bg-yellow-500/20 border-yellow-500' : 
               'bg-blue-500/20 border-blue-500';
  
  $: textColor = type === 'error' ? 'text-red-400' : 
                 type === 'warning' ? 'text-yellow-400' : 
                 'text-blue-400';
</script>

<div class={`p-4 rounded-md border ${bgColor} ${className}`}>
  <div class="flex items-start">
    <div class="flex-shrink-0">
      <i class={`${icon} ${textColor} text-lg`}></i>
    </div>
    <div class="ml-3 flex-1">
      <p class={`text-sm ${textColor}`}>{message}</p>
      
      {#if retryable}
        <div class="mt-3">
          <button 
            on:click={onRetry}
            class="text-sm px-3 py-1.5 rounded-md bg-gray-700 hover:bg-gray-600 text-white"
          >
            <i class="fas fa-sync-alt mr-1.5"></i>
            Retry
          </button>
        </div>
      {/if}
    </div>
  </div>
</div>
