<script lang="ts">
	import { replayLoading } from '$lib/stores/replay';
	import type { ExecuteReplayResponse } from '$lib/types/Replay';
	import { createEventDispatcher } from 'svelte';

	const dispatch = createEventDispatcher();

	export let isExpanded = false; // Controls the visibility of the response body area - now accepting as prop
	export let executionResult: ExecuteReplayResponse | null = null; // Add execution result prop with any type

	// For response content tabs
	let activeSection = 'response'; // 'response', 'headers', 'cookies'

	// Function to expand the footer (can be called from parent)
	export function expand() {
		isExpanded = true;
		dispatch('toggleExpand', { expanded: isExpanded });
	}

	// Function to collapse the footer (can be called from parent)
	export function collapse() {
		isExpanded = false;
		dispatch('toggleExpand', { expanded: isExpanded });
	}

	function toggleExpand() {
		isExpanded = !isExpanded;
		console.log('Response body toggled:', isExpanded);
		dispatch('toggleExpand', { expanded: isExpanded });
	}

	function showHistory() {
		dispatch('showHistory');
	}

	function setActiveSection(section: string) {
		activeSection = section;
	}

	// Format response time in ms
	function formatResponseTime(timeMs: number): string {
		if (timeMs < 1000) {
			return `${timeMs}ms`;
		} else {
			return `${(timeMs / 1000).toFixed(2)}s`;
		}
	}

	// Get color class based on status code
	function getStatusColor(statusCode: number): string {
		if (statusCode >= 200 && statusCode < 300) {
			return 'bg-green-600 text-white';
		} else if (statusCode >= 300 && statusCode < 400) {
			return 'bg-blue-600 text-white';
		} else if (statusCode >= 400 && statusCode < 500) {
			return 'bg-yellow-600 text-white';
		} else if (statusCode >= 500) {
			return 'bg-red-600 text-white';
		} else {
			return 'bg-gray-600 text-white';
		}
	}

	// Format JSON for display
	function formatJson(json: string): string {
		try {
			return JSON.stringify(JSON.parse(json), null, 2);
		} catch (e) {
			return json;
		}
	}
</script>

