<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	
	// Notification preferences
	let emailNotifications = true;
	let systemNotifications = true;
	let digestFrequency = 'daily';
	let isSaving = false;
	
	// Save preferences
	async function savePreferences() {
		isSaving = true;
		
		try {
			// Simulate API call
			await new Promise(r => setTimeout(r, 800));
			
			toast.success('Notification preferences saved');
		} catch (error) {
			toast.error('Failed to save preferences');
		} finally {
			isSaving = false;
		}
	}
</script>

<div class="space-y-6">
	<p class="theme-text-secondary text-sm">Choose how and when you'd like to be notified about Beo Echo activity.</p>
	
	<div class="space-y-4">
		<!-- Email Notifications -->
		<div class="flex items-center justify-between py-3 border-b theme-border">
			<div>
				<h4 class="font-medium theme-text-primary">Email Notifications</h4>
				<p class="text-sm theme-text-secondary">Receive updates via email</p>
			</div>
			<label class="inline-flex items-center cursor-pointer">
				<input type="checkbox" bind:checked={emailNotifications} class="sr-only peer">
				<div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
					rounded-full peer peer-focus:outline-none peer-focus:ring-2 
					peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 
					peer-checked:after:translate-x-full peer-checked:after:border-white 
					after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
					after:bg-white after:rounded-full after:h-5 after:w-5 
					after:transition-all dark:border-gray-600 relative">
				</div>
			</label>
		</div>
		
		<!-- System Notifications -->
		<div class="flex items-center justify-between py-3 border-b theme-border">
			<div>
				<h4 class="font-medium theme-text-primary">In-App Notifications</h4>
				<p class="text-sm theme-text-secondary">Receive updates in the app</p>
			</div>
			<label class="inline-flex items-center cursor-pointer">
				<input type="checkbox" bind:checked={systemNotifications} class="sr-only peer">
				<div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
					rounded-full peer peer-focus:outline-none peer-focus:ring-2 
					peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 
					peer-checked:after:translate-x-full peer-checked:after:border-white 
					after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
					after:bg-white after:rounded-full after:h-5 after:w-5 
					after:transition-all dark:border-gray-600 relative">
				</div>
			</label>
		</div>
		
		<!-- Digest Frequency -->
		<div class="py-3 border-b theme-border">
			<h4 class="font-medium theme-text-primary mb-2">Summary Digest</h4>
			<p class="text-sm theme-text-secondary mb-3">Frequency of activity summary emails</p>
			<div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
				<label class="inline-flex items-center p-3 theme-bg-primary border theme-border rounded-md cursor-pointer">
					<input type="radio" bind:group={digestFrequency} value="daily" class="mr-2">
					<span class="theme-text-primary">Daily</span>
				</label>
				<label class="inline-flex items-center p-3 theme-bg-primary border theme-border rounded-md cursor-pointer">
					<input type="radio" bind:group={digestFrequency} value="weekly" class="mr-2">
					<span class="theme-text-primary">Weekly</span>
				</label>
				<label class="inline-flex items-center p-3 theme-bg-primary border theme-border rounded-md cursor-pointer">
					<input type="radio" bind:group={digestFrequency} value="never" class="mr-2">
					<span class="theme-text-primary">Never</span>
				</label>
			</div>
		</div>
	</div>
	
	<!-- Save Button -->
	<div class="flex justify-end">
		<button 
			on:click={savePreferences}
			class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm flex items-center gap-2"
			disabled={isSaving}
		>
			{#if isSaving}
				<i class="fas fa-spinner fa-spin"></i>
				<span>Saving...</span>
			{:else}
				<i class="fas fa-save"></i>
				<span>Save Preferences</span>
			{/if}
		</button>
	</div>
</div>
