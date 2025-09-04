<script>
  import { navigate } from "svelte-routing";
  import { create } from "../../lib/controllers/playlists";
  import routes from "../routes";

  let playlistName = $state("");
  let submitting = $state(false);

  async function handleSubmit(e) {
    e.preventDefault();
    submitting = true;

    try {
      await create(playlistName);

      navigate(routes.playlists);
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
  <h1 class="title is-1">Create Playlist</h1>
  <form onsubmit={handleSubmit}>
    <div class="field">
      <label class="label" for="playlist-name">Playlist Name</label>
      <div class="control">
        <input
          class="input"
          name="playlist-name"
          bind:value={playlistName}
          placeholder="Playlist name"
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
