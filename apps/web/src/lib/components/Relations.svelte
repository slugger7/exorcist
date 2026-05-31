<script>
  import VideoCard from "./VideoCard.svelte";

  /**@import { MediaRelationDto } from "../types";

  /**
   * @typedef props
   * @type {object}
   * @property {MediaRelationDto[]} relations
   */
  /** @type {props}*/
  let { relations } = $props();
</script>

<div>
  <h2 class="title is-2 inline">Relations</h2>
  {#each relations as relation}
    <VideoCard video={{ id: relation.relatedToId, title: "" }} />
  {/each}
</div>
