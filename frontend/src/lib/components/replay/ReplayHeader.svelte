<script lang="ts">
	import { onMount, onDestroy, tick } from 'svelte'; // Import tick
	import * as ThemeUtils from '$lib/utils/themeUtils';
	export var activeTabId: string;
	export var switchTab: (tabId: string) => void;
	export var closeTab: (tabId: string) => void;
	export var createNewTab: () => void;
	export var tabs: {
		id: string;
		name: string;
		method: string;
		isUnsaved?: boolean;
	}[] = [];

	let scrollableContainer: HTMLDivElement | null = null;
	let isOverflowing = false;
	let mutationObserver: MutationObserver | null = null;

	function checkOverflow() {
		if (scrollableContainer) {
			const currentlyOverflowing = scrollableContainer.scrollWidth > scrollableContainer.clientWidth;
			if (isOverflowing !== currentlyOverflowing) {
				isOverflowing = currentlyOverflowing;
			}
		}
	}

	onMount(() => {
		// Initial check after DOM is ready for the scrollableContainer
		// The reactive block below handles the first check once scrollableContainer is bound.

		window.addEventListener('resize', checkOverflow);

		if (scrollableContainer) {
			mutationObserver = new MutationObserver(checkOverflow);
			// Observe changes that could affect overflow: child additions/removals, attribute changes (like style/class)
			mutationObserver.observe(scrollableContainer, {
				childList: true,
				attributes: true,
				subtree: true // If children's content changes width
			});
		}
	});

	onDestroy(() => {
		window.removeEventListener('resize', checkOverflow);
		if (mutationObserver) {
			mutationObserver.disconnect();
		}
	});

	// Reactive statement to check overflow when tabs change or scrollableContainer is bound/updated
	$: if (scrollableContainer && typeof tabs !== 'undefined') {
		// tabs is included to ensure this runs when the number of tabs changes
		const updateOverflowState = async () => {
			await tick(); // Wait for Svelte to update the DOM based on any changes (e.g., tabs)
			checkOverflow();
		};
		updateOverflowState();
	}
</script>

<!-- Header with tabs/actions -->
<div class={ThemeUtils.themeBgSecondary('border-b theme-border')}>
	<div class="flex items-center justify-between px-4 py-2 text-sm">
		<div class="flex items-center space-x-2 flex-1 min-w-0">
			<button
				class={ThemeUtils.themeBgAccent(
					'flex items-center space-x-1 px-2 py-1 rounded-md theme-text-primary flex-shrink-0'
				)}
				title="Replay mode"
				aria-label="Replay mode"
			>
				<i class="fas fa-play text-sm"></i>
				<span>Replay</span>
			</button>

			<div
				class="flex items-center space-x-1 overflow-x-auto flex-1 hide-scrollbar"
				bind:this={scrollableContainer}
			>
				{#each tabs as tab (tab.id)}
					<div
						class="flex items-center bg-gray-100 dark:bg-gray-750 rounded-lg transition-all duration-200 hover:shadow-md flex-shrink-0"
					>
						<button
							class="flex items-center space-x-2 px-3 py-2 {activeTabId === tab.id
								? 'bg-blue-600 text-white shadow-md'
								: 'hover:bg-gray-200 dark:hover:bg-gray-600 theme-text-primary'} rounded-l-lg transition-all duration-200 min-w-0"
							title="Switch to {tab.name} ({tab.method})"
							aria-label="Switch to tab {tab.name} using {tab.method} method"
							on:click={() => switchTab(tab.id)}
						>
							<span
								class={activeTabId === tab.id
									? 'px-2 py-0.5 rounded text-xs font-semibold bg-white/20 text-white'
									: ThemeUtils.methodBadge(tab.method, 'text-xs px-1.5 py-0.5')}
							>
								{tab.method}
							</span>
							<span class="max-w-24 truncate text-sm font-medium">{tab.name}</span>
							{#if tab.isUnsaved}
								<span
									class="w-2 h-2 bg-orange-500 rounded-full flex-shrink-0"
									title="Unsaved changes"
									aria-label="This tab has unsaved changes"
								></span>
							{/if}
						</button>
						<button
							class="p-2 hover:bg-red-500 hover:text-white rounded-r-lg transition-all duration-200 theme-text-muted hover:theme-text-white"
							title="Close {tab.name} tab"
							aria-label="Close tab {tab.name}"
							on:click={() => closeTab(tab.id)}
						>
							<i class="fas fa-times text-xs"></i>
						</button>
					</div>
				{/each}
				<!--- Button 1 (Inline add button - shown when not overflowing) -->
				{#if !isOverflowing}
					<button
						class="p-2 ml-2 hover:bg-blue-600 hover:text-white bg-gray-200 dark:bg-gray-700 theme-text-primary rounded-lg transition-all duration-200 flex items-center justify-center shadow-sm hover:shadow-md flex-shrink-0"
						title="Create new request tab"
						aria-label="Add new request tab"
						on:click={createNewTab}
					>
						<i class="fas fa-plus text-sm"></i>
					</button>
				{/if}
			</div>

			<!--- Button 2 (Fixed add button - shown when overflowing) -->
			{#if isOverflowing}
				<button
					class="p-2 ml-2 hover:bg-blue-600 hover:text-white bg-gray-200 dark:bg-gray-700 theme-text-primary rounded-lg transition-all duration-200 flex items-center justify-center shadow-sm hover:shadow-md flex-shrink-0"
					title="Create new request tab"
					aria-label="Add new request tab"
					on:click={createNewTab}
				>
					<i class="fas fa-plus text-sm"></i>
				</button>
			{/if}
		</div>
	</div>
</div>
