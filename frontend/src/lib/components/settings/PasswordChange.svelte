<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	import { featureToggles } from '$lib/stores/featureToggles';
	
	// Form state
	let currentPassword = '';
	let newPassword = '';
	let confirmPassword = '';
	let isSaving = false;
	
	// Password requirements state
	let passwordRequirements = {
		length: false,
		uppercase: false,
		lowercase: false,
		number: false,
		special: false
	};
	
	// Update password requirements check
	$: {
		passwordRequirements = {
			length: newPassword.length >= 8,
			uppercase: /[A-Z]/.test(newPassword),
			lowercase: /[a-z]/.test(newPassword),
			number: /[0-9]/.test(newPassword),
			special: /[^A-Za-z0-9]/.test(newPassword)
		};
	}
	
	// Check if password meets minimum requirements
	// If requirements are disabled, only check basic length requirement
	$: passwordValid = $featureToggles.showPasswordRequirements 
		? (passwordRequirements.length && 
		   passwordRequirements.uppercase && 
		   passwordRequirements.lowercase && 
		   passwordRequirements.number)
		: (newPassword.length >= 8); // Basic minimum length check when requirements are disabled
	
	// Change password function
	async function changePassword() {
		// Validation
		if (!currentPassword) {
			toast.error('Current password is required');
			return;
		}
		
		if (!passwordValid) {
			toast.error('New password does not meet requirements');
			return;
		}
		
		if (newPassword !== confirmPassword) {
			toast.error('New passwords do not match');
			return;
		}
		
		isSaving = true;
		
		try {
			// Simulate API call
			await new Promise(r => setTimeout(r, 800));
			
			// Reset form on success
			currentPassword = '';
			newPassword = '';
			confirmPassword = '';
			
			toast.success('Password changed successfully');
		} catch (error) {
			toast.error('Failed to change password');
			console.error('Password change error:', error);
		} finally {
			isSaving = false;
		}
	}
</script>

<form on:submit|preventDefault={changePassword} class="space-y-6">
	<!-- Current Password -->
	<div>
		<label for="currentPassword" class="block text-sm font-medium theme-text-primary mb-1">Current Password</label>
		<input 
			type="password" 
			id="currentPassword" 
			bind:value={currentPassword}
			class="block w-full p-2 theme-bg-primary theme-text-primary border theme-border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
		/>
	</div>
	
	<!-- New Password -->
	<div>
		<label for="newPassword" class="block text-sm font-medium theme-text-primary mb-1">New Password</label>
		<input 
			type="password" 
			id="newPassword" 
			bind:value={newPassword}
			class="block w-full p-2 theme-bg-primary theme-text-primary border theme-border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"
		/>
	</div>
	
	<!-- Password Requirements (conditionally shown based on feature toggle) -->
	{#if $featureToggles.showPasswordRequirements}
		<div class="p-3 theme-bg-primary border theme-border rounded-md" transition:fade={{ duration: 200 }}>
			<h4 class="text-sm font-medium theme-text-primary mb-2">Password Requirements</h4>
			<ul class="space-y-1 text-sm">
				<li class="flex items-center gap-2">
					<i class="fas fa-{passwordRequirements.length ? 'check text-green-500' : 'times text-gray-400'}"></i>
					<span class="theme-text-secondary">At least 8 characters</span>
				</li>
				<li class="flex items-center gap-2">
					<i class="fas fa-{passwordRequirements.uppercase ? 'check text-green-500' : 'times text-gray-400'}"></i>
					<span class="theme-text-secondary">At least one uppercase letter</span>
				</li>
				<li class="flex items-center gap-2">
					<i class="fas fa-{passwordRequirements.lowercase ? 'check text-green-500' : 'times text-gray-400'}"></i>
					<span class="theme-text-secondary">At least one lowercase letter</span>
				</li>
				<li class="flex items-center gap-2">
					<i class="fas fa-{passwordRequirements.number ? 'check text-green-500' : 'times text-gray-400'}"></i>
					<span class="theme-text-secondary">At least one number</span>
				</li>
				<li class="flex items-center gap-2">
					<i class="fas fa-{passwordRequirements.special ? 'check text-green-500' : 'times text-gray-400'}"></i>
					<span class="theme-text-secondary">Special character (recommended)</span>
				</li>
			</ul>
		</div>
	{/if}
	
	<!-- Confirm Password -->
	<div>
		<label for="confirmPassword" class="block text-sm font-medium theme-text-primary mb-1">Confirm New Password</label>
		<input 
			type="password" 
			id="confirmPassword" 
			bind:value={confirmPassword}
			class={`block w-full p-2 theme-bg-primary theme-text-primary border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 ${newPassword && confirmPassword && newPassword !== confirmPassword ? 'border-red-500' : 'theme-border'}`}
		/>
		{#if newPassword && confirmPassword && newPassword !== confirmPassword}
			<p class="mt-1 text-sm text-red-500">Passwords do not match. The change button will be disabled until passwords match.</p>
		{/if}
	</div>
	
	<!-- Submit Button -->
	<div class="flex justify-end">
		<button 
			type="submit"
			class={`px-4 py-2 rounded-md text-sm flex items-center gap-2 text-white ${isSaving || !currentPassword || !passwordValid || newPassword !== confirmPassword ? 'bg-gray-400 cursor-not-allowed opacity-70' : 'bg-blue-600 hover:bg-blue-700'}`}
			disabled={isSaving || !currentPassword || !passwordValid || newPassword !== confirmPassword}
		>
			{#if isSaving}
				<i class="fas fa-spinner fa-spin"></i>
				<span>Changing Password...</span>
			{:else}
				<i class="fas fa-key"></i>
				<span>Change Password</span>
			{/if}
		</button>
	</div>
</form>
