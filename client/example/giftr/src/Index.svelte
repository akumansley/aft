<script>
	import user from './user.js';
	import client from './client.js';
	import { Link } from "svelte-routing";

	const load = client.api.user.findMany({});
</script>
<style>
	h2 {
		font-size: var(--scale-0);
		margin: 0;
	}
	.box {
		background: var(--color-1-darker);
		padding: .5em 1em;
		border: 2px solid var(--color-3);
		border-radius: 5px;

	}
	.box :global(a) {

	}
	.spacer {
		height: 1em;
	}
</style>

<div class="box">
	<Link to="/list/{$user.id}">
			My gift list &rarr;
	</Link>
</div>
<div class="spacer"></div>

<h2>Lists</h2>

{#await load then users}
{#each users as u}
{#if u.id !== $user.id}
<div>
<Link to="/list/{u.id}">
		{u.email}
</Link>
</div>
{/if}
{/each}
{/await}
