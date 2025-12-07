<script>
    export let isOpen = false;
    export let toggle = () => {
        isOpen = !isOpen;
    };
    export let size = "md"; // sm, md, lg, xl, fullscreen
</script>

{#if isOpen}
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4"
        on:click|self={toggle}
    >
        <div
            class="bg-white dark:bg-gray-800 rounded-lg shadow-xl overflow-hidden max-h-full flex flex-col
      {size === 'fullscreen'
                ? 'w-full h-full'
                : size === 'lg'
                  ? 'w-full max-w-4xl'
                  : 'w-full max-w-md'}"
        >
            {#if $$slots.header}
                <div
                    class="p-4 border-b dark:border-gray-700 flex justify-between items-center"
                >
                    <div class="text-lg font-semibold dark:text-white">
                        <slot name="header" />
                    </div>
                    <button
                        on:click={toggle}
                        class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 text-2xl leading-none"
                        >&times;</button
                    >
                </div>
            {/if}

            <div class="p-4 overflow-y-auto flex-1 dark:text-gray-300">
                <slot />
            </div>

            {#if $$slots.footer}
                <div
                    class="p-4 border-t dark:border-gray-700 bg-gray-50 dark:bg-gray-900 flex justify-end gap-2"
                >
                    <slot name="footer" />
                </div>
            {/if}
        </div>
    </div>
{/if}
