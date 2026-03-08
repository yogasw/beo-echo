<script lang="ts">
	import InkMde from 'ink-mde/svelte';
	import { marked } from 'marked';
	import DOMPurify from 'isomorphic-dompurify';
	import { theme } from '$lib/stores/theme';

	const STORAGE_KEY = (id: string) => `beo_replay_folder_${id}_doc`;

	let { folder }: { folder: any } = $props();

	function loadDescription(folderId: string): string {
		try {
			const raw = localStorage.getItem(STORAGE_KEY(folderId));
			if (raw) return JSON.parse(raw).description || '';
		} catch {}
		return '';
	}

	function saveDescription(folderId: string, desc: string) {
		try {
			localStorage.setItem(
				STORAGE_KEY(folderId),
				JSON.stringify({ description: desc, updatedAt: new Date().toISOString() })
			);
		} catch {}
	}

	let description = $state('');
	let isEditing = $state(false);

	$effect(() => {
		if (folder?.id) {
			description = loadDescription(folder.id);
			isEditing = false;
		}
	});

	let saveTimer: ReturnType<typeof setTimeout> | null = null;
	
	// Track changes from InkMde
	$effect(() => {
		if (!isEditing || !folder?.id) return;
		const current = description;
		if (saveTimer) clearTimeout(saveTimer);
		saveTimer = setTimeout(() => saveDescription(folder.id, current), 400);
	});

	function exitEditMode() {
		if (saveTimer) clearTimeout(saveTimer);
		saveDescription(folder.id, description);
		isEditing = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' || (e.key === 'Enter' && (e.ctrlKey || e.metaKey))) {
			exitEditMode();
		}
	}

	// Render markdown for preview
	let renderedDescription = $derived.by(() => {
		if (!description.trim()) return '';
		const html = marked.parse(description, { async: false }) as string;
		return DOMPurify.sanitize(html);
	});

</script>

<div class="ink-wrapper h-full flex flex-col" onkeydown={handleKeydown} role="presentation">
	{#if isEditing}
		<!-- Edit mode -->
		<div class="flex items-center justify-between mb-2 flex-shrink-0">
			<span class="text-xs theme-text-muted flex items-center gap-1.5">
				<i class="fas fa-edit text-blue-400"></i>
				Editing —
				<kbd class="theme-bg-secondary theme-border border rounded px-1.5 py-0.5 text-[10px] font-mono theme-text-primary">Esc</kbd>
				to preview
			</span>
			<button
				onclick={exitEditMode}
				class="text-xs px-2.5 py-1 theme-bg-secondary hover:theme-bg-hover theme-text-primary rounded transition-colors flex items-center gap-1.5 theme-border border shadow-sm"
				title="Switch to preview"
				aria-label="Switch to preview mode"
			>
				<i class="fas fa-eye text-[10px]"></i>
				Preview
			</button>
		</div>

		<div class="flex-1 min-h-0 ink-editor-wrap rounded-lg overflow-hidden theme-border border theme-bg-primary">
			<InkMde
				bind:value={description}
				options={{
					interface: {
						appearance: $theme === 'dark' ? 'dark' : 'light',
						spellcheck: false,
						readonly: false,
					},
				}}
			/>
		</div>
	{:else}
		<!-- Preview mode — double-click to edit -->
		<div
			class="flex-1 flex flex-col min-h-0 cursor-text group"
			role="region"
			aria-label="Folder documentation — double-click to edit"
			ondblclick={() => (isEditing = true)}
		>
			<div class="flex-1 overflow-auto min-h-0 bg-transparent rounded-lg">
				{#if description.trim()}
					<!-- eslint-disable-next-line svelte/no-at-html-tags -->
					<div class="ink-viewer-wrap prose {$theme === 'dark' ? 'prose-invert' : ''} prose-sm max-w-none px-2 py-1">{@html renderedDescription}</div>
				{:else}
					<div class="flex flex-col items-center justify-center h-full text-center py-12 rounded-lg border-2 border-dashed theme-border theme-text-muted transition-colors">
						<i class="fas fa-file-alt text-3xl mb-3 opacity-50"></i>
						<p class="text-sm mb-1 theme-text-primary">No documentation yet</p>
						<p class="text-xs">Double-click to start writing</p>
					</div>
				{/if}
			</div>

			<p class="mt-3 flex-shrink-0 text-[11px] theme-text-muted transition-colors flex items-center gap-1.5 px-1 font-medium">
				<i class="fas fa-mouse-pointer"></i>
				Double-click to edit · Markdown supported by ink-mde
			</p>
		</div>
	{/if}
</div>

<style>
	/* ── Ink MDE attribution hide ── */
	:global(.ink-mde-details) {
		display: none !important;
	}
	
/* Hide internal frame styling */
:global(.ink-mde) {
--ink-internal-block-background-color: transparent !important;
--ink-internal-border-radius: 0 !important;
}

/* Hide editor focus outline */
:global(.ink-mde .cm-editor.cm-focused) {
	outline: none !important;
}
</style>