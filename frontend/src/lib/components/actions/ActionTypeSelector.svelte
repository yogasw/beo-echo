<script lang="ts">
	import { onMount } from 'svelte';
	import { actionsApi } from '$lib/api/actionsApi';
	import { toast } from '$lib/stores/toast';
	import type { ActionTypeInfo } from '$lib/types/Action';

	export let onSelectType: (typeId: string) => void;
	export let onCancel: () => void;

	let actionTypes: ActionTypeInfo[] = [];
	let isLoading = true;

	onMount(async () => {
		try {
			const response = await actionsApi.getActionTypes();
			actionTypes = response.data;
		} catch (err) {
			toast.error('Failed to load action types');
		} finally {
			isLoading = false;
		}
	});
</script>

<div>
	<!-- Header -->
	<div class="px-4 pt-4 pb-3 border-b theme-border">
		<div class="flex items-center gap-3">
			<!-- Back Button (Left) -->
			<button
				type="button"
				class="group flex items-center justify-center w-10 h-10 rounded-lg bg-gray-100 dark:bg-gray-700/50 hover:bg-gray-200 dark:hover:bg-gray-600 transition-all duration-200 flex-shrink-0"
				on:click={onCancel}
				title="Back to actions list"
				aria-label="Back to actions list"
			>
				<i class="fas fa-arrow-left text-base text-gray-600 dark:text-gray-300 group-hover:text-gray-800 dark:group-hover:text-gray-100 transition-colors"></i>
			</button>

			<!-- Header Content -->
			<div class="flex items-center">
				<div class="flex items-center justify-center w-10 h-10 bg-purple-600/10 dark:bg-purple-600/10 rounded-lg mr-3 flex-shrink-0">
					<i class="fas fa-bolt text-purple-500 text-lg"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">Select Action Type</h2>
					<p class="text-sm theme-text-muted">Choose an action to automate your requests</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Content -->
	<div class="flex-1 overflow-y-auto p-4 theme-bg-primary">
		<div class="max-w-6xl mx-auto">
			{#if isLoading}
				<div class="text-center py-12">
					<i class="fas fa-spinner fa-spin text-4xl theme-text-secondary"></i>
					<p class="mt-4 text-sm theme-text-secondary">Loading action types...</p>
				</div>
			{:else}
				<!-- Flat list without category grouping -->
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
					{#each actionTypes as actionType}
						<button
							type="button"
							class="group relative p-4 rounded-lg bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:border-purple-500 dark:hover:border-purple-500 hover:shadow-lg transition-all duration-200 text-left overflow-hidden"
							on:click={() => onSelectType(actionType.id)}
							title="Create {actionType.name} action"
							aria-label="Create {actionType.name} action"
						>
							<!-- Hover gradient overlay -->
							<div
								class="absolute inset-0 bg-gradient-to-br from-purple-500/0 to-blue-500/0 group-hover:from-purple-500/5 group-hover:to-blue-500/5 transition-all duration-200"
							></div>

							<!-- Content -->
							<div class="relative flex flex-col items-start">
								<!-- Icon with category color -->
								<div
									class="w-11 h-11 rounded-lg flex items-center justify-center mb-3 group-hover:scale-105 transition-transform duration-200 {actionType.category === 'Transform'
										? 'bg-blue-500/10 dark:bg-blue-500/20'
										: actionType.category === 'Network'
											? 'bg-green-500/10 dark:bg-green-500/20'
											: actionType.category === 'Timing'
												? 'bg-amber-500/10 dark:bg-amber-500/20'
												: 'bg-purple-500/10 dark:bg-purple-500/20'}"
								>
									<i
										class="fas {actionType.icon} text-lg {actionType.category === 'Transform'
											? 'text-blue-500'
											: actionType.category === 'Network'
												? 'text-green-500'
												: actionType.category === 'Timing'
													? 'text-amber-500'
													: 'text-purple-500'}"
									></i>
								</div>

								<!-- Name -->
								<h4 class="text-base font-bold theme-text-primary mb-1.5 group-hover:text-purple-600 dark:group-hover:text-purple-400 transition-colors">
									{actionType.name}
								</h4>

								<!-- Description -->
								<p class="text-xs theme-text-secondary mb-3 leading-relaxed line-clamp-2">
									{actionType.description}
								</p>

								<!-- Category Badge -->
								<div class="flex items-center gap-2 mt-auto">
									<span
										class="inline-flex items-center px-2 py-0.5 text-xs font-medium rounded {actionType.category === 'Transform'
											? 'bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300'
											: actionType.category === 'Network'
												? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-300'
												: actionType.category === 'Timing'
													? 'bg-amber-100 dark:bg-amber-900/30 text-amber-700 dark:text-amber-300'
													: 'bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-300'}"
									>
										{actionType.category}
									</span>
								</div>

								<!-- Arrow indicator on hover -->
								<div
									class="absolute top-3 right-3 opacity-0 group-hover:opacity-100 transition-opacity duration-200"
								>
									<i class="fas fa-arrow-right text-sm text-purple-500"></i>
								</div>
							</div>
						</button>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>
