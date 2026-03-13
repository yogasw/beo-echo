<script lang="ts">
	import type { ExecuteReplayResponse } from '$lib/types/Replay';

	let { executionResult }: { executionResult: ExecuteReplayResponse } = $props();
</script>

{#if executionResult.error}
	<div
		class="p-4 bg-red-100 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md text-red-800 dark:text-red-300 mb-4"
	>
		<h3 class="font-semibold">Error</h3>
		<p>{executionResult.error}</p>
	</div>
{:else if executionResult.response_body}
	<div>
		{executionResult.response_body}
	</div>
{:else}
	<div
		class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300"
	>
		No response body
	</div>
{/if}
