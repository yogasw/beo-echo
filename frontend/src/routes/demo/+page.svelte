<!-- Demo index page -->
<script lang="ts">
  import { page } from '$app/stores';
  
  const demoPages = [
    {
      title: 'Component Testing',
      description: 'Test enhanced StatusCodeInput and HttpMethodDropdown components with dynamic padding and multi-criteria search features',
      href: '/demo/components',
      icon: 'fas fa-cogs',
      tags: ['Components', 'Search', 'Padding', 'UX']
    },
    {
      title: 'Toggle Switch',
      description: 'Interactive toggle switch component demonstrations',
      href: '/demo/toggle-switch',
      icon: 'fas fa-toggle-on',
      tags: ['Toggle', 'Switch', 'Interactive']
    }
  ];
</script>

<svelte:head>
  <title>Demo Pages - Beo Echo</title>
</svelte:head>

<div class="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
  <div class="max-w-6xl mx-auto px-4">
    <div class="mb-8">
      <h1 class="text-4xl font-bold text-gray-900 dark:text-white mb-4">
        <i class="fas fa-flask mr-3 text-blue-600"></i>Demo Pages
      </h1>
      <p class="text-xl text-gray-600 dark:text-gray-400">
        Interactive demonstrations of Beo Echo components and features
      </p>
    </div>

    <!-- Demo Pages Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each demoPages as demo}
        <a 
          href={demo.href}
          class="group bg-white dark:bg-gray-800 rounded-lg shadow-md border border-gray-200 dark:border-gray-700 hover:shadow-lg hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-200 p-6 block"
        >
          <div class="flex items-start">
            <div class="flex-shrink-0">
              <div class="w-12 h-12 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                <i class="{demo.icon} text-blue-600 dark:text-blue-400 text-xl"></i>
              </div>
            </div>
            <div class="ml-4 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 dark:text-white group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors">
                {demo.title}
              </h3>
              <p class="mt-2 text-sm text-gray-600 dark:text-gray-400 line-clamp-3">
                {demo.description}
              </p>
              
              <!-- Tags -->
              <div class="mt-3 flex flex-wrap gap-2">
                {#each demo.tags as tag}
                  <span class="px-2 py-1 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 text-xs rounded-md">
                    {tag}
                  </span>
                {/each}
              </div>
            </div>
          </div>
          
          <!-- Arrow Icon -->
          <div class="mt-4 flex justify-end">
            <i class="fas fa-arrow-right text-gray-400 group-hover:text-blue-600 dark:group-hover:text-blue-400 transition-colors"></i>
          </div>
        </a>
      {/each}
    </div>

    <!-- Add New Demo Placeholder -->
    <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow-md border-2 border-dashed border-gray-300 dark:border-gray-600 p-8 text-center">
      <div class="w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-full flex items-center justify-center mx-auto mb-4">
        <i class="fas fa-plus text-gray-400 text-2xl"></i>
      </div>
      <h3 class="text-lg font-medium text-gray-900 dark:text-white mb-2">More Demos Coming Soon</h3>
      <p class="text-gray-600 dark:text-gray-400 text-sm">
        Additional component demonstrations and feature showcases will be added here
      </p>
    </div>

    <!-- Navigation -->
    <div class="mt-8 text-center">
      <a 
        href="/home" 
        class="inline-flex items-center px-6 py-3 bg-blue-600 hover:bg-blue-700 text-white rounded-lg transition-colors font-medium"
      >
        <i class="fas fa-home mr-2"></i>
        Back to Dashboard
      </a>
    </div>
  </div>
</div>
