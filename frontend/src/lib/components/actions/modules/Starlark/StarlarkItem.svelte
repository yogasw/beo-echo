<script lang="ts">
	import type { Action } from '$lib/types/Action';

	export let action: Action;
</script>

<div class="flex items-center gap-3">
	<div class="flex items-center justify-center w-8 h-8 bg-yellow-600/10 dark:bg-yellow-600/20 rounded">
		<i class="fas fa-python text-yellow-500 text-sm"></i>
	</div>
	<div class="flex-1 min-w-0">
		<div class="font-medium theme-text-primary text-sm truncate">
			{action.name}
		</div>
		<div class="text-xs theme-text-muted">
			Starlark Script • {action.execution_point === 'before_request' ? 'Before' : 'After'} Request
		</div>
	</div>
</div>
