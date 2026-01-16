<script>
    import { Route } from "svelte-routing";
    import routes from "./routes.js";
    import Home from "./Home.svelte";
    import People from "./People.svelte";
    import Login from "./Login.svelte";
    import Logout from "./Logout.svelte";
    import Libraries from "./Libraries.svelte";
    import CreateLibrary from "./create/Library.svelte";
    import CreateLibraryPath from "./create/LibraryPath.svelte";
    import Library from "./Library.svelte";
    import LibraryPath from "./LibraryPath.svelte";
    import EditLibrary from "./edit/Library.svelte";
    import Video from "./Video.svelte";
    import Jobs from "./Jobs.svelte";
    import Person from "./Person.svelte";
    import Tag from "./Tag.svelte";
    import Tags from "./Tags.svelte";
    import User from "./User.svelte";
    import DeleteMedia from "./delete/Media.svelte";
    import Playlists from "./Playlists.svelte";
    import Playlist from "./Playlist.svelte";
    import CreatePlaylist from "./create/Playlist.svelte";
    import PlaylistAdd from "./PlaylistAdd.svelte";
    import RefreshMetadata from "./RefreshMetadata.svelte";
    import GenerateChaptersJob from "./GenerateChaptersJob.svelte";
    import VideoAlt from "./VideoAlt.svelte";
    import GenerateThumbnailJob from "./GenerateThumbnailJob.svelte";
</script>

<main>
    <Route path={routes.login}><Login /></Route>
    <Route path={routes.logout}><Logout /></Route>
    <Route path={routes.people}><People /></Route>
    <Route path={routes.person} let:params
        ><Person id={params.id} name={params.name} /></Route
    >
    <Route path={routes.tags}><Tags /></Route>
    <Route path={routes.tag} let:params
        ><Tag id={params.id} name={params.name} /></Route
    >
    <Route path={routes.libraries}><Libraries /></Route>
    <Route path={routes.create.library}><CreateLibrary /></Route>
    <Route path={routes.edit.library} let:params
        ><EditLibrary id={params.libraryId} /></Route
    >
    <Route path={routes.library} let:params
        ><Library id={params.id} name={params.name} /></Route
    >
    <Route path={routes.create.libraryPath} let:params
        ><CreateLibraryPath libraryId={params.libraryId} /></Route
    >
    <Route path={routes.libraryPath} let:params
        ><LibraryPath id={params.id} name={params.name} /></Route
    >
    <Route path={routes.user} let:params
        ><User id={params.id} name={params.name} /></Route
    >
    <Route path={routes.jobs}><Jobs /></Route>
    <Route path={routes.home}><Home /></Route>
    <Route path={routes.delete.media} let:params
        ><DeleteMedia id={params.mediaId} /></Route
    >
    <Route path={routes.video} let:params><Video id={params.id} /></Route>
    <Route path={routes.videoAlt} let:params><VideoAlt id={params.id} /></Route>
    <Route path={routes.playlists}><Playlists /></Route>
    <Route path={routes.playlist} let:params
        ><Playlist id={params.id} name={params.name} /></Route
    >
    <Route path={routes.create.playlist}><CreatePlaylist /></Route>
    <Route path={routes.playlistAdd}><PlaylistAdd /></Route>
    <Route path={routes.refreshMetadata} let:params
        ><RefreshMetadata
            mediaId={params.id}
            redirect={params.redirect}
        /></Route
    >
    <Route path={routes.refreshLibraryMetadata} let:params
        ><RefreshMetadata
            libraryId={params.id}
            redirect={params.redirect}
        /></Route
    >
    <Route path={routes.generateChapters} let:params
        ><GenerateChaptersJob
            mediaId={params.id}
            redirect={params.redirect}
        /></Route
    >
    <Route path={routes.generateThumbnail} let:params
        ><GenerateThumbnailJob
            mediaId={params.id}
            redirect={params.redirect}
        /></Route
    >
</main>
