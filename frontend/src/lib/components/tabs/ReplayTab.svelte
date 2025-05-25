<script lang="ts">
	import ReplayManager from '$lib/components/replay/ReplayManager.svelte';
	import type { Project } from '$lib/api/BeoApi';

	export let selectedProject: Project | null;
</script>

<div class="replay-tab h-full">
	<ReplayManager />
</div>

<style>
	.replay-tab {
		display: flex;
		flex-direction: column;
	}
</style>
