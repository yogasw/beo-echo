<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let settings = {
		timeout: 30000,
		followRedirects: true,
		maxRedirects: 5,
		verifySsl: true,
		ignoreSslErrors: false,
		encoding: 'utf-8',
		sendCookies: true,
		storeCookies: true,
		keepAlive: true,
		userAgent: 'Beo-Echo/1.0',
		retryOnFailure: false,
		retryCount: 3,
		retryDelay: 1000
	};

	const dispatch = createEventDispatcher();

	function updateSetting(key: string, value: any) {
		settings = { ...settings, [key]: value };
		dispatch('settingsChange', { settings });
	}

	function resetToDefaults() {
		settings = {
			timeout: 30000,
			followRedirects: true,
			maxRedirects: 5,
			verifySsl: true,
			ignoreSslErrors: false,
			encoding: 'utf-8',
			sendCookies: true,
			storeCookies: true,
			keepAlive: true,
			userAgent: 'Beo-Echo/1.0',
			retryOnFailure: false,
			retryCount: 3,
			retryDelay: 1000
		};
		dispatch('settingsChange', { settings });
	}

	// Encoding options
	const encodingOptions = [
		{ value: 'utf-8', label: 'UTF-8' },
		{ value: 'ascii', label: 'ASCII' },
		{ value: 'iso-8859-1', label: 'ISO-8859-1' },
		{ value: 'windows-1252', label: 'Windows-1252' },
		{ value: 'utf-16', label: 'UTF-16' }
	];

	// User agent presets
	const userAgentPresets = [
		{ value: 'Beo-Echo/1.0', label: 'Beo-Echo Default' },
		{ value: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36', label: 'Chrome (Windows)' },
		{ value: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36', label: 'Chrome (macOS)' },
		{ value: 'Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36', label: 'Chrome (Linux)' },
		{ value: 'curl/7.68.0', label: 'cURL' },
		{ value: 'PostmanRuntime/7.29.0', label: 'Postman' }
	];

	$: timeoutSeconds = settings.timeout / 1000;
	$: retryDelaySeconds = settings.retryDelay / 1000;
</script>

<!-- Settings section -->
<div role="tabpanel" aria-labelledby="settings-tab" class="space-y-6">
	<div class="flex justify-between items-center mb-4">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-cogs text-gray-500 mr-2"></i>
			Request Settings
		</h2>
		<button
			class="text-sm bg-gray-100 dark:bg-gray-700 theme-text-secondary px-3 py-1 rounded hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors duration-200"
			title="Reset all settings to defaults"
			aria-label="Reset to default settings"
			on:click={resetToDefaults}
		>
			<i class="fas fa-undo mr-1"></i>
			Reset to Defaults
		</button>
	</div>

	<!-- Request Behavior Settings -->
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<div class="p-4 border-b theme-border bg-gray-50 dark:bg-gray-750">
			<h3 class="text-sm font-medium theme-text-primary flex items-center">
				<i class="fas fa-globe text-blue-500 mr-2"></i>
				Request Behavior
			</h3>
		</div>
		<div class="p-4 space-y-4">
			<!-- Timeout -->
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="request-timeout" class="text-sm font-medium theme-text-primary">
						Request timeout
					</label>
					<span class="text-xs theme-text-muted">Maximum time to wait for a response ({timeoutSeconds}s)</span>
				</div>
				<div class="flex items-center space-x-2">
					<input
						id="request-timeout"
						type="number"
						min="100"
						max="300000"
						step="1000"
						class="w-24 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
						title="Request timeout in milliseconds"
						aria-label="Request timeout"
						value={settings.timeout}
						on:input={(e) => updateSetting('timeout', parseInt(e.currentTarget.value))}
					/>
					<span class="text-xs theme-text-muted">ms</span>
				</div>
			</div>

			<!-- Follow Redirects -->
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="follow-redirects" class="text-sm font-medium theme-text-primary cursor-pointer">
						Follow redirects
					</label>
					<span class="text-xs theme-text-muted">Automatically follow HTTP redirects (3xx responses)</span>
				</div>
				<div class="flex items-center">
					<input
						id="follow-redirects"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Follow HTTP redirects"
						aria-label="Enable following redirects"
						checked={settings.followRedirects}
						on:change={(e) => updateSetting('followRedirects', e.currentTarget.checked)}
					/>
				</div>
			</div>

			<!-- Max Redirects (only show if follow redirects is enabled) -->
			{#if settings.followRedirects}
				<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
					<div class="flex flex-col">
						<label for="max-redirects" class="text-sm font-medium theme-text-primary">
							Maximum redirects
						</label>
						<span class="text-xs theme-text-muted">Maximum number of redirects to follow</span>
					</div>
					<div class="flex items-center">
						<input
							id="max-redirects"
							type="number"
							min="1"
							max="20"
							class="w-16 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
							title="Maximum number of redirects"
							aria-label="Maximum redirects"
							value={settings.maxRedirects}
							on:input={(e) => updateSetting('maxRedirects', parseInt(e.currentTarget.value))}
						/>
					</div>
				</div>
			{/if}

			<!-- User Agent -->
			<div class="space-y-2">
				<label for="user-agent" class="block text-sm font-medium theme-text-primary">
					User Agent
				</label>
				<select
					id="user-agent-preset"
					class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200 mb-2"
					title="Select user agent preset"
					aria-label="User agent preset"
					value={settings.userAgent}
					on:change={(e) => updateSetting('userAgent', e.currentTarget.value)}
				>
					{#each userAgentPresets as preset}
						<option value={preset.value}>{preset.label}</option>
					{/each}
					<option value="custom">Custom...</option>
				</select>
				<input
					id="user-agent"
					type="text"
					class="w-full theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
					placeholder="Custom User-Agent string"
					title="Custom user agent string"
					aria-label="User agent string"
					value={settings.userAgent}
					on:input={(e) => updateSetting('userAgent', e.currentTarget.value)}
				/>
			</div>
		</div>
	</div>

	<!-- SSL/TLS Settings -->
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<div class="p-4 border-b theme-border bg-gray-50 dark:bg-gray-750">
			<h3 class="text-sm font-medium theme-text-primary flex items-center">
				<i class="fas fa-shield-alt text-green-500 mr-2"></i>
				SSL/TLS Options
			</h3>
		</div>
		<div class="p-4 space-y-3">
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="verify-ssl" class="text-sm font-medium theme-text-primary cursor-pointer">
						Verify SSL certificates
					</label>
					<span class="text-xs theme-text-muted">Validate SSL certificates for HTTPS requests</span>
				</div>
				<div class="flex items-center">
					<input
						id="verify-ssl"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Verify SSL certificates"
						aria-label="Enable SSL certificate verification"
						checked={settings.verifySsl}
						on:change={(e) => updateSetting('verifySsl', e.currentTarget.checked)}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="ignore-ssl-errors" class="text-sm font-medium theme-text-primary cursor-pointer">
						Ignore SSL errors
					</label>
					<span class="text-xs theme-text-muted">Continue with requests even if SSL verification fails</span>
				</div>
				<div class="flex items-center">
					<input
						id="ignore-ssl-errors"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Ignore SSL certificate errors"
						aria-label="Enable ignoring SSL errors"
						checked={settings.ignoreSslErrors}
						on:change={(e) => updateSetting('ignoreSslErrors', e.currentTarget.checked)}
					/>
				</div>
			</div>
		</div>
	</div>

	<!-- Cookie Settings -->
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<div class="p-4 border-b theme-border bg-gray-50 dark:bg-gray-750">
			<h3 class="text-sm font-medium theme-text-primary flex items-center">
				<i class="fas fa-cookie-bite text-orange-500 mr-2"></i>
				Cookie Management
			</h3>
		</div>
		<div class="p-4 space-y-3">
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="send-cookies" class="text-sm font-medium theme-text-primary cursor-pointer">
						Send cookies
					</label>
					<span class="text-xs theme-text-muted">Include cookies in the request</span>
				</div>
				<div class="flex items-center">
					<input
						id="send-cookies"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Send cookies with request"
						aria-label="Enable sending cookies"
						checked={settings.sendCookies}
						on:change={(e) => updateSetting('sendCookies', e.currentTarget.checked)}
					/>
				</div>
			</div>
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="store-cookies" class="text-sm font-medium theme-text-primary cursor-pointer">
						Store cookies
					</label>
					<span class="text-xs theme-text-muted">Automatically store cookies from responses</span>
				</div>
				<div class="flex items-center">
					<input
						id="store-cookies"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Store cookies from responses"
						aria-label="Enable storing cookies"
						checked={settings.storeCookies}
						on:change={(e) => updateSetting('storeCookies', e.currentTarget.checked)}
					/>
				</div>
			</div>
		</div>
	</div>

	<!-- Advanced Settings -->
	<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg overflow-hidden">
		<div class="p-4 border-b theme-border bg-gray-50 dark:bg-gray-750">
			<h3 class="text-sm font-medium theme-text-primary flex items-center">
				<i class="fas fa-cogs text-purple-500 mr-2"></i>
				Advanced Options
			</h3>
		</div>
		<div class="p-4 space-y-4">
			<!-- Response Encoding -->
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="encoding" class="text-sm font-medium theme-text-primary">
						Response encoding
					</label>
					<span class="text-xs theme-text-muted">Character encoding for response interpretation</span>
				</div>
				<div class="flex items-center">
					<select
						id="encoding"
						class="theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
						title="Select response encoding"
						aria-label="Response encoding"
						value={settings.encoding}
						on:change={(e) => updateSetting('encoding', e.currentTarget.value)}
					>
						{#each encodingOptions as option}
							<option value={option.value}>{option.label}</option>
						{/each}
					</select>
				</div>
			</div>

			<!-- Keep Alive -->
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="keep-alive" class="text-sm font-medium theme-text-primary cursor-pointer">
						Keep connection alive
					</label>
					<span class="text-xs theme-text-muted">Reuse connection for multiple requests</span>
				</div>
				<div class="flex items-center">
					<input
						id="keep-alive"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Keep connection alive"
						aria-label="Enable keep-alive"
						checked={settings.keepAlive}
						on:change={(e) => updateSetting('keepAlive', e.currentTarget.checked)}
					/>
				</div>
			</div>

			<!-- Retry on Failure -->
			<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
				<div class="flex flex-col">
					<label for="retry-on-failure" class="text-sm font-medium theme-text-primary cursor-pointer">
						Retry on failure
					</label>
					<span class="text-xs theme-text-muted">Automatically retry failed requests</span>
				</div>
				<div class="flex items-center">
					<input
						id="retry-on-failure"
						type="checkbox"
						class="form-checkbox h-4 w-4 theme-bg-secondary theme-border text-blue-500 focus:ring-blue-500 rounded transition-colors duration-200"
						title="Retry failed requests"
						aria-label="Enable retry on failure"
						checked={settings.retryOnFailure}
						on:change={(e) => updateSetting('retryOnFailure', e.currentTarget.checked)}
					/>
				</div>
			</div>

			<!-- Retry Settings (only show if retry is enabled) -->
			{#if settings.retryOnFailure}
				<div class="pl-4 space-y-3 border-l-2 border-blue-400">
					<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
						<div class="flex flex-col">
							<label for="retry-count" class="text-sm font-medium theme-text-primary">
								Retry attempts
							</label>
							<span class="text-xs theme-text-muted">Number of retry attempts</span>
						</div>
						<div class="flex items-center">
							<input
								id="retry-count"
								type="number"
								min="1"
								max="10"
								class="w-16 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
								title="Number of retry attempts"
								aria-label="Retry count"
								value={settings.retryCount}
								on:input={(e) => updateSetting('retryCount', parseInt(e.currentTarget.value))}
							/>
						</div>
					</div>
					<div class="flex items-center justify-between p-3 theme-bg-secondary rounded-md">
						<div class="flex flex-col">
							<label for="retry-delay" class="text-sm font-medium theme-text-primary">
								Retry delay
							</label>
							<span class="text-xs theme-text-muted">Delay between retry attempts ({retryDelaySeconds}s)</span>
						</div>
						<div class="flex items-center space-x-2">
							<input
								id="retry-delay"
								type="number"
								min="100"
								max="30000"
								step="100"
								class="w-20 theme-bg-secondary p-2 focus:outline-none focus:ring-2 focus:ring-blue-500 rounded border theme-border theme-text-secondary text-sm transition-all duration-200"
								title="Retry delay in milliseconds"
								aria-label="Retry delay"
								value={settings.retryDelay}
								on:input={(e) => updateSetting('retryDelay', parseInt(e.currentTarget.value))}
							/>
							<span class="text-xs theme-text-muted">ms</span>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>

	<!-- Settings Summary -->
	<div class="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
		<h4 class="text-sm font-medium theme-text-primary mb-2 flex items-center">
			<i class="fas fa-info-circle text-blue-400 mr-2"></i>
			Settings Summary
		</h4>
		<div class="text-xs theme-text-secondary space-y-1">
			<div>Timeout: <span class="font-medium">{timeoutSeconds}s</span></div>
			<div>Redirects: <span class="font-medium">{settings.followRedirects ? `Follow up to ${settings.maxRedirects}` : 'Disabled'}</span></div>
			<div>SSL Verification: <span class="font-medium">{settings.verifySsl ? 'Enabled' : 'Disabled'}</span></div>
			<div>Cookies: <span class="font-medium">{settings.sendCookies ? 'Send & ' : ''}{settings.storeCookies ? 'Store' : 'Don\'t store'}</span></div>
			<div>Retry: <span class="font-medium">{settings.retryOnFailure ? `${settings.retryCount} attempts with ${retryDelaySeconds}s delay` : 'Disabled'}</span></div>
		</div>
	</div>
</div>
