<script>
	import { navStore } from '../stores.js';
	import client from '../../data/client.js';
	import {HLSelect} from '../../ui/form/form.js';
	import {HLBorder} from '../../ui/page/page.js';
	import TxEntry from './TxEntry.svelte';

	navStore.set("log");

	let logName = "db";


	let entries = [];
	$: load = client.rpc.scan({
		log: logName,
		count: 100,
		offset: 0,
	}).then(result => {
		entries = result;
	});


</script>

<style>
	.frame {
		max-width: 50em;
		margin: 1em 1em;
	}
	.logtable {
		overflow: auto;
	}
	.entry {
		padding: 1em;
	}

</style>

<div class="frame">
	View logs for:
	<HLSelect bind:value={logName}>
		<option value="db">Database</option>
		<option value="request">Requests</option>
	</HLSelect>
</div>
<HLBorder/>

<div class="logtable">
	{#await load then _}

	{#if logName === "db"}
	{#each entries as entry}
	<TxEntry value={entry}/>
	{/each}
	{:else if logName === "request"}
	{#each entries as entry}
	<div class="entry">{JSON.stringify(entry)}</div>
	{/each}
	{/if}

	{/await}
</div>

