<script>
  import { navigate } from "svelte-routing";
  import { create } from "../lib/controllers/job";

  /**
   * @typedef props
   * @type {object}
   * @property {string} [mediaId]
   * @property {string} [libraryId]
   * @property {string} redirect
   */
  /** @type {props}*/
  let { mediaId, libraryId, redirect } = $props();
  let submitting = $state(false);
  let size = $state(true);
  let checksum = $state(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    const jobType = mediaId
      ? "refresh_metadata"
      : libraryId
        ? "refresh_library_metadata"
        : null;
    submitting = true;
    try {
      await create(jobType, {
        libraryId,
        mediaId,
        batchSize: 50,
        refreshFields: { size, checksum },
      });

      navigate(redirect);
    } finally {
      submitting = false;
    }
  };

  const handleCancel = () => {
    navigate(redirect, { replace: true });
  };
</script>

<div class="container">
  <h1 class="title is-1">
    Refresh metadata for {mediaId ? "media" : libraryId ? "library" : "UNKNOWN"}
  </h1>
  <div class="block">
    <h3 class="title is-3">Select the data to refresh</h3>
  </div>
  <form onsubmit={handleSubmit}>
    <div class="field">
      <label class="checkbox">
        <input type="checkbox" bind:checked={size} /> Size
      </label>
    </div>

    <div class="field">
      <label class="checkbox">
        <input type="checkbox" bind:checked={checksum} /> Checksum
      </label>
    </div>

    <div class="field is-grouped">
      <p class="control">
        <input
          type="submit"
          class={`button is-primary ${submitting ? "is-loading" : ""}`}
          value="Submit"
          disabled={submitting}
        />
      </p>
      <p class="control">
        <button class="button" disabled={submitting} onclick={handleCancel}>
          Cancel
        </button>
      </p>
    </div>
  </form>
</div>
