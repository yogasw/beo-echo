import { writable, get } from 'svelte/store';
import { toast } from './toast';
import type { Endpoint, Project, Response } from '$lib/api/BeoApi';
import { selectedProject } from './selectedConfig';

export const showSaveButton = writable<boolean>(false);
export const saveInprogress = writable<boolean>(false);

// Endpoint update list only contains project ID and endpoint ID
export interface EndpointUpdate {
    projectId: string;
    endpointId: string;
    [key: string]: any; // For additional key-value pairs
}

// Response update interface
export interface ResponseUpdate {
    projectId: string;
    endpointId: string;
    responseId: string;
    [key: string]: any; // For additional key-value pairs
}

export const endpointsUpdateList = writable<(EndpointUpdate | ResponseUpdate)[]>([]);

/**
 * Handle save button functionality
 * This function is a more robust version of the save button handler
 * that can be called from different components
 */
export function saveButtonHandler() {
    // Get current save in progress state
    let currentValue = get(saveInprogress);
    
    // If a save is already in progress, don't start another one
    if (currentValue) {
        console.log('Save already in progress, ignoring request');
        return;
    }
    
    // Toggle the save button state to indicate a save is in progress
    saveInprogress.set(true);
    
    // Return the current list of updates for processing
    return getEndpointsUpdateList();
}

export function updateEndpoint(key: string, value: any, currentEndpoint: Endpoint) {
    let currentProject = get(selectedProject) as Project | null;
    if (!currentProject) {
        toast.error('No project selected');
        return currentEndpoint;
    }

    // Validate that key and value are provided
    if (!key) {
        toast.error('Invalid update: key is required');
        return currentEndpoint;
    }

    if (value === undefined || value === null) {
        // Allow null values for proxy_target_id
        if (key !== 'proxy_target_id') {
            toast.error(`Invalid update: value for '${key}' cannot be empty`);
            return currentEndpoint;
        }
    }

    if (!currentEndpoint) {
        toast.error('Invalid update: endpoint object is required');
        return currentEndpoint;
    }

    // Check if the value is different from what's already in the endpoint
    // Skip update if the value is the same
    if (key !== 'all' && currentEndpoint[key as keyof Endpoint] === value) {
        console.log(`Value for '${key}' hasn't changed, skipping update`);
        return currentEndpoint;
    }

    // Update the endpoints update list for saving later
    let listToUpdate = get(endpointsUpdateList);
    let endpointUpdateIndex = listToUpdate.findIndex((e) => e.endpointId === currentEndpoint.id);

    // If this endpoint is already in the list, we'll update its value
    if (endpointUpdateIndex !== -1) {
        const updatedList = [...listToUpdate];
        if (key === 'all') {
            // Update the whole endpoint in the list
            updatedList[endpointUpdateIndex] = {
                ...updatedList[endpointUpdateIndex],
                projectId: currentEndpoint.project_id,
                endpointId: currentEndpoint.id,
                endpoint: { ...currentEndpoint } // Store the entire endpoint object
            };
        } else {
            // Check if the value is actually different than what's already stored
            const storedValue = updatedList[endpointUpdateIndex][key];
            if (storedValue === value) {
                console.log(`Value for '${key}' hasn't changed, skipping update`);
                return currentEndpoint;
            }

            // Update just the specific field
            updatedList[endpointUpdateIndex] = {
                ...updatedList[endpointUpdateIndex],
                [key]: value
            };
        }
        endpointsUpdateList.set(updatedList);
        showSaveButton.set(true);
    } else {
        // Add new entry to the list
        const newEntry: EndpointUpdate = {
            projectId: currentEndpoint.project_id,
            endpointId: currentEndpoint.id
        };

        if (key === 'all') {
            newEntry.endpoint = { ...currentEndpoint }; // Store the entire endpoint object
        } else {
            newEntry[key] = value;
        }

        endpointsUpdateList.set([...listToUpdate, newEntry]);
        showSaveButton.set(true);
    }

    console.log('Updated endpoints list', get(endpointsUpdateList));
    
    // Return the updated endpoint
    return getUpdatedEndpoint(key, value, currentEndpoint);
}

/**
 * Get a copy of the updated endpoint with the new value
 */
function getUpdatedEndpoint(key: string, value: any, currentEndpoint: Endpoint): Endpoint {
    // Create a new endpoint object with the updated value
    const updatedEndpoint = {
        ...currentEndpoint
    };

    // Handle special fields
    switch(key) {
        case 'use_proxy':
            updatedEndpoint.use_proxy = Boolean(value);
            // If disabling proxy, also clear the target
            if (!value && updatedEndpoint.proxy_target_id) {
                updatedEndpoint.proxy_target_id = null;
                updatedEndpoint.proxy_target = null;
            }
            break;
        case 'proxy_target_id':
            updatedEndpoint.proxy_target_id = value;
            break;
        default:
            // For standard fields, simply update the value
            (updatedEndpoint as any)[key] = value;
    }

    return updatedEndpoint;
}

