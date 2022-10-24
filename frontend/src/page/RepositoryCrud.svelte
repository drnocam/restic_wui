<script lang="ts">
    import { ChooseRepository,AddUpdateRepository,DeleteRepositorySettings } from "/wailsjs/go/main/App";
    import { selected_repo_id,repositories } from '/src/store.js';
    import { myfetch_json } from "/src/myfuncs"; 




    let selected_repository: string = ""; 
    /* 
    0 : new record
    1 : update_id
    */
    let islem : number = 0;
    let repo_form = {
        name:null,path:null,password:null,args:null
    }

    function delete_from_repo() :void { 
        myfetch_json(DeleteRepositorySettings,$selected_repo_id) .then(r=>{
            repositories.set(r.names);
    }
      )
  }


    function choose_repository(): void {
        ChooseRepository().then(
            (r) => {
                repo_form["path"] = r;
            }
        );
    }

    function save_repo() { 
        myfetch_json(AddUpdateRepository,islem,JSON.stringify(repo_form) ).then(r=>{
            repositories.set(r.names);
    }
      )
    }
 
</script>

<form class="row g-3">
    <div class="col-12">
        <label for="i4" class="form-label">Alias Name</label>
        <input type="text"
            class="form-control"
            id="i4"
            placeholder='like: Home Folder'
            bind:value={repo_form["name"]}
        /> 
    </div>

    <div class="col-12">
        <label for="i1" class="form-label">Choose Repository</label>

        <div class="input-group mb-3">
            <input
                type="text"
                class="form-control"
                id="i1"
                bind:value={repo_form["path"]}
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
            bind:value={repo_form["password"]}
        />
    </div>
    <div class="col-12">
        <label for="i3" class="form-label">Other Args</label>
        <textarea
            class="form-control"
            id="i3"
            placeholder='argclone="--tpslimit=10" etc'
            bind:value={repo_form["args"]}
        />
        <span class="text-muted">You may put inputs as new line.</span>
    </div>
    <div class="col-12">
        <div class="btn btn-danger" on:click={delete_from_repo}>Delete Repository</div>
        <div class="btn btn-success" on:click={save_repo}>SAVE Repository</div>
    </div>
</form>

<style>
    /* your styles go here */
</style>
