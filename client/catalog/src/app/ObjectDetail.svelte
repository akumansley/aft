<script>
export let params;
import { onMount } from 'svelte';
import client from '../data/client.js';
let id = params.id;
let load = client.model.findOne({where: {id: id}});
</script>

<style>
	.box {
		padding-left: 1em;
	}
	.stuff {
		margin-top: 1.25em;

	}

</style>

<div class="box">
	<a href="/objects">&larr; back</a>
	{#await load}
		Loading..
	{:then object}
		<div>{object.Name}</div>
		{#each Object.entries(object.Attributes) as entry}
			<div>{entry[0]}</div>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
