<script>
  import { getLibraries } from "../lib/controllers/libraries";
  import HeaderIconLink from "../lib/components/HeaderIconLink.svelte";
  import routes from "./routes.js";
  import LibraryCard from "../lib/components/Card.svelte";
</script>

<div class="container">
  <h1 class="title is-1 inline">Libraries</h1>
  <HeaderIconLink
    icon="fa-solid fa-plus"
    ariaLabel="create library"
    to={routes.create.library}
  />

  <div class="section">
    {#await getLibraries()}
      <span>loading</span>
    {:then libraries}
      <div class="grid">
        {#each libraries as { id, name } (id)}
          <div class="cell">
            <LibraryCard
              {name}
              {id}
              toFunc={routes.libraryFunc}
              toEditFunc={routes.edit.libraryFunc}
            />
          </div>
        {:else}
          <p>No libraries. Create one</p>
        {/each}
      </div>
    {:catch}
      <p>Oops! Could not get libraries</p>
    {/await}
  </div>
</div>
