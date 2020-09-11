<script>
	import { onMount } from 'svelte';
	import { nonEmpty, clone } from '../lib/util.js';

	export let component = null;
	export let init = null;
	export let op = [];
	let added = [];
	let existing = [];

	export const add = () => {
		added = [...added, null];
	}

	$: {
		op = {
			create: added.filter(nonEmpty),
			update: existing.filter(nonEmpty),
		}
	}
</script>


{#each init as i, ix}
	<svelte:component this={component} init={i} bind:op={existing[ix]}/>
{/each}

{#each added as i, ix}
	<svelte:component this={component} bind:op={added[ix]}/>
{/each}

<slot/>