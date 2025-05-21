<script lang="ts">
	import { fade } from 'svelte/transition';
	import { onMount } from 'svelte';
	import { currentUser } from '$lib/stores/auth';
	import type { Workspace } from '$lib/types/Workspace';
	import { workspaces as workspaceStore } from '$lib/stores/workspace';

	// API imports
	import { workspaceApi } from '$lib/api/workspaceApi';
	import { toast } from '$lib/stores/toast';
	
	export let visible = false;
	
	let workspaces: Workspace[] = [];
	let filteredWorkspaces: Workspace[] = [];
	let searchQuery = '';
	let isLoading = false;
	let error: string | null = null;
	
	// Pagination
	let currentPage = 1;
	let itemsPerPage = 10;
	let totalPages = 1;
	let paginatedWorkspaces: Workspace[] = [];
	
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
			applyFiltersAndPagination();
		} catch (err) {
			console.error('Failed to load workspaces:', err);
			error = 'Failed to load workspaces. Please try again.';
		} finally {
			isLoading = false;
		}
	}
	
	// Handle filtering and pagination
	function applyFiltersAndPagination() {
		// Apply search filter if search query exists
		if (searchQuery) {
			filteredWorkspaces = workspaces.filter(workspace => 
				workspace.name.toLowerCase().includes(searchQuery.toLowerCase())
			);
		} else {
			filteredWorkspaces = [...workspaces];
		}
		
		// Calculate total pages
		totalPages = Math.ceil(filteredWorkspaces.length / itemsPerPage);
		
		// Reset to first page if current page exceeds total pages
		if (currentPage > totalPages) {
			currentPage = 1;
		}
		
		// Apply pagination
		const startIndex = (currentPage - 1) * itemsPerPage;
		paginatedWorkspaces = filteredWorkspaces.slice(startIndex, startIndex + itemsPerPage);
	}
	
	function goToPage(page: number) {
		if (page >= 1 && page <= totalPages) {
			currentPage = page;
			applyFiltersAndPagination();
		}
	}
	
	function search() {
		currentPage = 1; // Reset to first page when searching
		applyFiltersAndPagination();
	}
	
	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleDateString();
	}
	
	function openAutoInviteModal(event: CustomEvent<{ workspaceId: string, workspaceName: string }>) {

	}
	
	function closeAutoInviteModal() {
		
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
			
			// update the workspace list
			workspaces = workspaces.map(ws => 
				ws.id === workspace.id ? { ...ws, auto_invite_enabled: newStatus } : ws
			);

			console.log(workspaces);
		} catch (err) {
			console.error('Failed to toggle auto-invite:', err);
			toast.error(err)
		}
	}
	
	// Delete function removed in instance scope
	
	function selectWorkspace(workspace: Workspace) {
		// Callback to parent component to switch to the selected workspace
		// onWorkspaceSelect(id, name);
		workspaceStore.switchWorkspace(workspace);
		toast.info(`Switched to workspace: ${workspace.name}`);
	}
	
	function manageAutoInvite(id: string) {
		openAutoInviteModal({ detail: { workspaceId: id, workspaceName: workspaces.find(ws => ws.id === id)?.name || '' } });
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
				<!-- Search bar -->
				<div class="mb-4 relative">
					<div class="flex">
						<div class="relative w-full">
							<div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
								<i class="fas fa-search text-gray-400"></i>
							</div>
							<input 
								type="text" 
								bind:value={searchQuery} 
								on:input={search}
								class="block w-full p-3 ps-10 text-sm rounded-lg bg-gray-100 dark:bg-gray-700 border border-gray-200 dark:border-gray-700 theme-text-primary focus:ring-blue-500 focus:border-blue-500 placeholder-gray-400"
								placeholder="Search workspaces..." 
							/>
						</div>
					</div>
				</div>
                
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
						{#each paginatedWorkspaces as workspace}
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
											on:click={() => selectWorkspace(workspace)} 
											class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" 
											title="Switch to this workspace"
											aria-label="Switch to this workspace"
										>
											<i class="fas fa-exchange-alt theme-text-secondary"></i>
										</button>
										
										<button 
											on:click={() => manageAutoInvite(workspace.id)} 
											class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" 
											title="Manage Auto-Invite Settings"
											aria-label="Manage Auto-Invite Settings"
										>
											<i class="fas fa-envelope theme-text-secondary"></i>
										</button>
									</div>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			{/if}
		</div>
		
		<!-- Pagination controls -->
		{#if !isLoading && !error && workspaces.length > 0 && totalPages > 1}
			<div class="flex items-center justify-between mt-4">
				<div class="text-sm theme-text-secondary">
					Showing {Math.min((currentPage - 1) * itemsPerPage + 1, filteredWorkspaces.length)}-{Math.min(currentPage * itemsPerPage, filteredWorkspaces.length)} of {filteredWorkspaces.length} workspaces
				</div>
				<div class="flex gap-1">
					<button
						on:click={() => goToPage(1)}
						disabled={currentPage === 1}
						class="p-2 rounded theme-bg-secondary hover:bg-blue-500/20 disabled:opacity-50 disabled:hover:bg-transparent"
						aria-label="First page"
					>
						<i class="fas fa-angle-double-left theme-text-secondary"></i>
					</button>
					<button
						on:click={() => goToPage(currentPage - 1)}
						disabled={currentPage === 1}
						class="p-2 rounded theme-bg-secondary hover:bg-blue-500/20 disabled:opacity-50 disabled:hover:bg-transparent"
						aria-label="Previous page"
					>
						<i class="fas fa-angle-left theme-text-secondary"></i>
					</button>
					
					<!-- Page numbers -->
					{#each Array(Math.min(5, totalPages)) as _, i}
						{@const pageNum = (() => {
							// Show pages around current page
							const halfWindow = 2;
							let start = Math.max(1, currentPage - halfWindow);
							let end = Math.min(totalPages, start + 4);
							
							// Adjust start if we're near the end
							if (end === totalPages) {
								start = Math.max(1, end - 4);
							}
							
							return start + i;
						})()}
						
						{#if pageNum <= totalPages}
							<button
								on:click={() => goToPage(pageNum)}
								class={`w-8 h-8 rounded text-sm flex items-center justify-center ${
									currentPage === pageNum 
										? 'bg-blue-600 text-white' 
										: 'theme-bg-secondary theme-text-secondary hover:bg-blue-500/20'
								}`}
							>
								{pageNum}
							</button>
						{/if}
					{/each}
					
					<button
						on:click={() => goToPage(currentPage + 1)}
						disabled={currentPage === totalPages}
						class="p-2 rounded theme-bg-secondary hover:bg-blue-500/20 disabled:opacity-50 disabled:hover:bg-transparent"
						aria-label="Next page"
					>
						<i class="fas fa-angle-right theme-text-secondary"></i>
					</button>
					<button
						on:click={() => goToPage(totalPages)}
						disabled={currentPage === totalPages}
						class="p-2 rounded theme-bg-secondary hover:bg-blue-500/20 disabled:opacity-50 disabled:hover:bg-transparent"
						aria-label="Last page"
					>
						<i class="fas fa-angle-double-right theme-text-secondary"></i>
					</button>
				</div>
			</div>
		{/if}
		
		<!-- No workspace actions in instance scope -->
	</div>
</div>
