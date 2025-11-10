<script lang="ts">
	import type { RequestLog } from '$lib/api/BeoApi';
	import { logs } from '$lib/stores/logs';

	export let visible: boolean = false;
	export let script: string = '';
	export let executionPoint: 'before_request' | 'after_request' = 'after_request';

	let selectedLog: RequestLog | null = null;
	let modifiedContext: any = null;
	let executing: boolean = false;

	// Get logs from store
	$: availableLogs = $logs.slice(0, 20);

	function selectLog(log: RequestLog) {
		selectedLog = log;
		modifiedContext = null;
	}

	function runScript() {
		if (!selectedLog) return;

		executing = true;
		modifiedContext = null;

		console.log('%c▶ Executing JavaScript Action...', 'color: #4ec9b0; font-size: 14px; font-weight: bold;');

		try {
			// Parse log data
			const request = {
				method: selectedLog.method,
				path: selectedLog.path,
				query: selectedLog.query_params ? JSON.parse(selectedLog.query_params) : {},
				headers: selectedLog.request_headers ? JSON.parse(selectedLog.request_headers) : {},
				body: selectedLog.request_body || ''
			};

			const response = executionPoint === 'after_request' ? {
				status_code: selectedLog.response_status,
				headers: selectedLog.response_headers ? JSON.parse(selectedLog.response_headers) : {},
				body: selectedLog.response_body || ''
			} : undefined;

			console.log('%cContext:', 'color: #569cd6; font-weight: bold;');
			console.log('  • request:', request);
			if (response) {
				console.log('  • response:', response);
			}
			console.log('');

			// Execute script in real Chrome context
			// All console.log/error/warn will appear in Chrome DevTools Console
			// debugger; statement will work
			// Breakpoints can be set in DevTools Sources tab
			const scriptFunction = new Function('request', 'response', script);
			scriptFunction(request, response);

			// Store modified context
			modifiedContext = {
				request,
				...(response && { response })
			};

			console.log('');
			console.log('%c✅ Script completed successfully', 'color: #107c10; font-weight: bold;');
			console.log('');
			console.log('%cModified Context:', 'color: #569cd6; font-weight: bold;');
			console.log('  • request:', request);
			if (response) {
				console.log('  • response:', response);
			}

		} catch (error: any) {
			console.log('');
			console.log('%c━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━', 'color: #3e3e42;');
			console.error('%c❌ Script execution failed', 'color: #f48771; font-weight: bold;');
			console.error(error);
		} finally {
			executing = false;
		}
	}
</script>

