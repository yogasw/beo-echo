<script lang="ts">
	import { fade, scale } from 'svelte/transition';

	let {
		isOpen = false,
		title = 'Confirm Action',
		message = 'Are you sure you want to proceed?',
		confirmText = 'Confirm',
		cancelText = 'Cancel',
		confirmColor = 'bg-red-600 hover:bg-red-700 text-white border border-red-600 hover:border-red-700',
		onconfirm,
		oncancel
	}: {
		isOpen?: boolean;
		title?: string;
		message?: string;
		confirmText?: string;
		cancelText?: string;
		confirmColor?: string;
		onconfirm?: () => void;
		oncancel?: () => void;
	} = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (!isOpen) return;
		if (e.key === 'Escape') {
			if (oncancel) oncancel();
		} else if (e.key === 'Enter') {
			if (onconfirm) onconfirm();
		}
	}
</script>

<svelte:window onkeydown={handleKeydown} />

{#if isOpen}
	<div 
		class="fixed inset-0 z-[200] flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
		transition:fade={{ duration: 150 }}
	>
		<!-- svelte-ignore a11y_click_events_have_key_events -->
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<div 
			class="absolute inset-0"
			onclick={oncancel}
		></div>

		<div 
			class="relative w-full max-w-sm theme-bg-primary theme-border border rounded-xl shadow-2xl overflow-hidden flex flex-col"
			transition:scale={{ duration: 150, start: 0.95 }}
		>
			<!-- Header -->
			<div class="px-5 py-4 theme-border border-b flex items-center justify-between theme-bg-secondary">
				<h3 class="text-lg font-semibold theme-text-primary tracking-tight">
					{title}
				</h3>
				<button 
					onclick={oncancel}
					class="theme-text-muted hover:theme-text-primary transition-colors p-1.5 rounded-lg hover:bg-gray-500/10 active:scale-95"
					aria-label="Close"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>

			<!-- Body -->
			<div class="p-5">
				<p class="theme-text-secondary text-sm leading-relaxed">
					{message}
				</p>
			</div>

			<!-- Footer -->
			<div class="px-5 py-4 theme-bg-secondary flex justify-end gap-3 rounded-b-xl border-t theme-border">
				<button 
					onclick={oncancel}
					class="px-4 py-2 text-sm font-medium theme-text-primary bg-transparent theme-border border rounded-md hover:theme-bg-hover transition-colors focus:ring-2 focus:ring-gray-300 dark:focus:ring-gray-600 outline-none"
				>
					{cancelText}
				</button>
				<button 
					onclick={onconfirm}
					class="px-4 py-2 text-sm font-medium rounded-md shadow-sm transition-colors focus:ring-2 focus:ring-offset-2 focus:ring-red-500 outline-none {confirmColor}"
				>
					{confirmText}
				</button>
			</div>
		</div>
	</div>
{/if}
