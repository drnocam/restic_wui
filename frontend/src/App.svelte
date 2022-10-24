<script lang="ts">

import NavUrls from './NavUrls.svelte';
    import Header from './partials/Header.svelte';
     import { onMount, text } from 'svelte/internal';
    import {CmdCheck } from '../wailsjs/go/main/App.js'
    import { repositories } from './store';
    import ToastMessage from './partials/ToastMessage.svelte';
    import { myfetch_json } from './myfuncs';

    onMount(()=>{
      myfetch_json(CmdCheck,null).then(r=>{

        repositories.set(r.names);
      
    })
  })
</script>
<ToastMessage />
<Header />
<div class="container-fluid">
  <NavUrls />
</div>

<style>

</style>
