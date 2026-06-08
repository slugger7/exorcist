<script>
  /** @import { Item, MediaDTO, WSMessage, WSTopicMap, ChapterMetadadataDTO } from "../lib/types";*/
  import { onDestroy } from "svelte";
  import { imageUrlById } from "../lib/controllers/image";
  import {
    videoUrlById,
    get,
    updateProgress,
    updateMedia,
  } from "../lib/controllers/media";
  import Items from "../lib/components/Items.svelte";
  import {
    add as addTag,
    create as createTags,
    getAll as getAllTags,
    remove as removeTag,
  } from "../lib/controllers/tags";
  import {
    getAll as getAllPeople,
    add as addPerson,
    create as createPeople,
    remove as removePerson,
  } from "../lib/controllers/people";
  import routes from "./routes";
  import { nextFocusState } from "../lib/state/nextFocus.svelte";
  import { formatFileSize } from "../lib/formatting/filesize";
  import { formatRuntime } from "../lib/formatting/runtime";
  import HeaderIconButton from "../lib/components/HeaderIconButton.svelte";
  import EditHeading from "../lib/components/EditHeading.svelte";
  import { Link } from "svelte-routing";
  import { addFavourite, removeFavourite } from "../lib/controllers/users";
  import Chapters from "../lib/components/Chapters.svelte";
  import { wsState } from "../lib/state/wsState.svelte";
  import { PONG } from "../lib/constants/websocket";
  import Relations from "../lib/components/Relations.svelte";
  import { pageState } from "../lib/state/pageState.svelte";
  /** @type {{id: string}}*/
  let { id } = $props();
  /** @type {HTMLVideoElement | undefined}*/
  let videoNode = $state();
  /** @type {MediaDTO | undefined}*/
  let mediaEntity = $state();
  let loadingMedia = $state(false);
  let loadingProgress = $state(false);
  let editingTitle = $state(false);
  let loadingTitle = $state(false);
  let loadingFavourite = $state(false);

  let watchedPercentage = $derived(
    mediaEntity
      ? mediaEntity.progress /
          (mediaEntity.video ? mediaEntity.video.runtime : 1)
      : 0,
  );
  let thumbnailRelation = $derived(
    mediaEntity?.relations?.find(
      (relation) => relation.relationType === "thumbnail",
    ),
  );

  /** @param {string} mediaId */
  const fetchMedia = async (mediaId) => {
    loadingMedia = true;
    try {
      mediaEntity = await get(mediaId);
    } finally {
      loadingMedia = false;
    }
  };

  $effect(() => {
    fetchMedia(id);
  });

  onDestroy(() => {
    localStorage.setItem("item", id);
    nextFocusState.node = undefined;

    if (wsState.active) {
      wsState.connection.removeEventListener("message", onWsMessage);
    }
  });

  $effect(() => {
    if (wsState.active) {
      wsState.connection.removeEventListener("message", onWsMessage);
      wsState.connection.addEventListener("message", onWsMessage);
    }
  });

  /** @param {MessageEvent<string>} e*/
  const onWsMessage = (e) => {
    if (e.data === PONG) return;

    /** @type {WSMessage<MediaDTO>}*/
    const data = JSON.parse(e.data);

    if (data.data.id === id) {
      const topic = topicMap[data.topic];
      if (topic) {
        topic(data.data);
      }
    }
  };

  /** @type {WSTopicMap<MediaDTO>}*/
  const topicMap = {
    media_update: (updatedMedia) => {
      if (mediaEntity && updatedMedia.relations?.length > 0) {
        mediaEntity.relations = [
          ...(mediaEntity.relations || []),
          ...updatedMedia.relations,
        ];
      }
      if (mediaEntity && updatedMedia.relations?.length == 0) {
        mediaEntity.relations = [];
      }
    },
  };

  $effect(() => {
    if (nextFocusState.node === null && videoNode) {
      nextFocusState.node = videoNode;
    }
  });

  $effect(() => {
    if (videoNode) {
      videoNode.focus();
    }
  });

  /** @param {string} tagName
   * @returns {Promise<Item>}*/
  const createTagHandler = async (tagName) => {
    const createdTags = await createTags([tagName]);

    if (createdTags?.length > 0) {
      return createdTags[0];
    }

    throw Error("No tags returned after create");
  };

  /** @param {string} personName
   * @returns {Promise<Item>}
   */
  const createPersonHandler = async (personName) => {
    const createdPeople = await createPeople([personName]);

    if (createdPeople?.length > 0) {
      return createdPeople[0];
    }

    throw Error("No people returned after create");
  };

  /** @param {KeyboardEvent} e*/
  const handleOnKeyUp = (e) => {
    if (!videoNode) {
      return;
    }
    switch (e.code) {
      case "KeyL":
        videoNode.currentTime = videoNode.currentTime + 10;
        break;
      case "KeyJ":
        videoNode.currentTime = videoNode.currentTime - 10;
    }
  };

  /** @param {KeyboardEvent} e*/
  const handleOnKeyDown = (e) => {
    switch (e.code) {
      case "Escape":
        if (nextFocusState.node) {
          nextFocusState.node.focus();
          nextFocusState.node = videoNode;
        }
        break;
    }
  };

  const handleOnFocus = () => {
    if (!nextFocusState.node) {
      nextFocusState.node = videoNode;
    }
  };

  /** @type {number} */
  let progressTimeout;

  /** @param {Event} e*/
  const handleTimeUpdate = (e) => {
    clearTimeout(progressTimeout);
    progressTimeout = setTimeout(async () => {
      loadingProgress = true;
      try {
        if (videoNode && mediaEntity) {
          const prog = await updateProgress(id, videoNode.currentTime);

          mediaEntity.progress = prog.progress;
        }
      } finally {
        loadingProgress = false;
      }
    }, 1000);
  };

  const handleWatchedClick = async () => {
    let val = 0;
    if (watchedPercentage <= 0.9) {
      val = mediaEntity?.video?.runtime ?? 0;
    }

    loadingProgress = true;
    try {
      const prog = await updateProgress(id, val, true);
      if (mediaEntity) {
        mediaEntity.progress = prog.progress;
      }
    } finally {
      loadingProgress = false;
    }
  };

  /**
   * @param {Event} e
   * @param {String} updatedTitle
   */
  const handleTitleUpdate = async (e, updatedTitle) => {
    loadingTitle = true;
    try {
      const res = await updateMedia(id, { title: updatedTitle });

      if (mediaEntity) {
        mediaEntity.title = res.title ?? "";
      }
      editingTitle = false;
    } finally {
      loadingTitle = false;
    }
  };

  const handleFavouriteClick = async () => {
    loadingFavourite = true;

    try {
      if (mediaEntity) {
        if (mediaEntity.favourite) {
          await removeFavourite(id);
          mediaEntity.favourite = false;
        } else {
          await addFavourite(id);
          mediaEntity.favourite = true;
        }
      }
    } finally {
      loadingFavourite = false;
    }
  };

  /**
   * @param {Event} _e
   * @param {{metadata: ChapterMetadadataDTO}} chapter
   */
  const handleChapterClick = (_e, chapter) => {
    const newTime = chapter.metadata.timestamp;

    if (videoNode) {
      videoNode.currentTime = newTime;
    }
  };
