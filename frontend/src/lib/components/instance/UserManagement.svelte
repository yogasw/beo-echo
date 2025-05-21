<script lang="ts">
	import { fade } from 'svelte/transition';
	import { toast } from '$lib/stores/toast';
	import { onMount } from 'svelte';
	import type { User } from '$lib/types/User';
	import { userManagementApi } from '$lib/api/userManagement';

	let users: User[] = [];
	let loading = true;

	// Modal state management
	let showAddModal = false;
	let showEditModal = false;
	let showDeleteModal = false;
	let selectedUser: User | null = null;

	// Form data for add/edit
	let formData = {
		name: '',
		email: '',
		is_active: false
	};

	// Status style mapping for badge colors (now based on is_active)
	const activeStyles = {
		true: 'bg-green-500/20 text-green-400',
		false: 'bg-red-500/20 text-red-400'
	};

	// Owner badge style
	const ownerBadgeStyle = 'bg-amber-500/20 text-amber-400 ml-2 px-2 py-1 rounded-full text-xs';

	onMount(async () => {
		await loadUsers();
	});

	// Load all user from api
	async function loadUsers() {
		loading = true;
		try {
			users = await userManagementApi.getAllUsers();
		} catch (error) {
			// Error is already shown by the API function
			users = [];
		} finally {
			loading = false;
		}
	}

	// Open add user modal
	function openAddModal() {
		// Reset form data
		formData = {
			name: '',
			email: '',
			is_active: false
		};
		showAddModal = true;
	}

	// Open edit user modal
	function openEditModal(user: User) {
		selectedUser = user;
		formData = {
			name: user.name,
			email: user.email,
			is_active: user.is_active || false
		};
		showEditModal = true;
	}

	// Open delete confirmation modal
	function openDeleteModal(user: User) {
		selectedUser = user;
		showDeleteModal = true;
	}

	// Close all modals
	function closeModals() {
		showAddModal = false;
		showEditModal = false;
		showDeleteModal = false;
		selectedUser = null;
	}

	// Handle adding a new user
	async function handleAddUser() {
		if (!formData.email) {
			toast.error('Email is required');
			return;
		}

		try {
			// Reload users to get the updated list
			await loadUsers();
			toast.success('User added successfully');
			closeModals();
		} catch (error) {
			// Error is already shown by the API function
		}
	}

	// Handle updating an existing user
	async function handleEditUser() {
		if (selectedUser) {
			try {
				// Update user via API
				// Reload users to get the updated list
				await loadUsers();
				toast.success('User updated successfully');
				closeModals();
			} catch (error) {
				// Error is already shown by the API function
			}
		}
	}

	// Handle deleting a user
	async function handleDeleteUser() {
		if (selectedUser) {
			try {
				// Remove user

				// Reload users to get the updated list
				await loadUsers();
				toast.success('User removed successfully');
				closeModals();
			} catch (error) {
				// Error is already shown by the API function
			}
		}
	}

	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border mb-4">
		<h3 class="theme-text-primary font-medium mb-3">Users</h3>

		{#if loading}
			<div class="flex justify-center items-center p-8">
				<div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
			</div>
		{:else}
			<div class="overflow-x-auto mb-4">
				{#if users.length === 0}
					<div class="text-center p-4 theme-text-secondary">No users found in this instance.</div>
				{:else}
					<table class="w-full text-sm text-left">
						<thead class="text-xs uppercase theme-text-secondary">
							<tr>
								<th scope="col" class="px-4 py-3">User</th>
								<th scope="col" class="px-4 py-3">Email</th>
								<th scope="col" class="px-4 py-3">Status</th>
								<th scope="col" class="px-4 py-3">Actions</th>
							</tr>
						</thead>
						<tbody>
							{#each users as user}
								<tr class="theme-border-subtle border-b">
									<td class="px-4 py-3 flex items-center gap-3">
										<span class="theme-text-primary">{user.name}</span>
										{#if user.is_owner}
											<span class={ownerBadgeStyle}>
												<i class="fas fa-crown mr-1"></i>
												Owner
											</span>
										{/if}
									</td>
									<td class="px-4 py-3 theme-text-secondary">{user.email}</td>
									<td class="px-4 py-3">
										<span
											class="{activeStyles[
												String(user.is_active !== false)
											]} px-2 py-1 rounded-full text-xs"
										>
											{user.is_active !== false ? 'Active' : 'Inactive'}
										</span>
									</td>
									<td class="px-4 py-3">
										<div class="flex items-center gap-2">
											<button
												class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20"
												on:click={() => openEditModal(user)}
												aria-label="Edit user"
											>
												<i class="fas fa-edit theme-text-secondary"></i>
											</button>
											<button
												class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20"
												on:click={() => openDeleteModal(user)}
												aria-label="Delete user"
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
		{/if}
		<button
			class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center gap-2"
			on:click={openAddModal}
		>
			<i class="fas fa-plus"></i>
			<span>Add User</span>
		</button>
	</div>
</div>

<!-- Add User Modal -->
{#if showAddModal}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div
			class="theme-bg-primary rounded-lg shadow-lg max-w-md w-full mx-4"
			transition:fade={{ duration: 200 }}
		>
			<div class="p-4 border-b theme-border flex justify-between items-center">
				<h3 class="theme-text-primary font-medium">Add New User</h3>
				<button
					class="theme-text-secondary hover:text-gray-500 dark:hover:text-gray-400"
					on:click={closeModals}
					aria-label="Close modal"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
			<div class="p-4">
				<form on:submit|preventDefault={handleAddUser}>
					<div class="mb-4">
						<label for="email" class="block theme-text-secondary text-sm mb-2">Email</label>
						<input
							type="email"
							id="email"
							class="block w-full p-2 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							placeholder="Enter user's email"
							bind:value={formData.email}
							required
						/>
						<p class="mt-1 text-xs theme-text-secondary">
							The user will be invited to join this instance.
						</p>
					</div>
					<div class="flex justify-end gap-3 mt-6">
						<button
							type="button"
							class="px-4 py-2 theme-bg-secondary theme-text-primary rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
							on:click={closeModals}
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
						>
							Add User
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Edit User Modal -->
{#if showEditModal && selectedUser}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div
			class="theme-bg-primary rounded-lg shadow-lg max-w-md w-full mx-4"
			transition:fade={{ duration: 200 }}
		>
			<div class="p-4 border-b theme-border flex justify-between items-center">
				<h3 class="theme-text-primary font-medium">Edit User</h3>
				<button
					class="theme-text-secondary hover:text-gray-500 dark:hover:text-gray-400"
					on:click={closeModals}
					aria-label="Close modal"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
			<div class="p-4">
				<form on:submit|preventDefault={handleEditUser}>
					<div class="mb-4">
						<label for="edit-name" class="block theme-text-secondary text-sm mb-2">Name</label>
						<input
							type="text"
							id="edit-name"
							class="block w-full p-2 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							value={selectedUser.name}
							disabled
						/>
					</div>
					<div class="mb-4">
						<label for="edit-email" class="block theme-text-secondary text-sm mb-2">Email</label>
						<input
							type="email"
							id="edit-email"
							class="block w-full p-2 text-sm rounded-lg theme-bg-secondary border theme-border theme-text-primary focus:ring-blue-500 focus:border-blue-500"
							value={selectedUser.email}
							disabled
						/>
					</div>
					<div class="mb-4">
						<label for="edit-status" class="block theme-text-secondary text-sm mb-2">Status</label>
						<div class="flex items-center">
							<label class="inline-flex items-center cursor-pointer">
								<input
									type="checkbox"
									id="edit-status"
									checked={formData.status === 'Active'}
									on:change={() =>
										(formData.status = formData.status === 'Active' ? 'Inactive' : 'Active')}
									class="sr-only peer"
								/>
								<div
									class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-checked:bg-green-500
									rounded-full peer peer-focus:outline-none peer-focus:ring-2
									peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600
									peer-checked:after:translate-x-full peer-checked:after:border-white
									after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white
									after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5
									after:transition-all dark:border-gray-600"
								></div>
								<span class="ml-3 theme-text-primary text-sm">{formData.status}</span>
							</label>
						</div>
					</div>
					<div class="flex justify-end gap-3 mt-6">
						<button
							type="button"
							class="px-4 py-2 theme-bg-secondary theme-text-primary rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
							on:click={closeModals}
						>
							Cancel
						</button>
						<button
							type="submit"
							class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm"
						>
							Update User
						</button>
					</div>
				</form>
			</div>
		</div>
	</div>
{/if}

<!-- Delete User Modal -->
{#if showDeleteModal && selectedUser}
	<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
		<div
			class="theme-bg-primary rounded-lg shadow-lg max-w-md w-full mx-4"
			transition:fade={{ duration: 200 }}
		>
			<div class="p-4 border-b theme-border flex justify-between items-center">
				<h3 class="theme-text-primary font-medium">Confirm Removal</h3>
				<button
					class="theme-text-secondary hover:text-gray-500 dark:hover:text-gray-400"
					on:click={closeModals}
					aria-label="Close modal"
				>
					<i class="fas fa-times"></i>
				</button>
			</div>
			<div class="p-4">
				<div class="mb-4 theme-text-primary">
					<p>
						Are you sure you want to remove <strong>{selectedUser.name}</strong> from this instance?
					</p>
					<p class="mt-2 text-sm theme-text-secondary">This action cannot be undone.</p>
				</div>
				<div class="flex justify-end gap-3 mt-6">
					<button
						type="button"
						class="px-4 py-2 theme-bg-secondary theme-text-primary rounded-md text-sm hover:bg-gray-200 dark:hover:bg-gray-600"
						on:click={closeModals}
					>
						Cancel
					</button>
					<button
						type="button"
						class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md text-sm"
						on:click={handleDeleteUser}
					>
						Remove User
					</button>
				</div>
			</div>
		</div>
	</div>
{/if}
