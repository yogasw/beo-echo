<script lang="ts">
  import { onMount } from 'svelte';
  import { getAllSystemConfigs, updateSystemConfig } from '$lib/api/BeoApi';
  import { featureToggles, updateFeatureToggle } from '$lib/stores/featureToggles';
  import { toast } from '$lib/stores/toastStore';
  
  // Feature descriptions for UI
  const featureDescriptions = {
    showPasswordRequirements: 'Show password requirements when users create or update passwords',
    emailUpdatesEnabled: 'Allow users to update their email addresses'
  };

  // Track loading and error states
  let loading = true;
  let error: string | null = null;
  let features: { key: string; name: string; enabled: boolean; description: string }[] = [];
  
  // Load feature flags from the backend
  async function loadFeatureFlags() {
    try {
      loading = true;
      error = null;
      
      // Get all system configs
      const configs = await getAllSystemConfigs();
      
      // Filter out feature flags (starting with feature_ or FEATURE_)
      const featureConfigs = configs.filter(
        config => config.key.toLowerCase().startsWith('feature_') || 
                  config.key.startsWith('FEATURE_')
      );
      
      // Transform to our local format
      features = featureConfigs.map(config => {
        // Convert key to camelCase for frontend use
        const featureName = convertToCamelCase(config.key);
        
        return {
          key: config.key,
          name: featureName,
          enabled: config.value === 'true',
          description: featureDescriptions[featureName] || config.description || `${featureName} feature flag`
        };
      });
      
      // Also update the feature toggle store with these values
      features.forEach(feature => {
        updateFeatureToggle(feature.name, feature.enabled);
      });
      
    } catch (err) {
      console.error('Failed to load feature flags:', err);
      error = 'Failed to load feature flags. Please try again.';
    } finally {
      loading = false;
    }
  }
  
  // Function to toggle a feature flag
  async function toggleFeature(feature) {
    try {
      const newValue = !feature.enabled;
      
      // Update in backend
      await updateSystemConfig(feature.key, newValue ? 'true' : 'false');
      
      // Update local state
      feature.enabled = newValue;
      
      // Update feature toggle store
      updateFeatureToggle(feature.name, newValue);
      
      // Show success message
      toast.success(`${feature.name} ${newValue ? 'enabled' : 'disabled'} successfully`);
      
    } catch (err) {
      console.error(`Failed to update ${feature.name}:`, err);
      toast.error(`Failed to update ${feature.name}. Please try again.`);
    }
  }
  
  // Helper function to convert feature keys from DB format to frontend camelCase
  function convertToCamelCase(key: string): string {
    // Remove the "feature_" or "FEATURE_" prefix
    key = key.replace(/^(feature_|FEATURE_)/i, '');
    
    // Convert to lowercase
    key = key.toLowerCase();
    
    // Split by underscore
    const parts = key.split('_');
    
    // Convert to camelCase
    for (let i = 1; i < parts.length; i++) {
      if (parts[i].length > 0) {
        parts[i] = parts[i].charAt(0).toUpperCase() + parts[i].slice(1);
      }
    }
    
    return parts.join('');
  }
  
  // Load feature flags when component mounts
  onMount(() => {
    loadFeatureFlags();
  });
</script>

<div class="w-full bg-gray-800 p-4 rounded-lg">
  <div class="flex justify-between items-center mb-4">
    <h2 class="text-xl font-bold text-white flex items-center">
      <i class="fas fa-toggle-on mr-2 text-blue-500"></i>
      Feature Flags
    </h2>
    <button 
      on:click={loadFeatureFlags} 
      class="bg-gray-700 hover:bg-gray-600 text-white py-1 px-3 rounded flex items-center text-sm"
      disabled={loading}
    >
      <i class="fas fa-sync-alt mr-2"></i>
      Refresh
    </button>
  </div>
  
  {#if error}
    <div class="bg-red-600/20 border border-red-600 text-red-300 p-3 rounded-md mb-4">
      <i class="fas fa-exclamation-circle mr-2"></i>
      {error}
    </div>
  {/if}
  
  {#if loading}
    <div class="flex justify-center items-center p-8">
      <div class="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-blue-500"></div>
    </div>
  {:else if features.length === 0}
    <div class="bg-gray-700 p-4 rounded-md text-gray-300">
      <i class="fas fa-info-circle mr-2"></i>
      No feature flags found.
    </div>
  {:else}
    <div class="space-y-3">
      {#each features as feature (feature.key)}
        <div class="bg-gray-700 p-4 rounded-md shadow-md">
          <div class="flex justify-between items-center">
            <div>
              <h3 class="font-semibold text-white">{feature.name}</h3>
              <p class="text-sm text-gray-300">{feature.description}</p>
              <div class="text-xs text-gray-400 mt-1">Key: {feature.key}</div>
            </div>
            <label class="inline-flex items-center cursor-pointer">
              <input 
                type="checkbox" 
                checked={feature.enabled} 
                on:change={() => toggleFeature(feature)} 
                class="sr-only peer"
              />
              <div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
                peer-focus:ring-blue-800 rounded-full peer 
                peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full 
                peer-checked:after:border-white after:content-[''] after:absolute 
                after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 
                after:border after:rounded-full after:h-5 after:w-5 after:transition-all 
                peer-checked:bg-blue-600 border border-gray-600"
              ></div>
              <span class="ml-3 text-sm font-medium text-gray-300">
                {feature.enabled ? 'Enabled' : 'Disabled'}
              </span>
            </label>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
