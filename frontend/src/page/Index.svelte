<script lang="ts">
    import { selected_repo_id } from '/src/store.js';
     import { onMount, text } from 'svelte/internal';
import {CmdCheck, ChooseRepository,GetSnapshots, ReadWriteSettings} from '../../wailsjs/go/main/App.js'
    import RepositoryCrud from './RepositoryCrud.svelte';


let cmd_ready : boolean = false;

  let name: string
  let text_result: string = ""





  function get_snapshots() :void {
    GetSnapshots($selected_repo_id).then(r=>{
      let onceki = text_result;
      text_result = r+ onceki}
      )
  }

function read_write_settings() : void {

  ReadWriteSettings().then(r=>{
    console.log(r)
  })

}


</script>

<style>
    
</style>
<div class="btn btn-primary" on:click={read_write_settings}>Read Write Settings</div>
 
    <button class="btn btn-primary" on:click={get_snapshots}>Get Snapshots</button>

    <div style="min-height:200px ">{text_result}</div>


    <RepositoryCrud />