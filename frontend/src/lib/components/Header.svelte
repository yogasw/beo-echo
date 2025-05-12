<script lang="ts">
	import { activeTab } from '$lib/stores/activeTab';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import Settings from '$lib/components/settings/Settings.svelte';
	import SaveButton from './SaveButton.svelte';
	import WorkspaceManager from './workspace/WorkspaceManager.svelte';

	export let handleLogout: () => void;

	function handleTabClick(tab: string) {
		$activeTab = tab;
	}

	function toggleProfileMenu() {
		const menu = document.getElementById('profileMenu');
		menu?.classList.toggle('hidden');
	}

	function openSettingsModal() {
		const modal = document.getElementById('settingsModal');
		modal?.classList.remove('hidden');
	}

	function closeSettingsModal() {
		const modal = document.getElementById('settingsModal');
		modal?.classList.add('hidden');
	}

	async function handleDownload() {
		if (!$selectedProject) {
			toast.error('Please select a configuration first');
			return;
		}
	}
</script>

<div class="theme-bg-primary flex items-center p-4">

	<button
		class="relative group mr-4 flex flex-col items-center"
		on:click={() => handleTabClick('routes')}
	>
		<div
			class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 border-blue-500 flex items-center justify-center"
			class:bg-blue-500={$activeTab === 'routes'}
			class:theme-bg-secondary={$activeTab !== 'routes'}
		>
			<i class="fas fa-route"></i>
		</div>
		<span class="text-xs mt-1 theme-text-primary">Routes</span>
	</button>
	<button
		class="relative group mr-4 flex flex-col items-center"
		on:click={() => handleTabClick('logs')}
	>
		<div
			class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 border-yellow-500 flex items-center justify-center"
			class:bg-blue-500={$activeTab === 'logs'}
			class:theme-bg-secondary={$activeTab !== 'logs'}
		>
			<i class="fas fa-file-alt"></i>
		</div>
		<span class="text-xs mt-1 theme-text-primary">Logs</span>
	</button>
	<button
		class="relative group mr-auto flex flex-col items-center"
		on:click={() => handleTabClick('configuration')}
	>
		<div
			class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 border-purple-500 flex items-center justify-center"
			class:bg-blue-500={$activeTab === 'configuration'}
			class:theme-bg-secondary={$activeTab !== 'configuration'}
		>
			<i class="fas fa-cogs"></i>
		</div>
		<span class="text-xs mt-1 theme-text-primary">Configuration</span>
	</button>

	<button class="relative group mr-4 flex flex-col items-center" on:click={handleDownload}>
		<div
			class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-blue-500 flex items-center justify-center"
		>
			<i class="fas fa-download"></i>
		</div>
		<span class="text-xs mt-1 theme-text-primary">Download JSON</span>
	</button>

	<!-- Theme Toggle Button -->
	<button class="relative group mr-4 flex flex-col items-center" on:click={toggleTheme}>
		<div
			class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-amber-500 flex items-center justify-center transition-all"
		>
			{#if $theme === 'dark'}
				<i class="fas fa-sun text-amber-400"></i>
			{:else}
				<i class="fas fa-moon text-indigo-600"></i>
			{/if}
		</div>
		<span class="text-xs mt-1 theme-text-primary"
			>{$theme === 'dark' ? 'Light Mode' : 'Dark Mode'}</span
		>
	</button>

	<!-- Workspace Manager Button -->
	<WorkspaceManager className="mr-4" />

	<!-- Profile Button -->
	<div class="relative group flex flex-col items-center">
		<button
			class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-gray-500 flex items-center justify-center"
			on:click={() => {
				toggleProfileMenu();
			}}
			aria-label="Open profile menu"
		>
			<i class="fas fa-user-circle"></i>
		</button>
		<span class="text-xs mt-1 theme-text-primary">Profile</span>
		<!-- Profile Menu -->
		<div
			id="profileMenu"
			class="absolute top-full right-0 theme-bg-secondary theme-text-primary rounded shadow-lg mt-2 hidden w-48"
		>
			<button class="block w-full text-left px-4 py-2 theme-hover" on:click={openSettingsModal}>
				<i class="fas fa-cog mr-2"></i> Settings
			</button>
			<button class="block w-full text-left px-4 py-2 theme-hover" on:click={handleLogout}>
				<i class="fas fa-sign-out-alt mr-2"></i> Logout
			</button>
		</div>
	</div>
</div>

<!-- Settings Modal -->
<Settings />
<SaveButton />

<style>
	.fa-spin {
		animation: fa-spin 1s infinite linear;
	}

	@keyframes fa-spin {
		0% {
			transform: rotate(0deg);
		}
		100% {
			transform: rotate(360deg);
		}
	}
</style>
