<script lang="ts">
	// Props if needed, e.g., for response data, history handling, etc.
	// export let responseData: any = null;
	// export let responseHistory: any[] = [];

	// Event dispatcher for actions like toggling visibility or fetching history
	import { createEventDispatcher } from 'svelte';
	const dispatch = createEventDispatcher();

	let isExpanded = false; // Controls the visibility of the response body area

	function toggleExpand() {
		isExpanded = !isExpanded;
		dispatch('toggleExpand', { expanded: isExpanded });
	}

	function showHistory() {
		dispatch('showHistory');
	}
</script>

<footer class="bg-gray-900 dark:bg-gray-900 border-t border-gray-700 dark:border-gray-700">
	<div
		class="flex items-center justify-between p-2 border-b border-gray-700 dark:border-gray-700"
	>
		<div class="flex items-center space-x-2">
			<span class="text-sm font-medium text-gray-300 dark:text-gray-300">Response</span>
			<button
				on:click={showHistory}
				class="flex items-center space-x-1 text-sm text-gray-400 dark:text-gray-400 hover:text-gray-200 dark:hover:text-gray-200"
			>
				<i class="fas fa-history text-base"></i>
				<span>History</span>
			</button>
		</div>
		<button
			on:click={toggleExpand}
			class="text-gray-400 dark:text-gray-400 hover:text-gray-200 dark:hover:text-gray-200"
		>
			<i class="fas {isExpanded ? 'fa-chevron-down' : 'fa-chevron-up'}"></i>
		</button>
	</div>
	{#if isExpanded}
		<div class="flex flex-col items-center justify-center h-64 text-center p-4">
			<!-- Placeholder for response body or actual response display -->
			<div
				class="h-32 w-32 mb-4 opacity-75 bg-gray-700 dark:bg-gray-700 rounded-lg flex items-center justify-center"
			>
				<i class="fas fa-rocket text-4xl text-gray-500 dark:text-gray-500"></i>
			</div>
			<p class="text-gray-400 dark:text-gray-400">
				Enter the URL and click Send to get a response. Response content will appear here.
			</p>
			<!-- Example: Displaying actual response data -->
			<!-- {#if responseData} -->
			<!-- <pre class="text-left text-xs bg-gray-800 p-2 rounded-md overflow-auto max-h-56 w-full">{JSON.stringify(responseData, null, 2)}</pre> -->
			<!-- {:else} -->
			<!-- <p class="text-gray-400 dark:text-gray-400">No response data yet.</p> -->
			<!-- {/if} -->
		</div>
	{/if}
</footer>
