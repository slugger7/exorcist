<script>
  /**
   * @import { WSMessage, WSTopicMap } from "../lib/types"
   * @import { JobDTO, PageDTO } from '../dto'
   */
  import { onDestroy, onMount } from "svelte";
  import Pagination from "../lib/components/Pagination.svelte";
  import { getAll } from "../lib/controllers/job";
  import {
    getArrayOfStringsSearchParamOrDefault,
    getIntSearchParamOrDefault,
    getStringSearchParamOrDefault,
    setValueAndNavigate,
  } from "../lib/searchParams";
  import routes from "./routes";
  import { wsState } from "../lib/state/wsState.svelte";
  import { PONG } from "../lib/constants/websocket";
  import JobCard from "../lib/components/JobCard.svelte";

  let parent = $state(getStringSearchParamOrDefault("parent", ""));
  let statuses = $state(getArrayOfStringsSearchParamOrDefault("status", []));
  let types = $state(getArrayOfStringsSearchParamOrDefault("type", []));
  let page = $state(getIntSearchParamOrDefault("page", 1));
  let limit = $state(getIntSearchParamOrDefault("limit", 100));
  let loading = $state(false);
  let error = $state();
  /** @type {PageDTO<JobDTO>}*/
  let jobPage = $state();

  const onPopState = () => {
    parent = getStringSearchParamOrDefault("parent", "");
    statuses = getArrayOfStringsSearchParamOrDefault("status", []);
    types = getArrayOfStringsSearchParamOrDefault("type", []);
    page = getIntSearchParamOrDefault("page", 1);
    limit = getIntSearchParamOrDefault("limit", 100);
  };

  onMount(() => {
    window.addEventListener("popstate", onPopState);
  });

  onDestroy(() => {
    window.removeEventListener("popstate", onPopState);

    if (!wsState.active || !wsState.connection) return;
    wsState.connection.removeEventListener("message", onmessage);
  });

  const syncJobs = async () => {
    loading = true;
    try {
      jobPage = await getAll(page, limit, parent, statuses, types);
    } catch (e) {
      error = e;
    } finally {
      loading = false;
    }
  };

  $effect(() => {
    syncJobs();
  });

  $effect(() => {
    setValueAndNavigate("parent", parent, routes.jobs);
  });

  $effect(() => {
    if (!wsState.active || !wsState.connection) return;
    wsState.connection.addEventListener("message", onmessage);
  });

  $effect(() => {
    const params = new URLSearchParams(location.search);

    setValueAndNavigate("page", page.toString(), routes.jobs, {
      replace: params.get("page") !== "",
    });
  });

  $effect(() => {
    setValueAndNavigate("limit", limit.toString(), routes.jobs, {
      replace: true,
    });
  });

  /** @param {MessageEvent<string>} e */
  const onmessage = (e) => {
    if (e.data === PONG) return;

    /** @type {WSMessage<JobDTO>}*/
    const data = JSON.parse(e.data);

    const topicFn = topics[data.topic];
    if (topicFn) {
      topicFn(data.data);
    }
  };

  /** @type {WSTopicMap<JobDTO>}*/
  const topics = {
    job_update: (job) => {
      jobPage.data = jobPage.data.map((j) => {
        if (j.id === job.id) {
          return job;
        }
        return j;
      });
    },
    job_create: (job) => {
      jobPage.total = jobPage.total + 1;
      jobPage.data.push(job);
    },
  };
</script>

<div class="container">
  <h1 class="title is-1">Jobs</h1>
  {#if loading}
    <p>loading</p>
  {:else if !error && jobPage}
    {#each jobPage.data as job (job.id)}
      <JobCard {job} />
    {:else}
      <p>No jobs here</p>
    {/each}

    <Pagination bind:page bind:limit total={jobPage.total} url={routes.jobs} />
    <div class="section"></div>
  {:else if error}
    <p>something went wrong ${error.message}</p>
  {/if}
</div>
