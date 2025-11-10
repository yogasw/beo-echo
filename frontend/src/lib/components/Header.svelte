<script lang="ts">
	import { activeTab } from '$lib/stores/activeTab';
	import { theme, toggleTheme } from '$lib/stores/theme';
	import { logStatus } from '$lib/stores/logStatus';
	import SaveButton from './SaveButton.svelte';
	import WorkspaceManager from './workspace/WorkspaceManager.svelte';
	import { currentUser } from '$lib/stores/auth';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import Tooltip from './common/Tooltip.svelte';

	export let handleLogout: () => void;

	// Local state for profile menu toggle
	let profileMenuOpen = false;

	// Tooltip hover states
	let routesHover = false;
	let logsHover = false;
	let replayHover = false;
	let actionsHover = false;
	let configHover = false;

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
</script>

<div class="theme-bg-primary flex items-center p-4">
	<div class="relative mr-4">
		<button
			class="group flex flex-col items-center"
			class:opacity-50={!$selectedProject}
			class:cursor-not-allowed={!$selectedProject}
			on:click={() => $selectedProject && handleTabClick('routes')}
			on:mouseenter={() => routesHover = true}
			on:mouseleave={() => routesHover = false}
			disabled={!$selectedProject}
			aria-label={!$selectedProject ? "Please select project" : "Switch to Routes tab"}
		>
			<div
				class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 flex items-center justify-center"
				class:border-blue-500={$selectedProject}
				class:border-gray-400={!$selectedProject}
				class:bg-blue-500={$activeTab === 'routes' && $selectedProject}
				class:theme-bg-secondary={$activeTab !== 'routes' || !$selectedProject}
			>
				<i class="fas fa-route"></i>
			</div>
			<span class="text-xs mt-1 theme-text-primary">Routes</span>
		</button>
		<Tooltip
			text={!$selectedProject ? "Select a project to access Routes" : ""}
			show={routesHover && !$selectedProject}
			position="bottom"
		/>
	</div>
	<div class="relative mr-4">
		<button
			class="group flex flex-col items-center"
			class:opacity-50={!$selectedProject}
			class:cursor-not-allowed={!$selectedProject}
			on:click={() => {
				if ($selectedProject) {
					handleTabClick('logs');
					logStatus.resetUnread();
				}
			}}
			on:mouseenter={() => logsHover = true}
			on:mouseleave={() => logsHover = false}
			disabled={!$selectedProject}
			aria-label={!$selectedProject ? "Please select project" : "Switch to Logs tab"}
		>
			<div
				class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 flex items-center justify-center relative"
				class:border-yellow-500={$selectedProject}
				class:border-gray-400={!$selectedProject}
				class:bg-blue-500={$activeTab === 'logs' && $selectedProject}
				class:theme-bg-secondary={$activeTab !== 'logs' || !$selectedProject}
			>
				<i class="fas fa-file-alt"></i>

				<!-- Live connection indicator -->
				{#if $selectedProject && $logStatus.isConnected}
					<span class="absolute -top-1 -right-1 flex">
						<span
							class="animate-ping absolute inline-flex h-3 w-3 rounded-full bg-green-400 opacity-75"
						></span>
						<span class="relative inline-flex rounded-full h-3 w-3 bg-green-500"></span>
					</span>
				{:else if $selectedProject}
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
				{#if $selectedProject && $logStatus.unreadCount > 0 && $activeTab !== 'logs'}
					<span
						class="absolute -bottom-1 -right-1 inline-flex items-center justify-center px-2 py-0.5 text-xs font-bold leading-none text-white bg-red-600 rounded-full"
					>
						{$logStatus.unreadCount > 99 ? '99+' : $logStatus.unreadCount}
					</span>
				{/if}
			</div>
			<span class="text-xs mt-1 theme-text-primary">Logs</span>
		</button>
		<Tooltip
			text={!$selectedProject ? "Select a project to view Logs" : ""}
			show={logsHover && !$selectedProject}
			position="bottom"
		/>
	</div>
	<div class="relative mr-4">
		<button
			class="group flex flex-col items-center"
			class:opacity-50={!$selectedProject}
			class:cursor-not-allowed={!$selectedProject}
			on:click={() => $selectedProject && handleTabClick('replay')}
			on:mouseenter={() => replayHover = true}
			on:mouseleave={() => replayHover = false}
			disabled={!$selectedProject}
			aria-label={!$selectedProject ? "Please select project" : "Switch to Replay tab"}
			title="Replay Tab - Execute saved requests"
		>
			<div
				class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 flex items-center justify-center"
				class:border-green-500={$selectedProject}
				class:border-gray-400={!$selectedProject}
				class:bg-blue-500={$activeTab === 'replay' && $selectedProject}
				class:theme-bg-secondary={$activeTab !== 'replay' || !$selectedProject}
			>
				<i class="fas fa-play-circle"></i>
			</div>
			<span class="text-xs mt-1 theme-text-primary">Replay</span>
		</button>
		<Tooltip
			text={!$selectedProject ? "Select a project to manage Replays" : ""}
			show={replayHover && !$selectedProject}
			position="bottom"
		/>
	</div>
	<div class="relative mr-4">
		<button
			class="group flex flex-col items-center"
			class:opacity-50={!$selectedProject}
			class:cursor-not-allowed={!$selectedProject}
			on:click={() => $selectedProject && handleTabClick('actions')}
			on:mouseenter={() => actionsHover = true}
			on:mouseleave={() => actionsHover = false}
			disabled={!$selectedProject}
			aria-label={!$selectedProject ? "Please select project" : "Switch to Actions tab"}
			title="Actions Tab - Modify requests and responses"
		>
			<div
				class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 flex items-center justify-center"
				class:border-amber-500={$selectedProject}
				class:border-gray-400={!$selectedProject}
				class:bg-blue-500={$activeTab === 'actions' && $selectedProject}
				class:theme-bg-secondary={$activeTab !== 'actions' || !$selectedProject}
			>
				<i class="fas fa-bolt"></i>
			</div>
			<span class="text-xs mt-1 theme-text-primary">Actions</span>
		</button>
		<Tooltip
			text={!$selectedProject ? "Select a project to manage Actions" : ""}
			show={actionsHover && !$selectedProject}
			position="bottom"
		/>
	</div>
	<div class="relative mr-auto">
		<button
			class="group flex flex-col items-center"
			class:opacity-50={!$selectedProject}
			class:cursor-not-allowed={!$selectedProject}
			on:click={() => $selectedProject && handleTabClick('configuration')}
			on:mouseenter={() => configHover = true}
			on:mouseleave={() => configHover = false}
			disabled={!$selectedProject}
			aria-label={!$selectedProject ? "Please select project" : "Switch to Configuration tab"}
		>
			<div
				class="w-12 aspect-square theme-text-primary p-3 rounded-full border-2 flex items-center justify-center"
				class:border-purple-500={$selectedProject}
				class:border-gray-400={!$selectedProject}
				class:bg-blue-500={$activeTab === 'configuration' && $selectedProject}
				class:theme-bg-secondary={$activeTab !== 'configuration' || !$selectedProject}
			>
				<i class="fas fa-cogs"></i>
			</div>
			<span class="text-xs mt-1 theme-text-primary">Configuration</span>
		</button>
		<Tooltip
			text={!$selectedProject ? "Select a project to edit Configuration" : ""}
			show={configHover && !$selectedProject}
			position="bottom"
		/>
	</div>

	<!-- Theme Toggle Button -->
	<button class="relative group mr-4 flex flex-col items-center" on:click={toggleTheme}
		title={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}
		aria-label={$theme === 'dark' ? 'Switch to light mode' : 'Switch to dark mode'}>
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
		<button on:click={toggleProfileMenu} class="profile-button flex flex-col items-center"
			title="Open profile menu"
			aria-label="Open profile menu">
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
								<div class="theme-text-primary font-medium">{$currentUser?.name || ""}</div>
								<div class="theme-text-secondary text-xs">{$currentUser?.email || ""}</div>
							</div>
						</div>
					</div>

					<!-- Menu Actions -->
					<div class="border-t theme-border pt-2">
						<button
							class="w-full text-left p-3 rounded-md flex items-center theme-text-primary hover:theme-bg-secondary transition-colors"
							on:click={openSettings}
							title="Open settings"
							aria-label="Open settings"
						>
							<i class="fas fa-cog mr-2 text-blue-400"></i>
							<span>Settings</span>
						</button>
						<button
							class="w-full text-left p-3 rounded-md flex items-center theme-text-primary hover:theme-bg-secondary transition-colors"
							on:click={handleLogout}
							title="Logout"
							aria-label="Logout"
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
	/* Empty style section for future custom styles */
</style>
