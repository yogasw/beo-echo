<script lang="ts">
	import type { ReplaceTextConfig } from '$lib/types/Action';

	export let config: ReplaceTextConfig;

	// Get target label and icon for replace_text
	function getTargetInfo(target: string): { label: string; icon: string; color: string } {
		const info: Record<string, { label: string; icon: string; color: string }> = {
			request_body: {
				label: 'Response Body',
				icon: 'fa-arrow-left',
				color: 'text-green-500 dark:text-green-400'
			},
			response_body: {
				label: 'Response Body',
				icon: 'fa-arrow-left',
				color: 'text-green-500 dark:text-green-400'
			},
			request_header: {
				label: 'Request Header',
				icon: 'fa-list',
				color: 'text-purple-500 dark:text-purple-400'
			},
			response_header: {
				label: 'Response Header',
				icon: 'fa-list-alt',
				color: 'text-amber-500 dark:text-amber-400'
			}
		};
		return info[target] || { label: target, icon: 'fa-question', color: 'text-gray-400' };
	}

	$: targetInfo = getTargetInfo(config.target);
</script>

<div class="mt-2 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-md border border-gray-200 dark:border-gray-700">
	<!-- Compact Header: Target & Mode in one line -->
	<div class="flex items-center justify-between gap-2 mb-2.5 pb-2 border-b border-gray-200 dark:border-gray-700">
		<!-- Target -->
		<div class="flex items-center gap-2 min-w-0 flex-1">
			<i class="fas {targetInfo.icon} {targetInfo.color} text-sm flex-shrink-0"></i>
			<div class="min-w-0 flex-1">
				<span class="text-xs font-medium theme-text-primary uppercase tracking-wide">{targetInfo.label}</span>
				{#if config.header_key}
					<span class="text-xs theme-text-secondary ml-1">({config.header_key})</span>
				{/if}
			</div>
		</div>

		<!-- Compact Mode Badge -->
		<span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium flex-shrink-0 {config.use_regex
			? 'bg-purple-100 dark:bg-purple-900/40 text-purple-700 dark:text-purple-300'
			: 'bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300'}">
			<i class="fas {config.use_regex ? 'fa-code' : 'fa-text-width'} mr-1"></i>
			{config.use_regex ? 'Regex' : 'Text'}
		</span>
	</div>

	<!-- Horizontal Layout: Find → Replace -->
	<div class="grid grid-cols-2 gap-3">
		<!-- Find Column -->
		<div>
			<div class="flex items-center gap-1.5 mb-1">
				<i class="fas fa-search text-amber-500 dark:text-amber-400" style="font-size: 10px;"></i>
				<span class="text-xs theme-text-secondary font-medium uppercase" style="font-size: 10px; letter-spacing: 0.05em;">Find</span>
			</div>
			<code class="block px-2 py-1.5 bg-white dark:bg-gray-900/50 border border-gray-200 dark:border-gray-700 rounded text-xs theme-text-primary font-mono break-all">
				{config.pattern}
			</code>
		</div>

		<!-- Replace Column -->
		<div>
			<div class="flex items-center gap-1.5 mb-1">
				<i class="fas fa-exchange-alt text-green-500 dark:text-green-400" style="font-size: 10px;"></i>
				<span class="text-xs theme-text-secondary font-medium uppercase" style="font-size: 10px; letter-spacing: 0.05em;">Replace</span>
			</div>
			<code class="block px-2 py-1.5 bg-white dark:bg-gray-900/50 border border-gray-200 dark:border-gray-700 rounded text-xs theme-text-primary font-mono break-all">
				{config.replacement || '(empty)'}
			</code>
		</div>
	</div>
</div>
