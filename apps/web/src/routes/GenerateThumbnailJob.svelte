<script>
  import { navigate } from "svelte-routing";
  import { create } from "../lib/controllers/job";

  /**
   * @typedef props
   * @type {object}
   * @property {string} mediaId
   * @property {string} redirect
   */
  /** @type {props}*/
  let { mediaId, redirect = null } = $props();
  let timestamp = $state(0);
  let width = $state(400);
  let submitting = $state(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    submitting = true;
    try {
      await create("generate_thumbnail", {
        mediaId,
        timestamp: timestamp,
        width: width,
        relationType: "thumbnail",
      });

      navigate(redirect, { replace: true });
    } finally {
      submitting = false;
    }
  };
  const handleCancel = () => {
    navigate(redirect, { replace: true });
  };
</script>

<div class="container">
  <h1 class="title is-1">Generate chapters</h1>

  <form onsubmit={handleSubmit}>
    <label class="label" for="timestamp">Timestamp</label>
    <div class="field has-addons">
      <div class="control">
        <input
          class="input"
          name="timestamp"
          type="number"
          bind:value={timestamp}
        />
      </div>
      <div class="control">
        <button class="button is-static">seconds</button>
      </div>
    </div>

    <label class="label" for="max-dimension">Width</label>
    <div class="field has-addons">
      <div class="control">
        <input class="input" type="number" name="width" bind:value={width} />
      </div>
      <div class="control">
        <button class="button is-static">px</button>
      </div>
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
