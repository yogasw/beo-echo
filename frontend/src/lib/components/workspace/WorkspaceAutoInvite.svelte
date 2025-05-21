<script lang="ts">
    import { onMount } from 'svelte';
    import { autoInviteApi } from '$lib/api/autoInviteApi';
    import type { AutoInviteConfig } from '$lib/types/autoInviteTypes';
	import { toast } from '$lib/stores/toast';
    
    export let workspaceId: string;
    export let isOwner: boolean = false;
    export let onClose: () => void;
    
    let loading = false;
    let saving = false;
    let error: string | null = null;
    let autoInviteConfig: AutoInviteConfig | null = null;
    
    // Form state
    let enabled = false;
    let domains: string[] = [];
    let newDomain = '';
    let role: 'member' | 'admin' = 'member';
    
    onMount(async () => {
        if (!isOwner) return;
        await loadConfig();
    });
    
    async function loadConfig() {
        loading = true;
        error = null;
        
        try {
            autoInviteConfig = await autoInviteApi.getConfig(workspaceId);
            
            // Set form state from config
            enabled = autoInviteConfig.enabled;
            domains = [...autoInviteConfig.domains];
            role = autoInviteConfig.role as 'member' | 'admin';
        } catch (err) {
            console.error('Failed to load auto-invite config:', err);
            error = 'Failed to load auto-invite configuration.';
        } finally {
            loading = false;
        }
    }
    
    function addDomain() {
        if (!newDomain || newDomain.trim() === '') return;
        
        const domain = newDomain.trim();
        
        // Basic validation
        if (!domain.includes('.')) {
            error = 'Please enter a valid domain (e.g. example.com)';
            return;
        }
        
        if (domain.includes('@')) {
            error = 'Do not include @ in domain. Use example.com, not @example.com';
            return;
        }
        
        if (domains.includes(domain)) {
            error = 'This domain is already in the list.';
            return;
        }
        
        domains = [...domains, domain];
        newDomain = '';
        error = null;
    }
    
    function removeDomain(domain: string) {
        domains = domains.filter(d => d !== domain);
    }
    
    async function saveConfig() {
        if (saving) return;
        error = null;
        saving = true;
        
        try {
            const updatedConfig = await autoInviteApi.updateConfig(workspaceId, {
                enabled,
                domains,
                role
            });
            
            autoInviteConfig = updatedConfig;
            toast.info('Auto-invite configuration saved successfully');
            onClose();
        } catch (err) {
            console.error('Failed to save auto-invite config:', err);
            error = 'Failed to save configuration. Please try again.';
            toast.error('Failed to save auto-invite configuration.');
        } finally {
            saving = false;
        }
    }
</script>

