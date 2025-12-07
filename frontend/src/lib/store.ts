import { writable } from 'svelte/store';

export const currentRepoPath = writable<string>('');
export const isRepoOpen = writable<boolean>(false);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export interface Snapshot {
    id: string;
    time: string;
    tree: string;
    paths: string[];
    hostname: string;
    username: string;
    tags: string[];
    short_id: string;
}

export const snapshots = writable<Snapshot[]>([]);
