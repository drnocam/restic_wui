import { writable } from 'svelte/store';

export const currentRepoPath = writable<string>('');
export const isRepoOpen = writable<boolean>(false);
export const loading = writable<boolean>(false);
export const error = writable<string | null>(null);

export interface SnapshotSummary {
    files_new?: number;
    files_changed?: number;
    files_unmodified?: number;
    dirs_new?: number;
    dirs_changed?: number;
    dirs_unmodified?: number;
    data_blobs?: number;
    tree_blobs?: number;
    data_added?: number;
    total_files_processed?: number;
    total_bytes_processed?: number;
}

export interface Snapshot {
    id: string;
    time: string;
    tree: string;
    paths: string[];
    hostname: string;
    username: string;
    tags: string[];
    short_id: string;
    summary?: SnapshotSummary;
}

export const snapshots = writable<Snapshot[]>([]);
export const repoStats = writable<any>(null);
