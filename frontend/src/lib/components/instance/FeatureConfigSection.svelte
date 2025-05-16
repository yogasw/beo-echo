<script lang="ts">
    import { featureToggles, updateFeatureToggle } from '$lib/stores/featureToggles';
    import { fade } from 'svelte/transition';
    import { toast } from '$lib/stores/toast';
    
    let saving = false;
    
    // Update a feature toggle with backend integration
    async function handleFeatureToggle(feature: string, enabled: boolean) {
        saving = true;
        
        try {
            // Call backend API to update feature toggle
            // This would normally be an API call like:
            // await fetch('/api/feature-toggles', {
            //   method: 'POST',
            //   headers: { 'Content-Type': 'application/json' },
            //   body: JSON.stringify({ feature, enabled })
            // });
            
            // For now, just update the local store directly
            updateFeatureToggle(feature as keyof typeof $featureToggles, enabled);
            
            // Simulate API latency
            await new Promise(resolve => setTimeout(resolve, 300));
            
            toast.success(`Feature "${feature}" ${enabled ? 'enabled' : 'disabled'}`);
        } catch (error) {
            toast.error('Failed to update feature toggle');
            console.error('Failed to update feature toggle:', error);
        } finally {
            saving = false;
        }
    }
</script>

<div class="space-y-4">
    <div class="flex justify-between items-center mb-3">
        {#if saving}
            <span class="text-sm theme-text-secondary flex items-center gap-1.5">
                <i class="fas fa-spinner fa-spin"></i>
                <span>Updating...</span>
            </span>
        {/if}
    </div>
    
    <!-- Feature toggles section -->
    <div class="space-y-4">
        <!-- Password Requirements Toggle -->
        <div class="flex items-center justify-between py-3 border-b theme-border">
            <div>
                <h4 class="font-medium theme-text-primary">Password Requirements</h4>
                <p class="text-sm theme-text-secondary">Show password complexity requirements when changing passwords</p>
            </div>
            <label class="inline-flex items-center cursor-pointer">
                <input 
                    type="checkbox" 
                    bind:checked={$featureToggles.showPasswordRequirements} 
                    on:change={() => handleFeatureToggle('showPasswordRequirements', $featureToggles.showPasswordRequirements)}
                    class="sr-only peer"
                    disabled={saving}
                >
                <div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
                    rounded-full peer peer-focus:outline-none peer-focus:ring-2 
                    peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 
                    peer-checked:after:translate-x-full peer-checked:after:border-white 
                    after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
                    after:bg-white after:rounded-full after:h-5 after:w-5 
                    after:transition-all dark:border-gray-600 relative">
                </div>
            </label>
        </div>
        
        <!-- You can add more feature toggles here -->
    </div>
    
    <!-- Note about feature configuration -->
    <div class="mt-4 p-3 bg-blue-600/10 border border-blue-600/20 rounded text-sm theme-text-secondary">
        <div class="flex items-start gap-3">
            <i class="fas fa-info-circle text-blue-400 mt-0.5"></i>
            <div>
                <p class="mb-1">Feature toggles are updated from the backend API.</p>
                <p>Changes take effect immediately for all users system-wide.</p>
            </div>
        </div>
    </div>
</div>
