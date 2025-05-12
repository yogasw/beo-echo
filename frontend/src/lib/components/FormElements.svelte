<script lang="ts">
  import * as ThemeUtils from '$lib/utils/themeUtils';
  import { theme } from '$lib/stores/theme';
  
  // Input props
  export let label: string = '';
  export let id: string = '';
  export let name: string = '';
  export let value: string = '';
  export let type: string = 'text';
  export let placeholder: string = '';
  export let required: boolean = false;
  export let disabled: boolean = false;
  export let errorMessage: string = '';
  export let icon: string = '';
  export let inputClass: string = '';
  export let containerClass: string = '';
  export let labelClass: string = '';
  export let variant: 'default' | 'compact' = 'default';
</script>

<!-- Reusable Input with Icon -->
{#if type !== 'select' && type !== 'textarea' && type !== 'toggle'}
  <div class="mb-4 {containerClass}">
    {#if label}
      <label for={id} class={`block mb-2 text-sm font-medium theme-text-primary ${labelClass}`}>
        {label}
        {#if required}<span class="text-red-500 ml-1">*</span>{/if}
      </label>
    {/if}
    
    <div class="relative">
      {#if icon}
        <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
          <i class={`fas fa-${icon} theme-text-secondary`}></i>
        </div>
      {/if}
      
      <input 
        {id} 
        {name} 
        {type} 
        {placeholder} 
        {disabled} 
        {required} 
        bind:value 
        class={ThemeUtils.inputField(inputClass)}
        aria-invalid={errorMessage ? "true" : "false"}
      />
    </div>
    
    {#if errorMessage}
      <p class="mt-2 text-xs text-red-500">{errorMessage}</p>
    {/if}
  </div>
{:else if type === 'textarea'}
  <div class="mb-4 {containerClass}">
    {#if label}
      <label for={id} class={`block mb-2 text-sm font-medium theme-text-primary ${labelClass}`}>
        {label}
        {#if required}<span class="text-red-500 ml-1">*</span>{/if}
      </label>
    {/if}
    
    <div class="relative">
      {#if icon}
        <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
          <i class={`fas fa-${icon} theme-text-secondary`}></i>
        </div>
      {/if}
      
      <textarea 
        {id} 
        {name} 
        {placeholder} 
        {disabled} 
        {required} 
        bind:value 
        class={ThemeUtils.inputField(`resize-none ${inputClass}`)}
        aria-invalid={errorMessage ? "true" : "false"}
      ></textarea>
    </div>
    
    {#if errorMessage}
      <p class="mt-2 text-xs text-red-500">{errorMessage}</p>
    {/if}
  </div>
{:else if type === 'select'}
  <div class="mb-4 {containerClass}">
    {#if label}
      <label for={id} class={`block mb-2 text-sm font-medium theme-text-primary ${labelClass}`}>
        {label}
        {#if required}<span class="text-red-500 ml-1">*</span>{/if}
      </label>
    {/if}
    
    <div class="relative">
      {#if icon}
        <div class="absolute inset-y-0 start-0 flex items-center ps-3 pointer-events-none">
          <i class={`fas fa-${icon} theme-text-secondary`}></i>
        </div>
      {/if}
      
      <select 
        {id} 
        {name} 
        {disabled} 
        {required} 
        bind:value 
        class={ThemeUtils.inputField(`${icon ? 'ps-10' : ''} ${inputClass}`)}
        aria-invalid={errorMessage ? "true" : "false"}
      >
        <slot></slot>
      </select>
      
      <div class="absolute inset-y-0 end-3 flex items-center pointer-events-none">
        <i class="fas fa-chevron-down theme-text-secondary text-xs"></i>
      </div>
    </div>
    
    {#if errorMessage}
      <p class="mt-2 text-xs text-red-500">{errorMessage}</p>
    {/if}
  </div>
{:else if type === 'toggle'}
  <div class={`flex items-center justify-between mb-4 ${containerClass}`}>
    {#if label}
      <span class={`text-sm font-medium theme-text-primary ${labelClass}`}>
        {label}
        {#if required}<span class="text-red-500 ml-1">*</span>{/if}
      </span>
    {/if}
    
    <label class="inline-flex items-center cursor-pointer">
      <input 
        type="checkbox" 
        {id}
        {name}
        {disabled}
        bind:checked={value === 'true' || value === true}
        on:change={(e) => value = e.target.checked}
        class="sr-only peer"
      />
      <div class="w-11 h-6 bg-gray-700 peer-focus:outline-none peer-focus:ring-4 
        peer-focus:ring-blue-800 rounded-full peer peer-checked:after:translate-x-full 
        rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white 
        after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
        after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
        after:transition-all peer-checked:bg-blue-600"></div>
    </label>
  </div>
{/if}
