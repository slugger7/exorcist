<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import { getMediaWithParams, updatePerson } from "../lib/controllers/people";
  import { navigate } from "svelte-routing";

  let { id, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);

  const updateTitle = async (newTitle) => {
    submittingTitle = true;
    try {
      const updatedPerson = await updatePerson(id, newTitle);

      title = updatedPerson.name;

      const params = new URLSearchParams(location.search);

      navigate(`${routes.personFunc(id, title)}?${params.toString()}`, {
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
  route={routes.personFunc(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  disablePeople={true}
  {updateTitle}
  bind:submittingTitle
/>
