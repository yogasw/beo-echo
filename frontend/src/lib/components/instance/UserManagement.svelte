<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border mb-4">
		<h3 class="theme-text-primary font-medium mb-3">Users</h3>
		<div class="overflow-x-auto mb-4">
			<table class="w-full text-sm text-left">
				<thead class="text-xs uppercase theme-text-secondary">
					<tr>
						<th scope="col" class="px-4 py-3">User</th>
						<th scope="col" class="px-4 py-3">Email</th>
						<th scope="col" class="px-4 py-3">Role</th>
						<th scope="col" class="px-4 py-3">Status</th>
						<th scope="col" class="px-4 py-3">Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr class="theme-border-subtle border-b">
						<td class="px-4 py-3 flex items-center gap-3">
							<div class="w-8 h-8 rounded-full theme-bg-secondary flex items-center justify-center">
								<span class="theme-text-primary text-xs">JD</span>
							</div>
							<span class="theme-text-primary">John Doe</span>
						</td>
						<td class="px-4 py-3 theme-text-secondary">john@example.com</td>
						<td class="px-4 py-3">
							<span class="bg-blue-500/20 text-blue-400 px-2 py-1 rounded-full text-xs">
								Admin
							</span>
						</td>
						<td class="px-4 py-3">
							<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded-full text-xs">
								Active
							</span>
						</td>
						<td class="px-4 py-3">
							<div class="flex items-center gap-2">
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20">
									<i class="fas fa-edit theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20">
									<i class="fas fa-trash theme-text-secondary"></i>
								</button>
							</div>
						</td>
					</tr>
					<!-- Sample additional user -->
					<tr class="theme-border-subtle border-b">
						<td class="px-4 py-3 flex items-center gap-3">
							<div class="w-8 h-8 rounded-full theme-bg-secondary flex items-center justify-center">
								<span class="theme-text-primary text-xs">JS</span>
							</div>
							<span class="theme-text-primary">Jane Smith</span>
						</td>
						<td class="px-4 py-3 theme-text-secondary">jane@example.com</td>
						<td class="px-4 py-3">
							<span class="bg-purple-500/20 text-purple-400 px-2 py-1 rounded-full text-xs">
								User
							</span>
						</td>
						<td class="px-4 py-3">
							<span class="bg-green-500/20 text-green-400 px-2 py-1 rounded-full text-xs">
								Active
							</span>
						</td>
						<td class="px-4 py-3">
							<div class="flex items-center gap-2">
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20">
									<i class="fas fa-edit theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20">
									<i class="fas fa-trash theme-text-secondary"></i>
								</button>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
		<button class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center gap-2">
			<i class="fas fa-plus"></i>
			<span>Add User</span>
		</button>
	</div>
</div>
