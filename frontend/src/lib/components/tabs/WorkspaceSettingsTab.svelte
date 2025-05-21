<script lang="ts">
	import { currentWorkspace, workspaceStore } from '$lib/stores/workspace';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { toast } from '$lib/stores/toast';
	
	let editMode = false;
	let workspaceName = '';
	let isLoading = false;
	interface Member {
		id: string;
		name: string;
		email: string;
		role: string;
	}
	
	let members: Member[] = [];
	
	// Initialize form with current workspace data
	$: if ($currentWorkspace && !editMode) {
		workspaceName = $currentWorkspace.name;
		console.log('Current workspace:', $currentWorkspace);
	}
	
	// Toggle edit mode
	function toggleEditMode() {
		editMode = !editMode;
		if (!editMode) {
			// Reset form if canceling edit
			workspaceName = $currentWorkspace?.name || '';
		}
	}
	
	// Save workspace settings
	async function saveWorkspaceSettings() {
		if (!workspaceName.trim()) {
			toast.error('Workspace name is required');
			return;
		}
		
		isLoading = true;
		
		try {
			// TODO: Implement API call to update workspace
			// await updateWorkspace($currentWorkspace.id, { name: workspaceName });
			
			// For now, just show success message
			toast.success('Workspace settings updated successfully');
			editMode = false;
		} catch (error) {
			toast.error('Failed to update workspace settings');
			console.error('Failed to update workspace:', error);
		} finally {
			isLoading = false;
		}
	}
	
	// Load workspace members
	async function loadWorkspaceMembers() {
		if (!$currentWorkspace) return;
		
		try {
			// TODO: Implement API call to get workspace members
			// const response = await getWorkspaceMembers($currentWorkspace.id);
			// members = response.data;
			
			// For now, use mock data
			members = [
				{ id: '1', name: 'John Doe', email: 'john@example.com', role: 'admin' },
				{ id: '2', name: 'Jane Smith', email: 'jane@example.com', role: 'member' }
			];
		} catch (error) {
			console.error('Failed to load workspace members:', error);
		}
	}
	
	// Update member role
	async function updateMemberRole(memberId: string, role: string) {
		try {
			// TODO: Implement API call to update member role
			// await updateWorkspaceMemberRole($currentWorkspace.id, memberId, role);
			
			toast.success('Member role updated successfully');
			
			// Update local state
			members = members.map(member => 
				member.id === memberId ? { ...member, role } : member
			);
		} catch (error) {
			toast.error('Failed to update member role');
			console.error('Failed to update member role:', error);
		}
	}
	
	// Remove member from workspace
	async function removeMember(memberId: string) {
		try {
			// TODO: Implement API call to remove member
			// await removeWorkspaceMember($currentWorkspace.id, memberId);
			
			toast.success('Member removed successfully');
			
			// Update local state
			members = members.filter(member => member.id !== memberId);
		} catch (error) {
			toast.error('Failed to remove member');
			console.error('Failed to remove member:', error);
		}
	}
	
	// Load members when component mounts or workspace changes
	$: if ($currentWorkspace) {
		loadWorkspaceMembers();
	}
</script>

