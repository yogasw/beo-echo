<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { toast } from '$lib/stores/toast';
	import UserManagement from '$lib/components/instance/UserManagement.svelte';
	import WorkspaceManagement from '$lib/components/instance/WorkspaceManagement.svelte';
	import SecuritySettings from '$lib/components/instance/SecuritySettings.svelte';
	import CustomDomain from '$lib/components/instance/CustomDomain.svelte';
	import SsoIntegration from '$lib/components/instance/SsoIntegration.svelte';
	import GeneralSettings from '$lib/components/instance/GeneralSettings.svelte';
	import FeatureConfigSection from '$lib/components/instance/FeatureConfigSection.svelte';
	import FeatureDebug from '$lib/components/instance/FeatureDebug.svelte';

	// State for each section's visibility
	let sectionsVisible = {
		users: false,
		workspaces: false,
		security: false,
		domain: false,
		sso: false,
		features: false,
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
			toast.success('Instance settings saved successfully');

			// Reset success message after a delay
			setTimeout(() => {
				saveSuccess = false;
			}, 3000);
		}, 1000);
	}
</script>

<div class="w-full theme-bg-primary p-4">
	<!-- Header -->
	<div class="mb-6">
		<div class="flex justify-between items-center mb-4">
			<div class="flex items-center">
				<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
					<i class="fas fa-server text-blue-500 text-xl"></i>
				</div>
				<div>
					<h2 class="text-xl font-bold theme-text-primary">Instance Settings</h2>
					<p class="text-sm theme-text-muted">Manage system-wide configuration</p>
				</div>
			</div>

			<button
				on:click={handleSave}
				class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md flex items-center gap-2 text-sm"
				disabled={isSaving}
			>
				{#if isSaving}
					<i class="fas fa-spinner fa-spin"></i>
					<span>Saving...</span>
				{:else}
					<i class="fas fa-save"></i>
					<span>Save Changes</span>
				{/if}
			</button>
		</div>

		<!-- Info Message -->
		<div
			class="w-full p-4 mb-6 theme-bg-secondary rounded-lg border theme-border flex items-center gap-3"
		>
			<div class="text-blue-400 text-xl">
				<i class="fas fa-info-circle"></i>
			</div>
			<div>
				<h3 class="theme-text-primary font-medium">About Instance Settings</h3>
				<p class="theme-text-secondary text-sm mt-1">
					These settings affect the entire Beo Echo instance. Changes made here will apply to all
					workspaces and users.
				</p>
			</div>
		</div>
	</div>

	<!-- Settings Sections -->
	<div class="space-y-5">
		<!-- 6. Feature Configuration -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.features = !sectionsVisible.features)}
				on:keydown={(e) =>
					e.key === 'Enter' && (sectionsVisible.features = !sectionsVisible.features)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-amber-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-toggle-on text-amber-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Feature Configuration</h3>
				</div>
				<i
					class="fas {sectionsVisible.features
						? 'fa-chevron-up'
						: 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.features}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
					<FeatureConfigSection />

					<!-- Include debug component in development -->
					{#if import.meta.env?.DEV}
						<div class="mt-8 pt-4 border-t theme-border">
							<FeatureDebug />
						</div>
					{/if}
				</div>
			{/if}
		</div>

		<!-- 1. User Management -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.users = !sectionsVisible.users)}
				on:keydown={(e) => e.key === 'Enter' && (sectionsVisible.users = !sectionsVisible.users)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-blue-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-users text-blue-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">User Management</h3>
				</div>
				<i
					class="fas {sectionsVisible.users ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.users}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<UserManagement visible={true} />
				</div>
			{/if}
		</div>

		<!-- 2. Workspace Management -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.workspaces = !sectionsVisible.workspaces)}
				on:keydown={(e) =>
					e.key === 'Enter' && (sectionsVisible.workspaces = !sectionsVisible.workspaces)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-purple-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-building text-purple-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Workspace Management</h3>
				</div>
				<i
					class="fas {sectionsVisible.workspaces
						? 'fa-chevron-up'
						: 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.workspaces}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<WorkspaceManagement visible={true} />
				</div>
			{/if}
		</div>

		<!-- 3. Security Settings -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.security = !sectionsVisible.security)}
				on:keydown={(e) =>
					e.key === 'Enter' && (sectionsVisible.security = !sectionsVisible.security)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-green-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-shield-alt text-green-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Security Settings</h3>
				</div>
				<i
					class="fas {sectionsVisible.security
						? 'fa-chevron-up'
						: 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.security}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<SecuritySettings visible={true} />
				</div>
			{/if}
		</div>

		<!-- 4. Custom Domain -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.domain = !sectionsVisible.domain)}
				on:keydown={(e) => e.key === 'Enter' && (sectionsVisible.domain = !sectionsVisible.domain)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-yellow-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-globe text-yellow-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Custom Domain</h3>
				</div>
				<i
					class="fas {sectionsVisible.domain
						? 'fa-chevron-up'
						: 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.domain}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<CustomDomain visible={true} />
				</div>
			{/if}
		</div>

		<!-- 5. SSO Integration -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.sso = !sectionsVisible.sso)}
				on:keydown={(e) => e.key === 'Enter' && (sectionsVisible.sso = !sectionsVisible.sso)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-red-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-key text-red-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">SSO Integration</h3>
				</div>
				<i class="fas {sectionsVisible.sso ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.sso}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<SsoIntegration visible={true} />
				</div>
			{/if}
		</div>

		<!-- 7. General Settings -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => (sectionsVisible.general = !sectionsVisible.general)}
				on:keydown={(e) =>
					e.key === 'Enter' && (sectionsVisible.general = !sectionsVisible.general)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-blue-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-cog text-blue-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">General Settings</h3>
				</div>
				<i
					class="fas {sectionsVisible.general
						? 'fa-chevron-up'
						: 'fa-chevron-down'} theme-text-muted"
				></i>
			</div>

			{#if sectionsVisible.general}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border">
					<GeneralSettings visible={true} />
				</div>
			{/if}
		</div>
	</div>
</div>
