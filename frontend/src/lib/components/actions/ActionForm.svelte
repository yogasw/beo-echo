<script lang="ts">
	import { actionsApi } from '$lib/api/actionsApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import type { Action, ActionFilter, CreateActionRequest, UpdateActionRequest } from '$lib/types/Action';
	import ReplaceTextConfig from './ReplaceTextConfig.svelte';
	import ActionFilterForm from './ActionFilterForm.svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';

	export let action: Action | null = null;
	export let onClose: () => void;

	let isSubmitting = false;
	let name = action?.name || '';
	let executionPoint: 'before_request' | 'after_request' = action?.execution_point || 'after_request';
	let enabled = action?.enabled ?? true;
	let priority = action?.priority || 0;
	let actionType: 'replace_text' = 'replace_text';

	// Config for replace_text
	let config: any = null;

	// Filters
	let filters: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[] = [];

	// Initialize config and filters from action
	$: {
		if (action) {
			try {
				config = JSON.parse(action.config);
			} catch {
				config = null;
			}
			filters = action.filters?.map(f => ({
				type: f.type,
				key: f.key,
				operator: f.operator,
				value: f.value
			})) || [];
		}
	}

	function handleConfigChange(event: CustomEvent) {
		config = event.detail;
	}

	function handleFiltersChange(event: CustomEvent) {
		filters = event.detail;
	}

	async function handleSubmit() {
		if (!$selectedWorkspace || !$selectedProject) return;

		// Validation
		if (!name.trim()) {
			toast.error('Action name is required');
			return;
		}

		if (!config) {
			toast.error('Action configuration is required');
			return;
		}

		try {
			isSubmitting = true;

			const configString = JSON.stringify(config);

			if (action) {
				// Update existing action
				const updateData: UpdateActionRequest = {
					name: name.trim(),
					execution_point: executionPoint,
					enabled,
					priority,
					config: configString,
					filters
				};

				await actionsApi.updateAction(
					$selectedWorkspace.id,
					$selectedProject.id,
					action.id,
					updateData
				);
				toast.success('Action updated successfully');
			} else {
				// Create new action
				const createData: CreateActionRequest = {
					name: name.trim(),
					type: actionType,
					execution_point: executionPoint,
					enabled,
					priority,
					config: configString,
					filters
				};

				await actionsApi.createAction($selectedWorkspace.id, $selectedProject.id, createData);
				toast.success('Action created successfully');
			}

			onClose();
		} catch (err: any) {
			toast.error(err);
		} finally {
			isSubmitting = false;
		}
	}

	function handleCancel() {
		onClose();
	}
</script>

<!-- Modal Overlay -->
<div
	class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
	on:click={handleCancel}
	on:keydown={(e) => e.key === 'Escape' && handleCancel()}
	role="dialog"
	aria-modal="true"
	aria-labelledby="action-form-title"
	tabindex="-1"
