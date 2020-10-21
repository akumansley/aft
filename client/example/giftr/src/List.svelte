<script>
	import client from './client.js';
	import user from './user.js';
	export let userId;
	import { Link } from "svelte-routing";
	
</script>

<Link to="/">&larr; Back</Link>
<div class="box">
	You are {$user.id}; this is the list of {userId}.
</div>

