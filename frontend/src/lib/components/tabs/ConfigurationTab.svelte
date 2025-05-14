<script lang="ts">
	import { deleteProject, type Project } from '$lib/api/BeoApi';
	import { goto } from '$app/navigation';
	import { projects } from '$lib/stores/configurations';
	import { selectedProject as selectedProjectStore } from '$lib/stores/selectedConfig';
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	
	// Import the component modules
	import GeneralInfo from './Configuration/GeneralInfo.svelte';
	import AdvancedSettings from './Configuration/AdvancedSettings.svelte';
	import ProxyManagement from './Configuration/ProxyManagement.svelte';

	export let selectedProject: Project;

	let showDeleteConfirm = false;
	// Notification for save/delete operations
	let notification = { show: false, message: '', type: 'success' };
	
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
			</div>
			
			<div class="flex items-center space-x-3">
				<button 
					class={ThemeUtils.destructiveButton('py-2 px-4 rounded-md text-sm')}
					on:click={() => showDeleteConfirm = true}
				>
					<i class="fas fa-trash-alt mr-2"></i> Delete Project
				</button>
			</div>
		</div>
	</div>
	
	<div class="space-y-4">
		<!-- General Information Component -->
		<GeneralInfo project={selectedProject} {showNotification} />
		
		<!-- Advanced Settings Component -->
		<AdvancedSettings project={selectedProject} {showNotification} />
		
		<!-- Proxy Management Component -->
		<ProxyManagement project={selectedProject} {showNotification} />
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
