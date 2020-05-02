<script>
import client from '../data/client.js';
import HLTable from '../ui/HLTable.svelte';
import HLRow from '../ui/HLRow.svelte';
import { breadcrumbStore } from './breadcrumbStore.js';
breadcrumbStore.set(
	[{
		href: "/log",
		text: "Log",
	}]
);
let load = client.log.scan({
	count: 10,
	offset: 0,
});

function trunc(s) {
	return s
}
</script>

<style>
.box {
	margin:  1.5em; 
}
.op-closed {
	white-space: nowrap;
	overflow: hidden;
	max-height: 4em;
	text-overflow: ellipsis;
}

</style>



<div class="box">
	<h1>Log</h1>
	<HLTable>
	{#await load}
		<HLRow>
		</HLRow>
	{:then entries}
		{#each entries as entry}
		<HLRow>
			<div class="op-closed">{JSON.stringify(entry, null, 2)}</div>
		</HLRow>
		{/each}
	{:catch error}
		<div>Error..</div>
	{/await}
	</HLTable>
</div>
