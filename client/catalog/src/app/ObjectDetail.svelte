<script>
export let params;
import { onMount } from 'svelte';
import client from '../data/client.js';
import { breadcrumbStore } from './breadcrumbStore.js';
	let cap= (s) => { 
		if (!s) {
			return "";
		}
		return s.charAt(0).toUpperCase() + s.slice(1)
	};

let id = params.id;
let load = client.model.findOne({where: {id: id}, include: {relationships: true, attributes: true}});

load.then(obj => {
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}, {
		href: "/object/" + id,
		text: cap(obj.name),
	}]
);
});
</script>

<style>
	.box {
		margin: 1em 1.5em; 
	}
	.model-name {
		font-size: 1.728em;
		font-weight: 600;
	}


</style>

<div class="box">
	{#await load}
	{:then model}
		<div class="model-name">{cap(model.name)}</div>
		{#each model.attributes as attr}
			<div>{attr.name}</div>
		{/each}
		{#each model.relationships as rel}
			<div>&rarr; {rel.name}</div>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
