<script lang="ts">
  import { featureToggles } from '$lib/stores/featureToggles';
  
  export let showDebug = false;
</script>

{#if showDebug}
  <div class="bg-gray-900/50 p-4 rounded-lg border border-gray-600 my-4">
    <h3 class="text-sm font-semibold text-white mb-2 flex items-center">
      <i class="fas fa-flag mr-2 text-blue-400"></i>
      Active Feature Flags
    </h3>
    <div class="text-xs">
      {#each Object.entries($featureToggles) as [key, value]}
        <div class="flex justify-between items-center mb-1">
          <span class="text-gray-300">{key}:</span>
          {#if value}
            <span class="text-green-400 font-mono">true</span>
          {:else}
            <span class="text-gray-500 font-mono">false</span>
          {/if}
        </div>
      {/each}
    </div>
  </div>
{/if}
