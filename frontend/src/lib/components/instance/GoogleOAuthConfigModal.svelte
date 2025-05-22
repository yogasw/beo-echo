<!-- src/lib/components/instance/GoogleOAuthConfigModal.svelte -->
<script lang="ts">
    import { ssoStore } from '$lib/stores/ssoStore';
    import { ssoApi } from '$lib/api/ssoApi';
    import { onMount } from 'svelte';
    import { fade } from 'svelte/transition';
    import DomainManager from '$lib/components/common/DomainManager.svelte';
    import ToggleSwitch from '$lib/components/common/ToggleSwitch.svelte';

    export let visible = false;
    export let onClose = () => {};

    let clientId = '';
    let clientSecret = '';
    let allowedDomains: string[] = [];
    let isEnabled = false;
    let saving = false;
    let error = '';

    onMount(async () => {
        try {
            const config = await ssoApi.getGoogleConfig();
            if (config) {
                clientId = config.config.client_id;
                clientSecret = config.config.client_secret;
                allowedDomains = config.config.allow_domains || [];
                isEnabled = config.enabled;
            }
        } catch (err) {
            error = 'Failed to load configuration';
        }
    });

    async function handleSubmit() {
        saving = true;
        error = '';

        try {
            await ssoApi.updateGoogleConfig({
                client_id: clientId,
                client_secret: clientSecret,
                allow_domains: allowedDomains,
                instructions: ''
            });

            await ssoApi.updateGoogleState(isEnabled);

            // Update store
            ssoStore.setGoogleConfig({
                client_id: clientId,
                client_secret: clientSecret,
                allow_domains: allowedDomains,
                instructions: ''
            });
            ssoStore.setGoogleEnabled(isEnabled);

            onClose();
        } catch (err) {
            error = err.message || 'Failed to save configuration';
        } finally {
            saving = false;
        }
    }
</script>

{#if visible}
<div
    class="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center"
    transition:fade={{ duration: 200 }}
    on:click|self={onClose}
    role="dialog"
    aria-modal="true"
    tabindex="0"
    on:keydown={(e) => { if (e.key === 'Escape') onClose(); }}
>
    <div
        class="absolute inset-0"
        role="button"
        tabindex="0"
        aria-label="Close modal"
        on:click|self={onClose}
        on:keydown={(e) => { if (e.key === 'Enter' || e.key === ' ') onClose(); }}
        style="z-index: -1;"
    ></div>
    <div class="theme-bg-primary rounded-lg shadow-xl max-w-2xl w-full mx-4 overflow-hidden"
         on:click|stopPropagation 
         on:keydown|stopPropagation
         role="dialog"
         tabindex="0">
        <div class="p-4 border-b theme-border">
            <h3 class="theme-text-primary font-medium">Configure Google SSO</h3>
        </div>

        <form on:submit|preventDefault={handleSubmit} class="p-4 space-y-4">
            {#if error}
                <div class="bg-red-900/20 border border-red-700 text-red-500 p-3 rounded-md text-sm">
                    {error}
                </div>
            {/if}

            <div class="space-y-2">
                <label for="clientId" class="block theme-text-primary text-sm font-medium">
                    Client ID
                </label>
                <input 
                    type="text" 
                    id="clientId"
                    bind:value={clientId}
                    class="block w-full p-2.5 theme-bg-secondary theme-border border rounded-lg 
                           text-sm theme-text-primary focus:ring-1 focus:ring-blue-500"
                    placeholder="Google OAuth Client ID"
                />
            </div>

            <div class="space-y-2">
                <label for="clientSecret" class="block theme-text-primary text-sm font-medium">
                    Client Secret
                </label>
                <input 
                    type="password" 
                    id="clientSecret"
                    bind:value={clientSecret}
                    class="block w-full p-2.5 theme-bg-secondary theme-border border rounded-lg 
                           text-sm theme-text-primary focus:ring-1 focus:ring-blue-500"
                    placeholder="Google OAuth Client Secret"
                />
            </div>

            <div class="space-y-2">
                <DomainManager 
                    bind:domains={allowedDomains}
                    label="Allowed Domains"
                    helpText="Only users with these email domains will be allowed to sign in with Google SSO. Leave empty to allow all domains."
                    emptyText="No domains restrictions - all users can sign in."
                    placeholder="example.com, yourdomain.com"
                />
            </div>

            <div class="flex items-center space-x-3 pt-2">
                <ToggleSwitch bind:checked={isEnabled} />
                <span class="theme-text-primary text-sm">Enable Google SSO</span>
            </div>

            <div class="pt-4 flex justify-end space-x-3">
                <button
                    type="button"
                    class="px-4 py-2 theme-bg-secondary hover:bg-gray-600 
                           theme-text-primary rounded-md text-sm"
                    on:click={onClose}
                    disabled={saving}>
                    Cancel
                </button>
                <button
                    type="submit"
                    class="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md text-sm 
                           flex items-center gap-2"
                    disabled={saving}>
                    {#if saving}
                        <i class="fas fa-spinner fa-spin"></i>
                    {/if}
                    Save Changes
                </button>
            </div>
        </form>
    </div>
</div>
{/if}
