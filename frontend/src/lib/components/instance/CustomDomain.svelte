<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<div class="space-y-4">
			<div>
				<h3 class="theme-text-primary font-medium mb-3">Domain Configuration</h3>
				<div class="space-y-4">
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							Custom Domain
						</label>
						<div class="flex gap-2">
							<input
								type="text"
								placeholder="yourdomain.com"
								class="flex-1 px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
							/>
							<button class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm">
								Verify
							</button>
						</div>
						<p class="theme-text-muted text-xs mt-2">
							Enter your custom domain name to use for this Beo Echo instance
						</p>
					</div>
					
					<div>
						<h4 class="theme-text-primary text-sm font-medium mb-2">DNS Configuration</h4>
						<div class="theme-bg-secondary p-3 rounded-md">
							<p class="theme-text-secondary text-xs mb-2">Add the following records to your DNS provider:</p>
							<div class="theme-bg-primary border theme-border rounded overflow-hidden">
								<div class="overflow-x-auto">
									<table class="w-full text-xs text-left">
										<thead class="uppercase theme-text-secondary theme-bg-secondary">
											<tr>
												<th scope="col" class="px-3 py-2">Type</th>
												<th scope="col" class="px-3 py-2">Name</th>
												<th scope="col" class="px-3 py-2">Value</th>
											</tr>
										</thead>
										<tbody>
											<tr class="theme-border-subtle border-b">
												<td class="px-3 py-2 theme-text-primary">CNAME</td>
												<td class="px-3 py-2 theme-text-secondary">@</td>
												<td class="px-3 py-2 theme-text-secondary">beo-echo.example.com</td>
											</tr>
										</tbody>
									</table>
								</div>
							</div>
						</div>
					</div>
					
					<div>
						<h4 class="theme-text-primary text-sm font-medium mb-2">SSL Configuration</h4>
						<div class="space-y-3">
							<div class="flex items-center justify-between">
								<div>
									<p class="theme-text-primary text-sm">Automatic SSL</p>
									<p class="theme-text-muted text-xs">Let's Encrypt certificates will be automatically provisioned</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" class="sr-only peer" checked>
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
									peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
									rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
									after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
									after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
									after:transition-all peer-checked:bg-blue-600"></div>
								</label>
							</div>
							
							<div class="flex items-center justify-between">
								<div>
									<p class="theme-text-primary text-sm">Force HTTPS</p>
									<p class="theme-text-muted text-xs">Redirect all HTTP traffic to HTTPS</p>
								</div>
								<label class="relative inline-flex items-center cursor-pointer">
									<input type="checkbox" class="sr-only peer" checked>
									<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
									peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
									rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
									after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
									after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
									after:transition-all peer-checked:bg-blue-600"></div>
								</label>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
