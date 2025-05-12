<script lang="ts">
  import * as ThemeUtils from '$lib/utils/themeUtils';
  export let callbacks: any[] = [];
</script>

<div>
  <h3 class="text-sm font-semibold mb-2 {ThemeUtils.themeTextPrimary()}">Callbacks</h3>
  {#if callbacks && callbacks.length > 0}
    <pre class="{ThemeUtils.themeBgSecondary()} p-2 rounded text-xs {ThemeUtils.themeTextPrimary()}">{JSON.stringify(callbacks, null, 2)}</pre>
  {:else}
    <div class="{ThemeUtils.themeTextMuted()}">No callbacks available.</div>
  {/if}
</div> 