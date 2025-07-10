<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast';
	import { tick } from 'svelte';
	import type { Rule } from '$lib/api/rulesApi';
	import { getRules, createRule, updateRule, deleteRule } from '$lib/api/rulesApi';

	// Props
	export let projectId: string;
	export let endpointId: string;
	export let responseId: string;
	export let rules: Rule[] = [];
	export let rulesOperator = 'OR'; // Default operator

	// Local state
	let isLoading = false;
	let isEditing = false;
	let hasChanges = false;
	let containerElement: HTMLElement;
	let tableContainer: HTMLElement;
	let calculatedHeight = '300px'; // Default height
	let keyInputs: HTMLInputElement[] = [];
	let valueInputs: HTMLInputElement[] = [];
	
	// Local copy of rules to avoid modifying the original during edit
	let localRules: Rule[] = [];
	
	// Track previous rules to detect changes
	let previousRulesJSON = JSON.stringify(rules);
	let previousResponseId = responseId;

	// Update local rules when original rules change
	$: {
		const rulesJSON = JSON.stringify(rules);
		const rulesChanged = rulesJSON !== previousRulesJSON;
		
		// Reset to view mode if new rules loaded while editing
		if (rulesChanged && isEditing) {
			isEditing = false;
			hasChanges = false;
		}
		
		// Update local rules if not editing or rules changed
		if (!isEditing || rulesChanged) {
			localRules = rules.map(rule => ({ ...rule }));
			hasChanges = false;
			previousRulesJSON = rulesJSON;
		}
	}

	// Watch for responseId changes and reload rules
	$: {
		if (responseId !== previousResponseId) {
			previousResponseId = responseId;
			// Reset editing state when response changes
			isEditing = false;
			hasChanges = false;
			// Reload rules for new response
			loadRules();
		}
	}

	// Load rules on component mount
	onMount(() => {
		// Start async operations without waiting
		loadRules();
		calculateHeight();
		
		// Setup resize observer to recalculate height
		const resizeObserver = new ResizeObserver(() => {
			calculateHeight();
		});
		
		if (containerElement) {
			resizeObserver.observe(containerElement);
		}
		
		// Listen to window resize events
		window.addEventListener('resize', calculateHeight);
		
		return () => {
			if (resizeObserver) {
				resizeObserver.disconnect();
			}
			window.removeEventListener('resize', calculateHeight);
		};
	});

	// Calculate height based on parent container
	function calculateHeight() {
		if (!containerElement || !tableContainer) return;
		
		// Get container position relative to viewport
		const containerRect = containerElement.getBoundingClientRect();
		// Calculate available height
		const availableHeight = window.innerHeight - containerRect.top - 75;
		// Set height within reasonable bounds
		calculatedHeight = `${Math.max(100, Math.min(350, availableHeight))}px`;
	}

	// Load rules from API
	async function loadRules() {
		isLoading = true;
		try {
			rules = await getRules(projectId, endpointId, responseId);
			localRules = rules.map(rule => ({ ...rule }));
			previousRulesJSON = JSON.stringify(rules);
			console.log('Rules loaded:', rules);
		} catch (error) {
			console.error('Error loading rules:', error);
			toast.error('Failed to load rules');
		} finally {
			isLoading = false;
		}
	}

	// Add a new rule inline (directly in the edit mode)
	async function addRule() {
		const newIndex = localRules.length;
		const newRule = {
			responseId: responseId,
			type: 'header',
			key: '',
			operator: 'equals',
			value: '',
			id: '',  // Temporary ID until saved
			isNew: true // Flag to identify new unsaved rules
		};
		
		localRules = [...localRules, newRule];
		hasChanges = true;
		
		// Focus on the new rule field after DOM update
		await tick();
		if (keyInputs[newIndex]) {
			keyInputs[newIndex].focus();
		}
	}

	// Handle field changes
	function handleTypeChange(index: number, newType: string) {
		if (localRules[index].type === newType) return; // Prevent unnecessary updates
		
		localRules[index].type = newType;
		// Clear the key field when type changes to "body" since body rules don't need a key
		if (newType === 'body') {
			localRules[index].key = '';
		}
		hasChanges = true;
		// Use more efficient update instead of spreading the entire array
		localRules = localRules.slice();
	}
	
	function handleKeyChange(index: number, newKey: string) {
		if (localRules[index].key === newKey) return; // Prevent unnecessary updates
		
		localRules[index].key = newKey;
		hasChanges = true;
		localRules = localRules.slice();
	}
	
	function handleOperatorChange(index: number, newOperator: string) {
		if (localRules[index].operator === newOperator) return; // Prevent unnecessary updates
		
		localRules[index].operator = newOperator;
		hasChanges = true;
		localRules = localRules.slice();
	}
	
	function handleValueChange(index: number, newValue: string) {
		if (localRules[index].value === newValue) return; // Prevent unnecessary updates
		
		localRules[index].value = newValue;
		hasChanges = true;
		localRules = localRules.slice();
	}
	
	// Remove a rule in edit mode
	function removeLocalRule(index: number) {
		localRules = localRules.filter((_, i) => i !== index);
		hasChanges = true;
	}

	// Enable edit mode
	function startEditing() {
		isEditing = true;
	}
	
	// Cancel editing and revert changes
	function cancelEdit() {
		// Reset to original rules
		localRules = rules.map(rule => ({ ...rule }));
		isEditing = false;
		hasChanges = false;
	}
	
	// Save all changes to rules
	async function saveChanges() {
		try {
			// Remove empty rules (but keep body rules which don't need keys)
			const filteredRules = localRules.filter(rule => {
				// Keep body rules even with empty keys
				if (rule.type === 'body') {
					return rule.value.trim() !== '';
				}
				// For header/query rules, both key and value must be present
				return rule.key.trim() !== '' && rule.value.trim() !== '';
			});
			
			// Handle updates for existing rules
			const updatePromises = filteredRules
				.filter(rule => rule.id && !rule?.isNew)
				.map(rule => 
					updateRule(
						projectId,
						endpointId,
						responseId,
						rule?.id,
						{
							type: rule.type,
							key: rule.key,
							operator: rule.operator,
							value: rule.value
						}
					)
				);
				
			// Handle creation of new rules
			const createPromises = filteredRules
				.filter(rule => !rule.id || rule.isNew)
				.map(rule => 
					createRule(
						projectId,
						endpointId,
						responseId,
						{
							responseId: rule.responseId,
							type: rule.type,
							key: rule.key,
							operator: rule.operator,
							value: rule.value
						}
					)
				);
			
			// Execute all promises
			const updatedRules = await Promise.all(updatePromises);
			const createdRules = await Promise.all(createPromises);
			
			// Update the rules array with the results
			rules = [...updatedRules, ...createdRules];
			toast.success('Rules saved successfully');
			
			isEditing = false;
			hasChanges = false;
		} catch (error) {
			console.error('Error saving rules:', error);
			toast.error('Failed to save rules');
		}
	}

	// Delete a rule in view mode
	async function removeRule(ruleId: string, index: number) {
		try {
			await deleteRule(projectId, endpointId, responseId, ruleId);
			// Remove the rule from the local array
			rules = rules.filter((_, i) => i !== index);
			localRules = localRules.filter((_, i) => i !== index);
			toast.success('Rule deleted successfully');
		} catch (error) {
			console.error('Error deleting rule:', error);
			toast.error('Failed to delete rule');
		}
	}

	// Toggle between AND/OR logic
	function toggleLogic(newOperator: string) {
		rulesOperator = newOperator;
	}

	// Handle key press events for quick add
	function handleKeyDown(event: KeyboardEvent, index: number) {
		if (event.key === 'Enter' && index === localRules.length - 1) {
			// Add a new rule when pressing Enter on the last rule
			event.preventDefault();
			addRule();
		}
	}
