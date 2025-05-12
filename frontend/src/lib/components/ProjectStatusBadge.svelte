<script lang="ts">
  import * as ThemeUtils from '$lib/utils/themeUtils';

  export let status: string = 'stopped';
  export let showLabel: boolean = true;
  export let size: 'small' | 'default' | 'large' = 'default';
  
  // Determine size classes
  $: sizeClass = 
    size === 'small' ? 'h-2 w-2' : 
    size === 'large' ? 'h-4 w-4' : 
    'h-3 w-3';
  
  $: textSize = 
    size === 'small' ? 'text-xs' : 
    size === 'large' ? 'text-sm' : 
    'text-xs';
  
  // Determine colors based on status
  $: indicatorColor = 
    status === 'running' ? 'bg-green-500' :
    status === 'error' ? 'bg-red-500' :
    'bg-gray-500';
  
  $: pingColor = 
    status === 'running' ? 'bg-green-400' :
    status === 'error' ? 'bg-red-400' :
    '';
  
  $: textColor = 
    status === 'running' ? 'text-green-400' :
    status === 'error' ? 'text-red-400' :
    'theme-text-muted';
</script>

<div class="flex items-center">
  <span class="relative flex {sizeClass} mr-2">
    {#if status === 'running' || status === 'error'}
      <span class="animate-ping absolute inline-flex h-full w-full rounded-full {pingColor} opacity-75"></span>
    {/if}
    <span class="relative inline-flex rounded-full {sizeClass} {indicatorColor}"></span>
  </span>
  
  {#if showLabel}
    <span class="{textColor} {textSize}">
      {status.toUpperCase()}
    </span>
  {/if}
</div>
