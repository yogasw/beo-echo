<script lang="ts">
	import RoutesTab from './tabs/RoutesTab.svelte';
	import LogsTab from './tabs/LogsTab.svelte';
	import ReplayTab from './tabs/ReplayTab.svelte';
	import ConfigurationTab from './tabs/ConfigurationTab.svelte';
	import WorkspaceSettingsTab from './tabs/WorkspaceSettingsTab.svelte';
	import InstanceSettingsTab from './tabs/InstanceSettingsTab.svelte';
	import SettingsTab from './tabs/SettingsTab.svelte';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { activeTab } from '$lib/stores/activeTab';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { isLoadingContentArea } from '$lib/stores/loadingContentArea';
	import BeoEchoLoader from './common/BeoEchoLoader.svelte';
	export let activeContentTab = 'Status & Body';
	let error = '';
</script>

<div class={ThemeUtils.themeBgPrimary('content-area')}>
	{#if $isLoadingContentArea}
	<div class="flex items-center justify-center h-full min-h-[300px]">
		<BeoEchoLoader 
			message="Loading project..." 
			size="lg"
			animated={true}
			isLoading={$isLoadingContentArea}
			delay={500}
			minShowTime={300}
		/>
	</div>
	{:else if $activeTab === 'workspace-settings'}
		<!-- Always render workspace settings tab regardless of project selection -->
		<div class="tab-content">
			<WorkspaceSettingsTab />
		</div>
	{:else if $activeTab === 'instance-settings'}
		<!-- Always render instance settings tab regardless of project selection -->
		<div class="tab-content">
			<InstanceSettingsTab />
		</div>
	{:else if $activeTab === 'settings'}
		<!-- Always render settings tab regardless of project selection -->
		<div class="tab-content">
			<SettingsTab />
		</div>
	{:else if !$selectedProject}
		<div class="no-config-message theme-text-primary">
			<i class="fas fa-info-circle text-blue-500"></i>
			<h2 class="theme-text-primary">No Configuration Selected</h2>
			<p class="theme-text-secondary">
				Please select a configuration from the list to view its details.
			</p>
		</div>
	{:else if error}
		<div class="text-red-500 text-center p-4 bg-red-100/10 rounded-md border border-red-500/20">
			{error}
		</div>
	{:else}
		<div class="tab-content">
			{#if $activeTab === 'routes'}
				<RoutesTab {activeContentTab} />
			{:else if $activeTab === 'logs'}
				<LogsTab selectedProject={$selectedProject} />
			{:else if $activeTab === 'replay'}
				<ReplayTab selectedProject={$selectedProject} />
			{:else if $activeTab === 'configuration'}
				<ConfigurationTab selectedProject={$selectedProject} />
			{/if}
		</div>
	{/if}
</div>

<style>
	.content-area {
		display: flex;
		flex-direction: column;
		height: 100%;
		padding: 0 1rem 1rem 1rem;
		transition:
			background-color 0.3s ease,
			color 0.3s ease;
	}

	.tab-content {
		flex: 1;
		overflow-y: auto;
	}

	.no-config-message {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		text-align: center;
	}

	.no-config-message i {
		font-size: 3rem;
		margin-bottom: 1rem;
	}

	.no-config-message h2 {
		font-size: 1.5rem;
		font-weight: 600;
		margin-bottom: 0.5rem;
	}

	.no-config-message p {
		font-size: 1rem;
	}
</style>
