<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import type { Replay } from '$lib/types/Replay';

	let {
		isOpen = false,
		parsedResult = null,
		folders = [],
		initialFolderId = null,
		onclose,
		onImportWithoutSaving,
		onImportIntoCollection
	}: {
		isOpen?: boolean;
		parsedResult?: { parsed: Partial<Replay>; importType: string; displayName?: string; rawText: string } | null;
		folders?: any[];
		initialFolderId?: string | null;
		onclose?: () => void;
		onImportWithoutSaving?: (request: Partial<Replay>) => void;
		onImportIntoCollection?: (request: Partial<Replay>, folderId: string | null) => void;
	} = $props();

	let requestName = $state('');
	let selectedFolderId = $state<string | null>(null);

	// Sync when modal opens or result changes
	$effect(() => {
		if (isOpen && parsedResult?.parsed) {
			
			let suggestedName = parsedResult.parsed.name || '';
			
			// If it's the generic "Imported Request", try to be smarter
			if (suggestedName === 'Imported Request' && parsedResult.parsed.url) {
				try {
					const urlObj = new URL(parsedResult.parsed.url);
					suggestedName = urlObj.pathname.split('/').filter(Boolean).pop() || urlObj.host;
				} catch (e) {
					// fallback to just the url string if not a full valid URL
					const parts = parsedResult.parsed.url.split('/').filter(Boolean);
					if (parts.length > 0) suggestedName = parts[parts.length - 1];
				}
			}

			requestName = suggestedName || 'New Imported Request';
			
			// Set folder selection to initialFolderId if provided, otherwise root.
			selectedFolderId = initialFolderId || '';
		}
	});

	function handleKeydown(e: KeyboardEvent) {
		if (!isOpen) return;
		if (e.key === 'Escape') {
			if (onclose) onclose();
		}
	}

	function handleImportWithoutSaving() {
		if (!parsedResult || !onImportWithoutSaving) return;
		onImportWithoutSaving({
			...parsedResult.parsed,
			name: requestName
		});
	}

	function handleImportIntoCollection() {
		if (!parsedResult || !onImportIntoCollection) return;
		onImportIntoCollection(
			{
				...parsedResult.parsed,
				name: requestName
			},
			selectedFolderId === '' ? null : selectedFolderId
		);
	}

</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
	<div 
		class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/60 backdrop-blur-sm"
		transition:fade={{ duration: 150 }}
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div 
			class="absolute inset-0"
			onclick={onclose}
		></div>

		<div 
			class="relative w-full max-w-2xl theme-bg-primary theme-border border rounded-xl shadow-2xl flex flex-col"
			transition:scale={{ duration: 150, start: 0.98 }}
		>
			<!-- Header -->
			<div class="px-6 py-4 flex items-center justify-between border-b theme-border">
				<h2 class="text-xl font-semibold theme-text-primary">
					Import Replay from {parsedResult?.displayName ?? parsedResult?.importType} into a Folder
				</h2>
				<button 
					onclick={onclose}
					class="theme-text-muted hover:theme-text-primary transition-colors p-1 rounded-lg hover:bg-gray-500/10 active:scale-95"
					aria-label="Close"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>

			<!-- Body -->
			<div class="p-6 space-y-6 flex-1 overflow-y-auto max-h-[70vh]">
				
				<!-- Preview Block -->
				<div class="rounded-lg theme-bg-secondary border theme-border overflow-hidden">
					<div class="flex">
						<div class="py-3 px-3 text-right text-xs font-mono select-none text-gray-500 border-r border-[#ffffff10] theme-bg-primary/50">
							{#if parsedResult?.rawText}
								{#each (parsedResult?.rawText || '').split('\n') as _, i}
									<div>{i + 1}</div>
								{/each}
							{/if}
						</div>
						<textarea
							readonly
							class="w-full flex-1 p-3 bg-transparent text-xs font-mono theme-text-primary resize-y min-h-[120px] max-h-[300px] outline-none border-none whitespace-pre-wrap overflow-x-hidden"
							value={parsedResult?.rawText || ''}
							spellcheck="false"
						></textarea>
					</div>
				</div>

				<!-- Request Name -->
				<div class="space-y-2">
					<label for="requestName" class="block text-sm font-medium theme-text-primary">
						Request name
					</label>
					<input 
						id="requestName"
						type="text"
						bind:value={requestName}
						class="w-full bg-transparent border theme-border rounded-lg px-4 py-2 text-sm theme-text-primary focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-shadow"
						placeholder="E.g., Get User Data"
					/>
				</div>

				<!-- Folder Name -->
				<div class="space-y-2">
					<div class="flex items-center gap-2">
						<label for="collectionName" class="block text-sm font-medium theme-text-primary">
							Folder name
						</label>
						<i class="fas fa-info-circle text-[10px] theme-text-muted" title="Select the folder where this request should be saved."></i>
					</div>
					<div class="relative">
						<select 
							id="collectionName"
							bind:value={selectedFolderId}
							class="w-full bg-transparent border theme-border rounded-lg px-4 py-2 text-sm theme-text-primary appearance-none focus:outline-none focus:border-blue-500 focus:ring-1 focus:ring-blue-500 transition-shadow cursor-pointer"
						>
							<option value="" class="theme-bg-primary">Root (No Folder)</option>
							{#each folders as folder}
								<option value={folder.id} class="theme-bg-primary">{folder.name}</option>
							{/each}
						</select>
						<div class="absolute inset-y-0 right-0 flex items-center px-4 pointer-events-none theme-text-muted">
							<i class="fas fa-chevron-down text-xs"></i>
						</div>
					</div>
				</div>

			</div>

			<!-- Footer Actions -->
			<div class="px-6 py-4 flex items-center justify-end gap-3 border-t theme-border theme-bg-secondary/50 rounded-b-xl">
				<button 
					onclick={handleImportWithoutSaving}
					class="px-4 py-2 rounded-lg text-sm font-medium theme-text-primary theme-bg-secondary hover:bg-gray-500/20 transition-colors border theme-border"
				>
					Import Without Saving
				</button>
				<button 
					onclick={handleImportIntoCollection}
					class="px-4 py-2 rounded-lg text-sm font-medium text-white bg-orange-500 hover:bg-orange-600 transition-colors shadow-sm"
				>
					Import Into Folder
				</button>
			</div>
		</div>
	</div>
{/if}
