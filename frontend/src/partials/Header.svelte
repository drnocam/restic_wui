<script type="ts">
    import { dark_mode,repositories,selected_repo_id,snapshots } from "/src/store";
    import { Router,Link } from "svelte-routing";
    import {
    Button,
    ButtonGroup,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader
  } from 'sveltestrap';  
    import RepositoryCrud from "/src/page/RepositoryCrud.svelte";
    import { GetSnapshots } from "/wailsjs/go/main/App"; 
    import { myfetch_json } from "/src/myfuncs";

    let mode = 'lightbulb-off-fill';
    let dark = 'light';
    let settings_modal = false;
    let repo_add_modal = false;
    let islem = 0; // 0 , 1 add , update vs repo


    dark_mode.subscribe(m=>{
    if(m==0) {
        dark = 'light'
        mode = 'lightbulb'
        document.getElementsByTagName("body")[0].style = 'background:white;color:var(--bs-dark)';
    } else {
        dark = 'dark';
        mode = 'lightbulb-off-fill'
        document.getElementsByTagName("body")[0].style = 'background:var(--bs-dark);color:var(--bs-light)';

    }
})


    const change_mode = ()=>{
        if($dark_mode == 0) {
            dark_mode.set(1)
        } else {
            dark_mode.set(0)
        }
    }
    function get_snapshots() :void {
          myfetch_json(GetSnapshots,$selected_repo_id).then(r=>{
            snapshots.set(r) 
          }
            )
        }

    const settings_toggle = ()=>{
      settings_modal = !settings_modal;
    }

    const repo_toggle = ()=>{
      repo_add_modal = !repo_add_modal;
    }

    const repo_info = (e)=>{
      selected_repo_id.set(parseInt(e.currentTarget.value))

get_snapshots()
      
    }

</script>
<Router>
<nav class="navbar navbar-expand-lg bg-{dark}">
    <div class="container-fluid">
      <div class="navbar-brand"><Link to="/">Restic Command Line Parser UI</Link></div>

      {#if Object.keys($repositories).length > 0}
      <select class="form-select" aria-label="Repositories" on:change={repo_info}>
        <option selected>Choose Repository</option>
        {#each  Object.keys($repositories) as rp}
        <option value="{rp}">{$repositories[rp]}</option>
        {/each} 
        
      </select>{:else }
      <div class="alert alert-warning">Add Repository</div>
    {/if}
<div style="white-space:nowrap">
{#if $selected_repo_id != -1}
      <div class="btn btn-{dark}" title="Update Repository" on:click={repo_toggle}><i class="bi-pen"></i></div>
{/if}
      <div class="btn btn-{dark}" title="Add Repository" on:click={repo_toggle}><i class="bi-plus"></i></div>
      <div class="btn btn-{dark}" title="Settings" on:click={settings_toggle}><i class="bi-gear"></i></div>
        <div class="btn btn-{dark}" title="Dark/Light Mode"  on:click={change_mode}><i class="bi-{mode}"></i></div>
      </div>
    </div>
  </nav>
</Router>
<Modal isOpen={settings_modal}  toggle={settings_toggle} fullscreen>
  <ModalHeader {settings_toggle}>Modal title</ModalHeader>
  <ModalBody>
    Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
    tempor incididunt ut labore et dolore magna aliqua.
    <div class="btn btn-{dark}"  on:click={repo_toggle}>asdfasdfas</div>

  </ModalBody>
  <ModalFooter>
    <Button color="primary" on:click={settings_toggle}>Do Something</Button>
    <Button color="secondary" on:click={settings_toggle}>Cancel</Button>
  </ModalFooter>
</Modal>
<Modal isOpen={repo_add_modal} toggle={repo_toggle} size="lg">
  <ModalHeader {repo_toggle}>
    {#if islem == 0}
    Update / Delete Repository
    {/if}
  </ModalHeader>
  <ModalBody>
    {#if islem == 0}
    <RepositoryCrud />
    {/if}
  </ModalBody> 
</Modal>