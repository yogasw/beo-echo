<script lang="ts">
	import { getProxyTargets, updateEndpoint, type ProxyTarget, type Endpoint } from '$lib/api/BeoApi';
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast'; 
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { theme } from '$lib/stores/theme';
	import { currentWorkspace } from '$lib/stores/workspace';
	import { updateEndpoint as storeUpdateEndpoint } from '$lib/stores/saveButton';
	
	export let endpoint: Endpoint;
	export let proxyTargets: ProxyTarget[] = [];
	export let isLoading: boolean = false;
	export let onChange: (endpoint: Endpoint) => void;
	
	let error: string | null = null;

	// Load proxy targets on mount if we don't have them
	onMount(async () => {
		if (!proxyTargets || proxyTargets.length === 0) {
			try {
				isLoading = true;
				proxyTargets = await getProxyTargets(endpoint.project_id);
				isLoading = false;
			} catch (err) {
				console.error('Failed to load proxy targets:', err);
				toast.error('Failed to load proxy targets');
				isLoading = false;
			}
		}
	});
	
	// Handle proxy toggle
	async function handleProxyToggle() {
		// Toggle the use_proxy flag
		const useProxy = !endpoint.use_proxy;
		
		// If disabling proxy, also clear the target
		const proxyTargetId = useProxy ? endpoint.proxy_target_id : null;
		
		try {
			// Update the endpoint via API
			await updateEndpoint(endpoint.project_id, endpoint.id, {
				use_proxy: useProxy,
				proxy_target_id: proxyTargetId
			});
			
			// Update local state
			endpoint = storeUpdateEndpoint('use_proxy', useProxy, endpoint);
			if (!useProxy) {
				endpoint = storeUpdateEndpoint('proxy_target_id', null, endpoint);
			}
			
			// Notify parent component
			onChange(endpoint);
			
			// Confirm success
			toast.success(`Proxy ${useProxy ? 'enabled' : 'disabled'}`);
			
		} catch (err) {
			console.error('Failed to update proxy settings:', err);
			toast.error('Failed to update proxy settings');
		}
	}
	
	// Handle proxy target change
	async function handleProxyTargetChange(event: Event) {
		const select = event.target as HTMLSelectElement;
		const targetId = select.value;
		
		try {
			// Update the endpoint via API
			await updateEndpoint(endpoint.project_id, endpoint.id, {
				proxy_target_id: targetId
			});
			
			// Update local state
			endpoint = storeUpdateEndpoint('proxy_target_id', targetId, endpoint);
			
			// Find and attach the proxy target object for display
			const target = proxyTargets.find(target => target.id === targetId);
			if (target) {
				endpoint.proxy_target = target;
			}
			
			// Notify parent component
			onChange(endpoint);
			
			// Confirm success
			toast.success('Proxy target updated');
			
		} catch (err) {
			console.error('Failed to update proxy target:', err);
			toast.error('Failed to update proxy target');
		}
	}
</script>

<div class="space-y-4">
	<div class="flex flex-col">
		<div class="flex items-center justify-between mb-4">
			<h3 class="text-lg font-semibold {ThemeUtils.themeTextPrimary()}">Endpoint Proxy</h3>
			
			<label class="relative inline-flex items-center cursor-pointer">
				<input 
					type="checkbox" 
					class="sr-only peer"
					bind:checked={endpoint.use_proxy}
					on:change={handleProxyToggle}
				/>
				<div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
					peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
					rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
					after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
					after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
					after:transition-all peer-checked:bg-blue-600"></div>
				<span class="ml-3 text-sm font-medium {ThemeUtils.themeTextSecondary()}">
					{endpoint.use_proxy ? 'Enabled' : 'Disabled'}
				</span>
			</label>
		</div>
		
		{#if endpoint.use_proxy}
			<div class="space-y-4">
				<div class="mt-2">
					<label for="proxyTarget" class="block mb-2 text-sm font-medium {ThemeUtils.themeTextSecondary()}">
						Select Proxy Target
					</label>
					
					{#if isLoading}
						<div class="animate-pulse flex space-x-4">
							<div class="h-10 bg-gray-700 rounded w-full"></div>
						</div>
					{:else if proxyTargets.length === 0}
						<div class="text-yellow-500 text-sm mb-2">
							No proxy targets available. Configure proxy targets in the project settings first.
						</div>
					{:else}
						<select 
							id="proxyTarget"
							class="{ThemeUtils.themeBgSecondary('border')} {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
							value={endpoint.proxy_target_id || ''}
							on:change={handleProxyTargetChange}
						>
							<option value="" disabled>Select a proxy target</option>
							{#each proxyTargets as target}
								<option value={target.id}>{target.label} ({target.url})</option>
							{/each}
						</select>
						
						{#if endpoint.proxy_target_id && endpoint.proxy_target}
							<div class="mt-2 p-3 bg-gray-750 rounded-lg border border-gray-700">
								<div class="flex flex-col">
									<span class="text-sm {ThemeUtils.themeTextSecondary()}">
										Target: <span class="text-blue-400">{endpoint.proxy_target.label}</span>
									</span>
									<span class="text-sm {ThemeUtils.themeTextSecondary()}">
										URL: <span class="text-blue-400">{endpoint.proxy_target.url}</span>
									</span>
								</div>
							</div>
						{/if}
					{/if}
				</div>
				
				<div class="p-3 bg-gray-750 rounded-lg">
					<p class="text-sm {ThemeUtils.themeTextSecondary()}">
						<i class="fas fa-info-circle mr-2 text-blue-400"></i>
						Requests to this endpoint will be forwarded to the selected proxy target, 
						only when the project is in "Mock" mode. This will not work in other modes.
					</p>
				</div>
			</div>
		{:else}
			<div class="p-3 bg-gray-750 rounded-lg">
				<p class="text-sm {ThemeUtils.themeTextSecondary()}">
					<i class="fas fa-info-circle mr-2 text-blue-400"></i>
					Enable proxy to forward requests for this specific endpoint to a proxy target, 
					only when the project is in "Mock" mode. This will not work in other modes.
				</p>
			</div>
		{/if}
	</div>
</div>
