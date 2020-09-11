<script context="module">
	import client from '../../data/client.js';
	let load = client.api.datatype.findMany({where:{}});
</script>

<script>
	import { getContext } from 'svelte';
	import { key } from '../../api/api.js';

	export let value;
	export let op;

	load.then((dts) => {
		value = dts[0];
	});

	import HLSelect from '../../ui/form/HLSelect.svelte';
	import {restrictToIdent, cap, isObject} from '../../lib/util.js';

	let operation = getContext(key);

	$: {
		op = {
			connect: {id: value.id},
		}
	}

</script>

{#await load then datatypes}
<HLSelect bind:value={value}>
	{#each Object.entries(datatypes) as dt}
	<option value={dt[1]}>
		{cap(dt[1].name)}
	</option>
	{/each}
</HLSelect>
{/await}
