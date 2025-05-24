<script lang="ts">
  /**
   * A reusable skeleton loader component for displaying loading states
   * when fetching data from APIs.
   * 
   * @component
   * 
   * @prop {string} type - The type of skeleton to display: 'card', 'list', 'table', 'text', or 'custom'
   * @prop {number} count - Number of skeleton items to display (for list and table types)
   * @prop {string} height - Custom height for the skeleton
   * @prop {string} width - Custom width for the skeleton
   * @prop {string} className - Additional CSS classes to apply to the skeleton
   */
  export let type: 'card' | 'list' | 'table' | 'text' | 'custom' = 'card';
  export let count: number = 3;
  export let height: string = '';
  export let width: string = '';
  export let className: string = '';

  // Determine skeleton dimensions based on type
  let dimensions = '';
  if (height && width) {
    dimensions = `style="height: ${height}; width: ${width};"`;
  } else if (height) {
    dimensions = `style="height: ${height};"`;
  } else if (width) {
    dimensions = `style="width: ${width};"`;
  }

  // Define some preset animations for different skeleton types
  const getAnimation = () => {
    return 'animate-pulse';
  };
</script>

{#if type === 'card'}
  <div class="theme-bg-secondary rounded-md shadow {getAnimation()} {className}" {dimensions || 'style="height: 200px;"'}>
    <div class="h-40 rounded-t-md bg-gray-300 dark:bg-gray-600"></div>
    <div class="p-4 space-y-3">
      <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/4"></div>
      <div class="h-3 bg-gray-300 dark:bg-gray-700 rounded w-1/2"></div>
      <div class="h-3 bg-gray-300 dark:bg-gray-700 rounded w-5/6"></div>
    </div>
  </div>
{:else if type === 'list'}
  <div class="space-y-4 {className}" {dimensions}>
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
  <div class="theme-bg-secondary rounded-md shadow overflow-hidden {className}" {dimensions}>
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
  <div class="space-y-2 {className}" {dimensions}>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-full {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-5/6 {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-4/6 {getAnimation()}"></div>
    <div class="h-4 bg-gray-300 dark:bg-gray-700 rounded w-3/6 {getAnimation()}"></div>
  </div>
{:else}
  <!-- Custom skeleton - just a simple animated placeholder -->
  <div class="rounded-md bg-gray-300 dark:bg-gray-700 {getAnimation()} {className}" {dimensions || 'style="height: 100px;"'}></div>
{/if}
