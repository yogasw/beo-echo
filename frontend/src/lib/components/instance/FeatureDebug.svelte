<script lang="ts">
    import { featureToggles, updateFeatureToggle } from '$lib/stores/featureToggles';
</script>

<!-- Simple component to view and force feature toggle states -->
<div class="p-4 bg-red-500/10 border border-red-500/20 rounded-md mt-4">
    <h3 class="text-md font-medium theme-text-primary mb-2">
        <i class="fas fa-bug mr-2"></i>
        Feature Toggle Debug
    </h3>
    <p class="text-sm theme-text-secondary mb-3">
        This debug panel shows the current state of all feature toggles. Use the buttons to directly update toggles.
    </p>
    
    <div class="space-y-2">
        <div class="flex items-center justify-between p-2 bg-black/10 rounded">
            <span class="font-mono text-sm">showPasswordRequirements:</span>
            <div class="flex items-center space-x-2">
                <span class="text-sm font-semibold" class:text-green-500={$featureToggles.showPasswordRequirements} class:text-red-500={!$featureToggles.showPasswordRequirements}>
                    {$featureToggles.showPasswordRequirements ? 'TRUE' : 'FALSE'}
                </span>
                <button 
                    class="text-xs px-2 py-1 bg-blue-600 text-white rounded hover:bg-blue-700"
                    on:click={() => updateFeatureToggle('showPasswordRequirements', true)}
                >
                    Force On
                </button>
                <button 
                    class="text-xs px-2 py-1 bg-gray-600 text-white rounded hover:bg-gray-700"
                    on:click={() => updateFeatureToggle('showPasswordRequirements', false)}
                >
                    Force Off
                </button>
            </div>
        </div>
    </div>
    
    <div class="mt-4">
        <button 
            class="text-xs px-3 py-1 bg-yellow-500 text-white rounded hover:bg-yellow-600"
            on:click={() => localStorage.removeItem('featureToggles')}
        >
            Clear localStorage
        </button>
        <button 
            class="text-xs px-3 py-1 bg-green-500 text-white rounded hover:bg-green-600 ml-2"
            on:click={() => window.location.reload()}
        >
            Reload Page
        </button>
    </div>
</div>
