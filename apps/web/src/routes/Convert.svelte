<script>
  /** @import { ConvertData, MediaDTO } from "../dto"*/
  /** @import { Errors, GetErrorFn, Touched, Validators } from "../lib/types"*/
  import { onDestroy } from "svelte";
  import { pageState } from "../lib/state/pageState.svelte";
  import { get } from "../lib/controllers/media";
  import { handleValidation } from "../lib/forms/handlers";
  import { create } from "../lib/controllers/job";
  import {
    calculateScaledHeight,
    calculateScaledWidth,
  } from "../lib/forms/dimensions";

  /** @type {{id: string}}*/
  let { id } = $props();
  /** @type {MediaDTO | undefined}*/
  let media = $state();
  let loading = $state(false);
  let submitting = $state(false);
  let filename = $state("");
  let filenameErrors = $state([]);
  let filenameTouched = $state(false);
  let filenameValidators = [
    (value, state) => {
      if (value === state.originalFilename) {
        return "New file name can't be the same as the previous filename";
      }
    },
  ];
  let originalFilename = "";

  let constantRateFactor = $state(23);
  let constantRateFactorErrors = $state([]);
  let forcePixelFormat = $state("yuv420p");
  let copyPeople = $state(true);
  let copyTags = $state(true);
  let height = $state();
  let heightErrors = $state([]);
  let width = $state();
  let widthErrors = $state([]);
  let keepScale = $state(true);

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
    if (media) {
      if (filename === "") {
        filename = media.path.split("/").pop() ?? "";
        originalFilename = filename;
      }

      if (height === undefined) {
        height = media.video?.height;
      }

      if (width === undefined) {
        width = media.video?.width;
      }
    }
  });

  $effect(() => {
    filenameErrors = handleValidation({
      value: filename,
      state: {
        originalFilename,
      },
      validators: filenameValidators,
    });
  });

  /** @param {SubmitEvent} e*/
  const handleSubmit = async (e) => {
    e.preventDefault();

    filenameTouched = true;

    if (
      filenameErrors.length ||
      heightErrors.length ||
      widthErrors.length ||
      constantRateFactorErrors.length
    ) {
      return;
    }

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
            dimension: { height, width },
            copyPeople,
            copyTags,
          },
        });
        history.back();
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

  const scaleWidth = (h) => {
    if (keepScale) {
      width = calculateScaledWidth(
        media?.video?.height ?? 0,
        media?.video?.width ?? 0,
        h,
      );
    }
  };
</script>

<div class="container">
  {#if !loading}
    <h1 class="title is-1">Convert {media?.title}</h1>
    <form onsubmit={handleSubmit}>
      <div class="field">
        <label class="checkbox">
          <input type="checkbox" bind:checked={copyTags} /> Copy Tags
        </label>
      </div>

      <div class="field">
        <label class="checkbox">
          <input type="checkbox" bind:checked={copyPeople} /> Copy People
        </label>
      </div>

      <div class="field">
        <label class="label" for="filename">Filename</label>
        <input
          class={`input ${filenameTouched && filenameErrors && filenameErrors.length > 0 ? "is-danger" : ""}`}
          type="text"
          placeholder="Filename"
          name="filename"
          bind:value={filename}
          onfocus={() => (filenameTouched = true)}
        />
      </div>
      {#if filenameTouched && filenameErrors && filenameErrors.length > 0}
        {#each filenameErrors as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <label class="label" for="height">Height</label>
      <div class="field has-addons">
        <div class="control is-expanded">
          <input
            class={`input ${
              heightErrors && heightErrors.length > 0 ? "is-danger" : ""
            }`}
            type="number"
            name="height"
            bind:value={height}
            placeholder="Height"
            oninput={(e) => {
              const validationMessage = e.target.validationMessage;
              if (validationMessage.length > 0) {
                heightErrors = [validationMessage];
                return;
              } else {
                heightErrors = [];
              }

              scaleWidth(e.target.valueAsNumber);
            }}
          />
        </div>
        <div class="control">
          <button
            class={`button ${height === 1080 ? "is-primary" : ""}`}
            onclick={(e) => {
              e.preventDefault();
              height = 1080;
              scaleWidth(height);
            }}>1080p</button
          >
        </div>
        <div class="control">
          <button
            class={`button ${height === 720 ? "is-primary" : ""}`}
            onclick={(e) => {
              e.preventDefault();
              height = 720;
              scaleWidth(height);
            }}>720p</button
          >
        </div>
        <div class="control">
          <button
            class={`button ${height === 480 ? "is-primary" : ""}`}
            onclick={(e) => {
              e.preventDefault();
              height = 480;
              scaleWidth(height);
            }}>480p</button
          >
        </div>
      </div>
      {#if heightErrors && heightErrors.length > 0}
        {#each heightErrors as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field">
        <label class="label" for="width">Width</label>
        <input
          class={`input ${widthErrors && widthErrors.length > 0 ? "is-danger" : ""}`}
          type="number"
          name="width"
          bind:value={width}
          placeholder="Width"
          oninput={(e) => {
            const validationMessage = e.target.validationMessage;
            if (validationMessage.length > 0) {
              widthErrors = [validationMessage];
              return;
            } else {
              widthErrors = [];
            }

            if (keepScale) {
              height = calculateScaledHeight(
                media?.video?.height ?? 0,
                media?.video?.width ?? 0,
                e.target.valueAsNumber,
              );
            }
          }}
        />
      </div>
      {#if widthErrors && widthErrors.length > 0}
        {#each widthErrors as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field">
        <label class="label checkbox">
          <input class="checkbox" type="checkbox" bind:checked={keepScale} /> Keep
          scale
        </label>
      </div>

      <div class="field">
        <label class="label" for="constantRateFactor"
          >Constant Rate Factor</label
        >
        <input
          class={`input ${constantRateFactorErrors && constantRateFactorErrors.length > 0 ? "is-danger" : ""}`}
          type="number"
          placeholder="Constant Rate Factor"
          name="constantRateFactor"
          bind:value={constantRateFactor}
          oninput={(e) => {
            const validationMessage = e.target.validationMessage;
            if (validationMessage.length > 0) {
              constantRateFactorErrors = [validationMessage];
              return;
            } else {
              constantRateFactorErrors = [];
            }
          }}
        />
      </div>
      {#if constantRateFactorErrors && constantRateFactorErrors.length > 0}
        {#each constantRateFactorErrors as error}
          <p class="help is-danger">{error}</p>
        {/each}
      {/if}

      <div class="field">
        <label class="label" for="forcePixelFormat">Force Pixel Format</label>
        <input
          class={`input`}
          type="text"
          placeholder="Force Pixel Format"
          name="constantRateFactor"
          bind:value={forcePixelFormat}
        />
      </div>

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
