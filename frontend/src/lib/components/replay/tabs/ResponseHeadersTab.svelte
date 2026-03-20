<script lang="ts">
	import type { ExecuteReplayResponse } from '$lib/types/Replay';

	let { executionResult }: { executionResult: ExecuteReplayResponse } = $props();
</script>

{#if executionResult.response_headers && Object.keys(executionResult.response_headers).length > 0}
	<div class="border border-gray-200 dark:border-gray-700 rounded-md overflow-hidden">
		<table class="w-full">
			<thead class="bg-gray-100 dark:bg-gray-700">
				<tr>
					<th
						class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
					>Name</th>
					<th
						class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
					>Value</th>
				</tr>
			</thead>
			<tbody>
				{#each Object.entries(executionResult.response_headers || {}) as [name, value], i}
					<tr
						class={i % 2 === 0
							? 'bg-white dark:bg-gray-800'
							: 'bg-gray-50 dark:bg-gray-900'}
					>
						<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm font-mono"
							>{name}</td>
						<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm font-mono"
							>{value}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{:else}
	<div
		class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300"
	>
		No headers received
	</div>
{/if}
