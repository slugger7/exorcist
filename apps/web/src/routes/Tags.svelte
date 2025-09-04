<script>
  import { getAllWithParams, ordinals } from "../lib/controllers/tags";
  import routes from "./routes";
  import ItemsOverview from "../lib/components/ItemsOverview.svelte";
</script>

<ItemsOverview
  title="Tags"
  route={routes.tags}
  urlFn={routes.tagFunc}
  fetch={getAllWithParams}
  {ordinals}
/>
