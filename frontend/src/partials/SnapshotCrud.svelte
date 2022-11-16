<script lang="ts">
    import {
        Accordion,
    AccordionItem,
    Button,
    ButtonGroup,
    Modal,
    ModalBody,
    ModalFooter,
    ModalHeader
  } from 'sveltestrap';  
    import { ListFilesInSnapshots,RestoreRepo } from "/wailsjs/go/main/App";
    import { selected_repo_id  } from '/src/store.js';
    import { myfetch_json ,formatBytes, warning_notify } from "/src/myfuncs";   


/* == Variables == */
    let modal_info = "";
    let modal_open = false;
    let snapshot_files_list = []
    let cb_form ;
    export let snapshot_info = {'id':null};
    
    const restore_snapshot = () => {
        myfetch_json(RestoreRepo, $selected_repo_id,snapshot_info.id ).then(r=>{
            if(r){
                toggleModal()
                modal_info =  String(r).replace(/(?:\r\n|\r|\n)/g,"<br>") 

            }            
        }
      )
    }

    const list_files_in_snapshot = () => {
        myfetch_json(ListFilesInSnapshots, $selected_repo_id,snapshot_info.id ).then(r=>{
            if(r){
                toggleModal()
                snapshot_files_list = r.slice(1) /* except first index */ 
            }            
        }
      )
    }

const toggleModal = () => {
    /* 
    * if modal is open then it will be closed so reset modal_info;
    */
    if(modal_open ) {
        modal_info = "";
    }
    modal_open = !modal_open;

}

const extract_files = () => {
      if(!cb_form) {
        warning_notify("No Results Yet")
        return ;
      }
      let el = cb_form.elements
      if (el.length == 0) {
        warning_notify("Please Select Items First")
        return;
      }
      for (let i = 0; i < el.length; i++) {
        if(cb_form.elements[i].checked)
          console.log(cb_form.elements[i].value)
    }
}



</script>

<style>
    /* your styles go here */
</style>
{#if snapshot_info.id!=0}
<div class="btn btn-primary" on:click={()=>{
    restore_snapshot()
}}>Restore Snapshot Id:
    <span class="">{snapshot_info["short_id"]}</span> 
    </div>
<div class="btn btn-primary" on:click={()=>{
    list_files_in_snapshot()
}}>List Files In Snapshot Id:
    <span class="">{snapshot_info["short_id"]}</span> 
    </div>
{#each Object.keys(snapshot_info) as sp}
    <div><strong>{sp}</strong> : {snapshot_info[sp]}</div>
      {/each}

      {/if}
<Modal isOpen={modal_open} toggle={toggleModal} size="lg">
    <ModalBody>
        {#if snapshot_files_list.length > 0}
        <button class="btn btn-primary" on:click={extract_files}>Extract Selected Files</button>
        
        <hr>

<form class="table-responsive" bind:this={cb_form}>
<table class="table table-sm table-striped">
    <thead>
      <tr>
        <th scope="col">#</th>
        <th scope="col">Name</th>
        <th scope="col" >Path</th>
        <th scope="col">Type</th>
        <th scope="col">Size</th>
        <th scope="col">Created Time</th>
        <th scope="col">Select</th>
      </tr>
    </thead>
    <tbody>
{#each snapshot_files_list as tm,i}
      <tr>
        <th scope="row">{i+1}</th>
        <td>{tm.name}</td>
        <td style="max-width: 100px;overflow:scroll">{tm.path}</td>
        <td>{tm.type}</td>
        <td>{formatBytes(tm.size)}</td> 
        <td>{(new Date(Date.parse(tm.ctime))).toLocaleString("tr-TR")}</td> 
        <td>
        <input type="checkbox" value="{tm.path}" name="{`cb`+i}" >  
        </td> 
      </tr>
{/each}
        
    </tbody>
  </table>
</form>

{/if}
    </ModalBody>
</Modal>