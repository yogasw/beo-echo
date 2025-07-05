<script lang="ts">
  /**
   * A reusable skeleton loader component for displaying loading states
   * when fetching data from APIs. Features smart loading with delay to 
   * prevent jarring flash for quick operations.
   * 
   * @component
   * 
   * @prop {string} type - The type of skeleton to display: 'card', 'list', 'table', 'text', or 'custom'
   * @prop {number} count - Number of skeleton items to display (for list and table types)
   * @prop {string} height - Custom height for the skeleton
   * @prop {string} width - Custom width for the skeleton
   * @prop {string} className - Additional CSS classes to apply to the skeleton
   * @prop {number} delay - Delay before showing skeleton in ms (default: 400)
   * @prop {number} minShowTime - Minimum time to show skeleton once displayed (default: 300)
   * @prop {boolean} isLoading - External loading state (default: true)
   */
  export let type: 'card' | 'list' | 'table' | 'text' | 'custom' = 'card';
  export let count: number = 3;
  export let height: string = '';
  export let width: string = '';
  export let className: string = '';
  export let delay: number = 400;
  export let minShowTime: number = 300;
  export let isLoading: boolean = true;

  import { onDestroy } from 'svelte';

  // Smart loading state management
  let showSkeleton = false;
  let loadingStartTime: number | null = null;
  let showTimeout: ReturnType<typeof setTimeout> | null = null;
  let hideTimeout: ReturnType<typeof setTimeout> | null = null;

  // Determine skeleton dimensions based on type
  let styleString = '';
  if (height && width) {
    styleString = `height: ${height}; width: ${width};`;
  } else if (height) {
    styleString = `height: ${height};`;
  } else if (width) {
    styleString = `width: ${width};`;
  }

  // Define some preset animations for different skeleton types
  const getAnimation = () => {
    return 'animate-pulse';
  };

  // Watch for loading state changes with smart delay
  $: handleLoadingChange(isLoading);

  function handleLoadingChange(loading: boolean) {
    if (loading) {
      startSmartLoading();
    } else {
      stopSmartLoading();
    }
  }

  function startSmartLoading() {
    // Clear any existing timeouts
    clearTimeouts();
    
    loadingStartTime = Date.now();
    showSkeleton = false;

    // Set timeout to show skeleton after delay
    showTimeout = setTimeout(() => {
      if (isLoading) { // Only show if still loading
        showSkeleton = true;
      }
    }, delay);
  }

  function stopSmartLoading() {
    // Clear show timeout if loading finished before delay
    if (showTimeout) {
      clearTimeout(showTimeout);
      showTimeout = null;
    }

    if (!showSkeleton) {
      // Skeleton was never shown, just stop immediately
      return;
    }

    // Calculate how long the skeleton has been visible
    const now = Date.now();
    const totalLoadingTime = loadingStartTime ? now - loadingStartTime : 0;
    const skeletonVisibleTime = Math.max(0, totalLoadingTime - delay);
    const remainingMinTime = Math.max(0, minShowTime - skeletonVisibleTime);

    if (remainingMinTime > 0) {
      // Keep showing skeleton for remaining minimum time
      hideTimeout = setTimeout(() => {
        showSkeleton = false;
        loadingStartTime = null;
      }, remainingMinTime);
    } else {
      // Can hide immediately
      showSkeleton = false;
      loadingStartTime = null;
    }
  }

  function clearTimeouts() {
    if (showTimeout) {
      clearTimeout(showTimeout);
      showTimeout = null;
    }
    if (hideTimeout) {
      clearTimeout(hideTimeout);
      hideTimeout = null;
    }
  }

  onDestroy(() => {
    clearTimeouts();
  });

  // This makes TypeScript recognize the component has a default export
  // export default;
</script>

{#if showSkeleton}
{#if type === 'card'}
  <div class="theme-bg-secondary rounded-md shadow {getAnimation()} {className}" style={styleString || 'height: 200px;'}>
    <div class="h-40 rounded-t-md bg-gray-300 dark:bg-gray-600"></div>
    <div class="p-4 space-y-3">
      <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/4"></div>
      <div class="h-3 bg-gray-300 dark:bg-gray-700 rounded w-1/2"></div>
      <div class="h-3 bg-gray-300 dark:bg-gray-700 rounded w-5/6"></div>
    </div>
  </div>
{:else if type === 'list'}
  <div class="space-y-4 {className}" style={styleString}>
    {#each Array(count) as _, i}
      <div class="flex items-center p-4 theme-bg-secondary rounded-md shadow {getAnimation()}">
        <div class="w-12 h-12 rounded-full bg-gray-300 dark:bg-gray-600 mr-4"></div>
        <div class="space-y-2 flex-1">
          <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/4"></div>
          <div class="h-3 bg-gray-300 dark:bg-gray-700 rounded w-1/2"></div>
        </div>
      </div>
    {/each}
  </div>
{:else if type === 'table'}
  <div class="theme-bg-secondary rounded-md shadow overflow-hidden {className}" style={styleString}>
    <!-- Table header -->
    <div class="bg-gray-100 dark:bg-gray-750 p-4">
      <div class="h-6 bg-gray-300 dark:bg-gray-600 rounded w-full {getAnimation()}"></div>
    </div>
    
    <!-- Table rows -->
    <div class="divide-y theme-border">
      {#each Array(count) as _, i}
        <div class="p-4 {getAnimation()}">
          <div class="grid grid-cols-4 gap-4">
            <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded"></div>
            <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded"></div>
            <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded"></div>
            <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded"></div>
          </div>
        </div>
      {/each}
    </div>
  </div>
{:else if type === 'text'}
  <div class="space-y-2 {className}" style={styleString}>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-5/6 {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-4/6 {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/6 {getAnimation()}"></div>
  </div>
{:else}
  <!-- Custom skeleton - just a simple animated placeholder -->
  <div class="rounded-md bg-gray-300 dark:bg-gray-700 {getAnimation()} {className}" style={styleString || 'height: 100px;'}></div>
{/if}
{/if}
