<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import {
    getMediaWithParams,
    updatePlaylist,
  } from "../lib/controllers/playlists";
  import { navigate } from "svelte-routing";

  let { id, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);

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
</script>

<MediaView
  {title}
  route={routes.playlistFn(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  {updateTitle}
  bind:submittingTitle
/>