</script>

{#if loadingMedia}
  <p>loading</p>
{:else if mediaEntity}
  <div class="container is-fluid">
    {#if !mediaEntity.exists || mediaEntity.deleted}
      <section
        class={`hero ${mediaEntity.exists ? "is-warning" : "is-danger"}`}
      >
        <div class="hero-body">
          {#if !mediaEntity.exists && !mediaEntity.deleted}
            <p class="title">File deleted from disk</p>
            <p class="subtitle">
              The file has been deleted from disk outside of Exorcist
            </p>
          {:else if !mediaEntity.exists}
            <p class="title">File does not exist</p>
            <p class="subtitle">
              Not much we can do here but show you the information that remains
            </p>
          {:else}
            <p class="title">File exists but has been deleted</p>
            <p class="subtitle">
              Soon you will be able to restore deleted files that still exist.
              <br />
              You can permanently delete the files on disk by going through the delete
              flow again.
            </p>
          {/if}
        </div>
      </section>
      <br />
    {/if}
    <!-- svelte-ignore a11y_media_has_caption -->
    {#if mediaEntity.exists}
      <video
        src={videoUrlById(id)}
        controls
        poster={imageUrlById(thumbnailRelation?.relatedToId ?? "")}
        bind:this={videoNode}
        onkeyup={handleOnKeyUp}
        onkeydown={handleOnKeyDown}
        onfocus={handleOnFocus}
        ontimeupdate={handleTimeUpdate}
      ></video>
    {/if}

    <div class="container">
      <div class="field has-addons">
        {#if !mediaEntity.deleted && mediaEntity.exists}
          <p class="control">
            <Link class="button" to={routes.playlistAdd + `?media=${id}`}>
              <span class="icon">
                <i class="fas fa-plus"></i>
              </span>
            </Link>
          </p>

          <p class="control">
            <button
              class={`button ${loadingFavourite ? "is-loading" : ""}`}
              onclick={handleFavouriteClick}
              aria-label="toggle favourite"
              disabled={loadingFavourite}
            >
              <span class="icon">
                <i
                  class={`${mediaEntity.favourite ? "fas fa-heart" : "fa-regular fa-heart"}`}
                ></i>
              </span>
            </button>
          </p>
          <p class="control">
            <Link
              class={`button`}
              aria-label="refresh metadata"
              to={routes.refreshMetadataFn(id, routes.videoFunc(id))}
            >
              <span class="icon"><i class="fas fa-arrows-rotate"></i></span
              ></Link
            >
          </p>

          <p class="control">
            <Link
              class="button"
              aria-label="convert media"
              to={routes.convertFn(id)}
              on:click={() => (pageState.media = mediaEntity)}
              ><span class="icon"><i class="fas fa-arrows-spin"></i></span
              ></Link
            >
          </p>
          <p class="control">
            <Link
              class="button"
              aria-label="generate thumbnail"
              to={routes.generateThumbnailFn(id)}
            >
              <span class="icon"><i class="fas fa-image"></i></span>
            </Link>
          </p>
          <p class="control">
            <Link
              class="button"
              aria-label="generate chapters"
              to={routes.generateChaptersFn(id)}
            >
              <span class="icon"><i class="fas fa-images"></i></span>
            </Link>
          </p>
        {/if}
        <p class="control">
          <button
            class={`button ${loadingProgress ? "is-loading" : ""}`}
            onclick={handleWatchedClick}
            aria-label="toggle watched"
            disabled={loadingProgress}
          >
            <span class="icon">
              <i
                class={`fas ${watchedPercentage > 0.9 ? "fa-eye-slash" : "fa-eye"}`}
              ></i>
            </span>
          </button>
        </p>
        {#if !mediaEntity.deleted || mediaEntity.exists}
          <p class="control">
            <Link
              class="button"
              to={routes.delete.mediaFunc(id, mediaEntity.title)}
            >
              <span class="icon">
                <i class="fas fa-trash"></i>
              </span>
            </Link>
          </p>
        {/if}
      </div>
      {#if editingTitle}
        <EditHeading
          value={mediaEntity.title}
          oncancel={() => (editingTitle = false)}
          onsave={handleTitleUpdate}
          loading={loadingTitle}
        />
      {:else}
        <h1 class="title is-1">
          {mediaEntity.title}
          <HeaderIconButton
            icon={`fas fa-pen`}
            iconClass={editingTitle ? "has-text-info" : ""}
            ariaLabel="edit title toggle"
            onclick={() => {
              editingTitle = !editingTitle;
            }}
          />
        </h1>
      {/if}
    </div>
    <br />
    <div class="container">
      <Items
        title="Tags"
        items={mediaEntity.tags}
        fetch={getAllTags}
        remove={async (tagId) => removeTag(id, tagId)}
        add={async (tagId) => addTag(id, tagId)}
        create={createTagHandler}
        urlFn={routes.tagFunc}
        disableEdit={mediaEntity.deleted}
      />
    </div>
    <br />
    <div class="container">
      <Items
        title="People"
        items={mediaEntity.people}
        fetch={getAllPeople}
        remove={async (personId) => removePerson(id, personId)}
        add={async (personId) => addPerson(id, personId)}
        create={createPersonHandler}
        urlFn={routes.personFunc}
        disableEdit={mediaEntity.deleted}
      />
    </div>

    {#if mediaEntity.relations?.find((r) => r.relationType === "chapter")}
      <br />
      <div class="container">
        <Chapters
          chapters={mediaEntity.relations
            ?.filter((relation) => relation.relationType === "chapter")
            .sort((a, b) => a.metadata.timestamp - b.metadata.timestamp)}
          onclick={handleChapterClick}
        />
      </div>
    {/if}

    <br />
    <div class="container">
      <Relations
        {id}
        relations={mediaEntity.relations?.filter(
          (r) => r.relationType === "media",
        )}
      />
    </div>

    <br />
    <div class="container">
      <table class="table">
        <thead>
          <tr>
            <th>Key</th>
            <th>Value</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>Dimensions</td>
            <td>{mediaEntity?.video?.width}x{mediaEntity?.video?.height}</td>
          </tr>
          <tr>
            <td>Runtime</td>
            <td>{formatRuntime(mediaEntity?.video?.runtime ?? 0)}</td>
          </tr>
          <tr>
            <td>Size</td>
            <td>{formatFileSize(mediaEntity.size)}</td>
          </tr>
          <tr>
            <td>Path</td>
            <td>{mediaEntity.path}</td>
          </tr>
          <tr>
            <td>Added</td>
            <td>{mediaEntity.added}</td>
          </tr>
          <tr>
            <td>Created</td>
            <td>{mediaEntity.created}</td>
          </tr>
          <tr>
            <td>Modified</td>
            <td>{mediaEntity.modified}</td>
          </tr>
          <tr>
            <td>Checksum</td>
            <td>{mediaEntity.checksum}</td>
          </tr>
          <tr>
            <td>Deleted</td>
            <td>{mediaEntity.deleted}</td>
          </tr>
          <tr>
            <td>Exists</td>
            <td>{mediaEntity.exists}</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
{/if}

<style>
  video {
    height: 100%;
    width: 100%;
    max-height: 90vh;
  }
</style>
