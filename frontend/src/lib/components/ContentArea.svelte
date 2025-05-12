<script lang="ts">
	import RoutesTab from './tabs/RoutesTab.svelte';
	import LogsTab from './tabs/LogsTab.svelte';
	import ConfigurationTab from './tabs/ConfigurationTab.svelte';
	import ProjectStatusBadge from './ProjectStatusBadge.svelte';
	import { getProjectDetail, updateProjectStatus, type Endpoint } from '$lib/api/mockoonApi';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import { activeTab } from '$lib/stores/activeTab';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';

	export let endpoints: Endpoint[] = [];
	export let activeContentTab = 'Status & Body';

	let loading = false;
	let error = '';
	let updatingStatus = false;

	async function loadConfigData() {
		if (!$selectedProject) return;

		loading = true;
		try {
			// Download config using the configFile name
			const response = await getProjectDetail($selectedProject.id);

			// Parse routes from config
			endpoints = response.endpoints;
		} catch (err) {
			console.error('Failed to load config data:', err);
			error = 'Failed to load configuration data';
		} finally {
			loading = false;
		}
	}

	async function handleUpdateStatus(newStatus: string) {
		if (!$selectedProject || updatingStatus) return;
		
		updatingStatus = true;
		try {
			await updateProjectStatus($selectedProject.id, newStatus);
			// Update local project status
			selectedProject.update(current => {
				if (current) {
					return { ...current, status: newStatus };
				}
				return current;
			});
			
			toast.success(`Project ${newStatus === 'running' ? 'started' : 'stopped'} successfully`);
		} catch (err) {
			toast.error(`Failed to ${newStatus === 'running' ? 'start' : 'stop'} project`);
			console.error('Error updating project status:', err);
		} finally {
			updatingStatus = false;
		}
	}

	// Watch for selectedConfig changes
	$: if ($selectedProject) {
		loadConfigData();
	}
</script>

<div class={ThemeUtils.themeBgPrimary("content-area")}>
	{#if !$selectedProject}
		<div class="no-config-message theme-text-primary">
			<i class="fas fa-info-circle text-blue-500"></i>
			<h2 class="theme-text-primary">No Configuration Selected</h2>
			<p class="theme-text-secondary">Please select a configuration from the list to view its details.</p>
		</div>
	{:else if loading}
		<div class="flex items-center justify-center h-full min-h-[300px]">
			<div class="animate-spin rounded-full h-16 w-16 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if error}
		<div class="text-red-500 text-center p-4 bg-red-100/10 rounded-md border border-red-500/20">{error}</div>
	{:else}
		<div class="tab-content">
			{#if $activeTab === 'routes'}
				<RoutesTab selectedProject={$selectedProject} {endpoints} activeContentTab={activeContentTab} />
			{:else if $activeTab === 'logs'}
				<LogsTab selectedProject={$selectedProject} />
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
        padding: 1rem;
        transition: background-color 0.3s ease, color 0.3s ease;
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
