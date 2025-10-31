<script lang="ts">
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { toast } from '$lib/stores/toast';
	import { onMount } from 'svelte';
	import { getSystemConfig, updateSystemConfig } from '$lib/api/BeoApi';

	let provider: string = 'gemini';
	let apiKey: string = '';
	let apiEndpoint: string = '';
	let model: string = '';
	let isLoading: boolean = false;
	let isSaving: boolean = false;
	let showApiKey: boolean = false;
	let isCustomModel: boolean = false;

	// Popular models
	const geminiModels = [
		{ value: 'gemini-pro', label: 'Gemini Pro (Free)' },
		{ value: 'gemini-2.0-flash', label: 'Gemini 2.0 Flash (Fast)' }
	];

	const openaiModels = [
		{ value: 'gpt-3.5-turbo', label: 'GPT-3.5 Turbo' },
		{ value: 'gpt-4', label: 'GPT-4' },
		{ value: 'gpt-4-turbo', label: 'GPT-4 Turbo' }
	];

	// Use stored provider value
	$: activeProvider = provider || 'gemini';

	// Available models based on provider
	$: availableModels =
		activeProvider === 'gemini'
			? geminiModels
			: activeProvider === 'openai'
				? openaiModels
				: [];

	// Check if current model is in the list
	$: {
		if (model && availableModels.length > 0) {
			isCustomModel = !availableModels.some((m) => m.value === model);
		}
	}

	// Load current settings
	onMount(async () => {
		await loadSettings();
	});

	async function loadSettings() {
		try {
			isLoading = true;
			// Load all AI configs in parallel
			const [providerConfig, keyConfig, endpointConfig, modelConfig] = await Promise.all([
				getSystemConfig('AI_PROVIDER'),
				getSystemConfig('AI_API_KEY'),
				getSystemConfig('AI_API_ENDPOINT'),
				getSystemConfig('AI_MODEL')
			]);

			provider = providerConfig.value || 'gemini';
			apiKey = keyConfig.value || '';
			apiEndpoint = endpointConfig.value || '';
			model = modelConfig.value || '';
		} catch (error: any) {
			console.error('Failed to load AI settings:', error);
			toast.error('Failed to load AI settings');
		} finally {
			isLoading = false;
		}
	}

	async function saveSettings() {
		try {
			isSaving = true;

			// Save all configs in parallel
			await Promise.all([
				updateSystemConfig('AI_PROVIDER', provider),
				updateSystemConfig('AI_API_KEY', apiKey),
				updateSystemConfig('AI_API_ENDPOINT', apiEndpoint),
				updateSystemConfig('AI_MODEL', model)
			]);

			toast.success('AI settings saved successfully');
		} catch (error: any) {
			console.error('Failed to save AI settings:', error);
			toast.error('Failed to save AI settings');
		} finally {
			isSaving = false;
		}
	}

	function setGeminiDefaults() {
		provider = 'gemini';
		apiEndpoint = 'https://generativelanguage.googleapis.com/v1beta';
		model = 'gemini-pro';
		isCustomModel = false;
		toast.success('Gemini defaults applied');
	}

	function setOpenAIDefaults() {
		provider = 'openai';
		apiEndpoint = 'https://api.openai.com/v1';
		model = 'gpt-3.5-turbo';
		isCustomModel = false;
		toast.success('OpenAI defaults applied');
	}
</script>

