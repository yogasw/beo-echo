<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';

	/**
	 * Component for the "Send" button options dropdown.
	 * Handles execution source selection and URL host replacement.
	 */

	export let show = false;
	export let activeSource: 'server' | 'browser' | string = 'server';
	
	const dispatch = createEventDispatcher();
	
	let port = '';
	let showPortInput = false;
	let containerRef: HTMLDivElement;

	function handleSourceClick(source: 'server' | 'browser') {
		dispatch('selectSource', { source });
	}

	function togglePortInput() {
		showPortInput = !showPortInput;
		if (showPortInput) {
			// Auto focus input after it appears
			setTimeout(() => {
				const input = containerRef?.querySelector('input');
				if (input) input.focus();
			}, 0);
		}
	}

	function applyLocalhost() {
		dispatch('replaceHost', { port });
		showPortInput = false;
		dispatch('close');
	}

	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
            if (showPortInput) {
			    applyLocalhost();
            }
		} else if (e.key === 'Escape') {
			dispatch('close');
		}
	}

	// Close on outside click
	function handleClickOutside(e: MouseEvent) {
		if (show && containerRef && !containerRef.contains(e.target as Node)) {
            // Check if click was on the trigger button (handled by parent usually, but extra safety)
            const isTrigger = (e.target as HTMLElement).closest('[aria-label="Execute request"]');
            if (!isTrigger) {
			    dispatch('close');
            }
		}
	}

	onMount(() => {
		document.addEventListener('mousedown', handleClickOutside);
		return () => document.removeEventListener('mousedown', handleClickOutside);
	});
</script>

{#if show}
	<div 
		bind:this={containerRef}
		role="menu"
		tabindex="-1"
		class="absolute top-full right-0 mt-2 w-72 bg-white dark:bg-gray-800 rounded-md shadow-lg border theme-border z-50 flex flex-col py-1 overflow-hidden"
		on:keydown={handleKeydown}
	>
		<button 
			class="text-left px-4 py-3 text-sm hover:bg-blue-50 dark:hover:bg-blue-900/10 transition-colors flex items-center gap-3 {activeSource === 'server' ? 'text-blue-600 dark:text-blue-400 bg-blue-50/50 dark:bg-blue-900/10' : 'text-gray-700 dark:text-gray-300'}"
			on:click={() => handleSourceClick('server')}
		>
			<i class="fas fa-server w-5 text-center text-lg"></i>
			<div class="flex flex-col">
				<span class="font-medium">Execute from Beo Echo server</span>
				<span class="text-xs text-gray-400 dark:text-gray-500">Run from our backend infrastructure</span>
			</div>
		</button>
		
		<button 
			class="text-left px-4 py-3 text-sm hover:bg-green-50 dark:hover:bg-green-900/10 transition-colors flex items-center gap-3 {activeSource === 'browser' ? 'text-green-600 dark:text-green-400 bg-green-50/50 dark:bg-green-900/10' : 'text-gray-700 dark:text-gray-300'}"
			on:click={() => handleSourceClick('browser')}
		>
			<i class="fas fa-desktop w-5 text-center text-lg"></i>
			<div class="flex flex-col">
				<span class="font-medium">Execute from your local browser</span>
				<span class="text-xs text-gray-400 dark:text-gray-500">Run directly from this browser window</span>
			</div>
		</button>
		
		<div class="h-px bg-gray-200 dark:bg-gray-700 my-1 mx-2"></div>
		
		{#if !showPortInput}
			<button 
				class="text-left px-4 py-3 text-sm text-orange-600 dark:text-orange-400 hover:bg-orange-100/50 dark:hover:bg-orange-900/30 transition-colors flex items-center gap-3"
				on:click|stopPropagation={togglePortInput}
			>
				<i class="fas fa-laptop-code w-5 text-center text-lg"></i>
				<div class="flex flex-col">
					<span class="font-medium">Replace host to localhost</span>
					<span class="text-xs text-gray-400 dark:text-gray-500">Quickly point to local dev server</span>
				</div>
			</button>
		{:else}
			<div class="px-4 py-3 flex flex-col gap-2 bg-orange-50/30 dark:bg-orange-950/20">
				<div class="flex items-center justify-between">
					<span class="text-sm font-semibold text-orange-600 dark:text-orange-400">Localhost Port</span>
					<button 
						class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200 p-1"
						aria-label="Close port input"
						on:click={() => showPortInput = false}
					>
						<i class="fas fa-times text-xs"></i>
					</button>
				</div>
				<div class="flex items-center gap-2">
					<div class="relative flex-1">
						<span class="absolute left-2 top-1/2 -translate-y-1/2 text-gray-400 text-xs">:</span>
						<input 
							type="number" 
							bind:value={port}
							placeholder="Port (optional)"
							class="w-full pl-5 pr-2 py-2 text-sm rounded border theme-border bg-white dark:bg-gray-900 focus:outline-none focus:ring-1 focus:ring-orange-500"
						/>
					</div>
					<button 
						class="px-4 py-2 bg-orange-600 hover:bg-orange-700 text-white text-sm font-bold rounded shadow-sm transition-all"
						on:click={applyLocalhost}
					>
						Replace
					</button>
				</div>
			</div>
		{/if}
	</div>
{/if}

<style>
	/* Hide arrows/spinners in number input */
	input::-webkit-outer-spin-button,
	input::-webkit-inner-spin-button {
		-webkit-appearance: none;
		margin: 0;
	}
	input[type=number] {
		-moz-appearance: textfield;
		appearance: textfield;
	}
</style>
