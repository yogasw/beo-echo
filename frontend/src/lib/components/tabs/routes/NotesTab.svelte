<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	
	export let notes: string = '';
	export let onSaveNotes: (notes: string) => void;
	
	// Set the maximum character limit for notes
	const MAX_CHARS = 200;
	
	let charCount = notes.length;
	let isOverLimit = false;
	
	// Update character count whenever notes are changed
	$: {
		charCount = notes.length;
		isOverLimit = charCount > MAX_CHARS;
		
		// Auto-truncate if over limit
		if (isOverLimit) {
			notes = notes.substring(0, MAX_CHARS);
			charCount = MAX_CHARS;
			isOverLimit = false;
		}
	}
	
	function handleSave() {
		if (!isOverLimit) {
			onSaveNotes(notes);
		}
	}
</script>

<div class="flex flex-col h-full">
	<div class="mb-2 flex justify-between items-center"></div>
		<label for="response-notes" class="block text-sm font-medium {ThemeUtils.themeTextPrimary()}">
			Response Notes
		</label>
		<div class="flex items-center space-x-2">
			<span class="{isOverLimit ? 'text-red-500' : ThemeUtils.themeTextMuted()} text-xs">
				{charCount}/{MAX_CHARS}
			</span>
			<button
				on:click={handleSave}
				class="bg-blue-600 hover:bg-blue-700 text-white text-xs px-3 py-1 rounded disabled:opacity-50 disabled:cursor-not-allowed"
				disabled={isOverLimit}
				title="Save response notes"
				aria-label="Save response notes"
			>
				Save Notes
			</button>
		</div>
	</div>
	
	<textarea
		id="response-notes"
		bind:value={notes}
		maxlength={MAX_CHARS}
		class="w-full h-full min-h-[200px] rounded {ThemeUtils.themeBgSecondary()} px-4 py-3 {ThemeUtils.themeTextPrimary()} border {isOverLimit ? 'border-red-500' : ThemeUtils.themeBorder()}"
		placeholder="Add detailed notes about this response (max {MAX_CHARS} characters)..."
	></textarea>
	
	<div class="mt-2 flex justify-between items-center">
		<p class="text-xs {ThemeUtils.themeTextMuted()}">
			Notes are useful for documenting the purpose of this response, expected usage scenarios and any other important information for team members.
		</p>
		{#if isOverLimit}
			<p class="text-xs text-red-500">
				Character limit exceeded. Maximum is {MAX_CHARS} characters.
			</p>
		{/if}
	</div>
</div>
