<script lang="ts">

  import { text } from 'svelte/internal';
import {ChooseRepository,GetSnapshots} from '../wailsjs/go/main/App.js'

  let selected_repository: string = ""
  let name: string
  let repo_password: string
  let text_result: string = ""



  function choose_repository():void {
    ChooseRepository().then(r=>selected_repository = 'Selected repository: ' + r)
  }

  function get_snapshots() :void {
    GetSnapshots(repo_password).then(r=>{
      let onceki = text_result;
      text_result = r+ onceki}
      )
  }

</script>

<main>
    <input type="password" placeholder="Write Repository Password" bind:value={repo_password} />
    <button class="btn" on:click={choose_repository}>Choose Repository</button>
    <div class="result" id="result">{selected_repository}</div>
    {#if selected_repository && repo_password}
    <button class="btn" on:click={get_snapshots}>Get Snapshots</button>
    {:else}
    <div>Please write password first to get snapshots.</div>
    {/if}
    <div style="background-color: white;min-height:200px;color:black">{text_result}</div>
</main>

<style>

</style>
