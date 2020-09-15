<script>
	import { onMount } from 'svelte';
	import { nonEmpty, clone } from '../lib/util.js';

	export let component = null;
	export let init = null;
	export let op = [];
	 onMount(() => {
	 	if (init) {
	 		values = clone(init);
	 	}
	 });

	let values = [];

	export const add = () => {
		init = [...init, null];
	}

	$: {
		op = {};
		for (let value of values) {
			if (!value) {
				continue;
			}
			if (value.create) {
				if (op.create) {
					op.create = [...op.create, value.create];
				} else {
					op.create = [value.create];
				}
			} else if (value.update) {
				if (op.update) {
					op.update = [...op.update, value.update];
				} else {
					op.update = [value.update];
				}
			} else if (value.delete) {
				if (op.delete) {
					op.delete = [...op.delete, value.delete];
				} else {
					op.delete = [value.delete];
				}
			}
		}
	}

</script>


{#each init as i, ix}
<svelte:component this={component} init={i} bind:op={values[ix]}/>
{/each}

