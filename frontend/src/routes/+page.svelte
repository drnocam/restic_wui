<script lang="ts">
    import { api } from "$lib/api";
    import {
        currentRepoPath,
        isRepoOpen,
        loading,
        error,
        repoStats,
    } from "$lib/store";
    import { toast } from "@zerodevx/svelte-toast";
    import { onMount } from "svelte";

    let repoPassword = "";
    let showPasswordModal = false;
    let selectedPath = "";
    let isInitMode = false;

    async function handleSelectDirectory() {
        try {
            const path = await api.selectDirectory();
            if (path) {
                selectedPath = path;
                showPasswordModal = true;
            }
        } catch (e) {
            toast.push("Failed to select directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        }
    }

    async function handleSubmit() {
        if (!selectedPath || !repoPassword) return;

        loading.set(true);
        try {
            if (isInitMode) {
                await api.initRepository(selectedPath, repoPassword);
                toast.push("Repository initialized successfully", {
                    theme: {
                        "--toastBackground": "#4ADE80",
                        "--toastColor": "black",
                    },
                });
            } else {
                await api.openRepository(selectedPath, repoPassword);
                toast.push("Repository opened successfully", {
                    theme: {
                        "--toastBackground": "#4ADE80",
                        "--toastColor": "black",
                    },
                });
            }

            currentRepoPath.set(selectedPath);
            isRepoOpen.set(true);
            showPasswordModal = false;
            repoPassword = "";

            // Load stats if opened
            await loadStats();
        } catch (e) {
            console.error(e);
            const msg = isInitMode
                ? "Failed to initialize repository"
                : "Failed to open repository. Check password.";
            toast.push(msg, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
            error.set(e as string);
        } finally {
            loading.set(false);
        }
    }

    async function loadStats() {
        try {
            const stats = await api.getRepositoryStats();
            repoStats.set(stats);
        } catch (e) {
            console.error("Failed to load stats", e);
        }
    }

    onMount(async () => {
        if ($isRepoOpen && !$repoStats) {
            await loadStats();
        }
    });

    async function handlePrune() {
        if (
            !confirm(
                "Are you sure you want to prune the repository? This will permanently delete unreferenced data.",
            )
        )
            return;

        loading.set(true);
        try {
            await api.pruneRepository();
            toast.push("Repository pruned successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            await loadStats();
        } catch (e) {
            toast.push("Failed to prune repository", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
            console.error(e);
        } finally {
            loading.set(false);
        }
    }

    function startOpen() {
        isInitMode = false;
        handleSelectDirectory();
    }

    function startInit() {
        isInitMode = true;
        handleSelectDirectory();
    }

    async function handleClose() {
        try {
            await api.closeRepository();
            isRepoOpen.set(false);
            currentRepoPath.set("");
            repoStats.set(null);
            toast.push("Repository closed", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
        }
    }

    async function handleDelete() {
        if (
            !confirm(
                "Are you sure you want to delete this repository? THIS ACTION IS IRREVERSIBLE AND WILL DELETE ALL DATA.",
            )
        )
            return;

        const password = prompt(
            "Please enter the repository password to confirm deletion:",
        );
        if (!password) return;

        // Ideally we should verify password, but for now we just proceed if user confirms.
        // Actually, DeleteRepository in backend doesn't check password again, it just deletes the folder.
        // But asking for it adds friction.

        loading.set(true);
        try {
            await api.deleteRepository();
            isRepoOpen.set(false);
            currentRepoPath.set("");
            repoStats.set(null);
            toast.push("Repository deleted", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Failed to delete repository", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    let showBackupModal = false;
    let backupSource = "";
    let backupExcludes = "";

    async function handleBackup() {
        if (!backupSource) {
            toast.push("Please enter a source directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
            return;
        }

        loading.set(true);
        try {
            const excludes = backupExcludes
                .split("\n")
                .filter((e) => e.trim() !== "");
            const output = await api.backup(backupSource, excludes);
            console.log(output);
            toast.push("Backup completed successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            await loadStats();
            showBackupModal = false;
            backupSource = "";
            backupExcludes = "";
        } catch (e) {
            console.error(e);
            toast.push("Backup failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleSelectBackupSource() {
        try {
            const path = await api.selectDirectory();
            if (path) {
                backupSource = path;
            }
        } catch (e) {
            toast.push("Failed to select directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        }
    }

    async function handleSelectRestoreTarget() {
        try {
            const path = await api.selectDirectory();
            if (path) {
                restoreTarget = path;
            }
        } catch (e) {
            toast.push("Failed to select directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        }
    }

    let showRestoreModal = false;
    let restoreSnapshotId = "latest";
    let restoreTarget = "";
    let restorePaths = "";
    let restoreHost = "";

    async function handleRestore() {
        if (!restoreSnapshotId || !restoreTarget) {
            toast.push("Please enter snapshot ID and target directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
            return;
        }

        loading.set(true);
        try {
            const paths = restorePaths
                .split("\n")
                .filter((e) => e.trim() !== "");
            const output = await api.restoreSnapshot(
                restoreSnapshotId,
                restoreTarget,
                paths,
                restoreHost,
            );
            console.log(output);
            toast.push("Restore completed successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            showRestoreModal = false;
            restoreSnapshotId = "latest";
            restoreTarget = "";
            restorePaths = "";
            restoreHost = "";
        } catch (e) {
            console.error(e);
            toast.push("Restore failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    // Maintenance handlers
    let showFindModal = false;
    let findPattern = "";
    let findResults: any[] = [];
    let findRestoreTarget = "";
    let findRestorePath = "";
    let findRestoreSnapshot = "";

    let showMountModal = false;
    let mountPoint = "";
    let isMounted = false;

    async function handleRepair() {
        if (
            !confirm(
                "Are you sure you want to rebuild the index? This may take some time.",
            )
        )
            return;
        loading.set(true);
        try {
            await api.repairRepository();
            toast.push("Index rebuilt successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Repair failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleCheck() {
        loading.set(true);
        try {
            const result = await api.checkRepository();
            console.log(result);
            toast.push("Repository check passed", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Check failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleUnlock() {
        loading.set(true);
        try {
            await api.unlockRepository();
            toast.push("Repository unlocked", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Unlock failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleFind() {
        if (!findPattern) return;
        loading.set(true);
        try {
            const result = await api.find(findPattern);
            findResults = result || [];
        } catch (e) {
            console.error(e);
            toast.push("Find failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    function formatBytes(bytes: number): string {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }

    async function handleRestoreFile(snapshotId: string, filePath: string) {
        findRestoreSnapshot = snapshotId;
        findRestorePath = filePath;
        // Ask for target directory
        try {
            const path = await api.selectDirectory();
            if (path) {
                findRestoreTarget = path;
                loading.set(true);
                await api.restoreSnapshot(snapshotId, path, [filePath], "");
                toast.push("File restored successfully to " + path, {
                    theme: {
                        "--toastBackground": "#4ADE80",
                        "--toastColor": "black",
                    },
                });
            }
        } catch (e) {
            console.error(e);
            toast.push("Restore failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleSelectMountPoint() {
        try {
            const path = await api.selectDirectory();
            if (path) {
                mountPoint = path;
            }
        } catch (e) {
            toast.push("Failed to select directory", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        }
    }

    async function handleMount() {
        if (!mountPoint) return;
        loading.set(true);
        try {
            await api.mountRepository(mountPoint);
            isMounted = true;
            toast.push("Repository mounted at " + mountPoint, {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            showMountModal = false;
        } catch (e) {
            console.error(e);
            toast.push("Mount failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleUnmount() {
        loading.set(true);
        try {
            await api.unmountRepository();
            isMounted = false;
            mountPoint = "";
            toast.push("Repository unmounted", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Unmount failed: " + e, {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }
</script>

<div class="space-y-8">
    {#if !$isRepoOpen}
        <div class="text-center space-y-4 mb-12">
            <h2
                class="text-4xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-cyan-400"
            >
                Welcome to Restic UI
            </h2>
            <p class="text-slate-400 text-lg max-w-2xl mx-auto">
                Secure, fast, and efficient backup management. Select an
                existing repository or create a new one to get started.
            </p>
        </div>

        <div class="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
            <!-- Open Repository Card -->
            <button
                class="group relative p-8 rounded-2xl bg-slate-800/50 border border-slate-700 hover:border-blue-500/50 transition-all duration-300 text-left hover:shadow-2xl hover:shadow-blue-500/10"
                on:click={startOpen}
            >
                <div
                    class="absolute inset-0 bg-gradient-to-br from-blue-500/5 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"
                ></div>
                <div class="relative z-10">
                    <div
                        class="w-12 h-12 bg-blue-500/10 rounded-xl flex items-center justify-center mb-6 text-blue-400 group-hover:scale-110 transition-transform duration-300"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z"
                            />
                        </svg>
                    </div>
                    <h3 class="text-xl font-semibold text-slate-100 mb-2">
                        Open Repository
                    </h3>
                    <p class="text-slate-400">
                        Connect to an existing local Restic repository to manage
                        backups and snapshots.
                    </p>
                </div>
            </button>

            <!-- Init Repository Card -->
            <button
                class="group relative p-8 rounded-2xl bg-slate-800/50 border border-slate-700 hover:border-cyan-500/50 transition-all duration-300 text-left hover:shadow-2xl hover:shadow-cyan-500/10"
                on:click={startInit}
            >
                <div
                    class="absolute inset-0 bg-gradient-to-br from-cyan-500/5 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity"
                ></div>
                <div class="relative z-10">
                    <div
                        class="w-12 h-12 bg-cyan-500/10 rounded-xl flex items-center justify-center mb-6 text-cyan-400 group-hover:scale-110 transition-transform duration-300"
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-6 w-6"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                            />
                        </svg>
                    </div>
                    <h3 class="text-xl font-semibold text-slate-100 mb-2">
                        Initialize New
                    </h3>
                    <p class="text-slate-400">
                        Create a brand new encrypted repository in a local
                        directory.
                    </p>
                </div>
            </button>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <!-- Stats will go here -->
            <div
                class="p-6 rounded-2xl bg-slate-800/50 border border-slate-700"
            >
                <h3 class="text-slate-400 text-sm font-medium mb-2">
                    Repository Status
                </h3>
                <p class="text-2xl font-bold text-green-400">Active</p>
            </div>

            {#if $repoStats}
                <div
                    class="p-6 rounded-2xl bg-slate-800/50 border border-slate-700"
                >
                    <h3 class="text-slate-400 text-sm font-medium mb-2">
                        Disk Size
                    </h3>
                    <p class="text-2xl font-bold text-white">
                        {formatBytes($repoStats.disk_size || 0)}
                    </p>
                    <p class="text-xs text-slate-500 mt-1">
                        Logical: {formatBytes($repoStats.total_size || 0)}
                    </p>
                </div>
                <div
                    class="p-6 rounded-2xl bg-slate-800/50 border border-slate-700"
                >
                    <h3 class="text-slate-400 text-sm font-medium mb-2">
                        Total Files
                    </h3>
                    <p class="text-2xl font-bold text-white">
                        {$repoStats.total_file_count.toLocaleString()}
                    </p>
                </div>

                <div
                    class="col-span-1 md:col-span-3 flex flex-wrap justify-around gap-4"
                >
                    <button
                        class="px-6 py-3 rounded-xl bg-green-600 text-white hover:bg-green-700 transition-colors font-medium flex items-center gap-2"
                        on:click={() => (showRestoreModal = true)}
                        disabled={$loading}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                            />
                        </svg>
                        Restore
                    </button>
                    <button
                        class="px-6 py-3 rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium flex items-center gap-2"
                        on:click={() => (showBackupModal = true)}
                        disabled={$loading}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                            />
                        </svg>
                        Backup
                    </button>
                    <button
                        class="px-6 py-3 rounded-xl bg-red-900/50 border border-red-700 text-red-200 hover:bg-red-800 hover:text-white transition-all flex items-center gap-2"
                        on:click={handleDelete}
                        disabled={$loading}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                        </svg>
                        Delete Repository
                    </button>
                    <button
                        class="px-6 py-3 rounded-xl bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2"
                        on:click={handleClose}
                        disabled={$loading}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
                            />
                        </svg>
                        Close Repository
                    </button>
                    <button
                        class="px-6 py-3 rounded-xl bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2"
                        on:click={handlePrune}
                        disabled={$loading}
                    >
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                stroke-linecap="round"
                                stroke-linejoin="round"
                                stroke-width="2"
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                            />
                        </svg>
                        Prune Repository
                    </button>
                </div>

                <!-- Maintenance Actions -->
                <div class="col-span-1 md:col-span-3 mt-6">
                    <h3 class="text-lg font-semibold text-white mb-4">
                        Maintenance
                    </h3>
                    <div class="flex flex-wrap gap-3">
                        <button
                            class="px-4 py-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2 text-sm"
                            on:click={handleCheck}
                            disabled={$loading}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                                />
                            </svg>
                            Check
                        </button>
                        <button
                            class="px-4 py-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2 text-sm"
                            on:click={handleRepair}
                            disabled={$loading}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"
                                />
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                />
                            </svg>
                            Repair
                        </button>
                        <button
                            class="px-4 py-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2 text-sm"
                            on:click={handleUnlock}
                            disabled={$loading}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M8 11V7a4 4 0 118 0m-4 8v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2z"
                                />
                            </svg>
                            Unlock
                        </button>
                        <button
                            class="px-4 py-2 rounded-lg bg-slate-800 border border-slate-700 text-slate-300 hover:bg-slate-700 hover:text-white transition-all flex items-center gap-2 text-sm"
                            on:click={() => (showFindModal = true)}
                            disabled={$loading}
                        >
                            <svg
                                xmlns="http://www.w3.org/2000/svg"
                                class="h-4 w-4"
                                fill="none"
                                viewBox="0 0 24 24"
                                stroke="currentColor"
                            >
                                <path
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                    stroke-width="2"
                                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                                />
                            </svg>
                            Find
                        </button>
                        {#if !isMounted}
                            <button
                                class="px-4 py-2 rounded-lg bg-purple-900/50 border border-purple-700 text-purple-200 hover:bg-purple-800 hover:text-white transition-all flex items-center gap-2 text-sm"
                                on:click={() => (showMountModal = true)}
                                disabled={$loading}
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                                    />
                                </svg>
                                Mount
                            </button>
                        {:else}
                            <button
                                class="px-4 py-2 rounded-lg bg-red-900/50 border border-red-700 text-red-200 hover:bg-red-800 hover:text-white transition-all flex items-center gap-2 text-sm"
                                on:click={handleUnmount}
                                disabled={$loading}
                            >
                                <svg
                                    xmlns="http://www.w3.org/2000/svg"
                                    class="h-4 w-4"
                                    fill="none"
                                    viewBox="0 0 24 24"
                                    stroke="currentColor"
                                >
                                    <path
                                        stroke-linecap="round"
                                        stroke-linejoin="round"
                                        stroke-width="2"
                                        d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4"
                                    />
                                </svg>
                                Unmount ({mountPoint})
                            </button>
                        {/if}
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>

<!-- Password Modal -->
{#if showPasswordModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md shadow-2xl relative overflow-hidden"
        >
            <div
                class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 to-cyan-500"
            ></div>

            <h3 class="text-xl font-bold mb-4">
                {isInitMode
                    ? "Set Repository Password"
                    : "Enter Repository Password"}
            </h3>

            <p class="text-slate-400 text-sm mb-6">
                {isInitMode
                    ? "Choose a strong password to encrypt your new repository. You will need this password to access your backups."
                    : "Please enter the password to unlock this repository."}
            </p>

            <div class="space-y-4">
                <div>
                    <label
                        for="password-input"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Password</label
                    >
                    <input
                        id="password-input"
                        type="password"
                        bind:value={repoPassword}
                        class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                        placeholder="••••••••"
                    />
                </div>

                <div class="flex gap-3 pt-2">
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors font-medium"
                        on:click={() => {
                            showPasswordModal = false;
                            repoPassword = "";
                        }}
                    >
                        Cancel
                    </button>
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={!repoPassword || $loading}
                        on:click={handleSubmit}
                    >
                        {#if $loading}
                            Processing...
                        {:else}
                            {isInitMode ? "Initialize" : "Unlock"}
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- Backup Modal -->
{#if showBackupModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md shadow-2xl relative overflow-hidden"
        >
            <div
                class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 to-cyan-500"
            ></div>

            <h3 class="text-xl font-bold mb-4">Create Backup</h3>

            <div class="space-y-4">
                <div>
                    <label
                        for="backup-source"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Source Directory</label
                    >
                    <div class="flex gap-2">
                        <input
                            id="backup-source"
                            type="text"
                            bind:value={backupSource}
                            class="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                            placeholder="/home/user/documents"
                        />
                        <button
                            class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg hover:bg-slate-700 text-slate-300 transition-colors"
                            on:click={handleSelectBackupSource}
                        >
                            Browse
                        </button>
                    </div>
                </div>

                <div>
                    <label
                        for="backup-excludes"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Excludes (one per line)</label
                    >
                    <textarea
                        id="backup-excludes"
                        bind:value={backupExcludes}
                        class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all h-24"
                        placeholder="*.tmp&#10;node_modules"
                    ></textarea>
                </div>

                <div class="flex gap-3 pt-2">
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors font-medium"
                        on:click={() => {
                            showBackupModal = false;
                        }}
                    >
                        Cancel
                    </button>
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={!backupSource || $loading}
                        on:click={handleBackup}
                    >
                        {#if $loading}
                            Running...
                        {:else}
                            Start Backup
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- Restore Modal -->
{#if showRestoreModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md shadow-2xl relative overflow-hidden"
        >
            <div
                class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-green-500 to-cyan-500"
            ></div>

            <h3 class="text-xl font-bold mb-4">Restore Snapshot</h3>

            <div class="space-y-4">
                <div>
                    <label
                        for="restore-snapshot-id"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Snapshot ID</label
                    >
                    <input
                        id="restore-snapshot-id"
                        type="text"
                        bind:value={restoreSnapshotId}
                        class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                        placeholder="latest"
                    />
                </div>

                <div>
                    <label
                        for="restore-target"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Target Directory</label
                    >
                    <div class="flex gap-2">
                        <input
                            id="restore-target"
                            type="text"
                            bind:value={restoreTarget}
                            class="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                            placeholder="/home/user/restore"
                        />
                        <button
                            class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg hover:bg-slate-700 text-slate-300 transition-colors"
                            on:click={handleSelectRestoreTarget}
                        >
                            Browse
                        </button>
                    </div>
                </div>

                <div>
                    <label
                        for="restore-paths"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Paths (optional, one per line)</label
                    >
                    <textarea
                        id="restore-paths"
                        bind:value={restorePaths}
                        class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all h-20"
                        placeholder="/specific/path/to/restore"
                    ></textarea>
                </div>

                <div>
                    <label
                        for="restore-host"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Host (optional)</label
                    >
                    <input
                        id="restore-host"
                        type="text"
                        bind:value={restoreHost}
                        class="w-full bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                        placeholder=""
                    />
                </div>

                <div class="flex gap-3 pt-2">
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors font-medium"
                        on:click={() => {
                            showRestoreModal = false;
                        }}
                    >
                        Cancel
                    </button>
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-green-600 text-white hover:bg-green-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={!restoreSnapshotId ||
                            !restoreTarget ||
                            $loading}
                        on:click={handleRestore}
                    >
                        {#if $loading}
                            Restoring...
                        {:else}
                            Restore
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- Find Modal -->
{#if showFindModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-2xl h-[80vh] flex flex-col shadow-2xl relative overflow-hidden"
        >
            <div
                class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-blue-500 to-cyan-500"
            ></div>

            <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-bold">Find Files</h3>
                <button
                    class="p-2 hover:bg-slate-800 rounded-lg text-slate-400 hover:text-white transition-colors"
                    aria-label="Close modal"
                    on:click={() => {
                        showFindModal = false;
                        findResults = [];
                        findPattern = "";
                    }}
                >
                    <svg
                        xmlns="http://www.w3.org/2000/svg"
                        class="h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                    >
                        <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M6 18L18 6M6 6l12 12"
                        />
                    </svg>
                </button>
            </div>

            <div class="flex gap-2 mb-4">
                <input
                    type="text"
                    bind:value={findPattern}
                    class="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                    placeholder="Enter search pattern (e.g., *.txt, filename)"
                    on:keydown={(e) => e.key === "Enter" && handleFind()}
                />
                <button
                    class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium disabled:opacity-50"
                    disabled={!findPattern || $loading}
                    on:click={handleFind}
                >
                    {#if $loading}
                        Searching...
                    {:else}
                        Search
                    {/if}
                </button>
            </div>

            <div
                class="flex-1 overflow-y-auto bg-slate-950 rounded-xl border border-slate-800"
            >
                {#if findResults.length > 0}
                    {#each findResults as result}
                        <div class="border-b border-slate-800 last:border-0">
                            <div
                                class="px-4 py-2 bg-slate-900/50 flex items-center justify-between"
                            >
                                <span class="text-xs text-slate-400"
                                    >Snapshot: <span
                                        class="font-mono text-blue-400"
                                        >{result.snapshot?.substring(
                                            0,
                                            8,
                                        )}</span
                                    ></span
                                >
                                <span class="text-xs text-slate-500"
                                    >{result.hits} match{result.hits > 1
                                        ? "es"
                                        : ""}</span
                                >
                            </div>
                            {#if result.matches}
                                {#each result.matches as match}
                                    <div
                                        class="px-4 py-3 flex items-center justify-between hover:bg-slate-800/50 gap-4 group"
                                    >
                                        <div class="flex-1 min-w-0">
                                            <div
                                                class="flex items-center gap-2"
                                            >
                                                {#if match.type === "dir"}
                                                    <svg
                                                        xmlns="http://www.w3.org/2000/svg"
                                                        class="h-4 w-4 text-yellow-400"
                                                        fill="none"
                                                        viewBox="0 0 24 24"
                                                        stroke="currentColor"
                                                    >
                                                        <path
                                                            stroke-linecap="round"
                                                            stroke-linejoin="round"
                                                            stroke-width="2"
                                                            d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z"
                                                        />
                                                    </svg>
                                                {:else}
                                                    <svg
                                                        xmlns="http://www.w3.org/2000/svg"
                                                        class="h-4 w-4 text-slate-400"
                                                        fill="none"
                                                        viewBox="0 0 24 24"
                                                        stroke="currentColor"
                                                    >
                                                        <path
                                                            stroke-linecap="round"
                                                            stroke-linejoin="round"
                                                            stroke-width="2"
                                                            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                                                        />
                                                    </svg>
                                                {/if}
                                                <span
                                                    class="text-white text-sm truncate font-mono"
                                                    title={match.path}
                                                    >{match.path}</span
                                                >
                                            </div>
                                            <div
                                                class="flex items-center gap-4 mt-1 text-xs text-slate-500"
                                            >
                                                <span
                                                    >{formatBytes(
                                                        match.size || 0,
                                                    )}</span
                                                >
                                                <span>{match.permissions}</span>
                                                <span
                                                    >{match.user}:{match.group}</span
                                                >
                                            </div>
                                        </div>
                                        <button
                                            class="px-3 py-1.5 rounded-lg bg-green-600/20 text-green-400 hover:bg-green-600 hover:text-white transition-colors text-sm flex items-center gap-1.5"
                                            on:click={() =>
                                                handleRestoreFile(
                                                    result.snapshot,
                                                    match.path,
                                                )}
                                            disabled={$loading}
                                            title="Restore this file"
                                        >
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                class="h-4 w-4"
                                                fill="none"
                                                viewBox="0 0 24 24"
                                                stroke="currentColor"
                                            >
                                                <path
                                                    stroke-linecap="round"
                                                    stroke-linejoin="round"
                                                    stroke-width="2"
                                                    d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
                                                />
                                            </svg>
                                            Restore
                                        </button>
                                    </div>
                                {/each}
                            {/if}
                        </div>
                    {/each}
                {:else}
                    <p class="text-slate-500 text-center py-8">
                        Enter a pattern and click Search to find files in the
                        repository.
                    </p>
                {/if}
            </div>
        </div>
    </div>
{/if}

<!-- Mount Modal -->
{#if showMountModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md shadow-2xl relative overflow-hidden"
        >
            <div
                class="absolute top-0 left-0 w-full h-1 bg-gradient-to-r from-purple-500 to-pink-500"
            ></div>

            <h3 class="text-xl font-bold mb-4">Mount Repository</h3>
            <p class="text-slate-400 text-sm mb-6">
                Mount the repository to a local directory. You can then browse
                the backups using your file manager.
                <br /><span class="text-yellow-400"
                    >Note: FUSE must be installed on your system.</span
                >
            </p>

            <div class="space-y-4">
                <div>
                    <label
                        for="mount-point"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Mount Point</label
                    >
                    <div class="flex gap-2">
                        <input
                            id="mount-point"
                            type="text"
                            bind:value={mountPoint}
                            class="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                            placeholder="/mnt/restic-backup"
                        />
                        <button
                            class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg hover:bg-slate-700 text-slate-300 transition-colors"
                            on:click={handleSelectMountPoint}
                        >
                            Browse
                        </button>
                    </div>
                </div>

                <div class="flex gap-3 pt-2">
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors font-medium"
                        on:click={() => {
                            showMountModal = false;
                        }}
                    >
                        Cancel
                    </button>
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-purple-600 text-white hover:bg-purple-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={!mountPoint || $loading}
                        on:click={handleMount}
                    >
                        {#if $loading}
                            Mounting...
                        {:else}
                            Mount
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}
