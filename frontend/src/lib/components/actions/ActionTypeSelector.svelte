<script lang="ts">
	import { onMount } from 'svelte';
	import { actionsApi } from '$lib/api/actionsApi';
	import { toast } from '$lib/stores/toast';
	import type { ActionTypeInfo } from '$lib/types/Action';

	export let onSelectType: (typeId: string) => void;
	export let onCancel: () => void;

	let actionTypes: ActionTypeInfo[] = [];
	let isLoading = true;

	// Group action types by category
	$: groupedTypes = actionTypes.reduce(
		(acc, type) => {
			if (!acc[type.category]) {
				acc[type.category] = [];
			}
			acc[type.category].push(type);
			return acc;
		},
		{} as Record<string, ActionTypeInfo[]>
	);

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
		<div class="flex items-center gap-4">
			<!-- Back Button -->
			<button
				type="button"
				class="flex items-center justify-center w-10 h-10 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 transition-colors"
				on:click={onCancel}
				title="Back to actions list"
				aria-label="Back to actions list"
			>
				<i class="fas fa-arrow-left text-lg theme-text-primary"></i>
			</button>

			<!-- Header Content -->
			<div class="flex items-center flex-1">
				<div class="bg-purple-600/10 dark:bg-purple-600/10 p-2 rounded-lg mr-3">
					<i class="fas fa-bolt text-purple-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">Select Action Type</h2>
					<p class="text-sm theme-text-muted">Choose an action to automate your requests</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Content -->
	<div class="flex-1 overflow-y-auto p-6">
		<div class="max-w-6xl mx-auto">
			{#if isLoading}
				<div class="text-center py-12">
					<i class="fas fa-spinner fa-spin text-4xl theme-text-secondary"></i>
					<p class="mt-4 text-sm theme-text-secondary">Loading action types...</p>
				</div>
			{:else}
				{#each Object.entries(groupedTypes) as [category, types]}
					<div class="mb-8">
						<!-- Category Header -->
						<div class="flex items-center mb-4">
							<div class="flex-1 h-px theme-border"></div>
							<h3 class="px-4 text-sm font-semibold theme-text-primary uppercase tracking-wide">
								{category}
							</h3>
							<div class="flex-1 h-px theme-border"></div>
						</div>

						<!-- Action Type Cards -->
						<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
							{#each types as actionType}
								<button
									type="button"
									class="group p-6 rounded-lg border-2 theme-border theme-bg-secondary hover:border-blue-500 hover:bg-blue-900/10 transition-all text-left"
									on:click={() => onSelectType(actionType.id)}
									title="Create {actionType.name} action"
									aria-label="Create {actionType.name} action"
								>
									<div class="flex flex-col items-center text-center">
										<!-- Icon -->
										<div
											class="w-16 h-16 bg-blue-600/20 rounded-full flex items-center justify-center mb-4 group-hover:bg-blue-600/30 transition-colors"
										>
											<i class="fas {actionType.icon} text-3xl text-blue-500"></i>
										</div>

										<!-- Name -->
										<h4 class="text-lg font-semibold theme-text-primary mb-2">
											{actionType.name}
										</h4>

										<!-- Description -->
										<p class="text-sm theme-text-secondary mb-3">
											{actionType.description}
										</p>

										<!-- Category Badge -->
										<span
											class="inline-block px-3 py-1 text-xs rounded-full bg-gray-700 theme-text-secondary"
										>
											{actionType.category}
										</span>
									</div>
								</button>
							{/each}
						</div>
					</div>
				{/each}
			{/if}
		</div>
	</div>
</div>
