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
		padding-left: 1em;
	}
	.stuff {
		margin-top: 1.25em;
	}

</style>

<div class="box">
</div>
