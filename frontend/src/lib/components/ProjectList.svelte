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
	import { initializeLogsStream } from '$lib/services/logsService';
	import { logStatus } from '$lib/stores/logStatus';
	import { setCurrentWorkspaceId, getProjectPanelWidth, setProjectPanelWidth } from '$lib/utils/localStorage';

	export let searchTerm = '';
	export let panelWidth: number = getProjectPanelWidth(); // Panel width in rem units (w-72 = 18rem)

	const dispatch = createEventDispatcher<{
		selectedProject: Project;
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

	// Resizable panel variables
	let isResizing = false;
	let startX = 0;
	let startWidth = panelWidth;

	// Resize functions
	function startResize(event: MouseEvent) {
		isResizing = true;
		startX = event.clientX;
		startWidth = panelWidth;
		
		document.addEventListener('mousemove', handleResize);
		document.addEventListener('mouseup', stopResize);
		document.body.style.cursor = 'col-resize';
		document.body.style.userSelect = 'none';
	}

	function handleResize(event: MouseEvent) {
		if (!isResizing) return;
		
		const deltaX = event.clientX - startX;
		const containerWidth = window.innerWidth;
		const newWidthRem = startWidth + (deltaX / containerWidth) * 100; // Convert to rem-like units
		
		// Constrain between 12rem and 30rem (minimum 192px, maximum 480px)
		panelWidth = Math.min(Math.max(newWidthRem, 12), 30);
	}

	function stopResize() {
		isResizing = false;
		document.removeEventListener('mousemove', handleResize);
		document.removeEventListener('mouseup', stopResize);
		document.body.style.cursor = '';
		document.body.style.userSelect = '';
		
		// Save panel width to localStorage when resize is complete
		setProjectPanelWidth(panelWidth);
	}

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
		let isResetLogs = project.id != $selectedProject?.id;
		if (isResetLogs) {
			logStatus.reset();
			initializeLogsStream(project.id, 100, isResetLogs);
		}
		
		isLoadingContentArea.set(true);

		try {
			const fullConfig = await getProjectDetail(project.id);
			selectedProject.set(project);
			// Disable switching tabs when click project
			// activeTab.set('routes');

			// Reset endpoints update list when changing projects
			resetEndpointsList();
			dispatch('selectedProject', project);
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
			projects.set(await getProjects());
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
			setCurrentWorkspaceId(workspaceId);
			const projectsData = await getProjects();
			projects.set(projectsData);
		} catch (err) {
			console.error('Failed to fetch projects for workspace:', workspaceId, err);
			toast.error('Failed to load projects');
		}
	}
</script>

<div class="theme-bg-primary p-4 flex flex-col h-full border-r theme-border relative" style="width: {panelWidth}rem;">
	<!-- Brand Section - Clean & Responsive -->
	<div class="mb-4 mt-2">
		<!-- Main Brand Header -->
		<div class="flex items-center justify-between mb-4 p-4 rounded-xl bg-gradient-to-r from-blue-600/10 to-purple-600/10 border border-blue-500/20 dark:border-blue-400/20">
			<div class="flex items-center">
				<div class="w-14 h-14 mr-4 flex items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-purple-600 p-3 shadow-lg">
					<img 
						src="/favicon.svg" 
						alt="Beo Echo Logo" 
						class="w-full h-full object-contain filter brightness-110"
						title="Beo Echo - API Mocking Service"
						aria-label="Beo Echo API Mocking Service logo"
					/>
				</div>
				{#if panelWidth >= 12}
					<div class="flex flex-col">
						<h1 class="font-bold theme-text-primary leading-tight tracking-tight bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent {panelWidth >= 16 ? 'text-2xl' : 'text-lg'}">
							Beo Echo
						</h1>
						{#if panelWidth >= 16}
							<span class="text-sm theme-text-secondary font-medium opacity-80">
								API Mocking Service
							</span>
						{:else}
							<span class="text-xs theme-text-secondary font-medium opacity-70">
								API Mocking Service
							</span>
						{/if}
					</div>
				{/if}
			</div>
		</div>
		
		<!-- Version and Links -->
		<div class="flex items-center justify-between px-2">
			<span class="inline-flex items-center px-3 py-1.5 rounded-full text-xs font-semibold bg-gradient-to-r from-blue-500 to-blue-600 text-white shadow-sm">
				{#if panelWidth >= 16}
					<i class="fas fa-tag mr-1.5 text-xs"></i>
				{/if}
				v2.3.2
			</span>
			
			<!-- Action Links -->
			<div class="flex items-center space-x-2">
				<a 
					href="https://github.com/yogasw/beo-echo" 
					target="_blank" 
					rel="noopener noreferrer"
					class="group flex items-center justify-center w-8 h-8 rounded-lg theme-bg-secondary hover:bg-gray-600 dark:hover:bg-gray-500 transition-all duration-200 transform hover:scale-105"
					style="text-decoration: none !important;"
					title="View on GitHub"
					aria-label="View Beo Echo project on GitHub"
				>
					<i class="fab fa-github text-sm theme-text-secondary group-hover:text-white transition-colors"></i>
				</a>
			</div>
		</div>
		
	</div>
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
		title="Upload configuration file"
		aria-label="Upload configuration file"
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
		title="Add new project"
		aria-label="Add new project"
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
						title="Cancel"
						aria-label="Cancel"
					>
						<i class="fas fa-times mr-2"></i> Cancel
					</button>
					<button
						class={isAddingProject || !projectName.trim() || !projectAlias.trim()
							? ThemeUtils.secondaryButton('px-4 py-2 cursor-not-allowed opacity-70')
							: ThemeUtils.primaryButton('px-4 py-2')}
						on:click={handleAddProject}
						disabled={isAddingProject || !projectName.trim() || !projectAlias.trim()}
						title={isAddingProject ? 'Creating project...' : 'Create project'}
						aria-label={isAddingProject ? 'Creating project...' : 'Create project'}
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
							<span class="truncate" title={project.name}>
								{project.name.length > 15 ? project.name.slice(0, 15) + '…' : project.name}
							</span>
						</h2>
						<div class="flex items-center space-x-2">
							<span class={ThemeUtils.badge('info', 'text-xs px-2 py-0.5 uppercase')}>
								{project.mode}
							</span>
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
				</div>			{/each}
		</div>
	</div>
	
	<!-- Resizable handle -->
	<button
		class="absolute top-0 right-0 bottom-0 w-1 cursor-col-resize hover:bg-blue-500 transition-colors duration-200 group bg-transparent border-none"
		on:mousedown={startResize}
		title="Drag to resize panel"
		aria-label="Resize panel"
	>
		<div class="w-full h-full bg-transparent group-hover:bg-blue-500/30"></div>
	</button>
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
