<script lang="ts">
	import type { JavaScriptConfig } from '$lib/types/Action';

	export let config: JavaScriptConfig;

	// Count lines in script
	$: lineCount = config.script ? config.script.split('\n').length : 0;
	$: characterCount = config.script ? config.script.length : 0;
</script>

<div
	class="mt-2 p-3 bg-gray-50 dark:bg-gray-800/50 rounded-md border border-gray-200 dark:border-gray-700"
>
	<!-- JavaScript Info -->
	<div class="flex items-center gap-3">
		<!-- Icon -->
		<div class="flex-shrink-0">
			<i class="fab fa-js text-yellow-500 dark:text-yellow-400 text-2xl"></i>
		</div>

		<!-- Script Stats -->
		<div class="flex-1">
			<div class="flex items-center gap-2 mb-1">
				<span class="text-xs font-medium theme-text-primary uppercase tracking-wide"
					>JavaScript Function</span
				>
			</div>
			<div class="flex items-center gap-3 text-xs theme-text-secondary">
				<div class="flex items-center gap-1">
					<i class="fas fa-code text-xs"></i>
					<span>{lineCount} {lineCount === 1 ? 'line' : 'lines'}</span>
				</div>
				<div class="flex items-center gap-1">
					<i class="fas fa-font text-xs"></i>
					<span>{characterCount} characters</span>
				</div>
			</div>
		</div>

		<!-- Available Context Badge -->
		<div class="flex-shrink-0">
			<div
				class="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-blue-100 dark:bg-blue-900/40 text-blue-700 dark:text-blue-300"
				title="Available context: request, response, helpers"
				aria-label="Context information"
			>
				<i class="fas fa-info-circle mr-1"></i>
				Context Available
			</div>
		</div>
	</div>

	<!-- Note about preview -->
	<div class="mt-2 pt-2 border-t border-gray-200 dark:border-gray-700">
		<p class="text-xs theme-text-secondary italic">
			<i class="fas fa-info-circle mr-1"></i>
			Code preview hidden for brevity. Edit action to view/modify script.
		</p>
	</div>
</div>
