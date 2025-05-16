<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import UserProfile from '$lib/components/settings/UserProfile.svelte';
	import PasswordChange from '$lib/components/settings/PasswordChange.svelte';
	
	// State for each section's visibility
	let sectionsVisible = {
		profile: true,    // Profile open by default
		password: false
	};
</script>

<div class="w-full theme-bg-primary p-4">
	<!-- Header -->
	<div class="mb-6">
		<div class="flex items-center mb-4">
			<div class="bg-blue-600/10 dark:bg-blue-600/10 p-2 rounded-lg mr-3">
				<i class="fas fa-user-cog text-blue-500 text-xl"></i>
			</div>
			<div>
				<h2 class="text-xl font-bold theme-text-primary">Account Settings</h2>
				<p class="text-sm theme-text-muted">Manage your personal settings and preferences</p>
			</div>
		</div>
		
		<!-- Info Message -->
		<div class="w-full p-4 mb-6 theme-bg-secondary rounded-lg border theme-border flex items-center gap-3">
			<div class="text-blue-400 text-xl">
				<i class="fas fa-info-circle"></i>
			</div>
			<div>
				<h3 class="theme-text-primary font-medium">Account Settings</h3>
				<p class="theme-text-secondary text-sm mt-1">
					Update your personal information, password, and preferences for your Beo Echo account.
				</p>
			</div>
		</div>
	</div>

	<!-- Settings Sections -->
	<div class="space-y-5">
		<!-- 1. Profile Settings -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div 
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => sectionsVisible.profile = !sectionsVisible.profile}
				on:keydown={(e) => e.key === 'Enter' && (sectionsVisible.profile = !sectionsVisible.profile)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-blue-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-user text-blue-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Profile Information</h3>
				</div>
				<i class="fas {sectionsVisible.profile ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
			</div>
			
			{#if sectionsVisible.profile}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
					<UserProfile />
				</div>
			{/if}
		</div>

		<!-- 2. Password Settings -->
		<div class={ThemeUtils.card('overflow-hidden')}>
			<div 
				class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
				on:click={() => sectionsVisible.password = !sectionsVisible.password}
				on:keydown={(e) => e.key === 'Enter' && (sectionsVisible.password = !sectionsVisible.password)}
				tabindex="0"
				role="button"
			>
				<div class="flex items-center">
					<div class="bg-purple-500/20 p-1.5 rounded mr-2">
						<i class="fas fa-key text-purple-400"></i>
					</div>
					<h3 class="font-medium theme-text-primary">Change Password</h3>
				</div>
				<i class="fas {sectionsVisible.password ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
			</div>
			
			{#if sectionsVisible.password}
				<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
					<PasswordChange />
				</div>
			{/if}
		</div>
		
	</div>
</div>
