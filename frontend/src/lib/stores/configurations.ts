import { writable } from 'svelte/store';
import type { Project } from '$lib/api/BeoApi';
import { getProjects } from '$lib/api/BeoApi';
import { selectedProject } from '$lib/stores/selectedConfig';

export const projects = writable<Project[]>([]);

export async function fetchConfigsStore() {
  projects.set(await getProjects());
}

// Push the result of updateProject into every store that caches the project,
// so the Configuration tab does not re-mount with the old values after a tab switch.
// Only project-level fields are taken: endpoints stay as-is because the routes
// editor may hold unsaved changes there, and url/is_pinned are computed per request.
export function applyProjectUpdate(updated: Project) {
  const { endpoints: _endpoints, url: _url, is_pinned: _isPinned, ...fields } = updated;
  const merge = (current: Project): Project => {
    let url = current.url;
    if (url && fields.alias && fields.alias !== current.alias) {
      url = url.replace(/[^/]*\/?$/, fields.alias);
    }
    return { ...current, ...fields, url };
  };

  projects.update((list) => list.map((p) => (p.id === updated.id ? merge(p) : p)));
  selectedProject.update((p) => (p && p.id === updated.id ? merge(p) : p));
}
