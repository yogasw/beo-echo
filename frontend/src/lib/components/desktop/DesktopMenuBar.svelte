<script lang="ts">
  import { browser } from '$app/environment';
  import { onMount } from 'svelte';
  import { isDesktopMode } from '$lib/utils/desktopConfig';
  import wailsBridge from '$lib/utils/wails-bridge';
  
  let appVersion = '1.0.0';
  let platform = 'web';
  
  onMount(() => {
    if (isDesktopMode()) {
      appVersion = wailsBridge.getAppVersion();
      platform = wailsBridge.getPlatform();
    }
  });
  
  // Desktop menu actions
  const handleAbout = () => {
    if (isDesktopMode()) {
      wailsBridge.showDialog(
        'About BeoEcho Desktop',
        `BeoEcho Desktop v${appVersion}\nAPI Mocking Service\nPlatform: ${platform}`,
        'info'
      );
    }
  };
  
  const handleQuit = () => {
    if (isDesktopMode()) {
      wailsBridge.quit();
    }
  };
  
  const handleOpenDocs = () => {
    const docsUrl = 'https://github.com/your-org/beo-echo/docs';
    if (isDesktopMode()) {
      wailsBridge.openExternal(docsUrl);
    } else {
      window.open(docsUrl, '_blank');
    }
  };
  
  const handleOpenGitHub = () => {
    const githubUrl = 'https://github.com/your-org/beo-echo';
    if (isDesktopMode()) {
      wailsBridge.openExternal(githubUrl);
    } else {
      window.open(githubUrl, '_blank');
    }
  };
</script>

<!-- Desktop Menu Bar - Only shown in desktop mode -->
{#if browser && isDesktopMode()}
  <div class="bg-gray-900 border-b border-gray-700 px-4 py-2 flex items-center justify-between text-sm">
    <!-- App Title and Version -->
    <div class="flex items-center space-x-4">
      <div class="flex items-center space-x-2">
        <i class="fas fa-cube text-blue-400"></i>
        <span class="font-semibold text-white">BeoEcho Desktop</span>
        <span class="text-gray-400 text-xs">v{appVersion}</span>
      </div>
      
      <!-- Platform indicator -->
      <div class="flex items-center space-x-1">
        <span class="w-2 h-2 rounded-full bg-green-400"></span>
        <span class="text-gray-400 text-xs capitalize">{platform}</span>
      </div>
    </div>
    
    <!-- Desktop Menu Actions -->
    <div class="flex items-center space-x-2">
      <!-- Help Menu -->
      <button
        class="px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-800 rounded text-xs transition-colors"
        title="View documentation"
        aria-label="Open documentation"
        on:click={handleOpenDocs}
      >
        <i class="fas fa-book mr-1"></i>
        Docs
      </button>
      
      <!-- GitHub Link -->
      <button
        class="px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-800 rounded text-xs transition-colors"
        title="View source code on GitHub"
        aria-label="Open GitHub repository"
        on:click={handleOpenGitHub}
      >
        <i class="fab fa-github mr-1"></i>
        GitHub
      </button>
      
      <!-- About -->
      <button
        class="px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-800 rounded text-xs transition-colors"
        title="About BeoEcho Desktop"
        aria-label="Show about dialog"
        on:click={handleAbout}
      >
        <i class="fas fa-info-circle mr-1"></i>
        About
      </button>
      
      <!-- Quit (Desktop only) -->
      <button
        class="px-2 py-1 text-red-300 hover:text-red-200 hover:bg-red-900/20 rounded text-xs transition-colors ml-2"
        title="Quit application"
        aria-label="Quit BeoEcho Desktop"
        on:click={handleQuit}
      >
        <i class="fas fa-times mr-1"></i>
        Quit
      </button>
    </div>
  </div>
{/if}
