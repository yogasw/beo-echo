<!-- HttpMethodBadge - Reusable HTTP method badge component with consistent colors -->
<script lang="ts">
  import { getHttpMethodColor, getHttpMethodDescription, getBadgeSizeClasses, getBaseBadgeClasses } from '$lib/utils/badgeUtils';

  export let method: string;
  export let size: 'xs' | 'sm' | 'md' | 'lg' = 'sm';
  export let showDescription: boolean = false;
  export let className: string = '';

  // Get method styling and information
  $: methodColors = getHttpMethodColor(method);
  $: methodDescription = getHttpMethodDescription(method);
  $: sizeClasses = getBadgeSizeClasses(size);
  $: baseClasses = getBaseBadgeClasses();
  $: displayMethod = method.toUpperCase();
</script>

<div class="inline-flex items-center {className}">
  <!-- HTTP Method Badge -->
  <span class="{baseClasses} {methodColors.bgColor} {sizeClasses}">
    {displayMethod}
  </span>
  
  <!-- Optional Description -->
  {#if showDescription}
    <span class="ml-2 text-sm text-gray-600 dark:text-gray-400">
      {methodDescription}
    </span>
  {/if}
</div>
