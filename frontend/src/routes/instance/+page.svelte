<script lang="ts">
	import { onMount } from 'svelte';
	import { fade } from 'svelte/transition';
	import { goto } from '$app/navigation';
	import SettingsSection from '$lib/components/instance/SettingsSection.svelte';
	import UserManagement from '$lib/components/instance/UserManagement.svelte';
	import WorkspaceManagement from '$lib/components/instance/WorkspaceManagement.svelte';
	import SecuritySettings from '$lib/components/instance/SecuritySettings.svelte';
	import CustomDomain from '$lib/components/instance/CustomDomain.svelte';
	import SsoIntegration from '$lib/components/instance/SsoIntegration.svelte';
	import GeneralSettings from '$lib/components/instance/GeneralSettings.svelte';
	import { currentWorkspace } from '$lib/stores/workspace';
	
	// State for each section's visibility
	let sectionsVisible = {
		users: false,
		workspaces: false,
		security: false,
		domain: false,
		sso: false,
		general: true // General open by default
	};
	
	// State for save action
	let isSaving = false;
	let saveSuccess = false;
	
	// Handle save
	function handleSave() {
		isSaving = true;
		
		// Simulate API call
		setTimeout(() => {
			isSaving = false;
			saveSuccess = true;
			
			// Reset success message after a delay
			setTimeout(() => {
				saveSuccess = false;
			}, 3000);
		}, 1000);
	}
</script>

<div class="w-full min-h-screen theme-bg-primary">
	<!-- Header -->
	<div class="w-full theme-bg-secondary p-4 shadow-md">
		<div class="container mx-auto">
			<div class="flex items-center justify-between">
				<div class="flex items-center gap-4">
					<a href="/" class="theme-text-primary hover:text-blue-500 transition-colors">
						<i class="fas fa-arrow-left"></i>
					</a>
					<h1 class="text-xl font-bold theme-text-primary">Instance Settings</h1>
				</div>
				<div>
					<button 
						on:click={handleSave}
						class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center gap-2 relative"
						disabled={isSaving}
					>
						{#if isSaving}
							<i class="fas fa-spinner fa-spin"></i>
							<span>Saving...</span>
						{:else}
							<i class="fas fa-save"></i>
							<span>Save Changes</span>
						{/if}
						
						{#if saveSuccess}
							<span class="absolute top-full right-0 mt-1 px-2 py-1 bg-green-600 text-white text-xs rounded animate-fade-in-out" transition:fade>
								Saved successfully
							</span>
						{/if}
					</button>
				</div>
			</div>
		</div>
	</div>

	<!-- Content -->
	<div class="container mx-auto p-4 max-w-5xl">
		<!-- Info Message -->
		<div class="w-full p-4 mb-6 theme-bg-secondary rounded-lg border theme-border flex items-center gap-3">
			<div class="text-blue-400 text-xl">
				<i class="fas fa-info-circle"></i>
			</div>
			<div>
				<h3 class="theme-text-primary font-medium">About Instance Settings</h3>
				<p class="theme-text-secondary text-sm mt-1">
					These settings affect the entire Beo Echo instance. Changes made here will apply to all workspaces and users.
				</p>
			</div>
		</div>

		<!-- Settings Sections -->
		<div class="space-y-5">
			<!-- 1. User Management -->
			<SettingsSection title="User Management" icon="fa-users" iconBgColorClass="bg-blue-500/20" iconTextColorClass="text-blue-400" bind:open={sectionsVisible.users}>
				<UserManagement visible={sectionsVisible.users} />
			</SettingsSection>

			<!-- 2. Workspace Management -->
			<SettingsSection title="Workspace Management" icon="fa-building" iconBgColorClass="bg-purple-500/20" iconTextColorClass="text-purple-400" bind:open={sectionsVisible.workspaces}>
				<WorkspaceManagement visible={sectionsVisible.workspaces} />
			</SettingsSection>

			<!-- 3. Security Settings -->
			<SettingsSection title="Security Settings" icon="fa-shield-alt" iconBgColorClass="bg-green-500/20" iconTextColorClass="text-green-400" bind:open={sectionsVisible.security}>
				<SecuritySettings visible={sectionsVisible.security} />
			</SettingsSection>

			<!-- 4. Custom Domain -->
			<SettingsSection title="Custom Domain" icon="fa-globe" iconBgColorClass="bg-yellow-500/20" iconTextColorClass="text-yellow-400" bind:open={sectionsVisible.domain}>
				<CustomDomain visible={sectionsVisible.domain} />
			</SettingsSection>

			<!-- 5. SSO Integration -->
			<SettingsSection title="SSO Integration" icon="fa-key" iconBgColorClass="bg-red-500/20" iconTextColorClass="text-red-400" bind:open={sectionsVisible.sso}>
				<SsoIntegration visible={sectionsVisible.sso} />
			</SettingsSection>

			<!-- 6. General Settings -->
			<SettingsSection title="General Settings" icon="fa-cog" iconBgColorClass="bg-blue-500/20" iconTextColorClass="text-blue-400" bind:open={sectionsVisible.general}>
				<GeneralSettings visible={sectionsVisible.general} />
			</SettingsSection>
		</div>
	</div>
</div>
