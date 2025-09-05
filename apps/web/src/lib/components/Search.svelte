<script>
  /**
   * @import { ChangeEventHandler } from "svelte/elements";
   * @import { Ordinal, Item } from "../types"
   */
  import { getAll as getAllTags } from "../controllers/tags";
  import { getAll as getAllPeople } from "../controllers/people";
  import ItemSelector from "./ItemSelector.svelte";
  import Select from "./Select.svelte";
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";
  import { updateProgress } from "../controllers/media";
  import {
    WatchStatus_Watched,
    WatchStatus_InProgress,
    WatchStatus_Unwatched,
  } from "@exorcist-dto";
  import { bind } from "ramda";

  /** @type {Item[]}*/
  const watchStatuses = [
    {
      id: WatchStatus_Watched,
      name: "Watched",
    },
    {
      id: WatchStatus_Unwatched,
      name: "Unwatched",
    },
    {
      id: WatchStatus_InProgress,
      name: "In Progress",
    },
  ];

  /**
   * @typedef props
   * @type {object}
   * @property {string} [value]
   * @property {ChangeEventHandler<HTMLInputElement>} onkeyup
   * @property {Function} [onclear]
   * @property {Ordinal[]} ordinals
   * @property {string} orderBy
   * @property {boolean} [disableTags]
   * @property {boolean} [disablePeople]
   * @property {boolean} [ascending]
   * @property {string[]} [selectedTags]
   * @property {string[]} [selectedPeople]
   * @property {string[]} [selectedWatchStatuses]
   * @property {boolean} [selecting]
   * @property {() => void} [clearSelection]
   * @property {string[]} [selection]
   * @property {boolean} [favourites]
   * @property {boolean} [expanded]
   * @property {boolean} [deleted]
   * @property {boolean} [exists]
   */

  /** @type {props} */
  let {
    value = $bindable(""),
    onkeyup,
    onclear = () => {},
    ordinals,
    orderBy = $bindable(),
    ascending = $bindable(),
    selectedTags = $bindable([]),
    selectedPeople = $bindable([]),
    selectedWatchStatuses = $bindable([]),
    disablePeople = false,
    disableTags = false,
    selecting = $bindable(false),
    clearSelection = () => {},
    selection = [],
    favourites = $bindable(false),
    expanded = $bindable(false),
    deleted = $bindable(false),
    exists = $bindable(true),
  } = $props();
  /** @type {Item[]}*/
  let tags = $state([]);
  let loadingTags = $state(false);
  let people = $state([]);
  let loadingPeople = $state(false);
  let loadingWatchStatus = $state(false);

  const fetchTags = async () => {
    loadingTags = true;
    try {
      tags = await getAllTags();
    } finally {
      loadingTags = false;
    }
  };

  const fetchPeople = async () => {
    loadingPeople = true;
    try {
      people = await getAllPeople();
    } finally {
      loadingPeople = false;
    }
  };

  $effect(() => {
    if (expanded && tags.length === 0) {
      fetchTags();
    }
  });

  $effect(() => {
    if (expanded && people.length === 0) {
      fetchPeople();
    }
  });

  /**
   * @param {Item} item
   */
  const handleTagToggle = (item) => {
    const existing = !!selectedTags.find((t) => t === item.name);
    if (existing) {
      removeSelectedTag(item);
      return;
    }

    selectedTags.push(item.name);
  };

  /** @param {Item} item */
  const handlePersonToggle = (item) => {
    const existing = !!selectedPeople.find((t) => t === item.name);
    if (existing) {
      removeSelectedPerson(item);
      return;
    }

    selectedPeople.push(item.name);
  };

  /** @param {Item} item */
  const handleWatchStatusToggle = (item) => {
    const existing = !!selectedWatchStatuses.find((w) => w === item.id);
    if (existing) {
      removeWatchStatus(item);
      return;
    }

    selectedWatchStatuses.push(item.id);
  };

  /**
   * @param {{name: string}} item
   */
  const removeSelectedTag = (item) => {
    selectedTags = selectedTags.filter((t) => t !== item.name);
  };

  /**
   * @param {{name: string}} item
   */
  const removeSelectedPerson = (item) => {
    selectedPeople = selectedPeople.filter((t) => t !== item.name);
  };

  /**
   * @param {{id: string}} item
   */
  const removeWatchStatus = (item) => {
    selectedWatchStatuses = selectedWatchStatuses.filter((ws) => ws != item.id);
  };

  /**
   * @param {string[]} selectedItems
   * @param {Item[]} items
   * @returns {Item[]}
   */
  const hydrateItemsByName = (selectedItems, items) =>
    selectedItems.map((s) => items.find((i) => i.name === s)).filter((s) => s);

  /**
   * @param {string[]} selectedItems
   * @param {Item[]} items
   * @returns {Item[]}
   */
  const hydrateItemsById = (selectedItems, items) =>
    selectedItems.map((s) => items.find((i) => i.id === s)).filter((s) => s);

  const handleWatchedClick = async () => {
    let loadingWatchStatus = true;
    try {
      await Promise.all(
        selection.map(async (s) => {
          await updateProgress(s, -1, true);
        }),
      );
    } finally {
      loadingWatchStatus = false;
    }
  };
  const handleUnwatchClick = async () => {
    loadingWatchStatus = true;
    try {
      await Promise.all(
        selection.map(async (s) => {
          await updateProgress(s, 0, true);
        }),
      );
    } finally {
      loadingWatchStatus = false;
    }
  };
</script>

