<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<div class="space-y-4">
			<div>
				<h3 class="theme-text-primary font-medium mb-3">Instance Information</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							Instance Name
						</label>
						<input
							type="text"
							value="Beo Echo"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							Instance URL
						</label>
						<input
							type="text"
							value="https://beo-echo.example.com"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
				</div>
			</div>
			
			<div class="pt-4 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-3">Appearance</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Default Theme</p>
							<p class="theme-text-muted text-xs">Set the default theme for new users</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="dark">Dark Mode</option>
							<option value="light">Light Mode</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Custom Branding</p>
							<p class="theme-text-muted text-xs">Upload your own logo and customize colors</p>
						</div>
						<label class="relative inline-flex items-center cursor-pointer">
							<input type="checkbox" class="sr-only peer">
							<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
							peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
							rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
							after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
							after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
							after:transition-all peer-checked:bg-blue-600"></div>
						</label>
					</div>
					
					<!-- Custom branding options (conditionally shown) -->
					<div class="ml-6 pl-3 border-l theme-border space-y-4">
						<div>
							<label class="block text-sm font-medium theme-text-primary mb-2">
								Logo
							</label>
							<div class="flex items-center gap-4">
								<div class="h-16 w-16 bg-blue-600 rounded flex items-center justify-center">
									<i class="fas fa-mountain text-white text-xl"></i>
								</div>
								<button class="px-3 py-2 theme-bg-secondary hover:bg-gray-600 theme-text-primary rounded-md text-xs">
									Upload Logo
								</button>
							</div>
						</div>
						
						<div>
							<label class="block text-sm font-medium theme-text-primary mb-2">
								Primary Color
							</label>
							<div class="flex items-center gap-2">
								<input
									type="color"
									value="#3b82f6"
									class="h-8 w-8 rounded border theme-border overflow-hidden cursor-pointer"
								/>
								<input
									type="text"
									value="#3b82f6"
									class="w-24 px-3 py-1 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
								/>
							</div>
						</div>
					</div>
				</div>
			</div>
			
			<div class="pt-4 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-3">Email Settings</h3>
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							SMTP Server
						</label>
						<input
							type="text"
							placeholder="smtp.example.com"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							SMTP Port
						</label>
						<input
							type="text"
							placeholder="587"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							SMTP Username
						</label>
						<input
							type="text"
							placeholder="user@example.com"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
					<div>
						<label class="block text-sm font-medium theme-text-primary mb-2">
							SMTP Password
						</label>
						<input
							type="password"
							placeholder="••••••••••"
							class="w-full px-3 py-2 text-sm theme-bg-secondary theme-border border rounded theme-text-primary"
						/>
					</div>
					<div class="md:col-span-2">
						<div class="flex items-center gap-2">
							<button class="px-3 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm">
								Test Connection
							</button>
							<span class="theme-text-muted text-xs">Verify email configuration</span>
						</div>
					</div>
				</div>
			</div>
			
			<div class="pt-4 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-3">Data Retention</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Auto-delete Request Logs</p>
							<p class="theme-text-muted text-xs">Automatically remove old request logs</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="never">Never</option>
							<option value="7">After 7 days</option>
							<option value="30">After 30 days</option>
							<option value="90">After 90 days</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Database Backup</p>
							<p class="theme-text-muted text-xs">Schedule automatic database backups</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="never">Never</option>
							<option value="daily">Daily</option>
							<option value="weekly">Weekly</option>
							<option value="monthly">Monthly</option>
						</select>
					</div>
				</div>
			</div>
		</div>
	</div>
</div>
