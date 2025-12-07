<script lang="ts">
    import { onMount } from "svelte";
    import { api } from "$lib/api";
    import { snapshots, loading } from "$lib/store";
    import { toast } from "@zerodevx/svelte-toast";
    import type { Snapshot } from "$lib/store";

    let searchTerm = "";
    let selectedSnapshot: Snapshot | null = null;
    let showRestoreModal = false;
    let showFilesModal = false;
    let snapshotFiles: any[] = [];
    let restoreTargetDir = "";

    // Pagination
    let currentPage = 1;
    let itemsPerPage = 10;

    $: filteredSnapshots = $snapshots.filter(
        (s) =>
            s.short_id.includes(searchTerm) ||
            s.paths.some((p) => p.includes(searchTerm)) ||
            s.tags?.some((t) => t.includes(searchTerm)),
    );

    $: totalPages = Math.ceil(filteredSnapshots.length / itemsPerPage);
    $: paginatedSnapshots = filteredSnapshots.slice(
        (currentPage - 1) * itemsPerPage,
        currentPage * itemsPerPage,
    );

    // Reset page when search changes
    $: if (searchTerm) {
        currentPage = 1;
    }

    function formatRelativeTime(dateStr: string): string {
        const date = new Date(dateStr);
        const now = new Date();
        const diffMs = now.getTime() - date.getTime();
        const diffSecs = Math.floor(diffMs / 1000);
        const diffMins = Math.floor(diffSecs / 60);
        const diffHours = Math.floor(diffMins / 60);
        const diffDays = Math.floor(diffHours / 24);
        const diffWeeks = Math.floor(diffDays / 7);
        const diffMonths = Math.floor(diffDays / 30);
        const diffYears = Math.floor(diffDays / 365);

        if (diffYears > 0)
            return `${diffYears} year${diffYears > 1 ? "s" : ""} ago`;
        if (diffMonths > 0)
            return `${diffMonths} month${diffMonths > 1 ? "s" : ""} ago`;
        if (diffWeeks > 0)
            return `${diffWeeks} week${diffWeeks > 1 ? "s" : ""} ago`;
        if (diffDays > 0)
            return `${diffDays} day${diffDays > 1 ? "s" : ""} ago`;
        if (diffHours > 0)
            return `${diffHours} hour${diffHours > 1 ? "s" : ""} ago`;
        if (diffMins > 0)
            return `${diffMins} minute${diffMins > 1 ? "s" : ""} ago`;
        return "Just now";
    }

    async function loadSnapshots() {
        loading.set(true);
        try {
            const data = await api.listSnapshots();
            // Sort by time descending
            data.sort(
                (a: any, b: any) =>
                    new Date(b.time).getTime() - new Date(a.time).getTime(),
            );
            snapshots.set(data);
        } catch (e) {
            toast.push("Failed to load snapshots", {
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

    async function handleForget(id: string) {
        if (
            !confirm(
                "Are you sure you want to delete this snapshot? This action cannot be undone.",
            )
        )
            return;

        loading.set(true);
        try {
            await api.forgetSnapshot(id, true); // Prune by default for now
            toast.push("Snapshot deleted successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            await loadSnapshots();
        } catch (e) {
            toast.push("Failed to delete snapshot", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading.set(false);
        }
    }

    async function handleRestoreSelectDir() {
        try {
            const path = await api.selectDirectory();
            if (path) {
                restoreTargetDir = path;
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

    async function handleRestore() {
        if (!selectedSnapshot || !restoreTargetDir) return;

        loading.set(true);
        try {
            await api.restoreSnapshot(
                selectedSnapshot.id,
                restoreTargetDir,
                [],
                "",
            );
            toast.push("Snapshot restored successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
            showRestoreModal = false;
            restoreTargetDir = "";
            selectedSnapshot = null;
        } catch (e) {
            toast.push("Failed to restore snapshot", {
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

    async function loadSnapshotFiles(snapshot: Snapshot) {
        selectedSnapshot = snapshot;
        loading.set(true);
        try {
            const data = await api.listSnapshotFiles(snapshot.id);
            // Now data is already an array of LSNode objects
            snapshotFiles = data || [];
            showFilesModal = true;
        } catch (e) {
            toast.push("Failed to load snapshot files", {
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

    function formatBytes(bytes: number): string {
        if (bytes === 0) return "0 B";
        const k = 1024;
        const sizes = ["B", "KB", "MB", "GB", "TB"];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i];
    }

    async function handleRestoreFile(filePath: string) {
        if (!selectedSnapshot) return;
        try {
            const targetDir = await api.selectDirectory();
            if (targetDir) {
                loading.set(true);
                await api.restoreSnapshot(
                    selectedSnapshot.id,
                    targetDir,
                    [filePath],
                    "",
                );
                toast.push("File restored successfully to " + targetDir, {
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

    function openRestoreModal(snapshot: Snapshot) {
        selectedSnapshot = snapshot;
        showRestoreModal = true;
    }

    function closeFilesModal() {
        showFilesModal = false;
        snapshotFiles = [];
        selectedSnapshot = null;
    }

    function goToPage(page: number) {
        if (page >= 1 && page <= totalPages) {
            currentPage = page;
        }
    }

    onMount(() => {
        loadSnapshots();
    });
</script>

<div class="space-y-6">
    <div class="flex items-center justify-between">
        <h2 class="text-2xl font-bold text-white">Snapshots</h2>
        <div class="flex gap-4">
            <div class="relative">
                <input
                    type="text"
                    bind:value={searchTerm}
                    placeholder="Search snapshots..."
                    class="bg-slate-800 border border-slate-700 rounded-lg pl-10 pr-4 py-2 text-sm text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none w-64"
                />
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-4 w-4 text-slate-400 absolute left-3 top-3"
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
            </div>
            <select
                bind:value={itemsPerPage}
                on:change={() => (currentPage = 1)}
                class="bg-slate-800 border border-slate-700 rounded-lg px-3 py-2 text-sm text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
            >
                <option value={5}>5 per page</option>
                <option value={10}>10 per page</option>
                <option value={25}>25 per page</option>
                <option value={50}>50 per page</option>
            </select>
            <button
                class="p-2 bg-slate-800 border border-slate-700 rounded-lg hover:bg-slate-700 text-slate-300 transition-colors"
                on:click={loadSnapshots}
                title="Refresh"
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
            </button>
        </div>
    </div>

    {#if $loading && $snapshots.length === 0}
        <div class="flex justify-center py-12">
            <div
                class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-500"
            ></div>
        </div>
    {:else if filteredSnapshots.length === 0}
        <div
            class="text-center py-12 bg-slate-800/30 rounded-2xl border border-slate-800 border-dashed"
        >
            <p class="text-slate-400">No snapshots found.</p>
        </div>
    {:else}
        <!-- Table View -->
        <div
            class="bg-slate-800/50 border border-slate-700 rounded-xl overflow-x-auto"
        >
            <table class="w-full min-w-[1000px]">
                <thead class="bg-slate-800 border-b border-slate-700">
                    <tr>
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >ID</th
                        >
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >Time</th
                        >
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >Paths</th
                        >
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >Host</th
                        >
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >Size</th
                        >
                        <th
                            class="text-left text-sm font-medium text-slate-400 px-4 py-3"
                            >Tags</th
                        >
                        <th
                            class="text-right text-sm font-medium text-slate-400 px-4 py-3"
                            >Actions</th
                        >
                    </tr>
                </thead>
                <tbody>
                    {#each paginatedSnapshots as snapshot, i}
                        <tr
                            class="border-b border-slate-700/50 hover:bg-slate-800/50 transition-colors"
                        >
                            <td class="px-4 py-3">
                                <span
                                    class="font-mono text-blue-400 bg-blue-400/10 px-2 py-0.5 rounded text-sm"
                                >
                                    {snapshot.short_id}
                                </span>
                            </td>
                            <td class="px-4 py-3">
                                <div class="flex flex-col">
                                    <span class="text-white text-sm"
                                        >{formatRelativeTime(
                                            snapshot.time,
                                        )}</span
                                    >
                                    <span class="text-slate-500 text-xs"
                                        >{new Date(
                                            snapshot.time,
                                        ).toLocaleString()}</span
                                    >
                                </div>
                            </td>
                            <td class="px-4 py-3">
                                <div class="flex flex-wrap gap-1">
                                    {#each snapshot.paths.slice(0, 2) as path}
                                        <span
                                            class="text-slate-300 text-sm bg-slate-800 px-2 py-0.5 rounded border border-slate-700/50 truncate max-w-[200px]"
                                            title={path}
                                        >
                                            {path}
                                        </span>
                                    {/each}
                                    {#if snapshot.paths.length > 2}
                                        <span class="text-slate-400 text-sm"
                                            >+{snapshot.paths.length - 2}</span
                                        >
                                    {/if}
                                </div>
                            </td>
                            <td class="px-4 py-3">
                                <span class="text-slate-300 text-sm"
                                    >{snapshot.hostname || "-"}</span
                                >
                            </td>
                            <td class="px-4 py-3">
                                <span class="text-slate-300 text-sm">
                                    {#if snapshot.summary?.total_bytes_processed}
                                        {formatBytes(
                                            snapshot.summary
                                                .total_bytes_processed,
                                        )}
                                    {:else}
                                        -
                                    {/if}
                                </span>
                            </td>
                            <td class="px-4 py-3">
                                <div class="flex gap-1 flex-wrap">
                                    {#if snapshot.tags && snapshot.tags.length > 0}
                                        {#each snapshot.tags.slice(0, 2) as tag}
                                            <span
                                                class="text-xs text-cyan-400 bg-cyan-400/10 px-2 py-0.5 rounded-full"
                                            >
                                                #{tag}
                                            </span>
                                        {/each}
                                        {#if snapshot.tags.length > 2}
                                            <span class="text-slate-400 text-xs"
                                                >+{snapshot.tags.length -
                                                    2}</span
                                            >
                                        {/if}
                                    {:else}
                                        <span class="text-slate-500 text-sm"
                                            >-</span
                                        >
                                    {/if}
                                </div>
                            </td>
                            <td class="px-4 py-3">
                                <div
                                    class="flex items-center justify-end gap-1"
                                >
                                    <button
                                        class="p-2 hover:bg-slate-700 rounded-lg text-slate-400 hover:text-white transition-colors"
                                        title="View Files"
                                        on:click={() =>
                                            loadSnapshotFiles(snapshot)}
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
                                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                                            />
                                            <path
                                                stroke-linecap="round"
                                                stroke-linejoin="round"
                                                stroke-width="2"
                                                d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                                            />
                                        </svg>
                                    </button>
                                    <button
                                        class="p-2 hover:bg-slate-700 rounded-lg text-slate-400 hover:text-white transition-colors"
                                        title="Restore"
                                        on:click={() =>
                                            openRestoreModal(snapshot)}
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
                                    </button>
                                    <button
                                        class="p-2 hover:bg-red-900/30 rounded-lg text-slate-400 hover:text-red-400 transition-colors"
                                        title="Delete"
                                        on:click={() =>
                                            handleForget(snapshot.id)}
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
                                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                                            />
                                        </svg>
                                    </button>
                                </div>
                            </td>
                        </tr>
                    {/each}
                </tbody>
            </table>
        </div>

        <!-- Pagination Controls -->
        {#if totalPages > 1}
            <div class="flex items-center justify-between">
                <span class="text-sm text-slate-400">
                    Showing {(currentPage - 1) * itemsPerPage + 1} to {Math.min(
                        currentPage * itemsPerPage,
                        filteredSnapshots.length,
                    )} of {filteredSnapshots.length} snapshots
                </span>
                <div class="flex items-center gap-2">
                    <button
                        class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors {currentPage ===
                        1
                            ? 'bg-slate-800 text-slate-500 cursor-not-allowed'
                            : 'bg-slate-800 text-slate-300 hover:bg-slate-700'}"
                        on:click={() => goToPage(currentPage - 1)}
                        disabled={currentPage === 1}
                    >
                        Previous
                    </button>
                    {#each Array(totalPages) as _, i}
                        {#if i + 1 === 1 || i + 1 === totalPages || (i + 1 >= currentPage - 1 && i + 1 <= currentPage + 1)}
                            <button
                                class="w-8 h-8 rounded-lg text-sm font-medium transition-colors {currentPage ===
                                i + 1
                                    ? 'bg-blue-600 text-white'
                                    : 'bg-slate-800 text-slate-300 hover:bg-slate-700'}"
                                on:click={() => goToPage(i + 1)}
                            >
                                {i + 1}
                            </button>
                        {:else if i + 1 === currentPage - 2 || i + 1 === currentPage + 2}
                            <span class="text-slate-500">...</span>
                        {/if}
                    {/each}
                    <button
                        class="px-3 py-1.5 rounded-lg text-sm font-medium transition-colors {currentPage ===
                        totalPages
                            ? 'bg-slate-800 text-slate-500 cursor-not-allowed'
                            : 'bg-slate-800 text-slate-300 hover:bg-slate-700'}"
                        on:click={() => goToPage(currentPage + 1)}
                        disabled={currentPage === totalPages}
                    >
                        Next
                    </button>
                </div>
            </div>
        {/if}
    {/if}
</div>

<!-- Restore Modal -->
{#if showRestoreModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-md shadow-2xl"
        >
            <h3 class="text-xl font-bold mb-4">Restore Snapshot</h3>
            <p class="text-slate-400 text-sm mb-6">
                Select a destination directory to restore the contents of
                snapshot <span class="font-mono text-blue-400"
                    >{selectedSnapshot?.short_id}</span
                >.
            </p>

            <div class="space-y-4">
                <div>
                    <label
                        for="dest-dir-input"
                        class="block text-sm font-medium text-slate-300 mb-1"
                        >Destination Directory</label
                    >
                    <div class="flex gap-2">
                        <input
                            id="dest-dir-input"
                            type="text"
                            bind:value={restoreTargetDir}
                            readonly
                            class="flex-1 bg-slate-800 border border-slate-700 rounded-lg px-4 py-2 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none"
                            placeholder="Select directory..."
                        />
                        <button
                            class="px-4 py-2 bg-slate-800 border border-slate-700 rounded-lg hover:bg-slate-700 text-slate-300 transition-colors"
                            on:click={handleRestoreSelectDir}
                        >
                            Browse
                        </button>
                    </div>
                </div>

                <div class="flex gap-3 pt-4">
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors font-medium"
                        on:click={() => {
                            showRestoreModal = false;
                            restoreTargetDir = "";
                        }}
                    >
                        Cancel
                    </button>
                    <button
                        class="flex-1 px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
                        disabled={!restoreTargetDir || $loading}
                        on:click={handleRestore}
                    >
                        {#if $loading}
                            Restoring...
                        {:else}
                            Restore Files
                        {/if}
                    </button>
                </div>
            </div>
        </div>
    </div>
{/if}

<!-- Files Modal -->
{#if showFilesModal}
    <div
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm"
    >
        <div
            class="bg-slate-900 border border-slate-700 rounded-2xl p-6 w-full max-w-4xl h-[80vh] flex flex-col shadow-2xl"
        >
            <div class="flex items-center justify-between mb-4">
                <h3 class="text-xl font-bold">
                    Files in <span class="font-mono text-blue-400"
                        >{selectedSnapshot?.short_id}</span
                    >
                </h3>
                <button
                    class="p-2 hover:bg-slate-800 rounded-lg text-slate-400 hover:text-white transition-colors"
                    on:click={closeFilesModal}
                    aria-label="Close modal"
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

            <div
                class="flex-1 overflow-y-auto bg-slate-950 rounded-xl border border-slate-800 p-4 font-mono text-sm"
            >
                {#if snapshotFiles.length === 0}
                    <p class="text-slate-500 text-center py-8">
                        No files found or empty snapshot.
                    </p>
                {:else}
                    <div class="space-y-1">
                        {#each snapshotFiles as file}
                            <div
                                class="flex items-center justify-between px-3 py-2 hover:bg-slate-800/50 rounded-lg gap-4"
                            >
                                <div class="flex-1 min-w-0">
                                    <div class="flex items-center gap-2">
                                        {#if file.type === "dir"}
                                            <svg
                                                xmlns="http://www.w3.org/2000/svg"
                                                class="h-4 w-4 text-yellow-400 flex-shrink-0"
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
                                                class="h-4 w-4 text-slate-400 flex-shrink-0"
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
                                            class="text-white text-sm truncate"
                                            title={file.path}>{file.path}</span
                                        >
                                    </div>
                                    <div
                                        class="flex items-center gap-4 mt-1 text-xs text-slate-500 ml-6"
                                    >
                                        {#if file.size}
                                            <span>{formatBytes(file.size)}</span
                                            >
                                        {/if}
                                        {#if file.permissions}
                                            <span>{file.permissions}</span>
                                        {/if}
                                    </div>
                                </div>
                                <button
                                    class="px-3 py-1.5 rounded-lg bg-green-600/20 text-green-400 hover:bg-green-600 hover:text-white transition-colors text-sm flex items-center gap-1.5 flex-shrink-0"
                                    on:click={() =>
                                        handleRestoreFile(file.path)}
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
                    </div>
                {/if}
            </div>
        </div>
    </div>
{/if}
