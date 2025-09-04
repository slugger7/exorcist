<script>
  import { navigate } from "svelte-routing";
  import { create } from "../../lib/controllers/libraries";
  import routes from "../routes";

  let libraryName = $state("");
  let submitting = $state(false);

  async function handleSubmit(e) {
    e.preventDefault();
    submitting = true;

    try {
      await create(libraryName);

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
  <h1 class="title is-1">Create Library</h1>
  <form onsubmit={handleSubmit}>
    <div class="field">
      <label class="label" for="library-name">Library Name</label>
      <div class="control">
        <input
          class="input"
          name="library-name"
          bind:value={libraryName}
          placeholder="Library name"
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
