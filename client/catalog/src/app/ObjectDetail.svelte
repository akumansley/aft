<script>
export let params;
import { onMount } from 'svelte';
import client from '../data/client.js';
import { breadcrumbStore } from './breadcrumbStore.js';
	let cap= (s) => { 
		return s.charAt(0).toUpperCase() + s.slice(1)
	};

let id = params.id;
let load = client.model.findOne({where: {id: id}});

load.then(obj => {
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}, {
		href: "/object/" + id,
		text: cap(obj.Name),
	}]
);
});
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
	{#await load}
		Loading..
	{:then object}
		<div>{cap(object.Name)}</div>
		{#each Object.entries(object.Attributes) as entry}
			<div>{entry[0]}</div>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
