<script>
  import { imageUrlById } from "../controllers/image";

  /**
   * @import { ChapterDTO } from "../types"
   * @typedef props
   * @type {object}
   * @property {ChapterDTO[]} chapters
   * @property {(e: Event, chapter: ChapterDTO) => void} onclick
   */
  /** @type {props}*/
  let { chapters, onclick } = $props();
</script>

<div class="grid">
  {#each chapters as chapter (chapter.thumbnailId)}
    <div class="cell">
      <figure class="image">
        <button onclick={(e) => onclick(e, chapter)}>
          <img
            class="image chapter"
            src={imageUrlById(chapter.thumbnailId)}
            alt={`chapter for ${chapter.timestamp} milliseconds`}
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
