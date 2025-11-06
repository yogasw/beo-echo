<script lang="ts">
	import ActionTypeSelector from './ActionTypeSelector.svelte';
	import ReplaceTextActionForm from './modules/ReplaceText/ReplaceTextActionForm.svelte';
	import type { Action } from '$lib/types/Action';

	export let action: Action | null = null;
	export let onCancel: () => void;
	export let onSave: () => void;

	// Wizard state
	let selectedType: string | null = action?.type || null;

	// Type-to-component mapping
	const ActionForms = {
		replace_text: ReplaceTextActionForm
		// Future action types can be added here:
		// add_header: AddHeaderActionForm,
		// webhook: WebhookActionForm,
		// delay: DelayActionForm,
	};

	function handleSelectType(typeId: string) {
		selectedType = typeId;
	}

	function handleBackToSelector() {
		if (action) {
			// If editing, go back to list
			onCancel();
		} else {
			// If creating, go back to type selector
			selectedType = null;
		}
	}

	$: ActionFormComponent = selectedType ? ActionForms[selectedType as keyof typeof ActionForms] : null;
</script>

{#if !selectedType}
	<!-- Step 1: Select Action Type -->
	<ActionTypeSelector onSelectType={handleSelectType} {onCancel} />
{:else if ActionFormComponent}
	<!-- Step 2: Configure Selected Action Type -->
	<svelte:component
		this={ActionFormComponent}
		{action}
		onCancel={handleBackToSelector}
		{onSave}
	/>
{:else}
	<!-- Fallback if action type not found -->
	<div class="h-full flex flex-col items-center justify-center theme-bg-primary">
		<i class="fas fa-exclamation-triangle text-6xl text-yellow-500 mb-4"></i>
		<h3 class="text-xl font-bold theme-text-primary mb-2">Action Type Not Supported</h3>
		<p class="text-sm theme-text-secondary mb-4">
			The selected action type is not yet implemented
		</p>
		<button
			class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md text-sm"
			on:click={handleBackToSelector}
			title="Go back"
			aria-label="Go back to action type selection"
		>
			<i class="fas fa-arrow-left mr-2"></i>
			Go Back
		</button>
	</div>
{/if}
