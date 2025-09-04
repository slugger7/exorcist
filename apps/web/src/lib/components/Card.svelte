<script>
  /** @import { ItemUrlFn } from "../types"*/
  import { Link } from "svelte-routing";

  /** @typedef props
   * @type {object}
   * @property {string} name
   * @property {string} id
   * @property {ItemUrlFn} toFunc
   * @property {ItemUrlFn} [toEditFunc]
   */
  /** @type {props} */
  let { name, id, toFunc, toEditFunc = null } = $props();
</script>

<div class="card">
  <header class="card-header">
    <Link class="card-header-title" to={toFunc(id, name)}>{name}</Link>
    {#if toEditFunc}
      <Link
        class="card-header-icon"
        aria-label="more-options"
        to={toEditFunc(id, name)}
      >
        <span class="icon">
          <i class="fa-solid fa-pen" aria-hidden="true"></i>
        </span>
      </Link>
    {/if}
  </header>
</div>
