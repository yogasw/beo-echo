<script lang="ts">
	import { fade } from 'svelte/transition';
	
	export let visible = false;
</script>

<div class="p-4" transition:fade={{ duration: 200 }}>
	<div class="theme-bg-primary p-4 rounded-lg border theme-border">
		<div class="space-y-4">
			<div>
				<h3 class="theme-text-primary font-medium mb-2">Password Policy</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Minimum Password Length</p>
							<p class="theme-text-muted text-xs">Set the minimum character requirement for passwords</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="8">8 characters</option>
							<option value="10">10 characters</option>
							<option value="12">12 characters</option>
							<option value="14">14 characters</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Require Password Complexity</p>
							<p class="theme-text-muted text-xs">Require numbers, special characters, and mixed case</p>
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
							<p class="theme-text-primary text-sm">Password Expiration</p>
							<p class="theme-text-muted text-xs">Require users to change passwords periodically</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="never">Never</option>
							<option value="30">30 days</option>
							<option value="60">60 days</option>
							<option value="90">90 days</option>
							<option value="180">180 days</option>
						</select>
					</div>
				</div>
			</div>
			
			<div class="pt-3 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-2">Session Settings</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Session Timeout</p>
							<p class="theme-text-muted text-xs">Automatically log users out after inactivity</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="30">30 minutes</option>
							<option value="60">1 hour</option>
							<option value="120">2 hours</option>
							<option value="240">4 hours</option>
							<option value="480">8 hours</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Enforce Single Session</p>
							<p class="theme-text-muted text-xs">Restrict users to one active session at a time</p>
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
				</div>
			</div>
			
			<div class="pt-3 border-t theme-border">
				<h3 class="theme-text-primary font-medium mb-2">API Security</h3>
				<div class="space-y-3">
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">API Token Expiration</p>
							<p class="theme-text-muted text-xs">Set the lifetime for generated API tokens</p>
						</div>
						<select class="theme-bg-secondary theme-border border rounded px-3 py-2 text-sm theme-text-primary">
							<option value="never">Never</option>
							<option value="7">7 days</option>
							<option value="30">30 days</option>
							<option value="90">90 days</option>
							<option value="365">1 year</option>
						</select>
					</div>
					
					<div class="flex items-center justify-between">
						<div>
							<p class="theme-text-primary text-sm">Rate Limiting</p>
							<p class="theme-text-muted text-xs">Limit API requests to prevent abuse</p>
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
