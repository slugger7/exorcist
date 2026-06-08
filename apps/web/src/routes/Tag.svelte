<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import {
    deleteTag,
    getMediaWithParams,
    updateTag,
  } from "../lib/controllers/tags";
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
      const updatedTitle = await updateTag(id, newTitle);

      title = updatedTitle.name;

      const params = new URLSearchParams(location.search);

      navigate(`${routes.tagFunc(id, title)}?${params.toString()}`, {
        replace: true,
      });
    } finally {
      submittingTitle = false;
    }

    return newTitle;
  };

  const handleDelete = async () => {
    deleting = true;

    try {
      await deleteTag(id);

      window.history.back();
    } finally {
      deleting = false;
    }
  };
</script>

{#snippet headerAddons()}
  <DeleteAddons
    context="tag"
    bind:deleteConfirmation
    {deleting}
    {handleDelete}
  />
{/snippet}

<MediaView
  {title}
  route={routes.tagFunc(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  disableTags={true}
  {updateTitle}
  bind:submittingTitle
  {headerAddons}
/>
