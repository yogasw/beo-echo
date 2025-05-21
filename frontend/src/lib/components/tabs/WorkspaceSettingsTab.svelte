<script lang="ts">
	import { currentWorkspace, workspaceStore } from '$lib/stores/workspace';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { toast } from '$lib/stores/toast';
	import { workspaceApi } from '$lib/api/workspaceApi';
	
	let editMode = false;
	let workspaceName = '';
	let isLoading = false;
	let showAddMemberModal = false;
	let memberEmail = '';
	let memberRole = 'member';
	let isAddingMember = false;
	
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
		
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}
		
		isLoading = true;
		
		try {
			await workspaceApi.updateWorkspace($currentWorkspace.id, { name: workspaceName });
			toast.success('Workspace settings updated successfully');
			editMode = false;
			
			// Update workspace in store
			workspaceStore.update(state => ({
				...state,
				currentWorkspace: state.currentWorkspace ? {
					...state.currentWorkspace,
					name: workspaceName
				} : null,
				// Also update in the workspaces array
				workspaces: state.workspaces.map(w => 
					w.id === $currentWorkspace?.id ? { ...w, name: workspaceName } : w
				)
			}));
		} catch (error) {
			toast.error('Failed to update workspace settings');
			console.error('Failed to update workspace:', error);
		} finally {
			isLoading = false;
		}
	}
	
	// Load workspace 
	async function loadWorkspaceMembers() {
		if (!$currentWorkspace) return;
		
		try {
			workspaceApi.getWorkspaceMembers($currentWorkspace.id)
				.then(response => {
					members = response;
				})
				.catch(error => {
					toast.error('Failed to load workspace members');
					console.error('Failed to load workspace members:', error);
				});
		} catch (error) {
			console.error('Failed to load workspace members:', error);
		}
	}
	
	// Update member role
	async function updateMemberRole(memberId: string, role: string) {
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}
		
		try {
			await workspaceApi.updateUserRole($currentWorkspace.id, memberId, role);
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
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}
		
		try {
			await workspaceApi.removeUserFromWorkspace($currentWorkspace.id, memberId);
			toast.success('Member removed successfully');
			
			// Update local state
			members = members.filter(member => member.id !== memberId);
		} catch (error) {
			toast.error('Failed to remove member');
			console.error('Failed to remove member:', error);
		}
	}
	
	// Toggle add member modal
	function toggleAddMemberModal() {
		showAddMemberModal = !showAddMemberModal;
		// Reset form when closing
		if (!showAddMemberModal) {
			memberEmail = '';
			memberRole = 'member';
		}
	}
	
	// Add member to workspace
	async function addMember() {
		if (!memberEmail.trim() || !memberEmail.includes('@')) {
			toast.error('Please enter a valid email address');
			return;
		}
		
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}
		
		isAddingMember = true;
		
		try {
			await workspaceApi.addMember($currentWorkspace.id, {
				email: memberEmail,
				role: memberRole
			});
			
			toast.success(`Member ${memberEmail} added to workspace`);
			toggleAddMemberModal();
			
			// Refresh member list
			await loadWorkspaceMembers();
		} catch (error) {
			toast.error(error || 'Failed to add member. Only existing users can be added to workspaces.');
			console.error('Failed to add member:', error);
		} finally {
			isAddingMember = false;
		}
	}
	
	// Delete workspace
	async function deleteWorkspace() {
		if (!$currentWorkspace) {
			toast.error('No workspace selected');
			return;
		}
		
		// Show confirmation dialog
		if (!confirm(`Are you sure you want to delete the workspace "${$currentWorkspace.name}"? This action cannot be undone and will delete all associated data.`)) {
			return;
		}
		
		try {
			await workspaceApi.deleteWorkspace($currentWorkspace.id);
			toast.success('Workspace deleted successfully');
			
			// Remove workspace from store and redirect
			workspaceStore.update(state => ({
				...state,
				currentWorkspace: null,
				workspaces: state.workspaces.filter(w => w.id !== $currentWorkspace?.id)
			}));
			
			// Redirect to workspaces list
			window.location.href = '/workspaces';
		} catch (error) {
			toast.error('Failed to delete workspace');
			console.error('Failed to delete workspace:', error);
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
											aria-label="Remove member"
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
						on:click={toggleAddMemberModal}
						class={ThemeUtils.themeBgSecondary('px-4 py-2 rounded-md hover:bg-blue-500/20 flex items-center gap-2')}
					>
						<i class="fas fa-user-plus text-blue-400"></i>
						<span class="theme-text-primary">Add Member</span>
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
					on:click={deleteWorkspace}
					class="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md flex items-center gap-2"
				>
					<i class="fas fa-trash-alt"></i>
					Delete
				</button>
			</div>
		</div>
	{/if}
</div>

<!-- Add Member Modal -->
{#if showAddMemberModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div class={ThemeUtils.card('max-w-md w-full mx-4 p-6')}>
			<div class="flex justify-between items-center mb-4">
				<h3 class="text-lg font-semibold theme-text-primary">Add Member</h3>
				<button 
					on:click={toggleAddMemberModal}
					class="theme-text-primary hover:text-gray-500 dark:hover:text-gray-300"
					aria-label="Close modal"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
			
			<form on:submit|preventDefault={addMember}>
				<div class="mb-4">
					<label for="member-email" class="block theme-text-secondary text-sm mb-1">Email Address</label>
					<input
						id="member-email"
						type="email"
						bind:value={memberEmail}
						placeholder="Enter email address of existing user"
						required
						class={ThemeUtils.themeBgSecondary('w-full p-2 rounded theme-border theme-text-primary')}
					/>
					<p class="text-xs theme-text-secondary mt-1">
						Note: Only existing users can be added to workspaces
					</p>
				</div>
				
				<div class="mb-6">
					<label for="member-role" class="block theme-text-secondary text-sm mb-1">Role</label>
					<select
						id="member-role"
						bind:value={memberRole}
						class={ThemeUtils.themeBgSecondary('w-full p-2 rounded theme-border theme-text-primary')}
					>
						<option value="admin">Admin</option>
						<option value="member">Member</option>
						<option value="readonly">Read Only</option>
					</select>
					<p class="text-xs theme-text-secondary mt-1">
						<strong>Admin:</strong> Can manage workspace settings and members<br>
						<strong>Member:</strong> Can create and edit projects<br>
						<strong>Read Only:</strong> Can only view projects
					</p>
				</div>
				
				<div class="flex justify-end">
					<button
						type="button"
						on:click={toggleAddMemberModal}
						class={ThemeUtils.themeBgSecondary('px-4 py-2 rounded-md mr-2')}
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={isAddingMember}
						class={`bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md 
							flex items-center gap-2 ${isAddingMember ? 'opacity-50 cursor-not-allowed' : ''}`}
					>
						{#if isAddingMember}
							<i class="fas fa-spinner fa-spin"></i>
						{:else}
							<i class="fas fa-user-plus"></i>
						{/if}
						Add Member
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
