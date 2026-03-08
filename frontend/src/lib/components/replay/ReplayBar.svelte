<script lang="ts">
	// Ensure Svelte 5 is correctly installed and configured in your project
	// for these imports to work.
	import {  tick } from 'svelte'; // Corrected import
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import TabContextMenu from './tabs/TabContextMenu.svelte';
	import type { Tab } from './types';

	// Props using Svelte 5 runes mode
	let {
		activeTabId,
		switchTab,
		closeTab,
		createNewTab,
		closeOtherTabs,
		closeAllTabs,
		duplicateTab,
		tabs = []
	}: {
		activeTabId: string;
		switchTab: (tabId: string) => void;
		closeTab: (tabId: string) => void;
		createNewTab: () => void;
		closeOtherTabs: (tabId: string) => void;
		closeAllTabs: () => void;
		duplicateTab: (tabId: string) => void;
		tabs?: Tab[];
	} = $props();

	// Reactive state for the DOM element reference
	let scrollableContainer: HTMLDivElement | null = $state(null);

	// Context menu state
	let contextMenu = $state({ show: false, x: 0, y: 0, tabId: '' });

	// Reactive state for overflow status
	let isOverflowing = $state(false);

	// Helper function to check and update overflow state
	function checkOverflow() {
		const sc = scrollableContainer; // Read $state variable
		if (sc) {
			const currentlyOverflowing = sc.scrollWidth > sc.clientWidth;
			if (isOverflowing !== currentlyOverflowing) {
				isOverflowing = currentlyOverflowing; // Update $state variable
			}
		}
	}

	async function handleCreateAndSwitchToNewTab() {
		const oldTabsCount = tabs.length; // `tabs` is a prop
		createNewTab(); // This prop function is expected to update the `tabs` array

		await tick(); // Wait for Svelte to process DOM changes

		// `tabs` prop will be updated by the parent, Svelte 5's reactivity handles the rest.
		if (tabs.length > oldTabsCount && tabs.length > 0) {
			const newTab = tabs[tabs.length - 1];
			if (newTab && newTab.id !== activeTabId) { // `activeTabId` is a prop
				switchTab(newTab.id);
			}
		}
	}

	function handleContextMenu(e: MouseEvent, tabId: string) {
		e.preventDefault();
		contextMenu = {
			show: true,
			x: e.clientX,
			y: e.clientY,
			tabId
		};
	}

	function closeContextMenu() {
		contextMenu.show = false;
	}

	// Close context menu when clicking outside
	$effect(() => {
		const handleClick = () => closeContextMenu();
		window.addEventListener('click', handleClick);
		return () => {
			window.removeEventListener('click', handleClick);
		};
	});

	// Effect for MutationObserver and initial check
	$effect(() => {
		const sc = scrollableContainer; // Dependency: scrollableContainer
		if (sc) {
			const observer = new MutationObserver(checkOverflow);
			observer.observe(sc, {
				childList: true,
				attributes: true, // Consider if needed, can be performance intensive
				subtree: true // Consider if needed
			});
			checkOverflow(); // Initial check

			// Cleanup function for this effect
			return () => {
				observer.disconnect();
			};
		}
	});

	// Effect for window resize listener
	$effect(() => {
		window.addEventListener('resize', checkOverflow);
		checkOverflow(); // Also check on mount

		// Cleanup function for this effect
		return () => {
			window.removeEventListener('resize', checkOverflow);
		};
	});

	// Effect to react to changes in `tabs` prop and `scrollableContainer`
	$effect(() => {
		const sc = scrollableContainer; // Dependency: scrollableContainer
		// Reading `tabs` here makes it a dependency as well.
		if (sc && tabs) {
			const updateOverflowStateAsync = async () => {
				await tick(); // Ensure DOM is updated after `tabs` change
				checkOverflow();
			};
			updateOverflowStateAsync();
		}
	});
</script>

<!-- Header with tabs/actions -->
<div class={ThemeUtils.themeBgSecondary('border-b theme-border')}>
	<div class="flex items-center justify-between px-4 py-2 text-sm relative">
		<div class="flex items-center space-x-2 flex-1 min-w-0">
			<button
				class={`${ThemeUtils.themeBgAccent(
					'flex items-center space-x-1 px-2 py-1 rounded-md theme-text-primary flex-shrink-0'
				)}`}
				title="Replay mode"
				aria-label="Replay mode"
			>
				<i class="fas fa-play text-sm"></i>
				<span>Replay</span>
			</button>

			<div
				class="flex items-center space-x-1 overflow-x-auto flex-1 hide-scrollbar relative"
				bind:this={scrollableContainer}
			>
				{#each tabs as tab (tab.id)}
					<!-- svelte-ignore a11y_no_static_element_interactions -->
					<div
						class="flex items-center bg-gray-100 dark:bg-gray-750 rounded-lg transition-all duration-200 hover:shadow-md flex-shrink-0"
						oncontextmenu={(e) => handleContextMenu(e, tab.id)}
					>
						<button
							class={`flex items-center space-x-2 px-3 py-2 ${activeTabId === tab.id
								? 'bg-blue-600 text-white shadow-md'
								: 'hover:bg-gray-200 dark:hover:bg-gray-600 theme-text-primary'} rounded-l-lg transition-all duration-200 min-w-0`}
							title="Switch to {tab.name} ({tab.method})"
							aria-label="Switch to tab {tab.name} using {tab.method} method"
							onclick={() => switchTab(tab.id)}
						>
							<span
								class={`${activeTabId === tab.id
									? 'px-2 py-0.5 rounded text-xs font-semibold bg-white/20 text-white'
									: tab.itemType === 'folder' 
										? 'text-xs text-orange-500' // Folder icon style
										: ThemeUtils.methodBadge(tab.method, 'text-xs px-1.5 py-0.5')}`}
							>
								{#if tab.itemType === 'folder'}
									<i class="fas fa-folder"></i>
								{:else}
									{tab.method}
								{/if}
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
							onclick={(e) => {
								e.stopPropagation();
								closeTab(tab.id);
							}}
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
						onclick={handleCreateAndSwitchToNewTab}
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
					onclick={handleCreateAndSwitchToNewTab}
				>
					<i class="fas fa-plus text-sm"></i>
				</button>
			{/if}

			<!-- Context Menu -->
			<TabContextMenu
				show={contextMenu.show}
				x={contextMenu.x}
				y={contextMenu.y}
				tabId={contextMenu.tabId}
				onClose={closeContextMenu}
				onCreateNewTab={createNewTab}
				onDuplicateTab={duplicateTab}
				onCloseTab={closeTab}
				onCloseOtherTabs={closeOtherTabs}
				onCloseAllTabs={closeAllTabs}
			/>
		</div>
	</div>
</div>
