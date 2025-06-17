<script lang="ts">
	import { fade } from 'svelte/transition';
	import * as ThemeUtils from '$lib/utils/themeUtils';
	import { 
		type Project, 
		getProjectAdvanceConfig, 
		updateProjectAdvanceConfig
	} from '$lib/api/BeoApi';
	import { toast } from '$lib/stores/toast';
	import { onMount } from 'svelte';

	export let project: Project;

	let isExpanded = false;
	
	// Delay configuration
	let delayMs = 0;
	let isLoadingAdvanceConfig = false;
	let advanceConfig: any = null;
	
	// Function to toggle section expansion
	function toggleSection() {
		isExpanded = !isExpanded;
		if (isExpanded && !advanceConfig) {
			loadAdvanceConfig();
		}
	}
	
	// Load advance configuration
	async function loadAdvanceConfig() {
		if (!project?.id || isLoadingAdvanceConfig) return;
		
		try {
			isLoadingAdvanceConfig = true;
			advanceConfig = await getProjectAdvanceConfig(project.id);
			
			// Extract delay configuration
			if (advanceConfig?.advance_config?.delayMs) {
				delayMs = advanceConfig.advance_config.delayMs;
			} else {
				delayMs = 0;
			}
		} catch (error) {
			console.error('Failed to load advance config:', error);
			// Ignore error for now, config might not exist yet
			advanceConfig = null;
		} finally {
			isLoadingAdvanceConfig = false;
		}
	}
	
	// Handle save of delay configuration
	async function handleDelaySave() {
		if (!project?.id) return;
		
		try {
			await updateProjectAdvanceConfig(project.id, { delayMs });			
			// Reload advance config to reflect changes
			await loadAdvanceConfig();
		} catch (error) {
			console.error('Failed to update delay configuration:', error);
			toast.error(error);
		}
	}

	// Handle delay reset
	async function handleDelayReset() {
		delayMs = 0;
		handleDelayChange();
	}
	
	// Handle delay input change
	function handleDelayChange() {
		handleDelaySave();
	}
	
	// Validate delay value
	function validateDelay() {
		if (delayMs < 0) {
			delayMs = 0;
		} else if (delayMs > 120000) {
			delayMs = 120000;
		}
	}
	
	// Convert delay to seconds for display
	$: delaySeconds = delayMs / 1000;
	
	// Load advance config on mount if expanded
	onMount(() => {
		if (isExpanded) {
			loadAdvanceConfig();
		}
	});
</script>

<div class={ThemeUtils.card('overflow-hidden')}>
	<div 
		class="flex justify-between items-center p-3 hover:bg-gray-100 dark:hover:bg-gray-700 cursor-pointer bg-gray-100 dark:bg-gray-750"
		on:click={toggleSection}
		on:keydown={(e) => e.key === 'Enter' && toggleSection()}
		tabindex="0"
		role="button"
		aria-expanded={isExpanded}
	>
		<div class="flex items-center">
			<div class="bg-purple-600/10 p-1.5 rounded mr-2">
				<i class="fas fa-sliders-h text-purple-500"></i>
			</div>
			<h3 class="font-medium theme-text-primary">Advanced Configuration</h3>
		</div>
		<i class="fas {isExpanded ? 'fa-chevron-up' : 'fa-chevron-down'} theme-text-muted"></i>
	</div>
	
	{#if isExpanded}
		<div transition:fade={{ duration: 150 }} class="border-t theme-border p-4">
			<div class="bg-gray-50 dark:bg-gray-900/50 rounded-lg p-4 border theme-border">
				<div class="flex items-start justify-between mb-4">
					<div class="flex-1">
						<label for="delay-ms" class="text-sm font-medium theme-text-secondary block mb-1">Response Delay (milliseconds)</label>
						<p class="text-xs theme-text-muted">
							Add artificial delay to all responses from this project. Set to 0 to disable delay.
						</p>
						{#if isLoadingAdvanceConfig}
							<div class="flex items-center mt-2">
								<i class="fas fa-spinner fa-spin text-blue-400 mr-2"></i>
								<span class="text-xs theme-text-muted">Loading configuration...</span>
							</div>
						{/if}
					</div>
					<button
						type="button"
						class="ml-4 px-3 py-1 text-xs bg-red-600 hover:bg-red-700 text-white rounded-md flex items-center"
						on:click={handleDelayReset}
						title="Reset delay to 0"
						aria-label="Reset delay configuration"
					>
						<i class="fas fa-undo mr-1"></i>
						Reset
					</button>
				</div>
				
				<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
					<div>
						<div class="relative">
							<div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
								<i class="fas fa-clock theme-text-muted"></i>
							</div>
							<input
								type="number"
								id="delay-ms"
								min="0"
								max="120000"
								step="1000"
								class={ThemeUtils.inputField()}
								bind:value={delayMs}
								on:input={validateDelay}
								on:change={handleDelayChange}
								placeholder="Enter delay in milliseconds (0 = disabled)"
								title="Response delay in milliseconds (0-120000ms / 2 minutes max)"
								aria-label="Response delay in milliseconds"
							/>
						</div>
						<p class="text-xs theme-text-muted mt-1">
							Current: {delaySeconds}s (Max: 120 seconds)
						</p>
					</div>
					
					<div class="flex items-center">
						<div class="bg-orange-50 dark:bg-orange-900/20 border border-orange-200 dark:border-orange-800 rounded-lg p-3 flex-1">
							<div class="flex items-center">
								<i class="fas fa-info-circle text-orange-500 mr-2"></i>
								<div class="text-xs theme-text-secondary">
									<p class="font-medium">Delay Priority:</p>
									<p>Response → Endpoint → Project</p>
									<p class="mt-1 text-xs theme-text-muted">
										This project-level delay applies to all endpoints unless overridden at the endpoint or response level.
									</p>
								</div>
							</div>
						</div>
					</div>
				</div>
			</div>
		</div>
	{/if}
</div>
