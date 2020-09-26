<script>
export let options = [];
export let selected = [];

const state = {};

function handleChange() {
	let newVal = [];
	for (let v of options) {
		if (state[v.id]) {
			newVal.push(v)
		}
	}
	selected = newVal;
}

import HLCheckbox from './HLCheckbox.svelte';

</script>

<style>

</style>
{#each options as option}
	<HLCheckbox bind:checked={state[option.id]} on:change{handleChange}>{option.name}</HLCheckbox>
{/each}
