<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<div class="space-y-5">
			<div>
				<div class="flex items-center justify-between mb-3">
					<h3 class="theme-text-primary font-medium">SSO Providers</h3>
					<div>
						<button class="px-3 py-2 theme-bg-secondary hover:bg-gray-600 theme-text-primary rounded-md text-sm flex items-center gap-2">
							<i class="fas fa-plus"></i>
							<span>Add Provider</span>
						</button>
					</div>
				</div>
				
				<div class="space-y-4">
					<!-- Google SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-white rounded flex items-center justify-center">
									<i class="fab fa-google text-gray-800"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">Google</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
					
					<!-- GitHub SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-gray-900 rounded flex items-center justify-center">
									<i class="fab fa-github text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">GitHub</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
					
					<!-- Microsoft SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-blue-500 rounded flex items-center justify-center">
									<i class="fab fa-microsoft text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">Microsoft</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
					
					<!-- SAML SSO -->
					<div class="theme-bg-secondary p-3 rounded-lg">
						<div class="flex items-center justify-between">
							<div class="flex items-center gap-3">
								<div class="w-8 h-8 bg-blue-600 rounded flex items-center justify-center">
									<i class="fas fa-id-card text-white"></i>
								</div>
								<div>
									<h4 class="theme-text-primary text-sm font-medium">SAML</h4>
									<p class="theme-text-muted text-xs">Not configured</p>
								</div>
							</div>
							<div>
								<button class="px-3 py-1 bg-gray-600 hover:bg-gray-500 text-white rounded text-xs">
									Configure
								</button>
							</div>
						</div>
					</div>
				</div>
			</div>
			
			<div class="pt-4 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-3">SSO Settings</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Default Sign-in Method</p>
							<p class="theme-text-muted text-xs">Select which sign-in method to display by default</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="local">Local Authentication</option>
							<option value="google">Google</option>
							<option value="github">GitHub</option>
							<option value="microsoft">Microsoft</option>
							<option value="saml">SAML</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Allow Local Authentication</p>
							<p class="theme-text-muted text-xs">Enable login with username and password</p>
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
							<p class="theme-text-primary text-sm">Auto-create User Accounts</p>
							<p class="theme-text-muted text-xs">Automatically create accounts for new SSO users</p>
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
