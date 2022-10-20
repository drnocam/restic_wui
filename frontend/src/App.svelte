<script lang="ts">

import NavUrls from './NavUrls.svelte';
    import Header from './partials/Header.svelte';
     import { onMount, text } from 'svelte/internal';
    import {CmdCheck } from '../wailsjs/go/main/App.js'
    import { repositories } from './store';

    onMount(()=>{
    CmdCheck().then(r=>JSON.parse(r)).then(r=>{

        repositories.set(r.data.names);
      
    })
  })
</script>
<Header />
<div class="container-fluid">
  <NavUrls />
</div>

<style>

</style>
