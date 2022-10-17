import { writable } from 'svelte/store';

export const dark_mode = writable(0);

export const selected_repo = writable({
    directory:"", 
    name: ""
});