<script lang="ts">
	import { fade } from 'svelte/transition';

	export let text: string = '';
	export let show: boolean = false;
	export let position: 'top' | 'bottom' | 'left' | 'right' = 'top';
</script>

{#if show && text}
	<div
		class="tooltip-wrapper absolute z-50"
		class:tooltip-top={position === 'top'}
		class:tooltip-bottom={position === 'bottom'}
		class:tooltip-left={position === 'left'}
		class:tooltip-right={position === 'right'}
		transition:fade={{ duration: 150 }}
	>
		<div
			class="bg-gray-900 dark:bg-gray-700 text-white text-xs font-medium px-3 py-2 rounded-lg shadow-lg border border-gray-700 dark:border-gray-600 whitespace-nowrap"
		>
			{text}
		</div>
		<!-- Arrow -->
		<div
			class="tooltip-arrow absolute"
			class:arrow-top={position === 'top'}
			class:arrow-bottom={position === 'bottom'}
			class:arrow-left={position === 'left'}
			class:arrow-right={position === 'right'}
		></div>
	</div>
{/if}

<style>
	.tooltip-wrapper {
		pointer-events: none;
	}

	.tooltip-top {
		bottom: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-bottom {
		top: calc(100% + 8px);
		left: 50%;
		transform: translateX(-50%);
	}

	.tooltip-left {
		right: calc(100% + 8px);
		top: 50%;
		transform: translateY(-50%);
	}

	.tooltip-right {
		left: calc(100% + 8px);
		top: 50%;
		transform: translateY(-50%);
	}

	.tooltip-arrow {
		width: 0;
		height: 0;
		border-style: solid;
	}

	.arrow-top {
		top: 100%;
		left: 50%;
		transform: translateX(-50%);
		border-width: 6px 6px 0 6px;
		border-color: #374151 transparent transparent transparent;
	}

	:global(.dark) .arrow-top {
		border-color: #4b5563 transparent transparent transparent;
	}

	.arrow-bottom {
		bottom: 100%;
		left: 50%;
		transform: translateX(-50%);
		border-width: 0 6px 6px 6px;
		border-color: transparent transparent #374151 transparent;
	}

	:global(.dark) .arrow-bottom {
		border-color: transparent transparent #4b5563 transparent;
	}

	.arrow-left {
		left: 100%;
		top: 50%;
		transform: translateY(-50%);
		border-width: 6px 0 6px 6px;
		border-color: transparent transparent transparent #374151;
	}

	:global(.dark) .arrow-left {
		border-color: transparent transparent transparent #4b5563;
	}

	.arrow-right {
		right: 100%;
		top: 50%;
		transform: translateY(-50%);
		border-width: 6px 6px 6px 0;
		border-color: transparent #374151 transparent transparent;
	}

	:global(.dark) .arrow-right {
		border-color: transparent #4b5563 transparent transparent;
	}
</style>
