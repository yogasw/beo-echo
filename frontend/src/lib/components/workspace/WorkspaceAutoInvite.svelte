<script lang="ts">
    import { onMount } from 'svelte';
    import { autoInviteApi } from '$lib/api/autoInviteApi';
    import type { AutoInviteConfig } from '$lib/types/autoInviteTypes';
	import { toast } from '$lib/stores/toast';
    import ToggleSwitch from '$lib/components/common/ToggleSwitch.svelte';
    import DomainManager from '$lib/components/common/DomainManager.svelte';
    
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
    
    // Domain management is now handled by the DomainManager component
    
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
                aria-label="Close auto-invite configuration"
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
                        <ToggleSwitch bind:checked={enabled} />
                    </div>
                    
                    <!-- Role selection -->
                    <div class="mb-4">
                        <label for="role-group" class="block text-gray-800 dark:text-white font-medium mb-2">
                            Role for Auto-Invited Users
                        </label>
                        <div class="flex gap-4" id="role-group" role="radiogroup">
                            <label class="inline-flex items-center">
                                <input type="radio" bind:group={role} value="member" id="role-member" class="form-radio h-4 w-4 
                                    text-blue-600 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 
                                    focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-600">
                                <span class="ml-2 text-gray-700 dark:text-gray-300">Member</span>
                            </label>
                            <label class="inline-flex items-center">
                                <input type="radio" bind:group={role} value="admin" id="role-admin" class="form-radio h-4 w-4 
                                    text-blue-600 bg-gray-100 dark:bg-gray-700 border-gray-300 dark:border-gray-600 
                                    focus:ring-2 focus:ring-blue-500 dark:focus:ring-blue-600">
                                <span class="ml-2 text-gray-700 dark:text-gray-300">Admin</span>
                            </label>
                        </div>
                    </div>
                    
                    <!-- Domain list -->
                    <div>
                        <DomainManager 
                            bind:domains={domains}
                            label="Email Domains"
                            helpText="Users with these email domains will be automatically added to this workspace"
                            emptyText="No domains added yet."
                        />
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
