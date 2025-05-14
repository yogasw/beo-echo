<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { updateProject, type Project } from '$lib/api/BeoApi';

	export let project: Project;
	export let showNotification: (message: string, type?: 'success' | 'error') => void;

	let isExpanded = true;
	let name = project?.name || '';
	let alias = project?.alias || '';
	let mode = project?.mode || 'mock';
	
	$: url = project?.url || '';
	$: status = project?.status || 'active';

	// Function to toggle section expansion
	function toggleSection() {
		isExpanded = !isExpanded;
	}
	
	// Handle save of general info
	async function handleSave() {
		try {
			const updatedProject = await updateProject(project.id, {
				name,
				alias,
				mode
			});
			
			// Update local project with new values
			project = updatedProject;
			
			// Show success notification
			showNotification('General information updated successfully!', 'success');
		} catch (error) {
			console.error('Failed to update project:', error);
			showNotification('Failed to update project: ' + (error instanceof Error ? error.message : String(error)), 'error');
		}
	}

	// Handle input change and auto-save
	async function handleInputChange() {
		// Only update if values have changed
		if (name !== project.name || alias !== project.alias || mode !== project.mode) {
			await handleSave();
		}
	}

	// Mode options
	const modeOptions = [
		{ value: 'mock', label: 'Mock' },
		{ value: 'proxy', label: 'Proxy' }
	];
</script>

<div class={ThemeUtils.card('overflow-hidden')}>
	<div 
		class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
		on:click={toggleSection}
		on:keydown={(e) => e.key === 'Enter' && toggleSection()}
		tabindex="0"
		role="button"
		aria-expanded={isExpanded}
	>
		<div class="flex items-center">
			<div class="bg-blue-600/10 p-1.5 rounded mr-2">
				<i class="fas fa-info-circle text-blue-500"></i>
			</div>
			<h3 class="font-medium theme-text-primary">General Information</h3>
		</div>
		<i class="fas {isExpanded ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
	</div>
	
	{#if isExpanded}
		<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<div>
					<label for="config-name" class="block text-sm font-medium mb-2 theme-text-secondary">Project Name</label>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-tag theme-text-muted"></i>
						</div>
						<input
							type="text"
							id="config-name"
							class={ThemeUtils.inputField()}
							bind:value={name}
							on:change={handleInputChange}
							placeholder="Enter project name"
						/>
					</div>
				</div>
				
				<div>
					<label for="config-alias" class="block text-sm font-medium mb-2 theme-text-secondary">Project Alias</label>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-bookmark theme-text-muted"></i>
						</div>
						<input
							type="text"
							id="config-alias"
							class={ThemeUtils.inputField()}
							bind:value={alias}
							on:change={handleInputChange}
							placeholder="Enter project alias"
						/>
					</div>
				</div>
				
				<div>
					<label for="config-url" class="block text-sm font-medium mb-2 theme-text-secondary">Base URL</label>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-link theme-text-muted"></i>
						</div>
						<div class="flex items-center theme-bg-secondary theme-border border rounded-lg px-3 py-3 ps-10">
							<span class="theme-text-secondary">{url}</span>
							<span class="ml-2 text-xs bg-gray-300 dark:bg-gray-600 px-1.5 py-0.5 rounded theme-text-secondary">Read only</span>
						</div>
					</div>
				</div>
				
				<div>
					<label for="config-mode" class="block text-sm font-medium mb-2 theme-text-secondary">Mode</label>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-exchange-alt theme-text-muted"></i>
						</div>
						<select
							id="config-mode"
							class={ThemeUtils.inputField()}
							bind:value={mode}
							on:change={handleInputChange}
						>
							{#each modeOptions as option}
								<option value={option.value}>{option.label}</option>
							{/each}
						</select>
					</div>
				</div>

				<div>
					<label for="config-status" class="block text-sm font-medium mb-2 theme-text-secondary">Status</label>
					<div class="flex items-center theme-bg-secondary theme-border border rounded-lg p-3">
						<span class="inline-flex items-center">
							<span class="relative flex h-3 w-3 mr-2">
								<span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-green-400 opacity-75"></span>
								<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
							</span>
							<span class="text-green-500 dark:text-green-400 text-sm">Active</span>
						</span>
					</div>
				</div>
				
				<div>
					<label for="config-created" class="block text-sm font-medium mb-2 theme-text-secondary">Created On</label>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-calendar-alt theme-text-muted"></i>
						</div>
						<div class="flex items-center theme-bg-secondary theme-border border rounded-lg px-3 py-3 ps-10">
							<span class="theme-text-secondary">{new Date(project?.created_at || Date.now()).toLocaleString()}</span>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
