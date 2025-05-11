import { writable, get } from 'svelte/store';
import { toast } from './toast';
import type { Endpoint, Project, Response } from '$lib/api/mockoonApi';
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

export function saveButtonHandler() {
    let currentValue = get(saveInprogress);
    setTimeout(() => {
        saveInprogress.set(!currentValue);
    }, 1000);
    console.log('saveButtonHandler', currentValue);
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
        toast.error(`Invalid update: value for '${key}' cannot be empty`);
        return currentEndpoint;
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

    console.log('Updated endpoint list', get(endpointsUpdateList));
    return currentEndpoint;
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