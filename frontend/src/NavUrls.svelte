<script>
import {Router, Route, Link} from 'svelte-routing'
import Index from './page/Index.svelte';
    import Settings from './page/Settings.svelte';

</script>

<style>
    /* your styles go here */
</style>
<Router >
     <Route path="/"><Index /></Route>
     <Route path="settings"><Settings /></Route>
  </Router>