<div class="bg-white dark:bg-gray-800 border-t border-gray-300 dark:border-gray-600 relative">
	<!-- Loading bar similar to Postman -->
	{#if $replayLoading.execute}
		<div class="absolute top-0 left-0 right-0 h-1 bg-gray-200 dark:bg-gray-700 overflow-hidden">
			<!-- Background progress bar -->
			<div class="h-full bg-blue-500/30 w-full"></div>
			<!-- Animated sliding indicator -->
			<div
				class="absolute top-0 h-full w-1/4 bg-gradient-to-r from-transparent via-blue-500 to-transparent animate-pulse"
				style="animation: slide 1.5s ease-in-out infinite; transform: translateX(-100%);"
			></div>
		</div>
		<style>
			@keyframes slide {
				0% {
					transform: translateX(-100%);
				}
				100% {
					transform: translateX(500%);
				}
			}
		</style>
	{/if}

	<div class="flex items-center justify-between p-3 border-b border-gray-200 dark:border-gray-700">
		<div class="flex items-center space-x-3">
			<span class="text-sm font-semibold text-gray-800 dark:text-white">Response</span>

			{#if executionResult && executionResult.status_code}
				<div class="flex items-center space-x-2">
					<span
						class={`px-2 py-0.5 rounded-full text-xs font-medium ${getStatusColor(executionResult.status_code)}`}
					>
						{executionResult.status_code}
					</span>

					{#if executionResult.status_text}
						<span class="text-gray-600 dark:text-gray-300 text-xs">
							{executionResult.status_text}
						</span>
					{/if}

					{#if executionResult.latency_ms}
						<span class="text-gray-600 dark:text-gray-300 text-xs">
							{formatResponseTime(executionResult.latency_ms)}
						</span>
					{/if}

					{#if executionResult.size}
						<span class="text-gray-600 dark:text-gray-300 text-xs">
							{executionResult.size} bytes
						</span>
					{/if}
				</div>
			{/if}

			<button
				on:click={showHistory}
				title="View request history"
				aria-label="View request history"
				class="flex items-center space-x-2 text-sm text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
			>
				<i class="fas fa-history text-sm"></i>
				<span>History</span>
			</button>
		</div>
		<button
			on:click={toggleExpand}
			aria-label="Toggle response body"
			title="Toggle response body"
			class="text-gray-600 dark:text-gray-300 hover:text-blue-600 dark:hover:text-blue-400 transition-colors duration-200"
		>
			<i class="fas {isExpanded ? 'fa-chevron-down' : 'fa-chevron-up'}"></i>
		</button>
	</div>
	{#if isExpanded}
		{#if executionResult}
			<div class="flex flex-col h-64 overflow-auto">
				<!-- Response tabs -->
				<div class="border-b border-gray-200 dark:border-gray-700">
					<div class="flex space-x-4 text-sm px-4" role="tablist" aria-label="Response tabs">
						<button
							class="py-2 px-1 border-b-2 {activeSection === 'response'
								? 'border-orange-600 text-orange-600'
								: 'border-transparent hover:text-gray-800 dark:hover:text-white'} transition-colors duration-200"
							title="View response body"
							aria-label="Response body tab"
							role="tab"
							aria-selected={activeSection === 'response'}
							on:click={() => setActiveSection('response')}
						>
							Body
						</button>
						<button
							class="py-2 px-1 border-b-2 {activeSection === 'headers'
								? 'border-orange-600 text-orange-600'
								: 'border-transparent hover:text-gray-800 dark:hover:text-white'} transition-colors duration-200"
							title="View response headers"
							aria-label="Response headers tab"
							role="tab"
							aria-selected={activeSection === 'headers'}
							on:click={() => setActiveSection('headers')}
						>
							Headers
						</button>
						<!-- <button
							class="py-2 px-1 border-b-2 {activeSection === 'cookies'
								? 'border-orange-600 text-orange-600'
								: 'border-transparent hover:text-gray-800 dark:hover:text-white'} transition-colors duration-200"
							title="View response cookies"
							aria-label="Response cookies tab"
							role="tab"
							aria-selected={activeSection === 'cookies'}
							on:click={() => setActiveSection('cookies')}
						>
							Cookies
						</button> -->
					</div>
				</div>

				<!-- Response content -->
				<div class="p-4 bg-gray-50 dark:bg-gray-900 flex-grow overflow-auto">
					{#if activeSection === 'response'}
						{#if executionResult.error}
							<div
								class="p-4 bg-red-100 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md text-red-800 dark:text-red-300 mb-4"
							>
								<h3 class="font-semibold">Error</h3>
								<p>{executionResult.error}</p>
							</div>
						{:else if executionResult.response_body}
							<pre
								class="bg-gray-800 text-gray-200 p-4 rounded-md overflow-auto font-mono text-sm max-h-52">
								{formatJson(executionResult.response_body)}
							</pre>
						{:else}
							<div
								class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300"
							>
								No response body
							</div>
						{/if}
					{:else if activeSection === 'headers'}
						{#if executionResult.headers && Object.keys(executionResult.headers).length > 0}
							<div class="border border-gray-200 dark:border-gray-700 rounded-md overflow-hidden">
								<table class="w-full">
									<thead class="bg-gray-100 dark:bg-gray-700">
										<tr>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Name</th
											>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Value</th
											>
										</tr>
									</thead>
									<tbody>
										{#each Object.entries(executionResult.headers || {}) as [name, value], i}
											<tr
												class={i % 2 === 0
													? 'bg-white dark:bg-gray-800'
													: 'bg-gray-50 dark:bg-gray-900'}
											>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm font-mono"
													>{name}</td
												>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm font-mono"
													>{value}</td
												>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{:else}
							<div
								class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300"
							>
								No headers received
							</div>
						{/if}
						<!-- this feature under development
						{:else if activeSection === 'cookies'}
						{#if executionResult.cookies && executionResult.cookies.length > 0}
							<div class="border border-gray-200 dark:border-gray-700 rounded-md overflow-hidden">
								<table class="w-full">
									<thead class="bg-gray-100 dark:bg-gray-700">
										<tr>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Name</th
											>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Value</th
											>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Domain</th
											>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Path</th
											>
											<th
												class="py-2 px-4 text-left text-gray-800 dark:text-white text-sm font-medium"
												>Expires</th
											>
										</tr>
									</thead>
									<tbody>
										{#each executionResult.cookies || [] as cookie, i}
											<tr
												class={i % 2 === 0
													? 'bg-white dark:bg-gray-800'
													: 'bg-gray-50 dark:bg-gray-900'}
											>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm"
													>{cookie.name}</td
												>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm font-mono"
													>{cookie.value}</td
												>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm"
													>{cookie.domain || '-'}</td
												>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm"
													>{cookie.path || '/'}</td
												>
												<td class="py-2 px-4 text-gray-600 dark:text-gray-300 text-sm"
													>{cookie.expires || '-'}</td
												>
											</tr>
										{/each}
									</tbody>
								</table>
							</div>
						{:else}
							<div
								class="p-4 bg-gray-100 dark:bg-gray-700 rounded-md text-center text-gray-600 dark:text-gray-300"
							>
								No cookies received
							</div>
						{/if} -->
					{/if}
				</div>
			</div>
		{:else}
			<div
				class="flex flex-col items-center justify-center h-64 text-center p-6 bg-gray-50 dark:bg-gray-900"
			>
				<!-- Placeholder for response body or actual response display -->
				<div
					class="h-24 w-24 mb-4 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center shadow-sm"
				>
					<i class="fas fa-rocket text-3xl text-gray-500 dark:text-gray-400"></i>
				</div>
				<p class="text-gray-600 dark:text-gray-300 max-w-md">
					Enter the URL and click Send to get a response. Response content will appear here.
				</p>
			</div>
		{/if}
	{/if}
</div>
