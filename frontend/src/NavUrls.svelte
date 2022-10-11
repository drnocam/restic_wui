<script>
import {Router, Route, Link} from 'svelte-routing'
import Index from './page/Index.svelte';
    import Settings from './page/Settings.svelte';

</script>

<style>
    /* your styles go here */
</style>
<Router >
    <nav>
      <Link to="/">Home</Link>
      <Link to="settings">Settings</Link>
    </nav>
    <div>
     <Route path="/"><Index /></Route>
     <Route path="settings"><Settings /></Route>
    </div>
  </Router>