<script lang="ts">
	import { actionsApi } from '$lib/api/actionsApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { toast } from '$lib/stores/toast';
	import type {
		Action,
		ActionFilter,
		CreateActionRequest,
		UpdateActionRequest
	} from '$lib/types/Action';
	import ActionFilterForm from '../../ActionFilterForm.svelte';
	import BasicInfo from '$lib/components/actions/BasicInfo.svelte';
	import ReplaceTextAction from './ReplaceTextAction.svelte';

	export let action: Action | null = null;
	export let onCancel: () => void;
	export let onSave: () => void;

	let isSubmitting = false;

	// Basic info fields
	let name = action?.name || 'Replace Text';
	let executionPoint: 'before_request' | 'after_request' =
		action?.execution_point || 'after_request';
	let enabled = action?.enabled ?? true;

	// Config for replace text
	let config: any = null;

	// Filters
	let filters: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[] = [];

	// Initialize from existing action
	$: if (action) {
		try {
			config = JSON.parse(action.config);
		} catch {
			config = null;
		}
		filters =
			action.filters?.map((f) => ({
				type: f.type,
				key: f.key,
				operator: f.operator,
				value: f.value
			})) || [];
	}

	function handleBasicInfoChange(event: CustomEvent) {
		const { name: n, executionPoint: ep, enabled: e } = event.detail;
		name = n;
		executionPoint = ep;
		enabled = e;
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
					type: 'replace_text',
					execution_point: executionPoint,
					enabled,
					config: configString,
					filters
				};

				await actionsApi.createAction($selectedWorkspace.id, $selectedProject.id, createData);
				toast.success('Action created successfully');
			}

			onSave();
		} catch (err: any) {
			toast.error(err);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div class="h-full flex flex-col theme-bg-primary">
	<!-- Header -->
	<div class="px-4 pt-4 pb-3 border-b theme-border">
		<div class="flex items-center gap-3">
			<!-- Back Button (Left) -->
			<button
				type="button"
				class="group flex items-center justify-center w-10 h-10 rounded-lg bg-gray-100 dark:bg-gray-700/50 hover:bg-gray-200 dark:hover:bg-gray-600 transition-all duration-200 flex-shrink-0"
				on:click={onCancel}
				title="Go back to action type selection"
				aria-label="Go back to action type selection"
			>
				<i
					class="fas fa-arrow-left text-base text-gray-600 dark:text-gray-300 group-hover:text-gray-800 dark:group-hover:text-gray-100 transition-colors"
				></i>
			</button>

			<!-- Header Content -->
			<div class="flex items-center">
				<div
					class="flex items-center justify-center w-10 h-10 bg-blue-600/10 dark:bg-purple-600/10 rounded-lg mr-3 flex-shrink-0"
				>
					<i class="fas fa-exchange-alt text-blue-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">Create Replace Text Action</h2>
					<p class="text-sm theme-text-muted">Find and replace text in requests or responses</p>
				</div>
			</div>
		</div>
	</div>

	<!-- Form Content -->
	<div class="flex-1 overflow-y-auto">
		<form on:submit|preventDefault={handleSubmit} class="h-full">
			<div class="px-6 py-4 space-y-6">
				<!-- Basic Information -->
				<div>
					<div class="flex items-center gap-2 mb-3">
						<i class="fas fa-info-circle text-blue-500 text-sm"></i>
						<h3 class="text-sm font-semibold theme-text-primary">Basic Information</h3>
					</div>
					<BasicInfo
						bind:name
						bind:executionPoint
						bind:enabled
						on:change={handleBasicInfoChange}
					/>
				</div>

				<!-- Divider -->
				<div class="border-t theme-border"></div>

				<!-- Replace Text Configuration -->
				<div>
					<div class="flex items-center gap-2 mb-3">
						<i class="fas fa-exchange-alt text-blue-500 text-sm"></i>
						<h3 class="text-sm font-semibold theme-text-primary">Replace Text Configuration</h3>
					</div>
					<ReplaceTextAction {config} on:change={handleConfigChange} />
				</div>

				<!-- Divider -->
				<div class="border-t theme-border"></div>

				<!-- Filters -->
				<div>
					<div class="flex items-center gap-2 mb-3">
						<i class="fas fa-filter text-purple-500 text-sm"></i>
						<h3 class="text-sm font-semibold theme-text-primary">Filters <span class="text-xs theme-text-muted font-normal">(Optional)</span></h3>
					</div>
					<ActionFilterForm {filters} on:change={handleFiltersChange} />
				</div>
			</div>

			<!-- Form Actions - Sticky Footer -->
			<div class="sticky bottom-0 px-6 py-3 border-t theme-border theme-bg-primary">
				<div class="flex justify-end space-x-3">
					<button
						type="button"
						class="bg-gray-700 hover:bg-gray-600 text-white py-2 px-5 rounded-md text-sm"
						on:click={onCancel}
						disabled={isSubmitting}
						title="Cancel and go back"
						aria-label="Cancel and go back"
					>
						Cancel
					</button>
					<button
						type="submit"
						class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-5 rounded-md text-sm flex items-center"
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
			</div>
		</form>
	</div>
</div>
