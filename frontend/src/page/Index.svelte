<script lang="ts">
    import { snapshots } from '/src/store.js';
     import { onMount, text } from 'svelte/internal';
import {GetSnapshots, ReadWriteSettings} from '../../wailsjs/go/main/App.js'


let cmd_ready : boolean = false;

  let name: string
  let text_result: string = ""



</script>

<style>
    
</style>
<button class="btn btn-primary">Create New Repository</button>
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


   