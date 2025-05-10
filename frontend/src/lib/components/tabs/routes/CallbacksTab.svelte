<script lang="ts">
  export let callbacks: any[] = [];
</script>

<div>
  <h3 class="text-sm font-semibold mb-2">Callbacks</h3>
  {#if callbacks && callbacks.length > 0}
    <pre class="bg-gray-800 p-2 rounded text-xs">{JSON.stringify(callbacks, null, 2)}</pre>
  {:else}
    <div class="text-gray-400">No callbacks available.</div>
  {/if}
</div> 