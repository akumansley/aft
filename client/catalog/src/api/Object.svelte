<script>
	import { onMount } from 'svelte';
	import { clone } from "../lib/util.js";

	export let init = null;
	let state = {
		value: {},
	}
	export let op = null;

	export const object = new Proxy({}, {
		get: function(target, property) {
			return state.value[property];
		},
		set: function(target, property, newVal) {
			state.value[property] = newVal;

			if (init && newVal === init[property]) {
				delete op[property];
				op = op;
			} else {
				op[property] = newVal;
			}
			return true;
		},
	});

	onMount(() => {
		state.value = clone(init);
	});


</script>

<slot />