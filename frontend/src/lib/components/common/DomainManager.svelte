<script lang="ts">
    /**
     * A reusable component for managing a list of domains.
     * Supports adding multiple domains at once with comma separation,
     * validation, and removal.
     */
    
    // The list of domains to manage - will be updated via two-way binding
    export let domains: string[] = [];
    
    // Optional properties
    export let placeholder = "example.com, gmail.com, company.net";
    export let label = "Email Domains";
    export let helpText = "Users with these email domains will be automatically added";
    export let emptyText = "No domains added yet.";
    export let maxHeight = "60"; // Max height in rem for the domains list container
    
    // Input state
    let newDomain = '';
    let error: string | null = null;
    
    function addDomain() {
        if (!newDomain || newDomain.trim() === '') return;
        
        // Split the input by commas and process each domain
        const inputDomains = newDomain.split(',').map(d => d.trim()).filter(d => d !== '');
        let hasErrors = false;
        let newDomainsAdded = false;
        
        for (const domain of inputDomains) {
            // Basic validation
            if (!domain.includes('.')) {
                error = `Invalid domain: "${domain}" - missing dot (e.g. example.com)`;
                hasErrors = true;
                continue;
            }
            
            if (domain.includes('@')) {
                error = `Invalid domain: "${domain}" - contains @ symbol. Use example.com, not @example.com`;
                hasErrors = true;
                continue;
            }
            
            if (domains.includes(domain)) {
                continue; // Skip duplicates silently
            }
            
            domains = [...domains, domain];
            newDomainsAdded = true;
        }
        
        if (!hasErrors && newDomainsAdded) {
            error = null;
        } else if (!hasErrors && !newDomainsAdded) {
            error = 'All domains are already in the list.';
        }
        
        newDomain = ''; // Clear the input field after processing
    }
    
    function removeDomain(domain: string) {
        domains = domains.filter(d => d !== domain);
    }
    
    // Handle keydown events to allow Enter key for adding domains
    function handleKeydown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            event.preventDefault();
            addDomain();
        }
    }
</script>

<div>
    <!-- Component label -->
    {#if label}
        <label for="domainInput" class="block text-gray-800 dark:text-white font-medium mb-2">
            {label}
        </label>
    {/if}
    
    <!-- Help text -->
    {#if helpText}
        <p class="text-sm text-gray-500 dark:text-gray-400 mb-3">
            {helpText}
        </p>
    {/if}
    
    <!-- Error message -->
    {#if error}
        <div class="mb-3 text-sm text-red-600 dark:text-red-400">
            <i class="fas fa-exclamation-circle mr-1"></i> {error}
        </div>
    {/if}
    
    <!-- Add domain input -->
    <div class="flex mb-3">
        <input 
            id="domainInput"
            type="text" 
            bind:value={newDomain}
            on:keydown={handleKeydown}
            {placeholder}
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

    <!-- Helper text for comma separation -->
    <p class="text-xs text-gray-500 dark:text-gray-400 mb-3">
        <i class="fas fa-info-circle mr-1"></i>
        You can enter multiple domains separated by commas
    </p>
    
    <!-- Domains list -->
    {#if domains.length === 0}
        <p class="text-gray-500 dark:text-gray-400 text-sm italic">
            {emptyText}
        </p>
    {:else}
        <ul class="space-y-2 max-h-{maxHeight} overflow-y-auto">
            {#each domains as domain}
                <li class="flex justify-between items-center py-2 px-3 rounded-md
                    bg-gray-50 dark:bg-gray-750 border border-gray-200 dark:border-gray-700">
                    <span class="text-gray-800 dark:text-white">{domain}</span>
                    <button 
                        type="button"
                        on:click={() => removeDomain(domain)}
                        class="text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-500"
                        aria-label={`Remove domain ${domain}`}
                    >
                        <i class="fas fa-times"></i>
                    </button>
                </li>
            {/each}
        </ul>
    {/if}
</div>
