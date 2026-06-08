<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import {
    deletePlaylist,
    getMediaWithParams,
    updatePlaylist,
  } from "../lib/controllers/playlists";
  import { navigate } from "svelte-routing";

  let { id, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);
  let deleting = $state(false);
  let deleteConfirmation = $state(false);

  /** @param {string} newTitle*/
  const updateTitle = async (newTitle) => {
    submittingTitle = true;
    try {
      const updatedTitle = await updatePlaylist(id, newTitle);

      title = updatedTitle.name;

      const params = new URLSearchParams(location.search);

      navigate(`${routes.playlistFn(id, title)}?${params.toString()}`, {
        replace: true,
      });
    } finally {
      submittingTitle = false;
    }

    return title;
  };

  const handleDelete = async () => {
    deleting = true;
    try {
      await deletePlaylist(id);

      window.history.back();
    } finally {
      deleting = false;
    }
  };
</script>

{#snippet headerAddons()}
  {#if !deleteConfirmation}
    <p class="control">
      <button
        class="button"
        aria-label="delete playlist"
        onclick={() => (deleteConfirmation = true)}
        ><span class="icon"><i class="fas fa-trash"></i></span></button
      >
    </p>
  {/if}
  {#if deleteConfirmation}
    <p class="control">
      <button
        class="button"
        aria-label="cancel delete"
        onclick={() => (deleteConfirmation = false)}
        ><span class="icon"><i class="fas fa-xmark"></i></span></button
      >
    </p>
    <p class="control">
      <button
        class="button is-danger"
        aria-label="confirm delete"
        onclick={handleDelete}
        ><span class="icon"><i class="fas fa-check"></i></span></button
      >
    </p>
  {/if}
{/snippet}

<MediaView
  {title}
  route={routes.playlistFn(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  {updateTitle}
  bind:submittingTitle
  {headerAddons}
/>
