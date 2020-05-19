<script>
export let params;
import client from '../../data/client.js';
import Model from './Model.svelte';
import { breadcrumbStore } from '../breadcrumbStore.js';

let id = params.id;
let load = client.model.findOne({where: {id: id}, include: {rightRelationships: true, leftRelationships: true, attributes: true}});

let cap= (s) => { 
	if (!s) {
		return "";
	}
	return s.charAt(0).toUpperCase() + s.slice(1)
};

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
</style>

<div class="box">
	{#await load}
		&nbsp;
	{:then model}
		<Model model={model}/>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>

