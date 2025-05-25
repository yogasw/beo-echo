<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { replayActions } from '$lib/stores/replay';
	import { toast } from '$lib/stores/toast';
	import { ReplayApi } from '$lib/api/replayApi';
	import { HTTP_METHODS, PROTOCOLS } from '$lib/types/Replay';
	import type { Replay, CreateReplayRequest } from '$lib/types/Replay';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';

	export let replay: Replay | null = null;

	const dispatch = createEventDispatcher();

	// Form data
	let formData: CreateReplayRequest = {
		name: '',
		protocol: 'https',
		method: 'GET',
		url: '',
		headers: {},
		body: ''
	};

	// Form state
	let isSubmitting = false;
	let showHeaders = false;
	let showBody = false;
	let headersText = '';
	let bodyText = '';

	// Initialize form data when replay prop changes
	$: if (replay) {
		formData = {
			name: replay.name,
			protocol: replay.protocol,
			method: replay.method,
			url: replay.url,
			headers: replay.headers || {},
			body: replay.body || ''
		};
		headersText = JSON.stringify(replay.headers || {}, null, 2);
		bodyText = replay.body || '';
		showHeaders = Object.keys(replay.headers || {}).length > 0;
		showBody = Boolean(replay.body);
	} else {
		// Reset form for new replay
		formData = {
			name: '',
			protocol: 'https',
			method: 'GET',
			url: '',
			headers: {},
			body: ''
		};
		headersText = '{}';
		bodyText = '';
		showHeaders = false;
		showBody = false;
	}

	// Methods that typically have a body
	$: methodsWithBody = ['POST', 'PUT', 'PATCH'].includes(formData.method);

	function toggleHeaders() {
		showHeaders = !showHeaders;
		if (!showHeaders) {
			headersText = '{}';
			formData.headers = {};
		}
	}

	function toggleBody() {
		showBody = !showBody;
		if (!showBody) {
			bodyText = '';
			formData.body = '';
		}
	}

	function handleHeadersChange(event: CustomEvent) {
		headersText = event.detail.value;
		try {
			formData.headers = JSON.parse(headersText);
		} catch (e) {
			// Invalid JSON, keep the text but don't update formData.headers
		}
	}

	function handleBodyChange(event: CustomEvent) {
		bodyText = event.detail.value;
		formData.body = bodyText;
	}

	async function handleSubmit() {
		if (!$selectedWorkspace || !$selectedProject) return;

		// Validation
		if (!formData.name.trim()) {
			toast.error('Alias is required');
			return;
		}

		if (!formData.url.trim()) {
			toast.error('target url is required');
			return;
		}

		// Validate headers JSON
		if (showHeaders && headersText.trim()) {
			try {
				formData.headers = JSON.parse(headersText);
			} catch (e) {
				toast.error('Invalid JSON in headers');
				return;
			}
		} else {
			formData.headers = {};
		}

		// Set body
		formData.body = showBody ? bodyText : '';

		try {
			isSubmitting = true;
			replayActions.setLoading('create', true);

			if (replay) {
				// Update existing replay (if API supports it)
				toast.info('Update functionality not yet implemented');
				dispatch('updated');
			} else {
				// Create new replay
				const response = await ReplayApi.createReplay(
					$selectedWorkspace.id,
					$selectedProject.id,
					formData
				);
				
				replayActions.addReplay(response.replay);
				toast.success('Replay created successfully');
				dispatch('created');
			}
		} catch (err: any) {
			toast.error(err);
		} finally {
			isSubmitting = false;
			replayActions.setLoading('create', false);
		}
	}

	function handleCancel() {
		dispatch('cancel');
	}

	function generateSampleUrl() {
		if (formData.protocol && !formData.url) {
			formData.url = `${formData.protocol}://api.example.com/endpoint`;
		}
	}

	function addCommonHeader(header: string, value: string) {
		try {
			const headers = showHeaders && headersText ? JSON.parse(headersText) : {};
			headers[header] = value;
			headersText = JSON.stringify(headers, null, 2);
			formData.headers = headers;
			showHeaders = true;
		} catch (e) {
			// If current headers are invalid, start fresh
			const headers = { [header]: value };
			headersText = JSON.stringify(headers, null, 2);
			formData.headers = headers;
			showHeaders = true;
		}
	}
</script>

