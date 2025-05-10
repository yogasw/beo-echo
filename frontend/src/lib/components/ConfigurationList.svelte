<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { getProjects, startMockServer, stopMockServer, uploadConfig } from '$lib/api/mockoonApi';
  import { selectedConfig } from '$lib/stores/selectedConfig';
  import { configurations } from '$lib/stores/configurations';
  import { activeTab } from '$lib/stores/activeTab';
  import { toast } from '$lib/stores/toast';

  interface Config {
    uuid: string;
    name: string;
    configFile: string;
    port: number;
    url: string;
    size: string;
    modified: string;
    inUse: boolean;
  }

  export let searchTerm = '';

  const dispatch = createEventDispatcher<{
    selectConfiguration: Config;
    startConfiguration: Config;
    stopConfiguration: Config;
  }>();

  $: filteredConfigurations = $configurations.filter(config =>
    config.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  let uploading = false;
  let fileInput: HTMLInputElement | null = null;

  function handleConfigClick(config: Config) {
    console.log('1. ConfigurationList - Clicked config:', config);
    selectedConfig.set(config);
    activeTab.set('routes');
    dispatch('selectConfiguration', config);
  }

  async function handleStartStop(config: Config, event: MouseEvent) {
    event.stopPropagation();
    try {
      if (config.inUse) {
        await stopMockServer(config.port);
        dispatch('stopConfiguration', config);
      } else {
        await startMockServer(config.port, config.configFile, config.uuid);
        dispatch('startConfiguration', config);
      }
    } catch (err) {
      console.error('Failed to start/stop server:', err);
    }
  }

  async function handleUploadConfig(event: Event) {
    const files = (event.target as HTMLInputElement).files;
    if (!files || files.length === 0) return;
    const file = files[0];
    const formData = new FormData();
    formData.append('config', file);
    uploading = true;
    try {
      await uploadConfig(formData);
      // Refresh config list
      configurations.set(await getProjects());
      toast.success('Config uploaded successfully');
    } catch (err) {
      toast.error('Failed to upload config');
    } finally {
      uploading = false;
      if (fileInput) fileInput.value = '';
    }
  }

  function triggerFileInput() {
    if (fileInput) fileInput.click();
  }
</script>

<style>
  /* Hide scrollbar for Chrome, Safari and Opera */
  .hide-scrollbar::-webkit-scrollbar {
    display: none;
  }
  /* Hide scrollbar for IE, Edge and Firefox */
  .hide-scrollbar {
    -ms-overflow-style: none;  /* IE and Edge */
    scrollbar-width: none;  /* Firefox */
  }
</style>

<div class="w-72 bg-gray-800 p-4 flex flex-col h-full">
  <h1 class="text-xl font-bold mb-4 flex items-center">
    <i class="fas fa-server text-5xl mr-4"></i> Mockoon Control Panel
  </h1>
  <div class="flex items-center bg-gray-700 py-2 px-4 rounded mb-4">
    <i class="fas fa-search text-white text-lg mr-2"></i>
    <input
      type="text"
      bind:value={searchTerm}
      placeholder="Search Configuration"
      class="w-full bg-gray-700 text-white py-2 px-2 rounded"
    />
  </div>
  <button class="bg-blue-500 text-white py-2 px-4 rounded mb-2 w-full flex items-center justify-center" on:click={triggerFileInput} disabled={uploading}>
    <i class="fas fa-upload mr-2"></i> {uploading ? 'Uploading...' : 'Upload Config'}
  </button>
  <input type="file" accept=".json" class="hidden" bind:this={fileInput} on:change={handleUploadConfig} />
  <!-- Configuration List -->
  <div class="flex-1 min-h-0 overflow-auto hide-scrollbar">
    <div class="space-y-4">
      {#each filteredConfigurations as config}
        <div
          role="button"
          tabindex="0"
          class="bg-gray-700 p-4 rounded cursor-pointer hover:bg-gray-600 transition-colors"
          class:border-2={$selectedConfig?.uuid === config.uuid}
          class:border-blue-500={$selectedConfig?.uuid === config.uuid}
          on:click={() => handleConfigClick(config)}
          on:keydown={(e) => e.key === 'Enter' && handleConfigClick(config)}
        >
          <h2 class="text-sm font-bold flex items-center">
            {#if $selectedConfig?.uuid === config.uuid}
              <i class="fas fa-edit text-blue-500 mr-2"></i>
            {/if}
            {config.name}
          </h2>
					<p class="text-xs">
						URL:
						{#if config.inUse}
							<a href={config.url} class="text-blue-400 hover:underline" target="_blank">{config.url}</a>
						{:else}
							<span class="relative group text-gray-500 cursor-not-allowed">
								{config.url}
								<span class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 w-max bg-gray-800 text-white text-xs rounded px-2 py-1 opacity-0 group-hover:opacity-100 transition-opacity">
									This URL cannot be accessed because it is not in use.
								</span>
							</span>
						{/if}
					</p>
          <p class="text-xs">Port: {config.port}</p>
          <p class="text-xs">File: {config.configFile}</p>
          <p class="text-xs">Size: {config.size}</p>
          <p class="text-xs">Modified: {new Date(config.modified).toLocaleString()}</p>
          <button
            class="text-white py-1 px-4 rounded mt-2 flex items-center"
            class:bg-green-500={config.inUse}
            class:bg-red-500={!config.inUse}
            on:click={(e) => handleStartStop(config, e)}
          >
            {#if config.inUse}
              <i class="fas fa-stop mr-2"></i> Stop
            {:else}
              <i class="fas fa-play mr-2"></i> Start
            {/if}
          </button>
        </div>
      {/each}
    </div>
  </div>
</div>
