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
	let name = action?.name || '';
	let executionPoint: 'before_request' | 'after_request' =
		action?.execution_point || 'after_request';
	let enabled = action?.enabled ?? true;

	// Config for replace text
	let config: any = null;

	// Filters
	let filters: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[] = [];

	// Expandable sections
	let expandedSections = {
		basic: true,
		config: false,
		filters: false
	};

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

		expandedSections = {
			basic: true,
			config: true,
			filters: true
		};
	}

	// Auto-expand next section
	$: if (name.trim()) {
		expandedSections.config = true;
	}

	$: if (config) {
		expandedSections.filters = true;
	}

	function toggleSection(section: keyof typeof expandedSections) {
		expandedSections[section] = !expandedSections[section];
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
	<div class="flex-1 overflow-y-auto p-6">
		<form on:submit|preventDefault={handleSubmit} class="max-w-4xl mx-auto space-y-4">
			<!-- Section 1: Basic Information -->
			<div class="theme-bg-secondary rounded-lg border theme-border overflow-hidden">
				<button
					type="button"
					class="w-full p-4 flex items-center justify-between hover:bg-gray-700/30 transition-colors"
					on:click={() => toggleSection('basic')}
					title="Toggle basic information section"
					aria-label="Toggle basic information section"
					aria-expanded={expandedSections.basic}
				>
					<div class="flex items-center">
						<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center mr-3">
							<span class="text-white font-bold text-sm">1</span>
						</div>
						<div class="text-left">
							<h3 class="text-base font-semibold theme-text-primary">Basic Information</h3>
							<p class="text-xs theme-text-muted">Action name and execution settings</p>
						</div>
					</div>
					<i class="fas fa-chevron-{expandedSections.basic ? 'up' : 'down'} theme-text-secondary"
					></i>
				</button>

				{#if expandedSections.basic}
					<div class="p-4 border-t theme-border">
						<BasicInfo
							bind:name
							bind:executionPoint
							bind:enabled
							on:change={handleBasicInfoChange}
						/>
					</div>
				{/if}
			</div>

			<!-- Section 2: Replace Text Configuration -->
			<div class="theme-bg-secondary rounded-lg border theme-border overflow-hidden">
				<button
					type="button"
					class="w-full p-4 flex items-center justify-between hover:bg-gray-700/30 transition-colors"
					on:click={() => toggleSection('config')}
					title="Toggle configuration section"
					aria-label="Toggle configuration section"
					aria-expanded={expandedSections.config}
				>
					<div class="flex items-center">
						<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center mr-3">
							<span class="text-white font-bold text-sm">2</span>
						</div>
						<div class="text-left">
							<h3 class="text-base font-semibold theme-text-primary">Replace Text Configuration</h3>
							<p class="text-xs theme-text-muted">Set up find and replace patterns</p>
						</div>
					</div>
					<i class="fas fa-chevron-{expandedSections.config ? 'up' : 'down'} theme-text-secondary"
					></i>
				</button>

				{#if expandedSections.config}
					<div class="p-4 border-t theme-border">
						<ReplaceTextAction {config} on:change={handleConfigChange} />
					</div>
				{/if}
			</div>

			<!-- Section 3: Filters (Optional) -->
			<div class="theme-bg-secondary rounded-lg border theme-border overflow-hidden">
				<button
					type="button"
					class="w-full p-4 flex items-center justify-between hover:bg-gray-700/30 transition-colors"
					on:click={() => toggleSection('filters')}
					title="Toggle filters section"
					aria-label="Toggle filters section"
					aria-expanded={expandedSections.filters}
				>
					<div class="flex items-center">
						<div class="w-8 h-8 bg-purple-600 rounded-full flex items-center justify-center mr-3">
							<i class="fas fa-filter text-white text-sm"></i>
						</div>
						<div class="text-left">
							<h3 class="text-base font-semibold theme-text-primary">Filters (Optional)</h3>
							<p class="text-xs theme-text-muted">
								Add conditions to control when this action runs
							</p>
						</div>
					</div>
					<i class="fas fa-chevron-{expandedSections.filters ? 'up' : 'down'} theme-text-secondary"
					></i>
				</button>

				{#if expandedSections.filters}
					<div class="p-4 border-t theme-border">
						<ActionFilterForm {filters} on:change={handleFiltersChange} />
					</div>
				{/if}
			</div>

			<!-- Form Actions -->
			<div class="flex justify-end space-x-3 pt-4">
				<button
					type="button"
					class="bg-gray-700 hover:bg-gray-600 text-white py-3 px-6 rounded-md text-sm"
					on:click={onCancel}
					disabled={isSubmitting}
					title="Cancel and go back"
					aria-label="Cancel and go back"
				>
					Cancel
				</button>
				<button
					type="submit"
					class="bg-blue-600 hover:bg-blue-700 text-white py-3 px-6 rounded-md text-sm flex items-center"
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
