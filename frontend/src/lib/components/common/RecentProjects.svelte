<script lang="ts">
	import { recentProjects, type RecentProject } from '$lib/stores/recentProjects';
	import { toast } from '$lib/stores/toast';

	export let showTitle = true;
	export let maxItems = 5;
	export let onProjectSelect: ((project: RecentProject) => void) | null = null;

	$: limitedProjects = $recentProjects.slice(0, maxItems);

	function formatTimeAgo(dateString: string): string {
		const now = new Date();
		const date = new Date(dateString);
		const diffInMs = now.getTime() - date.getTime();
		const diffInMinutes = Math.floor(diffInMs / (1000 * 60));
		const diffInHours = Math.floor(diffInMs / (1000 * 60 * 60));
		const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24));

		if (diffInMinutes < 1) {
			return 'Just now';
		} else if (diffInMinutes < 60) {
			return `${diffInMinutes}m ago`;
		} else if (diffInHours < 24) {
			return `${diffInHours}h ago`;
		} else if (diffInDays < 7) {
			return `${diffInDays}d ago`;
		} else {
			return date.toLocaleDateString();
		}
	}

	async function handleProjectClick(project: RecentProject) {
		try {
			if (onProjectSelect) {
				onProjectSelect(project);
			}
		} catch (error) {
			toast.error('Failed to open project');
			console.error('Error navigating to project:', error);
		}
	}

	function handleRemoveProject(event: Event, projectId: string) {
		event.stopPropagation();
		recentProjects.removeProject(projectId);
		toast.success('Project removed from recent list');
	}

	function getModeColor(mode?: string): string {
		switch (mode?.toLowerCase()) {
			case 'mock':
				return 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200';
			case 'proxy':
				return 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200';
			case 'forwarder':
				return 'bg-orange-100 dark:bg-orange-900 text-orange-800 dark:text-orange-200';
			default:
				return 'bg-gray-100 dark:bg-gray-800 text-gray-800 dark:text-gray-200';
		}
	}
</script>

{#if limitedProjects.length > 0}
	<div
		class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-hidden {showTitle
			? ''
			: 'border-0 bg-transparent'}"
	>
		{#if showTitle}
			<div
				class="px-6 py-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-750"
			>
				<h3 class="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
					<i class="fas fa-history text-indigo-600 dark:text-indigo-400 mr-2"></i>
					Recently Accessed Projects
				</h3>
				<p class="text-sm text-gray-600 dark:text-gray-400 mt-1">
					Quick access to your most recently used mock servers
				</p>
			</div>
		{/if}

		<div class="divide-y divide-gray-200 dark:divide-gray-700 {showTitle ? '' : 'divide-y-0'} {!showTitle && limitedProjects.length > 2 ? 'max-h-32 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-400 dark:scrollbar-thumb-gray-600 scrollbar-track-transparent' : ''}">
			{#each limitedProjects as project (project.id)}
				<div
					class="group flex items-center justify-between {showTitle
						? 'p-4'
						: 'p-3'} hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors cursor-pointer {showTitle
						? ''
						: 'rounded-lg mb-2 last:mb-0'}"
					on:click={() => handleProjectClick(project)}
					on:keydown={(e) => {
						if (e.key === 'Enter' || e.key === ' ') {
							e.preventDefault();
							handleProjectClick(project);
						}
					}}
					role="button"
					tabindex="0"
					title="Open project {project.name}"
					aria-label="Open project {project.name}"
				>
					<div class="flex items-center flex-1 min-w-0">
						<!-- Project Icon -->
						<div
							class="flex-shrink-0 {showTitle
								? 'w-10 h-10'
								: 'w-8 h-8'} bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center mr-3"
						>
							<i class="fas fa-cube text-white {showTitle ? 'text-sm' : 'text-xs'}"></i>
						</div>

						<!-- Project Info -->
						<div class="flex-1 min-w-0">
							<div class="flex items-center gap-2 mb-1">
								<p
									class="{showTitle
										? 'text-sm'
										: 'text-xs'} font-medium text-gray-900 dark:text-white truncate"
								>
									{project.name}
								</p>
								{#if project.mode}
									<span
										class="inline-flex items-center px-2 py-0.5 rounded-full {showTitle
											? 'text-xs'
											: 'text-xs'} font-medium {getModeColor(project.mode)}"
									>
										{project.mode}
									</span>
								{/if}
							</div>

							<div
								class="flex items-center {showTitle
									? 'text-xs'
									: 'text-xs'} text-gray-500 dark:text-gray-400 space-x-2"
							>
								<span class="truncate">{project.url}</span>
								<span>•</span>
								<span class="flex-shrink-0">{formatTimeAgo(project.lastUsed)}</span>
							</div>

							{#if project.workspaceName && showTitle}
								<p class="text-xs text-gray-400 dark:text-gray-500 mt-1 flex items-center">
									<i class="fas fa-users mr-1"></i>
									{project.workspaceName}
								</p>
							{/if}
						</div>
					</div>

					<!-- Action Buttons -->
					<div
						class="flex items-center space-x-2 opacity-0 group-hover:opacity-100 transition-opacity"
					>
						<button
							on:click={(e) => handleRemoveProject(e, project.id)}
							class="p-1.5 text-gray-400 hover:text-red-500 dark:hover:text-red-400 rounded-full hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
							title="Remove from recent projects"
							aria-label="Remove {project.name} from recent projects"
						>
							<i class="fas fa-times text-xs"></i>
						</button>

						<button
							class="p-1.5 text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400 rounded-full hover:bg-gray-100 dark:hover:bg-gray-600 transition-colors"
							title="Open project"
							aria-label="Open project {project.name}"
						>
							<i class="fas fa-external-link-alt text-xs"></i>
						</button>
					</div>
				</div>
			{/each}
		</div>

		{#if $recentProjects.length > maxItems && showTitle}
			<div
				class="px-6 py-3 bg-gray-50 dark:bg-gray-750 border-t border-gray-200 dark:border-gray-700"
			>
				<button
					class="text-sm text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 dark:hover:text-indigo-300 font-medium"
					title="View all recent projects"
					aria-label="View all recent projects"
				>
					View all {$recentProjects.length} recent projects
				</button>
			</div>
		{/if}
	</div>
{:else}
	<div
		class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 {showTitle
			? 'p-8'
			: 'p-4'} text-center {showTitle ? '' : 'border-0 bg-transparent'}"
	>
		<div
			class="{showTitle
				? 'w-16 h-16'
				: 'w-12 h-12'} bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-4"
		>
			<i class="fas fa-history text-gray-400 {showTitle ? 'text-xl' : 'text-lg'}"></i>
		</div>
		<h3 class="{showTitle ? 'text-lg' : 'text-sm'} font-medium text-gray-900 dark:text-white mb-2">
			No Recent Projects
		</h3>
		<p class="text-gray-600 dark:text-gray-400 {showTitle ? 'text-sm' : 'text-xs'}">
			{showTitle
				? "Start working with mock servers and they'll appear here for quick access"
				: 'Recent projects will appear here'}
		</p>
	</div>
{/if}
