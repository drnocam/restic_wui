<script lang="ts">

import NavUrls from './NavUrls.svelte';
    import Header from './partials/Header.svelte';
     import { onMount, text } from 'svelte/internal';
    import {CmdCheck } from '../wailsjs/go/main/App.js'
    import { repositories } from './store';
    import ToastMessage from './partials/ToastMessage.svelte';
    import { myfetch_json } from './myfuncs';
    import { EventsOn } from '../wailsjs/runtime'

    onMount(()=>{
      myfetch_json(CmdCheck,null).then(r=>{

        repositories.set(r.names);
      
    })

    console.log(EventsOn);
    EventsOn("event_test", x => {
        if (x) {
            console.log(x);
        } 
    });

  })
</script>
<ToastMessage />
<Header />
<div class="container-fluid">
  <NavUrls />
</div>

<style>

</style>
