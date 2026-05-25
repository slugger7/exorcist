<script>
  import { imageUrlById } from "../controllers/image";

  /**
   * TODO: type the chapters properly from dtos from backend
   * @typedef props
   * @type {object}
   * @property {any[]} chapters
   * @property {(e: Event, chapter: any) => void} onclick
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
