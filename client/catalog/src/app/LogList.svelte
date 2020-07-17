<script>
import { navStore } from './stores.js';
import client from '../data/client.js';
import HLRow from '../ui/list/HLRow.svelte';
 
let load = client.log({
	count: 100,
	offset: 0
});
navStore.set("log");

function trunc(s) {
	return s
}
</script>

<style>
.op-closed {
	white-space: nowrap;
	overflow: hidden;
	max-height: 4em;
	text-overflow: ellipsis;
}

</style>


<h1>Log</h1>
{#await load then entries}
	{#each entries as entry}
	<HLRow>
		<div class="op-closed">{JSON.stringify(entry, null, 2)}</div>
	</HLRow>
	{/each}
{/await}

