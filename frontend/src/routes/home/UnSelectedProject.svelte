<script lang="ts">
	import { goto } from '$app/navigation';
	import RecentProjects from '$lib/components/common/RecentProjects.svelte';
	import { projects } from '$lib/stores/configurations';
	import type { RecentProject } from '$lib/stores/recentProjects';
	import workspaceStore from '$lib/stores/workspace';

	// Handle project selection from recent projects
	function handleProjectSelect(project: RecentProject) {
        goto(`/home/workspace/${project.workspaceId}/projects/${project.id}`);
	}

	// Stats computation
	let stats = $derived({
		totalProjects: $projects.length,
		totalWorkspaces: $workspaceStore.workspaces.length,
		mockMode: $projects.filter((p) => p.mode === 'mock').length,
		proxyMode: $projects.filter((p) => p.mode === 'proxy').length,
		forwarderMode: $projects.filter((p) => p.mode === 'forwarder').length,
		disabled: $projects.filter((p) => p.mode === 'disabled').length
	});
</script>

<div class="h-full overflow-y-auto p-3 md:p-4 bg-gray-50 dark:bg-gray-900">
	<div class="max-w-6xl mx-auto space-y-3">
		<!-- Welcome Header - Modern Compact -->
		<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
			<h1 class="text-lg font-bold text-gray-900 dark:text-white mb-2">
				Welcome to Beo Echo
			</h1>

			<!-- Features - Compact Modern -->
			<div class="grid grid-cols-4 gap-3">
				<!-- AI Assistant -->
				<div class="text-center">
					<div class="text-2xl mb-1">✨</div>
					<div class="text-xs font-semibold text-gray-900 dark:text-white">AI Assistant</div>
					<div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Generate mocks</div>
				</div>

				<!-- Lightning Fast -->
				<div class="text-center">
					<div class="text-2xl mb-1">⚡</div>
					<div class="text-xs font-semibold text-gray-900 dark:text-white">Lightning Fast</div>
					<div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Deploy instantly</div>
				</div>

				<!-- Multiple Modes -->
				<div class="text-center">
					<div class="text-2xl mb-1">⚙️</div>
					<div class="text-xs font-semibold text-gray-900 dark:text-white">Multiple Modes</div>
					<div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Mock & Proxy</div>
				</div>

				<!-- Collaborate -->
				<div class="text-center">
					<div class="text-2xl mb-1">👥</div>
					<div class="text-xs font-semibold text-gray-900 dark:text-white">Collaborate</div>
					<div class="text-xs text-gray-500 dark:text-gray-400 mt-0.5">Team sharing</div>
				</div>
			</div>
		</div>

		<!-- Stats Overview - Compact -->
		{#if stats.totalProjects > 0}
			<div class="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 p-3">
				<div class="flex flex-col md:flex-row md:items-center md:justify-between gap-3">
					<div class="flex items-center gap-2">
						<i class="fas fa-chart-pie text-indigo-600 dark:text-indigo-400 text-sm" aria-hidden="true"></i>
						<h2 class="text-sm font-semibold text-gray-900 dark:text-white">Overview</h2>
					</div>

					<div class="flex flex-wrap items-center gap-4">
						<!-- Total Projects -->
						<div
							class="flex items-center gap-1.5 cursor-help"
							title="Total number of mock API projects"
						>
							<i class="fas fa-cube text-blue-500 text-xs" aria-hidden="true"></i>
							<span class="text-lg font-bold text-gray-900 dark:text-white">{stats.totalProjects}</span>
							<span class="text-xs text-gray-600 dark:text-gray-400">Projects</span>
						</div>

						<!-- Divider -->
						<div class="h-4 w-px bg-gray-200 dark:bg-gray-700"></div>

						<!-- Total Workspaces -->
						<div
							class="flex items-center gap-1.5 cursor-help"
							title="Number of workspaces you belong to"
						>
							<i class="fas fa-users text-purple-500 text-xs" aria-hidden="true"></i>
							<span class="text-lg font-bold text-gray-900 dark:text-white">{stats.totalWorkspaces}</span>
							<span class="text-xs text-gray-600 dark:text-gray-400">Workspaces</span>
						</div>

						<!-- Divider -->
						<div class="h-4 w-px bg-gray-200 dark:bg-gray-700"></div>

						<!-- Modes Summary -->
						<div class="flex items-center gap-3">
							<div
								class="flex items-center gap-1 cursor-help"
								title="Mock Mode: Serves predefined mock responses"
							>
								<i class="fas fa-database text-indigo-500 text-xs" aria-hidden="true"></i>
								<span class="text-sm font-semibold text-gray-900 dark:text-white">{stats.mockMode}</span>
								<span class="text-xs text-gray-500 dark:text-gray-400 hidden sm:inline">Mock</span>
							</div>
							<div
								class="flex items-center gap-1 cursor-help"
								title="Proxy Mode: Uses mocks when available, forwards otherwise"
							>
								<i class="fas fa-sync text-green-500 text-xs" aria-hidden="true"></i>
								<span class="text-sm font-semibold text-gray-900 dark:text-white">{stats.proxyMode}</span>
								<span class="text-xs text-gray-500 dark:text-gray-400 hidden sm:inline">Proxy</span>
							</div>
							<div
								class="flex items-center gap-1 cursor-help"
								title="Forwarder Mode: Always forwards all requests to target"
							>
								<i class="fas fa-arrow-right text-orange-500 text-xs" aria-hidden="true"></i>
								<span class="text-sm font-semibold text-gray-900 dark:text-white">{stats.forwarderMode}</span>
								<span class="text-xs text-gray-500 dark:text-gray-400 hidden sm:inline">Forward</span>
							</div>
							<div
								class="flex items-center gap-1 cursor-help"
								title="Disabled: Projects that are currently inactive"
							>
								<i class="fas fa-ban text-gray-400 text-xs" aria-hidden="true"></i>
								<span class="text-sm font-semibold text-gray-900 dark:text-white">{stats.disabled}</span>
								<span class="text-xs text-gray-500 dark:text-gray-400 hidden sm:inline">Disabled</span>
							</div>
						</div>
					</div>
				</div>
			</div>
		{/if}

        <!-- Quick Start - Only for New Users -->
		{#if $projects.length === 0}
			<div
				class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-lg p-4 border border-blue-200 dark:border-blue-800"
			>
				<div class="flex items-center gap-3">
					<div class="w-10 h-10 bg-blue-600 rounded-lg flex items-center justify-center flex-shrink-0">
						<i class="fas fa-rocket text-white" aria-hidden="true"></i>
					</div>
					<div class="flex-1">
						<h3 class="text-sm font-semibold text-gray-900 dark:text-white mb-1">
							Ready to get started?
						</h3>
						<p class="text-xs text-gray-600 dark:text-gray-400">
							Create your first project from the sidebar to get started
						</p>
					</div>
					<button
						class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-2 rounded-lg text-xs font-medium transition-colors flex items-center gap-1.5"
						onclick={() => {
							const addButton = document.querySelector(
								'button[aria-label="Add new project"]'
							) as HTMLButtonElement;
							addButton?.click();
						}}
						title="Create your first project"
						aria-label="Create your first project"
					>
						<i class="fas fa-plus text-xs" aria-hidden="true"></i>
						Create Project
					</button>
				</div>
			</div>
		{/if}
        
		<!-- Recent Projects Section - Main Focus -->
		<div class="flex-1">
			<div class="mb-2 px-1">
				<h2 class="text-sm font-semibold text-gray-700 dark:text-gray-300 flex items-center gap-2">
					<i class="fas fa-history text-indigo-500" aria-hidden="true"></i>
					Recent Projects
				</h2>
			</div>
			<RecentProjects onProjectSelect={handleProjectSelect} />
		</div>
	</div>
</div>
