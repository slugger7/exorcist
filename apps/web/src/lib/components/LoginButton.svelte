<script>
import { Link } from 'svelte-routing'
import routes from '../../routes/routes.js'
</script>


<div class="buttons">
    <Link class="button is-primary" to={routes.login}>
        Log in
    </Link>
</div>

