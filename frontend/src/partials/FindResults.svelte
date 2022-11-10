<script lang="ts"> 
    import { formatBytes , myfetch_json, warning_notify } from "/src/myfuncs";
    import { selected_repo_id } from '/src/store.js';


    let find_data = []
    let search_value = null
    let cb_form ;
    let checked_items = []
    import { SearchInRepo } from "/wailsjs/go/main/App";



    const search_in = () => {
        myfetch_json(SearchInRepo, $selected_repo_id,search_value ).then(r=>{
            if(r){
                find_data = r
            }            
        }
      )
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

<div class="input-group mb-3">
    <input bind:value={search_value} type="text" class="form-control" placeholder="You may find by * symbols" aria-describedby="search-button">
    <button class="btn btn-outline-primary" id="search-button" on:click={search_in}>Search</button>
  </div>

  <button class="btn btn-primary" on:click={extract_files}>Extract Selected Files</button>
<hr>
{#if find_data.length > 0}

{#each find_data as t}

<h6><strong>Snapshot:</strong> {t.snapshot} </h6>
<span class="text-muted">Hits {t.hits}</span>

<form class="table-responsive" bind:this={cb_form}>
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
        <input type="checkbox" value="{tm.path}" name="{`cb`+i}" >  
        </td> 
      </tr>
{/each}
        
    </tbody>
  </table>
</form>
{/each}

{/if}
<style>

        .table-sm>:not(caption)>*>*{
            padding:0;
        }
    
</style>