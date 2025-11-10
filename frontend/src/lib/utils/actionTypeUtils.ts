import type { ActionType } from '$lib/types/Action';

/**
 * Action type information including icon, label, color, and icon class
 */
export interface ActionTypeInfo {
	label: string;
	icon: string;
	iconClass: string;
	color: string;
}

/**
 * Get action type information (icon, label, color, icon class)
 * Used across the application for consistent action type display
 *
 * @param type - The action type
 * @returns Action type information object
 */
export function getActionTypeInfo(type: ActionType): ActionTypeInfo {
	const typeInfoMap: Record<ActionType, ActionTypeInfo> = {
		replace_text: {
			label: 'Replace Text',
			icon: 'fa-exchange-alt',
			iconClass: 'fas',
			color: 'text-amber-500'
		},
		run_javascript: {
			label: 'Run JavaScript',
			icon: 'fa-js',
			iconClass: 'fab', // Brand icon
			color: 'text-yellow-500'
		}
	};

	return (
		typeInfoMap[type] || {
			label: type,
			icon: 'fa-question',
			iconClass: 'fas',
			color: 'text-gray-400'
		}
	);
}

/**
 * Get only the icon class and icon name for an action type
 * Convenience function for when you only need the icon
 *
 * @param type - The action type
 * @returns Icon class and name (e.g., "fas fa-exchange-alt")
 */
export function getActionIcon(type: ActionType): string {
	const info = getActionTypeInfo(type);
	return `${info.iconClass} ${info.icon}`;
}