</script>


<div class="w-full rounded-md overflow-hidden shadow-sm" bind:this={containerElement}>
	<div class="p-0">
		<!-- Header with logic operators and edit/save controls -->
		<div class="bg-gray-100 dark:bg-gray-700 border-b dark:border-gray-600 flex justify-between items-center">
			<div class="flex items-center space-x-4 px-4 py-3">
				<span class="text-xs font-semibold {ThemeUtils.themeTextPrimary()}">Rules Logic:</span>
				<button on:click={() => toggleLogic('OR')}
								class="logic-button {rulesOperator === 'OR' ? 'text-blue-500 border-blue-500' : ThemeUtils.themeTextMuted() + ' border ' + ThemeUtils.themeBorder()} px-2 py-1 rounded text-xs"
								aria-label="Set rules logic to OR"
								title="Set rules logic to OR">
					OR
				</button>
				<button on:click={() => toggleLogic('AND')}
								class="logic-button {rulesOperator === 'AND' ? 'text-blue-500 border-blue-500' : ThemeUtils.themeTextMuted() + ' border ' + ThemeUtils.themeBorder()} px-2 py-1 rounded text-xs"
								aria-label="Set rules logic to AND"
								title="Set rules logic to AND">
					AND
				</button>
				<span class="{ThemeUtils.themeTextSecondary()} text-xs">
					{localRules.length} rule{localRules.length !== 1 ? 's' : ''}
				</span>
			</div>
			
			<div class="flex pr-4 shrink-0">
				<button 
					on:click={isEditing ? saveChanges : startEditing} 
					class="{isEditing && !hasChanges ? ThemeUtils.secondaryButton('py-1 px-2 text-xs opacity-60') : isEditing ? ThemeUtils.primaryButton('py-1 px-2 text-xs') : ThemeUtils.utilityButton('py-1 px-2 text-xs')}"
					type="button"
					disabled={isEditing && !hasChanges}
					aria-label={isEditing ? "Save rules" : "Edit rules"}
					title={isEditing ? "Save rules" : "Edit rules"}
				>
					{#if isEditing}
						<i class="fas fa-save text-xs mr-1"></i>
						<span>Save</span>
					{:else}
						<i class="fas fa-edit text-xs mr-1"></i>
						<span>Edit</span>
					{/if}
				</button>
				
				{#if isEditing}
					<button 
						on:click={cancelEdit} 
						class="{ThemeUtils.utilityButton('py-1 px-2 text-xs ml-1')}"
						type="button"
						aria-label="Cancel editing"
						title="Cancel editing"
					>
						Cancel
					</button>
				{/if}
			</div>
		</div>

		<!-- Loading State -->
		{#if isLoading}
			<div class="flex justify-center py-8">
				<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
			</div>
		
		<!-- View Mode -->
		{:else if !isEditing}
			{#if localRules && localRules.length > 0}
				<div bind:this={tableContainer} class="overflow-y-auto" style="max-height: {calculatedHeight}">
					<table class="w-full border-collapse table-fixed">
						<thead class="bg-gray-50 dark:bg-gray-750">
							<tr>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-1/6">Type</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}" style="flex: 1">Key</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-1/6">Operator</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}" style="flex: 1">Value</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-16">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y dark:divide-gray-700">
							{#each localRules as rule, index}
								<tr class="{index % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-2 align-top">
										<span class="inline-flex px-2 py-1 text-xs font-medium rounded-full {rule.type === 'header' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/30 dark:text-blue-300' : rule.type === 'query' ? 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300' : rule.type === 'body' ? 'bg-purple-100 text-purple-800 dark:bg-purple-900/30 dark:text-purple-300' : 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-300'}">{rule.type}</span>
									</td>
									<td class="px-4 py-2 align-top">
										{#if rule.type === 'body'}
											<span class="text-xs {ThemeUtils.themeTextMuted()} italic">Not required</span>
										{:else}
											<span class="font-medium text-blue-600 dark:text-blue-400 text-xs">{rule.key}</span>
										{/if}
									</td>
									<td class="px-4 py-2 align-top">
										<span class="text-xs {ThemeUtils.themeTextSecondary()}">{rule.operator}</span>
									</td>
									<td class="px-4 py-2 align-top">
										<div class="break-all whitespace-pre-wrap text-xs {ThemeUtils.themeTextSecondary()}">{rule.value}</div>
									</td>
									<td class="px-4 py-2 align-top">
										<button 
											on:click={() => removeRule(rule?.id, index)} 
											class="h-6 w-6 flex items-center justify-center text-red-500 hover:bg-red-100 dark:hover:bg-red-900/30 dark:text-red-400 dark:hover:text-red-300 rounded-full"
											title="Delete rule"
											aria-label="Delete rule"
										>
											<i class="fas fa-trash-alt text-xs"></i>
										</button>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{:else}
				<div class="text-center py-6 rounded-lg border border-dashed {ThemeUtils.themeBorder()}">
					<i class="fas fa-info-circle mb-2 text-lg {ThemeUtils.themeTextMuted()}"></i>
					<div class="{ThemeUtils.themeTextMuted()} italic text-xs">No rules defined. Response will always match requests.</div>
				</div>
			{/if}
		
		<!-- Edit Mode -->
		{:else}
			{#if localRules.length === 0}
				<div class="text-center py-6">
					<button
						on:click={addRule}
						class="{ThemeUtils.secondaryButton('py-1.5 px-3 text-xs')}"
						type="button"
						aria-label="Add first rule"
						title="Add first rule"
					>
						<i class="fas fa-plus text-xs mr-1"></i>
						<span>Add First Rule</span>
					</button>
				</div>
			{:else}
				<div bind:this={tableContainer} class="overflow-y-auto" style="max-height: {calculatedHeight}">
					<table class="w-full border-collapse table-fixed">
						<thead class="bg-gray-50 dark:bg-gray-750">
							<tr>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-1/6">Type</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}" style="flex: 1">Key</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-1/6">Operator</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()}" style="flex: 1">Value</th>
								<th class="text-left px-4 py-3 text-xs font-semibold {ThemeUtils.themeTextPrimary()} w-16">Actions</th>
							</tr>
						</thead>
						<tbody class="divide-y dark:divide-gray-700">
							{#each localRules as rule, i}
								<tr class="{i % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50/50 dark:bg-gray-750'}">
									<td class="px-4 py-2 align-top">
										<select 
											value={rule.type}
											on:change={(e) => handleTypeChange(i, (e.target as HTMLSelectElement).value)}
											class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
											aria-label="Rule type"
											title="Select rule type"
										>
											<option value="header">Header</option>
											<option value="query">Query</option>
											<option value="body">Body</option>
										</select>
									</td>
									<td class="px-4 py-2 align-top">
										<input 
											type="text"
											value={rule.key}
											on:input={(e) => handleKeyChange(i, (e.target as HTMLInputElement).value)}
											placeholder={rule.type === 'body' ? 'Not required for body rules' : 'Key name'}
											disabled={rule.type === 'body'}
											class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500 {rule.type === 'body' ? 'opacity-50 cursor-not-allowed' : ''}"
											aria-label="Rule key"
											title={rule.type === 'body' ? 'Key is not required for body rules' : 'Enter the key name'}
											bind:this={keyInputs[i]}
										/>
									</td>
									<td class="px-4 py-2 align-top">
										<select 
											value={rule.operator}
											on:change={(e) => handleOperatorChange(i, (e.target as HTMLSelectElement).value)}
											class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
											aria-label="Rule operator"
											title="Select rule operator"
										>
											<option value="equals">equals</option>
											<option value="contains">contains</option>
											<option value="regex">regex</option>
										</select>
									</td>
									<td class="px-4 py-2 align-top">
										<div class="flex items-center gap-2 w-full">
											<div class="flex-1">
												<input 
													type="text"
													value={rule.value}
													on:input={(e) => handleValueChange(i, (e.target as HTMLInputElement).value)}
													placeholder="Value to match"
													class="block w-full py-1 px-2 text-xs rounded bg-white dark:bg-gray-700 border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} focus:ring-1 focus:ring-blue-500/50 focus:border-blue-500"
													aria-label="Rule value"
													bind:this={valueInputs[i]}
													on:keydown={(e) => {
														if (e.key === 'Enter') {
															e.preventDefault();
															if (i === localRules.length - 1) {
																addRule();
															}
														}
														handleKeyDown(e, i);
													}}
												/>
											</div>
										</div>
									</td>
									<td class="px-4 py-2 align-top">
										<div class="flex items-center gap-1">
											{#if i === localRules.length - 1}
												<button
													on:click={() => addRule()}
													title="Add new rule"
													class="h-6 w-6 flex-shrink-0 flex items-center justify-center text-blue-600 hover:bg-blue-100 dark:hover:bg-blue-900/30 dark:text-blue-400 dark:hover:text-blue-300 rounded-full"
													type="button"
													aria-label="Add new rule"
												>
													<i class="fas fa-plus text-xs"></i>
												</button>
											{/if}
											<button
												on:click={() => removeLocalRule(i)}
												title="Remove rule"
												class="h-6 w-6 flex-shrink-0 flex items-center justify-center text-red-500 hover:bg-red-100 dark:hover:bg-red-900/30 dark:text-red-400 dark:hover:text-red-300 rounded-full"
												type="button"
												aria-label="Remove rule"
											>
												<i class="fas fa-trash-alt text-xs"></i>
											</button>
										</div>
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		{/if}
	</div>
</div>