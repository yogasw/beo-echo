<script lang="ts">
	import type { Endpoint } from '$lib/api/BeoApi';
	import { updateEndpoint } from '$lib/stores/saveButton';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let selectedEndpoint: Endpoint | null = null;

	// Default timeout value in milliseconds (30 seconds)
	const DEFAULT_TIMEOUT = 30000;

	// Get timeout from endpoint's advance_config field, or use default
	function getTimeoutValue(): number {
		if (!selectedEndpoint) return DEFAULT_TIMEOUT;
		
		try {
			// Parse the advance_config field if it exists
			if (selectedEndpoint.advance_config) {
				const config = JSON.parse(selectedEndpoint.advance_config);
				return config.timeout || DEFAULT_TIMEOUT;
			}
		} catch (e) {
			// If parsing fails, return default
			console.warn('Failed to parse endpoint advance_config:', e);
		}
		
		return DEFAULT_TIMEOUT;
	}

	// Update timeout value in endpoint advance_config
	function updateTimeout(newTimeout: number): void {
		if (!selectedEndpoint) return;

		let config: any = {};
		
		try {
			// Parse existing advance_config if it exists
			if (selectedEndpoint.advance_config) {
				config = JSON.parse(selectedEndpoint.advance_config);
			}
		} catch (e) {
			// If parsing fails, start with empty config
			config = {};
		}

		// Update the timeout value
		config.timeout = newTimeout;

		// Update the endpoint with the new advance_config
		selectedEndpoint = updateEndpoint('advance_config', JSON.stringify(config), selectedEndpoint);
	}

	$: currentTimeout = getTimeoutValue();
	$: timeoutSeconds = currentTimeout / 1000;
</script>

<div class="space-y-6 h-full">
	{#if selectedEndpoint}
		<div class="flex justify-between items-center mb-4">
			<h2 class="text-sm font-semibold {ThemeUtils.themeTextPrimary()} flex items-center">
				<i class="fas fa-cogs text-gray-500 mr-2"></i>
				Advanced Settings
			</h2>
		</div>

		<!-- Timeout Settings -->
		<div class="{ThemeUtils.card()}">
			<div class="p-4 border-b {ThemeUtils.themeBorder()} {ThemeUtils.themeBgSecondary()}">
				<h3 class="text-sm font-medium {ThemeUtils.themeTextPrimary()} flex items-center">
					<i class="fas fa-clock text-blue-500 mr-2"></i>
					Request Timeout
				</h3>
			</div>
			<div class="p-4 space-y-4">
				<!-- Timeout Configuration -->
				<div class="flex items-center justify-between p-3 {ThemeUtils.themeBgSecondary()} rounded-md">
					<div class="flex flex-col flex-1">
						<label 
							for="endpoint-timeout" 
							class="text-sm font-medium {ThemeUtils.themeTextPrimary()}"
						>
							Timeout (ms)
						</label>
						<span class="text-xs {ThemeUtils.themeTextMuted()}">
							Maximum time to wait for a response ({timeoutSeconds}s)
						</span>
						<span class="text-xs {ThemeUtils.themeTextMuted()} mt-1">
							This setting applies only to this endpoint when using proxy mode
						</span>
					</div>
					<div class="flex items-center space-x-2">
						<input
							id="endpoint-timeout"
							type="number"
							min="100"
							max="300000"
							step="1000"
							class="w-24 {ThemeUtils.themeBgSecondary()} p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border {ThemeUtils.themeBorder()} {ThemeUtils.themeTextSecondary()} text-sm transition-all duration-200"
							title="Request timeout in milliseconds"
							aria-label="Request timeout for this endpoint"
							value={currentTimeout}
							on:input={(e) => {
								const value = parseInt(e.currentTarget.value);
								if (!isNaN(value) && value >= 100) {
									updateTimeout(value);
								}
							}}
							on:blur={(e) => {
								const value = parseInt(e.currentTarget.value);
								if (isNaN(value) || value < 100) {
									e.currentTarget.value = currentTimeout.toString();
								}
							}}
						/>
						<span class="text-xs {ThemeUtils.themeTextMuted()}">ms</span>
					</div>
				</div>

				<!-- Timeout Information -->
				<div class="p-3 {ThemeUtils.themeBgTertiary()} rounded-md border-l-4 border-blue-500">
					<div class="flex items-start">
						<i class="fas fa-info-circle text-blue-500 mr-2 mt-0.5"></i>
						<div class="text-xs {ThemeUtils.themeTextMuted()}">
							<p class="mb-2">
								<strong>Endpoint-level timeout:</strong> This timeout setting applies specifically to this endpoint when it's configured to use proxy mode.
							</p>
							<p class="mb-2">
								<strong>Fallback behavior:</strong> If not set, the project-level timeout will be used as a fallback.
							</p>
							<p>
								<strong>Recommended range:</strong> 1-30 seconds (1000-30000ms) for most APIs. Higher values may be needed for long-running operations.
							</p>
						</div>
					</div>
				</div>

				<!-- Reset to Default -->
				<div class="flex justify-end">
					<button
						class="text-sm {ThemeUtils.themeBgSecondary()} {ThemeUtils.themeTextSecondary()} px-3 py-1 rounded hover:{ThemeUtils.themeBgSecondary()} transition-colors duration-200 border {ThemeUtils.themeBorder()}"
						title="Reset timeout to default value (30 seconds)"
						aria-label="Reset timeout to default"
						on:click={() => updateTimeout(DEFAULT_TIMEOUT)}
					>
						<i class="fas fa-undo mr-1"></i>
						Reset to Default
					</button>
				</div>
			</div>
		</div>

		<!-- Future Settings Placeholder -->
		<div class="{ThemeUtils.card()}">
			<div class="p-4 border-b {ThemeUtils.themeBorder()} {ThemeUtils.themeBgSecondary()}">
				<h3 class="text-sm font-medium {ThemeUtils.themeTextPrimary()} flex items-center">
					<i class="fas fa-tools text-gray-500 mr-2"></i>
					Additional Settings
				</h3>
			</div>
			<div class="p-4">
				<div class="text-center py-8 {ThemeUtils.themeTextMuted()}">
					<i class="fas fa-wrench text-2xl mb-2 opacity-50"></i>
					<p class="text-sm">More advanced settings will be available here in future updates.</p>
				</div>
			</div>
		</div>
	{:else}
		<div class="flex items-center justify-center h-64 {ThemeUtils.themeTextMuted()}">
			<div class="text-center">
				<i class="fas fa-route text-3xl mb-4 opacity-50"></i>
				<p>Select an endpoint to configure advanced settings</p>
			</div>
		</div>
	{/if}
</div>
