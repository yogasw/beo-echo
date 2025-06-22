<script lang="ts">
	import { onMount } from 'svelte';
	import { getAllSystemConfigs, updateSystemConfig, type SystemConfigItem } from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';
	import ToggleSwitch from '$lib/components/common/ToggleSwitch.svelte';
	import SkeletonLoader from '$lib/components/common/SkeletonLoader.svelte';
	import ErrorDisplay from '$lib/components/common/ErrorDisplay.svelte';

	// Component state
	let configs: SystemConfigItem[] = [];
	let isLoading = true;
	let error: Error | null = null;
	let searchQuery = '';
	let editingConfig: string | null = null;
	let editingValue = '';
	let showHelp = false;

	// Reactive filtered configs based on search
	$: filteredConfigs = configs.filter(config => {
		return searchQuery === '' || 
			config.key.toLowerCase().includes(searchQuery.toLowerCase()) ||
			config.description.toLowerCase().includes(searchQuery.toLowerCase());
	});

	// Get unique categories from actual configs
	$: availableCategories = [...new Set(configs.map(config => getCategoryFromConfig(config)))].sort();

	// Load all system configurations
	async function loadConfigs() {
		try {
			isLoading = true;
			error = null;
			configs = await getAllSystemConfigs();
			console.log('Configs loaded:', configs.length);
		} catch (err) {
			console.error('Failed to load configs:', err);
			error = err as Error;
			toast.error(err);
		} finally {
			isLoading = false;
		}
	}

	// Get category from config (based on key patterns or metadata)
	function getCategoryFromConfig(config: SystemConfigItem): string {
		const key = config.key.toLowerCase();
		
		if (key.startsWith('feature_')) return 'Features';
		if (key.includes('subdomain') || key.includes('domain')) return 'Domains';
		if (key.includes('log')) return 'Logging';
		if (key.includes('max_') || key.includes('limit') || key.includes('workspace')) return 'Limits';
		if (key.includes('jwt') || key.includes('secret') || key.includes('auth')) return 'Security';
		if (key.includes('response') || key.includes('default_response')) return 'Responses';
		
		return 'Other';
	}

	// Validate config value based on type
	function validateConfigValue(value: string, type: string): { isValid: boolean; message?: string } {
		if (!value.trim()) {
			return { isValid: false, message: 'Value cannot be empty' };
		}

		switch (type.toLowerCase()) {
			case 'boolean':
				if (value !== 'true' && value !== 'false') {
					return { isValid: false, message: 'Boolean value must be "true" or "false"' };
				}
				break;
			case 'number':
				const num = parseFloat(value);
				if (isNaN(num)) {
					return { isValid: false, message: 'Value must be a valid number' };
				}
				if (num < 0) {
					return { isValid: false, message: 'Number must be non-negative' };
				}
				break;
			case 'string':
				if (value.length > 1000) {
					return { isValid: false, message: 'String value is too long (max 1000 characters)' };
				}
				break;
		}

		return { isValid: true };
	}

	// Copy value to clipboard
	async function copyToClipboard(value: string, configKey: string) {
		try {
			await navigator.clipboard.writeText(value);
			toast.success(`Copied ${configKey} value to clipboard`);
		} catch (err) {
			toast.error('Failed to copy to clipboard');
		}
	}

	// Handle config value update
	async function handleConfigUpdate(config: SystemConfigItem, newValue: string) {
		try {
			await updateSystemConfig(config.key, newValue);
			
			// Update local state
			const configIndex = configs.findIndex(c => c.key === config.key);
			if (configIndex !== -1) {
				configs[configIndex].value = newValue;
				configs = [...configs]; // Trigger reactivity
			}
			
			toast.success(`Configuration "${config.key}" updated successfully`);
			editingConfig = null;
		} catch (err) {
			toast.error(`Failed to update configuration: ${err}`);
		}
	}

	// Handle boolean toggle - send update to backend
	async function handleToggle(config: SystemConfigItem, event: CustomEvent<{checked: boolean}>) {
		const newValue = event.detail.checked ? 'true' : 'false';
		await handleConfigUpdate(config, newValue);
	}

	// Start editing a config
	function startEditing(config: SystemConfigItem) {
		editingConfig = config.key;
		editingValue = config.value;
	}

	// Cancel editing
	function cancelEditing() {
		editingConfig = null;
		editingValue = '';
	}

	// Save edited value
	function saveEdit(config: SystemConfigItem) {
		const trimmedValue = editingValue.trim();
		
		if (!trimmedValue) {
			toast.warning('Value cannot be empty');
			return;
		}
		
		const validation = validateConfigValue(trimmedValue, config.type);
		if (!validation.isValid) {
			toast.error(validation.message || `Invalid ${config.type} value`);
			return;
		}
		
		handleConfigUpdate(config, trimmedValue);
	}

	// Handle Enter key in edit mode
	function handleKeydown(event: KeyboardEvent, config: SystemConfigItem) {
		if (event.key === 'Enter') {
			saveEdit(config);
		} else if (event.key === 'Escape') {
			cancelEditing();
		}
	}

	// Format config key for display
	function formatConfigKey(key: string): string {
		return key.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase());
	}

	// Get appropriate input type for config
	function getInputType(type: string): string {
		switch (type.toLowerCase()) {
			case 'number':
				return 'number';
			case 'boolean':
				return 'checkbox';
			default:
				return 'text';
		}
	}

	// Export configurations to JSON file
	function exportConfigs() {
		try {
			const exportData = {
				exported_at: new Date().toISOString(),
				total_configs: configs.length,
				configurations: configs.map(config => ({
					key: config.key,
					value: config.hide_value ? '[HIDDEN]' : config.value,
					type: config.type,
					description: config.description,
					category: getCategoryFromConfig(config)
				}))
			};

			const dataStr = JSON.stringify(exportData, null, 2);
			const dataBlob = new Blob([dataStr], { type: 'application/json' });
			const url = URL.createObjectURL(dataBlob);
			
			const link = document.createElement('a');
			link.href = url;
			link.download = `beo-echo-system-configs-${new Date().toISOString().split('T')[0]}.json`;
			document.body.appendChild(link);
			link.click();
			document.body.removeChild(link);
			URL.revokeObjectURL(url);

			toast.success('Configurations exported successfully');
		} catch (err) {
			toast.error('Failed to export configurations');
		}
	}

	// Load configs on mount
	onMount(() => {
		loadConfigs();
	});
