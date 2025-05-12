<script lang="ts">
  import * as ThemeUtils from '$lib/utils/themeUtils';
  import { theme } from '$lib/stores/theme';
  export let callbacks: any[] = [];
</script>

<div class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full">
  <div class="{ThemeUtils.themeBgSecondary()} rounded-lg p-4">
    <h3 class="text-sm font-semibold mb-3 {ThemeUtils.themeTextPrimary()}">Callbacks</h3>
    {#if callbacks && callbacks.length > 0}
      <pre class="p-2 rounded text-xs {ThemeUtils.themeTextPrimary()}">{JSON.stringify(callbacks, null, 2)}</pre>
    {:else}
      <div class="{ThemeUtils.themeTextMuted()}">No callbacks available.</div>
    {/if}
  </div>
</div> 