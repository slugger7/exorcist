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
  let interval = $state(60);
  let overwrite = $state(false);
  let maxDimension = $state(400);
  let submitting = $state(false);

  const handleSubmit = async (e) => {
    e.preventDefault();

    submitting = true;
    try {
      await create("generate_chapters", {
        mediaId,
        interval,
        overwrite,
        maxDimension,
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
    <label class="label" for="interval">Interval</label>
    <div class="field has-addons">
      <div class="control">
        <input
          class="input"
          name="interval"
          type="number"
          bind:value={interval}
        />
      </div>
      <div class="control">
        <button class="button is-static">seconds</button>
      </div>
    </div>

    <label class="label" for="max-dimension">Max dimension</label>
    <div class="field has-addons">
      <div class="control">
        <input
          class="input"
          type="number"
          name="max-dimension"
          bind:value={maxDimension}
        />
      </div>
      <div class="control">
        <button class="button is-static">px</button>
      </div>
    </div>

    <div class="field">
      <label class="checkbox">
        <input type="checkbox" bind:checked={overwrite} /> Overwrite existing chapters
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
