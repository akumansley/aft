<script>
export let params;
import { onMount } from 'svelte';
import HLTable from '../ui/HLTable.svelte';
import HLRow from '../ui/HLRow.svelte';
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
		font-size: var(--scale-3);
		font-weight: 600;
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
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
		<HLTable>
		{#each model.attributes as attr}
			<HLRow>
				{attr.name}
				<div class="v-space"/>
				<dl>
					<dt>Type</dt>
				<dd>
					{AttrType[attr.attrType]}
				</dd>
				</dl>
			</HLRow>
		{/each}
		</HLTable>
		<div class="v-space"/>
		<h2>Relationships</h2>
		<HLTable>
		{#each model.relationships as rel}
			<HLRow>
			{rel.name}
			</HLRow>
		{/each}
		</HLTable>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>

