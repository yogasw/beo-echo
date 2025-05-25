<script lang="ts">
    import { createEventDispatcher } from 'svelte';
    
    export let workspaceId: string;
    export let workspaceName: string = '';
    export let compact: boolean = false;
    export let isOwner: boolean = false;
    
    const dispatch = createEventDispatcher<{
        click: { workspaceId: string, workspaceName: string }
    }>();
    
    function handleClick() {
        dispatch('click', { workspaceId, workspaceName });
    }
</script>

{#if compact}
    <!-- Compact view for table rows -->
    <button 
        on:click={handleClick} 
        class="p-2 theme-bg-secondary rounded-full hover:bg-blue-500/20" 
        title="Auto-Invite Settings"
        aria-label="Auto-Invite Settings"
        disabled={!isOwner}
    >
        <i class="fas fa-user-plus theme-text-secondary"></i>
    </button>
{:else}
    <!-- Full view for dedicated links -->
    <button 
        on:click={handleClick}
        class="flex items-center space-x-2 text-sm px-3 py-2 theme-bg-secondary theme-text-primary 
        rounded-md hover:bg-blue-500/20 disabled:opacity-50 disabled:cursor-not-allowed"
        disabled={!isOwner}
        aria-label="Configure auto-invite settings for workspace"
        title="Configure auto-invite settings for workspace"
    >
        <i class="fas fa-user-plus"></i>
        <span>Auto-Invite Settings</span>
    </button>
{/if}
