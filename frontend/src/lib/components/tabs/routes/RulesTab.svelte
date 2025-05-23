<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast';
	import type { Rule } from '$lib/api/rulesApi';
	import { getRules, createRule, updateRule, deleteRule } from '$lib/api/rulesApi';

	// Props
	export let projectId: string;
	export let endpointId: string;
	export let responseId: string;
	export let rules: Rule[] = [];
	export let rulesOperator = 'OR'; // Default operator

	// Local state
	let isModalVisible = false;
	let isLoading = false;
	let newRule: Omit<Rule, 'id'> = {
		responseId: responseId,
		type: 'header',
		key: '',
		operator: 'equals',
		value: ''
	};

	// Load rules on component mount
	onMount(async () => {
		await loadRules();
	});

	// Load rules from API
	async function loadRules() {
		isLoading = true;
		try {
			rules = await getRules(projectId, endpointId, responseId);
			console.log('Rules loaded:', rules);
		} catch (error) {
			console.error('Error loading rules:', error);
			toast.error('Failed to load rules');
		} finally {
			isLoading = false;
		}
	}

	// Toggle modal visibility
	function toggleModal() {
		isModalVisible = !isModalVisible;
		// Reset form when opening modal
		if (isModalVisible) {
			newRule = {
				responseId: responseId,
				type: 'header',
				key: '',
				operator: 'equals',
				value: ''
			};
		}
	}

	// Add a new rule
	async function addRule() {
		if (!newRule.key) {
			toast.error('Rule key is required');
			return;
		}

		try {
			const createdRule = await createRule(projectId, endpointId, responseId, newRule);
			rules = [...rules, createdRule];
			toast.success('Rule added successfully');
			toggleModal(); // Close the modal after adding the rule
		} catch (error) {
			console.error('Error adding rule:', error);
			toast.error('Failed to add rule');
		}
	}

	// Update an existing rule
	async function updateExistingRule(rule: Rule, index: number) {
		try {
			if (!rule.id) {
				toast.error('Rule ID is missing');
				return;
			}
			
			const updatedRule = await updateRule(
				projectId, 
				endpointId, 
				responseId, 
				rule.id, 
				{
					type: rule.type,
					key: rule.key,
					operator: rule.operator,
					value: rule.value
				}
			);
			
			// Update the rule in the local array
			rules[index] = updatedRule;
			rules = [...rules]; // Trigger reactivity
			toast.success('Rule updated successfully');
		} catch (error) {
			console.error('Error updating rule:', error);
			toast.error('Failed to update rule');
		}
	}

	// Delete a rule
	async function removeRule(ruleId: string, index: number) {
		try {
			await deleteRule(projectId, endpointId, responseId, ruleId);
			// Remove the rule from the local array
			rules = rules.filter((_, i) => i !== index);
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

	// Logs on component update
	$: {
		console.log('rules updated', rules);
		console.log('rules operator', rulesOperator);
	}
</script>


<div class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full max-w-4xl">
	<!-- Rules Configuration Section -->
	<div>
		<div class="{ThemeUtils.themeBgSecondary()} rounded-lg p-4">
			<div class="flex items-center justify-between mb-4">
				<div class="flex items-center space-x-4">
					<button on:click={() => toggleLogic('OR')}
									class="logic-button {rulesOperator === 'OR' ? 'text-blue-500 border-blue-500' : ThemeUtils.themeTextMuted() + ' border ' + ThemeUtils.themeBorder()} px-2 py-1 rounded">
						OR
					</button>
					<button on:click={() => toggleLogic('AND')}
									class="logic-button {rulesOperator === 'AND' ? 'text-blue-500 border-blue-500' : ThemeUtils.themeTextMuted() + ' border ' + ThemeUtils.themeBorder()} px-2 py-1 rounded">
						AND
					</button>
				</div>
				
				<div class="{ThemeUtils.themeTextSecondary()} text-sm">
					{#if rules.length === 0}
						No rules defined - response will always match
					{:else}
						{rules.length} rule{rules.length > 1 ? 's' : ''} defined - all must {rulesOperator === 'AND' ? 'match' : 'be evaluated'} for response selection
					{/if}
				</div>
			</div>
			
			{#if isLoading}
				<div class="flex justify-center py-4">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
				</div>
			{:else if rules.length === 0}
				<div class="py-4 text-center {ThemeUtils.themeTextMuted()}">
					No rules defined yet. Add a rule to conditionally match this response to specific requests.
				</div>
			{:else}
				<div class="space-y-4">
					{#each rules as rule, index}
						<div class="flex items-center space-x-4 w-full bg-opacity-50 {ThemeUtils.themeBgPrimary()} p-3 rounded">
							<div class="w-1/6">
								<select bind:value={rule.type}
												on:change={() => updateExistingRule(rule, index)}
												class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}">
									<option value="header">Header</option>
									<option value="query">Query</option>
									<option value="body">Body</option>
								</select>
							</div>
							<input type="text" bind:value={rule.key}
										on:blur={() => updateExistingRule(rule, index)}
										placeholder="Key name"
										class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
										
							<select bind:value={rule.operator}
											on:change={() => updateExistingRule(rule, index)}
											class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-1/6 border {ThemeUtils.themeBorder()}">
								<option value="equals">equals</option>
								<option value="contains">contains</option>
								<option value="regex">regex</option>
							</select>
							
							<input type="text" bind:value={rule.value}
										on:blur={() => updateExistingRule(rule, index)}
										placeholder="Value to match"
										class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
										
							<button on:click={() => removeRule(rule.id, index)} 
											class="text-red-500 hover:text-red-600">
								<i class="fas fa-trash"></i>
							</button>
						</div>
					{/each}
				</div>
			{/if}
			
			<button on:click={toggleModal} class="text-green-500 mt-4 flex items-center hover:text-green-600">
				<i class="fas fa-plus-circle mr-2"></i> Add rule
			</button>
		</div>
	</div>
</div>

{#if isModalVisible}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center z-50">
		<div id="addRuleModal" class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full max-w-md p-6">
			<h2 class="text-xl font-semibold mb-4 {ThemeUtils.themeTextPrimary()}">Add New Rule</h2>
			<div class="space-y-4">
				<div>
					<label class="block {ThemeUtils.themeTextMuted()} mb-2">Type</label>
					<select bind:value={newRule.type}
									class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}">
						<option value="header">Header</option>
						<option value="query">Query</option>
						<option value="body">Body</option>
					</select>
				</div>
				
				<div>
					<label class="block {ThemeUtils.themeTextMuted()} mb-2">Key</label>
					<input type="text" bind:value={newRule.key}
								 placeholder="e.g. Authorization, user.id, version"
								 class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
				</div>
				
				<div>
					<label class="block {ThemeUtils.themeTextMuted()} mb-2">Operator</label>
					<select bind:value={newRule.operator}
									class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}">
						<option value="equals">equals</option>
						<option value="contains">contains</option>
						<option value="regex">regex</option>
					</select>
				</div>
				
				<div>
					<label class="block {ThemeUtils.themeTextMuted()} mb-2">Value</label>
					<input type="text" bind:value={newRule.value}
								 placeholder="Value to match"
								 class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
				</div>
			</div>
			
			<div class="flex justify-end mt-4 space-x-2">
				<button on:click={toggleModal} 
								class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextMuted()} rounded px-4 py-2 hover:opacity-80">
					Cancel
				</button>
				<button on:click={addRule} 
								class="bg-green-500 text-white rounded px-4 py-2 hover:bg-green-600">
					Add
				</button>
			</div>
		</div>
	</div>
{/if}