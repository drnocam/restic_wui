<script>
  import { SvelteToast } from "@zerodevx/svelte-toast";
  import { loading } from "$lib/store";
  import { Circle2 } from "svelte-loading-spinners";

  const options = {
    duration: 4000,
    initial: 1,
    next: 0,
    pausable: false,
    dismissable: true,
    reversed: false,
    intro: { x: 256 },
    theme: {},
    classes: [],
  };
</script>

{#if $loading}
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
    <div
      class="flex flex-col items-center justify-center bg-white p-4 rounded-lg"
    >
      <Circle2 size="90" unit="px" />
      <span class="p-2 font-bold">Loading ...</span>
    </div>
  </div>
{/if}
<SvelteToast {options} />
