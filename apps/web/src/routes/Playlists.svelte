<script>
  import HeaderIconLink from "../lib/components/HeaderIconLink.svelte";
  import routes from "./routes.js";
  import Card from "../lib/components/Card.svelte";
  import { get } from "../lib/controllers/playlists";
</script>

<div class="container">
  <h1 class="title is-1 inline">Playlists</h1>
  <HeaderIconLink
    icon="fa-solid fa-plus"
    ariaLabel="create playlist"
    to={routes.create.playlist}
  />

  <div class="section">
    {#await get()}
      <span>loading</span>
    {:then playlists}
      <div class="grid">
        {#each playlists as { id, name } (id)}
          <div class="cell">
            <Card {name} {id} toFunc={routes.playlistFn} />
          </div>
        {:else}
          <p>No playlists. Create one</p>
        {/each}
      </div>
    {:catch}
      <p>Oops! Could not get playlists</p>
    {/await}
  </div>
</div>
