<script>
  /**
   * @typedef props
   * @type {object}
   * @property {boolean} deleteConfirmation
   * @property {boolean} deleting
   * @property {() => void} handleDelete
   * @property {string} [context]
   */
  /** @type {props}*/
  let {
    deleteConfirmation = $bindable(false),
    deleting,
    handleDelete,
    context,
  } = $props();
</script>

{#if !deleteConfirmation}
  <p class="control">
    <button
      class="button"
      aria-label={`delete ${context}`}
      onclick={() => (deleteConfirmation = true)}
      ><span class="icon"><i class="fas fa-trash"></i></span></button
    >
  </p>
{/if}
{#if deleteConfirmation}
  <p class="control">
    <button
      class="button"
      aria-label="cancel delete"
      onclick={() => (deleteConfirmation = false)}
      disabled={deleting}
      ><span class="icon"><i class="fas fa-xmark"></i></span></button
    >
  </p>
  <p class="control">
    <button
      class="button is-danger"
      aria-label="confirm delete"
      onclick={handleDelete}
      disabled={deleting}
      ><span class="icon"><i class="fas fa-check"></i></span></button
    >
  </p>
{/if}
