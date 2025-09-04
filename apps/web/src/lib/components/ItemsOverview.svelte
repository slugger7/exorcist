<script>
  import { Link } from "svelte-routing";
  import {
    getBoolParamOrDefault,
    getStringSearchParamOrDefault,
    setValueAndNavigate,
  } from "../searchParams";
  import Search from "./Search.svelte";
  /**
   * @import { SearchItemsWithParams, ItemUrlFn, Item, Ordinal } from "../types"
   *
   * @typedef props
   * @type {object}
   * @property {string} title
   * @property {string} route
   * @property {ItemUrlFn} urlFn
   * @property {SearchItemsWithParams} fetch
   * @property {Ordinal[]} ordinals
   */
  /** @type {props} */
  let { title, route, ordinals, fetch, urlFn } = $props();

  /** @type {Item[]}*/
  let items = $state([]);
  let value = $state("");
  let search = $state(getStringSearchParamOrDefault("search", ""));
  let orderBy = $state(getStringSearchParamOrDefault("orderBy", "count"));
  let ascending = $state(getBoolParamOrDefault("ascending", false));
  let loading = $state(false);

  const fetchItems = async () => {
    const params = new URLSearchParams();

    params.set("asc", ascending.toString());
    params.set("orderBy", orderBy);

    if (search !== "") {
      params.set("search", search);
    }

    loading = true;
    try {
      items = await fetch(params);
    } finally {
      loading = false;
    }
  };

  $effect(() => {
    fetchItems();
  });

  $effect(() => {
    setValueAndNavigate("asc", ascending, route, { replace: true });
  });

  $effect(() => {
    setValueAndNavigate("orderBy", orderBy, route, { replace: true });
  });

  let searchTimeout;
  const onkeyup = (event) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      search = event.target.value;
      setValueAndNavigate("search", event.target.value, route, {
        replace: true,
      });

      event.target.focus();
    }, 500);
  };

  const onclear = () => {
    value = "";
    search = value;
    setValueAndNavigate("search", value, route, { replace: false });
  };
</script>

<div class="container is-fluid">
  <h1 class="title is-1">{title}</h1>
  <Search
    {onkeyup}
    value={search}
    {onclear}
    bind:orderBy
    bind:ascending
    {ordinals}
  />
  {#if loading}
    <p>Loading items</p>
  {:else if items}
    <div class="field is-grouped is-grouped-multiline">
      {#each items as item (item.id)}
        <div class="control">
          <div class="tags">
            <Link class="tag is-link" to={urlFn(item.id, item.name)}
              >{item.name}</Link
            >
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
