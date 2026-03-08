<script lang="ts">
	import { fade, scale } from 'svelte/transition';
	import { parseImportText } from '$lib/components/replay/parsers';
	import { toast } from '$lib/stores/toast';

	let {
		isOpen = false,
		onclose,
		onpreview
	}: {
		isOpen?: boolean;
		onclose?: () => void;
		onpreview?: (parsedResult: any) => void;
	} = $props();

	let importText = $state('');
	let isDragging = $state(false);

	function handleKeydown(e: KeyboardEvent) {
		if (!isOpen) return;
		if (e.key === 'Escape') {
			if (onclose) onclose();
		}
	}

	function processImport(text: string) {
		if (!text.trim()) return;
		
		try {
			const parsedResult = parseImportText(text);
			if (onpreview) {
				onpreview(parsedResult);
			}
			importText = '';
			if (onclose) onclose();
		} catch (err: any) {
			toast.error(err.message || "Failed to parse import block.");
		}
	}

	// Trigger processing on paste
	function handlePaste(e: ClipboardEvent) {
		// Wait a tick for value to bind, or grab from event
		setTimeout(() => {
			if (importText.trim()) {
				processImport(importText);
			}
		}, 10);
	}

	function handleEnter(e: KeyboardEvent) {
		if (e.key === 'Enter' && !e.shiftKey) {
			e.preventDefault();
			processImport(importText);
		}
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		
		// For now we just handle drops if they result in text/files.
		// The requirement was JSON or CURL. If a user drops a file, we could read it.
		// Let's implement active file reading if dropped:
		if (e.dataTransfer?.files?.length) {
			const file = e.dataTransfer.files[0];
			const reader = new FileReader();
			reader.onload = (event) => {
				const text = event.target?.result as string;
				if (text) {
					processImport(text);
				}
			};
			reader.readAsText(file);
		} else if (e.dataTransfer?.items && e.dataTransfer.items.length > 0) {
			const item = e.dataTransfer.items[0];
			if (item.kind === 'string') {
				item.getAsString((text) => {
					processImport(text);
				});
			}
		}
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
			class="relative w-full max-w-3xl theme-bg-primary theme-border border rounded-xl shadow-2xl overflow-hidden flex flex-col h-[600px] max-h-[90vh]"
			transition:scale={{ duration: 150, start: 0.98 }}
		>
			<!-- Header / Close -->
			<div class="px-4 py-3 flex items-center justify-between">
				<div class="text-sm theme-text-muted">
					Check out our <span class="text-blue-500 hover:underline">documentation</span> for more details.
				</div>
				<button 
					onclick={onclose}
					class="theme-text-muted hover:theme-text-primary transition-colors p-1 rounded-lg hover:bg-gray-500/10 active:scale-95"
					aria-label="Close"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>

			<!-- Input Section -->
			<div class="px-6 mb-2">
				<div class="relative flex items-center bg-transparent border-2 border-blue-600 rounded-lg overflow-hidden focus-within:ring-1 focus-within:ring-blue-600">
					<input 
						type="text"
						bind:value={importText}
						onpaste={handlePaste}
						onkeydown={handleEnter}
						placeholder="Paste cURL, Raw text or URL..."
						class="w-full bg-transparent px-4 py-3 theme-text-primary text-sm focus:outline-none"
					/>
				</div>
			</div>

			<!-- Tip Banner -->
			<div class="px-6 mb-4">
				<div class="flex items-center justify-between theme-bg-secondary px-4 py-2.5 rounded-lg theme-border border text-xs theme-text-muted">
					<div class="flex items-center gap-2">
						<i class="far fa-lightbulb"></i>
						<span>Tip : You can also paste cURL in the request bar to import</span>
					</div>
					<button class="hover:theme-text-primary transition-colors">Dismiss</button>
				</div>
			</div>

			<!-- Dropzone Body -->
			<!-- svelte-ignore a11y_no_static_element_interactions -->
			<div 
				class="flex-1 mx-6 mb-6 rounded-xl border-2 border-dashed transition-colors flex flex-col items-center justify-center text-center
					{isDragging ? 'border-blue-500 bg-blue-500/5' : 'theme-border hover:border-gray-500/50'}"
				ondragover={handleDragOver}
				ondragleave={handleDragLeave}
				ondrop={handleDrop}
			>
				<div class="w-16 h-16 rounded-full theme-bg-secondary theme-border border flex items-center justify-center mb-4">
					<i class="fas fa-file-import text-2xl theme-text-muted"></i>
				</div>
				<h3 class="text-xl font-medium theme-text-primary mb-1">
					Drop anywhere to import
				</h3>
				<p class="theme-text-muted text-sm">
					Or select <button class="text-blue-500 hover:underline">files</button> or <button class="text-blue-500 hover:underline">folders</button>
				</p>
			</div>

			<!-- Footer -->
			<div class="flex items-center justify-between px-6 py-4 theme-border border-t text-sm">
				<div class="flex items-center gap-6">
					<span class="flex items-center gap-2 theme-text-muted">
						<i class="fas fa-file-import"></i>
						Import from other formats
					</span>
				</div>
				<span class="theme-text-muted transition-colors flex items-center gap-1">
					Supported formats: cURL, JSON (Default Map)
				</span>
			</div>
		</div>
	</div>
{/if}
