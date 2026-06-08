<script>
  /** @import {MediaRelationDto} from "../../dto"*/
  import { imageUrlById } from "../controllers/image";

  /**
   * @typedef props
   * @type {object}
   * @property {MediaRelationDto[]} chapters
   * @property {(e: Event, chapter: MediaRelationDto) => void} onclick
   */
  /** @type {props}*/
  let { chapters, onclick } = $props();
</script>

<div class="grid">
  {#each chapters as chapter (chapter.relatedToId)}
    <div class="cell">
      <figure class="image">
        <button onclick={(e) => onclick(e, chapter)}>
          <img
            class="image chapter"
            src={imageUrlById(chapter.relatedToId)}
            alt={`chapter for ${chapter.metadata?.timestamp} milliseconds`}
          />
        </button>
      </figure>
    </div>
  {/each}
</div>

<style>
  .chapter {
    cursor: pointer;
  }
</style>
