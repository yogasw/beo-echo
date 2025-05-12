<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { theme } from '$lib/stores/theme';
	export let headers: string;
	
	// Parse JSON headers string into an array of key-value objects
	$: parsedHeaders = (() => {
		try {
			if (!headers) return [];
			const headersObj = JSON.parse(headers);
			return Object.entries(headersObj).map(([key, value]) => ({ key, value }));
		} catch (error) {
			console.error('Error parsing headers:', error);
			return [];
		}
	})();
</script>

<div class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full">
	<div class="{ThemeUtils.themeBgSecondary()} rounded-lg p-4">
		<h3 class="text-sm font-semibold mb-3 {ThemeUtils.themeTextPrimary()}">Headers</h3>
		{#if parsedHeaders && parsedHeaders.length > 0}
			<ul class="text-xs space-y-2">
				{#each parsedHeaders as header}
					<li class="flex items-start break-all {ThemeUtils.themeTextPrimary()}">
						<span class="font-bold whitespace-nowrap mr-1">{header.key}:</span>
						<span class="break-all">{header.value}</span>
					</li>
				{/each}
			</ul>
		{:else}
			<div class="{ThemeUtils.themeTextMuted()}">No headers available.</div>
		{/if}
	</div>
</div>
