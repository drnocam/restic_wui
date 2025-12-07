// @ts-ignore
import * as App from '$lib/wailsjs/go/main/App';

export const api = {
    // Repository
    openRepository: async (path: string, password: string) => {
        return await App.OpenRepository(path, password);
    },
    initRepository: async (path: string, password: string) => {
        return await App.InitRepository(path, password);
    },
    getRepositoryStats: async () => {
        return await App.GetRepositoryStats();
    },
    isRepositoryOpen: async () => {
        return await App.IsRepositoryOpen();
    },
    getCurrentRepositoryPath: async () => {
        return await App.GetCurrentRepositoryPath();
    },
    pruneRepository: async () => {
        return await App.PruneRepository();
    },

    // Snapshots
    listSnapshots: async () => {
        return await App.ListSnapshots();
    },
    forgetSnapshot: async (id: string, prune: boolean) => {
        return await App.ForgetSnapshot(id, prune);
    },
    restoreSnapshot: async (id: string, targetDir: string) => {
        return await App.RestoreSnapshot(id, targetDir);
    },
    listSnapshotFiles: async (id: string) => {
        return await App.ListSnapshotFiles(id);
    },

    // Helpers
    selectDirectory: async () => {
        return await App.SelectDirectory();
    }
};
