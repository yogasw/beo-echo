<script lang="ts">
	import { highlightSegments } from '$lib/utils/logSearch';

	export let text: string = '';
	export let terms: string[] = [];

	$: segments = highlightSegments(text, terms);
</script>

{#each segments as segment}{#if segment.match}<mark
			class="bg-yellow-300 text-gray-900 dark:bg-yellow-400 rounded-sm px-0.5">{segment.text}</mark
		>{:else}{segment.text}{/if}{/each}
