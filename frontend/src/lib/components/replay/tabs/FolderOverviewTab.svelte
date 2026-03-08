<script lang="ts">
	import InkMde from 'ink-mde/svelte';
	import { marked } from 'marked';
	// @ts-expect-error isomorphic-dompurify typing issue
	import DOMPurify from 'isomorphic-dompurify';
	import { theme } from '$lib/stores/theme';
	import { replayApi } from '$lib/api/replayApi';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';

	let { folder }: { folder: any } = $props();

	async function saveDocumentation(doc: string) {
		if (!$selectedWorkspace || !$selectedProject || !folder?.id) return;
		
		try {
			await replayApi.updateFolder($selectedWorkspace.id, $selectedProject.id, folder.id, {
				doc: doc
			});
			folder.doc = doc;
		} catch (err: any) {
			toast.error(err.message || 'Failed to save documentation');
		}
	}

	let documentation = $state('');
	let isEditing = $state(false);

	$effect(() => {
		if (folder?.id) {
			documentation = folder.doc || '';
			isEditing = false;
		}
	});

	let saveTimer: ReturnType<typeof setTimeout> | null = null;
	
	// Track changes from InkMde
	$effect(() => {
		if (!isEditing || !folder?.id) return;
		const current = documentation;
		if (saveTimer) clearTimeout(saveTimer);
		saveTimer = setTimeout(() => saveDocumentation(current), 800);
	});

	function exitEditMode() {
		if (saveTimer) clearTimeout(saveTimer);
		saveDocumentation(documentation);
		isEditing = false;
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Escape' || (e.key === 'Enter' && (e.ctrlKey || e.metaKey))) {
			exitEditMode();
		}
	}

	// Render markdown for preview
	let renderedDocumentation = $derived.by(() => {
		if (!documentation.trim()) return '';
		const html = marked.parse(documentation, { async: false }) as string;
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
				bind:value={documentation}
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
				{#if documentation.trim()}
					<!-- eslint-disable-next-line svelte/no-at-html-tags -->
					<div class="ink-viewer-wrap prose {$theme === 'dark' ? 'prose-invert' : ''} prose-sm max-w-none px-2 py-1">{@html renderedDocumentation}</div>
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