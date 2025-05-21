<script lang="ts">
    import { fade } from 'svelte/transition';
    import ToggleSwitch from './ToggleSwitch.svelte';
    
    // Props for the component
    export let title: string;
    export let description: string;
    export let checked: boolean;
    export let disabled: boolean = false;
    export let featureKey: string = '';
    export let ariaLabel: string = '';
    
    // Events
    import { createEventDispatcher } from 'svelte';
    const dispatch = createEventDispatcher<{
        change: { key: string; value: boolean }
    }>();
    
    // Handle toggle change
    function handleChange() {
        dispatch('change', { key: featureKey, value: checked });
    }
</script>

<div class="flex items-center justify-between py-3 border-b theme-border" transition:fade={{ duration: 200 }}>
    <div>
        <h4 class="font-medium theme-text-primary">{title}</h4>
        <p class="text-sm theme-text-secondary">{description}</p>
    </div>
    <ToggleSwitch 
        bind:checked={checked}
        on:change={handleChange}
        {disabled}
        ariaLabel={ariaLabel || title}
    />
</div>
