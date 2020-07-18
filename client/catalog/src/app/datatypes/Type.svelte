<script>
export let types;
export let type;
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

Type: <span class="spacer"/>
<HLSelect bind:value={type} on:change={change(type)}>
	{#each Object.entries(types) as it, ix}
	<option value={it[1]}>
		{cap(it[1])}
	</option>
	{/each}
</HLSelect>