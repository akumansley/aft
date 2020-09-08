<script context="module">
	import client from '../../data/client.js';
	let load = client.api.model.findMany({where:{}});
</script>

<script>
	import { getContext } from 'svelte';
	import { key } from '../../api/api.js';

	export let value;
	let model = {};
	load.then((models) => {
		model = models[0];
	});

	import HLSelect from '../../ui/form/HLSelect.svelte';
	import {restrictToIdent, cap, isObject} from '../util.js';

	let operation = getContext(key);

	$: {
		value = {
			connect: {id: model.id},
		}
	}

</script>

{#await load then models}
<HLSelect bind:value={model}>
	{#each Object.entries(models) as m}
	<option value={m[1]}>
		{cap(m[1].name)}
	</option>
	{/each}
</HLSelect>
{/await}
