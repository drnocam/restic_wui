<script lang="ts">
    import { page } from "$app/state";
    import { isRepoOpen } from "$lib/store";
</script>

<aside class="w-64 bg-slate-900 border-r border-slate-800 flex flex-col h-full">
    <div class="p-6 flex items-center gap-3">
        <div
            class="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center shadow-lg shadow-blue-500/20"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5 text-white"
                viewBox="0 0 20 20"
                fill="currentColor"
            >
                <path
                    fill-rule="evenodd"
                    d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H8z"
                    clip-rule="evenodd"
                />
            </svg>
        </div>
        <h1
            class="text-xl font-bold bg-clip-text text-transparent bg-gradient-to-r from-blue-400 to-cyan-400"
        >
            Restic UI
        </h1>
    </div>

    <nav class="flex-1 px-4 space-y-2 mt-4">
        <a
            href="/"
            class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 {page
                .url.pathname === '/'
                ? 'bg-blue-600/10 text-blue-400'
                : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
        >
            <svg
                xmlns="http://www.w3.org/2000/svg"
                class="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
            >
                <path
                    d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"
                />
            </svg>
            <span class="font-medium">Dashboard</span>
        </a>

        {#if $isRepoOpen}
            <a
                href="/snapshots"
                class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 {page.url.pathname.startsWith(
                    '/snapshots',
                )
                    ? 'bg-blue-600/10 text-blue-400'
                    : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                >
                    <path
                        fill-rule="evenodd"
                        d="M4 2a1 1 0 011 1v2.101a7.002 7.002 0 0111.601 2.566 1 1 0 11-1.885.666A5.002 5.002 0 005.999 7H9a1 1 0 010 2H4a1 1 0 01-1-1V3a1 1 0 011-1zm.008 9.057a1 1 0 011.276.61A5.002 5.002 0 0014.001 13H11a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0v-2.101a7.002 7.002 0 01-11.601-2.566 1 1 0 01.61-1.276z"
                        clip-rule="evenodd"
                    />
                </svg>
                <span class="font-medium">Snapshots</span>
            </a>

            <a
                href="/settings"
                class="flex items-center gap-3 px-4 py-3 rounded-xl transition-all duration-200 {page.url.pathname.startsWith(
                    '/settings',
                )
                    ? 'bg-blue-600/10 text-blue-400'
                    : 'text-slate-400 hover:bg-slate-800 hover:text-slate-200'}"
            >
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                >
                    <path
                        fill-rule="evenodd"
                        d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z"
                        clip-rule="evenodd"
                    />
                </svg>
                <span class="font-medium">Settings</span>
            </a>
        {/if}
    </nav>

    <div class="p-4 border-t border-slate-800">
        <div class="bg-slate-800/50 rounded-xl p-4">
            <p class="text-xs text-slate-500 mb-1">Status</p>
            <div class="flex items-center gap-2">
                <div
                    class="w-2 h-2 rounded-full {$isRepoOpen
                        ? 'bg-green-500 shadow-[0_0_8px_rgba(34,197,94,0.6)]'
                        : 'bg-slate-600'}"
                ></div>
                <span class="text-sm font-medium text-slate-300"
                    >{$isRepoOpen ? "Connected" : "Disconnected"}</span
                >
            </div>
        </div>
    </div>
</aside>
