<script>
	import { navStore } from '../stores.js';
	import client from '../../data/client.js';
	import {HLSelect} from '../../ui/form/form.js';
	import TxEntry from './TxEntry.svelte';

	navStore.set("log");

	let logName = "db";


	let entries = [];
	$: load = client.rpc.scan({args: {
		log: logName,
		count: 100,
		offset: 0,
	}}).then(result => {
		entries = result;
	});


</script>

<style>
	.frame {
		max-width: 50em;
		margin: 1em auto 0;
	}
	.logtable {
		overflow: auto;
	}

</style>

<div class="frame">
	View logs for:
	<HLSelect bind:value={logName}>
		<option value="db">Database</option>
		<option value="request">Requests</option>
	</HLSelect>
</div>

<div class="logtable">
	{#await load then _}

	{#if logName === "db"}
	{#each entries as entry}
	<TxEntry value={entry}/>
	{/each}
	{:else if logName === "request"}
	{#each entries as entry}
	<div>{JSON.stringify(entry)}</div>
	{/each}
	{/if}

	{/await}
</div>

