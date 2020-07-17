<script>
export let storedAs;
export let storage;
export let change = null;
import { cap } from '../util.js';
import HLSelect from '../../ui/form/HLSelect.svelte';
</script>

<style>
	.spacer {
		width: .5em;
		display: inline-block;
	}
</style>

Stored as: <span class="spacer"/>
<HLSelect bind:value={storedAs} on:change={change}>
	{#each Object.entries(storage) as it, ix}
	<option value={it[1]["id"]}>
		{cap(it[1]["name"])}
	</option>
	{/each}
</HLSelect>