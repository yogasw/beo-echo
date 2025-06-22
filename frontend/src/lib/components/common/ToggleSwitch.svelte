<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    /**
     * A reusable toggle switch component that follows the project's design system.
     * 
     * Features:
     * - Consistent styling with theme support (dark/light mode)
     * - Supports disabled state
     * - Fully typed with TypeScript
     * - Customizable size
     */
    
    // Toggle state - default to false if not provided
    export let checked = false;
    
    // Optional parameters
    export let disabled = false;
    export let size: 'small' | 'default' | 'large' = 'default';
    export let name = '';
    export let id = '';
    export let ariaLabel = '';
    export let title = '';
    
    // Create event dispatcher
    const dispatch = createEventDispatcher<{
        change: { checked: boolean };
    }>();
    
    // Calculate size classes based on the size prop
    $: sizeClass = size === 'small' 
        ? 'w-8 h-4 after:h-3 after:w-3'
        : size === 'large'
            ? 'w-14 h-7 after:h-6 after:w-6'
            : 'w-11 h-6 after:h-5 after:w-5';
    
    // Handle change event - create a custom event with the checked state
    function handleChange() {
        // Dispatch the change event using Svelte's event dispatcher
        // No need to access event.target, as bind:checked takes care of updating the value
        dispatch('change', { checked });
    }
</script>

<label class="inline-flex items-center cursor-pointer {disabled ? 'opacity-60 cursor-not-allowed' : ''}">
    <input 
        type="checkbox" 
        bind:checked
        {disabled}
        {name}
        {id}
        {title}
        aria-label={ariaLabel}
        on:change={handleChange}
        class="sr-only peer"
        {...$$restProps} 
    />
    <div class="{sizeClass} bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
        rounded-full peer peer-focus:outline-none peer-focus:ring-2 
        peer-focus:ring-blue-300 dark:peer-focus:ring-blue-600 
        peer-checked:after:translate-x-full peer-checked:after:border-white 
        after:content-[''] after:absolute after:top-[2px] after:start-[2px] 
        after:bg-white after:rounded-full 
        after:transition-all dark:border-gray-600 relative">
    </div>
    
    {#if $$slots.default}
        <div class="ml-3">
            <slot></slot>
        </div>
    {/if}
</label>
