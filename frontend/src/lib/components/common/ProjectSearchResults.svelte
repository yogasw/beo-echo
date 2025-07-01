<script lang="ts">
	import { goto } from '$app/navigation';
	import type { ProjectSearchResult } from '$lib/api/BeoApi';

	export let searchResults: ProjectSearchResult[] = [];
	export let showResults = false;
	export let onProjectSelect: ((project: ProjectSearchResult) => void) | undefined = undefined;

	// Handle project selection
	function selectProject(project: ProjectSearchResult) {
		// Call custom handler if provided, otherwise use default navigation
		if (onProjectSelect) {
			onProjectSelect(project);
		} else {
			// Default: navigate to project management page
			goto(`/home/workspace/${project.workspace_id}/projects/${project.id}`);
		}
	}
</script>

{#if showResults && searchResults.length > 0}
	<div class="absolute z-50 w-full bg-white dark:bg-gray-700 border border-gray-200 dark:border-gray-600 rounded-lg shadow-lg mt-1 max-h-60 overflow-y-auto">
		<div class="p-2 border-b border-gray-200 dark:border-gray-600">
			<span class="text-xs text-gray-600 dark:text-gray-400 font-medium">
				<i class="fas fa-search mr-1"></i>
				Found {searchResults.length} existing project{searchResults.length !== 1 ? 's' : ''}
			</span>
		</div>
		{#each searchResults as project}
			<button
				on:click={() => selectProject(project)}
				class="w-full text-left p-3 hover:bg-gray-50 dark:hover:bg-gray-600 transition-colors border-b border-gray-100 dark:border-gray-600 last:border-b-0"
				title="Open project {project.name}"
				aria-label="Open project {project.name}"
			>
				<div class="flex items-center justify-between">
					<div class="flex-1 min-w-0">
						<div class="font-medium text-gray-900 dark:text-white text-sm truncate">
							{project.name}
						</div>
						<div class="text-xs text-gray-500 dark:text-gray-400 truncate">
							{project.workspace_name}
						</div>
					</div>
					<div class="flex items-center ml-2">
						<i class="fas fa-arrow-right text-gray-400 text-xs"></i>
					</div>
				</div>
			</button>
		{/each}
	</div>
{/if}
