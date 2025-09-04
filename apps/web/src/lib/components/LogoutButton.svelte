<script>
import { Link } from 'svelte-routing'
import routes from '../../routes/routes.js'
</script>

<div class="buttons">
    <Link class="button" to={routes.logout}>
        Logout
    </Link>
</div>

