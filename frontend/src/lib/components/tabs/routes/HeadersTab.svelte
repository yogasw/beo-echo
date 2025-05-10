<script lang="ts">
	export let headers: { key: string; value: string }[] = [];
</script>

<h3 class="text-sm font-semibold mb-2">Headers</h3>
{#if headers && headers.length > 0}
	<ul class="text-xs space-y-1">
		{#each headers as header}
			<li class="flex items-start break-all">
				<span class="font-bold whitespace-nowrap mr-1">{header.key}:</span>
				<span class="break-all">{header.value}</span>
			</li>
		{/each}
	</ul>
{:else}
	<div class="text-gray-400">No headers available.</div>
{/if}
