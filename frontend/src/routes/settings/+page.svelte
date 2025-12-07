<script lang="ts">
    import { api } from "$lib/api";
    import { onMount } from "svelte";
    import { toast } from "@zerodevx/svelte-toast";

    let config = {
        ResticCommand: "",
    };
    let loading = false;

    onMount(async () => {
        try {
            const settings = await api.getSettings();
            if (settings) {
                config = settings;
            }
        } catch (e) {
            console.error(e);
            toast.push("Failed to load settings", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        }
    });

    async function save() {
        loading = true;
        try {
            await api.saveSettings(config);
            toast.push("Settings saved successfully", {
                theme: {
                    "--toastBackground": "#4ADE80",
                    "--toastColor": "black",
                },
            });
        } catch (e) {
            console.error(e);
            toast.push("Failed to save settings", {
                theme: {
                    "--toastBackground": "#F87171",
                    "--toastColor": "white",
                },
            });
        } finally {
            loading = false;
        }
    }
</script>

<div class="max-w-4xl mx-auto p-6">
    <div class="mb-8 flex items-center justify-between">
        <div>
            <h1 class="text-3xl font-bold text-white mb-2">Settings</h1>
            <p class="text-slate-400">Configure global application settings</p>
        </div>
        <a
            href="/"
            class="px-4 py-2 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 transition-colors"
        >
            Back to Dashboard
        </a>
    </div>

    <div class="bg-slate-800/50 border border-slate-700 rounded-2xl p-8">
        <div class="space-y-6">
            <div>
                <label
                    for="restic-cmd"
                    class="block text-sm font-medium text-slate-300 mb-2"
                >
                    Restic Command Path
                </label>
                <div class="relative">
                    <input
                        id="restic-cmd"
                        type="text"
                        bind:value={config.ResticCommand}
                        class="w-full bg-slate-900 border border-slate-700 rounded-xl px-4 py-3 text-white focus:ring-2 focus:ring-blue-500 focus:border-transparent outline-none transition-all"
                        placeholder="restic"
                    />
                    <p class="mt-2 text-sm text-slate-500">
                        Path to the restic executable. Default is "restic"
                        (assumes it's in your PATH).
                    </p>
                </div>
            </div>

            <div class="pt-4 border-t border-slate-700 flex justify-end">
                <button
                    class="px-6 py-3 rounded-xl bg-blue-600 text-white hover:bg-blue-700 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                    on:click={save}
                    disabled={loading}
                >
                    {#if loading}
                        <svg
                            class="animate-spin h-5 w-5 text-white"
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                        >
                            <circle
                                class="opacity-25"
                                cx="12"
                                cy="12"
                                r="10"
                                stroke="currentColor"
                                stroke-width="4"
                            ></circle>
                            <path
                                class="opacity-75"
                                fill="currentColor"
                                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                            ></path>
                        </svg>
                        Saving...
                    {:else}
                        Save Changes
                    {/if}
                </button>
            </div>
        </div>
    </div>
</div>
