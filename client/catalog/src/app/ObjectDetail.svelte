<script>
	export let params;
	import { onMount } from 'svelte';
	import client from '../data/client.js';
	let id = params.id;
	let load = client.objects.info({id: id});
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
		<div>{object.name}</div>
		{#each object.fields as field}
			<div>{field.name}</div>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
