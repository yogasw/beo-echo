<script lang="ts">
	import { onMount } from 'svelte';
	import { generateContent } from '$lib/api/aiApi';
	import { DEFAULT_GENERATE_PROMPT } from '$lib/constants/aiTemplates';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let isOpen: boolean = false;
	export let initialContext: string = '';
	export let onSave: (content: string) => void = () => {};
	export let onClose: () => void = () => {};

	let MonacoEditor: any;
	let editorRef: any;
	let generatedContent: string = '{}';
	let customPrompt: string = '';
	let isGenerating: boolean = false;
	let showCustomPrompt: boolean = false;
	let editorLoaded: boolean = false;

	// Dynamically import Monaco Editor to avoid SSR issues
	onMount(async () => {
		try {
			const module = await import('$lib/components/MonacoEditor.svelte');
			MonacoEditor = module.default;
			editorLoaded = true;
		} catch (err) {
			console.error('Failed to load Monaco Editor:', err);
		}
	});

	// Initialize content when modal opens
	$: if (isOpen && initialContext) {
		generatedContent = initialContext;
	}

	async function handleGenerate() {
		if (isGenerating) return;

		try {
			isGenerating = true;

			const prompt = showCustomPrompt && customPrompt.trim()
				? customPrompt
				: DEFAULT_GENERATE_PROMPT;

			const response = await generateContent({
				template: prompt,
				context: generatedContent !== '{}' ? generatedContent : undefined
			});

			generatedContent = response.content;

			// Try to format as JSON
			try {
				const parsed = JSON.parse(response.content);
				generatedContent = JSON.stringify(parsed, null, 2);
			} catch (e) {
				// If not JSON, keep as is
			}

			toast.success(`Generated using ${response.model} (${response.token_used} tokens)`);
		} catch (error: any) {
			const errorMsg = error?.response?.data?.error || error?.message || 'Unknown error';
			toast.error(`Failed to generate: ${errorMsg}`);
		} finally {
			isGenerating = false;
		}
	}

	function handleSave() {
		const content = editorRef?.getValue ? editorRef.getValue() : generatedContent;
		onSave(content);
		handleClose();
	}

	function handleClose() {
		generatedContent = '{}';
		customPrompt = '';
		showCustomPrompt = false;
		onClose();
	}

	function formatContent() {
		if (editorRef?.format) {
			editorRef.format();
		}
	}

	function handleKeyDown(event: KeyboardEvent) {
		if (event.key === 'Escape' && isOpen) {
			handleClose();
		}
	}
</script>

{#if isOpen}
	<div
		class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-70"
		role="dialog"
		aria-label="AI Generator Modal"
		on:keydown={handleKeyDown}
		tabindex="0"
	>
		<div
			class="relative w-11/12 h-5/6 {ThemeUtils.themeBgPrimary()} rounded-lg shadow-xl flex flex-col"
		>
			<!-- Header -->
			<div class="flex items-center justify-between p-4 border-b {ThemeUtils.themeBorderColor()}">
				<div class="flex items-center space-x-3">
					<div
						class="w-10 h-10 bg-gradient-to-br from-purple-500 to-pink-500 rounded-lg flex items-center justify-center"
					>
						<i class="fas fa-magic text-white"></i>
					</div>
					<div>
						<h2 class="text-xl font-bold {ThemeUtils.themeTextPrimary()}">AI Generator</h2>
						<p class="text-sm {ThemeUtils.themeTextSecondary()}">
							Generate mock response data with AI
						</p>
					</div>
				</div>
				<button
					on:click={handleClose}
					class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors"
					aria-label="Close modal"
				>
					<i class="fas fa-times text-xl"></i>
				</button>
			</div>

			<!-- Custom Prompt Section (Optional) -->
			<div class="p-4 border-b {ThemeUtils.themeBorderColor()}">
				<div class="flex items-center justify-between mb-2">
					<label class="text-sm font-medium {ThemeUtils.themeTextPrimary()}">
						Custom Prompt (Optional)
					</label>
					<button
						on:click={() => (showCustomPrompt = !showCustomPrompt)}
						class="text-xs text-blue-600 dark:text-blue-400 hover:underline"
					>
						{showCustomPrompt ? 'Hide' : 'Show'} Custom Prompt
					</button>
				</div>

				{#if showCustomPrompt}
					<textarea
						bind:value={customPrompt}
						placeholder="Enter custom instructions for AI generation... (leave empty to use default)"
						class="w-full h-20 px-3 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorderColor()} rounded-lg resize-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors text-sm"
					></textarea>
				{/if}
			</div>

			<!-- Action Buttons -->
			<div class="flex items-center justify-between p-4 border-b {ThemeUtils.themeBorderColor()}">
				<div class="flex space-x-2">
					<button
						on:click={handleGenerate}
						disabled={isGenerating}
						class="px-4 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg hover:from-purple-700 hover:to-pink-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center space-x-2 text-sm font-medium"
					>
						{#if isGenerating}
							<i class="fas fa-spinner fa-spin"></i>
							<span>Generating...</span>
						{:else}
							<i class="fas fa-magic"></i>
							<span>{generatedContent === '{}' ? 'Generate' : 'Regenerate'}</span>
						{/if}
					</button>

					<button
						on:click={formatContent}
						class="px-3 py-2 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors"
						aria-label="Format JSON"
						title="Format JSON"
					>
						<i class="fas fa-code mr-1"></i>
						Prettify
					</button>
				</div>

				<div class="flex space-x-2">
					<button
						on:click={handleSave}
						class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm font-medium"
					>
						<i class="fas fa-check mr-2"></i>
						Save & Close
					</button>
					<button
						on:click={handleClose}
						class="px-4 py-2 bg-gray-600 text-white rounded-lg hover:bg-gray-700 transition-colors text-sm font-medium"
					>
						Cancel
					</button>
				</div>
			</div>

			<!-- Monaco Editor -->
			<div class="flex-grow p-4 overflow-hidden">
				{#if editorLoaded && MonacoEditor}
					<svelte:component
						this={MonacoEditor}
						bind:this={editorRef}
						value={generatedContent}
						language="json"
					/>
				{:else}
					<div class="w-full h-full flex items-center justify-center {ThemeUtils.themeBgSecondary()} rounded-lg">
						<div class="text-center">
							<i class="fas fa-spinner fa-spin text-3xl {ThemeUtils.themeTextSecondary()} mb-2"></i>
							<p class="{ThemeUtils.themeTextSecondary()}">Loading editor...</p>
						</div>
					</div>
				{/if}
			</div>
		</div>
	</div>
{/if}