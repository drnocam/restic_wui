<script lang="ts">
     import { onMount, text } from 'svelte/internal';
import {CmdCheck, ChooseRepository,GetSnapshots, ReadWriteSettings} from '../../wailsjs/go/main/App.js'


let cmd_ready : boolean = false;
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

function read_write_settings() : void {

  ReadWriteSettings().then(r=>{
    console.log(r)
  })

}

  onMount(()=>{
    CmdCheck().then(r=>JSON.parse(r)).then(r=>{
      if("data" in r){
        cmd_ready = r.data== "1" ? true : false ; 
      }       
    })
  })
</script>

<style>
    
</style>
<div class="btn btn-primary" on:click={read_write_settings}>Read Write Settings</div>
<i class="bi-alarm"></i> 
  <div></div>
    <input type="password" placeholder="Write Repository Password" bind:value={repo_password} />
    <button class="btn btn-primary" on:click={choose_repository}>Choose Repository</button>
    <div class="result" id="result">{selected_repository}</div>
    {#if selected_repository && repo_password}
    <button class="btn btn-primary" on:click={get_snapshots}>Get Snapshots</button>
    {:else}
    <div>Please write password first to get snapshots.</div>
    {/if}
    <div style="background-color: white;min-height:200px;color:black">{text_result}</div>