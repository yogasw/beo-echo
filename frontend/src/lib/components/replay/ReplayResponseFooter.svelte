<script>
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

<footer class="bg-white dark:bg-gray-800 border-t border-gray-300 dark:border-gray-600">
	<div
		class="flex items-center justify-between p-3 border-b border-gray-200 dark:border-gray-700"
	>
		<div class="flex items-center space-x-3">
			<span class="text-sm font-semibold text-gray-800 dark:text-white">Response</span>
			<button
				on:click={showHistory}
				class="flex items-center space-x-2 text-sm text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
			>
				<i class="fas fa-history text-sm"></i>
				<span>History</span>
			</button>
		</div>
		<button
			on:click={toggleExpand}
			class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
		>
			<i class="fas {isExpanded ? 'fa-chevron-down' : 'fa-chevron-up'}"></i>
		</button>
	</div>
	{#if isExpanded}
		<div class="flex flex-col items-center justify-center h-64 text-center p-6 bg-gray-50 dark:bg-gray-900">
			<!-- Placeholder for response body or actual response display -->
			<div
				class="h-24 w-24 mb-4 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center shadow-sm"
			>
				<i class="fas fa-rocket text-3xl text-gray-500 dark:text-gray-400"></i>
			</div>
			<p class="text-gray-600 dark:text-gray-300 max-w-md">
				Enter the URL and click Send to get a response. Response content will appear here.
			</p>
			<!-- Example: Displaying actual response data -->
			<!-- {#if responseData} -->
			<!-- <pre class="text-left text-xs bg-gray-100 dark:bg-gray-800 p-3 rounded-md overflow-auto max-h-56 w-full border border-gray-200 dark:border-gray-700">{JSON.stringify(responseData, null, 2)}</pre> -->
			<!-- {:else} -->
			<!-- <p class="text-gray-600 dark:text-gray-300">No response data yet.</p> -->
			<!-- {/if} -->
		</div>
	{/if}
</footer>
