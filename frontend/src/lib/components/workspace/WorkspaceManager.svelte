<script lang="ts">
  import { onMount } from 'svelte';
  import { workspaces, currentWorkspace, workspaceStore } from '$lib/stores/workspace';
  import WorkspaceDeleteConfirmation from './WorkspaceDeleteConfirmation.svelte';
  import { toast } from '$lib/stores/toast';
  
  export let className = '';
  
  // Local state
  let loading = false;
  let error: string | null = null;
  let modalOpen = false;
  let showDeleteConfirmation = false;
  let workspaceToDelete = null;
  
  // Handle modal toggle
  function toggleModal() {
    modalOpen = !modalOpen;
  }
  
  // Open delete confirmation dialog
  function openDeleteConfirmation(workspace) {
    workspaceToDelete = workspace;
    showDeleteConfirmation = true;
    modalOpen = false;
  }
  
  // Handle the deletion confirmation result
  async function handleDeleteConfirmation(confirmed: boolean) {
    if (confirmed && workspaceToDelete) {
      try {
        loading = true;
        await workspaces.delete(workspaceToDelete.id);
        toast.success(`Workspace "${workspaceToDelete.name}" deleted successfully`);
      } catch (err) {
        toast.error(`Failed to delete workspace: ${err.message || 'Unknown error'}`);
      } finally {
        loading = false;
      }
    }
    
    showDeleteConfirmation = false;
    workspaceToDelete = null;
  }
</script>

<!-- Workspace Manager Component -->
<div class="relative {className}">
  <!-- Company/Workspace Button -->
  <button 
    on:click={toggleModal}
    class="flex flex-col items-center"
  >
    <div class="w-12 aspect-square theme-bg-secondary theme-text-primary p-3 rounded-full border-2 border-green-500 flex items-center justify-center">
      <i class="fas fa-building"></i>
    </div>
    <span class="text-xs mt-1 theme-text-primary">Workspace</span>
  </button>
  
  <!-- Workspace Modal -->
  {#if modalOpen}
    <div class="absolute top-full right-0 mt-2 w-64 theme-bg-primary rounded-md shadow-lg z-40 border theme-border">
      <div class="p-2">
        <h3 class="theme-text-secondary text-sm font-medium px-3 py-2">Current Workspace</h3>
        
        <!-- Current Workspace Info -->
        {#if $currentWorkspace}
          <div class="p-3 theme-bg-secondary rounded-md mx-2 mb-2">
            <div class="flex items-center justify-between">
              <span class="theme-text-primary font-semibold">{$currentWorkspace.name}</span>
              
              <!-- Only allow deletion if user has appropriate role -->
              {#if $currentWorkspace.role === 'owner' || $currentWorkspace.role === 'admin'}
                <button 
                  on:click={() => openDeleteConfirmation($currentWorkspace)}
                  class="text-red-500 hover:text-red-400 p-1"
                  title="Delete workspace"
                >
                  <i class="fas fa-trash"></i>
                </button>
              {/if}
            </div>
            
            {#if $currentWorkspace.role}
              <div class="theme-text-muted text-xs mt-1">
                Role: <span class="theme-text-secondary">{$currentWorkspace.role}</span>
              </div>
            {/if}
            
            {#if $currentWorkspace.description}
              <div class="theme-text-muted text-xs mt-1 line-clamp-2">
                {$currentWorkspace.description}
              </div>
            {/if}
          </div>
        {:else}
          <div class="theme-text-muted text-sm px-3 py-2">
            No workspace selected
          </div>
        {/if}
        
        <!-- Footer with action links -->
        <div class="mt-2 pt-2 border-t theme-border">
          <button 
            class="flex items-center w-full px-3 py-2 text-left hover:bg-blue-500/20 rounded-md transition-colors"
            on:click={() => {
              modalOpen = false;
              document.getElementById('workspaceSelectorButton')?.click();
            }}
          >
            <i class="fas fa-exchange-alt mr-2 text-xs theme-text-secondary"></i>
            <span class="theme-text-primary text-sm">Switch Workspace</span>
          </button>
        </div>
      </div>
    </div>
  {/if}
</div>

<!-- Delete Confirmation Modal -->
{#if showDeleteConfirmation && workspaceToDelete}
  <WorkspaceDeleteConfirmation 
    workspace={workspaceToDelete}
    onConfirm={handleDeleteConfirmation}
  />
{/if}
