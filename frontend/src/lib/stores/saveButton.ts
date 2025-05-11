import { writable, get } from 'svelte/store';
import { toast } from './toast';
import type { Endpoint, Project } from '$lib/api/mockoonApi';
import { selectedProject } from './selectedConfig';

export const showSaveButton = writable<boolean>(false);
export const saveInprogress = writable<boolean>(false);

// Endpoint update list only contains project ID and endpoint ID
export interface EndpointUpdate {
    projectId: string;
    endpointId: string;
    [key: string]: any; // For additional key-value pairs
}
export const endpointsUpdateList = writable<EndpointUpdate[]>([]);

export function saveButtonHandler() {
    let currentValue = get(saveInprogress);
    setTimeout(() => {
        saveInprogress.set(!currentValue);
    }, 1000);
    console.log('saveButtonHandler', currentValue);
}

export function updateEndpoint(keyOrEndpoint: string | Endpoint, endpointObj?: Endpoint) {
    let currentProject = get(selectedProject) as Project | null;
    if (!currentProject) {
        toast.error('No project selected');
        return;
    }
    
    let key: string;
    let endpoint: Endpoint;
    
    // Handle the case where updateEndpoint is called with just the endpoint
    if (typeof keyOrEndpoint === 'object') {
        endpoint = keyOrEndpoint;
        // When called with just the endpoint, we update the whole endpoint
        key = 'all';
    } else {
        // When called with key and endpoint, we update only that specific field
        key = keyOrEndpoint;
        endpoint = endpointObj as Endpoint;
    }

    let listToUpdate = get(endpointsUpdateList);
    let endpointUpdateIndex = listToUpdate.findIndex((e) => e.endpointId === endpoint.id);
    
    // If this endpoint is already in the list, we'll update its value
    if (endpointUpdateIndex !== -1) {
        const updatedList = [...listToUpdate];
        if (key === 'all') {
            // Update the whole endpoint in the list
            updatedList[endpointUpdateIndex] = {
                ...updatedList[endpointUpdateIndex],
                projectId: endpoint.project_id,
                endpointId: endpoint.id,
                endpoint: { ...endpoint } // Store the entire endpoint object
            };
        } else {
            // Update just the specific field
            updatedList[endpointUpdateIndex] = {
                ...updatedList[endpointUpdateIndex],
                [key]: endpoint[key as keyof Endpoint]
            };
        }
        endpointsUpdateList.set(updatedList);
        showSaveButton.set(true);
    } else {
        // Add new entry to the list
        const newEntry: EndpointUpdate = {
            projectId: endpoint.project_id,
            endpointId: endpoint.id
        };
        
        if (key === 'all') {
            newEntry.endpoint = { ...endpoint }; // Store the entire endpoint object
        } else {
            newEntry[key] = endpoint[key as keyof Endpoint];
        }
        
        endpointsUpdateList.set([...listToUpdate, newEntry]);
        showSaveButton.set(true);
    }
    
    console.log('Updated endpoint list', get(endpointsUpdateList));
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