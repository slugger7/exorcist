<script>
  import { navigate } from "svelte-routing";
  import routes from "../routes";
  import { create } from "../../lib/controllers/libraryPaths";

  let libraryPath = $state("");
  let submitting = $state(false);

  /** @type {{libraryId: string}}*/
  let { libraryId } = $props();
  /** @type {{libraryName: string}}*/
  let { libraryName } = history.state;

  async function handleSubmit(e) {
    e.preventDefault();
    submitting = true;

    try {
      await create(libraryId, libraryPath);

      navigate(routes.libraries);
    } catch (e) {
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  const handelCancelClick = (e) => {
    e.preventDefault();

    history.back();
  };
</script>

<div class="container">
  <h1 class="title is-1">New Path for {libraryName}</h1>
  <form onsubmit={handleSubmit}>
    <div class="field">
      <label class="label" for="library-path">Path</label>
      <div class="control">
        <input
          class="input"
          name="library-path"
          bind:value={libraryPath}
          placeholder="/"
          disabled={submitting}
        />
      </div>
    </div>
    <div class="field is-grouped">
      <div class="control">
        <button
          class="button is-primary {submitting ? 'is-loading' : ''}"
          disabled={submitting}>Create</button
        >
      </div>
      <div class="control">
        <button
          class="button {submitting ? 'is-loading' : ''}"
          disabled={submitting}
          onclick={handelCancelClick}>Cancel</button
        >
      </div>
    </div>
  </form>
</div>
