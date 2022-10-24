import { writable } from 'svelte/store';

export const dark_mode = writable(0);
export const loading = writable(0);

export const repositories = writable({})
export const snapshots = writable({})

export const selected_repo_id = writable(-1);