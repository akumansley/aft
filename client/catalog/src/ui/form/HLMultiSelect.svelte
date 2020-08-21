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

</script>

<style>

</style>
{#each options as option}
	<label>
	<input type=checkbox bind:checked={state[option.id]} on:change={handleChange}>
{option.name}
	</label>
{/each}
