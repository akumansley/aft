<script>
	import user from './user.js';
	import { Link } from "svelte-routing";
	
</script>

<Link to="/list/{$user.id}">
		Make your gift list &rarr;
</Link>

