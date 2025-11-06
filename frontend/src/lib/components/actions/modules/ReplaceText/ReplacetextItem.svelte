<script lang="ts">
	import type { ReplaceTextConfig } from '$lib/types/Action';

	export let config: ReplaceTextConfig;

	// Get target label for replace_text
	function getTargetLabel(target: string): string {
		const labels: Record<string, string> = {
			request_body: 'Request Body',
			response_body: 'Response Body',
			request_header: 'Request Header',
			response_header: 'Response Header'
		};
		return labels[target] || target;
	}
</script>

<div class="mt-3 p-3 bg-gray-900/50 rounded border theme-border">
	<div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
		<div>
			<span class="theme-text-secondary">Target:</span>
			<span class="theme-text-primary font-medium ml-2">{getTargetLabel(config.target)}</span>
			{#if config.header_key}
				<span class="theme-text-secondary ml-2">({config.header_key})</span>
			{/if}
		</div>
		<div>
			<span class="theme-text-secondary">Mode:</span>
			<span class="theme-text-primary font-medium ml-2">
				{config.use_regex ? 'Regex' : 'Text'}
			</span>
		</div>
		<div class="md:col-span-2">
			<span class="theme-text-secondary">Pattern:</span>
			<code class="ml-2 px-2 py-1 bg-gray-800 rounded text-xs theme-text-primary font-mono">
				{config.pattern}
			</code>
		</div>
		<div class="md:col-span-2">
			<span class="theme-text-secondary">Replacement:</span>
			<code class="ml-2 px-2 py-1 bg-gray-800 rounded text-xs theme-text-primary font-mono">
				{config.replacement}
			</code>
		</div>
	</div>
</div>
