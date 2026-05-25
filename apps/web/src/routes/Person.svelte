<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import { getMediaWithParams, updatePerson } from "../lib/controllers/people";
  import { navigate } from "svelte-routing";
  import { pageState } from "../lib/state/pageState.svelte";
  import { onMount } from "svelte";

  let { id: paramId, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);
  let id = $derived(pageState.id ? pageState.id : paramId);

  onMount(() => {
    if (pageState.id === undefined) {
      pageState.id = paramId;
    }
  });

  /** @param {string} newTitle */
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
