<script lang="ts">
    import HeadersTab from '$lib/components/common/HeadersEditor.svelte';
    import { updateResponse } from '$lib/stores/saveButton';
    import type { Endpoint, Response } from '$lib/api/BeoApi';
    
    // Assume these are passed in as props or obtained elsewhere
    export let currentEndpoint: Endpoint;
    export let currentResponse: Response;
    
    // Define the onSave handler that will be passed to HeadersTab
    function handleHeadersSave(headersJson: string) {
        // Call updateResponse from saveButton.ts with the headers value
        updateResponse("headers", headersJson, currentEndpoint, currentResponse);
        
        // Additional actions can be performed here if needed
        console.log("Headers saved:", headersJson);
    }
</script>

<!-- Use the HeadersTab component with the onSave prop -->
<div class="headers-container">
    <HeadersTab 
        headers={currentResponse?.headers || "{}"}
        editable={true}
        onSave={handleHeadersSave}
    />
</div>

<style>
    .headers-container {
        margin-top: 1rem;
    }
</style>
