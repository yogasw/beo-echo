<script lang="ts">
  import type { Workspace } from '$lib/types/User';
  
  // Props
  export let workspace: Workspace;
  export let onConfirm: (confirmed: boolean) => void;
  
  // Handle delete confirmation
  function confirmDelete() {
    onConfirm(true);
  }
  
  // Handle cancel
  function cancelDelete() {
    onConfirm(false);
  }
</script>

<!-- Modal Backdrop -->
<div 
  class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center"
  on:click|self={cancelDelete}
>
  <!-- Modal Content -->
  <div class="theme-bg-primary rounded-lg shadow-xl w-full max-w-md mx-4 overflow-hidden">
    <!-- Modal Header -->
    <div class="theme-bg-secondary px-4 py-3 flex items-center justify-between border-b theme-border">
      <h3 class="theme-text-primary font-semibold">Delete Workspace</h3>
      <button
        class="theme-text-secondary hover:theme-text-primary"
        on:click={cancelDelete}
      >
        <i class="fas fa-times"></i>
      </button>
    </div>
    
    <!-- Modal Body -->
    <div class="p-5 theme-text-primary">
      <div class="flex items-center mb-4 text-red-500">
        <div class="rounded-full bg-red-500/10 p-3 mr-3">
          <i class="fas fa-exclamation-triangle text-xl"></i>
        </div>
        <div>
          <h4 class="font-semibold">Warning: This action cannot be undone</h4>
        </div>
      </div>
      
      <p class="mb-4 theme-text-secondary">
        Are you sure you want to delete the workspace <span class="font-semibold theme-text-primary">{workspace.name}</span>?
      </p>
      
      <p class="theme-text-muted text-sm bg-red-500/10 p-3 rounded-md border border-red-500/20">
        This will permanently delete all data related to this workspace including projects, routes, and logs.
      </p>
    </div>
    
    <!-- Modal Footer -->
    <div class="px-5 py-4 theme-bg-secondary border-t theme-border flex justify-end space-x-3">
      <button
        class="px-4 py-2 theme-text-primary theme-bg-secondary border theme-border rounded-md hover:bg-gray-600"
        on:click={cancelDelete}
      >
        Cancel
      </button>
      
      <button
        class="px-4 py-2 text-white bg-red-600 rounded-md hover:bg-red-700"
        on:click={confirmDelete}
      >
        Delete Workspace
      </button>
    </div>
  </div>
</div>
