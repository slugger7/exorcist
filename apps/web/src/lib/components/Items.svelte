<script>
  import ItemSelector from "./ItemSelector.svelte";
  import HeaderIconButton from "./HeaderIconButton.svelte";
  import { Link } from "svelte-routing";
  /**
   * @import { Item, FetchItems, RemoveItem, AddItem, CreateItem, ItemUrlFn } from '../types'
   *
   * @typedef props
   * @type {object}
   * @property {string} title
   * @property {Item[]} items
   * @property {FetchItems} fetch
   * @property {RemoveItem} remove
   * @property {AddItem} add
   * @property {CreateItem} create
   * @property {ItemUrlFn} urlFn
   * @property {boolean} [disableEdit]
   */
  /** @type {props}*/
  let {
    title,
    items,
    fetch,
    remove,
    add,
    create,
    urlFn,
    disableEdit = false,
  } = $props();

  let editing = $state(false);
  /** @type {Item[]} */
  let selectedItems = $state([]);
  let allItems = $state([]);
  let loadingItems = $state(false);

  /**
   * @param {Item[]} items
   * @returns {Item[]}
   */
  const sortItems = (items) =>
    items.sort((a, b) => (a.name === b.name ? 0 : a.name < b.name ? -1 : 1));

  $effect(() => {
    if (items !== null && items !== undefined) {
      selectedItems = [...sortItems(items)];
    }
  });

  const refreshItems = async () => {
    loadingItems = true;
    try {
      const items = await fetch();
      allItems = items;
    } finally {
      loadingItems = false;
    }
  };

  $effect(() => {
    if (editing === true) {
      refreshItems();
    }
  });

  /**
   * @param {Item} item
   */
  const handleItemToggle = (item) => {
    const existing = !!selectedItems.find((t) => t.id === item.id);
    if (existing) {
      removeItem(item);
      return;
    }

    addItem(item);
  };

  const addItem = async (item) => {
    selectedItems.push(item);

    try {
      await add(item.id);
    } catch (e) {
      console.error(e);

      await refreshItems();
    }
  };

  /** @param {Item} item */
  const removeItem = async (item) => {
    selectedItems = selectedItems.filter((t) => t.id !== item.id);

    try {
      await remove(item.id);
    } catch (e) {
      console.error(e);

      await refreshItems();
    }
  };

  /** @param {Item} item*/
  const handleRemoveItem = (item) => async () => {
    removeItem(item);
  };

  /** @param {string} itemName */
  const handleCreateItem = async (itemName) => {
    try {
      const createdTag = await create(itemName);

      allItems.push(createdTag);

      addItem(createdTag);
    } catch (e) {
      console.error(e);
    }
  };

  const onEditClick = () => {
    editing = !editing;
  };
</script>

<div>
  <h2 class="title is-2 inline">{title}</h2>
  {#if !disableEdit}
    <HeaderIconButton
      icon="fas fa-pen"
      ariaLabel="edit tags"
      iconClass={editing ? "has-text-info" : ""}
      buttonClass="mb-2"
      onclick={onEditClick}
    />
  {/if}

  {#if editing}
    <ItemSelector
      items={allItems}
      {selectedItems}
      loading={loadingItems}
      toggle={handleItemToggle}
      create={handleCreateItem}
      autofocus={true}
    />
  {/if}
  <div class="field is-grouped is-grouped-multiline">
    {#each selectedItems as item}
      <div class="control">
        <div class="tags has-addons is-medium">
          <Link class="tag is-link" to={urlFn(item.id, item.name)}
            >{item.name}</Link
          >
          {#if editing}
            <button
              class="tag is-delete"
              aria-label="delete {item.name} tag"
              onclick={handleRemoveItem(item)}
            ></button>
          {/if}
        </div>
      </div>
    {/each}
  </div>
</div>
