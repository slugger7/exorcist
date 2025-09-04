<script>
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";
  import { imageUrlById } from "../controllers/image";

  /** @import { MediaOverviewDTO } from "../types"*/

  /**
   * @typedef props
   * @type {object}
   * @property {MediaOverviewDTO} video
   * @property {boolean} [selected]
   * @property {boolean} [selecting]
   * @property {() => void} [onselect]
   */
  /** @type {props}*/
  let {
    video,
    selecting = false,
    selected = false,
    onselect = () => {},
  } = $props();
  /** @type {HTMLElement}*/
  let element = $state();

  $effect(() => {
    requestAnimationFrame(() => {
      const item = localStorage.getItem("item");
      if (item === video.id) {
        element.scrollIntoView({ behavior: "smooth" });
        localStorage.removeItem("item");
      }
    });
  });
</script>

<figure
  class={`image ${video.thumbnailId === "00000000-0000-0000-000000000000" ? "is-skeleton" : ""}`}
  bind:this={element}
>
  {#if selecting}
    <button class={`button ${selected ? "is-focused" : ""}`} onclick={onselect}
      ><img
        class="image"
        src={imageUrlById(video.thumbnailId)}
        alt={video.title}
      /></button
    >
  {:else}
    <Link to={routes.videoFunc(video.id, video.title)}
      ><img
        class="image"
        src={imageUrlById(video.thumbnailId)}
        alt={video.title}
      /></Link
    >
  {/if}
</figure>
