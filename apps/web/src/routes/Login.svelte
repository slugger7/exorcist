<script>
    import { navigate } from "svelte-routing";
    import { userState } from "../lib/state/userState.svelte.js";
    import { login } from "../lib/controllers/authentication.js";
    import routes from "./routes.js";
    import { setupWebsocket } from "../lib/controllers/websocket.js";

    let username = $state(import.meta.env.DEV ? "admin" : "");
    let password = $state(import.meta.env.DEV ? "admin" : "");
    let loading = $state(false);
    const returnUrl = new URLSearchParams(window.location.search).get(
        "returnUrl",
    );

    async function handleSubmit(e) {
        e.preventDefault();
        loading = true;

        try {
            const user = await login(username, password);
            userState.id = user.userId;
            userState.username = user.username;

            setupWebsocket();

            if (returnUrl === null) {
                navigate(routes.home);
            } else {
                navigate(returnUrl);
            }
        } catch (e) {
            console.error(e);
        } finally {
            loading = false;
        }
    }
</script>

<div class="container">
    <h1 class="title is-1">Login</h1>
    <form onsubmit={handleSubmit}>
        <div class="field">
            <label class="label" for="username">Username</label>
            <div class="control">
                <input
                    class="input"
                    bind:value={username}
                    name="username"
                    type="text"
                    placeholder="Username"
                />
            </div>
        </div>

        <div class="field">
            <label class="label" for="password">Password</label>
            <div class="control">
                <input
                    class="input"
                    bind:value={password}
                    name="password"
                    type="password"
                    placeholder="Password"
                />
            </div>
        </div>

        <div class="field">
            <div class="control">
                <button class="button is-primary {loading ? 'is-loading' : ''}"
                    >Login</button
                >
            </div>
        </div>
    </form>
</div>
