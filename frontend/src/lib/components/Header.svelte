<script lang="ts">
  import { currentRepoPath } from "$lib/store";
</script>

<header
  class="h-16 border-b border-slate-800 bg-slate-900/50 backdrop-blur-xl flex items-center justify-between px-8 sticky top-0 z-10"
>
  <div class="flex items-center gap-4">
    {#if $currentRepoPath}
      <div class="flex flex-col">
        <span
          class="text-xs text-slate-500 font-medium uppercase tracking-wider"
          >Current Repository</span
        >
        <span
          class="text-sm text-slate-200 font-mono bg-slate-800/50 px-2 py-0.5 rounded border border-slate-700/50"
          >{$currentRepoPath}</span
        >
      </div>
    {/if}
  </div>

  <div class="flex items-center gap-4">
    <!-- Add user profile or other actions here -->
  </div>
</header>
