<script lang="ts">
  import { fade, scale } from 'svelte/transition';

  // Core modal props
  export let isOpen: boolean;
  export let title: string;
  export let subtitle: string = '';
  export let maxWidth: string = 'max-w-4xl';
  export let onClose: () => void;

  // Step management props  
  export let currentStep: number;
  export let totalSteps: number;
  export let stepLabels: string[] = [];

  // Navigation props
  export let onNext: (() => void) | null = null;
  export let onPrevious: (() => void) | null = null;
  export let onSubmit: (() => void) | null = null;
  export let isSubmitting: boolean = false;
  export let canGoNext: boolean = true;
  export let canGoPrevious: boolean = true;
  export let canSubmit: boolean = true;

  // Button customization
  export let nextButtonText: string = 'Next';
  export let previousButtonText: string = 'Previous';
  export let submitButtonText: string = 'Create';
  export let showNavigationButtons: boolean = true;

  // Exit confirmation
  export let hasUnsavedChanges: (() => boolean) | null = null;
  export let showExitConfirmation: boolean = false;

  // Error handling
  export let error: string | null = null;

  // Handle backdrop click
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      handleClose();
    }
  }

  // Handle escape key
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      handleClose();
    }
  }

  // Handle close with unsaved changes check
  function handleClose() {
    if (hasUnsavedChanges && hasUnsavedChanges()) {
      showExitConfirmation = true;
    } else {
      onClose();
    }
  }

  // Direct close (X button)
  function handleDirectClose() {
    handleClose();
  }

  // Navigation handlers
  function handleNext() {
    if (onNext && canGoNext) {
      onNext();
    }
  }

  function handlePrevious() {
    if (onPrevious && canGoPrevious) {
      onPrevious();
    }
  }

  function handleSubmit() {
    if (onSubmit && canSubmit && !isSubmitting) {
      onSubmit();
    }
  }

  // Generate default step labels if not provided
  $: defaultStepLabels = stepLabels.length === totalSteps 
    ? stepLabels 
    : Array.from({ length: totalSteps }, (_, i) => `Step ${i + 1}`);
</script>

