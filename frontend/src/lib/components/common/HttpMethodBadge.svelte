<!-- HttpMethodBadge - Reusable HTTP method badge component with consistent colors -->
<script lang="ts">
  import { getHttpMethodColor, getHttpMethodDescription, getBadgeSizeClasses, getBaseBadgeClasses } from '$lib/utils/badgeUtils';

  let {
    method,
    size = 'sm',
    showDescription = false,
    className = ''
  }: {
    method: string;
    size?: 'xs' | 'sm' | 'md' | 'lg';
    showDescription?: boolean;
    className?: string;
  } = $props();

  // Get method styling and information
  let methodColors = $derived(getHttpMethodColor(method));
  let methodDescription = $derived(getHttpMethodDescription(method));
  let sizeClasses = $derived(getBadgeSizeClasses(size));
  let baseClasses = $derived(getBaseBadgeClasses());
  let displayMethod = $derived(method.toUpperCase());
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
