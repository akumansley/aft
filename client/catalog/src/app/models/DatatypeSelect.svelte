<script context="module">
	import client from '../../data/client.js';
	let load = client.api.datatype.findMany({where:{}});
</script>

<script>
	import { getContext } from 'svelte';
	import { key } from '../../api/api.js';

	export let value;
	let datatype = {};
	load.then((dts) => {
		datatype = dts[0];
	});

	import HLSelect from '../../ui/form/HLSelect.svelte';
	import {restrictToIdent, cap, isObject} from '../util.js';

	let operation = getContext(key);

	$: {
		value = {
			connect: {id: datatype.id},
		}
	}

</script>

{#await load then datatypes}
<HLSelect bind:value={datatype}>
	{#each Object.entries(datatypes) as dt}
	<option value={dt[1]}>
		{cap(dt[1].name)}
	</option>
	{/each}
</HLSelect>
{/await}
