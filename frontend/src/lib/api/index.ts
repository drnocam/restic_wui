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
    closeRepository: async () => {
        return await App.CloseRepository();
    },
    deleteRepository: async () => {
        return await App.DeleteRepository();
    },
    backup: async (source: string, excludes: string[]) => {
        return await App.Backup(source, excludes);
    },

    // Snapshots
    listSnapshots: async () => {
        return await App.ListSnapshots();
    },
    forgetSnapshot: async (id: string, prune: boolean) => {
        return await App.ForgetSnapshot(id, prune);
    },
    restoreSnapshot: async (id: string, targetDir: string, paths: string[], host: string) => {
        return await App.RestoreSnapshot(id, targetDir, paths, host);
    },
    listSnapshotFiles: async (id: string) => {
        return await App.ListSnapshotFiles(id);
    },

    // Helpers
    selectDirectory: async () => {
        return await App.SelectDirectory();
    },

    // Settings
    getSettings: async () => {
        return await App.GetSettings();
    },
    saveSettings: async (config: any) => {
        return await App.SaveSettings(config);
    },

    // Maintenance
    repairRepository: async () => {
        return await App.RepairRepository();
    },
    checkRepository: async () => {
        return await App.CheckRepository();
    },
    unlockRepository: async () => {
        return await App.UnlockRepository();
    },
    find: async (pattern: string) => {
        return await App.Find(pattern);
    },
    mountRepository: async (mountPoint: string) => {
        return await App.MountRepository(mountPoint);
    },
    unmountRepository: async () => {
        return await App.UnmountRepository();
    },
    isMounted: async () => {
        return await App.IsMounted();
    }
};
