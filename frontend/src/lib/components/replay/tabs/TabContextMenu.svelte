<script lang="ts">
	let {
		show = false,
		x = 0,
		y = 0,
		tabId = '',
		onClose,
		onCreateNewTab,
		onDuplicateTab,
		onCloseTab,
		onCloseOtherTabs,
		onCloseAllTabs
	}: {
		show: boolean;
		x: number;
		y: number;
		tabId: string;
		onClose: () => void;
		onCreateNewTab: () => void;
		onDuplicateTab: (tabId: string) => void;
		onCloseTab: (tabId: string) => void;
		onCloseOtherTabs: (tabId: string) => void;
		onCloseAllTabs: () => void;
	} = $props();

	// Close context menu when clicking outside
	$effect(() => {
		if (show) {
			const handleClick = (e: MouseEvent) => {
				// We don't want to close immediately if this was the right-click that opened it
				// but Svelte handles this normally.
				onClose();
			};
			const handleContextMenu = (e: MouseEvent) => {
				// Will be overridden by the tab's specific contextmenu if clicking on a tab
				// But we delay the close slightly to allow the new tab's handler to run first
				setTimeout(() => {
					if (!e.defaultPrevented) onClose();
				}, 10);
			};
			window.addEventListener('click', handleClick);
			window.addEventListener('contextmenu', handleContextMenu);
			return () => {
				window.removeEventListener('click', handleClick);
				window.removeEventListener('contextmenu', handleContextMenu);
			};
		}
	});
</script>

{#if show}
	<!-- Context Menu Popup Only -->
	<!-- svelte-ignore a11y_click_events_have_key_events -->
	<!-- svelte-ignore a11y_no_static_element_interactions -->
	<div 
		class="fixed z-50 min-w-48 bg-white dark:bg-gray-800 rounded-lg shadow-xl border border-gray-200 dark:border-gray-700 py-1 overflow-hidden transition-all duration-200"
		style="left: {x}px; top: {y}px;"
	>
		<button
			class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-blue-600 flex items-center space-x-2"
			onclick={() => {
				onCreateNewTab();
				onClose();
			}}
		>
			<i class="fas fa-plus w-4 text-center"></i>
			<span>New Request</span>
		</button>
		<button
			class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-blue-600 flex items-center space-x-2"
			onclick={() => {
				onDuplicateTab(tabId);
				onClose();
			}}
		>
			<i class="fas fa-copy w-4 text-center"></i>
			<span>Duplicate Tab</span>
		</button>
		<div class="h-px bg-gray-200 dark:bg-gray-700 my-1"></div>
		<button
			class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-red-500 flex items-center space-x-2"
			onclick={() => {
				onCloseTab(tabId);
				onClose();
			}}
		>
			<i class="fas fa-times w-4 text-center"></i>
			<span>Close Tab</span>
		</button>
		<button
			class="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-red-500 flex items-center space-x-2"
			onclick={() => {
				onCloseOtherTabs(tabId);
				onClose();
			}}
		>
			<i class="fas fa-expand-arrows-alt w-4 text-center"></i>
			<span>Close Other Tabs</span>
		</button>
		<button
			class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 dark:hover:bg-red-900/30 flex items-center space-x-2"
			onclick={() => {
				onCloseAllTabs();
				onClose();
			}}
		>
			<i class="fas fa-ban w-4 text-center"></i>
			<span>Close All Tabs</span>
		</button>
	</div>
{/if}
