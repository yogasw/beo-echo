<script>
	import { getGitConfig, saveAndTestSyncGit, saveGitConfig } from '$lib/api/BeoApi';

	let gitName = '';
	let gitEmail = '';
	let gitBranch = '';
	let sshKey = '';
	let gitUrl = '';
	let message = '';
	let isLoading = false;

	async function saveConfiguration() {
		isLoading = true;
		message = '';
		try {
			const result = await saveGitConfig({ gitName, gitEmail, gitBranch, sshKey, gitUrl });
			message = result.message;
		} catch (error) {
			message = 'Failed to save configuration.';
		} finally {
			isLoading = false;
		}
	}

	async function saveAndTestSync() {
		isLoading = true;
		message = '';
		try {
			const result = await saveAndTestSyncGit({ gitName, gitEmail, gitBranch, sshKey, gitUrl });
			message = result.message;
		} catch (error) {
			message = 'Failed to save and test sync.';
		} finally {
			isLoading = false;
		}
	}

	// Fetch Git configuration from the backend
	async function fetchGitConfig() {
		isLoading = true;
		message = '';
		try {
			const result = await getGitConfig();
			if (result.success && result.data) {
				gitName = result.data.gitName ?? '';
				gitEmail = result.data.gitEmail ?? '';
				gitBranch = result.data.gitBranch ?? '';
				sshKey = ''; // SSH key is not returned for security reasons
				gitUrl = result.data.gitUrl ?? '';
			} else {
				message = 'Failed to fetch Git configuration.';
			}
		} catch (error) {
			message = 'Error fetching Git configuration.';
		} finally {
			isLoading = false;
		}
	}

	fetchGitConfig()

	// Fetch Git configuration on component mount

</script>

<!-- Git Configuration Tab -->
<div id="gitConfigTab" class="tab-content hidden">
	<h3 class="text-xl font-bold mb-4">Git Configuration</h3>
	<div class="space-y-4">
		<div>
			<label class="block text-sm font-bold mb-2">Name</label>
			<input
				type="text"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="Name commit author"
				bind:value={gitName}
			/>
		</div>
		<div>
			<label class="block text-sm font-bold mb-2">Email</label>
			<input
				type="email"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="Email commit author"
				bind:value={gitEmail}
			/>
		</div>
		<div>
			<label class="block text-sm font-bold mb-2">Default Branch</label>
			<input
				type="text"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="Default branch (e.g., main)"
				bind:value={gitBranch}
			/>
		</div>
		<div>
			<label class="block text-sm font-bold mb-2">Private Key (id_rsa)</label>
			<p class="text-sm text-gray-400 mt-2">The private key will only be visible during the initial input.</p>
			<textarea
				id="privateKey"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded resize-none"
				rows="4"
				placeholder="Paste your private key (e.g., id_rsa), starting with -----BEGIN OPENSSH PRIVATE KEY----- and ending with -----END OPENSSH PRIVATE KEY-----"
				bind:value={sshKey}
			></textarea>
		</div>
		<div>
			<label class="block text-sm font-bold mb-2">SSH URL</label>
			<input
				type="text"
				class="w-full bg-gray-700 text-white py-2 px-4 rounded"
				placeholder="Enter SSH URL (e.g., git@github.com:user/repo.git)"
				bind:value={gitUrl}
			/>
		</div>
	</div>
	<div class="flex space-x-4 mt-4">
		<button
			class="bg-blue-500 text-white py-2 px-4 rounded w-full"
			on:click={saveConfiguration}
			disabled={isLoading}
		>
			{isLoading ? 'Saving...' : 'Save Configuration'}
		</button>
		<button
			class="bg-green-500 text-white py-2 px-4 rounded w-full"
			on:click={saveAndTestSync}
			disabled={isLoading}
		>
			{isLoading ? 'Syncing...' : 'Save and Test Sync'}
		</button>
	</div>
	{#if message}
		<p class="mt-4 text-sm text-gray-300">{message}</p>
	{/if}
</div>
