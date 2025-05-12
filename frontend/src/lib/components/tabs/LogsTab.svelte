<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { createLogStream, getLogs, type Project, type RequestLog } from '$lib/api/BeoApi';
	import { fade } from 'svelte/transition';
	import ModalCreateMock from './logs/ModalCreateMock.svelte';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { currentWorkspace } from '$lib/stores/workspace';

	export let selectedProject: Project;

	let logs: RequestLog[] = [];
	let isLoading = true;
	let error: string | null = null;
	let searchTerm = '';
	let total = 0;
	let page = 1;
	const pageSize = 100;
	let eventSource: EventSource | null = null;
	let autoScroll = true;
	// Map to track expanded logs
	let expandedLogs: Record<string, boolean> = {};
	// Map to track active tab (request/response)
	let activeTabs: Record<string, 'request' | 'response'> = {};
	// Notification for copy operations
	let copyNotification = { show: false, message: '' };
	// State for create mock modal
	let isCreateMockModalOpen = false;
	let selectedLogForMock: RequestLog | null = null;

	// Function to toggle log expansion
	function toggleLogExpansion(logId: string) {
		expandedLogs[logId] = !expandedLogs[logId];
		if (expandedLogs[logId] && !activeTabs[logId]) {
			activeTabs[logId] = 'request';
		}
		expandedLogs = expandedLogs; // Force Svelte reactivity update
		activeTabs = activeTabs; // Force Svelte reactivity update
	}

	// Function to switch between request and response tabs
	function switchTab(logId: string, tab: 'request' | 'response') {
		activeTabs[logId] = tab;
		activeTabs = activeTabs; // Force Svelte reactivity update
	}

	// Function to copy content to clipboard
	async function copyToClipboard(text: string, label: string) {
		try {
			await navigator.clipboard.writeText(text);
			copyNotification = { show: true, message: `${label} copied to clipboard!` };
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
			copyNotification = { show: true, message: 'Failed to copy to clipboard' };
			setTimeout(() => {
				copyNotification = { show: false, message: '' };
			}, 2000);
		}
	}

	// Function to pretty format JSON
	function formatJson(jsonStr: string): string {
		try {
			return JSON.stringify(JSON.parse(jsonStr), null, 2);
		} catch (e) {
			return jsonStr;
		}
	}

	// Function to minify JSON
	function minifyJson(jsonStr: string): string {
		try {
			return JSON.stringify(JSON.parse(jsonStr));
		} catch (e) {
			return jsonStr;
		}
	}

	// Function to check if all search terms are present in a log
	function matchesAllSearchTerms(log: RequestLog, searchTerms: string[]): boolean {
		if (searchTerms.length === 0) return true;

		// Combine all searchable fields into one string for easier searching
		const searchableText = [
			log.path.toLowerCase(),
			log.method.toLowerCase(),
			log.request_body.toLowerCase(),
			log.response_body.toLowerCase()
		].join(' ');

		// Check if all search terms are present in the searchable text
		return searchTerms.every((term) => searchableText.includes(term));
	}

	$: searchTerms = searchTerm
		.toLowerCase()
		.split(' ')
		.filter((term) => term.trim() !== '');
	$: filteredLogs = searchTerm
		? logs.filter((log) => matchesAllSearchTerms(log, searchTerms))
		: logs;

	// Convert JSON string to object for display
	function parseJson(jsonString: string): any {
		try {
			return JSON.parse(jsonString);
		} catch (e) {
			return jsonString;
		}
	}

	// Format timestamp for display
	function formatDate(dateString: string | Date): string {
		try {
			const date = typeof dateString === 'string' ? new Date(dateString) : dateString;
			return date.toLocaleString();
		} catch (e) {
			return String(dateString);
		}
	}

	async function loadInitialLogs() {
		try {
			if (!$currentWorkspace) {
				throw new Error('No workspace selected');
			}
			
			isLoading = true;
			const result = await getLogs(1, pageSize, selectedProject.id);
			logs = result.logs;
			total = result.total;
			isLoading = false;
		} catch (err) {
			console.error('Failed to load logs:', err);
			error = 'Failed to load logs: ' + (err instanceof Error ? err.message : String(err));
			isLoading = false;
		}
	}

	function setupLogStream() {
		// Close any existing connection
		if (eventSource) {
			eventSource.close();
		}

		if (!$currentWorkspace) {
			console.error('Cannot setup log stream: No workspace selected');
			return;
		}

		console.log('Setting up log stream for project:', selectedProject.id, 'in workspace:', $currentWorkspace.id);

		// Create new connection
		eventSource = createLogStream(selectedProject.id, pageSize);

		// Setup event handlers
		eventSource.addEventListener('log', (event) => {
			try {
				console.log('Log event received:', event.data);
				const newLog = JSON.parse(event.data);

				// Check if log already exists to prevent duplicates
				if (!logs.some((log) => log.id === newLog.id)) {
					// Add to beginning of array (newest first) and force Svelte reactivity
					logs = [newLog, ...logs].slice(0, 1000); // Limit to 1000 logs to prevent browser slowdown
					console.log('Added new log, total logs:', logs.length);

					// Auto-scroll to top if enabled
					if (autoScroll) {
						window.scrollTo(0, 0);
					}
				}
			} catch (err) {
				console.error('Error processing log event:', err, event.data);
			}
		});

		// Direct message event (fallback)
		eventSource.onmessage = (event) => {
			console.log('Generic message received:', event.data);
			try {
				const newLog = JSON.parse(event.data);
				if (newLog && newLog.id && !logs.some((log) => log.id === newLog.id)) {
					logs = [newLog, ...logs].slice(0, 1000);
				}
			} catch (err) {
				console.error('Error processing generic message:', err);
			}
		};

		eventSource.addEventListener('ping', (event) => {
			// Keep connection alive, no action needed
			console.log('Ping received from server:', event.data);
			isConnected = true;
			reconnectAttempts = 0; // Reset reconnect counter on successful ping
		});

		eventSource.onopen = () => {
			console.log('Log stream connection established');
			isConnected = true;
			reconnectAttempts = 0; // Reset reconnect counter on connection
		};

		eventSource.onerror = (err) => {
			console.error('EventSource error:', err);
			isConnected = false;

			// Implement smart reconnection strategy with backoff
			if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
				reconnectAttempts++;
				const delay = RECONNECT_DELAY_MS * reconnectAttempts; // Increase delay with each attempt
				console.log(
					`Attempting to reconnect log stream (attempt ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS}) in ${delay}ms...`
				);

				setTimeout(() => {
					setupLogStream();
				}, delay);
			} else {
				console.error('Max reconnection attempts reached. Please refresh manually.');
			}
		};
	}

	// Function to create a mock from a log entry
	function createMockFromLog(log: RequestLog) {
		// Open the modal for creating a mock from the log
		selectedLogForMock = log;
		isCreateMockModalOpen = true;
		console.log('Opening mock creation modal for log:', log.id);
	}

	// Handle success after mock creation
	function handleMockCreationSuccess() {
		copyNotification = { show: true, message: 'Mock endpoint created successfully!' };
		setTimeout(() => {
			copyNotification = { show: false, message: '' };
		}, 2000);
	}

	// Track connection status for UI feedback
	let isConnected = false;
	let reconnectAttempts = 0;
	const MAX_RECONNECT_ATTEMPTS = 5;
	const RECONNECT_DELAY_MS = 3000;

	// Initialize on component mount
	onMount(() => {
		loadInitialLogs().then(() => {
			setupLogStream();
		});
	});

	// Clean up on component destroy
	onDestroy(() => {
		if (eventSource) {
			eventSource.close();
		}
	});
