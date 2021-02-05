<script>
	import {HLSectionTitle} from '../../ui/page/page.js';
	import {cap} from '../../lib/util.js';
	import ChipCheckbox from './ChipCheckbox.svelte';
	import {onMount} from 'svelte';

	export let options = null;
	export let value = null;
	export let key = null;


	const state = {};
	const modState = {};
	let initialized = false;
	$: if (options && options.length && value && !initialized) {
		initialized = true;
		init();
	}

	function init() {
		// initialize state
		for (let m of options) {
			for (let o of m[key]) {
				state[o.id] = false;
			}
		}
		for (let v of value) {
			state[v.id] = true;
		}
		for (let m of options) {
			updateModState(m);
		}
		updateValue();
	}

	function updateModState(mod) {
		let all = true;
		for (let o of mod[key]) {
			if (!state[o.id]) {
				all = false;
				break
			}
		}
		modState[mod.id] = all;
	}

	function handleChange(mod, opt) {
		return (e) => {
			let newVal = e.target.checked;
			state[opt.id] = newVal;
			updateModState(mod);
			updateValue();
		}
	}

	function toggleMod(mod){
		return (e) => {
			let newVal = !modState[mod.id];
			modState[mod.id] = newVal;
			for (let opt of mod[key]) {
				state[opt.id] = newVal;
			}
			updateValue();
		}
	}

	function updateValue() {
		for (let [id, selected] of Object.entries(state)) {
			if (selected){
				includeValue(id);
			} else {
				excludeValue(id);
			}
		}
	}

	function excludeValue(id) {
		const found = value.some(v => v && v.id === id)
		if (found)  {
			value = value.removeBy(v => v && v.id === id)
		}
	}

	function includeValue(id) {
		const found = value.some(v => v && v.id === id)
		if (!found) {
			value = value.add({id})
		}
	}

</script>

<style>
	.grid {
		display: flex;
		flex-direction: row;
		flex-wrap: wrap;
	}
	.wrapper {
		margin-right: .5em;
		margin-bottom: 1em;
	}
	.link {
		color: var(--text-color-highlight);
		cursor: pointer;
	}
</style>

{#each options as mod}
<h2>{cap(mod.name)}
	-
	{#if modState[mod.id]}
	<span class="link" on:click={toggleMod(mod)}>Deselect All</span>
	{:else}
	<span class="link" on:click={toggleMod(mod)}>Select All</span>
	{/if}
</h2>
<div class="grid">
	{#each mod[key] as option}
	<div class="wrapper">
		<ChipCheckbox bind:checked={state[option.id]} on:change={handleChange(mod, option)}>{option.name}</ChipCheckbox>
	</div>
	{/each}
</div>
{/each}
