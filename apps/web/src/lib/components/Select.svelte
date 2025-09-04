<script>
  let { options, value = $bindable() } = $props();
</script>

<div class="select">
  <select bind:value>
    {#each options as o}
      <option value={o.value}>{o.text}</option>
    {/each}
  </select>
</div>
