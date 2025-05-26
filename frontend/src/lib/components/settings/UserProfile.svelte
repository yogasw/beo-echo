<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	import { currentUser } from '$lib/stores/auth';
	import { onMount } from 'svelte';
	import { updateUserProfile } from '$lib/api/BeoApi';
	import { FeatureFlags, getFeatureToggle } from '$lib/stores/featureToggles';
	
	// Feature flags
	let emailUpdatesEnabled = getFeatureToggle(FeatureFlags.FEATURE_EMAIL_UPDATES_ENABLED);

	// User profile state
	$: fullName = $currentUser?.name || 'User Name';
	$: email = $currentUser?.email || 'user@example.com';
	$: userId = $currentUser?.id || '';
	$: accountEnabled = $currentUser?.isEnabled !== false; // Default to true if undefined
	
	let isEditing = false;
	let isSaving = false;
	// Form state for editing
	let formName = fullName;
	let formEmail = email;
	
	// Load feature flags on component mount
	onMount(() => {
		// Also initialize form values
		formName = fullName;
		formEmail = email;
	});
	
	// Toggle edit mode
	function toggleEdit() {
		if (isEditing) {
			// Reset form if canceling
			formName = fullName;
			formEmail = email;
		} else {
			// Initialize form values when entering edit mode
			formName = fullName;
			formEmail = email;
		}
		isEditing = !isEditing;
	}
	
	// Save profile changes
	async function saveProfile() {
		if (!formName.trim()) {
			toast.error('Name cannot be empty');
			return;
		}
		
		isSaving = true;
		
		try {
			// Create payload with updated user data
			const payload: {
				name?: string;
				email?: string;
			} = {
				name: formName
			};
			
			// Add email to payload only if email updates are enabled and changed
			if (emailUpdatesEnabled && formEmail !== email) {
				payload.email = formEmail;
			}
			
			// Call the API to update the user profile
			await updateUserProfile(userId, payload);
			
			isEditing = false;
			toast.success('Profile updated successfully');
		} catch (error) {
			toast.error('Failed to update profile');
			console.error('Profile update error:', error);
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="space-y-6">
	<!-- User Profile Section -->
	<div class="flex flex-col sm:flex-row items-start sm:items-center gap-4 p-4 theme-bg-primary rounded-lg mb-6">
		<div class="w-20 h-20 rounded-full bg-blue-500 flex items-center justify-center text-white text-2xl">
			<i class="fas fa-user"></i>
		</div>
		<div class="flex-1">
			<h3 class="text-lg font-semibold theme-text-primary">{fullName}</h3>
			<p class="theme-text-secondary text-sm">{email}</p>
			<div class="flex items-center gap-2 mt-2">
				{#if accountEnabled}
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100">
						<span class="w-2 h-2 mr-1 bg-green-400 dark:bg-green-300 rounded-full"></span>
						Active Account
					</span>
				{:else}
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 dark:bg-red-800 text-red-800 dark:text-red-100">
						<span class="w-2 h-2 mr-1 bg-red-400 dark:bg-red-300 rounded-full"></span>
						Disabled Account
					</span>
				{/if}
				
				{#if $currentUser?.isOwner}
					<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 dark:bg-purple-800 text-purple-800 dark:text-purple-100">
						Owner
					</span>
				{/if}
			</div>
		</div>
	</div>
	
	{#if isEditing}
		<!-- Edit Profile Form -->
		<form on:submit|preventDefault={saveProfile} class="space-y-4">
			<div class="space-y-2">
				<label for="fullName" class="block text-sm font-medium theme-text-primary">Full Name</label>
				<input 
					type="text" 
					id="fullName" 
					bind:value={formName}
					class="block w-full p-2 theme-bg-primary theme-text-primary border theme-border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
				/>
			</div>
			
			<div class="space-y-2">
				<label class="block text-sm font-medium theme-text-primary">Email</label>
				<div class="flex items-center gap-2">
					<input 
						type="text" 
						bind:value={formEmail}
						disabled={!emailUpdatesEnabled}
						class="block w-full p-2 theme-bg-primary theme-text-secondary border theme-border rounded-md" class:opacity-75={!emailUpdatesEnabled}
					/>
					<!-- Feature flag indicator for email updates -->
					{#if !emailUpdatesEnabled}
						<span class="inline-flex items-center px-2 py-1 rounded text-xs font-medium bg-yellow-100 dark:bg-yellow-800 text-yellow-800 dark:text-yellow-100">
							<i class="fas fa-lock mr-1"></i> Disabled
						</span>
					{/if}
				</div>
				<p class="text-xs theme-text-muted">
					{emailUpdatesEnabled ? 
						'You can update your email address' : 
						'Email updates are currently disabled by system administrator'}
				</p>
			</div>
			
			<div class="flex justify-end gap-3 pt-3">
				<button 
					type="button"
					on:click={toggleEdit}
					class="theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
					title="Cancel profile editing changes"
					aria-label="Cancel profile editing changes"
				>
					Cancel
				</button>
				<button 
					type="submit"
					class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center gap-2"
					disabled={isSaving}
					title="Save profile changes"
					aria-label="Save profile changes"
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
		</form>
	{:else}
		<!-- View Profile -->
		<div class="flex justify-end">
			<button 
				on:click={toggleEdit}
				class="theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm flex items-center gap-2 hover:bg-gray-200 dark:hover:bg-gray-600"
				title="Edit user profile information"
				aria-label="Edit user profile information"
			>
				<i class="fas fa-pencil-alt"></i>
				<span>Edit Profile</span>
			</button>
		</div>
	{/if}
</div>
