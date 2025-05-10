<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { getProjects, startMockServer, stopMockServer, uploadConfig, type ProjectResponse } from '$lib/api/mockoonApi';
  import { selectedProject } from '$lib/stores/selectedConfig';
  import { activeTab } from '$lib/stores/activeTab';
	import { projects } from '$lib/stores/configurations';
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
    selectConfiguration: ProjectResponse;
  }>();

  $: filteredConfigurations = $projects.filter(project =>
    project.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  let uploading = false;
  let fileInput: HTMLInputElement | null = null;

  function handleConfigClick(project: ProjectResponse) {
    console.log('1. ConfigurationList - Clicked config:', project);
    selectedProject.set(project);
    activeTab.set('routes');
    dispatch('selectConfiguration', project);
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
      projects.set(await getProjects());
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
    <i class="fas fa-server text-5xl mr-4"></i> Beo Echo
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
      {#each filteredConfigurations as project}
        <div
          role="button"
          tabindex="0"
          class="bg-gray-700 p-4 rounded cursor-pointer hover:bg-gray-600 transition-colors"
          class:border-2={$selectedProject?.id === project.id}
          class:border-blue-500={$selectedProject?.id === project.id}
          on:click={() => handleConfigClick(project)}
          on:keydown={(e) => e.key === 'Enter' && handleConfigClick(project)}
        >
          <div class="flex justify-between items-start mb-2">
            <h2 class="text-sm font-bold flex items-center">
              {#if $selectedProject?.id === project.id}
                <i class="fas fa-edit text-blue-500 mr-2"></i>
              {/if}
              <span class="truncate">{project.name}</span>
            </h2>
            <span class="bg-blue-600 text-xs px-2 py-0.5 rounded-full text-white uppercase">{project.mode}</span>
          </div>
          
          <div class="mt-2 space-y-1.5">
            <div class="flex items-center text-xs">
              <i class="fas fa-link text-blue-400 mr-1.5 w-4"></i>
              <a href={project.url} class="text-blue-400 hover:underline truncate" target="_blank" title={project.url}>
                {project.url}
              </a>
            </div>
            
            <div class="flex items-center text-xs">
              <i class="fas fa-tag text-gray-400 mr-1.5 w-4"></i>
              <span class="text-gray-300 truncate" title={project.alias || "No alias"}>
                {project.alias || "—"}
              </span>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
