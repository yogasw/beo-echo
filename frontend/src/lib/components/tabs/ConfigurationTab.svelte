<script lang="ts">
	import { deleteProject, type Project } from '$lib/api/BeoApi';
	import { goto } from '$app/navigation';
	import { projects } from '$lib/stores/configurations';
	import { selectedProject as selectedProjectStore } from '$lib/stores/selectedConfig';
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let selectedProject: Project;

	let showDeleteConfirm = false;
	// Notification for save/delete operations
	let notification = { show: false, message: '', type: 'success' };
	// Map to track expanded sections
	let expandedSections = {
		general: true,
		advanced: false
	};
	
	// Function to toggle section expansion
	function toggleSection(section: string) {
		expandedSections[section] = !expandedSections[section];
		expandedSections = expandedSections; // Force Svelte reactivity update
	}
	
	function handleSave() {
		// Here you would typically save the configuration
		console.log('Saving configuration:', selectedProject);
		
		// Show success notification
		showNotification('Project saved successfully!', 'success');
	}

	async function handleDelete() {
		try {
			await deleteProject(selectedProject.id);
			// Update configurations store
			projects.update(configs => configs.filter(c => c.id !== selectedProject.id));
			// Set selectedProject to null
			selectedProjectStore.set(null);
			// Redirect to home
			await goto('/home');
			
			// Show success notification 
			// (though user will be redirected, keeping it for consistency)
			showNotification('Project deleted successfully!', 'success');
		} catch (error) {
			console.error('Failed to delete project:', error);
			showNotification('Failed to delete project: ' + (error instanceof Error ? error.message : String(error)), 'error');
		}
	}
	
	function showNotification(message: string, type: 'success' | 'error' = 'success') {
		notification = { show: true, message, type };
		setTimeout(() => {
			notification = { ...notification, show: false };
		}, 3000);
	}
</script>

<div class="w-full theme-bg-primary p-4 relative">
	<!-- Notification toast -->
	{#if notification.show}
		<div 
			transition:fade={{ duration: 200 }}
			class="fixed top-6 right-6 {notification.type === 'success' ? 'theme-bg-secondary' : 'bg-red-100 dark:bg-red-800'} theme-text-primary px-4 py-2 rounded shadow-lg z-50 flex items-center"
		>
			<i class="fas {notification.type === 'success' ? 'fa-check-circle text-green-400' : 'fa-exclamation-circle text-red-400'} mr-2"></i>
			<span>{notification.message}</span>
		</div>
	{/if}

	<div class="mb-6">
		<div class="flex justify-between items-center mb-4">
			<div class="flex items-center">
				<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
					<i class="fas fa-cogs text-blue-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">{selectedProject.name}</h2>
					<p class="text-sm theme-text-muted">Project Configuration</p>
				</div>
				{#if selectedProject.url}
				<div class="ml-4 flex items-center bg-gray-100/50 dark:bg-gray-900/50 px-3 py-1 rounded-full">
					<span class="text-xs font-medium text-blue-500 dark:text-blue-400">
						{selectedProject.url}
					</span>
				</div>
				{/if}
			</div>
			
			<div class="flex items-center space-x-3">
				<button 
					class={ThemeUtils.primaryButton('py-2 px-4 rounded-md text-sm')}
					on:click={handleSave}
				>
					<i class="fas fa-save mr-2"></i> Save Changes
				</button>
				
				<button 
					class={ThemeUtils.destructiveButton('py-2 px-4 rounded-md text-sm')}
					on:click={() => showDeleteConfirm = true}
				>
					<i class="fas fa-trash-alt mr-2"></i> Delete
				</button>
			</div>
		</div>
	</div>
	
	<div class="space-y-4">
		<!-- General Information Section -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div 
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => toggleSection('general')}
				on:keydown={(e) => e.key === 'Enter' && toggleSection('general')}
				tabindex="0"
				role="button"
				aria-expanded={expandedSections.general}
			>
				<div class="flex items-center">
					<div class="bg-blue-600/10 p-1.5 rounded mr-2">
						<i class="fas fa-info-circle text-blue-500"></i>
					</div>
					<h3 class="font-medium theme-text-primary">General Information</h3>
				</div>
				<i class="fas {expandedSections.general ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
			</div>
			
			{#if expandedSections.general}
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
									bind:value={selectedProject.name}
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
									bind:value={selectedProject.alias}
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
									<span class="theme-text-secondary">{selectedProject.url}</span>
									<span class="ml-2 text-xs bg-gray-300 dark:bg-gray-600 px-1.5 py-0.5 rounded theme-text-secondary">Read only</span>
								</div>
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
					</div>
				</div>
			{/if}
		</div>
		
		<!-- Advanced Settings Section -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div 
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => toggleSection('advanced')}
				on:keydown={(e) => e.key === 'Enter' && toggleSection('advanced')}
				tabindex="0"
				role="button"
				aria-expanded={expandedSections.advanced}
			>
				<div class="flex items-center">
					<div class="bg-purple-600/10 p-1.5 rounded mr-2">
						<i class="fas fa-sliders-h text-purple-500"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Advanced Settings</h3>
				</div>
				<i class="fas {expandedSections.advanced ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
			</div>
			
			{#if expandedSections.advanced}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
					<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label for="config-created" class="block text-sm font-medium mb-2 theme-text-secondary">Created On</label>
							<div class="relative">
								<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
									<i class="fas fa-calendar-alt theme-text-muted"></i>
								</div>
								<div class="flex items-center theme-bg-secondary theme-border border rounded-lg px-3 py-3 ps-10">
									<span class="theme-text-secondary">{new Date(selectedProject.created_at || Date.now()).toLocaleString()}</span>
								</div>
							</div>
						</div>
					</div>
					
					<!-- New options in the advanced section -->
					<div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
						<div>
							<label for="config-timeout" class="block text-sm font-medium mb-2 theme-text-secondary">Timeout (ms)</label>
							<div class="relative">
								<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
									<i class="fas fa-clock theme-text-muted"></i>
								</div>
								<input
									type="number"
									id="config-timeout"
									class={ThemeUtils.inputField()}
									value={selectedProject?.timeout || 10000}
									placeholder="Response timeout in ms"
								/>
							</div>
						</div>
						
						<div>
							<label for="config-cors" class="block text-sm font-medium mb-2 theme-text-secondary">CORS Enabled</label>
							<div class="mt-3">
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" value="" class="sr-only peer" checked>
									<div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600"></div>
									<span class="ms-3 text-sm font-medium theme-text-secondary">Enable CORS</span>
								</label>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
		
	</div>
</div>

{#if showDeleteConfirm}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" transition:fade={{ duration: 150 }}>
		<div class={ThemeUtils.card('p-6 rounded-lg max-w-md w-full')}>
			<h3 class="text-xl font-bold mb-4 theme-text-primary flex items-center">
				<i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-2"></i>
				Confirm Delete
			</h3>
			<p class="mb-6 theme-text-secondary">Are you sure you want to delete <span class="font-semibold theme-text-primary">{selectedProject.name}</span>? This action cannot be undone.</p>
			<div class="flex justify-end space-x-4">
				<button 
					class={ThemeUtils.secondaryButton('py-2 px-4 rounded')} 
					on:click={() => showDeleteConfirm = false}
				>
					<i class="fas fa-times mr-1"></i>
					Cancel
				</button>
				<button 
					class={ThemeUtils.destructiveButton('py-2 px-4 rounded')}
					on:click={handleDelete}
				>
					<i class="fas fa-trash-alt mr-1"></i>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}
