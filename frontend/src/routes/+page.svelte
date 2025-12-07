<script lang="ts">
    import { api } from "$lib/api";
    import { currentRepoPath, isRepoOpen, loading, error } from "$lib/store";
    import { toast } from "@zerodevx/svelte-toast";

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

    let repoStats: any = null;

    async function loadStats() {
        try {
            repoStats = await api.getRepositoryStats();
        } catch (e) {
            console.error("Failed to load stats", e);
        }
    }

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

            {#if repoStats}
                <div
                    class="p-6 rounded-2xl bg-slate-800/50 border border-slate-700"
                >
                    <h3 class="text-slate-400 text-sm font-medium mb-2">
                        Total Size
                    </h3>
                    <p class="text-2xl font-bold text-white">
                        {(repoStats.total_size / 1024 / 1024 / 1024).toFixed(2)}
                        GB
                    </p>
                </div>
                <div
                    class="p-6 rounded-2xl bg-slate-800/50 border border-slate-700"
                >
                    <h3 class="text-slate-400 text-sm font-medium mb-2">
                        Total Files
                    </h3>
                    <p class="text-2xl font-bold text-white">
                        {repoStats.total_file_count.toLocaleString()}
                    </p>
                </div>

                <div class="col-span-1 md:col-span-3 flex justify-end">
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