<div class="flex flex-col h-full">
	<div class="flex-1 overflow-auto p-6">
		<form on:submit|preventDefault={handleSubmit} class="space-y-6 max-w-4xl">
			<!-- Basic Information -->
			<div class="theme-bg-secondary border theme-border rounded-lg p-4">
				<h3 class="text-lg font-medium theme-text-primary mb-4">
					<i class="fas fa-info-circle mr-2"></i>
					Basic Information
				</h3>

				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<!-- Alias -->
					<div>
						<label for="alias" class="block text-sm font-medium theme-text-primary mb-1">
							Alias <span class="text-red-400">*</span>
						</label>
						<input
							id="alias"
							type="text"
							bind:value={formData.name}
							placeholder="e.g., Get Users API"
							class="w-full px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>

					<!-- Protocol -->
					<div>
						<label for="protocol" class="block text-sm font-medium theme-text-primary mb-1">
							Protocol
						</label>
						<select
							id="protocol"
							bind:value={formData.protocol}
							on:change={generateSampleUrl}
							class="w-full px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							{#each PROTOCOLS as protocol}
								<option value={protocol}>{protocol.toUpperCase()}</option>
							{/each}
						</select>
					</div>

					<!-- Method -->
					<div>
						<label for="method" class="block text-sm font-medium theme-text-primary mb-1">
							HTTP Method
						</label>
						<select
							id="method"
							bind:value={formData.method}
							class="w-full px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary focus:outline-none focus:ring-2 focus:ring-blue-500"
						>
							{#each HTTP_METHODS as method}
								<option value={method}>{method}</option>
							{/each}
						</select>
					</div>

					<!-- URL -->
					<div>
						<label for="url" class="block text-sm font-medium theme-text-primary mb-1">
							URL <span class="text-red-400">*</span>
						</label>
						<input
							id="url"
							type="url"
							bind:value={formData.url}
							placeholder="https://api.example.com/endpoint"
							class="w-full px-3 py-2 border theme-border rounded-md theme-bg-primary theme-text-primary placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500"
							required
						/>
					</div>
				</div>
			</div>

			<!-- Headers Section -->
			<div class="theme-bg-secondary border theme-border rounded-lg p-4">
				<div class="flex items-center justify-between mb-4">
					<h3 class="text-lg font-medium theme-text-primary">
						<i class="fas fa-list-alt mr-2"></i>
						Headers
					</h3>
					<button
						type="button"
						on:click={toggleHeaders}
						class="px-3 py-1 text-sm {showHeaders ? 'bg-blue-600 text-white' : 'theme-bg-tertiary theme-text-secondary'} rounded transition-colors"
					>
						{showHeaders ? 'Hide' : 'Show'} Headers
					</button>
				</div>

				{#if showHeaders}
					<div class="space-y-3">
						<!-- Common Headers Shortcuts -->
						<div class="flex flex-wrap gap-2">
							<button
								type="button"
								on:click={() => addCommonHeader('Content-Type', 'application/json')}
								class="px-2 py-1 text-xs theme-bg-tertiary theme-text-secondary hover:theme-text-primary rounded"
							>
								+ JSON Content-Type
							</button>
							<button
								type="button"
								on:click={() => addCommonHeader('Authorization', 'Bearer TOKEN')}
								class="px-2 py-1 text-xs theme-bg-tertiary theme-text-secondary hover:theme-text-primary rounded"
							>
								+ Authorization
							</button>
							<button
								type="button"
								on:click={() => addCommonHeader('User-Agent', 'BeoEcho/1.0')}
								class="px-2 py-1 text-xs theme-bg-tertiary theme-text-secondary hover:theme-text-primary rounded"
							>
								+ User-Agent
							</button>
						</div>

						<!-- Headers Editor -->
						<div class="border theme-border rounded">
							<MonacoEditor
								value={headersText}
								language="json"
								height="150px"
								on:change={handleHeadersChange}
								options={{
									minimap: { enabled: false },
									scrollBeyondLastLine: false,
									fontSize: 13,
									wordWrap: 'on'
								}}
							/>
						</div>
					</div>
				{/if}
			</div>

			<!-- Body Section -->
			<div class="theme-bg-secondary border theme-border rounded-lg p-4">
				<div class="flex items-center justify-between mb-4">
					<h3 class="text-lg font-medium theme-text-primary">
						<i class="fas fa-file-code mr-2"></i>
						Request Body
						{#if methodsWithBody}
							<span class="text-sm theme-text-secondary">(recommended for {formData.method})</span>
						{/if}
					</h3>
					<button
						type="button"
						on:click={toggleBody}
						class="px-3 py-1 text-sm {showBody ? 'bg-blue-600 text-white' : 'theme-bg-tertiary theme-text-secondary'} rounded transition-colors"
					>
						{showBody ? 'Hide' : 'Show'} Body
					</button>
				</div>

				{#if showBody}
					<div class="border theme-border rounded">
						<MonacoEditor
							value={bodyText}
							language="json"
							height="200px"
							on:change={handleBodyChange}
							options={{
								minimap: { enabled: false },
								scrollBeyondLastLine: false,
								fontSize: 13,
								wordWrap: 'on'
							}}
						/>
					</div>
				{/if}
			</div>

			<!-- Form Actions -->
			<div class="flex items-center justify-end space-x-3 pt-4 border-t theme-border">
				<button
					type="button"
					on:click={handleCancel}
					class="px-4 py-2 text-sm theme-bg-tertiary theme-text-secondary hover:theme-text-primary border theme-border rounded transition-colors"
					disabled={isSubmitting}
				>
					Cancel
				</button>
				<button
					type="submit"
					class="px-4 py-2 text-sm bg-blue-600 hover:bg-blue-700 disabled:bg-blue-800 disabled:opacity-50 text-white rounded transition-colors flex items-center"
					disabled={isSubmitting}
				>
					{#if isSubmitting}
						<i class="fas fa-spinner fa-spin mr-2"></i>
					{:else}
						<i class="fas fa-save mr-2"></i>
					{/if}
					{replay ? 'Update' : 'Create'} Replay
				</button>
			</div>
		</form>
	</div>
</div>
