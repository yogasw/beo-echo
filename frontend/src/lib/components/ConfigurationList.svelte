<script lang="ts">
	import { createEventDispatcher, onMount } from 'svelte';
	import {
		getProjects,
		uploadConfig,
		addProject,
		updateProjectStatus,
		type Project,
		getProjectDetail
	} from '$lib/api/BeoApi';
	import ProjectStatusBadge from '$lib/components/ProjectStatusBadge.svelte';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { activeTab } from '$lib/stores/activeTab';
	import { projects } from '$lib/stores/configurations';
	import { toast } from '$lib/stores/toast';
	import { resetEndpointsList } from '$lib/stores/saveButton';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { currentWorkspace } from '$lib/stores/workspace';
	import { isLoadingContentArea } from '$lib/stores/loadingContentArea';

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
		selectConfiguration: Project;
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

	// For project status updates
	let updatingStatus: string | null = null;

	// Function to generate alias from project name
	function generateAlias(name: string): string {
		// First convert to lowercase and replace spaces with hyphens
		// Then remove all characters except lowercase letters, numbers, underscores and hyphens
		return name
			.toLowerCase()
			.replace(/\s+/g, '-')
			.replace(/[^a-z0-9_-]/g, '');
	}

	// Update project alias when project name changes and user hasn't manually edited the alias
	$: if (projectName && !userEditedAlias) {
		projectAlias = generateAlias(projectName);
	}

	// Track when user manually edits the alias and enforce validation
	function handleAliasInput(event: Event) {
		userEditedAlias = true;

		// Get input element and current value
		const input = event.target as HTMLInputElement;
		const currentValue = input.value;

		// Apply validation rules: lowercase, only allow lowercase letters, numbers, underscores and hyphens
		const validatedValue = currentValue.toLowerCase().replace(/[^a-z0-9_-]/g, '');

		// Update the value if it was changed by validation
		if (currentValue !== validatedValue) {
			projectAlias = validatedValue;
		}
	}

	// Reset the tracking when modal is opened or closed
	function resetAliasTracking() {
		userEditedAlias = false;
	}

	async function handleConfigClick(project: Project) {
		console.log('1. ConfigurationList - Clicked config:', project);
		isLoadingContentArea.set(true);
		try {
			const fullConfig = await getProjectDetail(project.id);
			selectedProject.set(project);
			activeTab.set('routes');
			// Reset endpoints update list when changing projects
			resetEndpointsList();
			dispatch('selectConfiguration', project);
			// Parse routes from config
			// endpoints = fullConfig.endpoints;
			fullConfig.url = project.url;
			selectedProject.set(fullConfig);
			console.log('Config data loaded:', fullConfig);
		} catch (err) {
			console.error('Failed to load config data:', err);
		} finally {
			isLoadingContentArea.set(false);
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

		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}

		isAddingProject = true;
		try {
			await addProject(projectName.trim(), projectAlias.trim());
			// Refresh project list
			projects.set(await getProjects($currentWorkspace.id));
			toast.success('Project created successfully');
			closeAddProjectModal();
		} catch (err) {
			toast.error(err);
		} finally {
			isAddingProject = false;
		}
	}

	async function handleUpdateStatus(project: Project, newStatus: string) {
		if (updatingStatus === project.id) return; // Prevent double clicks
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}

		updatingStatus = project.id;
		try {
			await updateProjectStatus(project.id, newStatus);
			// Update local project status
			$projects = $projects.map((p) => (p.id === project.id ? { ...p, status: newStatus } : p));

			// If the currently selected project is modified, update that too
			if ($selectedProject?.id === project.id) {
				selectedProject.set({ ...$selectedProject, status: newStatus });
			}

			toast.success(`Project ${newStatus === 'running' ? 'started' : 'stopped'} successfully`);
		} catch (err) {
			toast.error(`Failed to ${newStatus === 'running' ? 'start' : 'stop'} project`);
			console.error('Error updating project status:', err);
		} finally {
			updatingStatus = null;
		}
	}

	// Refresh projects when current workspace changes
	$: if ($currentWorkspace) {
		refreshProjects($currentWorkspace.id);
	}

	async function refreshProjects(workspaceId: string) {
		try {
			const projectsData = await getProjects(workspaceId);
			projects.set(projectsData);
		} catch (err) {
			console.error('Failed to fetch projects for workspace:', workspaceId, err);
			toast.error('Failed to load projects');
		}
	}
</script>

