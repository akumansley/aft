<script>
	import { onMount } from 'svelte';
	import client from '../data/client.js';
	let objects = [];
	let load = client.model.findMany();
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
		background: #0d0a10;
	}
	a.object-box:hover {
		background: #130f17;
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
