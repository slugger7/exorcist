<script>
  import HeaderIconLink from "../../lib/components/HeaderIconLink.svelte";
  import LibraryPathCard from "../../lib/components/LibraryPathCard.svelte";
  import { getPaths } from "../../lib/controllers/libraries";
  import routes from "../routes";
  /** @type {{id:string}}*/
  let { id: libraryId } = $props();
  /** @type {{libraryName: string}}*/
  let { libraryName } = history.state;
</script>

<div class="container">
  <h1 class="title is-1 inline">Library {libraryName}</h1>
  <HeaderIconLink
    icon="fa-solid fa-plus"
    ariaLabel="create library"
    to={routes.create.libraryPathFunc(libraryId, "")}
    state={{ libraryName }}
  />

  <div class="section">
    {#await getPaths(libraryId)}
      <span>loading</span>
    {:then libraryPaths}
      <div class="grid">
        {#each libraryPaths as { path, id } (id)}
          <div class="cell">
            <LibraryPathCard {path} {id} />
          </div>
        {:else}
          <p>No library paths in library</p>
        {/each}
      </div>
    {:catch}
      <span>something went wrong</span>
    {/await}
  </div>
</div>
