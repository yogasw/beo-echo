<script lang="ts">
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

{#if parsedHeaders && parsedHeaders.length > 0}
	<ul class="text-xs space-y-1">
		{#each parsedHeaders as header}
			<li class="flex items-start break-all">
				<span class="font-bold whitespace-nowrap mr-1">{header.key}:</span>
				<span class="break-all">{header.value}</span>
			</li>
		{/each}
	</ul>
{:else}
	<div class="text-gray-400">No headers available.</div>
{/if}
