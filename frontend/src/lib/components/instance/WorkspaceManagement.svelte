<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { currentUser } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import type { Workspace } from '$lib/types/Workspace';
	
	// API imports
	import { workspaceApi } from '$lib/api/workspaceApi';
	import { toast } from '$lib/stores/toast';
	
	export let visible = false;
	
	let workspaces: Workspace[] = [];
	let isLoading = false;
	let error: string | null = null;
	let showAutoInviteModal = false;
	let selectedWorkspace: { id: string, name: string } | null = null;
	
	onMount(async () => {
		if (visible) {
			await loadWorkspaces();
		}
	});
	
	$: if (visible) {
		loadWorkspaces();
	}
	
	async function loadWorkspaces() {
		isLoading = true;
		error = null;
		
		try {
			// Assuming there's an API endpoint to fetch all workspaces for system owners
			const response = await workspaceApi.getAllWorkspaces();
			workspaces = response;
		} catch (err) {
			console.error('Failed to load workspaces:', err);
			error = 'Failed to load workspaces. Please try again.';
		} finally {
			isLoading = false;
		}
	}
	
	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString();
	}
	
	function openAutoInviteModal(event: CustomEvent<{ workspaceId: string, workspaceName: string }>) {
		selectedWorkspace = { 
			id: event.detail.workspaceId, 
			name: event.detail.workspaceName 
		};
		showAutoInviteModal = true;
	}
	
	function closeAutoInviteModal() {
		showAutoInviteModal = false;
		// Refresh data after changes
		loadWorkspaces();
	}
	
	async function toggleAutoInvite(workspace: Workspace) {
		try {
			const newStatus = !workspace.auto_invite_enabled;
			await workspaceApi.toggleAutoInvite(workspace.id, newStatus);
			
			// Update local state
			workspace.auto_invite_enabled = newStatus;
			
			// Show notification
			toast.info(`Auto-invite ${newStatus ? 'enabled' : 'disabled'} for "${workspace.name}"`);
			
		} catch (err) {
			console.error('Failed to toggle auto-invite:', err);
			toast.error(err)
		}
	}
	
	async function deleteWorkspace(id: string, name: string) {
		if (!confirm(`Are you sure you want to delete the workspace "${name}"? This action cannot be undone.`)) {
			return;
		}
		
		try {
			await workspaceApi.deleteWorkspace(id);
			toast.info("Workspace deleted successfully");
			await loadWorkspaces();
		} catch (err) {
			console.error('Failed to delete workspace:', err);
			toast.error(err)
		}
	}
	
	function editWorkspace(id: string) {
		goto(`/workspaces/${id}/edit`);
	}
	
	function manageMembers(id: string) {
		goto(`/workspaces/${id}/members`);
	}
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<h3 class="theme-text-primary font-medium mb-3">All Workspaces</h3>
		<div class="overflow-x-auto mb-4">
			{#if isLoading}
				<div class="flex justify-center py-8">
					<div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
				</div>
			{:else if error}
				<div class="p-4 mb-4 text-sm text-red-700 bg-red-100 dark:bg-red-900/30 dark:text-red-300 rounded-lg">
					<i class="fas fa-exclamation-circle mr-2"></i> {error}
				</div>
			{:else if workspaces.length === 0}
				<div class="p-4 mb-4 text-sm theme-text-secondary bg-gray-100 dark:bg-gray-800 rounded-lg">
					<i class="fas fa-info-circle mr-2"></i> No workspaces found.
				</div>
			{:else}
				<table class="w-full text-sm text-left">
					<thead class="text-xs uppercase theme-text-secondary">
						<tr>
							<th scope="col" class="px-4 py-3">Workspace</th>
							<th scope="col" class="px-4 py-3">Members</th>
							<th scope="col" class="px-4 py-3">Projects</th>
							<th scope="col" class="px-4 py-3">Created</th>
							<th scope="col" class="px-4 py-3">Auto-Invite</th>
							<th scope="col" class="px-4 py-3">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each workspaces as workspace}
							<tr class="theme-border-subtle border-b">
								<td class="px-4 py-3 theme-text-primary">{workspace.name}</td>
								<td class="px-4 py-3 theme-text-secondary">
									{workspace.members?.length || 0}
								</td>
								<td class="px-4 py-3 theme-text-secondary">
									{workspace.projects?.length || 0}
								</td>
								<td class="px-4 py-3 theme-text-secondary">
									{formatDate(workspace.created_at)}
								</td>
								<td class="px-4 py-3 theme-text-secondary">
									<div class="flex items-center gap-2">
										{#if workspace.auto_invite_enabled}
											<span class="bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400 py-1 px-2 rounded text-xs">
												Enabled
											</span>
										{:else}
											<span class="bg-gray-100 text-gray-600 dark:bg-gray-800 dark:text-gray-400 py-1 px-2 rounded text-xs">
												Disabled
											</span>
										{/if}
										
										{#if $currentUser?.is_owner}
											<button 
												on:click={() => toggleAutoInvite(workspace)} 
												class="ml-2 text-xs bg-gray-100 hover:bg-gray-200 dark:bg-gray-700 dark:hover:bg-gray-600 py-1 px-2 rounded"
												title={workspace.auto_invite_enabled ? "Disable Auto-Invite" : "Enable Auto-Invite"}
											>
												{workspace.auto_invite_enabled ? 'Disable' : 'Enable'}
											</button>
										{/if}
									</div>
								</td>
								<td class="px-4 py-3">
									<div class="flex items-center gap-2">
										<button 
											on:click={() => manageMembers(workspace.id)} 
											class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" 
											title="Manage Members"
										>
											<i class="fas fa-users theme-text-secondary"></i>
										</button>
										
										<button 
											on:click={() => editWorkspace(workspace.id)} 
											class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" 
											title="Edit"
										>
											<i class="fas fa-edit theme-text-secondary"></i>
										</button>
										
										<button 
											on:click={() => deleteWorkspace(workspace.id, workspace.name)} 
											class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20" 
											title="Delete"
										>
											<i class="fas fa-trash theme-text-secondary"></i>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
		
		<!-- Workspace Actions -->
		<div class="flex gap-3">
			<button 
				on:click={() => goto('/workspaces/create')}
				class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center gap-2"
			>
				<i class="fas fa-plus"></i>
				<span>Create Workspace</span>
			</button>
			
			<button 
				on:click={loadWorkspaces} 
				class="px-3 py-2 theme-bg-secondary hover:theme-bg-tertiary theme-text-primary rounded-md text-sm flex items-center gap-2"
			>
				<i class="fas fa-sync-alt"></i>
				<span>Refresh List</span>
			</button>
		</div>
	</div>
</div>
