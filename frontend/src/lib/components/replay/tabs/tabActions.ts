import type { Tab } from '../types';

export interface TabActionsContext {
	getTabs: () => Tab[];
	setTabs: (tabs: Tab[]) => void;
	getActiveTabId: () => string;
	setActiveTabId: (id: string) => void;
	getActiveTabContent: () => any;
	dispatchTabsChange: () => void;
	triggerResetActiveTabContent: () => void;
}

export function createTabActions(ctx: TabActionsContext) {
	return {
		closeOtherTabs: (tabIdToKeep: string) => {
			const currentTabs = ctx.getTabs();
			const newTabs = currentTabs.filter((tab) => tab.id === tabIdToKeep);
			ctx.setTabs(newTabs);
			ctx.setActiveTabId(tabIdToKeep);
			ctx.dispatchTabsChange();
		},

		closeAllTabs: () => {
			ctx.setTabs([]);
			ctx.setActiveTabId('');
			ctx.triggerResetActiveTabContent();
			// Note: The dispatch happens inside or after reset, as requested by the original implementation
			// It was: dispatch('tabschange', { tabs: [], activeTabId: '', activeTabContent: null });
			// We can trigger it in triggerResetActiveTabContent or directly
		},

		duplicateTab: (tabId: string) => {
			const currentTabs = ctx.getTabs();
			const tabToDuplicate = currentTabs.find((t) => t.id === tabId);
			if (!tabToDuplicate) return;

			// Generate a unique ID to prevent overlap
			const randomSuffix = Math.random().toString(36).substring(2, 9);
			const newTabId = `tab-${Date.now()}-${randomSuffix}`;
			
		const newTab: Tab = {
				...tabToDuplicate,
				id: newTabId,
				isUnsaved: true,
				content: tabToDuplicate.content ? JSON.parse(JSON.stringify(tabToDuplicate.content)) : undefined
			};
			// Set duplicated name on replay if it exists
			if (newTab.replay) {
				newTab.replay = { ...newTab.replay, name: `${tabToDuplicate.replay?.name || 'Request'} Copy` };
			}

			ctx.setTabs([...currentTabs, newTab]);
			ctx.setActiveTabId(newTabId);
			ctx.dispatchTabsChange();
		}
	};
}
