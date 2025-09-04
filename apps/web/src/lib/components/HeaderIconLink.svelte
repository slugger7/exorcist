<script>
  /**
   * @typedef props
   * @type {object}
   * @property {string} icon
   * @property {string} ariaLabel
   * @property {string} to
   * @property {any} [state]
   */
  import { Link } from "svelte-routing";

  /** @type {props}*/
  let { icon, ariaLabel, to, state = null } = $props();
</script>

<Link class="mt-2 ml-1 button inline" aria-label={ariaLabel} {to} {state}>
  <span class="icon">
    <i class={icon}></i>
  </span>
</Link>
