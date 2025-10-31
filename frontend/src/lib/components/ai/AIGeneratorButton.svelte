<script lang="ts">
	import AIGeneratorPanel from './AIGeneratorPanel.svelte';

	export let initialContent: string = '{}';
	export let buttonText: string = 'AI Generate';
	export let buttonClass: string = 'bg-purple-600 hover:bg-purple-700';
	export let iconClass: string = 'fas fa-magic';
	export let size: 'sm' | 'md' | 'lg' = 'md';
	export let onGenerated: (content: string) => void = () => {};
	export let contentType: string = '';

	let isPanelOpen = false;

	const sizeClasses = {
		sm: 'px-2 py-1 text-xs',
		md: 'px-3 py-2 text-sm',
		lg: 'px-4 py-2 text-base'
	};

	function openPanel() {
		isPanelOpen = true;
	}

	function handleSave(content: string) {
		onGenerated(content);
	}

	function handleClose() {
		isPanelOpen = false;
	}
</script>

<button
	on:click={openPanel}
	class="{buttonClass} {sizeClasses[
		size
	]} text-white rounded-lg transition-colors flex items-center space-x-2 font-medium shadow-sm"
	aria-label="Open AI Generator"
	title="Generate content with AI"
>
	<i class={iconClass}></i>
	<span>{buttonText}</span>
</button>

<AIGeneratorPanel
	bind:isOpen={isPanelOpen}
	initialContent={initialContent}
	onSave={handleSave}
	onClose={handleClose}
	contentType={contentType}
/>