>
	<!-- Modal Content -->
	<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
	<div
		class="theme-bg-primary rounded-lg shadow-xl w-full max-w-3xl max-h-[90vh] overflow-y-auto m-4"
		on:click|stopPropagation
		role="document"
	>
		<!-- Header -->
		<div class="flex justify-between items-center p-6 border-b theme-border">
			<h2 id="action-form-title" class="text-2xl font-bold theme-text-primary">
				{action ? 'Edit Action' : 'Create New Action'}
			</h2>
			<button
				class="text-gray-400 hover:text-white"
				on:click={handleCancel}
				title="Close modal"
				aria-label="Close modal"
			>
				<i class="fas fa-times text-xl"></i>
			</button>
		</div>

		<!-- Form -->
		<form on:submit|preventDefault={handleSubmit} class="p-6 space-y-6">
			<!-- Basic Settings -->
			<div class="space-y-4">
				<h3 class="text-lg font-semibold theme-text-primary flex items-center">
					<i class="fas fa-cog mr-2 text-blue-500"></i>
					Basic Settings
				</h3>

				<!-- Action Name -->
				<div>
					<label for="action-name" class="block text-sm font-medium theme-text-primary mb-2">
						Action Name <span class="text-red-500">*</span>
					</label>
					<input
						id="action-name"
						type="text"
						bind:value={name}
						class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
						placeholder="e.g., Replace API URL"
						required
						aria-required="true"
					/>
				</div>

				<!-- Execution Point -->
				<div>
					<label for="execution-point" class="block text-sm font-medium theme-text-primary mb-2">
						When to Execute <span class="text-red-500">*</span>
					</label>
					<select
						id="execution-point"
						bind:value={executionPoint}
						class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
						required
						aria-required="true"
					>
						<option value="after_request">After Request (Modify Response)</option>
						<option value="before_request">Before Request (Modify Request)</option>
					</select>
					<p class="mt-1 text-xs theme-text-secondary">
						{executionPoint === 'after_request'
							? 'Executes after receiving the response, allowing you to modify response data'
							: 'Executes before forwarding the request, allowing you to modify request data'}
					</p>
				</div>

				<!-- Priority and Enabled -->
				<div class="grid grid-cols-2 gap-4">
					<div>
						<label for="priority" class="block text-sm font-medium theme-text-primary mb-2">
							Priority
						</label>
						<input
							id="priority"
							type="number"
							bind:value={priority}
							min="0"
							class="block w-full p-3 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							placeholder="0"
							aria-describedby="priority-help"
						/>
						<p id="priority-help" class="mt-1 text-xs theme-text-secondary">
							Lower number = higher priority (0 is highest)
						</p>
					</div>

					<div>
						<label for="enabled" class="block text-sm font-medium theme-text-primary mb-2">
							Status
						</label>
						<label class="relative inline-flex items-center cursor-pointer mt-3">
							<input
								id="enabled"
								type="checkbox"
								bind:checked={enabled}
								class="sr-only peer"
								aria-label="Enable or disable this action"
							/>
							<div
								class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"
							></div>
							<span class="ml-3 text-sm font-medium theme-text-primary">
								{enabled ? 'Enabled' : 'Disabled'}
							</span>
						</label>
					</div>
				</div>
			</div>

			<!-- Replace Text Configuration -->
			<div class="space-y-4">
				<h3 class="text-lg font-semibold theme-text-primary flex items-center">
					<i class="fas fa-exchange-alt mr-2 text-amber-500"></i>
					Replace Text Configuration
				</h3>
				<ReplaceTextConfig {config} on:change={handleConfigChange} />
			</div>

			<!-- Filters (Optional) -->
			<div class="space-y-4">
				<h3 class="text-lg font-semibold theme-text-primary flex items-center">
					<i class="fas fa-filter mr-2 text-purple-500"></i>
					Filters (Optional)
				</h3>
				<p class="text-sm theme-text-secondary">
					Add filters to conditionally execute this action. If no filters are set, the action runs for all requests.
				</p>
				<ActionFilterForm {filters} on:change={handleFiltersChange} />
			</div>

			<!-- Form Actions -->
			<div class="flex justify-end space-x-3 pt-4 border-t theme-border">
				<button
					type="button"
					class="bg-gray-700 hover:bg-gray-600 text-white py-2 px-6 rounded-md text-sm"
					on:click={handleCancel}
					disabled={isSubmitting}
					title="Cancel and close"
					aria-label="Cancel and close"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-6 rounded-md text-sm flex items-center"
					disabled={isSubmitting}
					title={action ? 'Update action' : 'Create action'}
					aria-label={action ? 'Update action' : 'Create action'}
				>
					{#if isSubmitting}
						<i class="fas fa-spinner fa-spin mr-2"></i>
						{action ? 'Updating...' : 'Creating...'}
					{:else}
						<i class="fas fa-save mr-2"></i>
						{action ? 'Update Action' : 'Create Action'}
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
