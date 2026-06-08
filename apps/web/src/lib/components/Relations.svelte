<script>
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";
  import VideoCard from "./VideoCard.svelte";
  import HeaderIconButton from "./HeaderIconButton.svelte";
  import DeleteAddons from "./DeleteAddons.svelte";
  import { deleteRelations } from "../controllers/media";

  /**@import { MediaRelationDto } from "../types";

  /**
   * @typedef props
   * @type {object}
   * @property {string} id
   * @property {MediaRelationDto[]} [relations]
   */
  /** @type {props}*/
  let { relations, id } = $props();
  let selecting = $state(false);
  /** @type {string[]}*/
  let selected = $state([]);
  let deleting = $state(false);
  let deleteConfirmation = $state(false);

  /** @param {string} id*/
  const handleSelected = (id) => {
    if (isSelected(id)) {
      selected = selected.filter((m) => m !== id);
    } else {
      selected.push(id);
    }
  };

  /**
   * @param {string} id
   * @returns {boolean}
   */
  const isSelected = (id) => selected.includes(id);

  const handleDelete = async () => {
    deleting = true;

    try {
      await deleteRelations(id, { relatedToIds: selected });

      relations = relations?.filter(
        (rel) => !selected.includes(rel.relatedToId),
      );

      selected = [];
      deleteConfirmation = false;
    } finally {
      deleting = false;
    }
  };
</script>

<div>
  <div>
    <h2 class="title is-2 inline">Relations</h2>

    <HeaderIconButton
      icon="fas fa-pen"
      ariaLabel="edit relations"
      buttonClass={`${selecting ? "has-text-info" : ""}`}
      onclick={() => (selecting = !selecting)}
    />

    <Link to={routes.relateFunc(id)}>Add relations</Link>
  </div>
  <br />
  {#if selecting}
    <div class="block">
      <div class="field has-addons">
        <DeleteAddons
          {handleDelete}
          {deleting}
          bind:deleteConfirmation
          disabled={selected.length === 0}
        />
      </div>
    </div>
  {/if}
  <div class="grid">
    {#if relations?.length === 0}
      <span>No relations</span>
    {/if}
    {#each relations as relation (relation.id)}
      <div class="cell">
        <VideoCard
          video={{ id: relation.relatedToId, title: "" }}
          {selecting}
          selected={isSelected(relation.relatedToId)}
          onselect={() => handleSelected(relation.relatedToId)}
        />
      </div>
    {/each}
  </div>
</div>
