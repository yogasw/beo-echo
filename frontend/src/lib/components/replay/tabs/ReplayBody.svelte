<script lang="ts">
	// Effect for MutationObserver and initial check
	$effect(() => {});
	let { payload }: { payload?: string } = $props();
</script>

<!-- Body section -->
<div role="tabpanel" aria-labelledby="body-tab" class="space-y-4">
	<div class="flex items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-file-alt text-purple-500 mr-2"></i>
			Request Body
		</h2>
	</div>
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4">
		<div class="space-y-4">
			<fieldset>
				<legend class="sr-only">Body type selection</legend>
				<div class="flex flex-wrap gap-4 text-sm">
					<label class="flex items-center cursor-pointer">
						<input
							type="radio"
							name="bodyType"
							value="none"
							class="mr-2 text-blue-500 focus:ring-blue-500"
							checked
							title="No request body"
							aria-label="No body"
						/>
						<span class="theme-text-secondary">None</span>
					</label>
					<label class="flex items-center cursor-pointer">
						<input
							type="radio"
							name="bodyType"
							value="raw"
							class="mr-2 text-blue-500 focus:ring-blue-500"
							title="Raw text or JSON body"
							aria-label="Raw body"
						/>
						<span class="theme-text-secondary">Raw</span>
					</label>
					<label class="flex items-center cursor-pointer">
						<input
							type="radio"
							name="bodyType"
							value="form"
							class="mr-2 text-blue-500 focus:ring-blue-500"
							title="Form data body"
							aria-label="Form data body"
						/>
						<span class="theme-text-secondary">Form Data</span>
					</label>
					<label class="flex items-center cursor-pointer">
						<input
							type="radio"
							name="bodyType"
							value="urlencoded"
							class="mr-2 text-blue-500 focus:ring-blue-500"
							title="URL encoded form body"
							aria-label="URL encoded body"
						/>
						<span class="theme-text-secondary">x-www-form-urlencoded</span>
					</label>
				</div>
			</fieldset>

			<div>
				<label for="request-body" class="block text-sm font-medium theme-text-secondary mb-2">
					Body Content
				</label>
				<textarea
					id="request-body"
					class="w-full h-64 theme-bg-secondary p-3 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded-md border theme-border theme-text-secondary font-mono text-sm transition-all duration-200"
					placeholder="Enter request body content..."
					title="Request body content"
					aria-label="Request body textarea"
				></textarea>
				<div class="flex justify-between items-center mt-2">
					<p class="text-xs theme-text-muted">
						<i class="fas fa-info-circle mr-1"></i>
						Supports JSON, XML, text, and other formats
					</p>
					<button
						class="text-xs text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200"
						title="Format and beautify JSON content"
						aria-label="Beautify JSON"
					>
						<i class="fas fa-magic mr-1"></i>
						Beautify
					</button>
				</div>
			</div>
		</div>
	</div>
</div>
