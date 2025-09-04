<script>
  import { getAllWithParams, ordinals } from "../lib/controllers/people";
  import routes from "./routes";
  import ItemsOverview from "../lib/components/ItemsOverview.svelte";
</script>

<ItemsOverview
  title="People"
  route={routes.people}
  urlFn={routes.personFunc}
  fetch={getAllWithParams}
  {ordinals}
/>
