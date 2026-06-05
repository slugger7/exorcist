<script>
  /** @import { ConvertData, MediaDTO } from "../dto"*/
  /** @import { Errors, GetErrorFn, Touched, Validators } from "../lib/types"*/
  import { onDestroy } from "svelte";
  import { pageState } from "../lib/state/pageState.svelte";
  import { get } from "../lib/controllers/media";
  import { handleValidation } from "../lib/forms/handlers";
  import { create } from "../lib/controllers/job";

  /** @type {{id: string}}*/
  let { id } = $props();
  /** @type {MediaDTO | undefined}*/
  let media = $state();
  let loading = $state(false);
  let submitting = $state(false);
  let filename = $state("");
  let originalFilename = "";

  let constantRateFactor = $state(23);
  let forcePixelFormat = $state("yuv420p");

  /** @type {Touched<ConvertData>}*/
  let touched = $state({});

  /** @type {Errors<ConvertData>}*/
  let errors = $state({});

  /** @type {Validators<ConvertData>}*/
  const validators = {
    filename: [
      (value, state) => {
        if (value === state.originalFilename) {
          return "New file name can't be the same as the previous filename";
        }
      },
    ],
  };

  onDestroy(() => {
    pageState.media = undefined;
  });

  /** @param {string} mediaId */
  const fetchMedia = async (mediaId) => {
    loading = true;
    try {
      media = await get(mediaId);
    } finally {
      loading = false;
    }
  };

  $effect(() => {
    if (!pageState.media) {
      fetchMedia(id);
    } else {
      media = pageState.media;
    }
  });

  $effect(() => {
    if (media && filename === "") {
      filename = media.path.split("/").pop() ?? "";
      originalFilename = filename;
    }
  });

  $effect(() => {
    errors.filename = handleValidation({
      value: filename,
      state: {
        originalFilename,
      },
      validators: validators.filename,
    });
  });

  /**
   * @type {GetErrorFn<ConvertData>}
   */
  const getError = (key) => {
    return touched[key] && errors[key] && errors[key].length > 0;
  };

  /** @param {SubmitEvent} e*/
  const handleSubmit = async (e) => {
    e.preventDefault();

    if (media?.id) {
      submitting = true;
      try {
        await create({
          type: "convert",
          data: {
            mediaId: media.id,
            filename,
            constantRateFactor,
            forcePixelFormat,
            dimension: {},
          },
        });
      } finally {
        submitting = false;
      }
    }
  };

  /** @param {Event} e */
  const handleCancel = (e) => {
    e.preventDefault();

    history.back();
  };
</script>

<div class="container">
  {#if !loading}
    <h1 class="title is-1">Convert {media?.title}</h1>
    <form onsubmit={handleSubmit}>
      <div class="field">
        <label class="label" for="filename">Filename</label>
        <input
          class={`input ${getError("filename") ? "is-danger" : ""}`}
          type="text"
          placeholder="Filename"
          name="filename"
          bind:value={filename}
          onfocus={() => (touched.filename = true)}
        />
      </div>
      {#if getError("filename")}
        {#each errors.filename as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field">
        <label class="label" for="constantRateFactor"
          >Constant Rate Factor</label
        >
        <input
          class={`input ${getError("constantRateFactor") ? "is-danger" : ""}`}
          type="number"
          placeholder="Constant Rate Factor"
          name="constantRateFactor"
          bind:value={constantRateFactor}
          onfocus={() => (touched.constantRateFactor = true)}
        />
      </div>
      {#if getError("constantRateFactor")}
        {#each errors.constantRateFactor as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field">
        <label class="label" for="forcePixelFormat">Force Pixel Format</label>
        <input
          class={`input ${getError("forcePixelFormat") ? "is-danger" : ""}`}
          type="text"
          placeholder="Force Pixel Format"
          name="constantRateFactor"
          bind:value={forcePixelFormat}
          onfocus={() => (touched.forcePixelFormat = true)}
        />
      </div>
      {#if getError("forcePixelFormat")}
        {#each errors.forcePixelFormat as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field is-grouped">
        <p class="control">
          <input
            type="submit"
            class={`button is-primary ${submitting ? "is-loading" : ""}`}
            value="Submit"
            disabled={submitting}
          />
        </p>
        <p class="control">
          <button class="button" disabled={submitting} onclick={handleCancel}
            >Cancel</button
          >
        </p>
      </div>
    </form>
  {/if}
</div>