</script>

<div class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-md overflow-hidden">
	<!-- Header -->
	<div class="flex justify-between items-center p-4 bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700">
		<div class="flex items-center">
			<i class="fas fa-cogs text-blue-400 mr-2"></i>
			<h2 class="text-lg font-semibold text-gray-800 dark:text-white">Advanced System Configuration</h2>
		</div>
		<div class="flex gap-2">
			<button 
				on:click={() => showHelp = true}
				class="bg-gray-600 hover:bg-gray-700 text-white px-3 py-1.5 rounded text-sm flex items-center"
				title="Show keyboard shortcuts and help"
				aria-label="Show help and keyboard shortcuts"
			>
				<i class="fas fa-question-circle mr-1"></i>
				Help
			</button>
			<button 
				on:click={exportConfigs}
				class="bg-gray-600 hover:bg-gray-700 text-white px-3 py-1.5 rounded text-sm flex items-center"
				title="Export configurations to JSON"
				aria-label="Export system configurations to JSON file"
				disabled={isLoading || configs.length === 0}
			>
				<i class="fas fa-download mr-1"></i>
				Export
			</button>
			<button 
				on:click={loadConfigs}
				class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded text-sm flex items-center"
				title="Refresh configurations"
				aria-label="Refresh system configurations"
				disabled={isLoading}
			>
				<i class="fas fa-sync {isLoading ? 'animate-spin' : ''} mr-1"></i>
				Refresh
			</button>
		</div>
	</div>

	<!-- Search and Filter Controls -->
	<div class="p-4 bg-gray-50 dark:bg-gray-900/50 border-b border-gray-200 dark:border-gray-700">
		<!-- Search Input -->
		<div class="relative">
			<div class="absolute inset-y-0 left-0 flex items-center pl-3">
				<i class="fas fa-search text-gray-400"></i>
			</div>
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Search configurations by key or description..."
				class="block w-full pl-10 pr-4 py-2 text-sm rounded-lg bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
				title="Search system configurations"
				aria-label="Search system configurations by key or description"
			/>
		</div>

		<!-- Results Count and Statistics -->
		{#if !isLoading && !error}
			<div class="mt-2 flex justify-between items-center text-xs text-gray-500 dark:text-gray-400">
				<span>Showing {filteredConfigs.length} of {configs.length} configurations</span>
				<div class="flex gap-4">
					{#each availableCategories as category}
						<span>{category}: {configs.filter(c => getCategoryFromConfig(c) === category).length}</span>
					{/each}
				</div>
			</div>
		{/if}
	</div>

	<!-- Configuration List -->
	<div class="max-h-96 overflow-y-auto">
		{#if isLoading}
			<div class="p-4">
				<SkeletonLoader type="list" count={6} />
			</div>
		{:else if error}
			<div class="p-4">
				<ErrorDisplay 
					message="Failed to load system configurations" 
					type="error" 
					retryable={true}
					onRetry={loadConfigs}
				/>
			</div>
		{:else if filteredConfigs.length === 0}
			<div class="p-8 text-center">
				<i class="fas fa-search text-gray-400 text-2xl mb-2"></i>
				<p class="text-gray-500 dark:text-gray-400">
					{searchQuery 
						? 'No configurations match your search criteria' 
						: 'No configurations available'}
				</p>
			</div>
		{:else}
			<div class="divide-y divide-gray-200 dark:divide-gray-700">
				{#each filteredConfigs as config (config.key)}
					<div class="p-4 hover:bg-gray-50 dark:hover:bg-gray-750 transition-colors">
						<div class="flex flex-col sm:flex-row sm:items-start justify-between gap-4">
							<!-- Config Info -->
							<div class="flex-1 min-w-0">
								<div class="flex flex-wrap items-center gap-2 mb-1">
									<h3 class="font-medium text-gray-900 dark:text-white truncate">
										{formatConfigKey(config.key)}
									</h3>
									<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 dark:bg-gray-700 text-gray-800 dark:text-gray-300">
										{getCategoryFromConfig(config)}
									</span>
									<span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-900/30 text-blue-800 dark:text-blue-300">
										{config.type}
									</span>
								</div>
								<p class="text-sm text-gray-600 dark:text-gray-400 mb-2">
									{config.description}
								</p>
								<div class="text-xs text-gray-500 dark:text-gray-500">
									Key: <code class="bg-gray-100 dark:bg-gray-700 px-1 rounded">{config.key}</code>
								</div>
							</div>

							<!-- Config Value Control -->
							<div class="flex items-center gap-2 w-full sm:w-auto">
								{#if config.type === 'boolean'}
									<!-- Boolean Toggle -->
									<ToggleSwitch
										checked={config.value === 'true'}
										on:change={(event) => {
                                            handleToggle(config, event)
                                        }}
                                        
										ariaLabel="Toggle {config.key}"
									/>
								{:else if editingConfig === config.key}
									<!-- Edit Mode -->
									<div class="flex items-center gap-2 w-full sm:w-auto">
										<input
											type={getInputType(config.type)}
											bind:value={editingValue}
											on:keydown={(event) => handleKeydown(event, config)}
											class="flex-1 sm:w-32 px-2 py-1 text-sm rounded border border-gray-300 dark:border-gray-600 bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:ring-blue-500 focus:border-blue-500"
											title="Edit {config.key} value"
											aria-label="Edit {config.key} value"
										/>
										<button
											on:click={() => saveEdit(config)}
											class="bg-green-600 hover:bg-green-700 text-white p-1.5 rounded text-xs"
											title="Save changes"
											aria-label="Save changes to {config.key}"
										>
											<i class="fas fa-check"></i>
										</button>
										<button
											on:click={cancelEditing}
											class="bg-gray-600 hover:bg-gray-700 text-white p-1.5 rounded text-xs"
											title="Cancel editing"
											aria-label="Cancel editing {config.key}"
										>
											<i class="fas fa-times"></i>
										</button>
									</div>
								{:else}
									<!-- Display Mode -->
									<div class="flex items-center gap-2 w-full sm:w-auto">
										<div class="text-sm font-mono bg-gray-100 dark:bg-gray-700 px-2 py-1 rounded min-w-0 flex-1 sm:flex-none truncate">
											{config.hide_value && !config.value ? 
												'[Hidden]' : 
												config.value || '[Empty]'}
										</div>
										{#if !config.hide_value}
											<button
												on:click={() => startEditing(config)}
												class="bg-blue-600 hover:bg-blue-700 text-white p-1.5 rounded text-xs"
												title="Edit {config.key}"
												aria-label="Edit {config.key} value"
											>
												<i class="fas fa-edit"></i>
											</button>
										{/if}
										<!-- Copy to Clipboard Button -->
										{#if config.value && !config.hide_value}
											<button
												on:click={() => copyToClipboard(config.value, config.key)}
												class="bg-green-600 hover:bg-green-700 text-white p-1.5 rounded text-xs"
												title="Copy {config.key} value to clipboard"
												aria-label="Copy {config.key} value to clipboard"
											>
												<i class="fas fa-copy"></i>
											</button>
										{/if}
									</div>
								{/if}
							</div>
						</div>
					</div>
				{/each}
			</div>
		{/if}
	</div>
</div>

<!-- Help Modal -->
{#if showHelp}
	<div 
		class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" 
		on:click={() => showHelp = false}
		on:keydown={(e) => e.key === 'Escape' && (showHelp = false)}
		role="dialog"
		aria-modal="true"
		aria-labelledby="help-modal-title"
		tabindex="-1"
	>
		<div 
			class="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg shadow-xl max-w-2xl w-full mx-4 max-h-96 overflow-y-auto" 
			role="document"
		>
			<div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
				<h3 id="help-modal-title" class="text-lg font-semibold text-gray-900 dark:text-white">System Configuration Help</h3>
				<button 
					on:click={() => showHelp = false}
					class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200"
					title="Close help"
					aria-label="Close help modal"
				>
					<i class="fas fa-times text-xl"></i>
				</button>
			</div>
			
			<div class="p-4 space-y-4 text-sm">
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Configuration Types</h4>
					<ul class="list-disc list-inside space-y-1 text-gray-700 dark:text-gray-300">
						<li><strong>Boolean:</strong> Toggle switches for true/false values</li>
						<li><strong>String:</strong> Text inputs for string values (max 1000 characters)</li>
						<li><strong>Number:</strong> Numeric inputs for integer/decimal values</li>
					</ul>
				</div>
				
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Categories</h4>
					<ul class="list-disc list-inside space-y-1 text-gray-700 dark:text-gray-300">
						<li><strong>Features:</strong> Feature flags and toggles</li>
						<li><strong>Security:</strong> Authentication and security settings</li>
						<li><strong>Domains:</strong> Domain and subdomain configurations</li>
						<li><strong>Limits:</strong> Resource limits and quotas</li>
						<li><strong>Logging:</strong> Log management settings</li>
						<li><strong>Responses:</strong> Default response configurations</li>
					</ul>
				</div>
				
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Keyboard Shortcuts</h4>
					<ul class="list-disc list-inside space-y-1 text-gray-700 dark:text-gray-300">
						<li><kbd class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Enter</kbd> - Save changes when editing</li>
						<li><kbd class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Escape</kbd> - Cancel editing</li>
						<li><kbd class="bg-gray-100 dark:bg-gray-700 px-1 rounded">Ctrl/Cmd + C</kbd> - Copy value (when copy button is clicked)</li>
					</ul>
				</div>
				
				<div>
					<h4 class="font-semibold text-gray-900 dark:text-white mb-2">Actions</h4>
					<ul class="list-disc list-inside space-y-1 text-gray-700 dark:text-gray-300">
						<li><strong>Edit:</strong> Click the edit icon to modify string/number values</li>
						<li><strong>Toggle:</strong> Click boolean toggles to switch true/false</li>
						<li><strong>Copy:</strong> Click copy icon to copy values to clipboard</li>
						<li><strong>Export:</strong> Download all configurations as JSON file</li>
					</ul>
				</div>
				
				<div class="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded p-3">
					<div class="flex">
						<i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400 mr-2 mt-0.5"></i>
						<div>
							<p class="text-yellow-800 dark:text-yellow-200 font-medium">Important Notes</p>
							<ul class="text-yellow-700 dark:text-yellow-300 text-xs mt-1 space-y-1">
								<li>• Changes are applied immediately</li>
								<li>• Hidden values are restricted to system owners</li>
								<li>• Some features may require application restart</li>
							</ul>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
