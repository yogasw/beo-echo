<script lang="ts">
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast';
	import { selectedWorkspace } from '$lib/stores/workspace';
	import WorkspaceDeleteConfirm from './WorkspaceDeleteConfirm.svelte';
	
	export let className = '';

	let showWorkspaceMenu = false;
	let showDeleteConfirm = false;
	
	function toggleWorkspaceMenu() {
		showWorkspaceMenu = !showWorkspaceMenu;
	}
	
	function openDeleteConfirm() {
		showWorkspaceMenu = false;
		showDeleteConfirm = true;
	}
	
	function closeDeleteConfirm() {
		showDeleteConfirm = false;
	}
	
	function handleClickOutside(event: MouseEvent) {
		const menu = document.getElementById('workspaceInfoMenu');
		const button = document.getElementById('workspaceInfoButton');
		if (menu && button && !menu.contains(event.target as Node) && !button.contains(event.target as Node)) {
			showWorkspaceMenu = false;
		}
	}
	
	onMount(() => {
		document.addEventListener('click', handleClickOutside);
		return () => {
			document.removeEventListener('click', handleClickOutside);
		};
	});
</script>

<div class={`relative group flex flex-col items-center ${className}`}>
	<button
		id="workspaceInfoButton"
		class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-green-500 flex items-center justify-center"
		on:click={toggleWorkspaceMenu}
		aria-label="Workspace Information"
	>
		<i class="fas fa-building"></i>
	</button>
	<span class="text-xs mt-1 theme-text-primary">Workspace</span>

	<!-- Workspace Info Menu -->
	{#if showWorkspaceMenu}
		<div
			id="workspaceInfoMenu"
			class="absolute top-full right-0 theme-bg-secondary theme-text-primary rounded shadow-lg mt-2 z-10 w-64"
		>
			<div class="p-3 border-b theme-border">
				<h3 class="font-semibold theme-text-primary">Current Workspace</h3>
				<p class="theme-text-secondary text-sm mt-1 truncate">
					{$selectedWorkspace?.name || 'No workspace selected'}
				</p>
			</div>
			
			{#if $selectedWorkspace}
				<button
					class="block w-full text-left px-4 py-2 text-red-500 hover:bg-gray-700 hover:text-red-400 transition-colors"
					on:click={openDeleteConfirm}
				>
					<i class="fas fa-trash-alt mr-2"></i> Delete Workspace
				</button>
			{:else}
				<div class="p-3 text-sm theme-text-secondary italic">
					Please select a workspace first
				</div>
			{/if}
		</div>
	{/if}
</div>

<WorkspaceDeleteConfirm 
	show={showDeleteConfirm} 
	workspace={$selectedWorkspace} 
	onClose={closeDeleteConfirm} 
/>
