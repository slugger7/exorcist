<script>
  import { Link } from "svelte-routing";
  import HeaderIconButton from "../lib/components/HeaderIconButton.svelte";
  import { create } from "../lib/controllers/job";
  import { get } from "../lib/controllers/libraryPaths";
  import routes from "./routes";

  /** @type {{id: string, name: string}}*/
  let { id, name } = $props();
  let creatingScanJob = $state(false);

  const handleScanPathClick = async () => {
    creatingScanJob = true;

    try {
      await create("scan_path", { libraryPathId: id });
    } catch (e) {
      console.error(e);
    } finally {
      creatingScanJob = false;
    }
  };
</script>

<div class="container">
  <h1 class="title is-1 inline">Library path</h1>
  <HeaderIconButton
    icon="fa-solid fa-arrows-rotate"
    ariaLabel="scan library path"
    onclick={handleScanPathClick}
  />

  {#await get(id)}
    <p>loading</p>
  {:then libraryPath}
    <table class="table">
      <tbody>
        <tr>
          <th>ID</th>
          <td>{libraryPath.id}</td>
        </tr>
        <tr>
          <th>Library ID</th>
          <td
            ><Link to={routes.libraryFunc(libraryPath.libraryId, name)}
              >{libraryPath.libraryId}</Link
            ></td
          >
        </tr>
        <tr>
          <th>Path</th>
          <td>{libraryPath.path}</td>
        </tr>
        <tr>
          <th>Created</th>
          <td>{libraryPath.created}</td>
        </tr>
        <tr>
          <th>Modified</th>
          <td>{libraryPath.modified}</td>
        </tr>
      </tbody>
    </table>
  {:catch}
    <p>something went wrong</p>
  {/await}
</div>