<div class="w-72 theme-bg-primary p-4 flex flex-col h-full border-r theme-border">
	<h1 class="text-xl font-bold mb-4 flex items-center theme-text-primary">
		<i class="fas fa-server text-5xl mr-4"></i> Beo Echo
	</h1>
	<div class="relative mb-4">
		<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
			<i class="fas fa-search text-gray-400"></i>
		</div>
		<input
			type="text"
			bind:value={searchTerm}
			placeholder="Search Configuration"
			class={ThemeUtils.inputField('py-2')}
		/>
	</div>
	<button
		class={ThemeUtils.primaryButton('mb-2 w-full justify-center')}
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
		class={ThemeUtils.primaryButton('mb-4 w-full justify-center bg-green-600 hover:bg-green-700')}
		on:click={openAddProjectModal}
	>
		<i class="fas fa-plus mr-2"></i> Add Project
	</button>

	<!-- Add Project Modal -->
	{#if showAddProjectModal}
		<div class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
			<div class={ThemeUtils.card('p-6 max-w-md w-full mx-4')}>
				<h2 class={ThemeUtils.headerSection('text-xl font-bold mb-4 rounded-md')}>
					Add New Project
				</h2>

				<div class="mb-4">
					<label for="projectName" class="block text-sm font-medium theme-text-secondary mb-1"
						>Project Name</label
					>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-tag text-gray-400"></i>
						</div>
						<input
							id="projectName"
							type="text"
							class={ThemeUtils.inputField('')}
							bind:value={projectName}
							placeholder="Enter project name"
						/>
					</div>
				</div>

				<div class="mb-4">
					<label for="projectAlias" class="block text-sm font-medium theme-text-secondary mb-1"
						>Project Alias</label
					>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-link text-gray-400"></i>
						</div>
						<input
							id="projectAlias"
							type="text"
							class={ThemeUtils.inputField('')}
							bind:value={projectAlias}
							placeholder="Enter project alias"
							on:input={handleAliasInput}
						/>
					</div>
					<p class="text-xs theme-text-muted mt-1">
						Only lowercase letters, numbers, underscores (_) and hyphens (-) allowed
					</p>
				</div>

				<div class="mb-6">
					<p class="block text-sm font-medium theme-text-secondary mb-1">URL Preview</p>
					<div
						class={ThemeUtils.themeBgTertiary(
							'px-3 py-2 rounded theme-border border font-mono text-sm break-all theme-text-secondary'
						)}
					>
						http://BASE_URL/{projectAlias || '[alias]'}
					</div>
				</div>
				<div class="flex justify-end space-x-2">
					<button
						class={ThemeUtils.secondaryButton('px-4 py-2 rounded transition-colors')}
						on:click={closeAddProjectModal}
						disabled={isAddingProject}
					>
						<i class="fas fa-times mr-2"></i> Cancel
					</button>
					<button
						class={isAddingProject || !projectName.trim() || !projectAlias.trim()
							? ThemeUtils.secondaryButton('px-4 py-2 cursor-not-allowed opacity-70')
							: ThemeUtils.primaryButton('px-4 py-2')}
						on:click={handleAddProject}
						disabled={isAddingProject || !projectName.trim() || !projectAlias.trim()}
					>
						{#if isAddingProject}
							<i class="fas fa-spinner fa-spin mr-2"></i>
						{:else}
							<i class="fas fa-save mr-2"></i>
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
					class={ThemeUtils.themeBgSecondary(`p-4 rounded cursor-pointer transition-colors 
					${$selectedProject?.id === project.id ? 'border-2 border-blue-500' : 'theme-border border'}
					${$selectedProject?.id !== project.id ? ThemeUtils.themeHover('') : ''}`)}
					on:click={() => handleConfigClick(project)}
					on:keydown={(e) => e.key === 'Enter' && handleConfigClick(project)}
				>
					<div class="flex justify-between items-start mb-2">
						<h2 class="text-sm font-bold flex items-center theme-text-primary">
							{#if $selectedProject?.id === project.id}
								<i class="fas fa-edit text-blue-500 mr-2"></i>
							{/if}
							<span class="truncate">{project.name}</span>
						</h2>
						<div class="flex items-center space-x-2">
							<span class={ThemeUtils.badge('info', 'text-xs px-2 py-0.5 uppercase')}
								>{project.mode}</span
							>
						</div>
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
							<i class="fas fa-tag theme-text-muted mr-1.5 w-4"></i>
							<span class="theme-text-secondary truncate" title={project.alias || 'No alias'}>
								{project.alias || '—'}
							</span>
						</div>

						<!-- Status indicator with live animation -->
						<div class="flex items-center text-xs">
							<ProjectStatusBadge status={project.status || 'stopped'} size="small" />
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
