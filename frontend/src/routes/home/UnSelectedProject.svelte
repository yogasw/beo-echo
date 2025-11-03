<script lang="ts">
	import { goto } from '$app/navigation';
	import RecentProjects from '$lib/components/common/RecentProjects.svelte';
	import { projects } from '$lib/stores/configurations';
	import type { RecentProject } from '$lib/stores/recentProjects';

	// Handle project selection from recent projects
	function handleProjectSelect(project: RecentProject) {
        goto(`/home/workspace/${project.workspaceId}/projects/${project.id}`);
	}
</script>

<div class="h-full overflow-y-auto p-4 md:p-6 bg-gray-50 dark:bg-gray-900">
	<div class="max-w-5xl mx-auto space-y-6">
		<!-- Welcome Header - Compact -->
		<div class="text-center py-4">
			<h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
				Welcome to Beo Echo
			</h1>
			<p class="text-sm text-gray-600 dark:text-gray-400">
				Powerful API mocking service for rapid development
			</p>
		</div>

		<!-- Feature Highlights - Compact Design -->
		<div class="grid grid-cols-2 md:grid-cols-4 gap-3">
			<!-- AI Assistant - Featured -->
			<div
				class="relative overflow-hidden bg-gradient-to-br from-purple-500 via-pink-500 to-rose-500 rounded-lg p-4 text-white shadow hover:shadow-md transition-all group"
			>
				<div class="absolute top-0 right-0 w-20 h-20 bg-white/10 rounded-full -mr-10 -mt-10"></div>
				<div class="relative z-10">
					<div class="flex items-center justify-between mb-2">
						<div class="w-8 h-8 bg-white/20 backdrop-blur-sm rounded-lg flex items-center justify-center">
							<i class="fas fa-magic text-sm" aria-hidden="true"></i>
						</div>
						<span class="px-1.5 py-0.5 bg-white/20 backdrop-blur-sm rounded text-xs">✨</span>
					</div>
					<h3 class="text-sm font-bold mb-1">AI Assistant</h3>
					<p class="text-xs text-white/90 leading-relaxed">
						Generate mock responses with AI
					</p>
				</div>
			</div>

			<!-- Fast Setup -->
			<div
				class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all"
			>
				<div class="w-8 h-8 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center mb-2">
					<i class="fas fa-bolt text-blue-600 dark:text-blue-400 text-sm" aria-hidden="true"></i>
				</div>
				<h3 class="text-sm font-bold text-gray-900 dark:text-white mb-1">Lightning Fast</h3>
				<p class="text-xs text-gray-600 dark:text-gray-400">
					Deploy mocks in seconds
				</p>
			</div>

			<!-- Flexible Modes -->
			<div
				class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all"
			>
				<div class="w-8 h-8 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center mb-2">
					<i class="fas fa-sliders-h text-green-600 dark:text-green-400 text-sm" aria-hidden="true"></i>
				</div>
				<h3 class="text-sm font-bold text-gray-900 dark:text-white mb-1">Multiple Modes</h3>
				<p class="text-xs text-gray-600 dark:text-gray-400">
					Mock, Proxy & Forward
				</p>
			</div>

			<!-- Team Collaboration -->
			<div
				class="bg-white dark:bg-gray-800 rounded-lg p-4 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all"
			>
				<div class="w-8 h-8 bg-purple-100 dark:bg-purple-900 rounded-lg flex items-center justify-center mb-2">
					<i class="fas fa-users text-purple-600 dark:text-purple-400 text-sm" aria-hidden="true"></i>
				</div>
				<h3 class="text-sm font-bold text-gray-900 dark:text-white mb-1">Collaborate</h3>
				<p class="text-xs text-gray-600 dark:text-gray-400">
					Share with your team
				</p>
			</div>
		</div>

		<!-- Divider -->
		<div class="relative">
			<div class="absolute inset-0 flex items-center">
				<div class="w-full border-t border-gray-200 dark:border-gray-700"></div>
			</div>
			<div class="relative flex justify-center">
				<span class="px-3 bg-gray-50 dark:bg-gray-900 text-sm text-gray-500 dark:text-gray-400">
					Your Projects
				</span>
			</div>
		</div>

		<!-- Recent Projects Section -->
		<div>
			<RecentProjects onProjectSelect={handleProjectSelect} />
		</div>

		<!-- Quick Start - Compact -->
		{#if $projects.length === 0}
			<div
				class="bg-gradient-to-r from-blue-50 to-indigo-50 dark:from-blue-900/20 dark:to-indigo-900/20 rounded-lg p-6 border border-blue-200 dark:border-blue-800"
			>
				<div class="flex items-start gap-4">
					<div class="flex-shrink-0">
						<div
							class="w-12 h-12 bg-blue-600 rounded-lg flex items-center justify-center"
						>
							<i class="fas fa-rocket text-white text-xl" aria-hidden="true"></i>
						</div>
					</div>
					<div class="flex-1">
						<h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">
							Ready to get started?
						</h3>
						<p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
							Create your first project from the sidebar, define your endpoints, and start testing
							your APIs instantly.
						</p>
						<button
							class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors flex items-center gap-2"
							onclick={() => {
								const addButton = document.querySelector(
									'button[aria-label="Add new project"]'
								) as HTMLButtonElement;
								addButton?.click();
							}}
							title="Create your first project"
							aria-label="Create your first project"
						>
							<i class="fas fa-plus" aria-hidden="true"></i>
							Create First Project
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>
