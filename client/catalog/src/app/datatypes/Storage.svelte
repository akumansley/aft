<script>
export let storedAs;
export let storage;
export let change = null;
import { cap } from '../../lib/util.js';
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
	{#each storage.enumValues as ev}
	<option value={ev.id}>
		{cap(ev.name)}
	</option>
	{/each}
</HLSelect>
