<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { 
		listProxyTargets, 
		createProxyTarget, 
		updateProxyTarget, 
		deleteProxyTarget,
		updateProjectMode,
		type Project, 
		type ProxyTarget 
	} from '$lib/api/BeoApi';
	import { onMount } from 'svelte';

	export let project: Project;
	export let showNotification: (message: string, type?: 'success' | 'error') => void;

	let isExpanded = false;
	let proxyTargets: ProxyTarget[] = [];
	let isLoading = false;
	let showAddModal = false;
	let showDeleteConfirm = false;
	let selectedProxy: ProxyTarget | null = null;

	// Form states for adding/editing proxy
	let proxyLabel = '';
	let proxyUrl = '';
	let isEditing = false;

	// Initialize data
	onMount(async () => {
		await loadProxyTargets();
	});

	// Load proxy targets from API
	async function loadProxyTargets() {
		isLoading = true;
		try {
			proxyTargets = await listProxyTargets(project.id);
		} catch (error) {
			console.error('Failed to load proxy targets:', error);
			showNotification('Failed to load proxy targets: ' + (error instanceof Error ? error.message : String(error)), 'error');
		} finally {
			isLoading = false;
		}
	}

	// Function to toggle section expansion
	function toggleSection() {
		isExpanded = !isExpanded;
		if (isExpanded && proxyTargets.length === 0) {
			loadProxyTargets();
		}
	}

	// Open add proxy modal
	function showAddProxyModal() {
		resetForm();
		showAddModal = true;
		isEditing = false;
	}

	// Open edit proxy modal
	function showEditProxyModal(proxy: ProxyTarget) {
		selectedProxy = proxy;
		proxyLabel = proxy.label;
		proxyUrl = proxy.url;
		showAddModal = true;
		isEditing = true;
	}

	// Reset form fields
	function resetForm() {
		proxyLabel = '';
		proxyUrl = '';
		selectedProxy = null;
	}

	// Close modal
	function closeModal() {
		showAddModal = false;
		showDeleteConfirm = false;
		resetForm();
	}

	// Show delete confirmation
	function confirmDelete(proxy: ProxyTarget) {
		selectedProxy = proxy;
		showDeleteConfirm = true;
	}

	// Submit new proxy
	async function handleSubmitProxy() {
		try {
			if (isEditing && selectedProxy) {
				// Update existing proxy
				await updateProxyTarget(project.id, selectedProxy.id, {
					label: proxyLabel,
					url: proxyUrl
				});
				showNotification('Proxy target updated successfully!', 'success');
			} else {
				// Create new proxy
				await createProxyTarget(project.id, proxyLabel, proxyUrl);
				showNotification('Proxy target created successfully!', 'success');
			}
			
			// Refresh the list
			await loadProxyTargets();
			closeModal();
		} catch (error) {
			console.error('Failed to save proxy target:', error);
			showNotification('Failed to save proxy target: ' + (error instanceof Error ? error.message : String(error)), 'error');
		}
	}

	// Delete proxy
	async function handleDeleteProxy() {
		if (!selectedProxy) return;
		
		try {
			await deleteProxyTarget(project.id, selectedProxy.id);
			showNotification('Proxy target deleted successfully!', 'success');
			await loadProxyTargets();
			closeModal();
		} catch (error) {
			console.error('Failed to delete proxy target:', error);
			showNotification('Failed to delete proxy target: ' + (error instanceof Error ? error.message : String(error)), 'error');
		}
	}
</script>

