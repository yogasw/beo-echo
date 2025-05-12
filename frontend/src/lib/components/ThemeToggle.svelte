<script lang="ts">
  import { theme, toggleTheme } from '$lib/stores/theme';
  
  // Props for customizing appearance
  export let showLabel = false;
  export let size = 'default'; // 'small', 'default', 'large'
  
  // Compute size class based on the prop
  $: sizeClass = size === 'small' ? 'w-8 h-4' : 
                 size === 'large' ? 'w-14 h-7' : 'w-11 h-6';
                 
  // Calculate appropriate height and position for the toggle knob
  $: knobClass = size === 'small' ? 'h-3 w-3 after:h-3 after:w-3' : 
                 size === 'large' ? 'h-6 w-6 after:h-6 after:w-6' : 'h-5 w-5 after:h-5 after:w-5';
                 
  // Calculate top position
  $: topPosition = size === 'small' ? 'top-[2px]' : 
                   size === 'large' ? 'top-[2px]' : 'top-[2px]';
</script>

<div class="flex items-center">
  {#if showLabel}
    <span class="mr-2 text-sm theme-text-secondary">
      {$theme === 'dark' ? 'Dark' : 'Light'}
    </span>
  {/if}
  
  <label class="inline-flex items-center cursor-pointer">
    <input 
      type="checkbox" 
      checked={$theme === 'dark'} 
      on:change={toggleTheme} 
      class="sr-only peer"
    />
    <div class="{sizeClass} bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
      rounded-full peer peer-focus:outline-none peer-focus:ring-2 
      peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 
      peer-checked:after:translate-x-full peer-checked:after:border-white 
      after:content-[''] after:absolute after:{topPosition} after:start-[2px] 
      after:bg-white after:border-gray-300 dark:after:border-gray-600 
      after:border after:rounded-full after:{knobClass}
      after:transition-all dark:border-gray-600 relative"
    >
    </div>
    
    {#if showLabel}
      <i class="ml-2 {$theme === 'dark' ? 'fas fa-moon text-blue-400' : 'fas fa-sun text-yellow-400'}"></i>
    {/if}
  </label>
</div>
