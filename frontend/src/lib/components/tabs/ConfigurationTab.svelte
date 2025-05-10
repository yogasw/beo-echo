<script lang="ts">
	import { deleteConfig } from '$lib/api/mockoonApi';
	import { goto } from '$app/navigation';
	import { projects } from '$lib/stores/configurations';
	import { selectedProject as selectedConfigStore } from '$lib/stores/selectedConfig';

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

	export let selectedConfig: Config;

	// Dummy data for configuration
	let config = {
		name: selectedConfig.name,
		baseUrl: 'http://localhost',
		port: selectedConfig.port,
		configFile: selectedConfig.configFile,
		cors: true,
		proxyMode: false,
		proxyHost: '',
		proxyPort: 0,
		proxyPrefix: '',
		https: false,
		tlsOptions: {
			enabled: false,
			type: 'CERT',
			pfxPath: '',
			certPath: '',
			keyPath: '',
			caPath: '',
			passphrase: ''
		},
		headers: [
			{
				key: 'Content-Type',
				value: 'application/json'
			},
			{
				key: 'X-Powered-By',
				value: 'Mockoon'
			}
		]
	};

	let showDeleteConfirm = false;

	function handleSave() {
		// Here you would typically save the configuration
		console.log('Saving configuration:', config);
	}

	async function handleDelete() {
		try {
			await deleteConfig(selectedConfig.configFile);
			// Update configurations store
			projects.update(configs => configs.filter(c => c.id !== selectedConfig.uuid));
			// Set selectedConfig to null
			selectedConfigStore.set(null);
			// Redirect to home
			await goto('/home');
		} catch (error) {
			console.error('Failed to delete configuration:', error);
		}
	}
</script>

<div class="w-full bg-gray-800 p-4">
	<div class="max-w-2xl mx-auto">
		<div class="bg-gray-700 p-4 rounded mb-4 flex items-center">
			<i class="fas fa-info-circle text-blue-500 text-2xl mr-2"></i>
			<span class="text-xl font-bold text-blue-500">
        Configuration: {selectedConfig.name}
      </span>
		</div>
		<div class="space-y-4">
			<div>
				<label for="config-name" class="block text-sm font-medium mb-2">Name</label>
				<input
					type="text"
					id="config-name"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={config.name}
				/>
			</div>
			<div>
				<label for="config-url" class="block text-sm font-medium mb-2">Base URL</label>
				<input
					type="text"
					id="config-url"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={config.baseUrl}
				/>
			</div>
			<div>
				<label for="config-port" class="block text-sm font-medium mb-2">Port</label>
				<input
					type="number"
					id="config-port"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={config.port}
				/>
			</div>
			<div>
				<label for="config-file" class="block text-sm font-medium mb-2">JSON File</label>
				<input
					type="text"
					id="config-file"
					class="w-full bg-gray-700 text-white p-2 rounded"
					bind:value={config.configFile}
				/>
			</div>
			<div class="flex items-center space-x-2">
				<input
					type="checkbox"
					id="config-cors"
					class="bg-gray-700 text-blue-500 rounded"
					bind:checked={config.cors}
				/>
				<label for="config-cors" class="text-sm font-medium">Enable CORS</label>
			</div>
			<div class="flex items-center space-x-2">
				<input
					type="checkbox"
					id="config-proxy"
					class="bg-gray-700 text-blue-500 rounded"
					bind:checked={config.proxyMode}
				/>
				<label for="config-proxy" class="text-sm font-medium">Enable Proxy Mode</label>
			</div>
			<div class="flex items-center space-x-2">
				<input
					type="checkbox"
					id="config-https"
					class="bg-gray-700 text-blue-500 rounded"
					bind:checked={config.https}
				/>
				<label for="config-https" class="text-sm font-medium">Enable HTTPS</label>
			</div>
			<button class="bg-blue-500 text-white py-2 px-4 rounded" on:click={handleSave}>
				Save Configuration
			</button>
			<button class="bg-red-500 text-white py-2 px-4 rounded" on:click={() => showDeleteConfirm = true}>
				Delete Configuration
			</button>
		</div>
	</div>
</div>

{#if showDeleteConfirm}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center">
		<div class="bg-gray-800 p-6 rounded-lg max-w-md w-full">
			<h3 class="text-xl font-bold mb-4">Confirm Delete</h3>
			<p class="mb-4">Are you sure you want to delete this configuration? This action cannot be undone.</p>
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
