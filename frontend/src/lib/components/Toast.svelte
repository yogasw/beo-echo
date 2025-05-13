<script lang="ts">
  import { toast } from '$lib/stores/toast';
  import { theme } from '$lib/stores/theme';
  import { onDestroy } from 'svelte';
  import * as ThemeUtils from '$lib/utils/themeUtils';

  let visible = true;
  let timeoutId: ReturnType<typeof setTimeout>;

  function close() {
    visible = false;
    toast.clear();
  }

  $: if ($toast) {
    visible = true;
    const duration = $toast.duration || 5000; // Default to 5 seconds if not specified
    clearTimeout(timeoutId);
    timeoutId = setTimeout(close, duration);
  }

  onDestroy(() => {
    clearTimeout(timeoutId);
  });

  // Generate theme-responsive classes for the toast
  $: toastClasses = `toast ${$toast?.type || ''} ${$theme === 'dark' ? 'dark' : 'light'}`;
</script>

{#if $toast}
  <div class={toastClasses} class:visible>
    <span class="icon">
      {#if $toast.type === 'success'}<i class="fas fa-check-circle text-green-500 dark:text-green-400"></i>{/if}
      {#if $toast.type === 'error'}<i class="fas fa-times-circle text-red-600 dark:text-red-400"></i>{/if}
      {#if $toast.type === 'warning'}<i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400"></i>{/if}
      {#if $toast.type === 'info'}<i class="fas fa-info-circle text-blue-600 dark:text-blue-400"></i>{/if}
    </span>
    <span class="message">{$toast.message}</span>
    <button class="close-btn" on:click={close} aria-label="Close notification">
      <i class="fas fa-times"></i>
    </button>
  </div>
{/if}

<style>
  .toast {
    position: fixed;
    left: 50%;
    bottom: 2rem;
    transform: translateX(-50%);
    min-width: 320px;
    max-width: 90vw;
    display: flex;
    align-items: center;
    border-radius: 0.75rem;
    box-shadow: 0 2px 16px rgba(0,0,0,0.2);
    padding: 1rem 1.5rem;
    z-index: 1000;
    font-size: 1.1rem;
    gap: 0.75rem;
    border: 1.5px solid transparent;
    animation: fadeIn 0.2s;
  }
  /* Dark theme (default) */
  .toast.dark {
    background: #1e293b; /* bg-gray-800 */
    color: #fff;
  }
  /* Light theme */
  .toast.light {
    background: #ffffff; /* bg-white */
    color: #1e293b; /* text-gray-800 */
    box-shadow: 0 2px 16px rgba(0,0,0,0.1);
  }
  /* Success variant */
  .toast.success.dark {
    border-color: #22c55e; /* border-green-500 */
  }
  .toast.success.light {
    border-color: #16a34a; /* border-green-600 */
  }
  /* Error variant */
  .toast.error.dark {
    border-color: #ef4444; /* border-red-500 */
  }
  .toast.error.light {
    border-color: #dc2626; /* border-red-600 */
  }
  /* Warning variant */
  .toast.warning.dark {
    border-color: #f59e42; /* border-yellow-500 */
  }
  .toast.warning.light {
    border-color: #ca8a04; /* border-yellow-600 */
  }
  /* Info variant */
  .toast.info.dark {
    border-color: #3b82f6; /* border-blue-500 */
  }
  .toast.info.light {
    border-color: #2563eb; /* border-blue-600 */
  }
  
  .toast .icon {
    font-size: 1.5rem;
    display: flex;
    align-items: center;
  }
  
  .toast .message {
    flex: 1;
  }
  
  .toast .close-btn {
    background: none;
    border: none;
    font-size: 1.2rem;
    cursor: pointer;
    margin-left: 0.5rem;
    opacity: 0.7;
    transition: opacity 0.2s;
  }
  
  .toast.dark .close-btn {
    color: #fff;
  }
  
  .toast.light .close-btn {
    color: #1e293b;
  }
  
  .toast .close-btn:hover {
    opacity: 1;
  }
  
  @keyframes fadeIn {
    from { opacity: 0; transform: translateX(-50%) translateY(20px); }
    to { opacity: 1; transform: translateX(-50%) translateY(0); }
  }
</style>
