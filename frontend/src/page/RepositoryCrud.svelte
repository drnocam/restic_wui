<script lang="ts">
    import {
    Button,
    ButtonGroup,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader
  } from 'sveltestrap';  
    import { ChooseRepository,AddUpdateRepository,DeleteRepositorySetting,DeleteRepositoryFromDisk } from "/wailsjs/go/main/App";
    import { selected_repo_id,repositories } from '/src/store.js';
    import { myfetch_json } from "/src/myfuncs"; 



    export let parent_snapshot_function = null;
    /* 
    0 : new record
    1 : update_id
    */
    let islem : number = 0;

    export let repo_add_modal = false;
    export let repo_form = {
        Name:null,Path:null,Password:null,Args:null
    }

    function delete_from_settings() :void { 
        myfetch_json(DeleteRepositorySetting,$selected_repo_id) .then(r=>{
            repositories.set(r.names);
            repo_toggle()
    }
      )
  }
/* 
TODO: delete from disk
*/
  function delete_from_disk() :void { 
        myfetch_json(DeleteRepositoryFromDisk,$selected_repo_id) .then(r=>{
            repositories.set(r.names);
            repo_toggle()
    }
      )
  }


    function choose_repository(): void {
        ChooseRepository().then(
            (r) => {
                repo_form["Path"] = r;
            }
        );
    }

    const repo_toggle = ()=>{
      repo_add_modal = !repo_add_modal;
    }

    function save_repo() { 
        myfetch_json(AddUpdateRepository, $selected_repo_id,JSON.stringify(repo_form) ).then(r=>{
            repositories.set(r.names);
            repo_toggle() 
            parent_snapshot_function()
        }
      )
    }

    


 
</script> 
<div class="btn btn-primary">Backup</div>
<div class="btn btn-primary">Check Backup</div>
<div class="btn btn-primary">Show Difference</div>
<div class="btn btn-primary">Forget</div>
<div class="btn btn-primary">List In Snapshot</div>
<div class="btn btn-primary">Find In Snapshot</div>
<div class="btn btn-primary">Stats</div>
<div class="btn btn-primary">Unlock</div>
<div class="btn btn-primary">Tag</div>
<Modal isOpen={repo_add_modal} toggle={repo_toggle} size="lg">
    <ModalHeader {repo_toggle}>
      {#if islem == 0}
      Update / Delete Repository
      {/if}
    </ModalHeader>
    <ModalBody>
<form class="row g-3">
    <div class="col-12">
        <label for="i4" class="form-label">Alias Name</label>
        <input type="text"
            class="form-control"
            id="i4"
            placeholder='like: Home Folder'
            bind:value={repo_form["Name"]}
        /> 
    </div>

    <div class="col-12">
        <label for="i1" class="form-label">Choose Repository</label>

        <div class="input-group mb-3">
            <input
                type="text"
                class="form-control"
                id="i1"
                bind:value={repo_form["Path"]}
            />
            <button class="btn btn-primary" on:click|preventDefault={choose_repository}
                >Choose Repository</button
            > 
        </div>
        <span class="text-muted"
            >If you want to input like rclone you can write.</span
        >
    </div>
    <div class="col-12">
        <label for="i2" class="form-label">Password</label>
        <input
            type="password"
            class="form-control"
            id="i2"
            bind:value={repo_form["Password"]}
        />
    </div>
    <div class="col-12">
        <label for="i3" class="form-label">Other Args</label>
        <textarea
            class="form-control"
            id="i3"
            placeholder='argclone="--tpslimit=10" etc'
            bind:value={repo_form["Args"]}
        />
        <span class="text-muted">You may put inputs as new line.</span>
    </div>
    <div class="col-12 text-end">
        {#if $selected_repo_id == -1}
        <div class="btn btn-success" on:click={save_repo}>Save Repository Setting</div>
        {:else}
        <div class="btn btn-success" on:click={save_repo}>Update Repository Setting</div>
        {/if}
    </div>
    <div class="col-12">
        <div class="btn btn-danger" on:click={delete_from_settings}>Delete Setting</div>
        <div class="btn btn-danger" on:click={delete_from_disk}>Delete Repository From Disk ( Dangerous )</div> 
    </div>
</form>
    </ModalBody> 
  </Modal> 
<style>
    /* your styles go here */
</style>
