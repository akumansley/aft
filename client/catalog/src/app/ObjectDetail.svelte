<script>
export let params;
import { onMount } from 'svelte';
import client from '../data/client.js';
import { AttrType } from '../data/enums.js';
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

	h1 {
		font-size: 1.728em;
		font-weight: 600;
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}
	.hl-row {
		border-bottom: 1px solid var(--border-color);
		padding-left: .5em;
		padding-top: .25em;
		padding-bottom: .25em;
	}
	.hl-row:hover {
		background: var(--background-highlight);
	}
	.hl-table {
		border: 1px solid var(--border-color);
		max-width: 30em;
	}
	.hl-row:last-child {
		border-bottom: none;
	}
	.v-space{
		height: .5em;
	}
	dl {
		padding: 0; 
		margin:0;
	}
	dt {
		font-size: var(--scale--2);
		text-transform: uppercase;
		font-weight: 600;
	}
	dd {
		margin: 0;
	}


</style>

<div class="box">
	{#await load}
		&nbsp;
	{:then model}
		<h1>{cap(model.name)}</h1>

		<h2>Attributes</h2>
		<div class="hl-table">
		{#each model.attributes as attr}
			<div class="hl-row">
				{attr.name}
				<div class="v-space"/>
				<dl>
					<dt>Type</dt>
				<dd>
					{AttrType[attr.attrType]}
				</dd>
				</dl>
			</div>
		{/each}
		</div>
		<div class="v-space"/>
		<h2>Relationships</h2>
		<div class="hl-table">
		{#each model.relationships as rel}
			<div class="hl-row">{rel.name}</div>
		{/each}
		</div>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>