<div class="workspace-settings w-full">
	<!-- Header Section -->
	<div class="flex justify-between items-center mb-6">
		<div>
			<h1 class="text-2xl font-semibold theme-text-primary">Workspace Settings</h1>
			<p class="theme-text-secondary text-sm mt-1">
				Manage your workspace details and team members
			</p>
		</div>
		
		<button
			on:click={() => toggleEditMode()}
			class={ThemeUtils.themeBgSecondary(`px-4 py-2 rounded-md flex items-center gap-2 transition-colors
				${editMode ? 'hover:bg-red-600' : 'hover:bg-blue-500/20'}`)}
		>
			{#if editMode}
				<i class="fas fa-times text-red-400"></i>
				<span class="theme-text-primary">Cancel</span>
			{:else}
				<i class="fas fa-edit text-blue-400"></i>
				<span class="theme-text-primary">Edit Workspace</span>
			{/if}
		</button>
	</div>
	
	<!-- Workspace Details Section -->
	<div class={ThemeUtils.card('mb-6 p-4')}>
		<h2 class="text-lg font-semibold theme-text-primary mb-4 flex items-center">
			<i class="fas fa-building text-blue-400 mr-2"></i>
			Workspace Details
		</h2>
		
		<div class="mb-4">
			<label for="workspace-name" class="block theme-text-secondary text-sm mb-1">Workspace Name</label>
			{#if editMode}
				<input
					id="workspace-name"
					type="text"
					bind:value={workspaceName}
					class={ThemeUtils.themeBgSecondary('w-full p-2 rounded theme-border theme-text-primary')}
					placeholder="Enter workspace name"
				/>
			{:else}
				<div id="workspace-name" class={ThemeUtils.themeBgSecondary('w-full p-2 rounded theme-border theme-text-primary')}>
					{$currentWorkspace?.name || 'No workspace selected'}
				</div>
			{/if}
		</div>
		
		<div class="mb-4">
			<label for="workspace-role" class="block theme-text-secondary text-sm mb-1">Your Role</label>
			<div id="workspace-role" class={ThemeUtils.themeBgSecondary('w-full p-2 rounded theme-border theme-text-primary')}>
				{$currentWorkspace?.role || 'N/A'}
			</div>
		</div>
		
		{#if editMode}
			<div class="flex justify-end mt-4">
				<button
					on:click={saveWorkspaceSettings}
					disabled={isLoading}
					class={`bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
						flex items-center gap-2 ${isLoading ? 'opacity-50 cursor-not-allowed' : ''}`}
				>
					{#if isLoading}
						<i class="fas fa-spinner fa-spin"></i>
					{:else}
						<i class="fas fa-save"></i>
					{/if}
					Save Changes
				</button>
			</div>
		{/if}
	</div>
	
	<!-- Team Members Section -->
	<div class={ThemeUtils.card('mb-6 p-4')}>
		<h2 class="text-lg font-semibold theme-text-primary mb-4 flex items-center">
			<i class="fas fa-users text-blue-400 mr-2"></i>
			Team Members
		</h2>
		
		{#if members.length > 0}
			<div class="overflow-x-auto">
				<table class="w-full">
					<thead>
						<tr class="theme-text-secondary text-left">
							<th class="p-2">Name</th>
							<th class="p-2">Email</th>
							<th class="p-2">Role</th>
							<th class="p-2 text-right">Actions</th>
						</tr>
					</thead>
					<tbody>
						{#each members as member}
							<tr class="theme-border-subtle border-b">
								<td class="p-2 theme-text-primary">{member.name}</td>
								<td class="p-2 theme-text-primary">{member.email}</td>
								<td class="p-2">
									{#if editMode && $currentWorkspace?.role === 'admin'}
										<select
											value={member.role}
											on:change={(e) => {
												if (e.target) {
													const target = e.target as HTMLSelectElement;
													updateMemberRole(member.id, target.value);
												}
											}}
											class={ThemeUtils.themeBgSecondary('p-1 rounded theme-border theme-text-primary')}
										>
											<option value="admin">Admin</option>
											<option value="member">Member</option>
											<option value="readonly">Read Only</option>
										</select>
									{:else}
										<span class="theme-text-primary">
											{member.role === 'admin' ? 'Admin' : member.role === 'readonly' ? 'Read Only' : 'Member'}
										</span>
									{/if}
								</td>
								<td class="p-2 text-right">
									{#if editMode && $currentWorkspace?.role === 'admin'}
										<button
											on:click={() => removeMember(member.id)}
											class="text-red-400 hover:text-red-600 p-1"
											title="Remove member"
										>
											<i class="fas fa-trash-alt"></i>
										</button>
									{/if}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
			
			{#if editMode && $currentWorkspace?.role === 'admin'}
				<div class="mt-4">
					<button
						class={ThemeUtils.themeBgSecondary('px-4 py-2 rounded-md hover:bg-blue-500/20 flex items-center gap-2')}
					>
						<i class="fas fa-user-plus text-blue-400"></i>
						<span class="theme-text-primary">Invite Member</span>
					</button>
				</div>
			{/if}
		{:else}
			<div class="theme-text-secondary text-center py-4">
				No members found in this workspace
			</div>
		{/if}
	</div>
	
	<!-- Danger Zone -->
	{#if editMode && $currentWorkspace?.role === 'admin'}
		<div class="bg-red-600/10 border border-red-500/30 rounded-md p-4 mb-6">
			<h2 class="text-lg font-semibold text-red-400 mb-4 flex items-center">
				<i class="fas fa-exclamation-triangle mr-2"></i>
				Danger Zone
			</h2>
			
			<div class="flex justify-between items-center">
				<div>
					<h3 class="theme-text-primary font-medium">Delete Workspace</h3>
					<p class="text-sm theme-text-secondary">
						Permanently delete this workspace and all associated data
					</p>
				</div>
				
				<button
					class="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md flex items-center gap-2"
				>
					<i class="fas fa-trash-alt"></i>
					Delete
				</button>
			</div>
		</div>
	{/if}
</div>
