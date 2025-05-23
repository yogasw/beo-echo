<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let rules: any[] = [];
	export let rulesOperator = 'OR'; // Default operator

	let isModalVisible = false;
	const disableAddButton = true;

	function toggleModal() {
		isModalVisible = !isModalVisible;
	}

	function addRule() {
		// Logic to add a new rule
		const newRule = {
			type: 'Body', // Default values, can be updated based on modal inputs
			key: '',
			operator: 'equals',
			value: ''
		};
		rules = [...rules, newRule];
		toggleModal(); // Close the modal after adding the rule
	}

	function toggleLogic(button: any) {
		if(disableAddButton) return
		const buttons = document.querySelectorAll('.logic-button');
		buttons.forEach((btn) => {
			btn.classList.remove('text-blue-500', 'border-blue-500');
			btn.classList.add('text-gray-400', 'border-gray-500');
		});
		button.classList.remove('text-gray-400', 'border-gray-500');
		button.classList.add('text-blue-500', 'border-blue-500');
	}

	//logs on component update
	$: {
		console.log('rules updated', rules);
	}
</script>


<div class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full max-w-4xl">
	<!-- Multi Config Section -->
	<div>
		<div class="{ThemeUtils.themeBgSecondary()} rounded-lg p-4">
			<div class="flex items-center space-x-4 mb-4">
				<button on:click={toggleLogic}
								class="logic-button text-blue-500 border border-blue-500 px-2 py-1 rounded">OR
				</button>
				<button on:click={toggleLogic}
								class="logic-button {ThemeUtils.themeTextMuted()} border {ThemeUtils.themeBorder()} px-2 py-1 rounded">AND
				</button>
			</div>
			<div class="space-y-4">
				{#each rules as rule, index}
					<div class="flex items-center space-x-4 w-full">
						<div class="w-1/6 {ThemeUtils.themeTextMuted()}">{rule.target}</div>
						<input type="text" value="{rule.modifier}"
									 class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
						<button class="rounded px-2 py-1 border {ThemeUtils.themeBorder()}"
										class:text-blue-500={rule.invert}
										class:{ThemeUtils.themeTextMuted()}={!rule.invert}>!
						</button>
						<select class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-1/6 border {ThemeUtils.themeBorder()}"
										bind:value="{rule.operator}">
							<option>equals</option>
							<option>regex</option>
							<option>regex (i)</option>
							<option>null</option>
							<option>empty array</option>
						</select>
						<input type="text" value="{rule.value}"
									 class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
						<button class="text-gray-400" class:hidden={disableAddButton}>
							<i class="fas fa-trash"></i> Delete rule
						</button>
					</div>
				{/each}
			</div>
			<button on:click={toggleModal} class="text-green-500 mt-4 flex items-center" class:hidden={disableAddButton}>
				<i class="fas fa-plus-circle mr-2"></i> Add rule
			</button>
		</div>
	</div>
</div>
{#if isModalVisible}
	<div class="fixed inset-0 bg-black bg-opacity-75 flex items-center justify-center">
		<div id="addRuleModal" class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full max-w-md p-6">
			<div class="{ThemeUtils.themeBgPrimary()} rounded-lg w-full max-w-md p-6">
				<h2 class="text-xl font-semibold mb-4 {ThemeUtils.themeTextPrimary()}">Add New Rule</h2>
				<div class="space-y-4">
					<div>
						<label class="block {ThemeUtils.themeTextMuted()} mb-2">Type</label>
						<select class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}">
							<option>Body</option>
							<option>Query string</option>
							<option>Header</option>
							<option>Cookie</option>
							<option>Route params</option>
							<option>Request number</option>
						</select>
					</div>
					<div>
						<label class="block {ThemeUtils.themeTextMuted()} mb-2">Key</label>
						<input type="text" class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
					</div>
					<div>
						<label class="block {ThemeUtils.themeTextMuted()} mb-2">Operator</label>
						<select class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}">
							<option>equals</option>
							<option>regex</option>
							<option>regex (i)</option>
							<option>null</option>
							<option>empty array</option>
						</select>
					</div>
					<div>
						<label class="block {ThemeUtils.themeTextMuted()} mb-2">Value</label>
						<input type="text" class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} rounded px-2 py-1 w-full border {ThemeUtils.themeBorder()}" />
					</div>
				</div>
				<div class="flex justify-end mt-4 space-x-2">
					<button on:click={toggleModal} class="{ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextMuted()} rounded px-4 py-2">Cancel</button>
					<button on:click={addRule} class="bg-green-500 text-white rounded px-4 py-2">Add</button>
				</div>
			</div>
		</div>
	</div>
{/if}