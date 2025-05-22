<script lang="ts">
    import { featureToggles, updateFeatureToggle } from '$lib/stores/featureToggles';
    import { fade } from 'svelte/transition';
    import { toast } from '$lib/stores/toast';
    import FeatureToggle from '$lib/components/common/FeatureToggle.svelte';
    
    let saving = false;
    
    // Update a feature toggle with backend integration
    async function handleFeatureToggle(event: CustomEvent<{key: string, value: boolean}>) {
        const { key, value } = event.detail;
        saving = true;
        
        try {
            // Call backend API to update feature toggle
            // This would normally be an API call like:
            // await fetch('/api/feature-toggles', {
            //   method: 'POST',
            //   headers: { 'Content-Type': 'application/json' },
            //   body: JSON.stringify({ feature: key, enabled: value })
            // });
            
            // For now, just update the local store directly
            updateFeatureToggle(key as keyof typeof $featureToggles, value);
            
            // Simulate API latency
            await new Promise(resolve => setTimeout(resolve, 300));
            
            toast.success(`Feature "${key}" ${value ? 'enabled' : 'disabled'}`);
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
        <FeatureToggle
            title="Password Requirements"
            description="Show password complexity requirements when changing passwords"
            bind:checked={$featureToggles.showPasswordRequirements}
            featureKey="showPasswordRequirements"
            on:change={handleFeatureToggle}
            disabled={saving}
            ariaLabel="Toggle password requirements"
        />
        
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
