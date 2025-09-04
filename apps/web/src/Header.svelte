<script>
    import { Link } from "svelte-routing";
    import routes from "./routes/routes.js";
    import { userState } from "./lib/state/userState.svelte.js";

    let isMenuActive = $state(false);
    let userDropdownActive = $state(false);

    const onclick = () => {
        isMenuActive = false;
        userDropdownActive = false;
    };
</script>

<header>
    <nav class="navbar" aria-label="main navigation">
        {#if userState.id}
            <div class="navbar-brand">
                <Link class="navbar-item" to="/" {onclick}>Exorcist</Link>
                <button
                    class="navbar-burger {isMenuActive ? 'is-active' : ''}"
                    aria-label="menu"
                    aria-expanded={isMenuActive}
                    data-target="exorcistNavbar"
                    onclick={() => (isMenuActive = !isMenuActive)}
                >
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                </button>
            </div>

            <div
                id="exorcistNavbar"
                class="navbar-menu {isMenuActive ? 'is-active' : ''}"
            >
                <Link class="navbar-item" to={routes.libraries} {onclick}
                    >Libraries</Link
                >
                <Link class="navbar-item" to={routes.people} {onclick}
                    >People</Link
                >
                <Link class="navbar-item" to={routes.tags} {onclick}>Tags</Link>
                <Link class="navbar-item" to={routes.playlists} {onclick}
                    >Playlists</Link
                >
                <Link class="navbar-item" to={routes.jobs} {onclick}>Jobs</Link>
                <div class="navbar-end">
                    <div
                        class="navbar-item has-dropdown {userDropdownActive
                            ? 'is-active'
                            : ''}"
                    >
                        <button
                            class="navbar-link"
                            onclick={() =>
                                (userDropdownActive = !userDropdownActive)}
                            >{userState.username}</button
                        >
                        <div class="navbar-dropdown">
                            <Link
                                class="navbar-item"
                                to={routes.userFunc(
                                    userState.id,
                                    userState.username,
                                )}
                                {onclick}>Profile</Link
                            >
                            <hr class="navbar-divider" />
                            <Link class="navbar-item" to={routes.logout}
                                >Logout</Link
                            >
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    </nav>
</header>
