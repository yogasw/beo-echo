<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	
	export let searchTerm: string = '';
</script>

<div class="theme-bg-secondary border theme-border p-6 rounded-lg flex flex-col items-center justify-center text-center">
	<div class="bg-blue-600/10 dark:bg-blue-600/20 p-4 rounded-full mb-4">
		<i class="fas fa-search text-blue-500 text-3xl"></i>
	</div>
	<h3 class="text-lg font-semibold theme-text-primary mb-2">No logs match your search</h3>
	<p class="theme-text-secondary mb-3 max-w-lg">
		We couldn't find any logs matching "{searchTerm}". Try using different keywords or clearing your search.
	</p>
	<button
		class={ThemeUtils.secondaryButton('py-1.5 px-4 text-sm')}
		on:click={() => {
			// Create a custom event that the parent can listen to
			const event = new CustomEvent('clearSearch');
			document.dispatchEvent(event);
		}}
	>
		<i class="fas fa-times mr-2"></i> Clear Search
	</button>
</div>
