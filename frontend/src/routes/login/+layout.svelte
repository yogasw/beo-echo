<script lang="ts">
	import '../../app.css';
	import { onMount } from 'svelte';

	onMount(() => {
		console.log('onMount: login');
	});
</script>

<slot />

<style>
    :global(body) {
        margin: 0;
        padding: 0;
        min-height: 100vh;
        background: linear-gradient(135deg, #1a1a1a 0%, #2d2d2d 100%);
    }
</style>
