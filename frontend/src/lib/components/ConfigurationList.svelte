<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import {
		getProjects,
		startMockServer,
		stopMockServer,
		uploadConfig,
		addProject,
		type ProjectResponse,
		type Project
	} from '$lib/api/mockoonApi';
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

	$: filteredConfigurations = $projects.filter((project) =>
		project.name.toLowerCase().includes(searchTerm.toLowerCase())
	);

	let uploading = false;
	let fileInput: HTMLInputElement | null = null;

	// For Add Project modal
	let showAddProjectModal = false;
	let projectName = '';
	let projectAlias = '';
	let isAddingProject = false;
	let userEditedAlias = false;

	// Function to generate alias from project name
	function generateAlias(name: string): string {
		return name
			.toLowerCase()
			.replace(/\s+/g, '-')
			.replace(/[^a-z0-9-]/g, '');
	}

	// Update project alias when project name changes and user hasn't manually edited the alias
	$: if (projectName && !userEditedAlias) {
		projectAlias = generateAlias(projectName);
	}

	// Track when user manually edits the alias
	function handleAliasInput() {
		userEditedAlias = true;
	}

	// Reset the tracking when modal is opened or closed
	function resetAliasTracking() {
		userEditedAlias = false;
	}

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

	function openAddProjectModal() {
		showAddProjectModal = true;
		resetAliasTracking();
	}

	function closeAddProjectModal() {
		showAddProjectModal = false;
		projectName = '';
		projectAlias = '';
		resetAliasTracking();
	}

	async function handleAddProject() {
		if (!projectName.trim()) {
			toast.error('Project name is required');
			return;
		}

		if (!projectAlias.trim()) {
			toast.error('Project alias is required');
			return;
		}

		isAddingProject = true;
		try {
			await addProject(projectName.trim(), projectAlias.trim());
			// Refresh project list
			projects.set(await getProjects());
			toast.success('Project created successfully');
			closeAddProjectModal();
		} catch (err) {
			toast.error('Failed to create project');
		} finally {
			isAddingProject = false;
		}
	}
</script>

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
	<button
		class="bg-blue-500 text-white py-2 px-4 rounded mb-2 w-full flex items-center justify-center"
		on:click={triggerFileInput}
		disabled={uploading}
	>
		<i class="fas fa-upload mr-2"></i>
		{uploading ? 'Uploading...' : 'Upload Config'}
	</button>
	<input
		type="file"
		accept=".json"
		class="hidden"
		bind:this={fileInput}
		on:change={handleUploadConfig}
	/>

	<button
		class="bg-green-600 text-white py-2 px-4 rounded mb-4 w-full flex items-center justify-center"
		on:click={openAddProjectModal}
	>
		<i class="fas fa-plus mr-2"></i> Add Project
	</button>

	<!-- Add Project Modal -->
	{#if showAddProjectModal}
		<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
			<div class="bg-gray-800 p-6 rounded-lg max-w-md w-full mx-4">
				<h2 class="text-xl font-bold mb-4 text-white">Add New Project</h2>

				<div class="mb-4">
					<label for="projectName" class="block text-sm font-medium text-gray-300 mb-1"
						>Project Name</label
					>
					<input
						id="projectName"
						type="text"
						class="w-full bg-gray-700 text-white py-2 px-3 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						bind:value={projectName}
						placeholder="Enter project name"
					/>
				</div>

				<div class="mb-4">
					<label for="projectAlias" class="block text-sm font-medium text-gray-300 mb-1"
						>Project Alias</label
					>
					<input
						id="projectAlias"
						type="text"
						class="w-full bg-gray-700 text-white py-2 px-3 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
						bind:value={projectAlias}
						placeholder="Enter project alias"
						on:input={handleAliasInput}
					/>
				</div>

				<div class="mb-6">
					<p class="block text-sm font-medium text-gray-300 mb-1">URL Preview</p>
					<div
						class="bg-gray-900 px-3 py-2 rounded border border-gray-600 text-gray-300 font-mono text-sm break-all"
					>
						http://BASE_URL/{projectAlias || '[alias]'}
					</div>
				</div>
				<div class="flex justify-end space-x-2">
					<button
						class="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 transition-colors"
						on:click={closeAddProjectModal}
						disabled={isAddingProject}
					>
						Cancel
					</button>
					<button
						class={`px-4 py-2 text-white rounded transition-colors flex items-center ${isAddingProject || !projectName.trim() || !projectAlias.trim() ? 'bg-gray-700 cursor-not-allowed' : 'bg-blue-500 hover:bg-blue-600'}`}
						on:click={handleAddProject}
						disabled={isAddingProject || !projectName.trim() || !projectAlias.trim()}
					>
						{#if isAddingProject}
							<i class="fas fa-spinner fa-spin mr-2"></i>
						{/if}
						{isAddingProject ? 'Creating...' : 'Create Project'}
					</button>
				</div>
			</div>
		</div>
	{/if}
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
						<span class="bg-blue-600 text-xs px-2 py-0.5 rounded-full text-white uppercase"
							>{project.mode}</span
						>
					</div>

					<div class="mt-2 space-y-1.5">
						<div class="flex items-center text-xs">
							<i class="fas fa-link text-blue-400 mr-1.5 w-4"></i>
							<a
								href={project.url}
								class="text-blue-400 hover:underline truncate"
								target="_blank"
								title={project.url}
							>
								{project.url}
							</a>
						</div>

						<div class="flex items-center text-xs">
							<i class="fas fa-tag text-gray-400 mr-1.5 w-4"></i>
							<span class="text-gray-300 truncate" title={project.alias || 'No alias'}>
								{project.alias || '—'}
							</span>
						</div>
					</div>
				</div>
			{/each}
		</div>
	</div>
</div>

<style>
	/* Hide scrollbar for Chrome, Safari and Opera */
	.hide-scrollbar::-webkit-scrollbar {
		display: none;
	}
	/* Hide scrollbar for IE, Edge and Firefox */
	.hide-scrollbar {
		-ms-overflow-style: none; /* IE and Edge */
		scrollbar-width: none; /* Firefox */
	}
</style>
