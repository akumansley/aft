<script>
export let params;
import { onMount } from 'svelte';
import client from '../data/client.js';
import { breadcrumbStore } from './breadcrumbStore.js';
breadcrumbStore.set(
	[{
		href: "/objects",
		text: "Objects",
	}, {
		href: "/objects/new",
		text: "New",
	}]
);
</script>

<style>
	.box {
		margin: 1em 1.5em; 
	}
	.v-space{
		height: .5em;
	}
	input.h1 {
		font-size: var(--scale-3);
		font-weight: 600;
		background:inherit;
		border-color: var(--border-color);
		color: inherit;
		border-radius: 0;
	}
	input.h1::placeholder {
		color: var(--text-color-darker);
	}
	h2 {
		font-size: var(--scale--1);
		font-weight: 500;
		line-height: 1;
	}

</style>

<div class="box">
	<input class="h1" placeholder="Object name.." type="text"/>
	<h2>Attributes</h2>
	<h2>Relationships</h2>
</div>
