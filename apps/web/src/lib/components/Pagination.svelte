<script>
  /**
   * @typedef props
   * @type {object}
   * @property {number} page
   * @property {number} limit
   * @property {number} total
   * @property {string} [url]
   * @property {Function} [onchange]
   */
  import { Link } from "svelte-routing";

  /** @type {props}*/
  let {
    page = $bindable(1),
    limit = $bindable(50),
    total,
    url = "",
    onchange = () => {},
  } = $props();
  let pages = $derived(Math.ceil(total / limit));

  /**
   * @param {number} pn
   * @returns {string}
   */
  const urlForPage = (pn) => {
    const params = new URLSearchParams(window.location.search);
    params.set("page", pn.toString());

    return `${url}?${params.toString()}`;
  };

  const clickHandler = (pn) => () => {
    page = pn;
    onchange();
  };

  let pageChangeTimeout;
  const onPageChange = (event) => {
    clearTimeout(pageChangeTimeout);
    setTimeout(() => {
      const newPage = parseInt(event.target.value);
      if (!isNaN(newPage)) {
        clickHandler(newPage)();
      }
    }, 250);
  };
</script>

<nav class="pagination is-centered" aria-label="pagination">
  {#if page - 1 >= 1}
    <Link
      onclick={clickHandler(page - 1)}
      class="pagination-previous"
      to={urlForPage(page - 1)}>Previous</Link
    >
  {/if}
  {#if page + 1 <= pages}
    <Link
      onclick={clickHandler(page + 1)}
      class="pagination-next is-right"
      to={urlForPage(page + 1)}>Next page</Link
    >
  {/if}
  <div class="pagination-list">
    <div class="field has-addons">
      <div class="control">
        <div class="select">
          <select bind:value={limit} onchange={() => onchange()}>
            <option value={10}>10</option>
            <option value={30}>30</option>
            <option value={50}>50</option>
            <option value={100}>100</option>
            <option value={200}>200</option>
          </select>
        </div>
      </div>
      <p class="control">
        <input
          type="number"
          class="input"
          value={page}
          onchange={onPageChange}
        />
      </p>
      <p class="control">
        <button class="button is-static">of {pages} ({total})</button>
      </p>
    </div>
  </div>
</nav>
