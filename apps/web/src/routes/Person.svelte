<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import routes from "./routes";
  import { ordinals } from "../lib/controllers/media";
  import {
    deletePerson,
    getMediaWithParams,
    updatePerson,
  } from "../lib/controllers/people";
  import { navigate } from "svelte-routing";
  import DeleteAddons from "../lib/components/DeleteAddons.svelte";

  let { id, name } = $props();
  let title = $state(name);
  let submittingTitle = $state(false);
  let deleting = $state(false);
  let deleteConfirmation = $state(false);

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

  const handleDelete = async () => {
    deleting = true;

    try {
      await deletePerson(id);

      window.history.back();
    } finally {
      deleting = false;
    }
  };
</script>

{#snippet headerAddons()}
  <DeleteAddons
    context="person"
    bind:deleteConfirmation
    {deleting}
    {handleDelete}
  />
{/snippet}

<MediaView
  {title}
  route={routes.personFunc(id, name)}
  {ordinals}
  fetchFn={async (params) => await getMediaWithParams(id, params)}
  disablePeople={true}
  {updateTitle}
  bind:submittingTitle
  {headerAddons}
/>
