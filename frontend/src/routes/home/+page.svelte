<script lang="ts">
	import ContentArea from '$lib/components/ContentArea.svelte';
	import { onMount } from 'svelte';
	import { getProjects } from '$lib/api/BeoApi';
	import BeoEchoLoader from '$lib/components/common/BeoEchoLoader.svelte';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { projects } from '$lib/stores/configurations';
	import { activeTab } from '$lib/stores/activeTab';
	import UnSelectedProject from './UnSelectedProject.svelte';

	let loading = $state(true);
	let error = $state('');

	// Tabs that don't require a selected project (account settings, instance &
	// workspace settings). For these, ContentArea must render even when no
	// project is selected — otherwise the welcome screen swallows the click.
	const projectIndependentTabs = ['settings', 'instance-settings', 'workspace-settings'];
	let showWelcome = $derived(!$selectedProject && !projectIndependentTabs.includes($activeTab));

	onMount(async () => {
		console.log('onMount: home');
		loading = true;
		try {
			await getProjects();
		} catch (err) {
			error = 'Failed to fetch data';
		} finally {
			loading = false;
		}
	});

	// Quick stats derived from projects
	let stats = $derived({
		total: $projects.length,
		running: $projects.filter((p) => p.status === 'running').length,
		mock: $projects.filter((p) => p.mode === 'mock').length,
		proxy: $projects.filter((p) => p.mode === 'proxy').length
	});
</script>

<svelte:head>
	<title>Beo Echo - API Mock Service Dashboard</title>
	<meta name="description" content="Manage your API mocks and configurations with Beo Echo" />
</svelte:head>

{#if loading}
	<div class="flex items-center justify-center h-full min-h-[300px]">
		<BeoEchoLoader size="lg" animated={true} />
	</div>
{:else if error}
	<div class="text-red-500 text-center p-4">{error}</div>
{:else if showWelcome}
	<!-- Home Page - shown when no project is selected and the active tab needs one -->
	 <UnSelectedProject />
{:else}
	<!-- Show ContentArea when a project is selected, or for project-independent tabs -->
	<ContentArea />
{/if}
