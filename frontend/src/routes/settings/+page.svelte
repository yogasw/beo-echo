<script lang="ts">
  import { page } from '$app/stores';

  // Get state from parent layout
  export let selectedConfig = null;
  export let activeTab = 'settings';
</script>

<div class="p-4">
  <h1 class="text-2xl font-bold mb-4">Settings</h1>
  <div class="bg-gray-800 p-4 rounded-lg">
    <h2 class="text-xl font-semibold mb-2">Application Settings</h2>
    <div class="space-y-4">
      <div class="flex items-center justify-between">
        <span>Dark Mode</span>
        <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md">
          Toggle
        </button>
      </div>
      <div class="flex items-center justify-between">
        <span>Auto-save</span>
        <button class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md">
          Enable
        </button>
      </div>
    </div>
  </div>
</div> 