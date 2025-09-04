<script>
  import { onMount, tick } from "svelte";

  /**
   * @typedef props
   * @type {object}
   * @property {string} value
   * @property {boolean} loading
   * @property {(e: Event, update: string) => void} onsave
   * @property {(e: Event) => void} oncancel
   */
  /** @type {props}*/
  let { value, loading, onsave, oncancel } = $props();
  /** @type {HTMLInputElement}*/
  let inputNode = $state();

  onMount(async () => {
    await tick();

    if (inputNode) {
      inputNode.focus();
    }
  });
</script>

<form
  onsubmit={(e) => {
    e.preventDefault();
    onsave(e, value);
  }}
>
  <div class="field has-addons">
    <p class="control is-expanded">
      <input
        class="input is-large"
        type="text"
        bind:value
        bind:this={inputNode}
      />
    </p>
    <p class="control">
      <button
        class="button is-large"
        type="submit"
        aria-label="save"
        disabled={loading}
      >
        <span class="icon is-large">
          <i class="fas fa-save"></i>
        </span>
      </button>
    </p>
    <p class="control">
      <button
        class="button is-large"
        onclick={oncancel}
        aria-label="cancel"
        disabled={loading}
      >
        <span class="icon is-large">
          <i class="fas fa-xmark"></i>
        </span>
      </button>
    </p>
  </div>
</form>
