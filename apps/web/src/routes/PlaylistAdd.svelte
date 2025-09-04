<script>
  import { getArrayOfStringsSearchParamOrDefault } from "../lib/searchParams";
  import { addMedia, get } from "../lib/controllers/playlists";
  import { onMount } from "svelte";
  import routes from "./routes";
  import { Link } from "svelte-routing";

  let media = $state(getArrayOfStringsSearchParamOrDefault("media", []));
  let playlists = $state([]);
  let loadingPlaylists = $state(false);
  let selectedPlaylists = $state([]);
  let submitting = $state(false);

  const fetchPlaylists = async () => {
    loadingPlaylists = true;
    try {
      playlists = await get();
    } finally {
      loadingPlaylists = false;
    }
  };

  onMount(() => {
    fetchPlaylists();
  });

  const togglePlaylist = (id) => () => {
    const exists = !!selectedPlaylists.find((p) => p === id);
    if (exists) {
      selectedPlaylists = selectedPlaylists.filter((p) => p !== id);
    } else {
      selectedPlaylists.push(id);
    }
  };

  const oncancel = () => {
    history.back();
  };

  const onadd = async () => {
    submitting = true;
    try {
      await addMedia(selectedPlaylists, media);

      history.back();
    } finally {
      submitting = false;
    }
  };
</script>

<div class="container">
  <h1 class="title is-1">Add {media.length} media to playlists</h1>
  {#if playlists && loadingPlaylists == false}
    {#if playlists.length === 0}
      <p>No playlists found to add media to.</p>
      <p>First <Link to={routes.create.playlist}>create</Link> a playlist</p>
    {/if}
    <div class="grid">
      {#each playlists as playlist (playlist.id)}
        <div class="cell">
          <div class="card">
            <header class="card-header">
              <button
                class="card-header-title"
                disabled={submitting}
                onclick={togglePlaylist(playlist.id)}
                >{playlist.name}
              </button>
              <div class="card-header-icon">
                <span class="icon">
                  <i
                    class={`fas ${selectedPlaylists.includes(playlist.id) ? "fa-check" : "fa-plus"}`}
                  ></i>
                </span>
              </div>
            </header>
          </div>
        </div>
      {/each}
    </div>
    <div class="field is-grouped">
      <div class="control">
        <button
          class="button is-primary"
          disabled={submitting || media.length == 0}
          onclick={onadd}>Add</button
        >
      </div>
      <div class="control">
        <button class="button" onclick={oncancel} disabled={submitting}
          >Cancel</button
        >
      </div>
    </div>
  {/if}
</div>
