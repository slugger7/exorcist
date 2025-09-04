<script>
  import { Link } from "svelte-routing";
  import routes from "../../routes/routes";

  /** @type {{ path: string, id: string }}*/
  let { path, id } = $props();
</script>

<div class="card">
  <header class="card-header">
    <Link class="card-header-title" to={routes.libraryPathFunc(id)}>{path}</Link
    >
  </header>
</div>
