<script lang="ts">
	import { fade } from 'svelte/transition';
	
	// Component that creates a section wrapper with expandable content
	export let title: string;
	export let icon: string;
	export let iconBgColorClass: string = 'bg-blue-500/20';
	export let iconTextColorClass: string = 'text-blue-400';
	export let open: boolean = false;
	
	// Toggle the open state
	function toggleOpen() {
		open = !open;
	}
</script>

<div class="border theme-border rounded-lg overflow-hidden">
	<div 
		class="theme-bg-secondary p-4 flex items-center justify-between cursor-pointer"
		on:click={toggleOpen}
	>
		<div class="flex items-center gap-3">
			<div class="w-8 h-8 {iconBgColorClass} rounded-md flex items-center justify-center {iconTextColorClass}">
				<i class="fas {icon}"></i>
			</div>
			<h2 class="theme-text-primary font-medium">{title}</h2>
		</div>
		<div class="theme-text-secondary">
			<i class="fas fa-chevron-{open ? 'up' : 'down'}"></i>
		</div>
	</div>
	
	{#if open}
		<div transition:fade={{ duration: 200 }}>
			<slot></slot>
		</div>
	{/if}
</div>
