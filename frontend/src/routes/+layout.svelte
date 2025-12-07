<script>
    import "../app.css";
    import Sidebar from "$lib/components/Sidebar.svelte";
    import Header from "$lib/components/Header.svelte";
    import { SvelteToast } from "@zerodevx/svelte-toast";
</script>

<div
    class="flex h-screen bg-slate-950 text-slate-100 overflow-hidden font-sans antialiased selection:bg-blue-500/30"
>
    <Sidebar />

    <main class="flex-1 flex flex-col min-w-0 overflow-hidden relative">
        <div
            class="absolute inset-0 bg-center [mask-image:linear-gradient(180deg,white,rgba(255,255,255,0))] opacity-20 pointer-events-none"
        ></div>

        <Header />

        <div class="flex-1 overflow-y-auto p-8 scroll-smooth">
            <div class="max-w-6xl mx-auto">
                <slot />
            </div>
        </div>
    </main>

    <SvelteToast options={{ duration: 4000, intro: { y: 100 } }} />
</div>
