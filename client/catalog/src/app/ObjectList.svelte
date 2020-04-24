<script>
import { onMount } from 'svelte';
import client from '../data/client.js';
let objects = [];
let load = client.model.findMany({
	include: {
		attributes: true,
		relationships: true,
	}
});
let cap = (s) => { 
	if (!s) {
		return "";
	}
	return s.charAt(0).toUpperCase() + s.slice(1);
};

import { breadcrumbStore } from './breadcrumbStore.js';
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}]
);
</script>

<style>
	.box {
		display: flex;
	}
	.stuff {

	}
	a.object-box {
		display: flex;
		flex-direction: column;
		color: inherit;
		width: 150px;
		padding: 1em 1.5em;
	}
	a.object-box:hover {
		background: var(--background-highlight);
	}
	a.object-box.center {
		align-items: center;
		justify-content: center;
	}

	.spacer {
		width: 0;
	}
	.obj-title{
		font-weight: 600;
	}

</style>

<div class="box">
	{#await load}
	{:then models}
		{#each models as model}
			<a href="/object/{model.id}" class="object-box">
				<div class="obj-title">{cap(model.name)}</div>
				{#each model.attributes as attr}
					<div>{attr.name}</div>
				{/each}
			</a>
			<div class="spacer"/>
		{/each}
		<a href="/objects/new" class="object-box center">
			<div>+ Add</div>
		</a>
		<div class="spacer"/>
	{:catch error}
		<div>Error..</div>
	{/await}
</div>
