<script>
  import { onMount } from "svelte";

  /** @type {{id: string}}*/
  let { id } = $props();

  $effect(() => {
    console.log(history.state);
  });
</script>

<h1>{id}</h1>
