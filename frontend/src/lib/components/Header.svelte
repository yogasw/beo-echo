<script lang="ts">
	import { activeTab } from '$lib/stores/activeTab';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { toast } from '$lib/stores/toast';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { logStatus } from '$lib/stores/logStatus';
	import SaveButton from './SaveButton.svelte';
	import WorkspaceManager from './workspace/WorkspaceManager.svelte';

	export let handleLogout: () => void;

	// Local state for profile menu toggle
	let profileMenuOpen = false;

	function handleTabClick(tab: string) {
		$activeTab = tab;
	}

	function toggleProfileMenu() {
		profileMenuOpen = !profileMenuOpen;
	}

	// Handle clicks outside of profile menu to close it
	function handleClickOutside(event: MouseEvent) {
		const profileButton = document.querySelector('.profile-button');
		const profileMenu = document.querySelector('.profile-menu-container');

		// Don't close if clicking the button itself (that's handled by the toggle)
		if (profileButton && profileButton.contains(event.target as Node)) {
			return;
		}

		// Close if clicking outside the menu
		if (profileMenuOpen && profileMenu && !profileMenu.contains(event.target as Node)) {
			profileMenuOpen = false;
		}
	}

	import { onMount, onDestroy } from 'svelte';
	import { logsConnectionStatus } from '$lib/stores/logs';

	onMount(() => {
		if (typeof window !== 'undefined') {
			document.addEventListener('mousedown', handleClickOutside);
		}
	});

	onDestroy(() => {
		if (typeof window !== 'undefined') {
			document.removeEventListener('mousedown', handleClickOutside);
		}
	});

	// Navigate to settings tab
	function openSettings() {
		$activeTab = 'settings';
		profileMenuOpen = false;
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
		on:click={() => {
			handleTabClick('logs');
			logStatus.resetUnread();
		}}
	>
		<div
			class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 border-yellow-500 flex items-center justify-center relative"
			class:bg-blue-500={$activeTab === 'logs'}
			class:theme-bg-secondary={$activeTab !== 'logs'}
		>
			<i class="fas fa-file-alt"></i>

			<!-- Live connection indicator -->
			{#if $logStatus.isConnected}
				<span class="absolute -top-1 -right-1 flex">
					<span
						class="animate-ping absolute inline-flex h-3 w-3 rounded-full bg-green-400 opacity-75"
					></span>
					<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
				</span>
			{:else}
				<!-- Offline Or loading indicator -->
				{#if $logsConnectionStatus.isConnected === false}
					<span class="absolute -top-1 -right-1 flex">
						<span class="animate-ping absolute inline-flex h-3 w-3 rounded-full bg-red-400 opacity-75"></span>
						<span class="relative inline-flex rounded-full h-3 w-3 bg-red-500"></span>
					</span>
				{:else}
					<span class="absolute -top-1 -right-1 flex">
						<span class="animate-ping absolute inline-flex h-3 w-3 rounded-full bg-green-400 opacity-75"></span>
						<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
					</span>
				{/if}
			{/if}

			<!-- Unread logs counter -->
			{#if $logStatus.unreadCount > 0 && $activeTab !== 'logs'}
				<span
					class="absolute -bottom-1 -right-1 inline-flex items-center justify-center px-2 py-0.5 text-xs font-bold leading-none text-white bg-red-600 rounded-full"
				>
					{$logStatus.unreadCount > 99 ? '99+' : $logStatus.unreadCount}
				</span>
			{/if}
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
	<div class="relative flex flex-col items-center">
		<!-- Company/Profile Button -->
		<button on:click={toggleProfileMenu} class="profile-button flex flex-col items-center">
			<div
				class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-gray-500 flex items-center justify-center"
			>
				<i class="fas fa-user-circle"></i>
			</div>
			<span class="text-xs mt-1 theme-text-primary">Profile</span>
		</button>

		<!-- Profile Menu -->
		{#if profileMenuOpen}
			<div
				class="profile-menu-container absolute top-full right-0 mt-2 w-64 theme-bg-primary rounded-md shadow-lg z-40 border theme-border"
			>
				<div class="p-2">
					<div class="max-h-48 overflow-y-auto mb-2">
						<!-- User Info Section -->
						<div class="p-3 mb-2 theme-bg-secondary rounded-md flex items-center">
							<div
								class="w-10 h-10 rounded-full bg-blue-500 flex items-center justify-center text-white mr-3"
							>
								<i class="fas fa-user"></i>
							</div>
							<div class="flex-1">
								<div class="theme-text-primary font-medium">User Name</div>
								<div class="theme-text-secondary text-xs">user@example.com</div>
							</div>
						</div>
					</div>

					<!-- Menu Actions -->
					<div class="border-t theme-border pt-2">
						<button
							class="w-full text-left p-3 rounded-md flex items-center theme-text-primary hover:theme-bg-secondary transition-colors"
							on:click={openSettings}
						>
							<i class="fas fa-cog mr-2 text-blue-400"></i>
							<span>Settings</span>
						</button>
						<button
							class="w-full text-left p-3 rounded-md flex items-center theme-text-primary hover:theme-bg-secondary transition-colors"
							on:click={handleLogout}
						>
							<i class="fas fa-sign-out-alt mr-2 text-red-400"></i>
							<span>Logout</span>
						</button>
					</div>
				</div>
			</div>
		{/if}
	</div>
</div>

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
