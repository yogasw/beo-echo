<!-- StatusCodeBadge - Reusable status code badge component with consistent colors -->
<script lang="ts">
  import { getStatusCodeColor, getStatusCodeInfo, getBadgeSizeClasses, getBaseBadgeClasses } from '$lib/utils/badgeUtils';

  export let statusCode: number;
  export let size: 'xs' | 'sm' | 'md' | 'lg' = 'sm';
  export let showDescription: boolean = false;
  export let className: string = '';

  // Get status code styling and information
  $: statusColors = getStatusCodeColor(statusCode);
  $: statusInfo = getStatusCodeInfo(statusCode);
  $: sizeClasses = getBadgeSizeClasses(size);
  $: baseClasses = getBaseBadgeClasses();
</script>

<div class="inline-flex items-center {className}">
  <!-- Status Code Badge -->
  <span class="{baseClasses} {statusColors.bgColor} {sizeClasses}">
    {statusCode}
  </span>
  
  <!-- Optional Description -->
  {#if showDescription}
    <span class="ml-2 text-sm text-gray-600 dark:text-gray-400">
      {#if statusInfo}
        {statusInfo.name}
        {#if statusInfo.description}
          - {statusInfo.description}
        {/if}
      {:else}
        Custom {statusColors.category.toLowerCase()} status code
      {/if}
    </span>
  {/if}
</div>
