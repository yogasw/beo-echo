<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { BodyTypeHttp } from '$lib/types/Replay';
	import MonacoEditor from '$lib/components/MonacoEditor.svelte';

	let { 
		payload = '',
		metadata = '',
		protocol = ''
	}: { 
		payload?: string;
		metadata?: string;
		protocol?: string;
	} = $props();
	
	const dispatch = createEventDispatcher();

	// Parse bodyType from raw metadata JSON string
	let localBodyType = $state<BodyTypeHttp>('none');
	$effect(() => {
		try {
			const meta = metadata ? JSON.parse(metadata) : {};
			console.log("meta",meta)
			localBodyType = meta.bodyType ?? 'none';
			console.debug('[ReplayBody] props:', { payload, metadata, protocol });
			console.debug('[ReplayBody] parsed meta:', meta, '→ bodyType:', localBodyType);
		} catch {
			localBodyType = 'none';
			console.debug('[ReplayBody] failed to parse metadata:', metadata);
		}
	});

	function handleBodyTypeChange(type: BodyTypeHttp) {
		localBodyType = type;
		dispatch('change', { bodyType: type });
	}

	function handleBodyContentChange(event: Event) {
		const target = event.target as HTMLTextAreaElement;
		dispatch('change', { payload: target.value });
	}

	let editorLanguage = $derived.by(() => {
		if (localBodyType === 'raw') {
			if (!payload || !payload.trim()) return 'json';
			const trimmed = payload.trim();
			if (trimmed.startsWith('{') || trimmed.startsWith('[')) return 'json';
			if (trimmed.startsWith('<')) return 'xml';
			return 'plaintext';
		}
		return 'plaintext';
	});

	function beautifyJson() {
		try {
			const parsed = JSON.parse(payload);
			dispatch('change', { payload: JSON.stringify(parsed, null, 2) });
		} catch {
			// Not JSON, leave as is
		}
	}
</script>

<!-- Body section -->
<div role="tabpanel" aria-labelledby="body-tab" class="space-y-4 flex flex-col h-full">
	<div class="flex items-center mb-4 shrink-0">
		<h2 class="text-sm font-semibold theme-text-primary flex items-center">
			<i class="fas fa-file-alt text-purple-500 mr-2"></i>
			Request Body
		</h2>
	</div>

	{#if protocol && protocol !== 'http'}
		<!-- Protocol not supported -->
		<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-8 flex flex-col items-center justify-center text-center gap-3">
			<i class="fas fa-ban text-3xl text-gray-400 dark:text-gray-600"></i>
			<p class="text-sm font-medium theme-text-secondary">
				Body not supported for <span class="font-semibold text-orange-500">{protocol}</span> protocol
			</p>
			<p class="text-xs theme-text-muted">Only HTTP / HTTPS requests support a request body.</p>
		</div>
	{:else}
		<div class="bg-white dark:bg-gray-800 border theme-border rounded-lg p-4 flex-1 flex flex-col min-h-0">
			<div class="flex flex-col flex-1 min-h-0 gap-4">
				<fieldset class="shrink-0">
					<legend class="sr-only">Body type selection</legend>
					<div class="flex flex-wrap gap-4 text-sm">
						{#each [
							{ value: 'none', label: 'None', title: 'No request body' },
							{ value: 'raw', label: 'Raw', title: 'Raw text or JSON body' },
							{ value: 'form-data', label: 'Form Data', title: 'Multipart form data' },
							{ value: 'x-www-form-urlencoded', label: 'x-www-form-urlencoded', title: 'URL encoded form body' }
						] as opt (opt.value)}
							<label class="flex items-center cursor-pointer">
								<input
									type="radio"
									name="bodyType"
									value={opt.value}
									class="mr-2 text-blue-500 focus:ring-blue-500"
									checked={localBodyType === opt.value}
									title={opt.title}
									aria-label={opt.label}
									onchange={() => handleBodyTypeChange(opt.value as BodyTypeHttp)}
								/>
								<span class="theme-text-secondary">{opt.label}</span>
							</label>
						{/each}
					</div>
				</fieldset>

				{#if localBodyType !== 'none'}
					<div class="flex flex-col flex-1 min-h-0">
						<label for="request-body" class="block text-sm font-medium theme-text-secondary mb-2 shrink-0">
							Body Content
						</label>
						<div class="flex-1 w-full rounded-md border theme-border overflow-hidden relative z-10 min-h-[300px]">
							<MonacoEditor
								value={payload}
								on:change={(e: CustomEvent<string>) => {
									dispatch('change', { payload: e.detail });
								}}
								language={editorLanguage}
							/>
						</div>
						<div class="flex justify-between items-center mt-2 shrink-0">
							<p class="text-xs theme-text-muted">
								<i class="fas fa-info-circle mr-1"></i>
								{#if localBodyType === 'raw'}
									Supports JSON, XML, text, and other formats
								{:else if localBodyType === 'form-data'}
									Multipart form data
								{:else}
									URL-encoded key=value pairs
								{/if}
							</p>
							{#if localBodyType === 'raw'}
								<button
									onclick={beautifyJson}
									class="text-xs text-blue-400 hover:text-blue-300 hover:underline transition-colors duration-200"
									title="Format and beautify JSON content"
									aria-label="Beautify JSON"
								>
									<i class="fas fa-magic mr-1"></i>
									Beautify
								</button>
							{/if}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}
</div>