/**
 * Update a response for an endpoint
 * @param key Field to update
 * @param value New value
 * @param currentEndpoint Endpoint containing the response
 * @param currentResponse Response to update
 * @returns The current response object (unchanged)
 */
export function updateResponse(key: string, value: any, currentEndpoint: Endpoint | null, currentResponse: Response | null) {
    /// check end point and response 
    if (!currentEndpoint || !currentResponse) {
        console.log('Invalid endpoint or response');
        return currentResponse;
    }

    let currentProject = get(selectedProject) as Project | null;
    if (!currentProject) {
        toast.error('No project selected');
        return currentResponse;
    }

    // Validate that key and value are provided
    if (!key) {
        toast.error('Invalid update: key is required');
        return currentResponse;
    }

    if (value === undefined || value === null) {
        toast.error(`Invalid update: value for '${key}' cannot be empty`);
        return currentResponse;
    }

    if (!currentResponse) {
        toast.error('Invalid update: response object is required');
        return currentResponse;
    }

    // Check if the value is different from what's already in the response
    // Skip update if the value is the same
    if (key !== 'all' && currentResponse[key as keyof Response] === value) {
        console.log(`Value for '${key}' hasn't changed, skipping update`);
        return currentResponse;
    }

    let listToUpdate = get(endpointsUpdateList);
    let responseUpdateIndex = listToUpdate.findIndex(
        (item) => 'responseId' in item && item.responseId === currentResponse.id && item.endpointId === currentEndpoint.id
    );

    // If this response is already in the list, we'll update its value
    if (responseUpdateIndex !== -1) {
        const updatedList = [...listToUpdate];
        if (key === 'all') {
            // Update the whole response in the list
            updatedList[responseUpdateIndex] = {
                ...updatedList[responseUpdateIndex],
                projectId: currentEndpoint.project_id,
                endpointId: currentEndpoint.id,
                responseId: currentResponse.id,
                response: { ...currentResponse } // Store the entire response object
            };
        } else {
            // Check if the value is actually different than what's already stored
            const storedValue = updatedList[responseUpdateIndex][key];
            if (storedValue === value) {
                console.log(`Value for '${key}' hasn't changed, skipping update`);
                return currentResponse;
            }

            // Update just the specific field
            updatedList[responseUpdateIndex] = {
                ...updatedList[responseUpdateIndex],
                [key]: value
            };
        }
        endpointsUpdateList.set(updatedList);
        showSaveButton.set(true);
    } else {
        // Add new entry to the list
        const newEntry: ResponseUpdate = {
            projectId: currentEndpoint.project_id,
            endpointId: currentEndpoint.id,
            responseId: currentResponse.id
        };

        if (key === 'all') {
            newEntry.response = { ...currentResponse }; // Store the entire response object
        } else {
            newEntry[key] = value;
        }

        endpointsUpdateList.set([...listToUpdate, newEntry]);
        showSaveButton.set(true);
    }

    console.log('Updated endpoints/responses list', get(endpointsUpdateList));
    return currentResponse;
}

/**
 * Reset the endpoints update list
 * Use this function when changing endpoints or projects
 */
export function resetEndpointsList() {
    endpointsUpdateList.set([]);
    showSaveButton.set(false);
}

/**
 * Get the current endpoints update list
 * This can be used to retrieve the list for saving to the server
 */
export function getEndpointsUpdateList() {
    return get(endpointsUpdateList);
}

/**
 * Update the route status in both the local state and the selected project store
 * @param route Updated endpoint
 */
export function updateRouteStatus(route: Endpoint) {
    // Get the current selected project
    const project = get(selectedProject);
    if (!project || !project.endpoints) {
        console.error('No project selected or project has no endpoints');
        return;
    }

    // Find the endpoint in the selected project
    const projectEndpointIndex = project.endpoints.findIndex((e) => e.id === route.id);
    if (projectEndpointIndex !== -1) {
        // Update the endpoint in the selected project
        project.endpoints[projectEndpointIndex] = {
            ...route
        };
        
        // Update the store
        selectedProject.set(project);
        console.log('Updated route status in selected project store', route);
    } else {
        console.warn('Endpoint not found in selected project', route);
    }
}

/**
 * Handle route change like the provided handleRouteStatusChange function
 * This is a more general utility function that can be used to update any endpoint
 * in both the local state and the selected project store
 * @param route Updated endpoint
 * @param endpoints Array of endpoints to update (local state)
 * @returns Updated endpoints array
 */
export function handleRouteChange(route: Endpoint, endpoints: Endpoint[]): Endpoint[] {
    console.log('Route changed:', route);
    const index = endpoints.findIndex((r) => r.id === route.id);
    if (index !== -1) {
        endpoints[index] = {
            ...route
        };
        endpoints = [...endpoints]; // Trigger reactivity with a new array reference

        // Also update in the selectedProject store
        updateRouteStatus(route);
    }
    return endpoints;
}