{#if isOpen}
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-50 p-4"
    transition:fade={{ duration: 200 }}
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <div 
      class="bg-white dark:bg-gray-800 rounded-xl w-full {maxWidth} max-h-[90vh] shadow-2xl border border-gray-200 dark:border-gray-700 flex flex-col"
      transition:scale={{ duration: 200, start: 0.9 }}
      role="dialog"
      aria-modal="true"
    >
      <!-- Header with Progress -->
      <div class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-800 dark:to-gray-750 p-6 border-b border-gray-200 dark:border-gray-700 shrink-0">
        <div class="flex justify-between items-center mb-4">
          <div>
            <h2 class="text-2xl font-bold text-gray-800 dark:text-white mb-1">{title}</h2>
            {#if subtitle}
              <p class="text-gray-600 dark:text-gray-400 text-sm">{subtitle}</p>
            {:else}
              <p class="text-gray-600 dark:text-gray-400 text-sm">Step {currentStep} of {totalSteps}</p>
            {/if}
          </div>
          <button
            class="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
            on:click={handleDirectClose}
            aria-label="Close"
          >
            <i class="fas fa-times text-lg"></i>
          </button>
        </div>
        
        <!-- Progress Bar -->
        <div class="flex space-x-2">
          {#each Array.from({ length: totalSteps }, (_, i) => i + 1) as step}
            <div class="flex-1 h-2 rounded-full bg-gray-300 dark:bg-gray-600 overflow-hidden relative">
              <div 
                class="h-full bg-gradient-to-r from-blue-500 to-blue-600 transition-all duration-500 ease-out rounded-full"
                class:w-full={currentStep >= step}
                class:w-0={currentStep < step}
              ></div>
              {#if currentStep === step}
                <div class="absolute inset-0 bg-gradient-to-r from-blue-400 to-blue-500 opacity-30 animate-pulse rounded-full"></div>
              {/if}
            </div>
          {/each}
        </div>
        
        <!-- Step Labels -->
        <div class="flex justify-between mt-2 text-xs">
          {#each defaultStepLabels as label, index}
            <span class="font-medium transition-colors"
                  class:text-blue-600={currentStep >= index + 1}
                  class:dark:text-blue-400={currentStep >= index + 1}
                  class:text-gray-400={currentStep < index + 1}
                  class:dark:text-gray-500={currentStep < index + 1}>
              {label}
            </span>
          {/each}
        </div>
      </div>

      <!-- Body -->
      <div class="flex-1 overflow-y-auto min-h-0">
        {#if error}
          <div class="mx-6 mt-4 p-4 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <div class="flex items-center">
              <i class="fas fa-exclamation-triangle text-red-500 dark:text-red-400 mr-3"></i>
              <div>
                <p class="text-red-700 dark:text-red-300 font-medium">Error</p>
                <p class="text-red-600 dark:text-red-200 text-sm">{error}</p>
              </div>
            </div>
          </div>
        {/if}

        <!-- Content Slot -->
        <slot />
      </div>

      <!-- Footer with Navigation -->
      {#if showNavigationButtons}
        <div class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-750 dark:to-gray-800 px-6 py-4 border-t border-gray-200 dark:border-gray-700 shrink-0">
          <div class="flex justify-between items-center">
            <!-- Previous Button -->
            <div>
              {#if currentStep > 1 && onPrevious}
                <button
                  type="button"
                  class="bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white py-2 px-4 rounded-md text-sm flex items-center transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={!canGoPrevious}
                  on:click={handlePrevious}
                >
                  <i class="fas fa-chevron-left mr-2"></i>
                  {previousButtonText}
                </button>
              {/if}
            </div>

            <!-- Next/Submit Button -->
            <div>
              {#if currentStep < totalSteps && onNext}
                <button
                  type="button"
                  class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md text-sm flex items-center transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={!canGoNext}
                  on:click={handleNext}
                >
                  {nextButtonText}
                  <i class="fas fa-chevron-right ml-2"></i>
                </button>
              {:else if currentStep === totalSteps && onSubmit}
                <button
                  type="button"
                  class="bg-green-600 hover:bg-green-700 text-white py-2 px-4 rounded-md text-sm flex items-center transition-colors focus:outline-none focus:ring-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  disabled={!canSubmit || isSubmitting}
                  on:click={handleSubmit}
                >
                  {#if isSubmitting}
                    <i class="fas fa-spinner fa-spin mr-2"></i>
                    Creating...
                  {:else}
                    <i class="fas fa-check mr-2"></i>
                    {submitButtonText}
                  {/if}
                </button>
              {/if}
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
{/if}

<!-- Exit Confirmation Dialog -->
{#if showExitConfirmation}
  <div class="fixed inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center z-[60] p-4" transition:fade={{ duration: 200 }}>
    <div class="bg-white dark:bg-gray-800 rounded-xl p-6 shadow-2xl border border-gray-200 dark:border-gray-700 max-w-md w-full" transition:scale={{ duration: 200, start: 0.9 }}>
      <div class="flex items-center mb-4">
        <i class="fas fa-exclamation-triangle text-yellow-500 dark:text-yellow-400 mr-3 text-xl"></i>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-white">Unsaved Changes</h3>
      </div>
      <p class="text-gray-600 dark:text-gray-400 mb-6">
        You have unsaved changes. Are you sure you want to close without saving?
      </p>
      <div class="flex justify-end space-x-3">
        <button
          type="button"
          class="bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white py-2 px-4 rounded-md text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500"
          on:click={() => showExitConfirmation = false}
        >
          Cancel
        </button>
        <button
          type="button"
          class="bg-red-600 hover:bg-red-700 text-white py-2 px-4 rounded-md text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-red-500"
          on:click={() => {
            showExitConfirmation = false;
            onClose();
          }}
        >
          Close Without Saving
        </button>
      </div>
    </div>
  </div>
{/if}
