<script lang="ts">
	import { onMount } from 'svelte';
	import { isAuthenticated } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { isFirstOpenPage } from '$lib/stores/isFirstOpen';
	onMount(async () => {
		console.log('onMount: page');
		if ($isFirstOpenPage) {
			isFirstOpenPage.set(false)
			if ($isAuthenticated) {
				await goto('/home');
			} else {
				await goto('/login');
			}
		}
	});
</script>
