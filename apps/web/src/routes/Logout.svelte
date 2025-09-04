<script>
    import { onMount } from "svelte";
    import { navigate } from "svelte-routing";
    import { userState } from "../lib/state/userState.svelte.js";
    import { logout } from "../lib/controllers/authentication.js";
    import routes from "./routes.js";

    onMount(async () => {
        try {
            await logout();
        } catch (e) {
            console.error(e);
        }

        localStorage.clear();
        userState.id = null;
        userState.username = null;
        navigate(routes.login);
    });
</script>

<h1 class="title is-1">Logout</h1>
