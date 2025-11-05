// Export all action type configuration components
// This makes it easy to add new action types in the future

import ReplaceTextAction from './ReplaceTextAction.svelte';

export const ActionTypeComponents = {
	replace_text: ReplaceTextAction,
	// Future action types can be added here:
	// add_header: AddHeaderAction,
	// webhook: WebhookAction,
	// delay: DelayAction,
};

export { ReplaceTextAction };
