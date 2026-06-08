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
  import DeleteAddons from "../lib/components/DeleteAddons.svelte";

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
  <DeleteAddons
    context="playlist"
    bind:deleteConfirmation
    {deleting}
    {handleDelete}
  />
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
