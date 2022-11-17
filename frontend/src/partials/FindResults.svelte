<script lang="ts"> 
    import { formatBytes , myfetch_json, warning_notify,htmlEntities } from "/src/myfuncs";
    import { selected_repo_id } from '/src/store.js';
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


    let modal_info = ""; 
    let modal_open_sm = false;
    let find_data = []
    let search_value = null
    let cb_form ;
    let checked_items = []
    import { SearchInRepo,RestoreFilesInSnapshots } from "/wailsjs/go/main/App";



    const search_in = () => {
        myfetch_json(SearchInRepo, $selected_repo_id,search_value ).then(r=>{
            if(r){
                find_data = r
            }            
        }
      )
    }

    const toggleModalSm = () => {
    /* 
    * if modal is open then it will be closed so reset modal_info;
    */
    if(modal_open_sm ) {
        modal_info = "";
    }
    modal_open_sm = !modal_open_sm;

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
      let values ={}
      for (let i = 0; i < el.length; i++) {
        if(cb_form.elements[i].checked){
          if(values[cb_form.elements[i].dataset.snapshot] === undefined){
            values[cb_form.elements[i].dataset.snapshot] = []
          }
          values[cb_form.elements[i].dataset.snapshot].push(cb_form.elements[i].value)
        }
    } 
    
    if(Object.keys(values).length>0) {
      Object.keys(values).forEach(key => {
/* 
TODO: dont repeat to ask folder, so go should handle
*/
        myfetch_json(RestoreFilesInSnapshots, $selected_repo_id,key,JSON.stringify(values[key]) ).then(r=>{
          modal_open_sm = true
          modal_info =  String(r).replace(/(?:\r\n|\r|\n)/g,"<br>")        
          modal_info =  htmlEntities(modal_info)        
        })

      });

        
        
      }


    }
    
</script>
<Modal isOpen={modal_open_sm} toggle={toggleModalSm} >
  <ModalBody>
      {@html modal_info}
  </ModalBody>
</Modal>
<div class="input-group mb-3">
    <input bind:value={search_value} type="text" class="form-control" placeholder="You may find by * symbols" aria-describedby="search-button">
    <button class="btn btn-outline-primary" id="search-button" on:click={search_in}>Search</button>
  </div>

  <button class="btn btn-primary" on:click={extract_files}>Extract Selected Files</button>
<hr>
{#if find_data.length > 0}
<form class="table-responsive" bind:this={cb_form}>
{#each find_data as t}

<h6><strong>Snapshot:</strong> {t.snapshot} </h6>
<span class="text-muted">Hits {t.hits}</span>


<table class="table table-sm table-striped">
    <thead>
      <tr>
        <th scope="col">#</th>
        <th scope="col">User</th>
        <th scope="col">Path</th>
        <th scope="col">Size</th>
        <th scope="col">Created Time</th>
        <th scope="col">Select</th>
      </tr>
    </thead>
    <tbody>
{#each t.matches as tm,i}
      <tr>
        <th scope="row">{i+1}</th>
        <td>{tm.user}</td>
        <td>{tm.path}</td>
        <td>{formatBytes(tm.size)}</td> 
        <td>{(new Date(Date.parse(tm.ctime))).toLocaleString("tr-TR")}</td> 
        <td>
        <input type="checkbox" data-snapshot={t.snapshot} value="{tm.path}" name="{`cb`+i}" >  
        </td> 
      </tr>
{/each}
        
    </tbody>
  </table> 
{/each}
</form>
{/if}

<style>

        .table-sm>:not(caption)>*>*{
            padding:0;
        }
    
</style>