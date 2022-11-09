<script lang="ts">
    import { snapshots } from '/src/store.js';




</script>

<style>
    
</style>
    <div style="min-height:200px ">
      {#if Array.isArray($snapshots) }
      {#each $snapshots as s_arr}
      {#each Object.keys(s_arr) as sp}
    <div><strong>{sp}</strong> : {s_arr[sp]}</div>
      {/each}
      <hr>
      {/each}
      {/if}
    </div>


   