{#if visible}
	<div class="fixed inset-0 bg-black/50 backdrop-blur-sm z-50 flex items-center justify-center p-4">
		<div class="bg-white dark:bg-gray-800 rounded-lg shadow-2xl max-w-6xl w-full max-h-[90vh] overflow-hidden flex flex-col">
			<!-- Header -->
			<div class="flex items-center justify-between p-4 border-b theme-border">
				<div class="flex items-center gap-2">
					<i class="fas fa-play-circle text-green-500"></i>
					<h2 class="text-lg font-semibold theme-text-primary">Test JavaScript Action</h2>
				</div>
				<button
					type="button"
					on:click={() => (visible = false)}
					class="text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
					aria-label="Close test panel"
				>
					<i class="fas fa-times text-xl"></i>
				</button>
			</div>

			<!-- Content -->
			<div class="flex-1 overflow-hidden flex">
				<!-- Left Panel - Log Selection -->
				<div class="w-1/3 border-r theme-border overflow-y-auto p-4">
					<h3 class="text-sm font-semibold theme-text-primary mb-3 sticky top-0 bg-white dark:bg-gray-800 pb-2">
						Select Sample Log
					</h3>

					{#if availableLogs.length === 0}
						<div class="text-center py-8">
							<i class="fas fa-inbox text-4xl theme-text-secondary mb-2"></i>
							<p class="theme-text-secondary text-sm">No logs available</p>
							<p class="text-xs theme-text-secondary mt-1">Make requests to generate logs</p>
						</div>
					{:else}
						<div class="space-y-2">
							{#each availableLogs as log}
								<button
									type="button"
									on:click={() => selectLog(log)}
									class="w-full text-left p-3 rounded border transition-all {selectedLog?.id === log.id
										? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20 ring-1 ring-blue-500'
										: 'border-gray-200 dark:border-gray-700 hover:border-gray-300 dark:hover:border-gray-600'}"
								>
									<div class="flex items-start justify-between gap-2">
										<div class="flex-1 min-w-0">
											<div class="flex items-center gap-2 mb-1">
												<span class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium {
													log.method === 'GET' ? 'bg-blue-100 text-blue-800 dark:bg-blue-900/50 dark:text-blue-300' :
													log.method === 'POST' ? 'bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300' :
													log.method === 'PUT' ? 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/50 dark:text-yellow-300' :
													log.method === 'DELETE' ? 'bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300' :
													'bg-gray-100 text-gray-800 dark:bg-gray-900/50 dark:text-gray-300'
												}">
													{log.method}
												</span>
												<code class="text-xs theme-text-primary font-mono truncate">{log.path}</code>
											</div>
											<div class="flex items-center gap-2 text-xs theme-text-secondary">
												<span class="flex items-center gap-1">
													<i class="fas fa-circle text-[8px] {log.response_status >= 200 && log.response_status < 300 ? 'text-green-500' : log.response_status >= 400 ? 'text-red-500' : 'text-yellow-500'}"></i>
													{log.response_status}
												</span>
												<span>{log.latency_ms}ms</span>
											</div>
										</div>
										{#if selectedLog?.id === log.id}
											<i class="fas fa-check-circle text-blue-500 flex-shrink-0 text-sm"></i>
										{/if}
									</div>
								</button>
							{/each}
						</div>
					{/if}
				</div>

				<!-- Right Panel - Execution & Results -->
				<div class="flex-1 overflow-y-auto p-4">
					{#if !selectedLog}
						<div class="flex items-center justify-center h-full">
							<div class="text-center theme-text-secondary">
								<i class="fas fa-hand-pointer text-5xl mb-3"></i>
								<p class="font-medium">Select a log to test your script</p>
								<p class="text-sm mt-1">Choose a sample request/response from the left panel</p>
							</div>
						</div>
					{:else}
						<div class="space-y-4">
							<!-- Run Button -->
							<div class="flex items-center gap-2">
								<button
									type="button"
									on:click={runScript}
									disabled={executing}
									class="px-4 py-2 rounded bg-green-600 hover:bg-green-700 text-white transition-colors disabled:bg-gray-400 disabled:cursor-not-allowed flex items-center gap-2"
								>
									{#if executing}
										<i class="fas fa-spinner fa-spin"></i>
										Executing...
									{:else}
										<i class="fas fa-play"></i>
										Run Script
									{/if}
								</button>
								<div class="text-xs theme-text-secondary">
									<i class="fas fa-arrow-right mr-1"></i>
									Check Chrome DevTools Console for output
								</div>
							</div>

							<!-- Request Data -->
							<div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4 border theme-border">
								<div class="flex items-center gap-2 mb-2">
									<i class="fas fa-arrow-up text-blue-500"></i>
									<h4 class="text-sm font-semibold theme-text-primary">Request Data</h4>
								</div>
								<pre class="text-xs overflow-x-auto p-3 bg-white dark:bg-gray-800 rounded border theme-border"><code>{JSON.stringify({
	method: selectedLog.method,
	path: selectedLog.path,
	query: selectedLog.query_params ? JSON.parse(selectedLog.query_params) : {},
	headers: selectedLog.request_headers ? JSON.parse(selectedLog.request_headers) : {},
	body: selectedLog.request_body || ''
}, null, 2)}</code></pre>
							</div>

							<!-- Response Data -->
							{#if executionPoint === 'after_request'}
								<div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4 border theme-border">
									<div class="flex items-center gap-2 mb-2">
										<i class="fas fa-arrow-down text-green-500"></i>
										<h4 class="text-sm font-semibold theme-text-primary">Response Data</h4>
									</div>
									<pre class="text-xs overflow-x-auto p-3 bg-white dark:bg-gray-800 rounded border theme-border"><code>{JSON.stringify({
	status_code: selectedLog.response_status,
	headers: selectedLog.response_headers ? JSON.parse(selectedLog.response_headers) : {},
	body: selectedLog.response_body || ''
}, null, 2)}</code></pre>
								</div>
							{/if}

							<!-- Modified Context -->
							{#if modifiedContext}
								<div class="bg-gray-50 dark:bg-gray-900 rounded-lg p-4 border theme-border">
									<div class="flex items-center gap-2 mb-2">
										<i class="fas fa-edit text-orange-500"></i>
										<h4 class="text-sm font-semibold theme-text-primary">Modified Context (After Execution)</h4>
									</div>
									<pre class="text-xs overflow-x-auto p-3 bg-white dark:bg-gray-800 rounded border theme-border"><code>{JSON.stringify(modifiedContext, null, 2)}</code></pre>
								</div>
							{/if}
						</div>
					{/if}
				</div>
			</div>

			<!-- Footer -->
			<div class="flex items-center justify-between gap-3 p-4 border-t theme-border bg-gray-50 dark:bg-gray-900">
				<div class="text-xs theme-text-secondary">
					<i class="fas fa-info-circle mr-1"></i>
					{availableLogs.length} log{availableLogs.length !== 1 ? 's' : ''} available for testing
				</div>
				<button
					type="button"
					on:click={() => (visible = false)}
					class="px-4 py-2 text-sm rounded border theme-border hover:bg-gray-100 dark:hover:bg-gray-800 transition-colors"
				>
					Close
				</button>
			</div>
		</div>
	</div>
{/if}
