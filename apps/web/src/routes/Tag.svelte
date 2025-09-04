<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import { getMediaWithParams, updateTag } from "../lib/controllers/tags";
  import { navigate } from "svelte-routing";

  let { id, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);

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
</script>

<MediaView
  {title}
  route={routes.tagFunc(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  disableTags={true}
  {updateTitle}
  bind:submittingTitle
/>