<div class="p-6 space-y-6">
	<!-- Header -->
	<div class="flex items-start justify-between">
		<div>
			<h3 class="text-lg font-semibold {ThemeUtils.themeTextPrimary()} mb-2">
				AI Configuration
			</h3>
			<p class="text-sm {ThemeUtils.themeTextSecondary()}">
				Configure AI service for generating mock data. Supports Gemini (free) and OpenAI.
			</p>
		</div>
	</div>

	{#if isLoading}
		<div class="flex items-center justify-center py-8">
			<i class="fas fa-spinner fa-spin text-2xl {ThemeUtils.themeTextSecondary()}"></i>
		</div>
	{:else}
		<!-- Quick Setup with Active Indicator -->
		<div class="{ThemeUtils.themeBgSecondary()} rounded-lg p-4 border {ThemeUtils.themeBorder()}">
			<h4 class="text-sm font-medium {ThemeUtils.themeTextPrimary()} mb-3">Quick Setup</h4>
			<div class="flex gap-2">
				<!-- Gemini -->
				<button
					on:click={setGeminiDefaults}
					class="relative flex-1 px-4 py-2 bg-gradient-to-r from-blue-600 to-cyan-600 text-white rounded-lg hover:from-blue-700 hover:to-cyan-700 transition-all text-sm font-medium flex items-center justify-center gap-2 {activeProvider ===
					'gemini'
						? 'ring-2 ring-blue-400 ring-offset-1 dark:ring-offset-gray-800'
						: ''}"
				>
					{#if activeProvider === 'gemini'}
						<div
							class="absolute -top-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center shadow-lg"
						>
							<i class="fas fa-check text-white text-[10px]"></i>
						</div>
					{/if}
					<i class="fas fa-brain"></i>
					<span>Gemini (Free)</span>
				</button>

				<!-- OpenAI -->
				<button
					on:click={setOpenAIDefaults}
					class="relative flex-1 px-4 py-2 bg-gradient-to-r from-green-600 to-teal-600 text-white rounded-lg hover:from-green-700 hover:to-teal-700 transition-all text-sm font-medium flex items-center justify-center gap-2 {activeProvider ===
					'openai'
						? 'ring-2 ring-green-400 ring-offset-1 dark:ring-offset-gray-800'
						: ''}"
				>
					{#if activeProvider === 'openai'}
						<div
							class="absolute -top-1 -right-1 w-5 h-5 bg-green-500 rounded-full flex items-center justify-center shadow-lg"
						>
							<i class="fas fa-check text-white text-[10px]"></i>
						</div>
					{/if}
					<i class="fas fa-robot"></i>
					<span>OpenAI</span>
				</button>
			</div>
			{#if activeProvider === 'custom'}
				<div class="mt-2 text-xs {ThemeUtils.themeTextSecondary()} flex items-center gap-2">
					<i class="fas fa-info-circle"></i>
					<span>Custom AI provider configured</span>
				</div>
			{/if}
		</div>

		<!-- API Key -->
		<div>
			<label for="ai-api-key" class="block text-sm font-medium {ThemeUtils.themeTextPrimary()} mb-2">
				<i class="fas fa-key mr-2"></i>
				API Key
				<span class="text-red-500">*</span>
			</label>
			<div class="relative">
				<input
					id="ai-api-key"
					type={showApiKey ? 'text' : 'password'}
					bind:value={apiKey}
					placeholder="Enter your AI API key..."
					class="w-full px-4 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 pr-10"
				/>
				<button
					type="button"
					on:click={() => (showApiKey = !showApiKey)}
					class="absolute right-3 top-1/2 -translate-y-1/2 {ThemeUtils.themeTextSecondary()} hover:text-purple-500"
					aria-label="Toggle API key visibility"
				>
					<i class="fas {showApiKey ? 'fa-eye-slash' : 'fa-eye'}"></i>
				</button>
			</div>
			<p class="text-xs {ThemeUtils.themeTextSecondary()} mt-1">
				Get free Gemini API key from <a
					href="https://aistudio.google.com/app/apikey"
					target="_blank"
					class="text-blue-500 hover:underline">Google AI Studio</a
				>
			</p>
		</div>

		<!-- API Endpoint -->
		<div>
			<label for="ai-endpoint" class="block text-sm font-medium {ThemeUtils.themeTextPrimary()} mb-2">
				<i class="fas fa-server mr-2"></i>
				API Endpoint
				<span class="text-red-500">*</span>
			</label>
			<input
				id="ai-endpoint"
				type="text"
				bind:value={apiEndpoint}
				placeholder="https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent"
				class="w-full px-4 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500 font-mono text-sm"
			/>
			<p class="text-xs {ThemeUtils.themeTextSecondary()} mt-1">
				Full URL to the AI API endpoint
			</p>
		</div>

		<!-- Model -->
		<div>
			<label for="ai-model" class="block text-sm font-medium {ThemeUtils.themeTextPrimary()} mb-2">
				<i class="fas fa-cube mr-2"></i>
				Model
				<span class="text-red-500">*</span>
			</label>

			<div class="space-y-2">
				{#if availableModels.length > 0}
					<!-- Dropdown for popular models -->
					<select
						bind:value={model}
						on:change={() => (isCustomModel = false)}
						class="w-full px-4 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
						disabled={isCustomModel}
					>
						<option value="" disabled>Select a model...</option>
						{#each availableModels as modelOption}
							<option value={modelOption.value}>
								{modelOption.label}
							</option>
						{/each}
					</select>

					<!-- Toggle custom input -->
					<button
						type="button"
						on:click={() => {
							isCustomModel = !isCustomModel;
							if (isCustomModel) model = '';
						}}
						class="text-xs text-blue-500 hover:underline flex items-center gap-1"
					>
						<i class="fas {isCustomModel ? 'fa-list' : 'fa-edit'}"></i>
						{isCustomModel ? 'Use dropdown' : 'Use custom model'}
					</button>
				{/if}

				{#if isCustomModel || availableModels.length === 0}
					<!-- Manual input for custom model -->
					<input
						id="ai-model"
						type="text"
						bind:value={model}
						placeholder="Enter custom model name..."
						class="w-full px-4 py-2 {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextPrimary()} border {ThemeUtils.themeBorder()} rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
					/>
				{/if}
			</div>

			<p class="text-xs {ThemeUtils.themeTextSecondary()} mt-1">
				{#if activeProvider === 'gemini'}
					Gemini models available: gemini-pro (free), gemini-1.5-pro, gemini-1.5-flash
				{:else if activeProvider === 'openai'}
					OpenAI models available: gpt-3.5-turbo, gpt-4, gpt-4-turbo
				{:else}
					Enter your custom AI model name
				{/if}
			</p>
		</div>

		<!-- Info Box -->
		<div
			class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4"
		>
			<div class="flex items-start space-x-3">
				<i class="fas fa-info-circle text-blue-500 mt-1"></i>
				<div>
					<h4 class="text-sm font-medium text-blue-900 dark:text-blue-300 mb-1">
						How to use AI Chat
					</h4>
					<ul class="text-xs text-blue-700 dark:text-blue-400 space-y-1">
						<li>• Configure AI settings here first</li>
						<li>• Click "AI Chat" button in response body editor</li>
						<li>• Chat with AI to generate mock data</li>
						<li>• Responses automatically apply to editor</li>
					</ul>
				</div>
			</div>
		</div>

		<!-- Save Button -->
		<div class="flex justify-end pt-4 border-t {ThemeUtils.themeBorder()}">
			<button
				on:click={saveSettings}
				disabled={isSaving || !apiKey || !apiEndpoint || !model}
				class="px-6 py-2 bg-gradient-to-r from-purple-600 to-pink-600 text-white rounded-lg hover:from-purple-700 hover:to-pink-700 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center space-x-2 font-medium"
			>
				{#if isSaving}
					<i class="fas fa-spinner fa-spin"></i>
					<span>Saving...</span>
				{:else}
					<i class="fas fa-save"></i>
					<span>Save Settings</span>
				{/if}
			</button>
		</div>
	{/if}
</div>
