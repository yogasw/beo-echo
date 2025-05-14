<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { updateProject, type Project } from '$lib/api/BeoApi';

	export let project: Project;
	export let showNotification: (message: string, type?: 'success' | 'error') => void;

	let isExpanded = false;
	let timeoutEnabled = true;
	let timeout = project?.timeout || 10000;
	let corsEnabled = project?.cors_enabled !== undefined ? project?.cors_enabled : true;
	
	// Function to toggle section expansion
	function toggleSection() {
		isExpanded = !isExpanded;
	}
	
	// Handle save of advanced settings
	async function handleSave() {
		try {
			const updatedProject = await updateProject(project.id, {
				timeout: timeoutEnabled ? timeout : undefined,
				cors_enabled: corsEnabled
			});
			
			// Update local project with new values
			project = updatedProject;
			
			// Show success notification
			showNotification('Advanced settings updated successfully!', 'success');
		} catch (error) {
			console.error('Failed to update advanced settings:', error);
			showNotification('Failed to update advanced settings: ' + (error instanceof Error ? error.message : String(error)), 'error');
		}
	}

	// Handle input or toggle change with auto-save
	function handleChange() {
		handleSave();
	}
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
			<div class="bg-purple-600/10 p-1.5 rounded mr-2">
				<i class="fas fa-sliders-h text-purple-500"></i>
			</div>
			<h3 class="font-medium theme-text-primary">Advanced Settings</h3>
		</div>
		<i class="fas {isExpanded ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
	</div>
	
	{#if isExpanded}
		<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<!-- Timeout Settings with Toggle -->
				<div>
					<div class="flex items-center justify-between mb-2">
						<label for="config-timeout" class="text-sm font-medium theme-text-secondary">Timeout (ms)</label>
						<label class="relative inline-flex items-center cursor-pointer">
							<input 
								type="checkbox" 
								class="sr-only peer" 
								bind:checked={timeoutEnabled}
								on:change={handleChange}
							>
							<div class="w-9 h-5 bg-gray-300 dark:bg-gray-700 peer-focus:outline-none peer-focus:ring-2 
									  peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer 
									  peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full 
									  peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] 
									  after:left-[2px] after:bg-white after:border-gray-300 after:border 
									  after:rounded-full after:h-4 after:w-4 after:transition-all 
									  peer-checked:bg-blue-600">
							</div>
						</label>
					</div>
					<div class="relative">
						<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
							<i class="fas fa-clock theme-text-muted"></i>
						</div>
						<input
							type="number"
							id="config-timeout"
							class={ThemeUtils.inputField(`${!timeoutEnabled ? 'opacity-50' : ''}`)}
							bind:value={timeout}
							on:change={handleChange}
							disabled={!timeoutEnabled}
							placeholder="Response timeout in ms"
						/>
					</div>
				</div>
				
				<!-- CORS Settings -->
				<div>
					<label for="config-cors" class="block text-sm font-medium mb-2 theme-text-secondary">CORS Enabled</label>
					<div class="mt-3">
						<label class="relative inline-flex items-center cursor-pointer">
							<input 
								type="checkbox" 
								class="sr-only peer" 
								bind:checked={corsEnabled}
								on:change={handleChange}
							>
							<div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
									  peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 rounded-full peer 
									  peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full 
									  peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] 
									  after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full 
									  after:h-5 after:w-5 after:transition-all peer-checked:bg-blue-600">
							</div>
							<span class="ms-3 text-sm font-medium theme-text-secondary">Enable CORS</span>
						</label>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