</script>

<div class="w-full theme-bg-primary p-4 relative">
	<!-- Copy notification toast -->
	{#if copyNotification.show}
		<div
			transition:fade={{ duration: 200 }}
			class="fixed top-6 right-6 theme-bg-secondary theme-text-primary px-4 py-2 rounded shadow-lg z-50 flex items-center"
		>
			<i class="fas fa-check-circle text-green-400 mr-2"></i>
			<span>{copyNotification.message}</span>
		</div>
	{/if}

	{#if !isConnected && reconnectAttempts > 0}
		<div
			class="bg-red-100/30 dark:bg-red-900/30 border border-red-300 dark:border-red-700 p-2 rounded mb-4 flex items-center justify-between"
		>
			<div class="flex items-center">
				<i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 text-lg mr-2"
				></i>
				<span class="theme-text-primary">Live stream disconnected. Using manual refresh only.</span>
			</div>
			<button
				class={ThemeUtils.primaryButton('py-1 px-3 text-sm')}
				on:click={() => setupLogStream()}
			>
				<i class="fas fa-sync mr-1"></i> Reconnect Stream
			</button>
		</div>
	{/if}
	<div class="mb-6">
		<div class="flex justify-between items-center mb-4">
			<div class="flex items-center">
				<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
					<i class="fas fa-list-alt text-blue-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">{selectedProject.name}</h2>
					<p class="text-sm theme-text-muted">Request logs</p>
				</div>
				<div
					class="ml-4 flex items-center bg-gray-100/50 dark:bg-gray-900/50 px-3 py-1 rounded-full"
				>
					<!-- Stream status indicator -->
					<span class="relative flex h-3 w-3 mr-2">
						{#if isConnected}
							<span
								class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"
							></span>
							<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
						{:else}
							<span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
						{/if}
					</span>
					<span class="text-xs font-medium {isConnected ? 'text-green-400' : 'text-red-400'}">
						{isConnected ? 'Live' : 'Offline'}
					</span>
				</div>
			</div>

			<div class="flex items-center space-x-3">
				<div class="flex items-center bg-gray-100/50 dark:bg-gray-900/50 px-3 py-1 rounded-full">
					<span class="text-xs theme-text-secondary mr-2">Auto-scroll</span>
					<label class="inline-flex items-center cursor-pointer">
						<input type="checkbox" bind:checked={autoScroll} class="sr-only peer" />
						<div
							class="relative w-9 h-5 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-500 rounded-full peer peer-focus:outline-none peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:rounded-full after:h-4 after:w-4 after:transition-all"
						></div>
					</label>
				</div>

				<button
					class={ThemeUtils.primaryButton('py-2 px-4 text-sm')}
					on:click={() => {
						loadInitialLogs();
						if (!isConnected) {
							setupLogStream(); // Also try to reconnect if disconnected
						}
					}}
				>
					<i class="fas fa-sync mr-2"></i> Refresh Logs
				</button>
			</div>
		</div>

		<div class="relative mb-6">
			<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
				<i class="fas fa-search theme-text-muted"></i>
			</div>
			<input
				type="text"
				bind:value={searchTerm}
				placeholder="Search by keywords separated by spaces (e.g. 'GET users')..."
				class={ThemeUtils.inputField('p-3 ps-10 text-sm rounded-lg')}
			/>
		</div>
	</div>

	{#if isLoading}
		<div class="flex justify-center py-8">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if error}
		<div class="bg-red-100 dark:bg-red-800 p-4 rounded mb-4 text-center">
			<p class="text-red-700 dark:text-white">{error}</p>
			<button on:click={loadInitialLogs} class={ThemeUtils.primaryButton('mt-2 py-1 px-4 text-sm')}>
				Retry
			</button>
		</div>
	{:else if filteredLogs.length === 0}
		<div class="theme-bg-secondary p-4 rounded text-center">
			<p class="theme-text-primary">
				No logs found {searchTerm ? 'matching your search criteria' : 'for this project'}
			</p>
		</div>
	{:else}
		<div class="space-y-4">
			{#each filteredLogs as log (log.id)}
				<div class={ThemeUtils.card('overflow-hidden')}>
					<!-- Log header - clickable to expand/collapse -->
					<div class="flex flex-col">
						<div
							class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer"
							on:click={() => toggleLogExpansion(log.id)}
							on:keydown={(e) => e.key === 'Enter' && toggleLogExpansion(log.id)}
							tabindex="0"
							role="button"
							aria-expanded={!!expandedLogs[log.id]}
						>
							<div class="flex items-center space-x-2">
								<!-- Method badge -->
								<span
									class="px-2 py-0.5 text-sm font-mono rounded {log.method === 'GET'
										? 'bg-green-600 text-white'
										: log.method === 'POST'
											? 'bg-blue-600 text-white'
											: log.method === 'PUT'
												? 'bg-yellow-600 text-white'
												: log.method === 'DELETE'
													? 'bg-red-600 text-white'
													: 'bg-gray-600 text-white'}"
								>
									{log.method}
								</span>

								<!-- Path with truncation -->
								<span class="font-mono text-sm theme-text-primary truncate max-w-sm">
									{log.path}
								</span>

								<!-- Status code -->
								<span
									class="px-2 py-0.5 text-xs font-mono rounded {log.response_status < 300
										? 'bg-green-600 text-white'
										: log.response_status < 400
											? 'bg-blue-600 text-white'
											: log.response_status < 500
												? 'bg-yellow-600 text-white'
												: 'bg-red-600 text-white'}"
								>
									{log.response_status}
								</span>

								<!-- Match status badge -->
								<span
									class="px-2 py-0.5 text-xs font-mono rounded {log.matched
										? 'bg-green-600 text-white'
										: 'bg-red-600 text-white'}"
								>
									{log.matched ? 'Matched' : 'Unmatched'}
								</span>
							</div>

							<div class="flex items-center space-x-3">
								<span class="text-xs theme-text-muted">{formatDate(log.created_at)}</span>
								<span class="px-2 py-0.5 text-xs bg-blue-600 rounded text-white"
									>{log.latency_ms}ms</span
								>
								<i
									class="fas {expandedLogs[log.id]
										? 'fa-chevron-up'
										: 'fa-chevron-down'} theme-text-muted"
								></i>
							</div>
						</div>

						<!-- Create Mock button row - only for unmatched requests -->
						{#if !log.matched}
							<div class="flex justify-end px-3 py-1 border-t theme-border-light">
								<button
									class="bg-emerald-600 hover:bg-emerald-700 text-white py-1 px-3 rounded text-xs flex items-center transition-all duration-200 transform hover:scale-105"
									on:click|stopPropagation={() => createMockFromLog(log)}
									title="Create a new mock endpoint using the data from this request"
								>
									<i class="fas fa-magic mr-1"></i> Create Mock from this Request
								</button>
							</div>
						{/if}
					</div>

					<!-- Expanded details -->
					{#if expandedLogs[log.id]}
						<div transition:fade={{ duration: 150 }} class="border-t theme-border px-4 py-3">
							<!-- Tab navigation -->
							<div class="flex mb-4 border-b theme-border">
								<button
									class="px-4 py-2 font-medium text-sm {activeTabs[log.id] === 'request'
										? 'text-blue-500 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400'
										: 'theme-text-muted hover:text-gray-600 dark:hover:text-gray-300'}"
									on:click|stopPropagation={() => switchTab(log.id, 'request')}
								>
									Request
								</button>
								<button
									class="px-4 py-2 font-medium text-sm {activeTabs[log.id] === 'response'
										? 'text-blue-500 dark:text-blue-400 border-b-2 border-blue-500 dark:border-blue-400'
										: 'theme-text-muted hover:text-gray-600 dark:hover:text-gray-300'}"
									on:click|stopPropagation={() => switchTab(log.id, 'response')}
								>
									Response
								</button>
							</div>

							<!-- Request content -->
							{#if activeTabs[log.id] === 'request'}
								<div>
									<!-- General info -->
									<div class="mb-4 bg-gray-100 dark:bg-gray-850 rounded-md p-3">
										<h3 class="text-sm font-semibold theme-text-secondary mb-2">General</h3>
										<div class="grid grid-cols-2 gap-2 text-sm">
											<div>
												<span class="theme-text-muted">Request URL:</span>
												<span class="theme-text-primary font-mono">{log.path}</span>
											</div>
											<div>
												<span class="theme-text-muted">Method:</span>
												<span class="theme-text-primary font-mono">{log.method}</span>
											</div>
										</div>
									</div>

									<!-- Headers with copy button -->
									<div class="mb-4">
										<div class="flex justify-between items-center mb-2">
											<h3 class="text-sm font-semibold theme-text-secondary">Headers</h3>
											<div class="flex space-x-2">
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.request_headers), null, 2),
															'Headers'
														)}
												>
													<i class="fas fa-copy mr-1"></i> Copy
												</button>
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.request_headers)),
															'Headers (minified)'
														)}
												>
													<i class="fas fa-compress-alt mr-1"></i> Minify
												</button>
											</div>
										</div>
										<pre
											class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-48">{JSON.stringify(
												parseJson(log.request_headers),
												null,
												2
											)}</pre>
									</div>

									<!-- Request body if exists -->
									{#if log.request_body}
										<div>
											<div class="flex justify-between items-center mb-2">
												<h3 class="text-sm font-semibold theme-text-secondary">Body</h3>
												<div class="flex space-x-2">
													<button
														class={ThemeUtils.utilityButton()}
														on:click|stopPropagation={() =>
															copyToClipboard(
																JSON.stringify(parseJson(log.request_body), null, 2),
																'Body'
															)}
													>
														<i class="fas fa-copy mr-1"></i> Copy
													</button>
													<button
														class={ThemeUtils.utilityButton()}
														on:click|stopPropagation={() =>
															copyToClipboard(
																JSON.stringify(parseJson(log.request_body)),
																'Body (minified)'
															)}
													>
														<i class="fas fa-compress-alt mr-1"></i> Minify
													</button>
												</div>
											</div>
											<pre
												class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-64">{JSON.stringify(
													parseJson(log.request_body),
													null,
													2
												)}</pre>
										</div>
									{/if}

									<!-- Query parameters if exists -->
									{#if log.query_params}
										<div class="mt-4">
											<div class="flex justify-between items-center mb-2">
												<h3 class="text-sm font-semibold theme-text-secondary">Query Parameters</h3>
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(log.query_params, 'Query parameters')}
												>
													<i class="fas fa-copy mr-1"></i> Copy
												</button>
											</div>
											<pre
												class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-32">{log.query_params}</pre>
										</div>
									{/if}
								</div>
							{:else}
								<!-- Response content -->
								<div>
									<!-- General info -->
									<div class="mb-4 bg-gray-100 dark:bg-gray-850 rounded-md p-3">
										<h3 class="text-sm font-semibold theme-text-secondary mb-2">General</h3>
										<div class="grid grid-cols-2 gap-2 text-sm">
											<div>
												<span class="theme-text-muted">Status Code:</span>
												<span
													class="{log.response_status < 300
														? 'text-green-600 dark:text-green-400'
														: log.response_status < 400
															? 'text-blue-600 dark:text-blue-400'
															: log.response_status < 500
																? 'text-yellow-600 dark:text-yellow-400'
																: 'text-red-600 dark:text-red-400'} font-mono"
												>
													{log.response_status}
												</span>
											</div>
											<div>
												<span class="theme-text-muted">Execution:</span>
												<span class="theme-text-primary font-mono">{log.execution_mode}</span>
											</div>
										</div>
									</div>

									<!-- Headers with copy button -->
									<div class="mb-4">
										<div class="flex justify-between items-center mb-2">
											<h3 class="text-sm font-semibold theme-text-secondary">Headers</h3>
											<div class="flex space-x-2">
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.response_headers), null, 2),
															'Headers'
														)}
												>
													<i class="fas fa-copy mr-1"></i> Copy
												</button>
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.response_headers)),
															'Headers (minified)'
														)}
												>
													<i class="fas fa-compress-alt mr-1"></i> Minify
												</button>
											</div>
										</div>
										<pre
											class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-48">{JSON.stringify(
												parseJson(log.response_headers),
												null,
												2
											)}</pre>
									</div>

									<!-- Response body -->
									<div>
										<div class="flex justify-between items-center mb-2">
											<h3 class="text-sm font-semibold text-gray-300">Body</h3>
											<div class="flex space-x-2">
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.response_body), null, 2),
															'Body'
														)}
												>
													<i class="fas fa-copy mr-1"></i> Copy
												</button>
												<button
													class={ThemeUtils.utilityButton()}
													on:click|stopPropagation={() =>
														copyToClipboard(
															JSON.stringify(parseJson(log.response_body)),
															'Body (minified)'
														)}
												>
													<i class="fas fa-compress-alt mr-1"></i> Minify
												</button>
											</div>
										</div>

										<!-- Special handling for endpoint not found error -->
										{#if log.response_status >= 400 && parseJson(log.response_body)?.error === true}
											<div
												class="bg-red-100/30 dark:bg-red-900/30 border border-red-300 dark:border-red-700 p-3 rounded-md"
											>
												<div class="flex items-center">
													<i
														class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-2"
													></i>
													<span class="text-sm theme-text-primary">
														{parseJson(log.response_body)?.message || 'Error'}
													</span>
												</div>
											</div>
										{:else}
											<pre
												class="bg-gray-300/50 dark:bg-gray-700 p-3 rounded-md text-xs theme-text-secondary font-mono overflow-auto max-h-64">{JSON.stringify(
													parseJson(log.response_body),
													null,
													2
												)}</pre>
										{/if}
									</div>
								</div>
							{/if}
						</div>
					{/if}
				</div>
			{/each}
		</div>

		{#if filteredLogs.length < total}
			<div class="mt-4 text-center">
				<span class="text-xs text-gray-400">Showing {filteredLogs.length} of {total} logs</span>
			</div>
		{/if}
	{/if}

	<!-- Create Mock Modal -->
	<ModalCreateMock
		isOpen={isCreateMockModalOpen}
		log={selectedLogForMock}
		projectId={selectedProject.id}
		onClose={() => (isCreateMockModalOpen = false)}
		onSuccess={handleMockCreationSuccess}
	/>
</div>
