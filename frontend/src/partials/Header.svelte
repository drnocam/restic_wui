<script type="ts">
    import { dark_mode } from "/src/store";
    import { Router,Link } from "svelte-routing";

    let mode = 'lightbulb-off-fill';
    let dark = 'light';

    dark_mode.subscribe(m=>{
    if(m==0) {
        dark = 'light'
        mode = 'lightbulb'
        document.getElementsByTagName("body")[0].style = 'background:white;color:var(--bs-light)';
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
    

</script>
<Router>
<nav class="navbar navbar-expand-lg bg-{dark}">
    <div class="container-fluid">
      <div class="navbar-brand"><Link to="/">Restic Comman Line Parser UI</Link></div>
      <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
        <span class="navbar-toggler-icon"></span>
      </button>
      <div class="collapse navbar-collapse" id="navbarSupportedContent">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item"> </li>
        </ul> 

        <div class="btn btn-{dark}" title="Settings"><Link to="/settings"><i class="bi-gear"></i></Link></div>
        <div class="btn btn-{dark}" title="Dark/Light Mode" on:click={change_mode}><i class="bi-{mode}"></i></div>
      </div>
    </div>
  </nav>
</Router>