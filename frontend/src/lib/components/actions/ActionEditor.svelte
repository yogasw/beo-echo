<script lang="ts">
	import { actionsApi } from '$lib/api/actionsApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { toast } from '$lib/stores/toast';
	import type { Action, ActionFilter, ActionTypeInfo, CreateActionRequest, UpdateActionRequest } from '$lib/types/Action';
	import { ActionTypeComponents } from './types';
	import ActionFilterForm from './ActionFilterForm.svelte';
	import { onMount } from 'svelte';

	export let action: Action | null = null;
	export let onCancel: () => void;
	export let onSave: () => void;

	let isSubmitting = false;
	let actionTypes: ActionTypeInfo[] = [];
	let isLoadingTypes = true;

	// Form fields
	let name = action?.name || '';
	let selectedActionType: string = action?.type || '';
	let executionPoint: 'before_request' | 'after_request' = action?.execution_point || 'after_request';
	let enabled = action?.enabled ?? true;
	let priority = action?.priority || 0;

	// Config for action type
	let config: any = null;

	// Filters
	let filters: Omit<ActionFilter, 'id' | 'action_id' | 'created_at' | 'updated_at'>[] = [];

	// Expandable sections state
	let expandedSections = {
		basic: true,
		type: false,
		config: false,
		filters: false
	};

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

			// For edit mode, expand all sections
			expandedSections = {
				basic: true,
				type: true,
				config: true,
				filters: true
			};
		}
	}

	// Auto-expand next section when previous is filled
	$: if (name.trim()) {
		expandedSections.type = true;
	}

	$: if (selectedActionType) {
		expandedSections.config = true;
	}

	$: if (config) {
		expandedSections.filters = true;
	}

	// Load action types
	onMount(async () => {
		try {
			const response = await actionsApi.getActionTypes();
			actionTypes = response.data;
		} catch (err) {
			toast.error('Failed to load action types');
		} finally {
			isLoadingTypes = false;
		}
	});

	function toggleSection(section: keyof typeof expandedSections) {
		expandedSections[section] = !expandedSections[section];
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

		if (!selectedActionType) {
			toast.error('Please select an action type');
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
					type: selectedActionType as any,
					execution_point: executionPoint,
					enabled,
					priority,
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

	// Get selected action type info
	$: selectedTypeInfo = actionTypes.find(t => t.id === selectedActionType);

	// Get the component for the selected action type
	$: ActionConfigComponent = selectedActionType ? ActionTypeComponents[selectedActionType as keyof typeof ActionTypeComponents] : null;
</script>

<div class="h-full flex flex-col theme-bg-primary">
	<!-- Header -->
	<div class="p-6 theme-bg-primary border-b theme-border">
		<div class="flex items-center mb-4">
			<button
				on:click={onCancel}
				class="mr-4 p-2 rounded-lg theme-bg-secondary hover:bg-gray-600 transition-colors"
				title="Go back to actions list"
				aria-label="Go back to actions list"
			>
				<i class="fas fa-arrow-left theme-text-primary"></i>
			</button>
			<div class="flex items-center">
				<div class="bg-amber-600/10 dark:bg-amber-600/10 p-2 rounded-lg mr-3">
					<i class="fas fa-bolt text-amber-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">
						{action ? 'Edit Action' : 'Create New Action'}
					</h2>
					<p class="text-sm theme-text-muted">
						{action ? 'Modify your action configuration' : 'Configure a new action to transform requests or responses'}
					</p>
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
					<i class="fas fa-chevron-{expandedSections.basic ? 'up' : 'down'} theme-text-secondary"></i>
				</button>

				{#if expandedSections.basic}
					<div class="p-4 border-t theme-border space-y-4">
						<!-- Action Name -->
						<div>
							<label for="action-name" class="block text-sm font-medium theme-text-primary mb-2">
								Action Name <span class="text-red-500">*</span>
							</label>
							<input
								id="action-name"
								type="text"
								bind:value={name}
								class="block w-full p-3 text-sm rounded-lg theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
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
								class="block w-full p-3 text-sm rounded-lg theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
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
									class="block w-full p-3 text-sm rounded-lg theme-bg-primary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
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
				{/if}
			</div>

			<!-- Section 2: Select Action Type -->
			{#if !action}
				<div class="theme-bg-secondary rounded-lg border theme-border overflow-hidden">
					<button
						type="button"
						class="w-full p-4 flex items-center justify-between hover:bg-gray-700/30 transition-colors"
						on:click={() => toggleSection('type')}
						title="Toggle action type selection"
						aria-label="Toggle action type selection"
						aria-expanded={expandedSections.type}
					>
						<div class="flex items-center">
							<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center mr-3">
								<span class="text-white font-bold text-sm">2</span>
							</div>
							<div class="text-left">
								<h3 class="text-base font-semibold theme-text-primary">Select Action Type</h3>
								<p class="text-xs theme-text-muted">Choose what kind of transformation to apply</p>
							</div>
						</div>
						<i class="fas fa-chevron-{expandedSections.type ? 'up' : 'down'} theme-text-secondary"></i>
					</button>

					{#if expandedSections.type}
						<div class="p-4 border-t theme-border">
							{#if isLoadingTypes}
								<div class="text-center py-8">
									<i class="fas fa-spinner fa-spin text-3xl theme-text-secondary"></i>
									<p class="mt-2 text-sm theme-text-secondary">Loading action types...</p>
								</div>
							{:else}
								<div class="grid grid-cols-1 md:grid-cols-2 gap-3">
									{#each actionTypes as actionType}
										<button
											type="button"
											class="p-4 rounded-lg border-2 text-left transition-all {selectedActionType === actionType.id ? 'border-blue-500 bg-blue-900/20' : 'theme-border theme-bg-primary hover:border-blue-400'}"
											on:click={() => (selectedActionType = actionType.id)}
											title="Select {actionType.name}"
											aria-label="Select {actionType.name} action type"
										>
											<div class="flex items-start space-x-3">
												<div class="flex-shrink-0">
													<i class="fas {actionType.icon} text-2xl text-blue-500"></i>
												</div>
												<div class="flex-1">
													<h4 class="font-semibold theme-text-primary">{actionType.name}</h4>
													<p class="text-xs theme-text-secondary mt-1">{actionType.description}</p>
													<span class="inline-block mt-2 px-2 py-1 text-xs rounded bg-gray-700 theme-text-secondary">
														{actionType.category}
													</span>
												</div>
												{#if selectedActionType === actionType.id}
													<i class="fas fa-check-circle text-blue-500"></i>
												{/if}
											</div>
										</button>
									{/each}
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/if}

			<!-- Section 3: Configure Action -->
			{#if selectedActionType && ActionConfigComponent}
				<div class="theme-bg-secondary rounded-lg border theme-border overflow-hidden">
					<button
						type="button"
						class="w-full p-4 flex items-center justify-between hover:bg-gray-700/30 transition-colors"
						on:click={() => toggleSection('config')}
						title="Toggle action configuration"
						aria-label="Toggle action configuration"
						aria-expanded={expandedSections.config}
					>
						<div class="flex items-center">
							<div class="w-8 h-8 bg-blue-600 rounded-full flex items-center justify-center mr-3">
								<span class="text-white font-bold text-sm">{action ? '2' : '3'}</span>
							</div>
							<div class="text-left">
								<h3 class="text-base font-semibold theme-text-primary">
									Configure {selectedTypeInfo?.name || 'Action'}
								</h3>
								<p class="text-xs theme-text-muted">Set up the action behavior and parameters</p>
							</div>
						</div>
						<i class="fas fa-chevron-{expandedSections.config ? 'up' : 'down'} theme-text-secondary"></i>
					</button>

					{#if expandedSections.config}
						<div class="p-4 border-t theme-border">
							<svelte:component this={ActionConfigComponent} {config} on:change={handleConfigChange} />
						</div>
					{/if}
				</div>
			{/if}

			<!-- Section 4: Filters (Optional) -->
			{#if selectedActionType}
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
								<p class="text-xs theme-text-muted">Add conditions to control when this action runs</p>
							</div>
						</div>
						<i class="fas fa-chevron-{expandedSections.filters ? 'up' : 'down'} theme-text-secondary"></i>
					</button>

					{#if expandedSections.filters}
						<div class="p-4 border-t theme-border">
							<ActionFilterForm {filters} on:change={handleFiltersChange} />
						</div>
					{/if}
				</div>
			{/if}

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
					disabled={isSubmitting || !selectedActionType}
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
