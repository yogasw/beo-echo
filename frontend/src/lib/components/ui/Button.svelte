<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  const dispatch = createEventDispatcher();
  
  export let type: 'submit' | 'button' = 'button';
  export let variant: 'primary' | 'secondary' | 'outline' = 'primary';
  export let disabled = false;
  export let loading = false;
  export let fullWidth = false;
  export let onClick: (() => void) | undefined = undefined;

  const baseClasses = "flex items-center justify-center gap-2 p-3 rounded-lg border transition-colors duration-200 font-medium focus:outline-none focus:ring-2 focus:ring-offset-2";
  
  const variants = {
    primary: `${baseClasses}
      dark:bg-gray-800 dark:hover:bg-gray-700 dark:border-gray-600 dark:text-white
      bg-white hover:bg-gray-50 border-gray-300 text-gray-700
      disabled:bg-gray-100 dark:disabled:bg-gray-800/50
      disabled:border-gray-200 dark:disabled:border-gray-700
      disabled:text-gray-400 dark:disabled:text-gray-500
      disabled:cursor-not-allowed
      focus:ring-gray-500`,
    secondary: `${baseClasses}
      dark:bg-blue-600 dark:hover:bg-blue-700 dark:border-blue-500 dark:text-white
      bg-blue-500 hover:bg-blue-600 border-blue-400 text-white
      disabled:bg-blue-300 dark:disabled:bg-blue-500/50
      disabled:border-blue-200 dark:disabled:border-blue-400
      disabled:text-blue-100 dark:disabled:text-blue-200
      disabled:cursor-not-allowed
      focus:ring-blue-500`,
    outline: `${baseClasses}
      dark:bg-transparent dark:hover:bg-gray-800 dark:border-gray-600 dark:text-white
      bg-transparent hover:bg-gray-50 border-gray-300 text-gray-700
      disabled:bg-transparent dark:disabled:bg-transparent
      disabled:border-gray-200 dark:disabled:border-gray-700
      disabled:text-gray-400 dark:disabled:text-gray-500
      disabled:cursor-not-allowed
      focus:ring-gray-500`
  };
</script>

<button
  {type}
  class={`${variants[variant]} ${fullWidth ? 'w-full' : ''}`}
  on:click={onClick}
  {disabled}
  {...$$restProps}
>
  {#if loading}
    <svg
      class="animate-spin -ml-1 mr-2 h-4 w-4"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
  {/if}
  <slot />
</button>
