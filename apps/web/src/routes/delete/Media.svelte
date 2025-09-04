<script>
  import { navigate } from "svelte-routing";
  import { deleteMedia, get } from "../../lib/controllers/media";
  import routes from "../routes";

  /** @type {{id: string}}*/
  let { id } = $props();
  let deletePhysical = $state(false);
  let deleting = $state(false);

  const handleOnDeleteClick = async (e) => {
    e.preventDefault();

    deleting = true;
    try {
      await deleteMedia(id, deletePhysical);

      if (history.length > 2) {
        history.go(-2);
      } else {
        navigate(routes.home);
      }
    } catch (e) {
      console.error(e);
    } finally {
      deleting = false;
    }
  };

  const handelOnCancelClick = (e) => {
    e.preventDefault();

    history.back();
  };
</script>

<div class="container">
  {#await get(id)}
    <p>loading</p>
  {:then media}
    <h1 class="title is-1">Delete {media.title}</h1>
    <p>
      You are about to delete <strong>{media.title}</strong> with the id
      <strong>{id}</strong> located at <strong>{media.path}</strong>.
    </p>
    <br />
    <p>
      This will only mark the record as deleted. Check the box below if you also
      want to remove the files on disk. Checking this box will also delete any
      generated assets of the media.
    </p>
    <br />
    <form>
      <div class="field">
        <div class="control">
          <label class="checkbox">
            <input type="checkbox" bind:checked={deletePhysical} />
            Delete all physical files on disk
          </label>
        </div>
      </div>
      <br />
      <div class="field is-grouped">
        <div class="control">
          <button
            class="button is-danger"
            onclick={handleOnDeleteClick}
            disabled={deleting}>Delete</button
          >
        </div>
        <div class="control">
          <button class="button" onclick={handelOnCancelClick}>Cancel</button>
        </div>
      </div>
    </form>
  {/await}
</div>
