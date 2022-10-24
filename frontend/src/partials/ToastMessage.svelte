<script>
    import { SvelteToast } from '@zerodevx/svelte-toast'
    import { loading } from '/src/store';
	import { Circle2 } from 'svelte-loading-spinners';

    const options = {
  duration: 4000,       // duration of progress bar tween to the `next` value
  initial: 1,           // initial progress bar value
  next: 0,              // next progress value
  pausable: false,      // pause progress bar tween on mouse hover
  dismissable: true,    // allow dismiss with close button
  reversed: false,      // insert new toast to bottom of stack
  intro: { x: 256 },    // toast intro fly animation settings
  theme: {},            // css var overrides
  classes: []           // user-defined classes
}
  </script>
  
  <style>
    #backdrp{
        position: fixed;
        width: 100%;
        height: 100%;
        background-color: rgba(0,0,0,0.5);
        z-index: 1;
        display: flex;
        align-items: center;
        justify-content: center;
    }
  </style>

  {#if $loading}
  <div id="backdrp">
    <div>
    <Circle2 size="90" unit="px"   />
    <span class="p-2">Loading ...</span>
    </div>
</div>
  {/if}
    <SvelteToast  {options} />
