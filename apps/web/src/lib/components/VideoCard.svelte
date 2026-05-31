<script>
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";
  import { thumbnailUrl } from "../controllers/image";

  /**
   * @typedef props
   * @type {object}
   * @property {{id: string, title: string}} video
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
  /** @type {HTMLElement | undefined}*/
  let element = $state();

  $effect(() => {
    requestAnimationFrame(() => {
      const item = localStorage.getItem("item");
      if (item === video.id) {
        if (element) {
          element.scrollIntoView({ behavior: "smooth" });
        }
        localStorage.removeItem("item");
      }
    });
  });
</script>

<figure class={`image`} bind:this={element}>
  {#if selecting}
    <button class={`button ${selected ? "is-focused" : ""}`} onclick={onselect}
      ><img
        class="image"
        src={thumbnailUrl(video.id)}
        alt={video.title}
      /></button
    >
  {:else}
    <Link to={routes.videoFunc(video.id)}
      ><img
        class="image"
        src={thumbnailUrl(video.id)}
        alt={video.title}
      /></Link
    >
  {/if}
</figure>
