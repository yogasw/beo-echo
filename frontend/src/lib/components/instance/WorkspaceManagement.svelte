<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<h3 class="theme-text-primary font-medium mb-3">All Workspaces</h3>
		<div class="overflow-x-auto mb-4">
			<table class="w-full text-sm text-left">
				<thead class="text-xs uppercase theme-text-secondary">
					<tr>
						<th scope="col" class="px-4 py-3">Workspace</th>
						<th scope="col" class="px-4 py-3">Owner</th>
						<th scope="col" class="px-4 py-3">Members</th>
						<th scope="col" class="px-4 py-3">Created</th>
						<th scope="col" class="px-4 py-3">Actions</th>
					</tr>
				</thead>
				<tbody>
					<tr class="theme-border-subtle border-b">
						<td class="px-4 py-3 theme-text-primary">Marketing Team</td>
						<td class="px-4 py-3 theme-text-secondary">jane@example.com</td>
						<td class="px-4 py-3 theme-text-secondary">8</td>
						<td class="px-4 py-3 theme-text-secondary">2023-05-15</td>
						<td class="px-4 py-3">
							<div class="flex items-center gap-2">
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" title="Manage Members">
									<i class="fas fa-users theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" title="Edit">
									<i class="fas fa-edit theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20" title="Delete">
									<i class="fas fa-trash theme-text-secondary"></i>
								</button>
							</div>
						</td>
					</tr>
					<!-- Sample additional workspace -->
					<tr class="theme-border-subtle border-b">
						<td class="px-4 py-3 theme-text-primary">Development Team</td>
						<td class="px-4 py-3 theme-text-secondary">john@example.com</td>
						<td class="px-4 py-3 theme-text-secondary">12</td>
						<td class="px-4 py-3 theme-text-secondary">2023-03-22</td>
						<td class="px-4 py-3">
							<div class="flex items-center gap-2">
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" title="Manage Members">
									<i class="fas fa-users theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" title="Edit">
									<i class="fas fa-edit theme-text-secondary"></i>
								</button>
								<button class="p-2 theme-bg-secondary rounded-full hover:bg-red-500/20" title="Delete">
									<i class="fas fa-trash theme-text-secondary"></i>
								</button>
							</div>
						</td>
					</tr>
				</tbody>
			</table>
		</div>
		
		<!-- Workspace Actions -->
		<div class="flex gap-3">
			<button class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm flex items-center gap-2">
				<i class="fas fa-plus"></i>
				<span>Create Workspace</span>
			</button>
			
			<button class="px-3 py-2 theme-bg-secondary hover:bg-gray-600 theme-text-primary rounded-md text-sm flex items-center gap-2">
				<i class="fas fa-download"></i>
				<span>Export List</span>
			</button>
		</div>
	</div>
</div>
