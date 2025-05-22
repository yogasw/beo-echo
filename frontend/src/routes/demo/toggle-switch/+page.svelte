<script lang="ts">
    import ToggleSwitch from '$lib/components/common/ToggleSwitch.svelte';
    import { theme, toggleTheme } from '$lib/stores/theme';
    
    let isEnabled = false;
    let isOption1Enabled = true;
    let isOption2Enabled = false;
    
    function handleToggleChange(event: CustomEvent) {
        console.log('Toggle changed:', event.detail.checked);
    }
</script>

<div class="theme-bg-primary p-6 rounded-lg shadow-md space-y-8">
    <h1 class="text-xl font-bold theme-text-primary mb-6">Toggle Switch Component Demo</h1>
    
    <!-- Basic usage -->
    <div class="space-y-4">
        <h2 class="text-lg font-semibold theme-text-primary">Basic Usage</h2>
        
        <div class="flex items-center justify-between py-3 border-b theme-border">
            <div>
                <h4 class="font-medium theme-text-primary">Simple Toggle</h4>
                <p class="text-sm theme-text-secondary">Toggle switch with default styling</p>
            </div>
            <ToggleSwitch 
                bind:checked={isEnabled} 
                on:change={handleToggleChange}
            />
        </div>
        
        <div class="text-sm theme-text-secondary mt-2">
            Current state: <span class="font-semibold">{isEnabled ? 'Enabled' : 'Disabled'}</span>
        </div>
    </div>
    
    <!-- Different sizes -->
    <div class="space-y-4">
        <h2 class="text-lg font-semibold theme-text-primary">Size Variants</h2>
        
        <div class="space-y-6">
            <div class="flex items-center justify-between py-3 border-b theme-border">
                <span class="font-medium theme-text-primary">Small Size</span>
                <ToggleSwitch size="small" />
            </div>
            
            <div class="flex items-center justify-between py-3 border-b theme-border">
                <span class="font-medium theme-text-primary">Default Size</span>
                <ToggleSwitch size="default" />
            </div>
            
            <div class="flex items-center justify-between py-3 border-b theme-border">
                <span class="font-medium theme-text-primary">Large Size</span>
                <ToggleSwitch size="large" />
            </div>
        </div>
    </div>
    
    <!-- Disabled state -->
    <div class="space-y-4">
        <h2 class="text-lg font-semibold theme-text-primary">Disabled State</h2>
        
        <div class="flex items-center justify-between py-3 border-b theme-border">
            <div>
                <h4 class="font-medium theme-text-primary">Disabled Toggle</h4>
                <p class="text-sm theme-text-secondary">Toggle switch in disabled state</p>
            </div>
            <ToggleSwitch disabled={true} checked={true} />
        </div>
    </div>
    
    <!-- With slot content -->
    <div class="space-y-4">
        <h2 class="text-lg font-semibold theme-text-primary">With Label in Slot</h2>
        
        <div class="py-3 border-b theme-border space-y-4">
            <ToggleSwitch bind:checked={isOption1Enabled}>
                <span class="font-medium theme-text-primary">Enable feature 1</span>
                <p class="text-sm theme-text-secondary">This enables the first feature</p>
            </ToggleSwitch>
            
            <ToggleSwitch bind:checked={isOption2Enabled}>
                <span class="font-medium theme-text-primary">Enable feature 2</span>
                <p class="text-sm theme-text-secondary">This enables the second feature</p>
            </ToggleSwitch>
        </div>
    </div>
    
    <!-- Real world example: Theme toggler -->
    <div class="space-y-4">
        <h2 class="text-lg font-semibold theme-text-primary">Theme Toggle Example</h2>
        
        <div class="flex items-center justify-between py-3 border-b theme-border">
            <div>
                <h4 class="font-medium theme-text-primary">Dark Mode</h4>
                <p class="text-sm theme-text-secondary">Toggle between light and dark themes</p>
            </div>
            <div class="flex items-center">
                <i class="fas fa-sun text-yellow-400 mr-2"></i>
                <ToggleSwitch 
                    checked={$theme === 'dark'} 
                    on:change={toggleTheme}
                />
                <i class="fas fa-moon text-blue-400 ml-2"></i>
            </div>
        </div>
    </div>
    
    <!-- Code usage examples -->
    <div class="space-y-4 mt-8">
        <h2 class="text-lg font-semibold theme-text-primary">Code Examples</h2>
        
        <div class="bg-gray-900 rounded-md p-4 text-sm font-mono text-gray-200 overflow-x-auto">
            <pre>{`<!-- Basic usage -->
<ToggleSwitch bind:checked={isEnabled} />

<!-- With size -->
<ToggleSwitch size="small" bind:checked={isEnabled} />
<ToggleSwitch size="default" bind:checked={isEnabled} />
<ToggleSwitch size="large" bind:checked={isEnabled} />

<!-- Disabled state -->
<ToggleSwitch disabled={true} />

<!-- With event handling -->
<ToggleSwitch bind:checked={isEnabled} on:change={handleChange} />

<!-- With slot content (label to the right) -->
<ToggleSwitch bind:checked={isEnabled}>
  <span>Label text goes here</span>
</ToggleSwitch>`}</pre>
        </div>
    </div>
</div>