<div class={ThemeUtils.card('overflow-hidden')}>
	<div 
		class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
		on:click={toggleSection}
		on:keydown={(e) => e.key === 'Enter' && toggleSection()}
		tabindex="0"
		role="button"
		aria-expanded={isExpanded}
	>
		<div class="flex items-center">
			<div class="bg-green-600/10 p-1.5 rounded mr-2">
				<i class="fas fa-exchange-alt text-green-500"></i>
			</div>
			<h3 class="font-medium theme-text-primary">Proxy Management</h3>
		</div>
		<i class="fas {isExpanded ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
	</div>
	
	{#if isExpanded}
		<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
			<div class="flex justify-between items-center mb-4">
				<p class="theme-text-secondary text-sm">
					{#if project.mode === 'proxy'}
						<span class="bg-blue-600/20 text-blue-500 dark:text-blue-400 px-2 py-1 rounded text-xs font-semibold">
							Proxy Mode Active
						</span>
					{:else}
						<span class="bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 px-2 py-1 rounded text-xs font-semibold">
							Mock Mode Active
						</span>
					{/if}
				</p>
				<button 
					class={ThemeUtils.primaryButton('py-1.5 px-3 rounded-md text-xs')}
					on:click={showAddProxyModal}
				>
					<i class="fas fa-plus mr-2"></i> Add Proxy
				</button>
			</div>
			
			{#if isLoading}
				<div class="flex justify-center items-center p-8">
					<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
				</div>
			{:else if proxyTargets.length === 0}
				<div class="text-center p-8 theme-bg-secondary rounded-lg">
					<i class="fas fa-exchange-alt text-3xl mb-2 text-gray-400 dark:text-gray-500"></i>
					<p class="theme-text-secondary">No proxy targets defined yet.</p>
					<p class="theme-text-muted text-sm mt-1">Click the "Add Proxy" button to create one.</p>
				</div>
			{:else}
				<div class="space-y-3">
					{#each proxyTargets as proxy}
						<div class="theme-border border rounded-lg overflow-hidden">
							<div class="flex justify-between items-center p-3 theme-bg-secondary">
								<div class="flex items-center">
									<div class="bg-green-600/20 p-1.5 rounded-full">
										<i class="fas fa-globe text-green-500 text-sm"></i>
									</div>
									<span class="ml-2 theme-text-primary font-medium">{proxy.label}</span>
								</div>
								<div class="flex items-center space-x-2">
									<button class={ThemeUtils.iconButton('text-xs')} on:click={() => showEditProxyModal(proxy)}>
										<i class="fas fa-edit text-blue-500"></i>
									</button>
									<button class={ThemeUtils.iconButton('text-xs')} on:click={() => confirmDelete(proxy)}>
										<i class="fas fa-trash-alt text-red-500"></i>
									</button>
								</div>
							</div>
							<div class="p-3">
								<div class="flex items-center mb-1">
									<span class="text-xs theme-text-muted mr-2">URL:</span>
									<span class="text-sm theme-text-secondary overflow-ellipsis overflow-hidden">{proxy.url}</span>
								</div>
								<div class="flex items-center text-xs theme-text-muted">
									<i class="fas fa-clock mr-1"></i>
									<span>Created: {new Date(proxy.created_at).toLocaleString()}</span>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<!-- Add/Edit Proxy Modal -->
{#if showAddModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" transition:fade={{ duration: 150 }}>
		<div class={ThemeUtils.card('p-6 rounded-lg max-w-lg w-full')}>
			<h3 class="text-xl font-bold mb-4 theme-text-primary">
				{isEditing ? 'Edit Proxy Target' : 'Add Proxy Target'}
			</h3>
			<form on:submit|preventDefault={handleSubmitProxy}>
				<div class="mb-4">
					<label for="proxy-label" class="block text-sm font-medium mb-2 theme-text-secondary">Label</label>
					<input
						type="text"
						id="proxy-label"
						class={ThemeUtils.inputField()}
						bind:value={proxyLabel}
						placeholder="Production, Staging, etc."
						required
					/>
				</div>
				<div class="mb-6">
					<label for="proxy-url" class="block text-sm font-medium mb-2 theme-text-secondary">Target URL</label>
					<input
						type="url"
						id="proxy-url"
						class={ThemeUtils.inputField()}
						bind:value={proxyUrl}
						placeholder="https://api.example.com"
						required
					/>
				</div>
				<div class="flex justify-end space-x-4">
					<button 
						type="button"
						class={ThemeUtils.secondaryButton('py-2 px-4 rounded')} 
						on:click={closeModal}
					>
						Cancel
					</button>
					<button 
						type="submit"
						class={ThemeUtils.primaryButton('py-2 px-4 rounded')}
					>
						{isEditing ? 'Update' : 'Add'} Proxy
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirm && selectedProxy}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" transition:fade={{ duration: 150 }}>
		<div class={ThemeUtils.card('p-6 rounded-lg max-w-md w-full')}>
			<h3 class="text-xl font-bold mb-4 theme-text-primary flex items-center">
				<i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-2"></i>
				Confirm Delete
			</h3>
			<p class="mb-6 theme-text-secondary">Are you sure you want to delete the proxy target <span class="font-semibold theme-text-primary">"{selectedProxy.name}"</span>? This action cannot be undone.</p>
			<div class="flex justify-end space-x-4">
				<button 
					class={ThemeUtils.secondaryButton('py-2 px-4 rounded')} 
					on:click={closeModal}
				>
					Cancel
				</button>
				<button 
					class={ThemeUtils.destructiveButton('py-2 px-4 rounded')}
					on:click={handleDeleteProxy}
				>
					Delete Proxy
				</button>
			</div>
		</div>
	</div>
{/if}
