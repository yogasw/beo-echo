<script lang="ts">
	import {
		getProxyTargets,
		updateEndpoint,
		type ProxyTarget,
		type Endpoint
	} from '$lib/api/BeoApi';
	import { onMount } from 'svelte';
	import { toast } from '$lib/stores/toast';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { updateEndpoint as storeUpdateEndpoint } from '$lib/stores/saveButton';
	import { selectedProject } from '$lib/stores/selectedConfig';
	import ToggleSwitch from '$lib/components/common/ToggleSwitch.svelte';

	export let endpoint: Endpoint;
	export let isLoading: boolean = false;
	export let onChange: (endpoint: Endpoint) => void;
	export let use_proxy: boolean = false; // Bind to parent component
	let proxyTargetId: string | null = null;
	let proxyTargets = $selectedProject?.proxy_targets || [];
	let proxyTarget: ProxyTarget | null = null;

	$: {
		if (endpoint) {
			// Initialize our exported use_proxy from the endpoint
			use_proxy = endpoint.use_proxy || false;
			proxyTargetId = endpoint.proxy_target_id || null;
			proxyTarget = proxyTargets.find((target) => target.id === proxyTargetId) || null;
		}
	}
	// Handle proxy toggle
	async function handleProxyToggle() {
		console.log('Proxy toggle changed:', use_proxy);

		// Update local state
		endpoint = storeUpdateEndpoint('use_proxy', use_proxy, endpoint);
		console.log('Updated endpoint:', endpoint);
		// Notify parent component
		onChange(endpoint);
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
			const target = proxyTargets.find((target) => target.id === targetId);
			if (target) {
				endpoint.proxy_target = target;
			}

			// Notify parent component
			onChange(endpoint);
		} catch (err) {
			console.error('Failed to update proxy target:', err);
			toast.error('Failed to update proxy target');
		}
	}
</script>

<div class="space-y-4">
	<div class="flex flex-col">
		<div class="flex items-center justify-between">
			<h3 class="text-lg font-semibold {ThemeUtils.themeTextPrimary()}">Endpoint Proxy</h3>

			<div class="flex items-center">
				<ToggleSwitch 
					bind:checked={use_proxy}
					on:change={handleProxyToggle}
					aria-label="Toggle proxy"
					label="Use Proxy"
				/>
			</div>
		</div>

		{#if use_proxy}
			<div class="space-y-4">
				<div class="mt-2">
					<label
						for="proxyTarget"
						class="block mb-2 text-sm font-medium {ThemeUtils.themeTextSecondary()}"
					>
						Select Proxy Target
					</label>

					{#if isLoading}
						<div class="animate-pulse flex space-x-4">
							<div class="h-10 {ThemeUtils.themeBgSecondary()} rounded w-full"></div>
						</div>
					{:else if proxyTargets.length === 0}
						<div class="text-yellow-500 text-sm mb-2">
							No proxy targets available. Configure proxy targets in the project settings first.
						</div>
					{:else}
						<select
							id="proxyTarget"
							class="{ThemeUtils.themeBgSecondary(
								'border'
							)} {ThemeUtils.themeBorder()} {ThemeUtils.themeTextPrimary()} text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block w-full p-2.5"
							value={proxyTargetId || ''}
							on:change={handleProxyTargetChange}
						>
							<option value="" disabled>Select a proxy target</option>
							{#each proxyTargets as target}
								<option value={target.id}>{target.label} ({target.url})</option>
							{/each}
						</select>

						{#if proxyTargetId && proxyTarget}
							<div class="mt-2 p-3 {ThemeUtils.themeBgSecondary()} rounded-lg border {ThemeUtils.themeBorder()}">
								<div class="flex flex-col">
									<span class="text-sm {ThemeUtils.themeTextSecondary()}">
										Target: <span class="text-blue-400">{proxyTarget.label}</span>
									</span>
									<span class="text-sm {ThemeUtils.themeTextSecondary()}">
										URL: <span class="text-blue-400">{proxyTarget.url}</span>
									</span>
								</div>
							</div>
						{/if}
					{/if}
				</div>
				<div class="p-3 {ThemeUtils.themeBgSecondary()} rounded-lg">
					<p class="text-sm {ThemeUtils.themeTextSecondary()}">
						<i class="fas fa-info-circle mr-2 text-blue-400"></i>
						Requests to this endpoint will be forwarded to the selected proxy target, only when the project
						is in "Mock" mode. This will not work in other modes.
					</p>
				</div>
			</div>
		{/if}
	</div>
</div>
