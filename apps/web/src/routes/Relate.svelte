<script>
  import MediaView from "../lib/components/MediaView.svelte";
  import {
    getAllWithParams,
    ordinals,
    relateMedia,
  } from "../lib/controllers/media";
  import routes from "./routes";

  /**
   * @typedef props
   * @type {object}
   * @property {string} id
   */
  /** @type {props}*/
  const { id } = $props();
  /** @type {string[]}*/
  let selectedMedia = $state([]);
  let submitting = $state(false);
  let backrelate = $state(false);
  let interrelate = $state(false);

  /** @param {SubmitEvent} e */
  const handleSubmit = async (e) => {
    e.preventDefault();

    submitting = true;
    try {
      await relateMedia(id, {
        relatedToIds: selectedMedia,
        backrelate,
        interrelate,
      });
      history.back();
    } finally {
      submitting = false;
    }
  };

  /** @param {Event} e*/
  const handleCancel = (e) => {
    e.preventDefault();
    history.back();
  };
</script>

<div>
  <MediaView
    title="Relate media"
    route={routes.relateFunc(id)}
    {ordinals}
    fetchFn={getAllWithParams}
    picker={true}
    bind:selectedMedia
  />
  <div class="container">
    <form onsubmit={handleSubmit}>
      <div class="field">
        <label class="checkbox">
          <input type="checkbox" bind:checked={backrelate} /> Back Relate
        </label>
      </div>
      <div class="field">
        <label class="checkbox">
          <input type="checkbox" bind:checked={interrelate} /> Interrelate
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
          <button class="button" disabled={submitting} onclick={handleCancel}
            >Cancel</button
          >
        </p>
      </div>
    </form>
  </div>
  <br />
  <div class="container">
    <article class="message">
      <div class="message-header">Back relate</div>
      <div class="message-body">
        <p>
          Back relating will create the normal relation forward as well as
          adding a relation back to the initial media
        </p>
      </div>
    </article>
    <article class="message">
      <div class="message-header">Interrrelate</div>
      <div class="message-body">
        <p>
          Inter relating will create connections forwards and backwards between
          everything that was selected.
        </p>
      </div>
    </article>
    <article class="message">
      <div class="message-header">Note</div>
      <div class="message-body">
        <p>
          Selecting inter relate will create bi-directional relationships
          between selected media but only a one way direction with from the
          media to the selected media
        </p>
      </div>
    </article>
  </div>
</div>
