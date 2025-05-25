<!-- Reusable Step Modal Component -->
<script lang="ts">
  import { fade, scale } from 'svelte/transition';
  import { createEventDispatcher } from 'svelte';

  const dispatch = createEventDispatcher();

  export let isOpen: boolean = false;
  export let title: string = '';
  export let subtitle: string = '';
  export let currentStep: number = 1;
  export let totalSteps: number = 3;
  export let stepLabels: string[] = [];
  export let error: string | null = null;
  export let validationErrors: Record<string, string> = {};
  export let isSubmitting: boolean = false;
  export let canGoNext: boolean = true;
  export let canGoPrev: boolean = true;
  export let nextButtonText: string = 'Continue';
  export let submitButtonText: string = 'Submit';
  export let hasUnsavedChanges: boolean = false;
  export let showExitConfirmation: boolean = false;
  export let maxWidth: string = 'max-w-4xl';

  // Close handlers
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      dispatch('requestClose', { hasUnsavedChanges });
    }
  }

  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      dispatch('requestClose', { hasUnsavedChanges });
    }
  }

  function handleDirectClose() {
    dispatch('requestClose', { hasUnsavedChanges });
  }

  // Step navigation
  function nextStep() {
    dispatch('nextStep');
  }

  function prevStep() {
    dispatch('prevStep');
  }

  function handleSubmit() {
    dispatch('submit');
  }

  // Exit confirmation
  function confirmExit() {
    dispatch('confirmExit');
  }

  function cancelExit() {
    dispatch('cancelExit');
  }
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
            <p class="text-gray-600 dark:text-gray-400 text-sm">{subtitle || `Step ${currentStep} of ${totalSteps}`}</p>
          </div>
          <button
            class="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
            on:click={handleDirectClose}
            aria-label="Close"
            title="Close"
          >
            <i class="fas fa-times text-lg"></i>
          </button>
        </div>
        
        <!-- Progress Bar -->
        <div class="flex space-x-2">
          {#each Array(totalSteps) as _, index}
            {@const step = index + 1}
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
        {#if stepLabels.length === totalSteps}
          <div class="flex justify-between mt-2 text-xs">
            {#each stepLabels as label, index}
              {@const step = index + 1}
              <span class="font-medium transition-colors"
                    class:text-blue-600={currentStep >= step}
                    class:dark:text-blue-400={currentStep >= step}
                    class:text-gray-400={currentStep < step}
                    class:dark:text-gray-500={currentStep < step}>
                {label}
              </span>
            {/each}
          </div>
        {/if}
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

        <!-- Step Content Slot -->
        <slot name="step-content" {currentStep} />
      </div>

      <!-- Footer -->
      <div class="bg-gradient-to-r from-gray-50 to-gray-100 dark:from-gray-750 dark:to-gray-800 px-6 py-4 border-t border-gray-200 dark:border-gray-700 shrink-0">
        <div class="flex justify-between items-center">
          <!-- Left: Step Progress Indicator -->
          <div class="flex items-center space-x-2">
            <div class="flex space-x-1">
              {#each Array(totalSteps) as _, index}
                <div class="w-2 h-2 rounded-full transition-colors duration-200 
                  {index < currentStep ? 'bg-blue-600' : 
                   index === currentStep - 1 ? 'bg-blue-500' : 'bg-gray-300 dark:bg-gray-600'}">
                </div>
              {/each}
            </div>
            <span class="text-xs text-gray-500 dark:text-gray-400 ml-2">
              Step {currentStep} of {totalSteps}
            </span>
          </div>
          
          <!-- Right: Navigation & Primary Actions -->
          <div class="flex items-center space-x-3">
            <!-- Step Navigation -->
            {#if currentStep > 1}
              <button
                type="button"
                class="px-4 py-2 bg-white dark:bg-gray-700 hover:bg-gray-50 dark:hover:bg-gray-600 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600 rounded-lg text-sm flex items-center transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-gray-500 shadow-sm"
                on:click={prevStep}
                disabled={isSubmitting || !canGoPrev}
                title="Go back to previous step"
              >
                <i class="fas fa-chevron-left mr-2"></i>Back
              </button>
            {/if}
            
            <!-- Primary Action -->
            {#if currentStep < totalSteps}
              <button
                type="button"
                class="px-6 py-2 bg-blue-600 hover:bg-blue-700 disabled:bg-blue-400 text-white rounded-lg text-sm flex items-center transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500 shadow-md"
                on:click={nextStep}
                disabled={isSubmitting || !canGoNext}
                title="Continue to next step"
              >
                {nextButtonText}<i class="fas fa-chevron-right ml-2"></i>
              </button>
            {:else}
              <button
                type="button"
                class="px-6 py-2 bg-green-600 hover:bg-green-700 disabled:bg-green-400 text-white rounded-lg text-sm flex items-center transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-green-500 shadow-md"
                on:click={handleSubmit}
                disabled={isSubmitting}
                title="Submit the form"
              >
                {#if isSubmitting}
                  <div class="animate-spin h-4 w-4 border-2 border-white border-t-transparent rounded-full mr-2"></div>
                  Submitting...
                {:else}
                  <i class="fas fa-check mr-2"></i>{submitButtonText}
                {/if}
              </button>
            {/if}
          </div>
        </div>
        
        <!-- Validation Summary -->
        {#if Object.keys(validationErrors).length > 0}
          <div class="mt-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
            <div class="flex items-start">
              <i class="fas fa-exclamation-triangle text-red-500 dark:text-red-400 mr-2 mt-0.5"></i>
              <div>
                <p class="text-red-700 dark:text-red-300 font-medium text-sm">Please fix the following errors:</p>
                <ul class="text-red-600 dark:text-red-200 text-xs mt-1 space-y-1">
                  {#each Object.values(validationErrors) as error}
                    <li>• {error}</li>
                  {/each}
                </ul>
              </div>
            </div>
          </div>
        {/if}
      </div>
    </div>
  </div>
{/if}

<!-- Exit Confirmation Modal -->
{#if showExitConfirmation}
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div
    class="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-[60] p-4"
    transition:fade={{ duration: 150 }}
    on:click|self={cancelExit}
    on:keydown={(e) => e.key === 'Escape' && cancelExit()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <div 
      class="bg-white dark:bg-gray-800 rounded-lg w-full max-w-md shadow-xl border border-gray-200 dark:border-gray-700"
      transition:scale={{ duration: 150, start: 0.95 }}
    >
      <!-- Header -->
      <div class="p-6 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center">
          <div class="w-10 h-10 bg-yellow-100 dark:bg-yellow-900/30 rounded-full flex items-center justify-center mr-4">
            <i class="fas fa-exclamation-triangle text-yellow-600 dark:text-yellow-400"></i>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-800 dark:text-white">Exit without saving?</h3>
            <p class="text-sm text-gray-600 dark:text-gray-400">Your changes will be lost</p>
          </div>
        </div>
      </div>
      
      <!-- Body -->
      <div class="p-6">
        <p class="text-gray-700 dark:text-gray-300 text-sm leading-relaxed">
          You have unsaved changes. Are you sure you want to exit? 
          All changes will be lost and cannot be recovered.
        </p>
      </div>
      
      <!-- Footer -->
      <div class="p-6 pt-0 flex justify-end space-x-3">
        <button
          type="button"
          class="px-4 py-2 bg-gray-200 dark:bg-gray-600 hover:bg-gray-300 dark:hover:bg-gray-500 text-gray-800 dark:text-white rounded-lg text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-gray-500"
          on:click={cancelExit}
        >
          <i class="fas fa-arrow-left mr-2"></i>Continue Editing
        </button>
        <button
          type="button"
          class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg text-sm transition-colors focus:outline-none focus:ring-2 focus:ring-red-500"
          on:click={confirmExit}
        >
          <i class="fas fa-sign-out-alt mr-2"></i>Yes, Exit
        </button>
      </div>
    </div>
  </div>
{/if}
