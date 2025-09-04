<script>
  /** @import { JobStatusEnum } from '../../dto/model'*/

  /** @type {{ status: JobStatusEnum}}*/
  let { status } = $props();

  const getColourClass = () => {
    switch (status) {
      case "not_started":
        return "is-dark";
      case "in_progress":
        return "is-primary";
      case "failed":
        return "is-danger";
      case "completed":
        return "is-success";
      case "cancelled":
        return "is-warning";
    }
  };
</script>

<span class={`tag ${getColourClass()}`}>{status}</span>