<div class="block">
  <div class="field has-addons">
    <p class="control has-icons-left is-expanded">
      <input
        class="input"
        type="text"
        placeholder="Search"
        bind:value
        {onkeyup}
      />
      <span class="icon is-left">
        <i class="fas fa-search"></i>
      </span>
    </p>
    {#if value !== ""}
      <p class="control">
        <button
          class="button"
          aria-label="clear search"
          onclick={() => {
            value = "";
            onclear();
          }}
        >
          <span class="icon">
            <i class="fas fa-close"></i>
          </span>
        </button>
      </p>
    {/if}
    <p class="control">
      <button
        class="button"
        aria-label="extra filters"
        onclick={() => {
          expanded = !expanded;
        }}
      >
        {#if expanded}
          <span class="icon is-right">
            <i class="fas fa-chevron-up"></i>
          </span>
        {:else}
          <span class="icon is-right">
            <i class="fas fa-chevron-down"></i>
          </span>
        {/if}
      </button>
    </p>
  </div>
</div>

{#if expanded}
  <div class="block">
    <div class="field has-addons">
      <p class="control">
        <button
          class="button"
          aria-label="change sort order"
          onclick={() => {
            ascending = !ascending;
          }}
        >
          <span class="icon is-left">
            {#if ascending}
              <i class="fas fa-arrow-up-short-wide"></i>
            {:else}
              <i class="fas fa-arrow-up-wide-short"></i>
            {/if}
          </span>
        </button>
      </p>
      <p class="control">
        <Select options={ordinals} bind:value={orderBy} />
      </p>

      <p class="control">
        <button
          class="button"
          aria-label="filter by favourites"
          onclick={() => {
            favourites = !favourites;
          }}
        >
          <span class="icon">
            <i class={`${favourites ? "fa fa-heart" : "fa-regular fa-heart"}`}
            ></i>
          </span>
        </button>
      </p>
    </div>
  </div>
  <div class="block">
    <ItemSelector
      placeholder="watch status"
      items={watchStatuses}
      selectedItems={hydrateItemsById(selectedWatchStatuses, watchStatuses)}
      loading={false}
      toggle={handleWatchStatusToggle}
    />
    <div class="field is-grouped is-grouped-multiline">
      {#each selectedWatchStatuses as status}
        <div class="control">
          <div class="tags has-addons is-medium">
            <span class="tag"
              >{watchStatuses.find((ws) => ws.id == status).name}</span
            >

            <button
              class="tag is-delete"
              aria-label="remove {status} from filter"
              onclick={() => removeWatchStatus({ id: status })}
            ></button>
          </div>
        </div>
      {/each}
    </div>
    {#if !disableTags}
      <ItemSelector
        placeholder="tag"
        items={tags}
        selectedItems={hydrateItemsByName(selectedTags, tags)}
        loading={loadingTags}
        toggle={handleTagToggle}
        disableCreate={true}
      />
      <div class="field is-grouped is-grouped-multiline">
        {#each selectedTags as tag}
          <div class="control">
            <div class="tags has-addons is-medium">
              <span class="tag">{tag}</span>

              <button
                class="tag is-delete"
                aria-label="remove {tag} from filter"
                onclick={() => removeSelectedTag({ name: tag })}
              ></button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
    {#if !disablePeople}
      <ItemSelector
        placeholder="people"
        items={people}
        selectedItems={hydrateItemsByName(selectedPeople, people)}
        loading={loadingPeople}
        toggle={handlePersonToggle}
        disableCreate={true}
      />
      <div class="field is-grouped is-grouped-multiline">
        {#each selectedPeople as person}
          <div class="control">
            <div class="tags has-addons is-medium">
              <span class="tag">{person}</span>

              <button
                class="tag is-delete"
                aria-label="remove {person} from filter"
                onclick={() => removeSelectedPerson({ name: person })}
              ></button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
    <div class="field">
      <div class="checkboxes">
        <label class="checkbox">
          <input type="checkbox" bind:checked={deleted} />
          Deleted
        </label>
        <label class="checkbox">
          <input type="checkbox" bind:checked={exists} />
          Exists
        </label>
      </div>
    </div>
    <div class="field is-grouped">
      <button
        class={`button ${selecting ? "is-primary" : ""}`}
        onclick={() => (selecting = !selecting)}
        ><span>Select</span><span class="icon"
          ><i class={`fas ${selecting ? "fa-xmark" : "fa-check"}`}></i></span
        >
      </button>
      {#if selection.length > 0}
        <button class="button" onclick={clearSelection}>Clear</button>
        <div class="field has-addons">
          <p class="control">
            <Link class="button" to={routes.playlistAddFn(selection)}>
              <span class="icon">
                <i class="fas fa-plus"></i>
              </span>
            </Link>
          </p>
          <p class="control">
            <button
              class={`button ${loadingWatchStatus ? "is-loading" : ""}`}
              aria-label="mark selection watched"
              onclick={handleWatchedClick}
              disabled={loadingWatchStatus}
            >
              <span class="icon">
                <i class="fas fa-eye"></i>
              </span>
            </button>
          </p>
          <p class="control">
            <button
              class={`button ${loadingWatchStatus ? "is-loading" : ""}`}
              aria-label="mark selection unwatched"
              onclick={handleUnwatchClick}
              disabled={loadingWatchStatus}
            >
              <span class="icon">
                <i class="fas fa-eye-slash"></i>
              </span>
            </button>
          </p>
        </div>
      {/if}
    </div>
  </div>
{/if}
