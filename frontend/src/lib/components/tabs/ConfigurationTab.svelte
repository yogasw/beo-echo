<script lang="ts">
	import { deleteProject, type Project } from '$lib/api/mockoonApi';
	import { goto } from '$app/navigation';
	import { projects } from '$lib/stores/configurations';
	import { selectedProject as selectedProjectStore } from '$lib/stores/selectedConfig';

	export let selectedProject: Project;

	let showDeleteConfirm = false;
	
	function handleSave() {
		// Here you would typically save the configuration
		console.log('Saving configuration:', selectedProject);
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
		} catch (error) {
			console.error('Failed to delete project:', error);
		}
	}
</script>

<div class="w-full bg-gray-800 p-4">
	<div class="max-w-2xl mx-auto">
		<div class="bg-gray-700 p-4 rounded mb-4 flex items-center">
			<i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
			<span class="text-xl font-bold text-blue-500">
        Project: {selectedProject.name}
      </span>
		</div>
		<div class="space-y-4">
			<div>
				<label for="config-name" class="block text-sm font-medium mb-2">Name</label>
				<input
					type="text"
					id="config-name"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={selectedProject.name}
				/>
			</div>
			<div>
				<label for="config-url" class="block text-sm font-medium mb-2">Base URL</label>
				<div class="w-full bg-gray-700 text-white p-2 rounded flex items-center">
					<span class="text-gray-300">{selectedProject.url}</span>
					<span class="ml-2 text-xs bg-gray-600 px-1.5 py-0.5 rounded text-gray-300">Read only</span>
				</div>
			</div>

			<div>
				<label for="config-port" class="block text-sm font-medium mb-2">Alias</label>
				<input
					type="text"
					id="config-port"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={selectedProject.alias}
				/>
			</div>

			<button class="bg-blue-500 text-white py-2 px-4 rounded" on:click={handleSave}>
				Save Project
			</button>
			<button class="bg-red-500 text-white py-2 px-4 rounded" on:click={() => showDeleteConfirm = true}>
				Delete Project
			</button>
		</div>
	</div>
</div>

{#if showDeleteConfirm}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
		<div class="bg-gray-800 p-6 rounded-lg max-w-md w-full">
			<h3 class="text-xl font-bold mb-4">Confirm Delete</h3>
			<p class="mb-4">Are you sure you want to delete this project? This action cannot be undone.</p>
			<div class="flex justify-end space-x-4">
				<button class="bg-gray-700 text-white py-2 px-4 rounded" on:click={() => showDeleteConfirm = false}>
					Cancel
				</button>
				<button class="bg-red-500 text-white py-2 px-4 rounded" on:click={handleDelete}>
					Delete
				</button>
			</div>
		</div>
	</div>
{/if}