<div class="fixed inset-0 z-50 flex items-center justify-center bg-gray-900 bg-opacity-75 overflow-y-auto">
    <div class="relative bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-3xl mx-4 max-h-[90vh] overflow-y-auto">
        <!-- Header -->
        <div class="p-4 bg-gray-50 dark:bg-gray-750 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-800 dark:text-white">
                <i class="fas fa-user-plus mr-2"></i> Auto-Invite Configuration
            </h3>
            <button 
                on:click={onClose}
                class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-white"
            >
                <i class="fas fa-times"></i>
            </button>
        </div>

        <!-- Content -->
        <div class="p-6">
            {#if !isOwner}
                <div class="p-4 bg-gray-100 dark:bg-gray-800 rounded-md text-center text-gray-600 dark:text-gray-300 text-sm">
                    Auto-invite configuration can only be accessed by system owners.
                </div>
            {:else if loading}
                <div class="flex justify-center py-6">
                    <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"></div>
                </div>
            {:else}
                {#if error}
                    <div class="mb-4 bg-red-100 dark:bg-red-900/30 border border-red-300 dark:border-red-800 text-red-700 dark:text-red-400 px-4 py-3 rounded">
                        <i class="fas fa-exclamation-circle mr-2"></i> {error}
                    </div>
                {/if}
                
                <form on:submit|preventDefault={saveConfig} class="space-y-4">
                    <!-- Enable/Disable toggle -->
                    <div class="flex items-center justify-between">
                        <div>
                            <span class="text-gray-800 dark:text-white font-medium">Enable Auto-Invite</span>
                            <p class="text-sm text-gray-500 dark:text-gray-400">
                                Automatically invite new users based on email domains
                            </p>
                        </div>
                        <label class="inline-flex items-center cursor-pointer">
                            <input type="checkbox" bind:checked={enabled} class="sr-only peer">
                            <div class="w-11 h-6 bg-gray-300 dark:bg-gray-700 peer-checked:bg-blue-600 
                                rounded-full peer peer-focus:ring-2 peer-focus:ring-blue-300 dark:peer-focus:ring-blue-800 
                                peer-checked:after:translate-x-full peer-checked:after:border-white 
                                after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white 
                                after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 
                                after:transition-all dark:border-gray-600">
                            </div>
                        </label>
                    </div>
                    
                    <!-- Role selection -->
                    <div class="mb-4">
                        <label class="block text-gray-800 dark:text-white font-medium mb-2">
                            Role for Auto-Invited Users
                        </label>
                        <div class="flex gap-4">
                            <label class="inline-flex items-center">
                                <input type="radio" bind:group={role} value="member" class="form-radio h-4 w-4 
                                    text-blue-600 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 
                                    focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-600">
                                <span class="ml-2 text-gray-700 dark:text-gray-300">Member</span>
                            </label>
                            <label class="inline-flex items-center">
                                <input type="radio" bind:group={role} value="admin" class="form-radio h-4 w-4 
                                    text-blue-600 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 
                                    focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-600">
                                <span class="ml-2 text-gray-700 dark:text-gray-300">Admin</span>
                            </label>
                        </div>
                    </div>
                    
                    <!-- Domain list -->
                    <div>
                        <label class="block text-gray-800 dark:text-white font-medium mb-2">
                            Email Domains
                        </label>
                        <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">
                            Users with these email domains will be automatically added to this workspace
                        </p>
                        
                        <!-- Add domain input -->
                        <div class="flex mb-3">
                            <input 
                                type="text" 
                                bind:value={newDomain} 
                                placeholder="example.com" 
                                class="flex-1 rounded-l-lg p-2.5 bg-gray-50 dark:bg-gray-700 border 
                                    border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white
                                    focus:ring-blue-500 focus:border-blue-500"
                            />
                            <button 
                                type="button" 
                                on:click={addDomain} 
                                class="rounded-r-lg px-4 py-2.5 bg-blue-600 hover:bg-blue-700 text-white"
                            >
                                Add
                            </button>
                        </div>
                        
                        <!-- Domains list -->
                        {#if domains.length === 0}
                            <p class="text-gray-500 dark:text-gray-400 text-sm italic">
                                No domains added yet.
                            </p>
                        {:else}
                            <ul class="space-y-2 max-h-60 overflow-y-auto">
                                {#each domains as domain}
                                    <li class="flex justify-between items-center py-2 px-3 rounded-md
                                        bg-gray-50 dark:bg-gray-750 border border-gray-200 dark:border-gray-700">
                                        <span class="text-gray-800 dark:text-white">{domain}</span>
                                        <button 
                                            type="button"
                                            on:click={() => removeDomain(domain)}
                                            class="text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-500"
                                        >
                                            <i class="fas fa-times"></i>
                                        </button>
                                    </li>
                                {/each}
                            </ul>
                        {/if}
                    </div>
                    
                    <div class="mt-3 p-3 bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded text-blue-700 dark:text-blue-300 text-sm">
                        <i class="fas fa-info-circle mr-2"></i>
                        <strong>Note:</strong> This feature will automatically add users with matching email domains to this workspace when they sign up or log in.
                    </div>
                    
                    <!-- Buttons -->
                    <div class="flex justify-end pt-4 border-t border-gray-200 dark:border-gray-700">
                        <button
                            type="button"
                            on:click={onClose}
                            class="mr-2 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-800 dark:text-gray-300 rounded-md hover:bg-gray-300 dark:hover:bg-gray-600"
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={saving}
                            class="bg-blue-600 hover:bg-blue-700 text-white py-2 px-4 rounded-md
                                flex items-center disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                            {#if saving}
                                <div class="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                                Saving...
                            {:else}
                                <i class="fas fa-save mr-2"></i> Save Configuration
                            {/if}
                        </button>
                    </div>
                </form>
            {/if}
        </div>
    </div>
</div>
