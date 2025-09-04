<script>
  import { Link } from "svelte-routing";
  import JobStatusTag from "./JobStatusTag.svelte";
  import routes from "../../routes/routes";

  /** @import { JobDTO } from "../../dto" */
  /** @type {{job: JobDTO}}*/
  let { job } = $props();
</script>

<div class="card">
  <header class="card-header">
    <p class="card-header-title">
      {job.jobType}&nbsp;<JobStatusTag status={job.status} />
    </p>
  </header>
  <table class="table is-fullwidth">
    <tbody>
      <tr>
        <th>ID</th>
        <td><Link to={`${routes.jobs}?parent=${job.id}`}>{job.id}</Link></td>
      </tr>
      <tr>
        <th>Parent</th>
        <td
          ><Link to={`${routes.jobs}?parent=${job.parent}`}>{job.parent}</Link
          ></td
        >
      </tr>
      <tr>
        <th>Priority</th>
        <td>{job.priority}</td>
      </tr>
      <tr>
        <th>Data</th>
        <td><pre>{job.data}</pre></td>
      </tr>
      <tr>
        <th>Outcome</th>
        <td><pre>{job.outcome}</pre></td>
      </tr>
      <tr>
        <th>Created</th>
        <td>{job.created}</td>
      </tr>
      <tr>
        <th>Modified</th>
        <td>{job.modified}</td>
      </tr>
    </tbody>
  </table>
</div>
