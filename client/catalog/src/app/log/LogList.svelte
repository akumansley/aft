<script>
	export let params = null;
	import {router} from '../router.js'
	import { navStore } from '../stores.js';
	import client from '../../data/client.js';
	import {HLSelect} from '../../ui/form/form.js';
	import {HLBorder} from '../../ui/page/page.js';
	import TxEntry from './TxEntry.svelte';

	navStore.set("log");

	let logName = params.source;
	let loadedLogName = null;

	function navigate() {
		router.route("/log/" + logName);
		loadData();
	}

	function loadData() {
		return client.rpc.scan({
			log: logName,
			count: 100,
			offset: 0,
		}).then(result => {
			entries = result;
			loadedLogName = logName;
		});

	}
	let entries = [];
	loadData()

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
	<HLSelect bind:value={logName} on:change={navigate}>
		<option value="request">Requests</option>
		<option value="db">Database</option>
	</HLSelect>
</div>
<HLBorder/>

<div class="logtable">

	{#if loadedLogName === "db"}
	{#each entries as entry}
	<TxEntry value={entry}/>
	{/each}
	{:else if loadedLogName === "request"}
	{#each entries as entry}
	<div class="entry">{JSON.stringify(entry)}</div>
	{/each}
	{/if}

</div>

