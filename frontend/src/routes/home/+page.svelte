<script lang="ts">
	import ContentArea from '$lib/components/ContentArea.svelte';
	import { onMount } from 'svelte';
	import { getProjects, getMockStatus } from '$lib/api/mockoonApi';

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
      mockStatus = await getMockStatus();
    } catch (err) {
      error = 'Failed to fetch data';
    } finally {
      loading = false;
    }
  });
</script>

{#if loading}
  <div class="flex items-center justify-center h-full min-h-[300px]">
    <div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-blue-500"></div>
  </div>
{:else if error}
  <div class="text-red-500 text-center p-4">{error}</div>
{:else}
  <ContentArea
    {routes}
    {activeTab}
    {activeContentTab}
  />
{/if}
