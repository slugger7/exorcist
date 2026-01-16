<script>
  import { onMount, tick } from "svelte";
  import { nextFocusState } from "../state/nextFocus.svelte";

  /**
   * @import { Item } from "../types"
   *
   * @callback ToggleItem
   * @param {Item} item
   *
   * @callback CreateItemNoReturn
   * @param {string} itemName
   *
   * @typedef props
   * @type {object}
   * @property {boolean} loading
   * @property {Item[]} items
   * @property {Item[]} selectedItems
   * @property {ToggleItem} toggle
   * @property {CreateItemNoReturn} [create]
   * @property {boolean} [disableCreate]
   * @property {string} [placeholder]
   * @property {boolean} [autofocus]
   */
  /** @type {props}*/
  let {
    loading,
    items,
    selectedItems,
    toggle,
    create,
    disableCreate = false,
    placeholder = "",
    autofocus = false,
  } = $props();

  const itemsInViewCount = 5;

  let query = $state("");
  let active = $state(false);
  let itemsInView = $state([]);
  let selectedIndex = $state(null);
  let showCreateItem = $derived(
    !disableCreate &&
      query.length >= 1 &&
      !items.find((item) => item.name.toLowerCase() === query.toLowerCase()),
  );
  let selectionBoundary = $derived(
    showCreateItem ? itemsInView.length : itemsInView.length - 1,
  );
  /** @type {HTMLInputElement}*/
  let inputNode = $state();

  onMount(async () => {
    await tick();

    if (inputNode && autofocus) {
      inputNode.focus();
    }
  });

  $effect(() => {
    itemsInView = items.slice(0, itemsInViewCount);
  });

  const onDropdownFocus = () => {
    active = true;
  };

  const onDropdownBlur = () => {
    active = false;
  };

  const reset = () => {
    query = "";
    itemsInView = itemsInView = items.slice(0, itemsInViewCount);
    selectedIndex = null;
  };

  const onMouseDown = (item) => (e) => {
    e.preventDefault();
    toggle(item);
  };

  const handleUpArrow = () => {
    if (selectedIndex === null) {
      selectedIndex = selectionBoundary;
    } else if (selectedIndex === 0) {
      selectedIndex = selectionBoundary;
    } else {
      selectedIndex--;
    }
  };

  const handleDownArrow = () => {
    if (selectedIndex === null) {
      selectedIndex = 0;
    } else if (selectedIndex === selectionBoundary) {
      selectedIndex = 0;
    } else {
      selectedIndex++;
    }
  };

  const handleEnter = () => {
    if (selectedIndex < itemsInView.length) {
      const item = itemsInView[selectedIndex];

      toggle(item);
    } else if (showCreateItem) {
      create(query);
    }
    reset();
  };

  const handleEscape = () => {
    active = false;
    inputNode.blur();

    if (nextFocusState.node) {
      nextFocusState.node.focus();
      nextFocusState.node = inputNode;
    }
  };

  const onKeyDown = (e) => {
    if (e.code === "ArrowUp") {
      e.preventDefault();
      return handleUpArrow();
    }
    if (e.code === "ArrowDown") {
      e.preventDefault();
      return handleDownArrow();
    }
    if (e.code === "Enter") {
      e.preventDefault();
      return handleEnter();
    }
    if (e.code === "Escape") {
      e.preventDefault();
      return handleEscape();
    }
  };

  /** @param {string} str */
  const makeSearchableString = (str) =>
    str.replaceAll(/[^a-zA-Z]/g, "").toLowerCase();

  const onQueryChange = (e) => {
    e.preventDefault();
    query = e.target.value;

    selectedIndex = 0;

    itemsInView = items
      .filter((t) =>
        makeSearchableString(t.name).includes(
          makeSearchableString(e.target.value),
        ),
      )
      .slice(0, itemsInViewCount);

    const exact = items.find(
      (tag) => makeSearchableString(tag.name) === makeSearchableString(query),
    );
    if (exact) {
      itemsInView = itemsInView.filter(
        (t) => t.name.toLowerCase() !== exact.name.toLowerCase(),
      );
      itemsInView = [exact, ...itemsInView];

      if (itemsInView.length > itemsInViewCount) {
        itemsInView.pop();
      }
    }

    active = query.length > 0;

    if (items.length === 1) {
      selectedIndex = 0;
    }
  };

  /**
   * @param {Item} item
   */
  const itemInSelection = (item) =>
    selectedItems.find((selectedItem) => item.id === selectedItem.id);

  const onMouseDownCreate = () => {
    create(query);
    reset();
  };
</script>

<div class="field is-grouped">
  <div class="control">
    <div class="dropdown {active ? 'is-active' : ''}">
      <div class="dropdown-trigger">
        <input
          class="input"
          type="text"
          {placeholder}
          value={query}
          onfocus={onDropdownFocus}
          onblur={onDropdownBlur}
          onkeydown={onKeyDown}
          oninput={onQueryChange}
          bind:this={inputNode}
        />
      </div>
      <div class="dropdown-menu" role="menu">
        <div class="dropdown-content">
          {#if loading}
            <p class="dropdown-item">Loading tags</p>
          {:else}
            {#each itemsInView as item, i (item.id)}
              <button
                class={`dropdown-item ${i === selectedIndex ? "is-active" : ""}`}
                onmousedown={onMouseDown(item)}
              >
                {#if itemInSelection(item)}
                  <span class="icon has-text-success"
                    ><i class="fas fa-square-check"></i></span
                  >
                {:else}
                  <span class="icon"><i class="far fa-square-check"></i></span>
                {/if}
                {item.name}</button
              >
            {/each}
            {#if showCreateItem}
              <hr class="dropdown-divider" />
              <button
                class={`dropdown-item ${selectedIndex === itemsInView.length ? "is-active" : ""}`}
                onmousedown={onMouseDownCreate}>Create {query}</button
              >
            {/if}
          {/if}
        </div>
      </div>
    </div>
  </div>
</div>
