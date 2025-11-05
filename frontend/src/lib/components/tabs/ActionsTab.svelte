<script lang="ts">
	import type { Project } from '$lib/types/beoConfig';
	import ActionList from '$lib/components/actions/ActionList.svelte';

	export let selectedProject: Project;
</script>

<div class="h-full">
	<ActionList />
</div>
