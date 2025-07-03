<script lang="ts">
	import ContentArea from '$lib/components/ContentArea.svelte';
	import { onMount } from 'svelte';
	import { getProjects } from '$lib/api/BeoApi';

	interface Route {
    path: string;
    method: string;
    status: 'enabled' | 'disabled';
  }

  let routes: Route[] = [
    { path: "/api/v1/resource/long/path/example", method: "GET", status: "enabled" },
    { path: "/api/v1/resource/another/long/path", method: "POST", status: "disabled" },
    { path: "/api/v1/example/long/path", method: "PUT", status: "enabled" }
  ];

  let configurations = [];
  let mockStatus = [];
  let loading = true;
  let error = '';

  let activeContentTab = 'Status & Body';

  // Get state from parent layout
  export let activeTab = 'routes';

  onMount(async () => {
    console.log("onMount: home");
    loading = true;
    try {
      configurations = await getProjects();
    } catch (err) {
      error = 'Failed to fetch data';
    } finally {
      loading = false;
    }
  });
</script>

<svelte:head>
	<title>Beo Echo - API Mock Service Dashboard</title>
	<meta name="description" content="Manage your API mocks and configurations with Beo Echo" />
</svelte:head>

<!-- Mobile Not Supported Message -->
<div class="fixed inset-0 bg-gray-900 text-white z-50 flex flex-col justify-center items-center text-center p-8 md:hidden">
	<div class="text-5xl mb-6 text-red-500">
		<i class="fas fa-desktop"></i>
	</div>
	<h1 class="text-2xl font-bold mb-4 text-red-500">Desktop Only</h1>
	<p class="text-lg mb-2 leading-relaxed">Sorry, this application only supports desktop devices.</p>
	<p class="text-lg mb-2 leading-relaxed">Please access using a computer or laptop with a larger screen.</p>
	<p class="text-base text-gray-300">Minimum resolution: 768px width</p>
</div>

<!-- Main Content (hidden on mobile) -->
<div class="hidden md:block">
	{#if loading}
	  <div class="flex items-center justify-center h-full min-h-[300px]">
	    <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-blue-500"></div>
	  </div>
	{:else if error}
	  <div class="text-red-500 text-center p-4">{error}</div>
	{:else}
	  <ContentArea
	    {activeContentTab}
	  />
	{/if}
</div>
