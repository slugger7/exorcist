<script>
  /**
   * @import { MediaOverviewDTO } from "../dto"
   * @import { WSTopicMapView } from "../lib/types"
   */
  import { getAllWithParams, ordinals } from "../lib/controllers/media";
  import routes from "./routes";
  import MediaView from "../lib/components/MediaView.svelte";

  /** @type {WSTopicMapView<MediaOverviewDTO>}*/
  const topics = {
    media_create: (mediaPage, newItems, media) => {
      newItems.push(media);

      mediaPage.total++;

      return [mediaPage, newItems];
    },
    media_overview_update: (mediaPage, newItems, updatedMedia) => {
      newItems = newItems.map((v) => {
        if (v.id === updatedMedia.id) {
          return { ...v, ...updatedMedia };
        }
        return v;
      });

      if (mediaPage) {
        mediaPage.data = mediaPage.data.map((v) => {
          if (v.id === updatedMedia.id) {
            return {
              ...v,
              ...updatedMedia,
            };
          }
          return v;
        });
      }

      return [mediaPage, newItems];
    },
  };
</script>

<MediaView
  title="Home"
  fetchFn={getAllWithParams}
  route={routes.home}
  {ordinals}
  topicMap={topics}
/>
