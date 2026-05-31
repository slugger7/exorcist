<script>
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";
  import VideoCard from "./VideoCard.svelte";

  /**@import { MediaRelationDto } from "../types";

  /**
   * @typedef props
   * @type {object}
   * @property {string} id
   * @property {MediaRelationDto[]} relations
   */
  /** @type {props}*/
  let { relations, id } = $props();
</script>

<div>
  <div>
    <h2 class="title is-2 inline">Relations</h2>
    <Link to={routes.relateFunc(id)}>Add relations</Link>
  </div>
  <div class="grid">
    {#if relations.length === 0}
      <span>No relations</span>
    {/if}
    {#each relations as relation}
      <div class="cell">
        <VideoCard video={{ id: relation.relatedToId, title: "" }} />
      </div>
    {/each}
  </div>
</div>
