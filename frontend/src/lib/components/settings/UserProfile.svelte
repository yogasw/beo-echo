<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	
	// User profile state
	let fullName = 'User Name';
	let email = 'user@example.com';
	let isEditing = false;
	let isSaving = false;
	
	// Form state for editing
	let formName = fullName;
	
	// Toggle edit mode
	function toggleEdit() {
		if (isEditing) {
			// Reset form if canceling
			formName = fullName;
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
			// Simulate API call
			await new Promise(r => setTimeout(r, 800));
			
			// Update local state on success
			fullName = formName;
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
				<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 dark:bg-green-800 text-green-800 dark:text-green-100">
					<span class="w-2 h-2 mr-1 bg-green-400 dark:bg-green-300 rounded-full"></span>
					Active Account
				</span>
				<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 dark:bg-blue-800 text-blue-800 dark:text-blue-100">
					Admin
				</span>
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
				<input 
					type="text" 
					value={email}
					disabled
					class="block w-full p-2 theme-bg-primary theme-text-secondary border theme-border rounded-md opacity-75"
				/>
				<p class="text-xs theme-text-muted">Email address cannot be changed</p>
			</div>
			
			<div class="flex justify-end gap-3 pt-3">
				<button 
					type="button"
					on:click={toggleEdit}
					class="theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
				>
					Cancel
				</button>
				<button 
					type="submit"
					class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center gap-2"
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
		</form>
	{:else}
		<!-- View Profile -->
		<div class="flex justify-end">
			<button 
				on:click={toggleEdit}
				class="theme-bg-secondary theme-text-primary px-4 py-2 rounded-md text-sm flex items-center gap-2 hover:bg-gray-200 dark:hover:bg-gray-600"
			>
				<i class="fas fa-pencil-alt"></i>
				<span>Edit Profile</span>
			</button>
		</div>
	{/if}
</div>
